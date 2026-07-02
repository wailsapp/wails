//go:build darwin && purego && !ios && !server

package application

import (
	"sync"
	"time"
	"unsafe"

	"github.com/ebitengine/purego/objc"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// macosApp is the CGO-free implementation of platformApp. It mirrors the field
// layout of the cgo macosApp so the rest of the package is agnostic to which
// backend is compiled in.
type macosApp struct {
	applicationMenu unsafe.Pointer
	parent          *App
}

// The shared application delegate instance and the terminate-behaviour flag
// (stored Go-side rather than as an Objective-C property).
var (
	appDelegate                            id
	appShouldTerminateAfterLastWindowClose bool
	// nsWindowToID maps an NSWindow pointer to its Wails window id. Populated by
	// the window backend on creation and used by getCurrentWindowID.
	nsWindowToID sync.Map // uintptr -> uint
)

func getNativeApplication() *macosApp {
	return globalApplication.impl.(*macosApp)
}

func newPlatformApp(app *App) *macosApp {
	appInit()
	return &macosApp{parent: app}
}

// ---------------------------------------------------------------------------
// Application delegate
// ---------------------------------------------------------------------------

// lifecycleEvent pairs an NSApplicationDelegate notification selector with the
// Wails application event it emits. Generated 1:1 from the cgo delegate.
type lifecycleEvent struct {
	sel string
	ev  events.ApplicationEventType
}

func lifecycleEvents() []lifecycleEvent {
	return []lifecycleEvent{
		{"applicationDidBecomeActive:", events.Mac.ApplicationDidBecomeActive},
		{"applicationDidChangeBackingProperties:", events.Mac.ApplicationDidChangeBackingProperties},
		{"applicationDidChangeEffectiveAppearance:", events.Mac.ApplicationDidChangeEffectiveAppearance},
		{"applicationDidChangeIcon:", events.Mac.ApplicationDidChangeIcon},
		{"applicationDidChangeOcclusionState:", events.Mac.ApplicationDidChangeOcclusionState},
		{"applicationDidChangeScreenParameters:", events.Mac.ApplicationDidChangeScreenParameters},
		{"applicationDidChangeStatusBarFrame:", events.Mac.ApplicationDidChangeStatusBarFrame},
		{"applicationDidChangeStatusBarOrientation:", events.Mac.ApplicationDidChangeStatusBarOrientation},
		{"applicationDidFinishLaunching:", events.Mac.ApplicationDidFinishLaunching},
		{"applicationDidHide:", events.Mac.ApplicationDidHide},
		{"applicationDidResignActive:", events.Mac.ApplicationDidResignActive},
		{"applicationDidUnhide:", events.Mac.ApplicationDidUnhide},
		{"applicationDidUpdate:", events.Mac.ApplicationDidUpdate},
		{"applicationWillBecomeActive:", events.Mac.ApplicationWillBecomeActive},
		{"applicationWillFinishLaunching:", events.Mac.ApplicationWillFinishLaunching},
		{"applicationWillHide:", events.Mac.ApplicationWillHide},
		{"applicationWillResignActive:", events.Mac.ApplicationWillResignActive},
		{"applicationWillTerminate:", events.Mac.ApplicationWillTerminate},
		{"applicationWillUnhide:", events.Mac.ApplicationWillUnhide},
		{"applicationWillUpdate:", events.Mac.ApplicationWillUpdate},
	}
}

var registerAppDelegateOnce sync.Once
var appDelegateClass id

func registerAppDelegateClass() id {
	registerAppDelegateOnce.Do(func() {
		methods := []objc.MethodDef{}

		// Generated lifecycle notifications -> application events.
		for _, le := range lifecycleEvents() {
			ev := le.ev
			methods = append(methods, objc.MethodDef{
				Cmd: sel_(le.sel),
				Fn: func(self objc.ID, cmd objc.SEL, notif objc.ID) {
					pushAppEvent(ev, nil)
				},
			})
		}

		// Custom notification observers (theme + workspace power events).
		custom := []lifecycleEvent{
			{"themeChanged:", events.Mac.ApplicationDidChangeTheme},
			{"workspaceWillSleep:", events.Mac.ApplicationWillSleep},
			{"workspaceDidWake:", events.Mac.ApplicationDidWake},
			{"workspaceScreensDidSleep:", events.Mac.ApplicationScreensDidSleep},
			{"workspaceScreensDidWake:", events.Mac.ApplicationScreensDidWake},
		}
		for _, le := range custom {
			ev := le.ev
			methods = append(methods, objc.MethodDef{
				Cmd: sel_(le.sel),
				Fn: func(self objc.ID, cmd objc.SEL, notif objc.ID) {
					pushAppEvent(ev, nil)
				},
			})
		}

		// application:openFile:
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("application:openFile:"),
			Fn: func(self objc.ID, cmd objc.SEL, app objc.ID, filename objc.ID) bool {
				HandleOpenFile(id(filename).string())
				return true
			},
		})

		// applicationShouldTerminateAfterLastWindowClosed:
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("applicationShouldTerminateAfterLastWindowClosed:"),
			Fn: func(self objc.ID, cmd objc.SEL, app objc.ID) bool {
				return appShouldTerminateAfterLastWindowClose
			},
		})

		// applicationShouldTerminate:
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("applicationShouldTerminate:"),
			Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) int {
				const nsTerminateCancel = 0
				const nsTerminateNow = 1
				if !shouldQuitApplication() {
					return nsTerminateCancel
				}
				cleanup()
				return nsTerminateNow
			},
		})

		// applicationShouldHandleReopen:hasVisibleWindows:
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("applicationShouldHandleReopen:hasVisibleWindows:"),
			Fn: func(self objc.ID, cmd objc.SEL, notif objc.ID, flag bool) bool {
				pushAppEvent(events.Mac.ApplicationShouldHandleReopen, map[string]any{"hasVisibleWindows": flag})
				return true
			},
		})

		// applicationSupportsSecureRestorableState:
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("applicationSupportsSecureRestorableState:"),
			Fn: func(self objc.ID, cmd objc.SEL, app objc.ID) bool {
				return true
			},
		})

		// handleSecondInstanceNotification: (single instance)
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("handleSecondInstanceNotification:"),
			Fn: func(self objc.ID, cmd objc.SEL, note objc.ID) {
				obj := id(note).send("object")
				if !obj.isNil() {
					handleSecondInstanceData(obj.string())
				}
			},
		})

		appDelegateClass = registerDelegateClass("WailsAppDelegate", "NSResponder", nil, methods)
	})
	return appDelegateClass
}

