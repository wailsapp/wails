//go:build darwin && purego && !ios && !server

package application

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/ebitengine/purego/objc"
	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/runtime"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// macosWebviewWindow is the CGO-free implementation of webviewWindowImpl. Field
// names match the cgo backend (nsWindow) so cross-file references keep working;
// additional handles are kept for the webview and its delegate.
type macosWebviewWindow struct {
	nsWindow  unsafe.Pointer // NSWindow*
	parent    *WebviewWindow
	wkWebView unsafe.Pointer // WKWebView*
	delegate  unsafe.Pointer // WailsWebviewWindowDelegate*

	// invisibleTitleBarHeight enables native window dragging from the top
	// `height` points of a frameless / transparent-titlebar window (see the
	// local mouse monitor installed in appInit).
	invisibleTitleBarHeight uint
	// showToolbarWhenFullscreen controls the fullscreen presentation options
	// returned by the delegate.
	showToolbarWhenFullscreen bool
}

// windowImplForID resolves the macOS window backend for a Wails window id.
func windowImplForID(windowID uint) *macosWebviewWindow {
	if cached, ok := windowImplCache.Load(windowID); ok {
		if impl, ok := cached.(*macosWebviewWindow); ok {
			return impl
		}
	}
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return nil
	}
	ww, ok := window.(*WebviewWindow)
	if !ok || ww == nil {
		return nil
	}
	impl, ok := ww.impl.(*macosWebviewWindow)
	if !ok {
		return nil
	}
	windowImplCache.Store(windowID, impl)
	return impl
}

func (w *macosWebviewWindow) win() id     { return id(uintptr(w.nsWindow)) }
func (w *macosWebviewWindow) webview() id { return id(uintptr(w.wkWebView)) }

// ---------------------------------------------------------------------------
// Window/webview delegate class
// ---------------------------------------------------------------------------

var (
	delegateToWindowID sync.Map // uintptr(delegate) -> uint
	macEventsVal       = reflect.ValueOf(events.Mac)
)

// macWindowEventID resolves an events.Mac field (e.g. "WindowDidResize") to its
// numeric id, so the delegate notification table can be generated from the
// selector names rather than hand-copying ~60 constants.
func macWindowEventID(field string) uint {
	f := macEventsVal.FieldByName(field)
	if !f.IsValid() {
		panic("wails/purego: unknown events.Mac field " + field)
	}
	return uint(f.Uint())
}

func windowIDForDelegate(self objc.ID) (uint, bool) {
	if v, ok := delegateToWindowID.Load(uintptr(self)); ok {
		return v.(uint), true
	}
	return 0, false
}

// windowNotificationSelectors are the NSWindowDelegate notifications that map
// 1:1 to an events.Mac.<Field> event (field name = selector without the
// trailing colon, first letter upper-cased).
var windowNotificationSelectors = []string{
	"windowDidBecomeKey:", "windowDidBecomeMain:", "windowDidBeginSheet:",
	"windowDidChangeAlpha:", "windowDidChangeBackingLocation:", "windowDidChangeBackingProperties:",
	"windowDidChangeCollectionBehavior:", "windowDidChangeEffectiveAppearance:", "windowDidChangeOrderingMode:",
	"windowDidChangeScreen:", "windowDidChangeScreenParameters:", "windowDidChangeScreenProfile:",
	"windowDidChangeScreenSpace:", "windowDidChangeScreenSpaceProperties:", "windowDidChangeSharingType:",
	"windowDidChangeSpace:", "windowDidChangeSpaceOrderingMode:", "windowDidChangeTitle:",
	"windowDidChangeToolbar:", "windowDidDeminiaturize:", "windowDidEndSheet:",
	"windowDidEnterFullScreen:", "windowDidEnterVersionBrowser:", "windowDidExitFullScreen:",
	"windowDidExitVersionBrowser:", "windowDidExpose:", "windowDidFocus:",
	"windowDidMiniaturize:", "windowDidMove:", "windowDidOrderOffScreen:",
	"windowDidOrderOnScreen:", "windowDidResignKey:", "windowDidResignMain:",
	"windowDidResize:", "windowDidUpdate:", "windowDidUpdateAlpha:",
	"windowDidUpdateCollectionBehavior:", "windowDidUpdateCollectionProperties:", "windowDidUpdateShadow:",
	"windowDidUpdateTitle:", "windowDidUpdateToolbar:", "windowWillBecomeKey:",
	"windowWillBecomeMain:", "windowWillBeginSheet:", "windowWillChangeOrderingMode:",
	"windowWillDeminiaturize:", "windowWillEnterFullScreen:", "windowWillEnterVersionBrowser:",
	"windowWillExitFullScreen:", "windowWillExitVersionBrowser:", "windowWillFocus:",
	"windowWillMiniaturize:", "windowWillMove:",
}

func selectorToEventField(sel string) string {
	s := sel[:len(sel)-1] // drop trailing ':'
	return string(s[0]-'a'+'A') + s[1:]
}

var registerWindowDelegateOnce sync.Once
var windowDelegateClass id

