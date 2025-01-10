package application

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
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
							result.handleFatalError(fmt.Errorf("unable to serve capabilities: %s", err.Error()))
						}
					case "/wails/flags":
						updatedOptions := result.impl.GetFlags(appOptions)
						flags, err := json.Marshal(updatedOptions)
						if err != nil {
							result.handleFatalError(fmt.Errorf("invalid flags provided to application: %s", err.Error()))
						}
						err = assetserver.ServeFile(rw, path, flags)
						if err != nil {
							result.handleFatalError(fmt.Errorf("unable to serve flags: %s", err.Error()))
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
		result.handleFatalError(fmt.Errorf("Fatal error in application initialisation: " + err.Error()))
	}

	result.assets = srv
	result.assets.LogDetails()

	result.bindings, err = NewBindings(appOptions.Services, appOptions.BindAliases)
	if err != nil {
		result.handleFatalError(fmt.Errorf("Fatal error in application initialisation: " + err.Error()))
	}

	for i, service := range appOptions.Services {
		if thisService, ok := service.instance.(ServiceStartup); ok {
			err := thisService.OnStartup(result.ctx, service.options)
			if err != nil {
				name := service.options.Name
				if name == "" {
					name = getServiceName(service.instance)
				}
				globalApplication.Logger.Error("OnStartup() failed shutting down application:", "service", name, "error", err.Error())
				// Run shutdown on all services that have already started
				for _, service := range appOptions.Services[:i] {
					if thisService, ok := service.instance.(ServiceShutdown); ok {
						err := thisService.OnShutdown()
						if err != nil {
							globalApplication.Logger.Error("Error shutting down service: " + err.Error())
						}
					}
				}
				// Shutdown the application
				os.Exit(1)
			}
		}
	}

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
					globalApplication.error("Failed to notify first instance: " + err.Error())
				}
				os.Exit(appOptions.SingleInstance.ExitCode)
			}
			result.handleFatalError(fmt.Errorf("failed to initialize single instance manager: %w", err))
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

func processPanicHandlerRecover() {
	h := globalApplication.options.PanicHandler
	if h == nil {
		return
	}

	if err := recover(); err != nil {
		h(err)
	}
}

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

	// Running
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

	contextMenus     map[string]*Menu
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
	performingShutdown bool

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
	var buffer strings.Builder
	buffer.WriteString("*********************** FATAL ***********************")
	buffer.WriteString("There has been a catastrophic failure in your application.")
	buffer.WriteString("Please report this error at https://github.com/wailsapp/wails/issues")
	buffer.WriteString("******************** Error Details ******************")
	buffer.WriteString(fmt.Sprintf("Message: " + err.Error()))
	buffer.WriteString(fmt.Sprintf("Stack: " + string(debug.Stack())))
	buffer.WriteString("*********************** FATAL ***********************")
	a.handleError(fmt.Errorf(buffer.String()))
	os.Exit(1)
}

func (a *App) init() {
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.applicationEventHooks = make(map[uint][]*eventHook)
	a.applicationEventListeners = make(map[uint][]*EventListener)
	a.windows = make(map[uint]Window)
	a.systemTrays = make(map[uint]*SystemTray)
	a.contextMenus = make(map[string]*Menu)
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
		go a.impl.on(eventID)
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
		go a.Logger.Info(message, args...)
	}
}

func (a *App) debug(message string, args ...any) {
	if a.Logger != nil {
		go a.Logger.Debug(message, args...)
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

	// Setup panic handler
	defer processPanicHandlerRecover()

	// Call post-create hooks
	err := a.preRun()
	if err != nil {
		return err
	}

	a.impl = newPlatformApp(a)
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

	for _, systray := range a.pendingRun {
		go systray.Run()
	}
	a.pendingRun = nil

	a.runLock.Unlock()

	// set the application menu
	if runtime.GOOS == "darwin" {
		a.impl.setApplicationMenu(a.ApplicationMenu)
	}
	if a.options.Icon != nil {
		a.impl.setIcon(a.options.Icon)
	}

	err = a.impl.run()
	if err != nil {
		return err
	}

	// Cancel the context
	a.cancel()

	for _, service := range a.options.Services {
		// If it conforms to the ServiceShutdown interface, call the Shutdown method
		if thisService, ok := service.instance.(ServiceShutdown); ok {
			err := thisService.OnShutdown()
			if err != nil {
				a.error("Error shutting down service: " + err.Error())
			}
		}
	}

	return nil
}

func (a *App) handleApplicationEvent(event *ApplicationEvent) {
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
			if event.Cancelled {
				return
			}
		}
	}

	for _, listener := range listeners {
		go listener.callback(event)
	}
}

func (a *App) handleDragAndDropMessage(event *dragAndDropMessage) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.HandleDragAndDropMessage(event.filenames)
}

func (a *App) handleWindowMessage(event *windowMessage) {
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.RUnlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
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
	a.assets.ServeWebViewRequest(request)
}

func (a *App) handleWindowEvent(event *windowEvent) {
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.WindowID]
	a.windowsLock.RUnlock()
	if !ok {
		log.Printf("Window #%d not found", event.WindowID)
		return
	}
	window.HandleWindowEvent(event.EventID)
}

func (a *App) handleMenuItemClicked(menuItemID uint) {
	menuItem := getMenuItemByID(menuItemID)
	if menuItem == nil {
		log.Printf("MenuItem #%d not found", menuItemID)
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
	a.shutdownTasks = append(a.shutdownTasks, f)
}

func (a *App) cleanup() {
	if a.performingShutdown {
		return
	}
	a.performingShutdown = true
	for _, shutdownTask := range a.shutdownTasks {
		InvokeSync(shutdownTask)
	}
	InvokeSync(func() {
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
	})
	// Cleanup single instance manager
	if a.singleInstanceManager != nil {
		a.singleInstanceManager.cleanup()
	}
}

func (a *App) Quit() {
	if a.impl != nil {
		InvokeSync(a.impl.destroy)
		a.postQuit()
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
		window.DispatchWailsEvent(event)
	}

	for _, listener := range listeners {
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

func (a *App) RegisterContextMenu(name string, menu *Menu) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	a.contextMenus[name] = menu
}

func (a *App) getContextMenu(name string) (*Menu, bool) {
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
	running := a.running
	if !running {
		a.pendingRun = append(a.pendingRun, r)
	}
	a.runLock.Unlock()

	if running {
		r.Run()
	}
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
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.RUnlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
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

// Path returns the path for the given selector
func (a *App) Path(selector Path) string {
	return paths[selector]
}

// Paths returns the paths for the given selector
func (a *App) Paths(selector Paths) []string {
	return pathdirs[selector]
}

// OpenFileManager opens the file manager at the specified path, optionally selecting the file.
func (a *App) OpenFileManager(path string, selectFile bool) error {
	return InvokeSyncWithError(func() error {
		return fileexplorer.OpenFileManager(path, selectFile)
	})
}
