package application

import (
	"context"
	"embed"
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
	"sync/atomic"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/capabilities"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/updater"
)

//go:embed assets/*
var alphaAssets embed.FS

var globalApplication *App
var startEventLoopsOnce sync.Once
var activeEventLoopApp atomic.Pointer[App]

// AlphaAssets is the default assets for the alpha application
var AlphaAssets = AssetOptions{
	Handler: BundledAssetFileServer(alphaAssets),
}

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
	updater.HandleHelperMode()

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

	messageProc := NewMessageProcessor(result.Logger)
	result.messageProcessor = messageProc

	// Initialize transport (default to HTTP if not specified)
	transport := appOptions.Transport
	if transport == nil {
		transport = NewHTTPTransport(HTTPTransportWithLogger(result.Logger))
	}

	err := transport.Start(result.ctx, messageProc)
	if err != nil {
		result.fatal("failed to start custom transport: %w", err)
	}
	// Register shutdown task to stop transport
	result.OnShutdown(func() {
		if err := transport.Stop(); err != nil {
			result.error("failed to stop custom transport: %w", err)
		}
	})

	// Auto-wire events if transport supports event delivery
	if eventTransport, ok := transport.(WailsEventListener); ok {
		result.wailsEventListeners = append(result.wailsEventListeners, eventTransport)
	} else {
		// otherwise fallback to IPC
		result.wailsEventListeners = append(result.wailsEventListeners, &EventIPCTransport{
			app: result,
		})
	}

	middlewares := []assetserver.Middleware{
		func(next http.Handler) http.Handler {
			if m := appOptions.Assets.Middleware; m != nil {
				return m(next)
			}
			return next
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				path := req.URL.Path
				// Oversized events are parked host-side and fetched here rather
				// than being spliced into an evaluateJavaScript source string.
				if strings.HasPrefix(path, eventPayloadPath) {
					result.serveEventPayload(rw, req)
					return
				}
				// GoStream: the poll is held here for as long as the frontend
				// has nothing to collect. Safe because every webview request
				// gets its own goroutine (see the dispatchWorkers note in
				// assetserver_webview.go).
				if strings.HasPrefix(path, streamPath) {
					result.serveStream(rw, req)
					return
				}
				switch path {
				case "/wails/runtime.js":
					// The prelude, where there is one, picks the stream
					// transport before any module body runs. It cannot be
					// deferred to custom.js — see stream_prelude_server.go.
					err := assetserver.ServeFile(rw, path, runtimeJSWithPrelude())
					if err != nil {
						result.fatal("unable to serve runtime.js: %w", err)
					}
				case "/wails/transport.js":
					err := assetserver.ServeFile(rw, path, transport.JSClient())
					if err != nil {
						result.fatal("unable to serve transport.js: %w", err)
					}
				case "/wails/custom.js":
					// custom.js is only served in server mode.
					// Return 404 so the runtime's loadOptionalScript skips it.
					http.NotFound(rw, req)
				default:
					next.ServeHTTP(rw, req)
				}
			})
		},
	}

	if handler, ok := transport.(TransportHTTPHandler); ok {
		middlewares = append(middlewares, handler.Handler())
	}

	opts := &assetserver.Options{
		Handler:    appOptions.Assets.Handler,
		Middleware: assetserver.ChainMiddleware(middlewares...),
		Logger:     result.Logger,
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

	// If transport implements AssetServerTransport, configure it to serve assets
	if assetTransport, ok := transport.(AssetServerTransport); ok {
		err := assetTransport.ServeAssets(srv)
		if err != nil {
			result.fatal("failed to configure transport for serving assets: %w", err)
		}
		result.debug("Transport configured to serve assets")
	}

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

var _ webview.Request = &webViewAssetRequest{}

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

const webViewRequestHeaderWindowId = "x-wails-window-id"
const webViewRequestHeaderWindowName = "x-wails-window-name"

type webViewAssetRequest struct {
	Request    webview.Request
	windowId   uint
	windowName string
}

var windowKeyEvents = make(chan *windowKeyEvent, 5)

type windowKeyEvent struct {
	windowId          uint
	acceleratorString string
}

func (r *webViewAssetRequest) URL() (string, error) {
	return r.Request.URL()
}

func (r *webViewAssetRequest) Method() (string, error) {
	return r.Request.Method()
}

func (r *webViewAssetRequest) Header() (http.Header, error) {
	h, err := r.Request.Header()
	if err != nil {
		return nil, err
	}

	hh := h.Clone()
	hh.Set(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(r.windowId), 10))
	if r.windowName != "" {
		hh.Set(webViewRequestHeaderWindowName, r.windowName)
	}
	return hh, nil
}

