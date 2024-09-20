//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "webview_window_bindings_darwin.h"
*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/runtime"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type macosWebviewWindow struct {
	nsWindow unsafe.Pointer
	parent   *WebviewWindow
}

func (w *macosWebviewWindow) handleKeyEvent(acceleratorString string) {
	// Parse acceleratorString
	accelerator, err := parseAccelerator(acceleratorString)
	if err != nil {
		globalApplication.error("unable to parse accelerator: %s", err.Error())
		return
	}
	w.parent.processKeyBinding(accelerator.String())
}

func (w *macosWebviewWindow) getBorderSizes() *LRTB {
	return &LRTB{}
}

func (w *macosWebviewWindow) isFocused() bool {
	return bool(C.windowIsFocused(w.nsWindow))
}

func (w *macosWebviewWindow) setPosition(x int, y int) {
	C.windowSetPosition(w.nsWindow, C.int(x), C.int(y))
}

func (w *macosWebviewWindow) print() error {
	C.windowPrint(w.nsWindow)
	return nil
}

func (w *macosWebviewWindow) startResize(_ string) error {
	// Never called. Handled natively by the OS.
	return nil
}

func (w *macosWebviewWindow) focus() {
	// Make the window key and main
	C.windowFocus(w.nsWindow)
}

func (w *macosWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	C.windowShowMenu(w.nsWindow, thisMenu.nsMenu, C.int(data.X), C.int(data.Y))
}

func (w *macosWebviewWindow) getZoom() float64 {
	return float64(C.windowZoomGet(w.nsWindow))
}

func (w *macosWebviewWindow) setZoom(zoom float64) {
	C.windowZoomSet(w.nsWindow, C.double(zoom))
}

func (w *macosWebviewWindow) setFrameless(frameless bool) {
	C.windowSetFrameless(w.nsWindow, C.bool(frameless))
	if frameless {
		C.windowSetTitleBarAppearsTransparent(w.nsWindow, C.bool(true))
		C.windowSetHideTitle(w.nsWindow, C.bool(true))
	} else {
		macOptions := w.parent.options.Mac
		appearsTransparent := macOptions.TitleBar.AppearsTransparent
		hideTitle := macOptions.TitleBar.HideTitle
		C.windowSetTitleBarAppearsTransparent(w.nsWindow, C.bool(appearsTransparent))
		C.windowSetHideTitle(w.nsWindow, C.bool(hideTitle))
	}
}

func (w *macosWebviewWindow) setHasShadow(hasShadow bool) {
	C.windowSetShadow(w.nsWindow, C.bool(hasShadow))
}

func (w *macosWebviewWindow) getScreen() (*Screen, error) {
	return getScreenForWindow(w)
}

func (w *macosWebviewWindow) show() {
	C.windowShow(w.nsWindow)
}

func (w *macosWebviewWindow) hide() {
	C.windowHide(w.nsWindow)
}

func (w *macosWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	C.setFullscreenButtonEnabled(w.nsWindow, C.bool(enabled))
}

func (w *macosWebviewWindow) disableSizeConstraints() {
	C.windowDisableSizeConstraints(w.nsWindow)
}

func (w *macosWebviewWindow) unfullscreen() {
	C.windowUnFullscreen(w.nsWindow)
}

func (w *macosWebviewWindow) fullscreen() {
	C.windowFullscreen(w.nsWindow)
}

func (w *macosWebviewWindow) unminimise() {
	C.windowUnminimise(w.nsWindow)
}

func (w *macosWebviewWindow) unmaximise() {
	C.windowUnmaximise(w.nsWindow)
}

func (w *macosWebviewWindow) maximise() {
	C.windowMaximise(w.nsWindow)
}

func (w *macosWebviewWindow) minimise() {
	C.windowMinimise(w.nsWindow)
}

func (w *macosWebviewWindow) on(eventID uint) {
	//C.registerListener(C.uint(eventID))
}

func (w *macosWebviewWindow) zoom() {
	C.windowZoom(w.nsWindow)
}

func (w *macosWebviewWindow) windowZoom() {
	C.windowZoom(w.nsWindow)
}

func (w *macosWebviewWindow) close() {
	C.windowClose(w.nsWindow)
}

func (w *macosWebviewWindow) zoomIn() {
	C.windowZoomIn(w.nsWindow)
}

func (w *macosWebviewWindow) zoomOut() {
	C.windowZoomOut(w.nsWindow)
}

func (w *macosWebviewWindow) zoomReset() {
	C.windowZoomReset(w.nsWindow)
}

func (w *macosWebviewWindow) reload() {
	//TODO: Implement
	globalApplication.debug("reload called on WebviewWindow", "parentID", w.parent.id)
}