func registerWindowDelegateClass() id {
	registerWindowDelegateOnce.Do(func() {
		methods := []objc.MethodDef{}

		// Generic notifications -> window events.
		for _, sel := range windowNotificationSelectors {
			evID := macWindowEventID(selectorToEventField(sel))
			methods = append(methods, objc.MethodDef{
				Cmd: sel_(sel),
				Fn: func(self objc.ID, cmd objc.SEL, notif objc.ID) {
					if wid, ok := windowIDForDelegate(self); ok {
						processWindowEvent(wid, evID)
					}
				},
			})
		}

		// Occlusion state -> Show/Hide.
		showID := macWindowEventID("WindowShow")
		hideID := macWindowEventID("WindowHide")
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("windowDidChangeOcclusionState:"),
			Fn: func(self objc.ID, cmd objc.SEL, notif objc.ID) {
				wid, ok := windowIDForDelegate(self)
				if !ok {
					return
				}
				win := id(notif).send("object")
				const nsWindowOcclusionStateVisible = 1 << 1
				state := get[uint](win, "occlusionState")
				if state&nsWindowOcclusionStateVisible != 0 {
					processWindowEvent(wid, showID)
				} else {
					processWindowEvent(wid, hideID)
				}
			},
		})

		// windowShouldClose: honour hidden / unconditional-close semantics.
		shouldCloseID := macWindowEventID("WindowShouldClose")
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("windowShouldClose:"),
			Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) bool {
				wid, ok := windowIDForDelegate(self)
				if !ok {
					return true
				}
				if windowShouldUnconditionallyClose(wid) {
					return true
				}
				if windowIsHidden(wid) {
					return false
				}
				processWindowEvent(wid, shouldCloseID)
				return false
			},
		})

		// WKScriptMessageHandler: userContentController:didReceiveScriptMessage:
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("userContentController:didReceiveScriptMessage:"),
			Fn: func(self objc.ID, cmd objc.SEL, ucc objc.ID, message objc.ID) {
				wid, ok := windowIDForDelegate(self)
				if !ok {
					return
				}
				msg := id(message)
				origin := ""
				frame := msg.send("frameInfo")
				if !frame.isNil() {
					req := frame.send("request")
					if !req.isNil() {
						u := req.send("URL")
						if !u.isNil() && !u.send("scheme").isNil() && !u.send("host").isNil() {
							origin = u.send("absoluteString").string()
						}
					}
				}
				body := msg.send("body")
				var bodyStr string
				if get[bool](body, "isKindOfClass:", class("NSString")) {
					bodyStr = body.string()
				} else {
					bodyStr = body.send("description").string()
				}
				isMain := get[bool](frame, "isMainFrame")
				processMessage(wid, bodyStr, origin, isMain)
			},
		})

		// WKURLSchemeHandler: webView:startURLSchemeTask:
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("webView:startURLSchemeTask:"),
			Fn: func(self objc.ID, cmd objc.SEL, wv objc.ID, task objc.ID) {
				if wid, ok := windowIDForDelegate(self); ok {
					processURLRequest(wid, unsafe.Pointer(uintptr(task)))
				}
			},
		})
		// WKURLSchemeHandler: webView:stopURLSchemeTask: — mark the task stopped so
		// the asset writer skips further messaging (avoids WebKit's NSException).
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("webView:stopURLSchemeTask:"),
			Fn: func(self objc.ID, cmd objc.SEL, wv objc.ID, task objc.ID) {
				webview.MarkTaskStopped(unsafe.Pointer(uintptr(task)))
			},
		})

		// WKNavigationDelegate: didFinishNavigation.
		finishNavID := macWindowEventID("WebViewDidFinishNavigation")
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("webView:didFinishNavigation:"),
			Fn: func(self objc.ID, cmd objc.SEL, wv objc.ID, nav objc.ID) {
				if wid, ok := windowIDForDelegate(self); ok {
					processWindowEvent(wid, finishNavID)
				}
			},
		})

		// window:willUseFullScreenPresentationOptions: — hide the toolbar in
		// fullscreen unless the window opted to keep it.
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("window:willUseFullScreenPresentationOptions:"),
			Fn: func(self objc.ID, cmd objc.SEL, window objc.ID, proposed uint) uint {
				const autoHideToolbar = 1 << 11 // NSApplicationPresentationAutoHideToolbar
				if wid, ok := windowIDForDelegate(self); ok {
					if impl := windowImplForID(wid); impl != nil && impl.showToolbarWhenFullscreen {
						return proposed
					}
				}
				return proposed | autoHideToolbar
			},
		})

		windowDelegateClass = registerDelegateClass("WailsWebviewWindowDelegate", "NSObject", nil, methods)
	})
	return windowDelegateClass
}

// ---------------------------------------------------------------------------
// Construction
// ---------------------------------------------------------------------------

func newWindowImpl(parent *WebviewWindow) *macosWebviewWindow {
	result := &macosWebviewWindow{parent: parent}
	result.parent.RegisterHook(events.Mac.WebViewDidFinishNavigation, func(event *WindowEvent) {
		js := runtime.Core(globalApplication.impl.GetFlags(globalApplication.options))
		js += fmt.Sprintf("window._wails.flags.enableFileDrop=%v;", result.parent.options.EnableFileDrop)
		result.execJS(js)
	})
	return result
}

func (w *macosWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}
	globalApplication.dispatchOnMainThread(func() {
		options := w.parent.options
		w.createWindow(options)

		w.setTitle(options.Title)
		w.setResizable(!options.DisableResize)
		if options.MinWidth != 0 || options.MinHeight != 0 {
			w.setMinSize(options.MinWidth, options.MinHeight)
		}
		if options.MaxWidth != 0 || options.MaxHeight != 0 {
			w.setMaxSize(options.MaxWidth, options.MaxHeight)
		}
		w.enableDevTools()
		w.setContentProtection(options.ContentProtectionEnabled)
		w.setBackgroundColour(options.BackgroundColour)
		w.applyMacOptions(options)

		switch options.StartState {
		case WindowStateMaximised:
			w.maximise()
		case WindowStateMinimised:
			w.minimise()
		case WindowStateFullscreen:
			w.fullscreen()
		case WindowStateNormal:
		}

		if options.InitialPosition == WindowCentered {
			w.center()
		} else {
			w.setPosition(options.X, options.Y)
		}

		startURL, err := assetserver.GetStartURL(options.URL)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
		w.setURL(startURL)

		w.parent.OnWindowEvent(events.Mac.WebViewDidFinishNavigation, func(_ *WindowEvent) {
			InvokeAsync(func() {
				if options.JS != "" {
					w.execJS(options.JS)
				}
				if options.CSS != "" {
					runOnMain(func() { w.injectCSS(options.CSS) })
				}
				if !options.Hidden {
					w.parent.Show()
					w.setHasShadow(!options.Mac.DisableShadow)
					w.setAlwaysOnTop(options.AlwaysOnTop)
				}
			})
		})

		if options.HTML != "" {
			w.setHTML(options.HTML)
		}
	})
}

// applyMacOptions applies the macOS-specific window options, mirroring the tail
// of the cgo run(): backdrop, window level, collection behavior, window
// buttons, ignore-mouse, and (for non-frameless windows) titlebar presets and
// appearance.
func (w *macosWebviewWindow) applyMacOptions(options WebviewWindowOptions) {
	macOptions := options.Mac

	if macOptions.DisableEscapeExitsFullscreen {
		windowDisableEscape.Store(w.win().ptr(), true)
	}

	switch macOptions.Backdrop {
	case MacBackdropTransparent:
		w.setTransparentBackdrop()
		w.setWebviewTransparent()
	case MacBackdropTranslucent:
		w.setTranslucentBackdrop()
		w.setWebviewTransparent()
	case MacBackdropLiquidGlass:
		w.applyLiquidGlass()
	case MacBackdropNormal:
	}

	level := macOptions.WindowLevel
	if level == "" {
		level = MacWindowLevelNormal
	}
	w.setWindowLevel(level)
	w.setCollectionBehavior(macOptions.CollectionBehavior)

	// Window buttons. Maximise and Fullscreen both drive the zoom button, so
	// apply the more restrictive of the two.
	w.setMinimiseButtonState(options.MinimiseButtonState)
	w.setCloseButtonState(options.CloseButtonState)
	zoomState := options.MaximiseButtonState
	if options.FullscreenButtonState > zoomState {
		zoomState = options.FullscreenButtonState
	}
	w.setMaximiseButtonState(zoomState)

	w.setIgnoreMouseEvents(options.IgnoreMouseEvents)

	titleBar := macOptions.TitleBar
	w.showToolbarWhenFullscreen = titleBar.ShowToolbarWhenFullscreen
	if !options.Frameless {
		w.setTitleBarAppearsTransparent(titleBar.AppearsTransparent)
		w.setHideTitleBar(titleBar.Hide)
		w.setTitleVisibility(titleBar.HideTitle)
		w.setFullSizeContent(titleBar.FullSizeContent)
		w.setUseToolbar(titleBar.UseToolbar)
		w.setToolbarStyle(int(titleBar.ToolbarStyle))
		w.setHideToolbarSeparator(titleBar.HideToolbarSeparator)
	}

	// Enable native drag from the invisible title-bar strip when configured for
	// a frameless or transparent-titlebar window (matches cgo run()).
	if macOptions.InvisibleTitleBarHeight != 0 && (options.Frameless || titleBar.AppearsTransparent) {
		w.invisibleTitleBarHeight = uint(macOptions.InvisibleTitleBarHeight)
	}

	if macOptions.Appearance != "" {
		w.setAppearanceByName(string(macOptions.Appearance))
	}
}

