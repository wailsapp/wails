//go:build windows
// +build windows

package edge

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/webview2/internal/w32"
	"github.com/wailsapp/wails/v3/internal/webview2/webviewloader"
	"golang.org/x/sys/windows"
)

type Rect = w32.Rect

func globalErrorHandler(err error) {
	if err == nil {
		return
	}

	log.Printf("[WebView2 Error] %v\n", err)

	stackBuf := make([]uintptr, 64)
	stackSize := runtime.Callers(2, stackBuf)
	frames := runtime.CallersFrames(stackBuf[:stackSize])

	log.Printf("\nStack trace:")
	stackIndex := 1
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		log.Printf("%d: %s\n\t%s:%d\n", stackIndex, frame.Function, frame.File, frame.Line)
		stackIndex++
	}
}

type Chromium struct {
	hwnd    uintptr
	padding struct {
		Left   int32
		Top    int32
		Right  int32
		Bottom int32
	}

	controller                       *ICoreWebView2Controller
	compositionController            *ICoreWebView2CompositionController
	compositionController4           *ICoreWebView2CompositionController4
	webview                          *ICoreWebView2
	inited                           uintptr
	envCompleted                     *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
	controllerCompleted              *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler
	compositionControllerCompleted   *iCoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler
	webMessageReceived               *iCoreWebView2WebMessageReceivedEventHandler
	containsFullScreenElementChanged *ICoreWebView2ContainsFullScreenElementChangedEventHandler
	permissionRequested              *iCoreWebView2PermissionRequestedEventHandler
	webResourceRequested             *iCoreWebView2WebResourceRequestedEventHandler
	acceleratorKeyPressed            *ICoreWebView2AcceleratorKeyPressedEventHandler
	cursorChanged                    *iCoreWebView2CursorChangedEventHandler
	navigationStarting               *ICoreWebView2NavigationStartingEventHandler
	navigationCompleted              *ICoreWebView2NavigationCompletedEventHandler
	processFailed                    *ICoreWebView2ProcessFailedEventHandler

	environment            *ICoreWebView2Environment
	webview2RuntimeVersion string
	compositionHost        *compositionHost

	// Settings
	Debug                         bool
	DataPath                      string
	BrowserPath                   string
	AdditionalBrowserArgs         []string
	NonClientRegionSupportEnabled bool
	CompositionControllerEnabled  bool

	// permissions
	permissions      map[CoreWebView2PermissionKind]CoreWebView2PermissionState
	globalPermission *CoreWebView2PermissionState

	// Callbacks
	MessageCallback                          func(message string, sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs)
	MessageWithAdditionalObjectsCallback     func(message string, sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs)
	WebResourceRequestedCallback             func(request *ICoreWebView2WebResourceRequest, args *ICoreWebView2WebResourceRequestedEventArgs)
	NavigationStartingCallback               func(sender *ICoreWebView2)
	NavigationCompletedCallback              func(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs)
	ProcessFailedCallback                    func(sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs)
	ContainsFullScreenElementChangedCallback func(sender *ICoreWebView2, args *ICoreWebView2ContainsFullScreenElementChangedEventArgs)
	AcceleratorKeyCallback                   func(uint) bool
	CursorChangedCallback                    func(cursor HCURSOR, systemCursorID uint32)

	// Error handling
	globalErrorCallback func(error)

	shuttingDown bool

	// Resize debouncing
	lastBounds  *w32.Rect
	resizeTimer *time.Timer
}