func (r *webViewAssetRequest) Body() (io.ReadCloser, error) {
	return r.Request.Body()
}

func (r *webViewAssetRequest) Response() webview.ResponseWriter {
	return r.Request.Response()
}

func (r *webViewAssetRequest) Close() error {
	return r.Request.Close()
}

var webviewRequests = make(chan *webViewAssetRequest, 256)

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
	Window         *WindowManager
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
	Updater        *updater.Updater

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

	// The main application menu (private - use app.Menu.GetApplicationMenu/SetApplicationMenu)
	applicationMenu *Menu

	clipboard            *Clipboard
	customEventProcessor *EventProcessor
	Logger               *slog.Logger

	contextMenus     map[string]*ContextMenu
	contextMenusLock sync.RWMutex

	assets *assetserver.AssetServer

	// eventPayloads holds oversized Go→JS event bodies awaiting a one-shot
	// fetch from the webview, keeping them out of evaluateJavaScript source.
	eventPayloads *eventPayloadStore

	// streams holds registered stream handlers and the per-page sessions that
	// carry their connections. See stream.go.
	streams *streamManager

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
	a.systemTrays = make(map[uint]*SystemTray)
	a.contextMenus = make(map[string]*ContextMenu)
	a.keyBindings = make(map[string]func(window Window))
	a.Logger = a.options.Logger
	a.pid = os.Getpid()
	a.wailsEventListeners = make([]WailsEventListener, 0)

	// Initialize managers
	a.Window = newWindowManager(a)
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
	a.Updater = updater.New(newUpdaterHost(a))
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

	// Start event handling before entering the native event loop. The remaining
	// initialisation is triggered by the platform's ApplicationStarted event, so
	// native window operations are available while services start in the
	// background.
	a.startEventLoops()
	activeEventLoopApp.Store(a)
	defer activeEventLoopApp.CompareAndSwap(a, nil)

	var startupOnce sync.Once
	var startupErr error
	var startupErrLock sync.Mutex
	startupStarted := make(chan struct{})
	startupDone := make(chan struct{})
	a.Event.RegisterApplicationEventHook(events.Common.ApplicationStarted, func(*ApplicationEvent) {
		startupOnce.Do(func() {
			close(startupStarted)
			go func() {
				defer close(startupDone)
				if err := a.finishStartup(); err != nil {
					startupErrLock.Lock()
					startupErr = err
					startupErrLock.Unlock()
					// Startup now happens asynchronously after the platform loop
					// begins, so report the fatal failure immediately as well as
					// returning it once the platform loop exits.
					a.handleError(&FatalError{err: err})
					a.Quit()
				}
			}()
		})
	})

	a.runLock.Lock()
	a.running = true
	a.runLock.Unlock()

	// Bind any global shortcuts that were registered before the app started.
	a.GlobalShortcut.flushPending()

	// The menu and icon must be configured before the platform loop starts on
	// platforms that consume them during native activation.
	if runtime.GOOS == "darwin" {
		a.impl.setApplicationMenu(a.applicationMenu)
	}
	if a.options.Icon != nil {
		a.impl.setIcon(a.options.Icon)
	}

	runErr := a.impl.run()
	// Prevent a late native ready event from starting services after the
	// platform loop has already begun shutting down.
	startupOnce.Do(func() {})
	select {
	case <-startupStarted:
		<-startupDone
	default:
	}
	startupErrLock.Lock()
	defer startupErrLock.Unlock()
	if startupErr != nil {
		return startupErr
	}
	return runErr
}