// createWindow builds the NSWindow, WKWebView, configuration, delegate and
// scheme/message handlers. Port of the cgo windowNew().
func (w *macosWebviewWindow) createWindow(options WebviewWindowOptions) {
	const (
		styleTitled      = 1 << 0
		styleClosable    = 1 << 1
		styleMiniaturize = 1 << 2
		styleResizable   = 1 << 3
		styleBorderless  = 0
		backingBuffered  = 2
	)
	styleMask := styleTitled | styleClosable | styleMiniaturize | styleResizable
	if options.Frameless {
		styleMask = styleBorderless | styleResizable | styleMiniaturize
	}
	width, height := options.Width, options.Height
	if width == 0 {
		width = 800
	}
	if height == 0 {
		height = 600
	}

	win := registerWebviewWindowClass().send("alloc").send("initWithContentRect:styleMask:backing:defer:",
		rect(0, 0, CGFloat(width-1), CGFloat(height-1)), uint(styleMask), uint(backingBuffered), false)
	w.nsWindow = unsafe.Pointer(win.ptr())

	// Delegate (also the script/scheme/navigation handler).
	del := registerWindowDelegateClass().send("new")
	w.delegate = unsafe.Pointer(del.ptr())
	win.send("setDelegate:", del)
	delegateToWindowID.Store(del.ptr(), w.parent.id)
	nsWindowToID.Store(win.ptr(), w.parent.id)

	// Content view.
	view := class("NSView").send("alloc").send("initWithFrame:", rect(0, 0, CGFloat(width-1), CGFloat(height-1)))
	const autoWidth, autoHeight = 1 << 1, 1 << 4
	view.send("setAutoresizingMask:", uint(autoWidth|autoHeight))
	win.send("setContentView:", view)

	// WebView configuration.
	config := class("WKWebViewConfiguration").send("alloc").send("init")
	config.send("setSuppressesIncrementalRendering:", true)
	appName := options.Mac.WebviewPreferences.ApplicationNameForUserAgent
	if appName == "" {
		appName = "wails.io"
	}
	config.send("setApplicationNameForUserAgent:", nsString(appName))
	config.send("setURLSchemeHandler:forURLScheme:", del, nsString("wails"))

	// User content controller + external message bridge.
	ucc := class("WKUserContentController").send("new")
	ucc.send("addScriptMessageHandler:name:", del, nsString("external"))
	config.send("setUserContentController:", ucc)

	webView := class("WKWebView").send("alloc").send("initWithFrame:configuration:",
		rect(0, 0, CGFloat(width), CGFloat(height)), config)
	w.wkWebView = unsafe.Pointer(webView.ptr())
	view.send("addSubview:", webView)
	webView.send("setNavigationDelegate:", del)
	webView.send("setUIDelegate:", del)
	webView.send("setAutoresizingMask:", uint(autoWidth|autoHeight))

	applyWebviewPreferences(webView, config, options.Mac.WebviewPreferences, options.Mac.EnableFraudulentWebsiteWarnings)

	if options.EnableFileDrop {
		w.installFileDropView(view, width, height)
	}
}

// ---------------------------------------------------------------------------
// NSWindow subclass
//
// A minimal NSWindow subclass supplying the overrides the cgo WebviewWindow
// relies on: it must be able to become key/main even when borderless
// (frameless), and it can optionally swallow the Escape key so it doesn't exit
// fullscreen.
// ---------------------------------------------------------------------------

var (
	registerWindowClassOnce sync.Once
	webviewWindowClass      id
	windowDisableEscape     sync.Map // uintptr(window) -> bool
)

// keyCodeToString maps macOS virtual key codes to the Wails accelerator key
// names, mirroring the cgo keyStringFromEvent switch (US layout).
var keyCodeToString = map[uint]string{
	122: "f1", 120: "f2", 99: "f3", 118: "f4", 96: "f5", 97: "f6", 98: "f7",
	100: "f8", 101: "f9", 109: "f10", 103: "f11", 111: "f12", 105: "f13",
	107: "f14", 113: "f15", 106: "f16", 64: "f17", 79: "f18", 80: "f19", 90: "f20",
	0: "a", 11: "b", 8: "c", 2: "d", 14: "e", 3: "f", 5: "g", 4: "h", 34: "i",
	38: "j", 40: "k", 37: "l", 46: "m", 45: "n", 31: "o", 35: "p", 12: "q",
	15: "r", 1: "s", 17: "t", 32: "u", 9: "v", 13: "w", 7: "x", 16: "y", 6: "z",
	29: "0", 18: "1", 19: "2", 20: "3", 21: "4", 23: "5", 22: "6", 26: "7",
	28: "8", 25: "9",
	51: "delete", 117: "forward delete", 123: "left", 124: "right", 126: "up",
	125: "down", 48: "tab", 53: "escape", 49: "space",
	33: "[", 30: "]", 43: ",", 27: "-", 39: "'", 44: "/", 47: ".", 41: ";",
	24: "=", 50: "`", 42: "\\",
}

// keyStringFromEvent replicates the cgo keyStringFromEvent: mapping.
func keyStringFromEvent(ev id) string {
	characters := ev.send("characters").string()
	if characters == "" {
		return ""
	}
	switch characters {
	case "\r":
		return "enter"
	case "\b":
		return "backspace"
	case "\x1b":
		return "escape"
	case "\x0b":
		return "page down"
	case "\x0e":
		return "page up"
	case "\x01":
		return "home"
	case "\x04":
		return "end"
	case "\x0c":
		return "clear"
	}
	return keyCodeToString[get[uint](ev, "keyCode")]
}

func registerWebviewWindowClass() id {
	registerWindowClassOnce.Do(func() {
		yes := func(self objc.ID, cmd objc.SEL) bool { return true }
		methods := []objc.MethodDef{
			{Cmd: sel_("canBecomeKeyWindow"), Fn: yes},
			{Cmd: sel_("canBecomeMainWindow"), Fn: yes},
			{Cmd: sel_("acceptsFirstResponder"), Fn: yes},
			{Cmd: sel_("becomeFirstResponder"), Fn: yes},
			{Cmd: sel_("resignFirstResponder"), Fn: yes},
			{Cmd: sel_("cancelOperation:"), Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) {
				const fullScreen = 1 << 14 // NSWindowStyleMaskFullScreen
				if v, ok := windowDisableEscape.Load(uintptr(self)); ok && v.(bool) {
					if get[uint](id(self), "styleMask")&fullScreen == fullScreen {
						return
					}
				}
				objc.ID(self).SendSuper(sel_("cancelOperation:"), sender)
			}},
			{Cmd: sel_("keyDown:"), Fn: func(self objc.ID, cmd objc.SEL, event objc.ID) {
				v, ok := nsWindowToID.Load(uintptr(self))
				if !ok {
					return
				}
				const (
					flagShift   = 1 << 17
					flagControl = 1 << 18
					flagOption  = 1 << 19
					flagCommand = 1 << 20
				)
				ev := id(event)
				mods := get[uint](ev, "modifierFlags")
				parts := make([]string, 0, 5)
				if mods&flagShift != 0 {
					parts = append(parts, "shift")
				}
				if mods&flagControl != 0 {
					parts = append(parts, "ctrl")
				}
				if mods&flagOption != 0 {
					parts = append(parts, "option")
				}
				if mods&flagCommand != 0 {
					parts = append(parts, "cmd")
				}
				if key := keyStringFromEvent(ev); key != "" {
					parts = append(parts, key)
				}
				processWindowKeyDownEvent(v.(uint), strings.Join(parts, "+"))
			}},
		}
		webviewWindowClass = registerDelegateClass("WailsWebviewNSWindow", "NSWindow", nil, methods)
	})
	return webviewWindowClass
}

