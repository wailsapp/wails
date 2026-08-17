package application

import (
	"context"
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

	"github.com/wailsapp/wails/v3/internal/capabilities"
)

var globalApplication *App

type EventListener struct {
	callback func(app *ApplicationEvent)
}

func Get() *App {
	return globalApplication
}

func New(appOptions Options) *App {
	// If we were spawned as an updater helper the process must perform the
	// swap and exit before any application machinery touches the disk. This
	// is a no-op when the sentinel env vars are absent, so normal startup is
	// unaffected.
	prepareUpdaterApplicationProcess()

	if globalApplication != nil {
		return globalApplication
	}
	if nativeBuild {
		// A native build has no frontend transport or WebView backend. Make the
		// matching runtime mode automatic so tagged applications cannot
		// accidentally start infrastructure that was compiled out.
		appOptions.NativeOnly = true
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

	// Set up signal handling (platform-specific)
	result.setupSignalHandler(appOptions)

	result.logStartup()
	result.logPlatformInfo()

	result.customEventProcessor = NewWailsEventProcessor(result.Event.dispatch)
	result.eventPayloads = newEventPayloadStore()
	// The store owns a reaper goroutine; window close only drops that window's
	// entries, so the store itself has to be shut down with the app.
	result.OnShutdown(result.eventPayloads.close)

	// Streams own session state and handler goroutines, so like the payload
	// store they are shut down with the app rather than per window.
	result.streams = newStreamManager(result)
	result.OnShutdown(result.streams.close)

	configureFrontend(result, appOptions)
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
		getAccentColor() string
	}

	runnable interface {
		Run()
	}
)

// Messages sent from javascript get routed here
type windowMessage struct {
	windowId   uint
	message    string
	originInfo *OriginInfo
}

type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}

var windowMessageBuffer = make(chan *windowMessage, 64)

