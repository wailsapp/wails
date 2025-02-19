package application

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/internal/fileexplorer"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"

	"github.com/pkg/browser"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/signal"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/capabilities"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets/*
var alphaAssets embed.FS

var globalApplication *App

// AlphaAssets is the default assets for the alpha application
var AlphaAssets = AssetOptions{
	Handler: BundledAssetFileServer(alphaAssets),
}

func init() {
	runtime.LockOSThread()
}

type EventListener struct {
	callback func(app *ApplicationEvent)
}

func Get() *App {
	return globalApplication
}

func New(appOptions Options) *App {
	if globalApplication != nil {
		return globalApplication
	}

	mergeApplicationDefaults(&appOptions)

	result := newApplication(appOptions)
	globalApplication = result
	fatalHandler(result.handleFatalError)

	if result.Logger == nil {
		if result.isDebugMode {
			result.Logger = DefaultLogger(result.options.LogLevel)
		} else {
			result.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
		}
	}

	if !appOptions.DisableDefaultSignalHandler {
		result.signalHandler = signal.NewSignalHandler(result.Quit)
		result.signalHandler.Logger = result.Logger
		result.signalHandler.ExitMessage = func(sig os.Signal) string {
			return "Quitting application..."
		}
	}

	result.logStartup()
	result.logPlatformInfo()

	result.customEventProcessor = NewWailsEventProcessor(result.dispatchEventToListeners)

	messageProc := NewMessageProcessor(result.Logger)
	opts := &assetserver.Options{
		Handler: appOptions.Assets.Handler,
		Middleware: assetserver.ChainMiddleware(
			func(next http.Handler) http.Handler {
				if m := appOptions.Assets.Middleware; m != nil {
					return m(next)
				}
				return next
			},
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					path := req.URL.Path
					switch path {
					case "/wails/runtime":
						messageProc.ServeHTTP(rw, req)
					case "/wails/capabilities":
						err := assetserver.ServeFile(rw, path, globalApplication.capabilities.AsBytes())
						if err != nil {
							result.fatal("unable to serve capabilities: %w", err)
						}
					case "/wails/flags":
						updatedOptions := result.impl.GetFlags(appOptions)
						flags, err := json.Marshal(updatedOptions)
						if err != nil {
							result.fatal("invalid flags provided to application: %w", err)
						}
						err = assetserver.ServeFile(rw, path, flags)
						if err != nil {
							result.fatal("unable to serve flags: %w", err)
						}
					default:
						next.ServeHTTP(rw, req)
					}
				})
			},
		),
		Logger: result.Logger,
	}

	if appOptions.Assets.DisableLogging {
		opts.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	srv, err := assetserver.NewAssetServer(opts)
	if err != nil {
		result.fatal("application initialisation failed: %w", err)
	}

	result.assets = srv
	result.assets.LogDetails()

	result.bindings = NewBindings(appOptions.MarshalError, appOptions.BindAliases)
	result.options.Services = slices.Clone(appOptions.Services)

	// Process keybindings
	if result.options.KeyBindings != nil {
		result.keyBindings = processKeyBindingOptions(result.options.KeyBindings)
	}

	if appOptions.OnShutdown != nil {
		result.OnShutdown(appOptions.OnShutdown)
	}

	// Initialize single instance manager if enabled
	if appOptions.SingleInstance != nil {
		manager, err := newSingleInstanceManager(result, appOptions.SingleInstance)
		if err != nil {
			if errors.Is(err, alreadyRunningError) && manager != nil {
				err = manager.notifyFirstInstance()
				if err != nil {
					globalApplication.error("failed to notify first instance: %w", err)
				}
				os.Exit(appOptions.SingleInstance.ExitCode)
			}
			result.fatal("failed to initialize single instance manager: %w", err)
		} else {
			result.singleInstanceManager = manager
		}
	}

	return result
}

func mergeApplicationDefaults(o *Options) {
	if o.Name == "" {
		o.Name = "My Wails Application"
	}
	if o.Description == "" {
		o.Description = "An application written using Wails"
	}
	if o.Windows.WndClass == "" {
		o.Windows.WndClass = "WailsWebviewWindow"
	}
}

type (
	platformApp interface {
		run() error
		destroy()
		setApplicationMenu(menu *Menu)
		name() string
		getCurrentWindowID() uint
		showAboutDialog(name string, description string, icon []byte)
		setIcon(icon []byte)
		on(id uint)
		dispatchOnMainThread(id uint)
		hide()
		show()
		getPrimaryScreen() (*Screen, error)
		getScreens() ([]*Screen, error)
		GetFlags(options Options) map[string]any
		isOnMainThread() bool
		isDarkMode() bool
	}

	runnable interface {
		Run()
	}
)

// Messages sent from javascript get routed here
type windowMessage struct {
	windowId uint
	message  string
}

var windowMessageBuffer = make(chan *windowMessage, 5)

type dragAndDropMessage struct {
	windowId  uint
	filenames []string
}

var windowDragAndDropBuffer = make(chan *dragAndDropMessage, 5)

func addDragAndDropMessage(windowId uint, filenames []string) {
	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:  windowId,
		filenames: filenames,
	}
}

var _ webview.Request = &webViewAssetRequest{}

const webViewRequestHeaderWindowId = "x-wails-window-id"
const webViewRequestHeaderWindowName = "x-wails-window-name"

type webViewAssetRequest struct {
	webview.Request
	windowId   uint
	windowName string
}

var windowKeyEvents = make(chan *windowKeyEvent, 5)

type windowKeyEvent struct {
	windowId          uint
	acceleratorString string
}

func (r *webViewAssetRequest) Header() (http.Header, error) {
	h, err := r.Request.Header()
	if err != nil {
		return nil, err
	}

	hh := h.Clone()
	hh.Set(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(r.windowId), 10))
	return hh, nil
}

var webviewRequests = make(chan *webViewAssetRequest, 5)

type eventHook struct {
	callback func(event *ApplicationEvent)
}

type App struct {
	ctx                           context.Context
	cancel                        context.CancelFunc
	options                       Options
	applicationEventListeners     map[uint][]*EventListener
	applicationEventListenersLock sync.RWMutex
	applicationEventHooks         map[uint][]*eventHook
	applicationEventHooksLock     sync.RWMutex

	// Screens layout manager (handles DIP coordinate system)
	screenManager ScreenManager

	// Windows
	windows     map[uint]Window
	windowsLock sync.RWMutex

	// System Trays
	systemTrays      map[uint]*SystemTray
	systemTraysLock  sync.Mutex
	systemTrayID     uint
	systemTrayIDLock sync.RWMutex

	// MenuItems
	menuItems     map[uint]*MenuItem
	menuItemsLock sync.Mutex

	// Starting and running
	starting   bool
	running    bool
	runLock    sync.Mutex
	pendingRun []runnable

	bindings *Bindings

	// platform app
	impl platformApp

	// The main application menu
	ApplicationMenu *Menu

	clipboard            *Clipboard
	customEventProcessor *EventProcessor
	Logger               *slog.Logger

	contextMenus     map[string]*ContextMenu
	contextMenusLock sync.Mutex

	assets   *assetserver.AssetServer
	startURL string

	// Hooks
	windowCreatedCallbacks []func(window Window)
	pid                    int

	// Capabilities
	capabilities capabilities.Capabilities
	isDebugMode  bool

	// Keybindings
	keyBindings     map[string]func(window *WebviewWindow)
	keyBindingsLock sync.RWMutex

	// Shutdown
	performingShutdown  bool
	shutdownLock        sync.Mutex
	serviceShutdownLock sync.Mutex

	// Shutdown tasks are run when the application is shutting down.
	// They are run in the order they are added and run on the main thread.
	// The application option `OnShutdown` is run first.
	shutdownTasks []func()

	// signalHandler is used to handle signals
	signalHandler *signal.SignalHandler

	// Wails ApplicationEvent Listener related
	wailsEventListenerLock sync.Mutex
	wailsEventListeners    []WailsEventListener

	// singleInstanceManager handles single instance functionality
	singleInstanceManager *singleInstanceManager
}

func (a *App) handleWarning(msg string) {
	if a.options.WarningHandler != nil {
		a.options.WarningHandler(msg)
	} else {
		a.Logger.Warn(msg)
	}
}

func (a *App) handleError(err error) {
	if a.options.ErrorHandler != nil {
		a.options.ErrorHandler(err)
	} else {
		a.Logger.Error(err.Error())
	}
}

// RegisterService appends the given service to the list of bound services.
// Registered services will be bound and initialised
// in registration order upon calling [App.Run].
//
// RegisterService will log an error message
// and discard the given service
// if called after [App.Run].
func (a *App) RegisterService(service Service) {
	a.runLock.Lock()
	defer a.runLock.Unlock()

	if a.starting || a.running {
		a.error("services must be registered before running the application. Service '%s' will not be registered.", getServiceName(service))
		return
	}

	a.options.Services = append(a.options.Services, service)
}

// EmitEvent will emit an event
func (a *App) EmitEvent(name string, data ...any) {
	a.customEventProcessor.Emit(&CustomEvent{
		Name: name,
		Data: data,
	})
}

// EmitEvent will emit an event
func (a *App) emitEvent(event *CustomEvent) {
	a.customEventProcessor.Emit(event)
}

// OnEvent will listen for events
func (a *App) OnEvent(name string, callback func(event *CustomEvent)) func() {
	return a.customEventProcessor.On(name, callback)
}

// OffEvent will remove an event listener
func (a *App) OffEvent(name string) {
	a.customEventProcessor.Off(name)
}

// OnMultipleEvent will listen for events a set number of times before unsubscribing.
func (a *App) OnMultipleEvent(name string, callback func(event *CustomEvent), counter int) {
	a.customEventProcessor.OnMultiple(name, callback, counter)
}

// ResetEvents will remove all event listeners and hooks
func (a *App) ResetEvents() {
	a.customEventProcessor.OffAll()
}

func (a *App) handleFatalError(err error) {
	a.handleError(&FatalError{err: err})
	os.Exit(1)
}

func (a *App) init() {
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.applicationEventHooks = make(map[uint][]*eventHook)
	a.applicationEventListeners = make(map[uint][]*EventListener)
	a.windows = make(map[uint]Window)
	a.systemTrays = make(map[uint]*SystemTray)
	a.contextMenus = make(map[string]*ContextMenu)
	a.keyBindings = make(map[string]func(window *WebviewWindow))
	a.Logger = a.options.Logger
	a.pid = os.Getpid()
	a.wailsEventListeners = make([]WailsEventListener, 0)
}

func (a *App) getSystemTrayID() uint {
	a.systemTrayIDLock.Lock()
	defer a.systemTrayIDLock.Unlock()
	a.systemTrayID++
	return a.systemTrayID
}

func (a *App) getWindowForID(id uint) Window {
	a.windowsLock.RLock()
	defer a.windowsLock.RUnlock()
	return a.windows[id]
}

func (a *App) deleteWindowByID(id uint) {
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	delete(a.windows, id)
}

func (a *App) Capabilities() capabilities.Capabilities {
	return a.capabilities
}

func (a *App) OnApplicationEvent(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func() {
	eventID := uint(eventType)
	a.applicationEventListenersLock.Lock()
	defer a.applicationEventListenersLock.Unlock()
	listener := &EventListener{
		callback: callback,
	}
	a.applicationEventListeners[eventID] = append(a.applicationEventListeners[eventID], listener)
	if a.impl != nil {
		go func() {
			defer handlePanic()
			a.impl.on(eventID)
		}()
	}

	return func() {
		// lock the map
		a.applicationEventListenersLock.Lock()
		defer a.applicationEventListenersLock.Unlock()
		// Remove listener
		a.applicationEventListeners[eventID] = lo.Without(a.applicationEventListeners[eventID], listener)
	}
}

// RegisterApplicationEventHook registers a hook for the given application event.
// Hooks are called before the event listeners and can cancel the event.
// The returned function can be called to remove the hook.
func (a *App) RegisterApplicationEventHook(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func() {
	eventID := uint(eventType)
	a.applicationEventHooksLock.Lock()
	defer a.applicationEventHooksLock.Unlock()
	thisHook := &eventHook{
		callback: callback,
	}
	a.applicationEventHooks[eventID] = append(a.applicationEventHooks[eventID], thisHook)

	return func() {
		a.applicationEventHooksLock.Lock()
		a.applicationEventHooks[eventID] = lo.Without(a.applicationEventHooks[eventID], thisHook)
		a.applicationEventHooksLock.Unlock()
	}
}

//func (a *App) RegisterListener(listener WailsEventListener) {
//	a.wailsEventListenerLock.Lock()
//	a.wailsEventListeners = append(a.wailsEventListeners, listener)
//	a.wailsEventListenerLock.Unlock()
//}
//
//func (a *App) RegisterServiceHandler(prefix string, handler http.Handler) {
//	a.assets.AttachServiceHandler(prefix, handler)
//}

func (a *App) NewWebviewWindow() *WebviewWindow {
	return a.NewWebviewWindowWithOptions(WebviewWindowOptions{})
}

func (a *App) GetPID() int {
	return a.pid
}

func (a *App) info(message string, args ...any) {
	if a.Logger != nil {
		go func() {
			defer handlePanic()
			a.Logger.Info(message, args...)
		}()
	}
}

func (a *App) debug(message string, args ...any) {
	if a.Logger != nil {
		go func() {
			defer handlePanic()
			a.Logger.Debug(message, args...)
		}()
	}
}

func (a *App) fatal(message string, args ...any) {
	err := fmt.Errorf(message, args...)
	a.handleFatalError(err)
}
func (a *App) warning(message string, args ...any) {
	msg := fmt.Sprintf(message, args...)
	a.handleWarning(msg)
}

func (a *App) error(message string, args ...any) {
	a.handleError(fmt.Errorf(message, args...))
}

func (a *App) NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow {
	newWindow := NewWindow(windowOptions)
	id := newWindow.ID()

	a.windowsLock.Lock()
	a.windows[id] = newWindow
	a.windowsLock.Unlock()

	// Call hooks
	for _, hook := range a.windowCreatedCallbacks {
		hook(newWindow)
	}

	a.runOrDeferToAppRun(newWindow)

	return newWindow
}

func (a *App) NewSystemTray() *SystemTray {
	id := a.getSystemTrayID()
	newSystemTray := newSystemTray(id)

	a.systemTraysLock.Lock()
	a.systemTrays[id] = newSystemTray
	a.systemTraysLock.Unlock()

	a.runOrDeferToAppRun(newSystemTray)

	return newSystemTray
}

func (a *App) Run() error {
	a.runLock.Lock()
	// Prevent double invocations.
	if a.starting || a.running {
		a.runLock.Unlock()
		return errors.New("application is running or a previous run has failed")
	}
	// Block further service registrations.
	a.starting = true
	a.runLock.Unlock()

	// Ensure application context is canceled in case of failures.
	defer a.cancel()

	// Call post-create hooks
	err := a.preRun()
	if err != nil {
		return err
	}

	a.impl = newPlatformApp(a)

	// Ensure services are shut down in case of failures.
	defer a.shutdownServices()
	// Ensure application context is canceled before service shutdown (duplicate calls don't hurt).
	defer a.cancel()

	// Startup services before dispatching any events.
	// No need to hold the lock here because a.options.Services may only change when a.running is false.
	services := a.options.Services
	a.options.Services = nil
	for i, service := range services {
		if err := a.startupService(service); err != nil {
			return fmt.Errorf("error starting service '%s': %w", getServiceName(service), err)
		}
		// Schedule started services for shutdown.
		a.options.Services = services[:i+1]
	}

	go func() {
		for {
			event := <-applicationEvents
			go a.handleApplicationEvent(event)
		}
	}()
	go func() {
		for {
			event := <-windowEvents
			go a.handleWindowEvent(event)
		}
	}()
	go func() {
		for {
			request := <-webviewRequests
			go a.handleWebViewRequest(request)
		}
	}()
	go func() {
		for {
			event := <-windowMessageBuffer
			go a.handleWindowMessage(event)
		}
	}()
	go func() {
		for {
			event := <-windowKeyEvents
			go a.handleWindowKeyEvent(event)
		}
	}()
	go func() {
		for {
			dragAndDropMessage := <-windowDragAndDropBuffer
			go a.handleDragAndDropMessage(dragAndDropMessage)
		}
	}()

	go func() {
		for {
			menuItemID := <-menuItemClicked
			go a.handleMenuItemClicked(menuItemID)
		}
	}()

	a.runLock.Lock()
	a.running = true
	a.runLock.Unlock()

	// No need to hold the lock here because
	//   - a.pendingRun may only change while a.running is false.
	//   - runnables are scheduled asynchronously anyway.
	for _, pending := range a.pendingRun {
		go func() {
			defer handlePanic()
			pending.Run()
		}()
	}
	a.pendingRun = nil

	// set the application menu
	if runtime.GOOS == "darwin" {
		a.impl.setApplicationMenu(a.ApplicationMenu)
	}
	if a.options.Icon != nil {
		a.impl.setIcon(a.options.Icon)
	}

	return a.impl.run()
}

func (a *App) startupService(service Service) error {
	err := a.bindings.Add(service)
	if err != nil {
		return fmt.Errorf("cannot bind service methods: %w", err)
	}

	if service.options.Route != "" {
		handler, ok := service.Instance().(http.Handler)
		if !ok {
			handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				http.Error(
					rw,
					fmt.Sprintf("Service '%s' does not handle HTTP requests", getServiceName(service)),
					http.StatusServiceUnavailable,
				)
			})
		}
		a.assets.AttachServiceHandler(service.options.Route, handler)
	}

	if s, ok := service.instance.(ServiceStartup); ok {
		a.debug("Starting up service:", "name", getServiceName(service))
		return s.ServiceStartup(a.ctx, service.options)
	}

	return nil
}

func (a *App) shutdownServices() {
	// Acquire lock to prevent double calls (defer in Run() + OnShutdown)
	a.serviceShutdownLock.Lock()
	defer a.serviceShutdownLock.Unlock()

	// Ensure app context is canceled first (duplicate calls don't hurt).
	a.cancel()

	for len(a.options.Services) > 0 {
		last := len(a.options.Services) - 1
		service := a.options.Services[last]
		a.options.Services = a.options.Services[:last] // Prevent double shutdowns

		if s, ok := service.instance.(ServiceShutdown); ok {
			a.debug("Shutting down service:", "name", getServiceName(service))
			if err := s.ServiceShutdown(); err != nil {
				a.error("error shutting down service '%s': %w", getServiceName(service), err)
			}
		}
	}
}

func (a *App) handleApplicationEvent(event *ApplicationEvent) {
	defer handlePanic()
	a.applicationEventListenersLock.RLock()
	listeners, ok := a.applicationEventListeners[event.Id]
	a.applicationEventListenersLock.RUnlock()
	if !ok {
		return
	}

	// Process Hooks
	a.applicationEventHooksLock.RLock()
	hooks, ok := a.applicationEventHooks[event.Id]
	a.applicationEventHooksLock.RUnlock()
	if ok {
		for _, thisHook := range hooks {
			thisHook.callback(event)
			if event.IsCancelled() {
				return
			}
		}
	}

	for _, listener := range listeners {
		go func() {
			if event.IsCancelled() {
				return
			}
			defer handlePanic()
			listener.callback(event)
		}()
	}
}

func (a *App) handleDragAndDropMessage(event *dragAndDropMessage) {
	defer handlePanic()
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.HandleDragAndDropMessage(event.filenames)
}

func (a *App) handleWindowMessage(event *windowMessage) {
	defer handlePanic()
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.RUnlock()
	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Check if the message starts with "wails:"
	if strings.HasPrefix(event.message, "wails:") {
		window.HandleMessage(event.message)
	} else {
		if a.options.RawMessageHandler != nil {
			a.options.RawMessageHandler(window, event.message)
		}
	}
}

func (a *App) handleWebViewRequest(request *webViewAssetRequest) {
	defer handlePanic()
	a.assets.ServeWebViewRequest(request)
}

func (a *App) handleWindowEvent(event *windowEvent) {
	defer handlePanic()
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.WindowID]
	a.windowsLock.RUnlock()
	if !ok {
		a.warning("Window #%d not found", event.WindowID)
		return
	}
	window.HandleWindowEvent(event.EventID)
}

func (a *App) handleMenuItemClicked(menuItemID uint) {
	defer handlePanic()

	menuItem := getMenuItemByID(menuItemID)
	if menuItem == nil {
		a.warning("MenuItem #%d not found", menuItemID)
		return
	}
	menuItem.handleClick()
}

func (a *App) CurrentWindow() *WebviewWindow {
	if a.impl == nil {
		return nil
	}
	id := a.impl.getCurrentWindowID()
	a.windowsLock.RLock()
	defer a.windowsLock.RUnlock()
	result := a.windows[id]
	if result == nil {
		return nil
	}
	return result.(*WebviewWindow)
}

// OnShutdown adds a function to be run when the application is shutting down.
func (a *App) OnShutdown(f func()) {
	if f == nil {
		return
	}

	a.shutdownLock.Lock()

	if !a.performingShutdown {
		defer a.shutdownLock.Unlock()
		a.shutdownTasks = append(a.shutdownTasks, f)
		return
	}

	a.shutdownLock.Unlock()
	InvokeAsync(f)
}

func (a *App) destroySystemTray(tray *SystemTray) {
	// Remove the system tray from the a.systemTrays map
	a.systemTraysLock.Lock()
	delete(a.systemTrays, tray.id)
	a.systemTraysLock.Unlock()
	tray.destroy()
}

func (a *App) cleanup() {
	a.shutdownLock.Lock()
	if a.performingShutdown {
		a.shutdownLock.Unlock()
		return
	}
	a.cancel() // Cancel app context before running shutdown hooks.
	a.performingShutdown = true
	a.shutdownLock.Unlock()

	// No need to hold the lock here because a.shutdownTasks
	// may only change while a.performingShutdown is false.
	for _, shutdownTask := range a.shutdownTasks {
		InvokeSync(shutdownTask)
	}
	InvokeSync(func() {
		a.shutdownServices()
		a.windowsLock.RLock()
		for _, window := range a.windows {
			window.Close()
		}
		a.windows = nil
		a.windowsLock.RUnlock()
		a.systemTraysLock.Lock()
		for _, systray := range a.systemTrays {
			systray.destroy()
		}
		a.systemTrays = nil
		a.systemTraysLock.Unlock()

		// Cleanup single instance manager
		if a.singleInstanceManager != nil {
			a.singleInstanceManager.cleanup()
		}

		a.postQuit()

		if a.options.PostShutdown != nil {
			a.options.PostShutdown()
		}
	})
}

func (a *App) Quit() {
	if a.impl != nil {
		InvokeSync(a.impl.destroy)
	}
}

func (a *App) SetIcon(icon []byte) {
	if a.impl != nil {
		a.impl.setIcon(icon)
	}
}

func (a *App) SetMenu(menu *Menu) {
	a.ApplicationMenu = menu
	if a.impl != nil {
		a.impl.setApplicationMenu(menu)
	}
}
func (a *App) ShowAboutDialog() {
	if a.impl != nil {
		a.impl.showAboutDialog(a.options.Name, a.options.Description, a.options.Icon)
	}
}

func InfoDialog() *MessageDialog {
	return newMessageDialog(InfoDialogType)
}

func QuestionDialog() *MessageDialog {
	return newMessageDialog(QuestionDialogType)
}

func WarningDialog() *MessageDialog {
	return newMessageDialog(WarningDialogType)
}

func ErrorDialog() *MessageDialog {
	return newMessageDialog(ErrorDialogType)
}

func OpenFileDialog() *OpenFileDialogStruct {
	return newOpenFileDialog()
}

func SaveFileDialog() *SaveFileDialogStruct {
	return newSaveFileDialog()
}

// NOTE: should use screenManager directly after DPI is implemented in all platforms
// (should also get rid of the error return)
func (a *App) GetScreens() ([]*Screen, error) {
	return a.impl.getScreens()
	// return a.screenManager.screens, nil
}

// NOTE: should use screenManager directly after DPI is implemented in all platforms
// (should also get rid of the error return)
func (a *App) GetPrimaryScreen() (*Screen, error) {
	return a.impl.getPrimaryScreen()
	// return a.screenManager.primaryScreen, nil
}

func (a *App) Clipboard() *Clipboard {
	if a.clipboard == nil {
		a.clipboard = newClipboard()
	}
	return a.clipboard
}

func (a *App) dispatchOnMainThread(fn func()) {
	// If we are on the main thread, just call the function
	if a.impl.isOnMainThread() {
		fn()
		return
	}

	mainThreadFunctionStoreLock.Lock()
	id := generateFunctionStoreID()
	mainThreadFunctionStore[id] = fn
	mainThreadFunctionStoreLock.Unlock()
	// Call platform specific dispatch function
	a.impl.dispatchOnMainThread(id)
}

func OpenFileDialogWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct {
	result := OpenFileDialog()
	result.SetOptions(options)
	return result
}

func SaveFileDialogWithOptions(s *SaveFileDialogOptions) *SaveFileDialogStruct {
	result := SaveFileDialog()
	result.SetOptions(s)
	return result
}

func (a *App) dispatchEventToListeners(event *CustomEvent) {
	listeners := a.wailsEventListeners

	for _, window := range a.windows {
		if event.IsCancelled() {
			return
		}
		window.DispatchWailsEvent(event)
	}

	for _, listener := range listeners {
		if event.IsCancelled() {
			return
		}
		listener.DispatchWailsEvent(event)
	}
}

func (a *App) IsDarkMode() bool {
	if a.impl == nil {
		return false
	}
	return a.impl.isDarkMode()
}

func (a *App) Hide() {
	if a.impl != nil {
		a.impl.hide()
	}
}

func (a *App) Show() {
	if a.impl != nil {
		a.impl.show()
	}
}

func (a *App) registerContextMenu(menu *ContextMenu) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	a.contextMenus[menu.name] = menu
}

func (a *App) unregisterContextMenu(name string) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	delete(a.contextMenus, name)
}

func (a *App) getContextMenu(name string) (*ContextMenu, bool) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	menu, ok := a.contextMenus[name]
	return menu, ok

}

func (a *App) OnWindowCreation(callback func(window Window)) {
	a.windowCreatedCallbacks = append(a.windowCreatedCallbacks, callback)
}

func (a *App) GetWindowByName(name string) Window {
	a.windowsLock.RLock()
	defer a.windowsLock.RUnlock()
	for _, window := range a.windows {
		if window.Name() == name {
			return window
		}
	}
	return nil
}

func (a *App) runOrDeferToAppRun(r runnable) {
	a.runLock.Lock()

	if !a.running {
		defer a.runLock.Unlock() // Defer unlocking for panic tolerance.
		a.pendingRun = append(a.pendingRun, r)
		return
	}

	// Unlock immediately to prevent deadlocks.
	// No TOC/TOU risk here because a.running can never switch back to false.
	a.runLock.Unlock()
	r.Run()
}

func (a *App) processKeyBinding(acceleratorString string, window *WebviewWindow) bool {
	if len(a.keyBindings) == 0 {
		return false
	}

	a.keyBindingsLock.RLock()
	defer a.keyBindingsLock.RUnlock()

	// Check key bindings
	callback, ok := a.keyBindings[acceleratorString]
	if !ok {
		return false
	}

	// Execute callback
	go callback(window)

	return true
}

func (a *App) addKeyBinding(acceleratorString string, callback func(window *WebviewWindow)) {
	a.keyBindingsLock.Lock()
	defer a.keyBindingsLock.Unlock()
	a.keyBindings[acceleratorString] = callback
}

func (a *App) removeKeyBinding(acceleratorString string) {
	a.keyBindingsLock.Lock()
	defer a.keyBindingsLock.Unlock()
	delete(a.keyBindings, acceleratorString)
}

func (a *App) handleWindowKeyEvent(event *windowKeyEvent) {
	defer handlePanic()
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.RUnlock()
	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.HandleKeyEvent(event.acceleratorString)
}

func (a *App) AssetServerHandler() func(rw http.ResponseWriter, req *http.Request) {
	return a.assets.ServeHTTP
}

func (a *App) BrowserOpenURL(url string) error {
	return browser.OpenURL(url)
}

func (a *App) BrowserOpenFile(path string) error {
	return browser.OpenFile(path)
}

func (a *App) Environment() EnvironmentInfo {
	info, _ := operatingsystem.Info()
	result := EnvironmentInfo{
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		Debug:  a.isDebugMode,
		OSInfo: info,
	}
	result.PlatformInfo = a.platformEnvironment()
	return result
}

func (a *App) shouldQuit() bool {
	if a.options.ShouldQuit != nil {
		return a.options.ShouldQuit()
	}
	return true
}

// OpenFileManager opens the file manager at the specified path, optionally selecting the file.
func (a *App) OpenFileManager(path string, selectFile bool) error {
	return InvokeSyncWithError(func() error {
		return fileexplorer.OpenFileManager(path, selectFile)
	})
}