// ---------------------------------------------------------------------------
// File-drop overlay (NSView <NSDraggingDestination>)
// ---------------------------------------------------------------------------

const nsFilenamesPboardType = "NSFilenamesPboardType"

var (
	dragViewToWindowID   sync.Map // uintptr(view) -> uint
	registerDragViewOnce sync.Once
	dragViewClass        id
	nsDragOperationNone  = uint(0)
	nsDragOperationCopy  = uint(1)
)

// dropPointToContentXY converts a dragging-info location into content-view
// top-left coordinates, matching the cgo WebviewDrag conversion.
func dropPointToContentXY(self, sender id) (int, int) {
	loc := get[NSPoint](sender, "draggingLocation")
	inView := get[NSPoint](self, "convertPoint:fromView:", loc, objc.ID(0))
	contentView := self.send("window").send("contentView")
	contentHeight := get[NSRect](contentView, "frame").Size.Height
	return int(inView.X), int(contentHeight - inView.Y)
}

func pasteboardHasFiles(sender id) bool {
	pb := sender.send("draggingPasteboard")
	return get[bool](pb.send("types"), "containsObject:", nsString(nsFilenamesPboardType))
}

func registerDragViewClass() id {
	registerDragViewOnce.Do(func() {
		enteredID := macWindowEventID("WindowFileDraggingEntered")
		exitedID := macWindowEventID("WindowFileDraggingExited")
		performedID := macWindowEventID("WindowFileDraggingPerformed")

		widFor := func(self objc.ID) (uint, bool) {
			if v, ok := dragViewToWindowID.Load(uintptr(self)); ok {
				return v.(uint), true
			}
			return 0, false
		}

		methods := []objc.MethodDef{
			{Cmd: sel_("draggingEntered:"), Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) uint {
				wid, ok := widFor(self)
				if ok && pasteboardHasFiles(id(sender)) {
					processWindowEvent(wid, enteredID)
					macosOnDragEnter(wid)
					return nsDragOperationCopy
				}
				return nsDragOperationNone
			}},
			{Cmd: sel_("draggingUpdated:"), Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) uint {
				wid, ok := widFor(self)
				if ok && pasteboardHasFiles(id(sender)) {
					x, y := dropPointToContentXY(id(self), id(sender))
					macosOnDragOver(wid, x, y)
					return nsDragOperationCopy
				}
				return nsDragOperationNone
			}},
			{Cmd: sel_("draggingExited:"), Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) {
				if wid, ok := widFor(self); ok {
					processWindowEvent(wid, exitedID)
					macosOnDragExit(wid)
				}
			}},
			{Cmd: sel_("prepareForDragOperation:"), Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) bool {
				return true
			}},
			{Cmd: sel_("performDragOperation:"), Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) bool {
				wid, ok := widFor(self)
				if !ok {
					return false
				}
				processWindowEvent(wid, performedID)
				s := id(sender)
				if !pasteboardHasFiles(s) {
					return false
				}
				files := s.send("draggingPasteboard").send("propertyListForType:", nsString(nsFilenamesPboardType))
				count := get[uint](files, "count")
				if count == 0 {
					return false
				}
				filenames := make([]string, 0, count)
				for i := uint(0); i < count; i++ {
					filenames = append(filenames, files.send("objectAtIndex:", i).string())
				}
				x, y := dropPointToContentXY(id(self), s)
				processDragItems(wid, filenames, x, y)
				return true
			}},
		}
		dragViewClass = registerDelegateClass("WailsWebviewDrag", "NSView", nil, methods)
	})
	return dragViewClass
}

func (w *macosWebviewWindow) installFileDropView(contentView id, width, height int) {
	const autoWidth, autoHeight = 1 << 1, 1 << 4
	dragView := registerDragViewClass().send("alloc").
		send("initWithFrame:", rect(0, 0, CGFloat(width-1), CGFloat(height-1)))
	dragView.send("setAutoresizingMask:", uint(autoWidth|autoHeight))
	dragView.send("registerForDraggedTypes:",
		class("NSArray").send("arrayWithObject:", nsString(nsFilenamesPboardType)))
	dragViewToWindowID.Store(dragView.ptr(), w.parent.id)
	contentView.send("addSubview:", dragView)
}

func applyWebviewPreferences(webView, config id, prefs MacWebviewPreferences, fraudulentWarnings bool) {
	wkPrefs := config.send("preferences")

	// Version-gated preferences. In cgo these were #if/@available gates; here we
	// feature-detect with respondsToSelector: so a build simply skips a setter
	// the running OS's WebKit does not have (calling a missing selector would
	// raise an uncatchable NSException).
	if prefs.TabFocusesLinks.IsSet() && respondsTo(wkPrefs, "setTabFocusesLinks:") {
		wkPrefs.send("setTabFocusesLinks:", prefs.TabFocusesLinks.Get())
	}
	if prefs.TextInteractionEnabled.IsSet() && respondsTo(wkPrefs, "setTextInteractionEnabled:") { // macOS 11.3+
		wkPrefs.send("setTextInteractionEnabled:", prefs.TextInteractionEnabled.Get())
	}
	if prefs.FullscreenEnabled.IsSet() && respondsTo(wkPrefs, "setElementFullscreenEnabled:") { // macOS 12.3+
		wkPrefs.send("setElementFullscreenEnabled:", prefs.FullscreenEnabled.Get())
	}
	if prefs.JavaScriptCanOpenWindowsAutomatically.IsSet() {
		wkPrefs.send("setJavaScriptCanOpenWindowsAutomatically:", prefs.JavaScriptCanOpenWindowsAutomatically.Get())
	}
	if prefs.MinimumFontSize.IsSet() {
		wkPrefs.send("setMinimumFontSize:", prefs.MinimumFontSize.Get())
	}
	if respondsTo(wkPrefs, "setFraudulentWebsiteWarningEnabled:") { // macOS 10.15+
		wkPrefs.send("setFraudulentWebsiteWarningEnabled:", fraudulentWarnings)
	}

	// Configuration-level preferences.
	if prefs.AllowsAirPlayForMediaPlayback.IsSet() {
		config.send("setAllowsAirPlayForMediaPlayback:", prefs.AllowsAirPlayForMediaPlayback.Get())
	}
	if prefs.EnableAutoplayWithoutUserAction.IsSet() && prefs.EnableAutoplayWithoutUserAction.Get() {
		const wkAudiovisualMediaTypeNone = 0
		config.send("setMediaTypesRequiringUserActionForPlayback:", uint(wkAudiovisualMediaTypeNone))
	}

	// WebView-level preferences.
	if prefs.AllowsBackForwardNavigationGestures.IsSet() {
		webView.send("setAllowsBackForwardNavigationGestures:", prefs.AllowsBackForwardNavigationGestures.Get())
	}
	if prefs.AllowsMagnification.IsSet() {
		webView.send("setAllowsMagnification:", prefs.AllowsMagnification.Get())
	}
}