func (w *macosWebviewWindow) forceReload() {
	//TODO: Implement
	globalApplication.debug("force reload called on WebviewWindow", "parentID", w.parent.id)
}

func (w *macosWebviewWindow) center() {
	C.windowCenter(w.nsWindow)
}

func (w *macosWebviewWindow) isMinimised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return bool(C.windowIsMinimised(w.nsWindow))
	})
}

func (w *macosWebviewWindow) isMaximised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return bool(C.windowIsMaximised(w.nsWindow))
	})
}

func (w *macosWebviewWindow) isFullscreen() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return bool(C.windowIsFullscreen(w.nsWindow))
	})
}

func (w *macosWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}

func (w *macosWebviewWindow) isVisible() bool {
	return bool(C.isVisible(w.nsWindow))
}

func (w *macosWebviewWindow) syncMainThreadReturningBool(fn func() bool) bool {
	var wg sync.WaitGroup
	wg.Add(1)
	var result bool
	globalApplication.dispatchOnMainThread(func() {
		result = fn()
		wg.Done()
	})
	wg.Wait()
	return result
}

func (w *macosWebviewWindow) restore() {
	// restore window to normal size
	C.windowRestore(w.nsWindow)
}

func (w *macosWebviewWindow) restoreWindow() {
	C.windowRestore(w.nsWindow)
}

func (w *macosWebviewWindow) setEnabled(enabled bool) {
	C.windowSetEnabled(w.nsWindow, C.bool(enabled))
}

func (w *macosWebviewWindow) execJS(js string) {
	InvokeAsync(func() {
		if globalApplication.performingShutdown {
			return
		}
		if w.nsWindow == nil {
			return
		}
		C.windowExecJS(w.nsWindow, C.CString(js))
	})
}

func (w *macosWebviewWindow) setURL(uri string) {
	C.navigationLoadURL(w.nsWindow, C.CString(uri))
}

func (w *macosWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	C.windowSetAlwaysOnTop(w.nsWindow, C.bool(alwaysOnTop))
}

func newWindowImpl(parent *WebviewWindow) *macosWebviewWindow {
	result := &macosWebviewWindow{
		parent: parent,
	}
	result.parent.RegisterHook(events.Mac.WebViewDidFinishNavigation, func(event *WindowEvent) {
		result.execJS(runtime.Core())
	})
	return result
}

func (w *macosWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		cTitle := C.CString(title)
		C.windowSetTitle(w.nsWindow, cTitle)
	}
}

func (w *macosWebviewWindow) flash(_ bool) {
	// Not supported on macOS
}

func (w *macosWebviewWindow) setSize(width, height int) {
	C.windowSetSize(w.nsWindow, C.int(width), C.int(height))
}

func (w *macosWebviewWindow) setMinSize(width, height int) {
	C.windowSetMinSize(w.nsWindow, C.int(width), C.int(height))
}
func (w *macosWebviewWindow) setMaxSize(width, height int) {
	C.windowSetMaxSize(w.nsWindow, C.int(width), C.int(height))
}

func (w *macosWebviewWindow) setResizable(resizable bool) {
	C.windowSetResizable(w.nsWindow, C.bool(resizable))
}

func (w *macosWebviewWindow) size() (int, int) {
	var width, height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		C.windowGetSize(w.nsWindow, &width, &height)
		wg.Done()
	})
	wg.Wait()
	return int(width), int(height)
}

func (w *macosWebviewWindow) setRelativePosition(x, y int) {
	C.windowSetRelativePosition(w.nsWindow, C.int(x), C.int(y))
}

func (w *macosWebviewWindow) setWindowLevel(level MacWindowLevel) {
	switch level {
	case MacWindowLevelNormal:
		C.setNormalWindowLevel(w.nsWindow)
	case MacWindowLevelFloating:
		C.setFloatingWindowLevel(w.nsWindow)
	case MacWindowLevelTornOffMenu:
		C.setTornOffMenuWindowLevel(w.nsWindow)
	case MacWindowLevelModalPanel:
		C.setModalPanelWindowLevel(w.nsWindow)
	case MacWindowLevelMainMenu:
		C.setMainMenuWindowLevel(w.nsWindow)
	case MacWindowLevelStatus:
		C.setStatusWindowLevel(w.nsWindow)
	case MacWindowLevelPopUpMenu:
		C.setPopUpMenuWindowLevel(w.nsWindow)
	case MacWindowLevelScreenSaver:
		C.setScreenSaverWindowLevel(w.nsWindow)
	}
}