// appInit performs the equivalent of the cgo init(): create the shared
// NSApplication, install the delegate, and register the theme/power observers.
func appInit() {
	loadFrameworks()
	app := class("NSApplication").send("sharedApplication")
	cls := registerAppDelegateClass()
	appDelegate = cls.send("new")
	app.send("setDelegate:", appDelegate)

	// Theme change notification (distributed centre).
	dnc := class("NSDistributedNotificationCenter").send("defaultCenter")
	dnc.send("addObserver:selector:name:object:", appDelegate, sel_("themeChanged:"),
		nsString("AppleInterfaceThemeChangedNotification"), objc.ID(0))

	// Workspace power/screen notifications live on NSWorkspace's own centre.
	wsCenter := class("NSWorkspace").send("sharedWorkspace").send("notificationCenter")
	for sel, name := range map[string]string{
		"workspaceWillSleep:":       "NSWorkspaceWillSleepNotification",
		"workspaceDidWake:":         "NSWorkspaceDidWakeNotification",
		"workspaceScreensDidSleep:": "NSWorkspaceScreensDidSleepNotification",
		"workspaceScreensDidWake:":  "NSWorkspaceScreensDidWakeNotification",
	} {
		wsCenter.send("addObserver:selector:name:object:", appDelegate, sel_(sel),
			nsString(name), objc.ID(0))
	}

	// Local mouse-down monitor: native window dragging from the invisible
	// title-bar strip of frameless / transparent-titlebar windows. Mirrors the
	// cgo handleLeftMouseDown: logic, implemented directly in Go.
	installFramelessDragMonitor()

	startCustomProtocolHandler()
}

func installFramelessDragMonitor() {
	const nsEventMaskLeftMouseDown = 1 << 1
	block := objc.NewBlock(func(b objc.Block, event objc.ID) objc.ID {
		ev := id(event)
		win := ev.send("window")
		if win.isNil() {
			return event
		}
		v, ok := nsWindowToID.Load(win.ptr())
		if !ok {
			return event
		}
		impl := windowImplForID(v.(uint))
		if impl == nil || impl.invisibleTitleBarHeight == 0 {
			return event
		}
		loc := get[NSPoint](ev, "locationInWindow")
		frame := get[NSRect](win, "frame")
		if loc.Y > frame.Size.Height-CGFloat(impl.invisibleTitleBarHeight) {
			// Skip near the left/right edges so native corner resize still works.
			const resizeThreshold = 5.0
			if loc.X < resizeThreshold || loc.X > frame.Size.Width-resizeThreshold {
				return event
			}
			win.send("performWindowDragWithEvent:", event)
		}
		return event
	})
	class("NSEvent").send("addLocalMonitorForEventsMatchingMask:handler:",
		uint(nsEventMaskLeftMouseDown), block)
}