// ---------------------------------------------------------------------------
// Core operations
// ---------------------------------------------------------------------------

func (w *macosWebviewWindow) setTitle(title string) {
	runOnMain(func() { w.win().send("setTitle:", nsString(title)) })
}

func (w *macosWebviewWindow) setURL(uri string) {
	runOnMain(func() {
		u := nsURL(uri)
		req := class("NSURLRequest").send("requestWithURL:", u)
		w.webview().send("loadRequest:", req)
	})
}

func (w *macosWebviewWindow) setHTML(html string) {
	runOnMain(func() { w.webview().send("loadHTMLString:baseURL:", nsString(html), objc.ID(0)) })
}

func (w *macosWebviewWindow) execJS(js string) {
	runOnMain(func() {
		w.webview().send("evaluateJavaScript:completionHandler:", nsString(js), objc.ID(0))
	})
}

func (w *macosWebviewWindow) execJSDragOver(buffer []byte) {
	// buffer is NUL-terminated; trim before wrapping in an NSString.
	s := buffer
	if n := len(s); n > 0 && s[n-1] == 0 {
		s = s[:n-1]
	}
	js := string(s)
	runOnMain(func() {
		w.webview().send("evaluateJavaScript:completionHandler:", nsString(js), objc.ID(0))
	})
}

func (w *macosWebviewWindow) show() {
	runOnMain(func() {
		w.win().send("makeKeyAndOrderFront:", objc.ID(0))
		class("NSApplication").send("sharedApplication").send("activateIgnoringOtherApps:", true)
	})
}

func (w *macosWebviewWindow) hide() { runOnMain(func() { w.win().send("orderOut:", objc.ID(0)) }) }

func (w *macosWebviewWindow) close() {
	// Mirror cgo windowClose: send close directly. performClose: consults the
	// close button and is a documented no-op (beep) on frameless windows,
	// which have no NSWindowStyleMaskClosable bit.
	atomic.StoreUint32(&w.parent.unconditionallyClose, 1)
	runOnMain(func() { w.win().send("close") })
}

func (w *macosWebviewWindow) destroy() {
	// Mirror cgo: mark destroyed BEFORE the NSWindow deallocs
	// (releasedWhenClosed defaults to YES) so the public API guards reject
	// any further calls instead of messaging freed memory.
	w.parent.markAsDestroyed()
	runOnMain(func() {
		delegateToWindowID.Delete(uintptr(w.delegate))
		nsWindowToID.Delete(uintptr(w.nsWindow))
		clearWindowDragCache(w.parent.id)
		w.win().send("close")
	})
}

func (w *macosWebviewWindow) center() { runOnMain(func() { w.win().send("center") }) }

func (w *macosWebviewWindow) focus() {
	runOnMain(func() {
		w.win().send("makeKeyAndOrderFront:", objc.ID(0))
		class("NSApplication").send("sharedApplication").send("activateIgnoringOtherApps:", true)
	})
}

func (w *macosWebviewWindow) reload() { runOnMain(func() { w.webview().send("reload") }) }
func (w *macosWebviewWindow) forceReload() {
	runOnMain(func() { w.webview().send("reloadFromOrigin") })
}

func (w *macosWebviewWindow) minimise() {
	runOnMain(func() { w.win().send("miniaturize:", objc.ID(0)) })
}
func (w *macosWebviewWindow) unminimise() {
	runOnMain(func() { w.win().send("deminiaturize:", objc.ID(0)) })
}
func (w *macosWebviewWindow) maximise() {
	runOnMain(func() {
		if !w.isMaximised() {
			w.win().send("zoom:", objc.ID(0))
		}
	})
}
func (w *macosWebviewWindow) unmaximise() {
	runOnMain(func() {
		if w.isMaximised() {
			w.win().send("zoom:", objc.ID(0))
		}
	})
}
func (w *macosWebviewWindow) fullscreen() {
	runOnMain(func() {
		if !w.isFullscreen() {
			w.win().send("toggleFullScreen:", objc.ID(0))
		}
	})
}
func (w *macosWebviewWindow) unfullscreen() {
	runOnMain(func() {
		if w.isFullscreen() {
			w.win().send("toggleFullScreen:", objc.ID(0))
		}
	})
}

func (w *macosWebviewWindow) isMinimised() bool {
	return w.syncMainThreadReturningBool(func() bool { return get[bool](w.win(), "isMiniaturized") })
}
func (w *macosWebviewWindow) isMaximised() bool {
	return w.syncMainThreadReturningBool(func() bool { return get[bool](w.win(), "isZoomed") })
}
func (w *macosWebviewWindow) isFullscreen() bool {
	return w.syncMainThreadReturningBool(func() bool {
		const nsWindowStyleMaskFullScreen = 1 << 14
		return get[uint](w.win(), "styleMask")&nsWindowStyleMaskFullScreen != 0
	})
}
func (w *macosWebviewWindow) isVisible() bool {
	return w.syncMainThreadReturningBool(func() bool { return get[bool](w.win(), "isVisible") })
}
func (w *macosWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}
func (w *macosWebviewWindow) isFocused() bool {
	return w.syncMainThreadReturningBool(func() bool { return get[bool](w.win(), "isKeyWindow") })
}

func (w *macosWebviewWindow) syncMainThreadReturningBool(fn func() bool) bool {
	var result bool
	runOnMain(func() { result = fn() })
	return result
}

func (w *macosWebviewWindow) setResizable(resizable bool) {
	runOnMain(func() {
		const mask = 1 << 3 // NSWindowStyleMaskResizable
		cur := get[uint](w.win(), "styleMask")
		if resizable {
			cur |= mask
		} else {
			cur &^= mask
		}
		w.win().send("setStyleMask:", cur)
	})
}

func (w *macosWebviewWindow) setSize(width, height int) {
	runOnMain(func() {
		frame := get[NSRect](w.win(), "frame")
		frame.Size = CGSize{Width: CGFloat(width), Height: CGFloat(height)}
		w.win().send("setFrame:display:", frame, true)
	})
}

func (w *macosWebviewWindow) size() (int, int) {
	var wd, ht int
	runOnMain(func() {
		cv := w.win().send("contentView")
		frame := get[NSRect](cv, "frame")
		wd, ht = int(frame.Size.Width), int(frame.Size.Height)
	})
	return wd, ht
}

func (w *macosWebviewWindow) width() int  { wd, _ := w.size(); return wd }
func (w *macosWebviewWindow) height() int { _, ht := w.size(); return ht }

func (w *macosWebviewWindow) setMinSize(width, height int) {
	runOnMain(func() { w.win().send("setContentMinSize:", CGSize{Width: CGFloat(width), Height: CGFloat(height)}) })
}
func (w *macosWebviewWindow) setMaxSize(width, height int) {
	mw, mh := CGFloat(width), CGFloat(height)
	if width == 0 {
		mw = 1 << 30
	}
	if height == 0 {
		mh = 1 << 30
	}
	runOnMain(func() { w.win().send("setContentMaxSize:", CGSize{Width: mw, Height: mh}) })
}