func (w *macosWebviewWindow) width() int {
	var width C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		width = C.windowGetWidth(w.nsWindow)
		wg.Done()
	})
	wg.Wait()
	return int(width)
}
func (w *macosWebviewWindow) height() int {
	var height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		height = C.windowGetHeight(w.nsWindow)
		wg.Done()
	})
	wg.Wait()
	return int(height)
}

func bool2CboolPtr(value bool) *C.bool {
	v := C.bool(value)
	return &v
}

func (w *macosWebviewWindow) getWebviewPreferences() C.struct_WebviewPreferences {
	wvprefs := w.parent.options.Mac.WebviewPreferences

	var result C.struct_WebviewPreferences

	if wvprefs.TextInteractionEnabled.IsSet() {
		result.TextInteractionEnabled = bool2CboolPtr(wvprefs.TextInteractionEnabled.Get())
	}
	if wvprefs.TabFocusesLinks.IsSet() {
		result.TabFocusesLinks = bool2CboolPtr(wvprefs.TabFocusesLinks.Get())
	}
	if wvprefs.FullscreenEnabled.IsSet() {
		result.FullscreenEnabled = bool2CboolPtr(wvprefs.FullscreenEnabled.Get())
	}

	return result
}

func (w *macosWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}
	globalApplication.dispatchOnMainThread(func() {
		options := w.parent.options
		macOptions := options.Mac

		w.nsWindow = C.windowNew(C.uint(w.parent.id),
			C.int(options.Width),
			C.int(options.Height),
			C.bool(macOptions.EnableFraudulentWebsiteWarnings),
			C.bool(options.Frameless),
			C.bool(options.EnableDragAndDrop),
			w.getWebviewPreferences(),
		)

		w.setup(&options, &macOptions)
	})
}

func (w *macosWebviewWindow) setup(options *WebviewWindowOptions, macOptions *MacWindow) {
	w.setTitle(options.Title)
	w.setAlwaysOnTop(options.AlwaysOnTop)
	w.setResizable(!options.DisableResize)
	if options.MinWidth != 0 || options.MinHeight != 0 {
		w.setMinSize(options.MinWidth, options.MinHeight)
	}
	if options.MaxWidth != 0 || options.MaxHeight != 0 {
		w.setMaxSize(options.MaxWidth, options.MaxHeight)
	}
	//w.setZoom(options.Zoom)
	w.enableDevTools()

	w.setBackgroundColour(options.BackgroundColour)

	switch macOptions.Backdrop {
	case MacBackdropTransparent:
		C.windowSetTransparent(w.nsWindow)
		C.webviewSetTransparent(w.nsWindow)
	case MacBackdropTranslucent:
		C.windowSetTranslucent(w.nsWindow)
		C.webviewSetTransparent(w.nsWindow)
	case MacBackdropNormal:
	}

	if macOptions.WindowLevel == "" {
		macOptions.WindowLevel = MacWindowLevelNormal
	}
	w.setWindowLevel(macOptions.WindowLevel)

	// Initialise the window buttons
	w.setMinimiseButtonState(options.MinimiseButtonState)
	w.setMaximiseButtonState(options.MaximiseButtonState)
	w.setCloseButtonState(options.CloseButtonState)

	// Ignore mouse events if requested
	w.setIgnoreMouseEvents(options.IgnoreMouseEvents)

	titleBarOptions := macOptions.TitleBar
	if !w.parent.options.Frameless {
		C.windowSetTitleBarAppearsTransparent(w.nsWindow, C.bool(titleBarOptions.AppearsTransparent))
		C.windowSetHideTitleBar(w.nsWindow, C.bool(titleBarOptions.Hide))
		C.windowSetHideTitle(w.nsWindow, C.bool(titleBarOptions.HideTitle))
		C.windowSetFullSizeContent(w.nsWindow, C.bool(titleBarOptions.FullSizeContent))
		C.windowSetUseToolbar(w.nsWindow, C.bool(titleBarOptions.UseToolbar))
		C.windowSetToolbarStyle(w.nsWindow, C.int(titleBarOptions.ToolbarStyle))
		C.windowSetShowToolbarWhenFullscreen(w.nsWindow, C.bool(titleBarOptions.ShowToolbarWhenFullscreen))
		C.windowSetHideToolbarSeparator(w.nsWindow, C.bool(titleBarOptions.HideToolbarSeparator))
	}

	if macOptions.Appearance != "" {
		C.windowSetAppearanceTypeByName(w.nsWindow, C.CString(string(macOptions.Appearance)))
	}

	if macOptions.InvisibleTitleBarHeight != 0 {
		C.windowSetInvisibleTitleBar(w.nsWindow, C.uint(macOptions.InvisibleTitleBarHeight))
	}

	switch w.parent.options.StartState {
	case WindowStateMaximised:
		w.maximise()
	case WindowStateMinimised:
		w.minimise()
	case WindowStateFullscreen:
		w.fullscreen()
	case WindowStateNormal:
	}
	C.windowCenter(w.nsWindow)

	startURL, err := assetserver.GetStartURL(options.URL)
	if err != nil {
		globalApplication.fatal(err.Error())
	}

	w.setURL(startURL)

	// We need to wait for the HTML to load before we can execute the javascript
	w.parent.OnWindowEvent(events.Mac.WebViewDidFinishNavigation, func(_ *WindowEvent) {
		InvokeAsync(func() {
			if options.JS != "" {
				w.execJS(options.JS)
			}
			if options.CSS != "" {
				C.windowInjectCSS(w.nsWindow, C.CString(options.CSS))
			}
			if !options.Hidden {
				C.windowShow(w.nsWindow)
				w.setHasShadow(!options.Mac.DisableShadow)
			} else {
				// We have to wait until the window is shown before we can remove the shadow
				var cancel func()
				cancel = w.parent.OnWindowEvent(events.Mac.WindowDidBecomeKey, func(_ *WindowEvent) {
					w.setHasShadow(!options.Mac.DisableShadow)
					cancel()
				})
			}
		})
	})

	// Translate ShouldClose to common WindowClosing event
	w.parent.OnWindowEvent(events.Mac.WindowShouldClose, func(_ *WindowEvent) {
		w.parent.emit(events.Common.WindowClosing)
	})

	// Translate WindowDidResignKey to common WindowLostFocus event
	w.parent.OnWindowEvent(events.Mac.WindowDidResignKey, func(_ *WindowEvent) {
		w.parent.emit(events.Common.WindowLostFocus)
	})
	w.parent.OnWindowEvent(events.Mac.WindowDidResignMain, func(_ *WindowEvent) {
		w.parent.emit(events.Common.WindowLostFocus)
	})
	w.parent.OnWindowEvent(events.Mac.WindowDidResize, func(_ *WindowEvent) {
		w.parent.emit(events.Common.WindowDidResize)
	})

	if options.HTML != "" {
		w.setHTML(options.HTML)
	}
}

