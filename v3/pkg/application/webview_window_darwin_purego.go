//go:build darwin && purego && !ios && !server

package application

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/ebitengine/purego/objc"
	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/runtime"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// macosWindowNative is the purego backend's per-window native state, which the
// cgo backend keeps in its Objective-C classes.
type macosWindowNative struct {
	wkWebView unsafe.Pointer // WKWebView*
	delegate  unsafe.Pointer // WailsWebviewWindowDelegate*
	dragView  unsafe.Pointer // WailsWebviewDrag* (only when EnableFileDrop)

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

		// windowWillClose: — emit the event, then run the per-window teardown
		// the cgo backend performs in the WebviewWindow dealloc: break the
		// userContentController→delegate retain cycle, nil the assign/weak
		// delegate pointers before the delegate can die, release our owning
		// delegate reference (from `new` in createWindow) and clear all
		// per-window registries. Both the user-initiated close and destroy()
		// funnel through here.
		willCloseID := macWindowEventID("WindowWillClose")
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("windowWillClose:"),
			Fn: func(self objc.ID, cmd objc.SEL, notif objc.ID) {
				wid, ok := windowIDForDelegate(self)
				if !ok {
					return
				}
				processWindowEvent(wid, willCloseID)

				if impl := windowImplForID(wid); impl != nil {
					if wv := impl.webview(); !wv.isNil() {
						wv.send("configuration").send("userContentController").
							send("removeScriptMessageHandlerForName:", nsString("external"))
						wv.send("setNavigationDelegate:", objc.ID(0))
						wv.send("setUIDelegate:", objc.ID(0))
					}
					if impl.dragView != nil {
						dragViewToWindowID.Delete(uintptr(impl.dragView))
					}
				}
				win := id(notif).send("object")
				if !win.isNil() {
					win.send("setDelegate:", objc.ID(0))
					windowDisableEscape.Delete(win.ptr())
					nsWindowToID.Delete(win.ptr())
					clearLeftMouseDown(win)
				}
				delegateToWindowID.Delete(uintptr(self))
				windowImplCache.Delete(wid)
				clearWindowDragCache(wid)
				// The configuration still retains us as its URL-scheme
				// handler, so this autorelease only becomes the final release
				// once the webview graph is torn down.
				id(self).send("autorelease")
			},
		})

		// NSDraggingDestination — the NSWindow forwards dragging messages to
		// its delegate. cgo registers EVERY window for file drags, so the
		// WindowFileDragging* events fire even with EnableFileDrop disabled
		// (the overlay view, when installed, sits above and takes precedence).
		dragEnteredID := macWindowEventID("WindowFileDraggingEntered")
		dragExitedID := macWindowEventID("WindowFileDraggingExited")
		dragPerformedID := macWindowEventID("WindowFileDraggingPerformed")
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("draggingEntered:"),
			Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) uint {
				if wid, ok := windowIDForDelegate(self); ok && pasteboardHasFiles(id(sender)) {
					processWindowEvent(wid, dragEnteredID)
					return nsDragOperationCopy
				}
				return nsDragOperationNone
			},
		})
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("draggingExited:"),
			Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) {
				if wid, ok := windowIDForDelegate(self); ok {
					processWindowEvent(wid, dragExitedID)
				}
			},
		})
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("prepareForDragOperation:"),
			Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) bool {
				return true
			},
		})
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("performDragOperation:"),
			Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) bool {
				wid, ok := windowIDForDelegate(self)
				if !ok {
					return false
				}
				processWindowEvent(wid, dragPerformedID)
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
				// cgo converts the drop point into the webview's coordinate
				// space and flips Y for the web coordinate system.
				var x, y int
				if impl := windowImplForID(wid); impl != nil {
					if wv := impl.webview(); !wv.isNil() {
						pt := get[NSPoint](s, "draggingLocation")
						inView := get[NSPoint](wv, "convertPoint:fromView:", pt, objc.ID(0))
						h := get[NSRect](wv, "frame").Size.Height
						x, y = int(inView.X), int(h-inView.Y)
					}
				}
				processDragItems(wid, filenames, x, y)
				// cgo returns NO here (the drop is fully handled Go-side).
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
				// cgo cancelURLRequest: end the handler's request context, so a
				// streaming response notices the fetch was aborted.
				webview.CancelRequest(unsafe.Pointer(uintptr(task)))
				webview.MarkTaskStopped(unsafe.Pointer(uintptr(task)))
				// Close a streamed request body like the cgo handler does, so an
				// in-flight body reader observes the cancellation promptly
				// instead of consuming the rest of a cancelled upload.
				if req := id(task).send("request"); !req.isNil() {
					if stream := req.send("HTTPBodyStream"); !stream.isNil() {
						stream.send("close")
					}
				}
			},
		})

		// WKNavigationDelegate: navigation lifecycle events (same set the cgo
		// delegate forwards).
		for sel, field := range map[string]string{
			"webView:didStartProvisionalNavigation:":                    "WebViewDidStartProvisionalNavigation",
			"webView:didReceiveServerRedirectForProvisionalNavigation:": "WebViewDidReceiveServerRedirectForProvisionalNavigation",
			"webView:didCommitNavigation:":                              "WebViewDidCommitNavigation",
			"webView:didFinishNavigation:":                              "WebViewDidFinishNavigation",
		} {
			evID := macWindowEventID(field)
			methods = append(methods, objc.MethodDef{
				Cmd: sel_(sel),
				Fn: func(self objc.ID, cmd objc.SEL, wv objc.ID, nav objc.ID) {
					if wid, ok := windowIDForDelegate(self); ok {
						processWindowEvent(wid, evID)
					}
				},
			})
		}

		// WKNavigationDelegate: recover from WebContent (renderer) process
		// termination — notably after macOS sleep. Without the reload the
		// WKWebView is left showing an unresponsive blank page (cgo parity).
		terminatedID := macWindowEventID("WebViewWebContentProcessDidTerminate")
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("webViewWebContentProcessDidTerminate:"),
			Fn: func(self objc.ID, cmd objc.SEL, wv objc.ID) {
				if wid, ok := windowIDForDelegate(self); ok {
					processWindowEvent(wid, terminatedID)
				}
				id(wv).send("reload")
			},
		})

		// WKUIDelegate: camera/microphone (getUserMedia) answered from the
		// window's Permissions; see permissions_darwin_purego.go.
		methods = append(methods, mediaCapturePermissionMethod())

		// WKUIDelegate: handle <input type="file"> clicks with an NSOpenPanel,
		// mirroring the cgo runOpenPanelWithParameters: implementation.
		methods = append(methods, objc.MethodDef{
			Cmd: sel_("webView:runOpenPanelWithParameters:initiatedByFrame:completionHandler:"),
			Fn: func(self objc.ID, cmd objc.SEL, wv objc.ID, params objc.ID, frame objc.ID, handler objc.ID) {
				// Copy the completion handler so it survives until the panel
				// completes (we invoke it asynchronously).
				completion := objc.Block(handler).Copy()
				panel := class("NSOpenPanel").send("openPanel")
				panel.send("setAllowsMultipleSelection:", get[bool](id(params), "allowsMultipleSelection"))
				if respondsTo(id(params), "allowsDirectories") {
					panel.send("setCanChooseDirectories:", get[bool](id(params), "allowsDirectories"))
				}
				window := id(wv).send("window")
				block := objc.NewBlock(func(b objc.Block, result int) {
					var urls objc.ID
					if result == nsModalResponseOK {
						urls = panel.send("URLs").raw()
					}
					invokeForeignBlock(objc.ID(completion), uintptr(urls))
					completion.Release()
				})
				panel.send("beginSheetModalForWindow:completionHandler:", window, block)
				block.Release()
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

		// Notch windows are placed on their camera housing by configureNotchWebView;
		// a generic position here would move them off it (cgo parity).
		if w.parent.notchWindow == nil {
			if options.Screen != nil {
				if options.InitialPosition == WindowCentered {
					w.centerOnScreen(options.Screen)
				} else {
					w.setPositionOnScreen(options.X, options.Y, options.Screen.ID)
				}
			} else if options.InitialPosition == WindowCentered {
				w.center()
			} else {
				w.setPosition(options.X, options.Y)
			}
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
				} else {
					// Shadow/level can only be applied once the window has
					// been shown; defer until it first becomes key (cgo parity).
					var cancel func()
					cancel = w.parent.OnWindowEvent(events.Mac.WindowDidBecomeKey, func(_ *WindowEvent) {
						InvokeAsync(func() {
							if !w.isVisible() {
								w.parent.Show()
							}
							w.setHasShadow(!options.Mac.DisableShadow)
							w.setAlwaysOnTop(options.AlwaysOnTop)
							cancel()
						})
					})
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
	w.applyWindowButtonStates()

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
	const backingBuffered = 2
	cornerSquare := options.Mac.CornerType == MacWindowCornerTypeSquare
	styleMask := int(macWindowStyleMask(options.Frameless, cornerSquare, options.Mac.CornerRadius))
	// A notch window is always a panel (cgo builds it only in the panel path).
	isNotch := w.parent.notchWindow != nil
	isPanel := options.Mac.WindowClass == MacWindowClassPanel || isNotch
	windowClass := registerWebviewWindowClass()
	if isPanel {
		styleMask = int(panelStyleMask(uint(styleMask), options.Mac.PanelPreferences))
		windowClass = registerWebviewPanelClass()
	}
	if isNotch {
		windowClass = registerNotchWindowClass()
	}
	width, height := options.Width, options.Height
	if width == 0 {
		width = 800
	}
	if height == 0 {
		height = 600
	}

	win := windowClass.send("alloc").send("initWithContentRect:styleMask:backing:defer:",
		rect(0, 0, CGFloat(width-1), CGFloat(height-1)), uint(styleMask), uint(backingBuffered), false)
	w.nsWindow = unsafe.Pointer(win.ptr())

	// cgo's WebviewWindow initWithContentRect: override applies these to every
	// window (transparent/frameless rendering + drag-by-background).
	win.send("setAlphaValue:", CGFloat(1.0))
	win.send("setBackgroundColor:", class("NSColor").send("clearColor"))
	win.send("setOpaque:", false)
	win.send("setMovableByWindowBackground:", true)
	if isPanel {
		configureNativePanel(win, options.Mac.PanelPreferences)
	}

	// cgo registers every window for file drags (window-level
	// NSDraggingDestination on the delegate) regardless of EnableFileDrop.
	win.send("registerForDraggedTypes:",
		class("NSArray").send("arrayWithObject:", nsString(nsFilenamesPboardType)))

	// Delegate (also the script/scheme/navigation handler). The `new` +1 is
	// our owning reference, released in the windowWillClose: teardown.
	del := registerWindowDelegateClass().send("new")
	w.delegate = unsafe.Pointer(del.ptr())
	win.send("setDelegate:", del)
	delegateToWindowID.Store(del.ptr(), w.parent.id)
	nsWindowToID.Store(win.ptr(), w.parent.id)
	windowImplCache.Store(w.parent.id, w)

	// Content view (owned by the window after setContentView:).
	view := class("NSView").send("alloc").send("initWithFrame:", rect(0, 0, CGFloat(width-1), CGFloat(height-1)))
	const autoWidth, autoHeight = 1 << 1, 1 << 4
	view.send("setAutoresizingMask:", uint(autoWidth|autoHeight))
	applyContentLayer(view, macCreationContentLayer(options.Frameless, cornerSquare, options.Mac.CornerRadius), options.Mac.CornerRadius)
	win.send("setContentView:", view)
	view.send("release")
	if macCreationHidesNativeTitlebar(options.Frameless, cornerSquare, options.Mac.CornerRadius) {
		w.setTitleBarAppearsTransparent(true)
		w.setTitleVisibility(true)
	}

	// WebView configuration.
	config := class("WKWebViewConfiguration").send("alloc").send("init")
	config.send("setSuppressesIncrementalRendering:", true)
	appName := options.Mac.WebviewPreferences.ApplicationNameForUserAgent
	if appName == "" {
		appName = "wails.io"
	}
	config.send("setApplicationNameForUserAgent:", nsString(appName))
	config.send("setURLSchemeHandler:forURLScheme:", del, nsString("wails"))

	// User content controller + external message bridge (retained by the
	// configuration; the handler registration is undone in windowWillClose:
	// to break the controller→delegate retain cycle, like cgo's dealloc).
	ucc := class("WKUserContentController").send("new")
	ucc.send("addScriptMessageHandler:name:", del, nsString("external"))
	// The asset-request body shim, as the cgo backend installs it
	// (darwin_body_fetch_shim.go). WKUserScriptInjectionTimeAtDocumentStart = 0.
	bodyShim := class("WKUserScript").send("alloc").send("initWithSource:injectionTime:forMainFrameOnly:",
		nsString(darwinBodyFetchShimJS), 0, false)
	ucc.send("addUserScript:", bodyShim)
	bodyShim.send("release")
	config.send("setUserContentController:", ucc)
	ucc.send("release")

	webView := class("WKWebView").send("alloc").send("initWithFrame:configuration:",
		rect(0, 0, CGFloat(width), CGFloat(height)), config)
	w.wkWebView = unsafe.Pointer(webView.ptr())
	view.send("addSubview:", webView)
	webView.send("setNavigationDelegate:", del)
	webView.send("setUIDelegate:", del)
	webView.send("setAutoresizingMask:", uint(autoWidth|autoHeight))

	applyWebviewPreferences(webView, config, options.Mac.WebviewPreferences, options.Mac.EnableFraudulentWebsiteWarnings)
	// The webview copied the configuration and the content view retains the
	// webview; drop our creation references (cgo autoreleases both).
	config.send("release")
	webView.send("release")

	if options.EnableFileDrop {
		w.installFileDropView(view, width, height)
	}
	if isNotch {
		configureNotchWebView(win, webView, w.parent.notchWindow)
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

func registerWebviewWindowClass() id {
	registerWindowClassOnce.Do(func() {
		webviewWindowClass = registerDelegateClass("WailsWebviewNSWindow", "NSWindow", nil, webviewWindowMethods(&webviewWindowClass))
	})
	return webviewWindowClass
}

// webviewWindowMethods is the method set shared by the NSWindow subclass and
// the NSPanel subclass (webview_panel_darwin_purego.go), mirroring cgo's
// WebviewWindow/WebviewPanel overrides. Callers must not register the same
// selector twice on one class. definer points at the class the methods are
// registered on; super calls start from it (see sendSuperFrom).
func webviewWindowMethods(definer *id) []objc.MethodDef {
	yes := func(self objc.ID, cmd objc.SEL) bool { return true }
	return []objc.MethodDef{
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
			sendSuperFrom(*definer, self, "cancelOperation:", sender)
		}},
		{Cmd: sel_("keyDown:"), Fn: func(self objc.ID, cmd objc.SEL, event objc.ID) {
			v, ok := nsWindowToID.Load(uintptr(self))
			if !ok {
				return
			}
			if accelerator := acceleratorFromNSEvent(id(event)); accelerator != "" {
				processWindowKeyDownEvent(v.(uint), accelerator)
			}
		}},
		// performKeyEquivalent: mirrors cgo -[WebviewWindow performKeyEquivalent:]. See keyequivalent_darwin_purego.go.
		{Cmd: sel_("performKeyEquivalent:"), Fn: func(self objc.ID, cmd objc.SEL, event objc.ID) bool {
			if v, ok := nsWindowToID.Load(uintptr(self)); ok {
				if processKeyEquivalent(v.(uint), acceleratorFromNSEvent(id(event))) {
					return true
				}
			}
			return sendSuperFromBool(*definer, self, "performKeyEquivalent:", event)
		}},
	}
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
	w.dragView = unsafe.Pointer(dragView.ptr())
	contentView.send("addSubview:", dragView)
	// The content view retains the drag view; drop the creation reference.
	dragView.send("release")
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
	// cgo parity: frameless windows have no titlebar to set.
	if w.parent.options.Frameless {
		return
	}
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
	// cgo parity: asynchronous, and guarded against shutdown / destroyed
	// windows so JS is never dispatched into a closing webview.
	InvokeAsync(func() {
		globalApplication.shutdownLock.Lock()
		performingShutdown := globalApplication.performingShutdown
		globalApplication.shutdownLock.Unlock()
		if performingShutdown {
			return
		}
		if w.nsWindow == nil || w.parent.isDestroyed() {
			return
		}
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
	// As in cgo, every show (notch windows too) reapplies the shadow option.
	defer w.setHasShadow(!w.parent.options.Mac.DisableShadow)
	if w.parent.notchWindow != nil {
		w.showNotch()
		return
	}
	// cgo windowShow is makeKeyAndOrderFront: only — Show() must not steal
	// application focus from other apps. Non-activating panels only order front.
	runOnMain(func() {
		if isNonActivatingPanel(w.win()) {
			w.win().send("orderFrontRegardless")
			return
		}
		w.win().send("makeKeyAndOrderFront:", objc.ID(0))
	})
}

func (w *macosWebviewWindow) hide() {
	if w.parent.notchWindow != nil {
		w.hideNotch()
		return
	}
	runOnMain(func() { w.win().send("orderOut:", objc.ID(0)) })
}

func (w *macosWebviewWindow) close() {
	// Mirror cgo windowClose: send close directly. performClose: consults the
	// close button and is a documented no-op (beep) on borderless windows
	// (frameless with square or custom-radius corners), which have no
	// NSWindowStyleMaskClosable bit.
	atomic.StoreUint32(&w.parent.unconditionallyClose, 1)
	runOnMain(func() { w.win().send("close") })
}

func (w *macosWebviewWindow) destroy() {
	// Mirror cgo: mark destroyed BEFORE the NSWindow deallocs
	// (releasedWhenClosed defaults to YES) so the public API guards reject
	// any further calls instead of messaging freed memory. The per-window
	// registry/delegate teardown runs in the windowWillClose: handler, which
	// this close triggers.
	w.parent.markAsDestroyed()
	runOnMain(func() { w.win().send("close") })
}

func (w *macosWebviewWindow) center() {
	// cgo windowCenter: exact visibleFrame centre of the window's screen
	// ([NSWindow center] biases toward the top third).
	runOnMain(func() {
		vf := get[NSRect](w.windowScreen(), "visibleFrame")
		frame := get[NSRect](w.win(), "frame")
		x := vf.Origin.X + (vf.Size.Width-frame.Size.Width)/2
		y := vf.Origin.Y + (vf.Size.Height-frame.Size.Height)/2
		w.win().send("setFrame:display:", rect(x, y, frame.Size.Width, frame.Size.Height), true)
	})
}

func (w *macosWebviewWindow) focus() {
	// cgo windowFocus: only activate the app when it is inactive, then make
	// the window key.
	runOnMain(func() {
		if isNonActivatingPanel(w.win()) {
			w.win().send("orderFrontRegardless")
			w.win().send("makeKeyWindow")
			return
		}
		app := class("NSApplication").send("sharedApplication")
		if !get[bool](app, "isActive") {
			app.send("activateIgnoringOtherApps:", true)
		}
		w.win().send("makeKeyAndOrderFront:", objc.ID(0))
		w.win().send("makeKeyWindow")
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
	// cgo parity: occlusion state, not NSWindow.isVisible (a fully occluded
	// window reports not-visible).
	return w.syncMainThreadReturningBool(func() bool {
		const nsWindowOcclusionStateVisible = uint(1 << 1)
		return get[uint](w.win(), "occlusionState")&nsWindowOcclusionStateVisible == nsWindowOcclusionStateVisible
	})
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
		// cgo windowSetSize: derive the content size from the requested frame
		// size, set it, then set the frame (animated).
		content := get[NSRect](w.win(), "contentRectForFrameRect:", rect(0, 0, CGFloat(width), CGFloat(height)))
		w.win().send("setContentSize:", content.Size)
		frame := get[NSRect](w.win(), "frame")
		w.win().send("setFrame:display:animate:",
			rect(frame.Origin.X, frame.Origin.Y, CGFloat(width), CGFloat(height)), true, true)
	})
}

func (w *macosWebviewWindow) size() (int, int) {
	// cgo windowGetSize reads the window FRAME size (not the content view),
	// keeping Set/Get symmetric so bounds() round-trips don't shrink the
	// window by the titlebar height.
	var wd, ht int
	runOnMain(func() {
		frame := get[NSRect](w.win(), "frame")
		wd, ht = int(frame.Size.Width), int(frame.Size.Height)
	})
	return wd, ht
}

func (w *macosWebviewWindow) width() int  { wd, _ := w.size(); return wd }
func (w *macosWebviewWindow) height() int { _, ht := w.size(); return ht }

func (w *macosWebviewWindow) setMinSize(width, height int) {
	runOnMain(func() {
		// cgo treats width/height as FRAME size: convert to content size and
		// set both limits.
		content := get[NSRect](w.win(), "contentRectForFrameRect:", rect(0, 0, CGFloat(width), CGFloat(height)))
		w.win().send("setContentMinSize:", content.Size)
		w.win().send("setMinSize:", CGSize{Width: CGFloat(width), Height: CGFloat(height)})
	})
}

func (w *macosWebviewWindow) setMaxSize(width, height int) {
	// cgo uses FLT_MAX for unset axes.
	const fltMax = CGFloat(3.4028234663852886e+38)
	mw, mh := CGFloat(width), CGFloat(height)
	if width <= 0 {
		mw = fltMax
	}
	if height <= 0 {
		mh = fltMax
	}
	runOnMain(func() {
		content := get[NSRect](w.win(), "contentRectForFrameRect:", rect(0, 0, mw, mh))
		w.win().send("setContentMaxSize:", content.Size)
		w.win().send("setMaxSize:", CGSize{Width: mw, Height: mh})
	})
}

// primaryScreen returns [[NSScreen screens] firstObject] (the primary screen,
// whose top-left is the origin of the Wails coordinate space), falling back to
// mainScreen like the cgo helpers.
func primaryScreen() id {
	s := class("NSScreen").send("screens").send("firstObject")
	if s.isNil() {
		s = class("NSScreen").send("mainScreen")
	}
	return s
}

// windowScreen returns the screen the window is on, falling back to mainScreen.
func (w *macosWebviewWindow) windowScreen() id {
	screen := w.win().send("screen")
	if screen.isNil() {
		screen = class("NSScreen").send("mainScreen")
	}
	return screen
}

func (w *macosWebviewWindow) setPosition(x, y int) {
	runOnMain(func() {
		// Canonical Wails coordinates: logical points, Y-down, origin at the
		// top-left of the PRIMARY screen (cgo windowSetPosition).
		primaryHeight := get[NSRect](primaryScreen(), "frame").Size.Height
		frame := get[NSRect](w.win(), "frame")
		frame.Origin.X = CGFloat(x)
		frame.Origin.Y = primaryHeight - frame.Size.Height - CGFloat(y)
		w.win().send("setFrame:display:", frame, true)
	})
}

func (w *macosWebviewWindow) position() (int, int) {
	var x, y int
	runOnMain(func() {
		primaryHeight := get[NSRect](primaryScreen(), "frame").Size.Height
		frame := get[NSRect](w.win(), "frame")
		x = int(frame.Origin.X)
		y = int(primaryHeight - frame.Origin.Y - frame.Size.Height)
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

// Webview zoom uses WKWebView magnification like cgo (not pageZoom), with the
// same >=1.0 floor and 0.05 steps.
func (w *macosWebviewWindow) getZoom() float64 {
	var z float64
	runOnMain(func() { z = get[float64](w.webview(), "magnification") })
	return z
}
func (w *macosWebviewWindow) zoomIn() {
	runOnMain(func() {
		w.webview().send("setMagnification:", get[float64](w.webview(), "magnification")+0.05)
	})
}
func (w *macosWebviewWindow) zoomOut() {
	runOnMain(func() {
		z := get[float64](w.webview(), "magnification")
		if z > 1.05 {
			w.webview().send("setMagnification:", z-0.05)
		} else {
			w.webview().send("setMagnification:", float64(1.0))
		}
	})
}
func (w *macosWebviewWindow) zoomReset() {
	runOnMain(func() { w.webview().send("setMagnification:", float64(1.0)) })
}

// zoom toggles window maximise ([window zoom:nil]), matching cgo windowZoom —
// it is the WINDOW zoom, not the webview zoom.
func (w *macosWebviewWindow) zoom() {
	runOnMain(func() { w.win().send("zoom:", objc.ID(0)) })
}

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

// relativePosition returns the window position relative to its own screen's
// frame origin: X from the screen's left edge, Y from its top edge (Y-down).
// Mirrors cgo windowGetRelativePosition (issue #5408: must round-trip with
// setRelativePosition on every screen).
func (w *macosWebviewWindow) relativePosition() (int, int) {
	var x, y int
	runOnMain(func() {
		frame := get[NSRect](w.win(), "frame")
		sf := get[NSRect](w.windowScreen(), "frame")
		x = int(frame.Origin.X - sf.Origin.X)
		y = int((sf.Origin.Y + sf.Size.Height) - frame.Origin.Y - frame.Size.Height)
	})
	return x, y
}

func (w *macosWebviewWindow) setRelativePosition(x, y int) {
	runOnMain(func() {
		frame := get[NSRect](w.win(), "frame")
		sf := get[NSRect](w.windowScreen(), "frame")
		frame.Origin.X = sf.Origin.X + CGFloat(x)
		frame.Origin.Y = (sf.Origin.Y + sf.Size.Height) - frame.Size.Height - CGFloat(y)
		w.win().send("setFrame:display:animate:", frame, true, false)
	})
}

// screenByID resolves an NSScreen by the Wails Screen.ID (its CGDirectDisplayID
// rendered as a decimal string), falling back to mainScreen like cgo. Both the
// signed (cgo %d) and unsigned renderings are accepted.
func screenByID(screenID string) id {
	screens := class("NSScreen").send("screens")
	count := get[uint](screens, "count")
	for i := uint(0); i < count; i++ {
		s := screens.send("objectAtIndex:", i)
		num := s.send("deviceDescription").send("objectForKey:", nsString("NSScreenNumber"))
		if num.isNil() {
			continue
		}
		displayID := get[uint32](num, "unsignedIntValue")
		if fmt.Sprintf("%d", int32(displayID)) == screenID || fmt.Sprintf("%d", displayID) == screenID {
			return s
		}
	}
	return class("NSScreen").send("mainScreen")
}

func (w *macosWebviewWindow) centerOnScreen(screen *Screen) {
	if screen == nil {
		w.center()
		return
	}
	runOnMain(func() {
		vf := get[NSRect](screenByID(screen.ID), "visibleFrame")
		frame := get[NSRect](w.win(), "frame")
		x := vf.Origin.X + (vf.Size.Width-frame.Size.Width)/2
		y := vf.Origin.Y + (vf.Size.Height-frame.Size.Height)/2
		w.win().send("setFrameOrigin:", CGPoint{X: x, Y: y})
	})
}

// setPositionOnScreen positions the window relative to the given screen's
// visible frame; x/y are DIP top-origin coordinates (cgo
// windowSetPositionOnScreen parity).
func (w *macosWebviewWindow) setPositionOnScreen(x, y int, screenID string) {
	runOnMain(func() {
		target := screenByID(screenID)
		scale := get[CGFloat](target, "backingScaleFactor")
		vf := get[NSRect](target, "visibleFrame")
		frame := get[NSRect](w.win(), "frame")
		newX := vf.Origin.X + CGFloat(x)/scale
		newY := vf.Origin.Y + vf.Size.Height - frame.Size.Height - CGFloat(y)/scale
		w.win().send("setFrameOrigin:", CGPoint{X: newX, Y: newY})
	})
}

func (w *macosWebviewWindow) getScreen() (*Screen, error) { return getScreenForWindow(w) }

// ---------------------------------------------------------------------------
// Remaining window operations
// ---------------------------------------------------------------------------

// getBorderSizes returns zero insets, matching the cgo backend.
func (w *macosWebviewWindow) getBorderSizes() *LRTB { return &LRTB{} }

func (w *macosWebviewWindow) print() error {
	runOnMain(func() {
		// printOperationWithPrintInfo: is macOS 11+; cgo gates the whole body
		// with @available(macOS 11.0,*) — an unguarded miss would raise an
		// uncatchable NSException.
		if !respondsTo(w.webview(), "printOperationWithPrintInfo:") {
			return
		}
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

func (w *macosWebviewWindow) windowZoom() {
	runOnMain(func() { w.win().send("zoom:", objc.ID(0)) })
}

// restore mirrors cgo windowRestore: exit fullscreen, un-zoom, deminiaturize.
func (w *macosWebviewWindow) restore() {
	runOnMain(func() {
		const nsWindowStyleMaskFullScreen = 1 << 14
		if get[uint](w.win(), "styleMask")&nsWindowStyleMaskFullScreen != 0 {
			w.win().send("toggleFullScreen:", objc.ID(0))
		}
		if get[bool](w.win(), "isZoomed") {
			w.win().send("zoom:", objc.ID(0))
		}
		if get[bool](w.win(), "isMiniaturized") {
			w.win().send("deminiaturize:", objc.ID(0))
		}
	})
}
func (w *macosWebviewWindow) restoreWindow() { w.restore() }

// setEnabled matches the cgo backend, where windowSetEnabled is intentionally a
// no-op.
func (w *macosWebviewWindow) setEnabled(enabled bool) {
	// Kept identical to cgo windowSetEnabled, which performs no AppKit work.
	_ = enabled
}

// nsInformationalRequest is NSInformationalRequest, the NSRequestUserAttentionType
// value cgo windowFlash passes.
const nsInformationalRequest = 10

// flash matches cgo windowFlash: request user attention so the Dock icon
// bounces, drawing the user back to a window that needs them.
// NSInformationalRequest bounces the icon once; the request completes on its
// own, so disabling is a no-op.
func (w *macosWebviewWindow) flash(enabled bool) {
	if !enabled {
		return
	}
	runOnMain(func() { nsApp().send("requestUserAttention:", nsInformationalRequest) })
}

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
		// Anchor at the original mouse-down like cgo's retained
		// leftMouseEvent; NSApp.currentEvent has usually moved on by the time
		// the JS bridge message arrives.
		ev := takeLeftMouseDown(w.win())
		if ev.isNil() {
			ev = class("NSApplication").send("sharedApplication").send("currentEvent")
		}
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
	// Both MaximiseButtonState and FullscreenButtonState target the zoom
	// button; apply the more restrictive so neither silently overrides the
	// other (cgo parity).
	setStdButtonState(w, nsWindowZoomButton, effectiveZoomButtonState(state, w.parent.options.FullscreenButtonState))
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
		// The window retains its toolbar; drop the creation reference
		// (cgo autoreleases).
		toolbar.send("release")
	} else {
		w.win().send("setToolbar:", objc.ID(0))
	}
}

func (w *macosWebviewWindow) setToolbarStyle(style int) {
	// setToolbarStyle: is macOS 11+ (cgo guards with @available); a missed
	// selector is an uncatchable NSException.
	if respondsTo(w.win(), "setToolbarStyle:") && !w.win().send("toolbar").isNil() {
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
	wailsPrivateSetWebviewTransparent(w.webview())
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
	// The content view retains the subview; drop the creation reference.
	effectView.send("release")
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
	wailsPrivateSetGlassStyle(glass, o.Style)
	wailsPrivateSetGlassGrouping(glass, o.GroupID, o.GroupSpacing)

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
	wailsPrivateSetWebviewTransparent(wv)
	wailsPrivateSetWebviewBackgroundColour(wv, 0, 0, 0, 0)
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

// platformTitlebarDoubleClickPreference returns the user's
// AppleActionOnDoubleClick setting ("" when unset), as the cgo build does.
func platformTitlebarDoubleClickPreference() string {
	ud := class("NSUserDefaults").send("standardUserDefaults")
	return ud.send("stringForKey:", nsString("AppleActionOnDoubleClick")).string()
}

// setNonClientHitTestRegions is a no-op on macOS, matching the cgo build.
func (w *macosWebviewWindow) setNonClientHitTestRegions([]nonClientHitTestRegion) {}

// setNativeZoom sets the WKWebView magnification on the main thread.
func (w *macosWebviewWindow) setNativeZoom(zoom float64) {
	runOnMain(func() { w.webview().send("setMagnification:", zoom) })
}