func (w *macosWebviewWindow) setPosition(x, y int) {
	runOnMain(func() {
		// Convert top-left (Wails) to bottom-left (Cocoa) using the main screen.
		screen := class("NSScreen").send("mainScreen")
		sf := get[NSRect](screen, "frame")
		frame := get[NSRect](w.win(), "frame")
		top := sf.Size.Height - CGFloat(y)
		w.win().send("setFrameTopLeftPoint:", CGPoint{X: CGFloat(x), Y: top})
		_ = frame
	})
}

func (w *macosWebviewWindow) position() (int, int) {
	var x, y int
	runOnMain(func() {
		screen := class("NSScreen").send("mainScreen")
		sf := get[NSRect](screen, "frame")
		frame := get[NSRect](w.win(), "frame")
		x = int(frame.Origin.X)
		y = int(sf.Size.Height - (frame.Origin.Y + frame.Size.Height))
	})
	return x, y
}

func (w *macosWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	runOnMain(func() {
		const nsFloatingWindowLevel = 3
		const nsNormalWindowLevel = 0
		if alwaysOnTop {
			w.win().send("setLevel:", nsFloatingWindowLevel)
		} else {
			w.win().send("setLevel:", nsNormalWindowLevel)
		}
	})
}

func (w *macosWebviewWindow) setBackgroundColour(colour RGBA) {
	runOnMain(func() {
		c := class("NSColor").send("colorWithRed:green:blue:alpha:",
			CGFloat(colour.Red)/255, CGFloat(colour.Green)/255, CGFloat(colour.Blue)/255, CGFloat(colour.Alpha)/255)
		w.win().send("setBackgroundColor:", c)
	})
}

func (w *macosWebviewWindow) setContentProtection(enabled bool) {
	runOnMain(func() {
		const nsWindowSharingNone = 0
		const nsWindowSharingReadOnly = 1
		if enabled {
			w.win().send("setSharingType:", nsWindowSharingNone)
		} else {
			w.win().send("setSharingType:", nsWindowSharingReadOnly)
		}
	})
}

func (w *macosWebviewWindow) getZoom() float64 {
	var z float64
	runOnMain(func() { z = get[float64](w.webview(), "pageZoom") })
	if z == 0 {
		z = 1
	}
	return z
}
func (w *macosWebviewWindow) setZoom(zoom float64) {
	runOnMain(func() { w.webview().send("setPageZoom:", zoom) })
}
func (w *macosWebviewWindow) zoomIn() { w.setZoom(w.getZoom() + 0.1) }
func (w *macosWebviewWindow) zoomOut() {
	z := w.getZoom() - 0.1
	if z < 0.1 {
		z = 0.1
	}
	w.setZoom(z)
}
func (w *macosWebviewWindow) zoomReset() { w.setZoom(1) }
func (w *macosWebviewWindow) zoom()      { w.zoomReset() }

func (w *macosWebviewWindow) nativeWindow() unsafe.Pointer { return w.nsWindow }

func (w *macosWebviewWindow) on(eventID uint) { /* hasListeners() is always true */ }

// Standard edit actions dispatched through the responder chain.
func (w *macosWebviewWindow) sendAction(sel string) {
	runOnMain(func() {
		class("NSApplication").send("sharedApplication").
			send("sendAction:to:from:", sel_(sel), objc.ID(0), w.win())
	})
}
func (w *macosWebviewWindow) cut()       { w.sendAction("cut:") }
func (w *macosWebviewWindow) copy()      { w.sendAction("copy:") }
func (w *macosWebviewWindow) paste()     { w.sendAction("paste:") }
func (w *macosWebviewWindow) delete()    { w.sendAction("delete:") }
func (w *macosWebviewWindow) selectAll() { w.sendAction("selectAll:") }
func (w *macosWebviewWindow) undo()      { w.sendAction("undo:") }
func (w *macosWebviewWindow) redo()      { w.sendAction("redo:") }

// ---------------------------------------------------------------------------
// Bounds helpers
// ---------------------------------------------------------------------------

func (w *macosWebviewWindow) bounds() Rect {
	x, y := w.position()
	wd, ht := w.size()
	return Rect{X: x, Y: y, Width: wd, Height: ht}
}
func (w *macosWebviewWindow) setBounds(bounds Rect) {
	w.setPosition(bounds.X, bounds.Y)
	w.setSize(bounds.Width, bounds.Height)
}
func (w *macosWebviewWindow) physicalBounds() Rect          { return w.bounds() }
func (w *macosWebviewWindow) setPhysicalBounds(bounds Rect) { w.setBounds(bounds) }
func (w *macosWebviewWindow) relativePosition() (int, int)  { return w.position() }
func (w *macosWebviewWindow) setRelativePosition(x, y int)  { w.setPosition(x, y) }
func (w *macosWebviewWindow) centerOnScreen(screen *Screen) { w.center() }

func (w *macosWebviewWindow) getScreen() (*Screen, error) { return getScreenForWindow(w) }

// ---------------------------------------------------------------------------
// Remaining window operations
// ---------------------------------------------------------------------------

func (w *macosWebviewWindow) handleKeyEvent(acceleratorString string) {
	accelerator, err := parseAccelerator(acceleratorString)
	if err != nil {
		globalApplication.error("unable to parse accelerator: %w", err)
		return
	}
	w.parent.processKeyBinding(accelerator.String())
}

// getBorderSizes returns zero insets, matching the cgo backend.
func (w *macosWebviewWindow) getBorderSizes() *LRTB { return &LRTB{} }

func (w *macosWebviewWindow) print() error {
	runOnMain(func() {
		const paginationAutomatic = 0
		const orientationLandscape = 1
		pInfo := class("NSPrintInfo").send("sharedPrintInfo")
		pInfo.send("setHorizontalPagination:", paginationAutomatic)
		pInfo.send("setVerticalPagination:", paginationAutomatic)
		pInfo.send("setVerticallyCentered:", true)
		pInfo.send("setHorizontallyCentered:", true)
		pInfo.send("setOrientation:", orientationLandscape)
		pInfo.send("setLeftMargin:", CGFloat(30))
		pInfo.send("setRightMargin:", CGFloat(30))
		pInfo.send("setTopMargin:", CGFloat(30))
		pInfo.send("setBottomMargin:", CGFloat(30))
		po := w.webview().send("printOperationWithPrintInfo:", pInfo)
		po.send("setShowsPrintPanel:", true)
		po.send("setShowsProgressPanel:", true)
		po.send("view").send("setFrame:", get[NSRect](w.webview(), "bounds"))
		// runOperation does not work with WKWebView; must run modal for window.
		po.send("runOperationModalForWindow:delegate:didRunSelector:contextInfo:",
			w.win(), id(uintptr(w.delegate)), objc.SEL(0), unsafe.Pointer(nil))
	})
	return nil
}

// startResize is never called; native resize is handled by the OS.
func (w *macosWebviewWindow) startResize(_ string) error { return nil }

func (w *macosWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	if menu.impl == nil {
		menu.impl = newMenuImpl(menu)
	}
	thisMenu := menu.impl.(*macosMenu)
	thisMenu.update()
	runOnMain(func() {
		nsMenu := id(uintptr(thisMenu.nsMenu))
		nsMenu.send("popUpMenuPositioningItem:atLocation:inView:",
			objc.ID(0), CGPoint{X: CGFloat(data.X), Y: CGFloat(data.Y)}, w.webview())
	})
}