func (w *macosWebviewWindow) nativeWindowHandle() uintptr {
	return uintptr(w.nsWindow)
}

func (w *macosWebviewWindow) setBackgroundColour(colour RGBA) {

	C.windowSetBackgroundColour(w.nsWindow, C.int(colour.Red), C.int(colour.Green), C.int(colour.Blue), C.int(colour.Alpha))
}

func (w *macosWebviewWindow) relativePosition() (int, int) {
	var x, y C.int
	InvokeSync(func() {
		C.windowGetRelativePosition(w.nsWindow, &x, &y)
	})

	return int(x), int(y)
}

func (w *macosWebviewWindow) position() (int, int) {
	var x, y C.int
	InvokeSync(func() {
		C.windowGetPosition(w.nsWindow, &x, &y)
	})

	return int(x), int(y)
}

func (w *macosWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
	C.windowDestroy(w.nsWindow)
}

func (w *macosWebviewWindow) setHTML(html string) {
	// Convert HTML to C string
	cHTML := C.CString(html)
	// Render HTML
	C.windowRenderHTML(w.nsWindow, cHTML)
}

func (w *macosWebviewWindow) startDrag() error {
	C.startDrag(w.nsWindow)
	return nil
}

func (w *macosWebviewWindow) setMinimiseButtonState(state ButtonState) {
	C.setMinimiseButtonState(w.nsWindow, C.int(state))
}

func (w *macosWebviewWindow) setMaximiseButtonState(state ButtonState) {
	C.setMaximiseButtonState(w.nsWindow, C.int(state))
}

func (w *macosWebviewWindow) setCloseButtonState(state ButtonState) {
	C.setCloseButtonState(w.nsWindow, C.int(state))
}

func (w *macosWebviewWindow) isIgnoreMouseEvents() bool {
	return bool(C.isIgnoreMouseEvents(w.nsWindow))
}

func (w *macosWebviewWindow) setIgnoreMouseEvents(ignore bool) {
	C.setIgnoreMouseEvents(w.nsWindow, C.bool(ignore))
}

func (w *macosWebviewWindow) cut() {
}

func (w *macosWebviewWindow) paste() {
}

func (w *macosWebviewWindow) copy() {
}

func (w *macosWebviewWindow) selectAll() {
}

func (w *macosWebviewWindow) undo() {
}

func (w *macosWebviewWindow) delete() {
}

func (w *macosWebviewWindow) redo() {
}