// pushAppEvent constructs and enqueues an application event, mirroring the cgo
// processApplicationEvent post-processing (dark-mode annotation on theme change).
func pushAppEvent(id events.ApplicationEventType, data map[string]any) {
	event := newApplicationEvent(id)
	if data != nil {
		event.Context().setData(data)
	}
	if uint(event.Id) == uint(events.Mac.ApplicationDidChangeTheme) {
		event.Context().setIsDarkMode(globalApplication.Env.IsDarkMode())
	}
	applicationEvents <- event
}

// ---------------------------------------------------------------------------
// platformApp implementation
// ---------------------------------------------------------------------------

func (m *macosApp) isDarkMode() bool {
	ud := class("NSUserDefaults").send("standardUserDefaults")
	if ud.isNil() {
		return false
	}
	style := ud.send("stringForKey:", nsString("AppleInterfaceStyle"))
	if style.isNil() {
		return false
	}
	return get[bool](style, "isEqualToString:", nsString("Dark"))
}

func (m *macosApp) getAccentColor() string {
	var result string
	runOnMain(func() {
		withAutoreleasePool(func() {
			accent := class("NSColor").send("controlAccentColor")
			rgb := accent.send("colorUsingColorSpace:", class("NSColorSpace").send("sRGBColorSpace"))
			if rgb.isNil() {
				rgb = accent
			}
			var r, g, b, a CGFloat
			objc.ID(rgb).Send(sel_("getRed:green:blue:alpha:"),
				unsafe.Pointer(&r), unsafe.Pointer(&g), unsafe.Pointer(&b), unsafe.Pointer(&a))
			result = "rgb(" + itoa(int(r*255)) + "," + itoa(int(g*255)) + "," + itoa(int(b*255)) + ")"
		})
	})
	return result
}

func (m *macosApp) hide() {
	runOnMain(func() { class("NSApplication").send("sharedApplication").send("hide:", objc.ID(0)) })
}
func (m *macosApp) show() {
	runOnMain(func() { class("NSApplication").send("sharedApplication").send("unhide:", objc.ID(0)) })
}

func (m *macosApp) on(eventID uint) {
	// hasListeners() is always true in the cgo backend, so registration is a
	// no-op here; retained for interface parity.
}

func (m *macosApp) setIcon(icon []byte) {
	if len(icon) == 0 {
		return
	}
	runOnMain(func() {
		image := class("NSImage").send("alloc").send("initWithData:", nsData(icon))
		class("NSApplication").send("sharedApplication").send("setApplicationIconImage:", image)
	})
}

func (m *macosApp) name() string {
	running := class("NSRunningApplication").send("currentApplication").send("localizedName")
	if running.isNil() {
		return class("NSProcessInfo").send("processInfo").send("processName").string()
	}
	return running.string()
}

func (m *macosApp) getCurrentWindowID() uint {
	var result uint
	runOnMain(func() {
		app := class("NSApplication").send("sharedApplication")
		win := app.send("keyWindow")
		if win.isNil() {
			win = app.send("mainWindow")
		}
		if win.isNil() {
			return
		}
		if v, ok := nsWindowToID.Load(win.ptr()); ok {
			result = v.(uint)
		}
	})
	return result
}

func (m *macosApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		menu = DefaultApplicationMenu()
	}
	menu.Update()
	m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	runOnMain(func() {
		class("NSApplication").send("sharedApplication").
			send("setMainMenu:", id(uintptr(m.applicationMenu)))
	})
}

func (m *macosApp) run() error {
	if m.parent.options.SingleInstance != nil {
		startSingleInstanceListener(m.parent.options.SingleInstance.UniqueID)
	}
	m.parent.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(*ApplicationEvent) {
		appShouldTerminateAfterLastWindowClose = m.parent.options.Mac.ApplicationShouldTerminateAfterLastWindowClosed
		app := class("NSApplication").send("sharedApplication")
		app.send("setActivationPolicy:", int(m.parent.options.Mac.ActivationPolicy))
		app.send("activateIgnoringOtherApps:", true)
		if err := m.processAndCacheScreens(); err != nil {
			m.parent.handleError(err)
		}
	})
	m.parent.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeScreenParameters, func(*ApplicationEvent) {
		if err := m.processAndCacheScreens(); err != nil {
			m.parent.handleError(err)
		}
	})
	m.setupCommonEvents()
	for eventID := range m.parent.applicationEventListeners {
		m.on(eventID)
	}
	class("NSApplication").send("sharedApplication").send("run")
	return nil
}