func (w *macosWebviewWindow) setFrameless(frameless bool) {
	runOnMain(func() {
		w.setFullSizeContent(frameless)
		if frameless {
			w.setTitleBarAppearsTransparent(true)
			w.setTitleVisibility(true)
		} else {
			tb := w.parent.options.Mac.TitleBar
			w.setTitleBarAppearsTransparent(tb.AppearsTransparent)
			w.setTitleVisibility(tb.HideTitle)
		}
	})
}

func (w *macosWebviewWindow) setHasShadow(hasShadow bool) {
	runOnMain(func() { w.win().send("setHasShadow:", hasShadow) })
}

func (w *macosWebviewWindow) setFullscreenButtonState(state ButtonState) {
	eff := effectiveZoomButtonState(state, w.parent.options.MaximiseButtonState)
	setStdButtonState(w, nsWindowZoomButton, eff)
}

func (w *macosWebviewWindow) disableSizeConstraints() {
	runOnMain(func() {
		w.win().send("setContentMinSize:", CGSize{})
		w.win().send("setContentMaxSize:", CGSize{})
	})
}

func (w *macosWebviewWindow) windowZoom()    { w.maximise() }
func (w *macosWebviewWindow) restore()       { w.unminimise() }
func (w *macosWebviewWindow) restoreWindow() { w.unminimise() }

// setEnabled matches the cgo backend, where windowSetEnabled is intentionally a
// no-op.
func (w *macosWebviewWindow) setEnabled(enabled bool) {}

func (w *macosWebviewWindow) flash(_ bool) {}

func (w *macosWebviewWindow) setWindowLevel(level MacWindowLevel) {
	// Classic AppKit NSWindowLevel values.
	lvl := 0
	switch level {
	case MacWindowLevelNormal:
		lvl = 0
	case MacWindowLevelFloating, MacWindowLevelTornOffMenu:
		lvl = 3
	case MacWindowLevelModalPanel:
		lvl = 8
	case MacWindowLevelMainMenu:
		lvl = 24
	case MacWindowLevelStatus:
		lvl = 25
	case MacWindowLevelPopUpMenu:
		lvl = 101
	case MacWindowLevelScreenSaver:
		lvl = 1000
	}
	runOnMain(func() { w.win().send("setLevel:", lvl) })
}

func (w *macosWebviewWindow) setCollectionBehavior(behavior MacWindowCollectionBehavior) {
	runOnMain(func() {
		const fullScreenPrimary = 1 << 7 // NSWindowCollectionBehaviorFullScreenPrimary
		b := uint(behavior)
		if b == 0 {
			b = fullScreenPrimary
		}
		w.win().send("setCollectionBehavior:", b)
	})
}

func (w *macosWebviewWindow) startDrag() error {
	runOnMain(func() {
		ev := class("NSApplication").send("sharedApplication").send("currentEvent")
		if !ev.isNil() {
			w.win().send("performWindowDragWithEvent:", ev)
		}
	})
	return nil
}

func (w *macosWebviewWindow) setMinimiseButtonState(state ButtonState) {
	setStdButtonState(w, nsWindowMiniaturizeButton, state)
}
func (w *macosWebviewWindow) setMaximiseButtonState(state ButtonState) {
	setStdButtonState(w, nsWindowZoomButton, state)
}
func (w *macosWebviewWindow) setCloseButtonState(state ButtonState) {
	setStdButtonState(w, nsWindowCloseButton, state)
}

func (w *macosWebviewWindow) isIgnoreMouseEvents() bool {
	return w.syncMainThreadReturningBool(func() bool { return get[bool](w.win(), "ignoresMouseEvents") })
}
func (w *macosWebviewWindow) setIgnoreMouseEvents(ignore bool) {
	runOnMain(func() { w.win().send("setIgnoresMouseEvents:", ignore) })
}

func (w *macosWebviewWindow) attachModal(modalWindow *WebviewWindow) {
	if modalWindow == nil || modalWindow.impl == nil || modalWindow.isDestroyed() {
		return
	}
	modalNativeWindow := modalWindow.impl.nativeWindow()
	if modalNativeWindow == nil {
		return
	}
	modal := id(uintptr(modalNativeWindow))
	runOnMain(func() {
		block := objc.NewBlock(func(b objc.Block, returnCode int) {})
		w.win().send("beginSheet:completionHandler:", modal, block)
		// beginSheet: copies the block; release our +1.
		block.Release()
	})
}

func (w *macosWebviewWindow) showMenuBar()    {}
func (w *macosWebviewWindow) hideMenuBar()    {}
func (w *macosWebviewWindow) toggleMenuBar()  {}
func (w *macosWebviewWindow) setMenu(_ *Menu) {}
func (w *macosWebviewWindow) snapAssist()     {}

// ---------------------------------------------------------------------------
// Titlebar / backdrop / appearance helpers (used by run()).
// ---------------------------------------------------------------------------

const (
	nsWindowCloseButton       = 0
	nsWindowMiniaturizeButton = 1
	nsWindowZoomButton        = 2
)

func (w *macosWebviewWindow) setTitleVisibility(hidden bool) {
	const nsWindowTitleVisible = 0
	const nsWindowTitleHidden = 1
	if hidden {
		w.win().send("setTitleVisibility:", nsWindowTitleHidden)
	} else {
		w.win().send("setTitleVisibility:", nsWindowTitleVisible)
	}
}

func (w *macosWebviewWindow) setTitleBarAppearsTransparent(transparent bool) {
	w.win().send("setTitlebarAppearsTransparent:", transparent)
}

func (w *macosWebviewWindow) setHideTitleBar(hide bool) {
	const titled = 1 << 0 // NSWindowStyleMaskTitled
	cur := get[uint](w.win(), "styleMask")
	if hide {
		cur &^= titled
	} else {
		cur |= titled
	}
	w.win().send("setStyleMask:", cur)
}

func (w *macosWebviewWindow) setFullSizeContent(fullSize bool) {
	const fullSizeContent = 1 << 15
	cur := get[uint](w.win(), "styleMask")
	if fullSize {
		cur |= fullSizeContent
	} else {
		cur &^= fullSizeContent
	}
	w.win().send("setStyleMask:", cur)
}

func (w *macosWebviewWindow) setUseToolbar(use bool) {
	if use {
		toolbar := class("NSToolbar").send("alloc").send("initWithIdentifier:", nsString("wails.toolbar"))
		w.win().send("setToolbar:", toolbar)
	} else {
		w.win().send("setToolbar:", objc.ID(0))
	}
}

func (w *macosWebviewWindow) setToolbarStyle(style int) {
	if !w.win().send("toolbar").isNil() {
		w.win().send("setToolbarStyle:", style)
	}
}

func (w *macosWebviewWindow) setHideToolbarSeparator(hide bool) {
	toolbar := w.win().send("toolbar")
	if !toolbar.isNil() {
		toolbar.send("setShowsBaselineSeparator:", !hide)
	}
}