// DropTargetDetails contains information about the HTML element
// where files were dropped (the element with data-file-drop-target attribute).
type DropTargetDetails struct {
	X          int               `json:"x"`
	Y          int               `json:"y"`
	ElementID  string            `json:"id"`
	ClassList  []string          `json:"classList"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type dragAndDropMessage struct {
	windowId   uint
	filenames  []string
	X          int
	Y          int
	DropTarget *DropTargetDetails
}

var windowDragAndDropBuffer = make(chan *dragAndDropMessage, 5)

func addDragAndDropMessage(windowId uint, filenames []string, dropTarget *DropTargetDetails) {
	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:   windowId,
		filenames:  filenames,
		DropTarget: dropTarget,
	}
}

// serveEventPayload delivers an oversized event body that was parked by
// DispatchWailsEvent. Payloads are one-shot and bound to the window they were
// dispatched to, so a stale or cross-window id simply 404s.
func (a *App) serveEventPayload(rw http.ResponseWriter, req *http.Request) {
	if a.eventPayloads == nil {
		http.NotFound(rw, req)
		return
	}

	// Read-only endpoint; anything else is not something we serve.
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		rw.Header().Set("Allow", "GET, HEAD")
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ids are always 32 hex chars. Checking the shape first keeps a stream of
	// junk requests from doing map work on arbitrarily long keys.
	id := strings.TrimPrefix(req.URL.Path, eventPayloadPath)
	if len(id) != eventPayloadIDLen || !isHexString(id) {
		http.NotFound(rw, req)
		return
	}

	// Bind to the requesting window where the platform tags the request.
	// Parsed at uint width so the conversion cannot truncate on 32-bit builds.
	var windowID uint
	if raw := req.Header.Get(webViewRequestHeaderWindowId); raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, strconv.IntSize); err == nil {
			windowID = uint(parsed)
		}
	}

	data, ok := a.eventPayloads.take(id, windowID)
	if !ok {
		http.NotFound(rw, req)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Content-Length", strconv.Itoa(len(data)))
	_, _ = rw.Write(data)
}

var windowKeyEvents = make(chan *windowKeyEvent, 5)

type windowKeyEvent struct {
	windowId          uint
	acceleratorString string
}

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

	// Manager pattern for organized API
	Window *WindowManager
	// NativeWindow creates windows whose content is composed entirely from
	// platform-native views. This API is experimental in v3 and currently has
	// an AppKit implementation on macOS.
	NativeWindow   *NativeWindowManager
	ContextMenu    *ContextMenuManager
	KeyBinding     *KeyBindingManager
	Browser        *BrowserManager
	Env            *EnvironmentManager
	Dialog         *DialogManager
	Event          *EventManager
	Menu           *MenuManager
	Screen         *ScreenManager
	Clipboard      *ClipboardManager
	SystemTray     *SystemTrayManager
	Autostart      *AutostartManager
	GlobalShortcut *GlobalShortcutManager
	Updater        *applicationUpdater

	// Windows
	windows     map[uint]Window
	windowsLock sync.RWMutex

	// Native windows deliberately live outside the legacy Window registry:
	// Window currently includes WebView-only methods such as ExecJS and SetURL.
	nativeWindows     map[uint]*NativeWindow
	nativeWindowsLock sync.RWMutex

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

	// The main application menu (private - use app.Menu.GetApplicationMenu/SetApplicationMenu)
	applicationMenu *Menu

	clipboard            *Clipboard
	customEventProcessor *EventProcessor
	Logger               *slog.Logger

	contextMenus     map[string]*ContextMenu
	contextMenusLock sync.RWMutex

	// eventPayloads holds oversized Go→JS event bodies awaiting a one-shot
	// fetch from the webview, keeping them out of evaluateJavaScript source.
	eventPayloads *eventPayloadStore

	// streams holds registered stream handlers and the per-page sessions that
	// carry their connections. See stream.go.
	streams *streamManager

	frontendState
	startURL string

	// Hooks
	windowCreatedCallbacks []func(window Window)
	pid                    int

	// Capabilities
	capabilities capabilities.Capabilities
	isDebugMode  bool

	// Keybindings
	keyBindings     map[string]func(window Window)
	keyBindingsLock sync.RWMutex

	// Shutdown
	performingShutdown  bool
	shutdownLock        sync.Mutex
	serviceShutdownLock sync.Mutex

	// Shutdown tasks are run when the application is shutting down.
	// They are run in the order they are added and run on the main thread.
	// The application option `OnShutdown` is run first.
	shutdownTasks []func()

	// Platform-specific fields (includes signal handler on desktop)
	platformSignalHandler

	// Wails ApplicationEvent Listener related
	wailsEventListenerLock sync.Mutex
	wailsEventListeners    []WailsEventListener

	// singleInstanceManager handles single instance functionality
	singleInstanceManager *singleInstanceManager

	// messageProcessor handles runtime messages
	messageProcessor *MessageProcessor
}

func (a *App) Config() Options {
	return a.options
}

// Context returns the application context that is canceled when the application shuts down.
// This context should be used for graceful shutdown of goroutines and long-running operations.
func (a *App) Context() context.Context {
	return a.ctx
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
		a.error(
			"services must be registered before running the application. Service '%s' will not be registered.",
			getServiceName(service),
		)
		return
	}

	a.options.Services = append(a.options.Services, service)
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
	a.nativeWindows = make(map[uint]*NativeWindow)
	a.systemTrays = make(map[uint]*SystemTray)
	a.contextMenus = make(map[string]*ContextMenu)
	a.keyBindings = make(map[string]func(window Window))
	a.Logger = a.options.Logger
	a.pid = os.Getpid()
	a.wailsEventListeners = make([]WailsEventListener, 0)

	// Initialize managers
	a.Window = newWindowManager(a)
	a.NativeWindow = newNativeWindowManager(a)
	a.ContextMenu = newContextMenuManager(a)
	a.KeyBinding = newKeyBindingManager(a)
	a.Browser = newBrowserManager(a)
	a.Env = newEnvironmentManager(a)
	a.Dialog = newDialogManager(a)
	a.Event = newEventManager(a)
	a.Menu = newMenuManager(a)
	a.Screen = newScreenManager(a)
	a.Clipboard = newClipboardManager(a)
	a.SystemTray = newSystemTrayManager(a)
	a.Autostart = newAutostartManager(a)
	a.GlobalShortcut = newGlobalShortcutManager(a)
	a.Updater = newApplicationUpdater(a)
}

func (a *App) Capabilities() capabilities.Capabilities {
	return a.capabilities
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

	// Ensure application context is cancelled in case of failures.
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

	// startup performs the remaining startup sequence: start services, spawn the
	// event-handling reader goroutines, run any pending windows, and apply the
	// menu/icon. On desktop this runs inline on the main goroutine. On iOS it is
	// deferred to a background goroutine (see below).
	startup := func() error {
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

		// Start the MCP server when the application is built with -tags mcp.
		// All configuration is read from environment variables (WAILS_MCP_HOST,
		// WAILS_MCP_PORT, WAILS_MCP_TIMEOUT, WAILS_MCP_HIDE_CURSOR).
		if err := startMCPServer(a); err != nil {
			return fmt.Errorf("mcp: %w", err)
		}

		go func() {
			for {
				event := <-applicationEvents
				go a.Event.handleApplicationEvent(event)
			}
		}()
		startFrontendEventLoops(a)

		go func() {
			for {
				menuItemID := <-menuItemClicked
				go a.Menu.handleMenuItemClicked(menuItemID)
			}
		}()

		go func() {
			for {
				itemID := <-toolbarItemClicked
				go handleToolbarItemClicked(itemID)
			}
		}()
		go func() {
			for {
				event := <-toolbarSearchTriggered
				go handleToolbarSearch(event.itemID, event.query)
			}
		}()
		go func() {
			for {
				event := <-toolbarShareCompleted
				go handleToolbarShareResult(event)
			}
		}()
		go func() {
			for {
				event := <-splitPaneCollapseEvents
				handleMacSplitPaneCollapsed(event.paneID, event.collapsed)
			}
		}()
		go func() {
			for {
				itemID := <-macSidebarItemSelected
				go handleMacSidebarItemSelected(itemID)
			}
		}()
		go func() {
			for {
				event := <-macInspectorControlEvents
				go handleMacInspectorControlEvent(event)
			}
		}()
		go func() {
			for {
				editorID := <-macTextEditorChanged
				go handleMacTextEditorChanged(editorID)
			}
		}()
		go func() {
			for {
				windowID := <-nativeWindowClosed
				if window, ok := a.NativeWindow.GetByID(windowID); ok {
					go window.Close()
				}
			}
		}()

		a.runLock.Lock()
		a.running = true
		a.runLock.Unlock()

		// Bind any global shortcuts that were registered before the app started.
		a.GlobalShortcut.flushPending()

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
			a.impl.setApplicationMenu(a.applicationMenu)
		}
		if a.options.Icon != nil {
			a.impl.setIcon(a.options.Icon)
		}
		return nil
	}

	if err := startup(); err != nil {
		return err
	}
	return a.impl.run()
}

func (a *App) startupService(service Service) error {
	if err := bindFrontendService(a, service); err != nil {
		return fmt.Errorf("cannot bind service methods: %w", err)
	}

	if err := attachServiceRoute(a, service); err != nil {
		return err
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

	// Ensure app context is cancelled first (duplicate calls don't hurt).
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

func (a *App) handleDragAndDropMessage(event *dragAndDropMessage) {
	defer handlePanic()
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	window.handleDragAndDropMessage(event.filenames, event.DropTarget)
}

func (a *App) handleWindowMessage(event *windowMessage) {
	defer handlePanic()
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.windowId]
	// Debug: log all window IDs
	var ids []uint
	for id := range a.windows {
		ids = append(ids, id)
	}
	a.windowsLock.RUnlock()

	a.debug("handleWindowMessage: Looking for window", "windowId", event.windowId, "availableIDs", ids)

	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Check if the message starts with "wails:"
	if strings.HasPrefix(event.message, "wails:") {
		a.debug("handleWindowMessage: Processing wails message", "message", event.message)
		window.HandleMessage(event.message)
	} else {
		if a.options.RawMessageHandler != nil {
			a.options.RawMessageHandler(window, event.message, event.originInfo)
		}
	}
}

func (a *App) handleWindowEvent(event *windowEvent) {
	defer handlePanic()
	// Get window from window map
	a.windowsLock.RLock()
	window, ok := a.windows[event.WindowID]
	a.windowsLock.RUnlock()
	if !ok {
		// Post-removal lifecycle notifications are expected: the default
		// WindowClosing listener removes the window from the manager, then
		// AppKit (or the equivalent on other platforms) keeps posting
		// windowWillClose / windowDidResignKey / etc. for the same window.
		// On darwin hasListeners always returns true today, so those
		// notifications are queued unconditionally and would warn here on
		// every window close. The same applies to App.cleanup nilling the
		// map during shutdown. None of these are bugs — just log them.
		a.debug("Window event for unknown window", "windowID", event.WindowID, "eventID", event.EventID)
		return
	}
	window.HandleWindowEvent(event.EventID)
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
	// Release any global shortcuts the application registered with the OS.
	if a.GlobalShortcut != nil {
		if err := a.GlobalShortcut.UnregisterAll(); err != nil {
			a.handleError(err)
		}
	}
	InvokeSync(func() {
		a.shutdownServices()
		a.windowsLock.Lock()
		for _, window := range a.windows {
			window.Close()
		}
		a.windows = nil
		a.windowsLock.Unlock()
		a.nativeWindowsLock.Lock()
		nativeWindows := make([]*NativeWindow, 0, len(a.nativeWindows))
		for _, window := range a.nativeWindows {
			nativeWindows = append(nativeWindows, window)
		}
		a.nativeWindows = nil
		a.nativeWindowsLock.Unlock()
		for _, window := range nativeWindows {
			window.Close()
		}
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

func (a *App) shouldQuit() bool {
	if a.options.ShouldQuit != nil {
		return a.options.ShouldQuit()
	}
	return true
}