func (a *App) finishStartup() error {
	// Dispatch framework lifecycle events directly so their listeners are
	// scheduled before the next startup phase begins. Native application events
	// continue to arrive through applicationEvents.
	a.Event.handleApplicationEventSync(newApplicationEvent(events.Common.ApplicationStarting))

	// No need to hold the lock here because a.options.Services may only change
	// before Run is called. Services remain sequential even though the startup
	// sequence itself runs away from the native UI thread.
	a.serviceShutdownLock.Lock()
	services := a.options.Services
	a.options.Services = nil
	for i, service := range services {
		if err := a.startupService(service); err != nil {
			a.serviceShutdownLock.Unlock()
			return fmt.Errorf("error starting service '%s': %w", getServiceName(service), err)
		}
		// Schedule started services for shutdown.
		a.options.Services = services[:i+1]
	}
	a.serviceShutdownLock.Unlock()

	// Start the MCP server when the application is built with -tags mcp.
	// All configuration is read from environment variables (WAILS_MCP_HOST,
	// WAILS_MCP_PORT, WAILS_MCP_TIMEOUT, WAILS_MCP_HIDE_CURSOR).
	if err := startMCPServer(a); err != nil {
		return fmt.Errorf("mcp: %w", err)
	}

	a.Event.handleApplicationEventSync(newApplicationEvent(events.Common.ApplicationInitialized))

	// Normal windows remain pending until all subsystems and services have
	// initialised. A listener for ApplicationStarting may explicitly Show a
	// pending window (for example a splash screen) while this work is running.
	a.runLock.Lock()
	pendingRun := a.pendingRun
	a.pendingRun = nil
	a.runLock.Unlock()
	for _, pending := range pendingRun {
		go func() {
			defer handlePanic()
			pending.Run()
		}()
	}
	return nil
}

func (a *App) startEventLoops() {
	startEventLoopsOnce.Do(func() {
		go func() {
			for event := range applicationEvents {
				if app := activeEventLoopApp.Load(); app != nil {
					go app.Event.handleApplicationEvent(event)
				}
			}
		}()
		go func() {
			for event := range windowEvents {
				if app := activeEventLoopApp.Load(); app != nil {
					go app.handleWindowEvent(event)
				}
			}
		}()
		go func() {
			for request := range webviewRequests {
				if app := activeEventLoopApp.Load(); app != nil {
					go app.handleWebViewRequest(request)
				}
			}
		}()
		go func() {
			for event := range windowMessageBuffer {
				if app := activeEventLoopApp.Load(); app != nil {
					go app.handleWindowMessage(event)
				}
			}
		}()
		go func() {
			for event := range windowKeyEvents {
				if app := activeEventLoopApp.Load(); app != nil {
					go app.handleWindowKeyEvent(event)
				}
			}
		}()
		go func() {
			for message := range windowDragAndDropBuffer {
				if app := activeEventLoopApp.Load(); app != nil {
					go app.handleDragAndDropMessage(message)
				}
			}
		}()
		go func() {
			for menuItemID := range menuItemClicked {
				if app := activeEventLoopApp.Load(); app != nil {
					go app.Menu.handleMenuItemClicked(menuItemID)
				}
			}
		}()
	})
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
					fmt.Sprintf(
						"Service '%s' does not handle HTTP requests",
						getServiceName(service),
					),
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
	// Cancel first so an in-progress ServiceStartup can stop before shutdown
	// waits for the startup goroutine to release serviceShutdownLock.
	a.cancel()

	// Acquire lock to prevent double calls (defer in Run() + OnShutdown)
	a.serviceShutdownLock.Lock()
	defer a.serviceShutdownLock.Unlock()

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

func (a *App) handleWebViewRequest(request *webViewAssetRequest) {
	defer handlePanic()
	// Log that we're processing the request
	url, _ := request.Request.URL()
	a.debug("handleWebViewRequest: Processing request", "url", url)
	// IMPORTANT: pass the wrapper request so our injected headers (x-wails-window-id/name) are used
	a.assets.ServeWebViewRequest(request)
	a.debug("handleWebViewRequest: Request processing complete", "url", url)
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