func (w *macosWebviewWindow) setTransparentBackdrop() {
	w.win().send("setOpaque:", false)
	w.win().send("setBackgroundColor:", class("NSColor").send("clearColor"))
}

func (w *macosWebviewWindow) setWebviewTransparent() {
	w.webview().send("setValue:forKey:", nsNumberBool(false), nsString("drawsBackground"))
}

func (w *macosWebviewWindow) setTranslucentBackdrop() {
	contentView := w.win().send("contentView")
	bounds := get[NSRect](contentView, "bounds")
	const behindWindow = 0
	const stateActive = 1
	const belowWindow = -1
	const autoWidth, autoHeight = 1 << 1, 1 << 4
	effectView := class("NSVisualEffectView").send("alloc").send("initWithFrame:", bounds)
	effectView.send("setAutoresizingMask:", uint(autoWidth|autoHeight))
	effectView.send("setBlendingMode:", behindWindow)
	effectView.send("setState:", stateActive)
	contentView.send("addSubview:positioned:relativeTo:", effectView, belowWindow, objc.ID(0))
}

func (w *macosWebviewWindow) setAppearanceByName(name string) {
	appearance := class("NSAppearance").send("appearanceNamed:", nsString(name))
	w.win().send("setAppearance:", appearance)
}

func (w *macosWebviewWindow) injectCSS(css string) {
	js := "(function() { var style = document.createElement('style'); style.appendChild(document.createTextNode('" +
		css + "')); document.head.appendChild(style); })();"
	w.webview().send("evaluateJavaScript:completionHandler:", nsString(js), objc.ID(0))
}

// applyLiquidGlass applies Apple's Liquid Glass effect (macOS 26+), falling back
// to a translucent backdrop on older systems. Port of the cgo applyLiquidGlass +
// windowSetLiquidGlass.
func (w *macosWebviewWindow) applyLiquidGlass() {
	options := w.parent.options.Mac.LiquidGlass
	if options.CornerRadius < 0 {
		options.CornerRadius = 0
	}
	if !classExists("NSGlassEffectView") {
		runOnMain(func() {
			w.setTranslucentBackdrop()
			w.setWebviewTransparent()
		})
		globalApplication.debug("Liquid Glass not supported on this macOS version, falling back to translucent", "window", w.parent.id)
		return
	}
	runOnMain(func() { w.windowSetLiquidGlass(options) })
}

func (w *macosWebviewWindow) windowSetLiquidGlass(o MacLiquidGlass) {
	w.removeVisualEffects()

	glass := class("NSGlassEffectView").send("alloc").send("init").send("autorelease")
	if o.CornerRadius > 0 && respondsTo(glass, "setCornerRadius:") {
		glass.send("setValue:forKey:", nsNumberDouble(o.CornerRadius), nsString("cornerRadius"))
	}
	if o.TintColor != nil && o.TintColor.Alpha > 0 && respondsTo(glass, "setTintColor:") {
		tint := class("NSColor").send("colorWithRed:green:blue:alpha:",
			CGFloat(o.TintColor.Red)/255, CGFloat(o.TintColor.Green)/255,
			CGFloat(o.TintColor.Blue)/255, CGFloat(o.TintColor.Alpha)/255)
		glass.send("setTintColor:", tint)
	}
	if respondsTo(glass, "setStyle:") {
		lightStyle := int(o.Style)
		if o.Style == LiquidGlassStyleVibrant {
			lightStyle = int(LiquidGlassStyleLight)
		}
		glass.send("setValue:forKey:", nsNumberInt(lightStyle), nsString("style"))
	}
	if o.GroupID != "" {
		switch {
		case respondsTo(glass, "setGroupIdentifier:"):
			glass.send("setGroupIdentifier:", nsString(o.GroupID))
		case respondsTo(glass, "setGroupName:"):
			glass.send("setGroupName:", nsString(o.GroupID))
		}
	}
	if o.GroupSpacing > 0 && respondsTo(glass, "setGroupSpacing:") {
		glass.send("setValue:forKey:", nsNumberDouble(o.GroupSpacing), nsString("groupSpacing"))
	}

	const autoWidth, autoHeight = 1 << 1, 1 << 4
	const belowWindow = -1
	contentView := w.win().send("contentView")
	glass.send("setFrame:", get[NSRect](contentView, "bounds"))
	glass.send("setAutoresizingMask:", uint(autoWidth|autoHeight))
	contentView.send("addSubview:positioned:relativeTo:", glass, belowWindow, objc.ID(0))

	// A real NSGlassEffectView hosts the webview in its own contentView.
	if respondsTo(glass, "contentView") {
		webView := w.webview()
		glassContent := glass.send("contentView")
		if !webView.isNil() && !glassContent.isNil() {
			webView.send("removeFromSuperview")
			glassContent.send("addSubview:", webView)
			webView.send("setFrame:", get[NSRect](glassContent, "bounds"))
			webView.send("setAutoresizingMask:", uint(autoWidth|autoHeight))
		}
	}

	w.configureWebViewForLiquidGlass()
	w.win().send("setOpaque:", false)
	w.win().send("setBackgroundColor:", class("NSColor").send("clearColor"))
}

func (w *macosWebviewWindow) configureWebViewForLiquidGlass() {
	wv := w.webview()
	wv.send("setValue:forKey:", nsNumberBool(false), nsString("drawsBackground"))
	wv.send("setValue:forKey:", class("NSColor").send("clearColor"), nsString("backgroundColor"))
	layer := wv.send("layer")
	if !layer.isNil() {
		layer.send("setZPosition:", CGFloat(1.0))
		layer.send("setShouldRasterize:", true)
		scale := get[CGFloat](class("NSScreen").send("mainScreen"), "backingScaleFactor")
		layer.send("setRasterizationScale:", scale)
	}
}

// removeVisualEffects removes any NSVisualEffectView / NSGlassEffectView backdrop
// previously added to the content view.
func (w *macosWebviewWindow) removeVisualEffects() {
	contentView := w.win().send("contentView")
	subviews := contentView.send("subviews")
	count := get[uint](subviews, "count")
	veClass := class("NSVisualEffectView")
	glassClass := class("NSGlassEffectView") // id(0) if unavailable
	for i := uint(0); i < count; i++ {
		sv := subviews.send("objectAtIndex:", i)
		if get[bool](sv, "isKindOfClass:", veClass) ||
			(!glassClass.isNil() && get[bool](sv, "isKindOfClass:", glassClass)) {
			sv.send("removeFromSuperview")
		}
	}
}

func classExists(name string) bool {
	loadFrameworks()
	return objc.GetClass(name) != 0
}

func respondsTo(o id, sel string) bool {
	return get[bool](o, "respondsToSelector:", sel_(sel))
}

func nsNumberDouble(d float64) id {
	return class("NSNumber").send("numberWithDouble:", d)
}

// setStdButtonState toggles a standard window button, matching the cgo
// setButtonState semantics exactly: hidden when state==ButtonHidden, disabled
// when state==ButtonDisabled, else enabled.
func setStdButtonState(w *macosWebviewWindow, buttonType int, state ButtonState) {
	runOnMain(func() {
		btn := w.win().send("standardWindowButton:", buttonType)
		if btn.isNil() {
			return
		}
		btn.send("setHidden:", state == ButtonHidden)
		btn.send("setEnabled:", state != ButtonDisabled)
	})
}