func NewChromium() *Chromium {
	e := &Chromium{}
	/*
	 All these handlers are passed to native code through syscalls with 'uintptr(unsafe.Pointer(handler))' and we know
	 that a pointer to those will be kept in the native code. Furthermore these handlers als contain pointer to other Go
	 structs like the vtable.
	 This violates the unsafe.Pointer rule '(4) Conversion of a Pointer to a uintptr when calling syscall.Syscall.' because
	 theres no guarantee that Go doesn't move these objects.
	 AFAIK currently the Go runtime doesn't move HEAP objects, so we should be safe with these handlers. But they don't
	 guarantee it, because in the future Go might use a compacting GC.
	 There's a proposal to add a runtime.Pin function, to prevent moving pinned objects, which would allow to easily fix
	 this issue by just pinning the handlers. The https://go-review.googlesource.com/c/go/+/367296/ should land in Go 1.19.
	*/
	e.envCompleted = newICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(e)
	e.controllerCompleted = newICoreWebView2CreateCoreWebView2ControllerCompletedHandler(e)
	e.compositionControllerCompleted = newICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler(e)
	e.webMessageReceived = newICoreWebView2WebMessageReceivedEventHandler(e)
	e.permissionRequested = newICoreWebView2PermissionRequestedEventHandler(e)
	e.webResourceRequested = newICoreWebView2WebResourceRequestedEventHandler(e)
	e.acceleratorKeyPressed = newICoreWebView2AcceleratorKeyPressedEventHandler(e)
	e.cursorChanged = newICoreWebView2CursorChangedEventHandler(e)
	e.navigationStarting = newICoreWebView2NavigationStartingEventHandler(e)
	e.navigationCompleted = newICoreWebView2NavigationCompletedEventHandler(e)
	e.processFailed = newICoreWebView2ProcessFailedEventHandler(e)
	e.containsFullScreenElementChanged = newICoreWebView2ContainsFullScreenElementChangedEventHandler(e)
	/*
		// Pinner seems to panic in some cases as reported on Discord, maybe during shutdown when GC detects pinned objects
		// to be released that have not been unpinned.
		// It would also be better to use our ComBridge for this event handlers implementation instead of pinning them.
		// So all COM Implementations on the go-side use the same code.
		var pinner runtime.Pinner
		pinner.Pin(e.envCompleted)
		pinner.Pin(e.controllerCompleted)
		pinner.Pin(e.webMessageReceived)
		pinner.Pin(e.permissionRequested)
		pinner.Pin(e.webResourceRequested)
		pinner.Pin(e.acceleratorKeyPressed)
		pinner.Pin(e.navigationCompleted)
		pinner.Pin(e.processFailed)
		pinner.Pin(e.containsFullScreenElementChanged)
	*/
	e.permissions = make(map[CoreWebView2PermissionKind]CoreWebView2PermissionState)
	e.globalErrorCallback = globalErrorHandler
	return e
}

func (e *Chromium) ShuttingDown() {
	e.shuttingDown = true
}

func (e *Chromium) errorCallback(err error) {
	e.globalErrorCallback(err)
	os.Exit(1)
}

func (e *Chromium) SetErrorCallback(callback func(error)) {
	if callback != nil {
		e.globalErrorCallback = callback
	}
}

func (e *Chromium) SetCursorChangedCallback(callback func(cursor HCURSOR, systemCursorID uint32)) {
	if callback != nil {
		e.CursorChangedCallback = callback
	}
}