func (m *macosApp) destroy() {
	runOnMain(func() { class("NSApplication").send("sharedApplication").send("terminate:", objc.ID(0)) })
}

func (m *macosApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

// startSingleInstanceListener registers this app to receive second-instance
// notifications for the given unique id (mirrors the cgo helper of the same
// role).
func startSingleInstanceListener(uniqueID string) {
	class("NSDistributedNotificationCenter").send("defaultCenter").
		send("addObserver:selector:name:object:", appDelegate,
			sel_("handleSecondInstanceNotification:"), nsString(uniqueID), objc.ID(0))
}

// ---------------------------------------------------------------------------
// Custom URL scheme (Apple Event) handler
// ---------------------------------------------------------------------------

var registerProtocolOnce sync.Once

func startCustomProtocolHandler() {
	registerProtocolOnce.Do(func() {
		cls := registerDelegateClass("WailsProtocolHandler", "NSObject", nil, []objc.MethodDef{{
			Cmd: sel_("handleGetURLEvent:withReplyEvent:"),
			Fn: func(self objc.ID, cmd objc.SEL, event objc.ID, reply objc.ID) {
				const keyDirectObject = 0x2d2d2d2d // '----'
				desc := id(event).send("paramDescriptorForKeyword:", uint(keyDirectObject))
				urlStr := desc.send("stringValue")
				if !urlStr.isNil() {
					HandleOpenURL(urlStr.string())
				}
			},
		}})
		handler := cls.send("new")
		const kInternetEventClass = 0x4755524c // 'GURL'
		const kAEGetURL = 0x4755524c           // 'GURL'
		class("NSAppleEventManager").send("sharedAppleEventManager").
			send("setEventHandler:andSelector:forEventClass:andEventID:",
				handler, sel_("handleGetURLEvent:withReplyEvent:"),
				uint(kInternetEventClass), uint(kAEGetURL))
	})
}

// ---------------------------------------------------------------------------
// Event / message plumbing (Go-native ports of the cgo //export callbacks)
// ---------------------------------------------------------------------------

func processWindowEvent(windowID uint, eventID uint) {
	windowEvents <- &windowEvent{WindowID: windowID, EventID: eventID}
}

func processMessage(windowID uint, message string, origin string, isMainFrame bool) {
	windowMessageBuffer <- &windowMessage{
		windowId:   windowID,
		message:    message,
		originInfo: &OriginInfo{Origin: origin, IsMainFrame: isMainFrame},
	}
}

func processURLRequest(windowID uint, wkURLSchemeTask unsafe.Pointer) {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		globalApplication.debug("could not find window with id", "windowID", windowID)
		return
	}
	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(wkURLSchemeTask),
		windowId:   windowID,
		windowName: window.Name(),
	}
}

func processWindowKeyDownEvent(windowID uint, acceleratorString string) {
	windowKeyEvents <- &windowKeyEvent{windowId: windowID, acceleratorString: acceleratorString}
}

func shouldQuitApplication() bool { return globalApplication.shouldQuit() }

func cleanup() { globalApplication.cleanup() }

func HandleOpenFile(goFilepath string) {
	eventContext := newApplicationEventContext()
	eventContext.setOpenedWithFile(goFilepath)
	applicationEvents <- &ApplicationEvent{Id: uint(events.Common.ApplicationOpenedWithFile), ctx: eventContext}
}

func HandleOpenURL(urlString string) {
	eventContext := newApplicationEventContext()
	eventContext.setURL(urlString)
	applicationEvents <- &ApplicationEvent{Id: uint(events.Common.ApplicationLaunchedWithUrl), ctx: eventContext}
}

func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("error getting OS info: %w", err)
		return
	}
	a.info("Platform Info:", info.AsLogSlice()...)
}

func (a *App) platformEnvironment() map[string]any { return map[string]any{} }

func fatalHandler(errFunc func(error)) {}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

// ---------------------------------------------------------------------------
// Drag-and-drop throttling (verbatim port of the pure-Go cgo logic)
// ---------------------------------------------------------------------------

var (
	dragOverJSBuffer = make([]byte, 128)
	dragOverJSMutex  sync.Mutex
	dragOverJSPrefix = []byte("window._wails.handleDragOver(")
	windowImplCache  sync.Map // windowID -> *macosWebviewWindow
	dragThrottle     sync.Map // windowID -> *dragThrottleState
)

type dragThrottleState struct {
	mu           sync.Mutex
	lastX, lastY int
	timer        *time.Timer
	pendingX     int
	pendingY     int
	hasPending   bool
}

func clearWindowDragCache(windowID uint) {
	windowImplCache.Delete(windowID)
	if throttleVal, ok := dragThrottle.Load(windowID); ok {
		if throttle, ok := throttleVal.(*dragThrottleState); ok {
			throttle.mu.Lock()
			if throttle.timer != nil {
				throttle.timer.Stop()
			}
			throttle.mu.Unlock()
		}
	}
	dragThrottle.Delete(windowID)
}

func writeInt(buf []byte, n int) int {
	if n < 0 {
		if len(buf) == 0 {
			return 0
		}
		buf[0] = '-'
		return 1 + writeInt(buf[1:], -n)
	}
	if n == 0 {
		if len(buf) == 0 {
			return 0
		}
		buf[0] = '0'
		return 1
	}
	tmp := n
	digits := 0
	for tmp > 0 {
		digits++
		tmp /= 10
	}
	if digits > len(buf) {
		return 0
	}
	for i := digits - 1; i >= 0; i-- {
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return digits
}

func macosOnDragEnter(windowID uint) {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return
	}
	window.ExecJS("window._wails.handleDragEnter();")
}

func macosOnDragExit(windowID uint) {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return
	}
	window.ExecJS("window._wails.handleDragLeave();")
}

func macosOnDragOver(windowID uint, x int, y int) {
	winID := windowID
	intX, intY := x, y
	throttleVal, _ := dragThrottle.LoadOrStore(winID, &dragThrottleState{lastX: intX, lastY: intY})
	throttle := throttleVal.(*dragThrottleState)
	throttle.mu.Lock()
	throttle.pendingX = intX
	throttle.pendingY = intY
	throttle.hasPending = true
	if throttle.timer != nil {
		throttle.mu.Unlock()
		return
	}
	dx := intX - throttle.lastX
	dy := intY - throttle.lastY
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	if dx >= 5 || dy >= 5 {
		throttle.lastX = intX
		throttle.lastY = intY
		throttle.hasPending = false
		throttle.mu.Unlock()
		sendDragUpdate(winID, intX, intY)
		throttle.mu.Lock()
	}
	throttle.timer = time.AfterFunc(50*time.Millisecond, func() {
		InvokeSync(func() {
			throttle.mu.Lock()
			throttle.timer = nil
			if throttle.hasPending {
				pendingX, pendingY := throttle.pendingX, throttle.pendingY
				throttle.lastX = pendingX
				throttle.lastY = pendingY
				throttle.hasPending = false
				throttle.mu.Unlock()
				sendDragUpdate(winID, pendingX, pendingY)
			} else {
				throttle.mu.Unlock()
			}
		})
	})
	throttle.mu.Unlock()
}

func sendDragUpdate(winID uint, x, y int) {
	var darwinImpl *macosWebviewWindow
	var needsExecJS bool
	if cached, found := windowImplCache.Load(winID); found {
		darwinImpl = cached.(*macosWebviewWindow)
		if darwinImpl != nil && darwinImpl.nsWindow != nil {
			needsExecJS = true
		} else {
			windowImplCache.Delete(winID)
		}
	}
	if !needsExecJS {
		window, ok := globalApplication.Window.GetByID(winID)
		if !ok || window == nil {
			return
		}
		webviewWindow, ok := window.(*WebviewWindow)
		if !ok || webviewWindow == nil {
			return
		}
		darwinImpl, ok = webviewWindow.impl.(*macosWebviewWindow)
		if !ok {
			return
		}
		windowImplCache.Store(winID, darwinImpl)
		needsExecJS = true
	}
	if !needsExecJS || darwinImpl == nil {
		return
	}
	dragOverJSMutex.Lock()
	n := copy(dragOverJSBuffer[:], dragOverJSPrefix)
	n += writeInt(dragOverJSBuffer[n:], x)
	if n < len(dragOverJSBuffer) {
		dragOverJSBuffer[n] = ','
		n++
	}
	n += writeInt(dragOverJSBuffer[n:], y)
	if n < len(dragOverJSBuffer) {
		dragOverJSBuffer[n] = ')'
		n++
	}
	if n < len(dragOverJSBuffer) {
		dragOverJSBuffer[n] = 0
	} else {
		dragOverJSMutex.Unlock()
		return
	}
	darwinImpl.execJSDragOver(dragOverJSBuffer[:n+1])
	dragOverJSMutex.Unlock()
}

func processDragItems(windowID uint, filenames []string, x, y int) {
	targetWindow, ok := globalApplication.Window.GetByID(windowID)
	if !ok || targetWindow == nil {
		return
	}
	targetWindow.InitiateFrontendDropProcessing(filenames, x, y)
}