func (e *Chromium) Embed(hwnd uintptr) bool {

	var err error

	e.hwnd = hwnd

	dataPath := e.DataPath
	if dataPath == "" {
		currentExePath := make([]uint16, windows.MAX_PATH)
		_, err = windows.GetModuleFileName(windows.Handle(0), &currentExePath[0], windows.MAX_PATH)
		if err != nil {
			e.errorCallback(err)
		}
		currentExeName := filepath.Base(windows.UTF16ToString(currentExePath))
		dataPath = filepath.Join(os.Getenv("AppData"), currentExeName)
	}

	if e.BrowserPath != "" {
		if _, err = os.Stat(e.BrowserPath); errors.Is(err, os.ErrNotExist) {
			e.errorCallback(fmt.Errorf("browser path '%s' does not exist", e.BrowserPath))
		}
	}

	browserArgs := strings.Join(e.AdditionalBrowserArgs, " ")
	if err := createCoreWebView2EnvironmentWithOptions(e.BrowserPath, dataPath, e.envCompleted, browserArgs); err != nil {
		e.errorCallback(fmt.Errorf("error calling Webview2Loader: %s", err.Error()))
	}

	e.webview2RuntimeVersion, err = webviewloader.GetAvailableCoreWebView2BrowserVersionString(e.BrowserPath)
	if err != nil {
		e.errorCallback(fmt.Errorf("error getting Webview2 runtime version: %s", err.Error()))
	}

	var msg w32.Msg
	for {
		if atomic.LoadUintptr(&e.inited) != 0 {
			break
		}
		r, _, _ := w32.User32GetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if r == 0 {
			break
		}
		w32.User32TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		w32.User32DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
	e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")
	return true
}

func (e *Chromium) SetPadding(padding Rect) {
	if e.padding.Left == padding.Left && e.padding.Top == padding.Top &&
		e.padding.Right == padding.Right && e.padding.Bottom == padding.Bottom {

		return
	}

	e.padding.Left = padding.Left
	e.padding.Top = padding.Top
	e.padding.Right = padding.Right
	e.padding.Bottom = padding.Bottom
	e.Resize()
}

func (e *Chromium) ResizeWithBounds(bounds *Rect) {
	if e.hwnd == 0 {
		return
	}

	bounds.Top += e.padding.Top
	bounds.Bottom -= e.padding.Bottom
	bounds.Left += e.padding.Left
	bounds.Right -= e.padding.Right

	e.SetSize(*bounds)
}

func (e *Chromium) Resize() {
	if e.hwnd == 0 {
		return
	}

	bounds, err := w32.GetClientRect(e.hwnd)
	if err != nil {
		// GetClientRect can fail transiently while the window is being torn
		// down or reconfigured during DPI churn. Skipping a resize frame is
		// recoverable; killing the process (errorCallback) is not.
		log.Printf("[WebView2] Resize failed to get client rect: %v", err)
		return
	}

	e.ResizeWithBounds(&bounds)
}

func (e *Chromium) Navigate(url string) {
	err := e.webview.Navigate(url)
	if err != nil {
		// A failed navigation is recoverable (the previous content stays
		// visible); killing the process is not.
		log.Printf("[WebView2] Navigate failed: %v", err)
	}
}

func (e *Chromium) NavigateToString(content string) {
	err := e.webview.NavigateToString(content)
	if err != nil {
		log.Printf("[WebView2] NavigateToString failed: %v", err)
	}
}

func (e *Chromium) Init(script string) {
	err := e.webview.AddScriptToExecuteOnDocumentCreated(script, nil)
	if err != nil {
		log.Printf("[WebView2] Init script registration failed: %v", err)
	}
}

func (e *Chromium) Eval(script string) {
	if e.webview == nil || e.shuttingDown {
		return
	}

	err := e.webview.ExecuteScript(script, nil)
	if err != nil && !errors.Is(err, windows.ERROR_IO_PENDING) {
		// ExecuteScript fails transiently while the browser process is busy
		// reconfiguring — e.g. RESOURCE_NOT_IN_CORRECT_STATE during a DPI
		// transition when the window is dragged between mixed-DPI monitors
		// (wailsapp/wails#5544). Script execution is fire-and-forget; a
		// dropped script during a transition is recoverable, killing the
		// process (errorCallback) is not.
		log.Printf("[WebView2] Eval failed: %v", err)
	}
}

func (e *Chromium) Show() error {
	return e.controller.PutIsVisible(true)
}

func (e *Chromium) Hide() error {
	return e.controller.PutIsVisible(false)
}

func (e *Chromium) QueryInterface(_, _ uintptr) uintptr {
	return 0
}

func (e *Chromium) AddRef() uintptr {
	return 1
}

func (e *Chromium) Release() uintptr {
	return 1
}

func (e *Chromium) EnvironmentCompleted(res uintptr, env *ICoreWebView2Environment) uintptr {
	if env == nil {
		err := syscall.Errno(res)
		log.Printf("[WebView2] Environment creation failed with error code %v: %v\n", res, err)
		if e.globalErrorCallback != nil {
			e.globalErrorCallback(fmt.Errorf("failed to create WebView2 environment: %w", err))
		}
		return res
	}

	log.Printf("[WebView2] Environment created successfully\n")

	env.AddRef()
	e.environment = env

	var err error
	if !e.CompositionControllerEnabled {
		err = env.CreateCoreWebView2Controller(e.hwnd, e.controllerCompleted)
	} else {
		err = e.createCoreWebView2CompositionController(env)
		if err != nil {
			err = e.fallbackToCoreWebView2Controller(fmt.Errorf("composition controller setup failed: %w", err))
		}
	}
	if err != nil {
		e.errorCallback(err)
	}
	return 0
}

func (e *Chromium) CreateCoreWebView2ControllerCompleted(res uintptr, controller *ICoreWebView2Controller) uintptr {
	if int32(res) < 0 {
		e.errorCallback(fmt.Errorf("error creating controller with %08x: %s", res, syscall.Errno(res)))
	}

	return e.initializeController(controller)
}

func (e *Chromium) createCoreWebView2CompositionController(env *ICoreWebView2Environment) error {
	env3 := env.GetICoreWebView2Environment3()
	if env3 == nil {
		return UnsupportedCapabilityError
	}
	defer env3.Release()

	host, err := newCompositionHost(e.hwnd)
	if err != nil {
		return err
	}
	e.compositionHost = host

	return env3.CreateCoreWebView2CompositionController(e.hwnd, e.compositionControllerCompleted)
}

func (e *Chromium) CreateCoreWebView2CompositionControllerCompleted(res uintptr, compositionController *ICoreWebView2CompositionController) uintptr {
	if int32(res) < 0 {
		if err := e.fallbackToCoreWebView2Controller(fmt.Errorf("error creating composition controller with %08x: %s", res, syscall.Errno(res))); err != nil {
			e.errorCallback(err)
		}
		return 0
	}

	if compositionController == nil {
		if err := e.fallbackToCoreWebView2Controller(fmt.Errorf("composition controller completed without a controller")); err != nil {
			e.errorCallback(err)
		}
		return 0
	}

	compositionController.AddRef()
	e.compositionController = compositionController
	e.compositionController4 = compositionController.GetICoreWebView2CompositionController4()

	if err := e.compositionHost.attachController(e.compositionController); err != nil {
		if fallbackErr := e.fallbackToCoreWebView2Controller(fmt.Errorf("attaching composition controller failed: %w", err)); fallbackErr != nil {
			e.errorCallback(fallbackErr)
		}
		return 0
	}

	controller := compositionController.GetICoreWebView2Controller()
	if controller == nil {
		if err := e.fallbackToCoreWebView2Controller(fmt.Errorf("error getting controller from composition controller")); err != nil {
			e.errorCallback(err)
		}
		return 0
	}

	return e.initializeController(controller)
}

func (e *Chromium) fallbackToCoreWebView2Controller(reason error) error {
	log.Printf("[WebView2] Composition hosting unavailable, falling back to HWND controller: %v\n", reason)

	e.releaseCompositionController()
	e.releaseCompositionHost()
	e.CompositionControllerEnabled = false

	if e.environment == nil {
		return fmt.Errorf("cannot fall back to WebView2 HWND controller without an environment: %w", reason)
	}
	if err := e.environment.CreateCoreWebView2Controller(e.hwnd, e.controllerCompleted); err != nil {
		return fmt.Errorf("composition hosting failed (%v), and HWND controller fallback failed: %w", reason, err)
	}
	return nil
}

func (e *Chromium) releaseCompositionController() {
	if e.compositionController4 != nil {
		e.compositionController4.Release()
		e.compositionController4 = nil
	}
	if e.compositionController != nil {
		if controller := e.compositionController.GetICoreWebView2Controller(); controller != nil {
			_ = controller.Close()
			controller.Release()
		}
		e.compositionController.Release()
		e.compositionController = nil
	}
}

func (e *Chromium) releaseCompositionHost() {
	if e.compositionHost != nil {
		e.compositionHost.release()
		e.compositionHost = nil
	}
}

func (e *Chromium) initializeController(controller *ICoreWebView2Controller) uintptr {
	var err error

	controller.AddRef()
	e.controller = controller

	// Try to get ICoreWebView2Controller3 interface for better performance.
	// Keep composition-hosted WebViews on WebView2-managed scale detection:
	// the DirectComposition surface is owned by WebView2, and manually
	// reasserting RasterizationScale during monitor transitions can leave that
	// surface black after DPI increases.
	if controller3 := e.controller.GetICoreWebView2Controller3(); controller3 != nil {
		if !e.CompositionControllerEnabled {
			// Use raw pixels mode for better performance during resize.
			if err := controller3.PutBoundsMode(COREWEBVIEW2_BOUNDS_MODE_USE_RAW_PIXELS); err != nil {
				e.errorCallback(err)
			}

			// ShouldDetectMonitorScaleChanges is deliberately left at its default
			// (enabled): WebView2 tracks monitor DPI changes and updates its own
			// rasterization scale. This module previously disabled it and left the
			// scale to the host's WM_DPICHANGED handling, but an externally
			// written scale races the browser process's internal display
			// bookkeeping during a mixed-DPI monitor cross — the browser can land
			// on a degenerate scale(0,0) compositor transform, which kills the
			// GPU process on every frame until the browser process itself exits
			// (wailsapp/wails#5732). Detection stays on so the scale has exactly
			// one writer, on the code path every mainstream embedder exercises;
			// bounds remain raw pixels, which is orthogonal.
		}
	}
	var token _EventRegistrationToken
	e.webview, err = e.controller.GetCoreWebView2()
	if err != nil {
		e.errorCallback(err)
	}

	e.webview.AddRef()
	if e.NonClientRegionSupportEnabled {
		if err := e.PutIsNonClientRegionSupportEnabled(true); err != nil {
			if !errors.Is(err, UnsupportedCapabilityError) {
				e.errorCallback(err)
			}
		}
	}
	err = e.webview.AddWebMessageReceived(e.webMessageReceived, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddPermissionRequested(e.permissionRequested, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddWebResourceRequested(e.webResourceRequested, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddNavigationStarting(e.navigationStarting, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddNavigationCompleted(e.navigationCompleted, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddProcessFailed(e.processFailed, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddContainsFullScreenElementChanged(e.containsFullScreenElementChanged, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.controller.AddAcceleratorKeyPressed(e.acceleratorKeyPressed, &token)
	if err != nil {
		e.errorCallback(err)
	}
	if e.compositionController != nil {
		err = e.compositionController.AddCursorChanged(e.cursorChanged, &token)
		if err != nil {
			e.errorCallback(err)
		}
	}

	atomic.StoreUintptr(&e.inited, 1)

	return 0
}

func (e *Chromium) ContainsFullScreenElementChanged(sender *ICoreWebView2, args *ICoreWebView2ContainsFullScreenElementChangedEventArgs) uintptr {
	if e.ContainsFullScreenElementChangedCallback != nil {
		e.ContainsFullScreenElementChangedCallback(sender, args)
	}
	return 0
}

func (e *Chromium) CursorChanged(sender *ICoreWebView2CompositionController, _ *IUnknown) uintptr {
	if e.CursorChangedCallback == nil {
		return 0
	}

	cursor, err := sender.GetCursor()
	if err != nil {
		e.errorCallback(err)
		return 0
	}

	systemCursorID, err := sender.GetSystemCursorId()
	if err != nil {
		e.errorCallback(err)
		return 0
	}

	e.CursorChangedCallback(cursor, systemCursorID)
	return 0
}

func (e *Chromium) MessageReceived(sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	message, err := args.TryGetWebMessageAsString()
	if err != nil {
		// The web message originates from (potentially untrusted) page
		// content. A malformed message must never be able to take the whole
		// process down — drop it and keep running.
		log.Printf("[WebView2] dropping malformed web message: %v", err)
		return 0
	}

	if HasCapability(e.webview2RuntimeVersion, GetAdditionalObjects) {
		obj, err := args.GetAdditionalObjects()
		if err != nil {
			// Fall back to delivering the plain string message below rather
			// than killing the process.
			log.Printf("[WebView2] failed to read additional objects, delivering message without them: %v", err)
			obj = nil
		}

		if obj != nil && e.MessageWithAdditionalObjectsCallback != nil {
			defer obj.Release()
			e.MessageWithAdditionalObjectsCallback(message, sender, args)
		} else if e.MessageCallback != nil {
			e.MessageCallback(message, sender, args)
		}
	} else if e.MessageCallback != nil {
		e.MessageCallback(message, sender, args)
	}

	err = sender.PostWebMessageAsString(message)
	if err != nil && !errors.Is(err, windows.ERROR_IO_PENDING) {
		log.Printf("[WebView2] PostWebMessageAsString failed: %v", err)
	}
	return 0
}

func (e *Chromium) SetPermission(kind CoreWebView2PermissionKind, state CoreWebView2PermissionState) {
	e.permissions[kind] = state
}

func (e *Chromium) SetBackgroundColour(R, G, B, A uint8) {
	controller := e.GetController()
	controller2 := controller.GetICoreWebView2Controller2()

	backgroundCol := COREWEBVIEW2_COLOR{
		A: A,
		R: R,
		G: G,
		B: B,
	}

	// WebView2 only has 0 and 255 as valid values.
	if backgroundCol.A > 0 && backgroundCol.A < 255 {
		backgroundCol.A = 255
	}

	err := controller2.PutDefaultBackgroundColor(backgroundCol)
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) SetGlobalPermission(state CoreWebView2PermissionState) {
	e.globalPermission = &state
}

func (e *Chromium) PermissionRequested(_ *ICoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr {
	kind, err := args.GetPermissionKind()
	if err != nil {
		e.errorCallback(err)
	}
	var result CoreWebView2PermissionState
	if e.globalPermission != nil {
		result = *e.globalPermission
	} else {
		var ok bool
		result, ok = e.permissions[kind]
		if !ok {
			result = CoreWebView2PermissionStateDefault
		}
	}
	err = args.PutState(result)
	if err != nil {
		e.errorCallback(err)
	}
	return 0
}

func (e *Chromium) WebResourceRequested(sender *ICoreWebView2, args *ICoreWebView2WebResourceRequestedEventArgs) uintptr {
	req, err := args.GetRequest()
	if err != nil {
		log.Fatal(err)
	}
	defer req.Release()

	if e.WebResourceRequestedCallback != nil {
		e.WebResourceRequestedCallback(req, args)
	}
	return 0
}

func (e *Chromium) AddWebResourceRequestedFilter(filter string, ctx COREWEBVIEW2_WEB_RESOURCE_CONTEXT) {
	err := e.webview.AddWebResourceRequestedFilter(filter, ctx)
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) Environment() *ICoreWebView2Environment {
	return e.environment
}

// AcceleratorKeyPressed is called when an accelerator key is pressed.
// If the AcceleratorKeyCallback method has been set, it will defer handling of the keypress
// to the callback. That callback returns a bool indicating if the event was handled.
func (e *Chromium) AcceleratorKeyPressed(sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr {
	if e.AcceleratorKeyCallback == nil {
		return 0
	}
	eventKind, _ := args.GetKeyEventKind()
	if eventKind == COREWEBVIEW2_KEY_EVENT_KIND_KEY_DOWN ||
		eventKind == COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN {
		virtualKey, _ := args.GetVirtualKey()
		status, _ := args.GetPhysicalKeyStatus()
		if !status.WasKeyDown {
			err := args.PutHandled(e.AcceleratorKeyCallback(virtualKey))
			if err != nil {
				e.errorCallback(err)
			}
		} else {
			return 0
		}
	}
	err := args.PutHandled(false)
	if err != nil {
		e.errorCallback(err)
	}
	return 0
}

func (e *Chromium) GetSettings() (*ICoreWebViewSettings, error) {
	return e.webview.GetSettings()
}

func (e *Chromium) GetController() *ICoreWebView2Controller {
	return e.controller
}

// IsReady reports whether the WebView2 controller has been fully initialised.
// e.controller is assigned partway through CreateCoreWebView2ControllerCompleted,
// before the controller's COM setup has finished, so a non-nil controller is
// not sufficient to safely call into it: COM calls made in that window fail
// with E_INVALIDARG ("The parameter is incorrect"). The inited flag is set
// only after setup completes. Safe to call from any goroutine.
func (e *Chromium) IsReady() bool {
	return atomic.LoadUintptr(&e.inited) != 0
}

func boolToInt(input bool) int {
	if input {
		return 1
	}
	return 0
}

func (e *Chromium) NavigationStarting(sender *ICoreWebView2, _ *IUnknown) uintptr {
	if e.NavigationStartingCallback != nil {
		e.NavigationStartingCallback(sender)
	}
	return 0
}

func (e *Chromium) NavigationCompleted(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
	if e.NavigationCompletedCallback != nil {
		e.NavigationCompletedCallback(sender, args)
	}
	return 0
}

func (e *Chromium) ProcessFailed(sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr {
	if e.ProcessFailedCallback != nil {
		e.ProcessFailedCallback(sender, args)
	}
	return 0
}

func (e *Chromium) Bounds() *Rect {
	if e == nil || e.controller == nil {
		return nil
	}

	rect, err := e.controller.GetBounds()
	if err != nil {
		e.errorCallback(err)
		return nil
	}
	return rect
}

func (e *Chromium) NotifyParentWindowPositionChanged() error {
	//It looks like the wndproc function is called before the controller initialization is complete.
	//Because of this the controller is nil
	if e.controller == nil {
		return nil
	}
	return e.controller.NotifyParentWindowPositionChanged()
}

func (e *Chromium) Focus() {
	// The WndProc can dispatch WM_SETFOCUS re-entrantly while the controller
	// is still being configured in CreateCoreWebView2ControllerCompleted
	// (issue #5446). Callers' GetController() != nil checks cannot exclude
	// that window, so guard here: dropping a focus request during startup is
	// harmless, calling MoveFocus on a partially-initialised controller is
	// fatal (errorCallback exits the process).
	if !e.IsReady() {
		return
	}
	err := e.controller.MoveFocus(COREWEBVIEW2_MOVE_FOCUS_REASON_PROGRAMMATIC)
	if err != nil {
		// MoveFocus can legitimately fail after initialisation too — e.g.
		// E_INVALIDARG when the window is hidden or minimised to the tray
		// (wailsapp/wails#4158 reproduces this on tray-click restore). A
		// failed focus request is never worth killing the process, which is
		// what errorCallback does; log it instead.
		log.Printf("[WebView2] Focus failed: %v", err)
	}
}

func (e *Chromium) PutZoomFactor(zoomFactor float64) {
	err := e.controller.PutZoomFactor(zoomFactor)
	if err != nil {
		log.Printf("[WebView2] PutZoomFactor failed: %v", err)
	}
}

func (e *Chromium) OpenDevToolsWindow() {
	err := e.webview.OpenDevToolsWindow()
	if err != nil {
		log.Printf("[WebView2] OpenDevToolsWindow failed: %v", err)
	}
}

func (e *Chromium) HasCapability(c Capability) bool {
	return HasCapability(e.webview2RuntimeVersion, c)
}

func (e *Chromium) CompositionControllerReady() bool {
	return e != nil && e.compositionController != nil
}

func (e *Chromium) NonClientRegionHitTestReady() bool {
	return e != nil &&
		e.compositionController4 != nil &&
		HasCapability(e.webview2RuntimeVersion, NonClientRegion)
}

func (e *Chromium) GetIsSwipeNavigationEnabled() (bool, error) {
	if !HasCapability(e.webview2RuntimeVersion, SwipeNavigation) {
		return false, UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return false, err
	}
	webview2Settings6 := webview2Settings.GetICoreWebView2Settings6()
	var result bool
	result, err = webview2Settings6.GetIsSwipeNavigationEnabled()
	if err != nil {
		return false, err
	}
	return result, nil
}

// PutIsGeneralAutofillEnabled controls whether autofill for information
// like names, street and email addresses, phone numbers, and arbitrary input
// is enabled. This excludes password and credit card information. When
// IsGeneralAutofillEnabled is false, no suggestions appear, and no new information
// is saved. When IsGeneralAutofillEnabled is true, information is saved, suggestions
// appear and clicking on one will populate the form fields.
// It will take effect immediately after setting.
// The default value is `FALSE`.
func (e *Chromium) PutIsGeneralAutofillEnabled(value bool) error {
	if !HasCapability(e.webview2RuntimeVersion, GeneralAutofillEnabled) {
		return UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return err
	}
	webview2Settings4 := webview2Settings.GetICoreWebView2Settings4()
	return webview2Settings4.PutIsGeneralAutofillEnabled(value)
}

// PutIsPasswordAutosaveEnabled sets whether the browser should offer to save passwords and other
// identifying information entered into forms automatically.
// The default value is `FALSE`.
func (e *Chromium) PutIsPasswordAutosaveEnabled(value bool) error {
	if !HasCapability(e.webview2RuntimeVersion, PasswordAutosaveEnabled) {
		return UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return err
	}
	webview2Settings4 := webview2Settings.GetICoreWebView2Settings4()
	return webview2Settings4.PutIsPasswordAutosaveEnabled(value)
}

func (e *Chromium) PutIsSwipeNavigationEnabled(enabled bool) error {
	if !HasCapability(e.webview2RuntimeVersion, SwipeNavigation) {
		return UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return err
	}
	webview2Settings6 := webview2Settings.GetICoreWebView2Settings6()
	err = webview2Settings6.PutIsSwipeNavigationEnabled(enabled)
	if err != nil {
		return err
	}
	return nil
}

func (e *Chromium) PutIsNonClientRegionSupportEnabled(enabled bool) error {
	if !HasCapability(e.webview2RuntimeVersion, NonClientRegion) {
		return UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return err
	}
	webview2Settings9 := webview2Settings.GetICoreWebView2Settings9()
	if webview2Settings9 == nil {
		return UnsupportedCapabilityError
	}
	return webview2Settings9.PutIsNonClientRegionSupportEnabled(enabled)
}

func (e *Chromium) GetNonClientRegionAtPoint(x, y int32) (COREWEBVIEW2_NON_CLIENT_REGION_KIND, bool, error) {
	if !e.NonClientRegionHitTestReady() {
		return COREWEBVIEW2_NON_CLIENT_REGION_KIND_NOWHERE, false, nil
	}

	region, err := e.compositionController4.GetNonClientRegionAtPoint(POINT{X: x, Y: y})
	if err != nil {
		return COREWEBVIEW2_NON_CLIENT_REGION_KIND_NOWHERE, false, err
	}
	return region, true, nil
}

func (e *Chromium) SendMouseInput(
	eventKind COREWEBVIEW2_MOUSE_EVENT_KIND,
	virtualKeys COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS,
	mouseData uint32,
	x, y int,
) error {
	if !e.CompositionControllerReady() {
		return errors.New("webview2 composition controller is not initialized")
	}

	return e.compositionController.SendMouseInput(eventKind, virtualKeys, mouseData, POINT{X: int32(x), Y: int32(y)})
}

func (e *Chromium) AllowExternalDrag(allow bool) error {
	if !HasCapability(e.webview2RuntimeVersion, AllowExternalDrop) {
		return UnsupportedCapabilityError
	}
	controller := e.GetController()
	controller4 := controller.GetICoreWebView2Controller4()
	err := controller4.PutAllowExternalDrop(allow)
	if err != nil {
		return err
	}
	return nil
}

func (e *Chromium) GetAllowExternalDrag() (bool, error) {
	if !HasCapability(e.webview2RuntimeVersion, AllowExternalDrop) {
		return false, UnsupportedCapabilityError
	}
	controller := e.GetController()
	controller4 := controller.GetICoreWebView2Controller4()
	result, err := controller4.GetAllowExternalDrop()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (e *Chromium) GetCookieManager() (*ICoreWebView2CookieManager, error) {
	if e.webview == nil {
		return nil, errors.New("webview not initialized")
	}

	// Check WebView2 version
	if e.webview2RuntimeVersion == "" {
		return nil, errors.New("WebView2 runtime version not available")
	}

	// Get ICoreWebView2_2 interface
	webview2, err := e.webview.QueryInterface2()
	if err != nil {
		return nil, fmt.Errorf("failed to get ICoreWebView2_2: %w\nThis functionality requires WebView2 Runtime version 89.0.721.0 or later. Current version: %s", err, e.webview2RuntimeVersion)
	}
	defer webview2.Release()

	// Get cookie manager
	cookieManager, err := webview2.GetCookieManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookie manager: %w", err)
	}

	// Note: The caller is responsible for calling Release() on the returned cookieManager
	return cookieManager, nil
}
