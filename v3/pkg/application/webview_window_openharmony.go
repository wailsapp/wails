//go:build openharmony && !server

package application

import (
	"fmt"
	"math"
	"unsafe"
)

// openharmonyWebviewWindow implements webviewWindowImpl for OHOS. The ArkTS
// shell owns the ArkWeb component; window geometry and chrome are therefore
// either fixed (mobile fullscreen) or delegated to the shell, exactly like
// androidWebviewWindow. Only execJS crosses the bridge.
type openharmonyWebviewWindow struct {
	windowID  uint32 // Wails window ID for tracking
	parent    *WebviewWindow
	zoomLevel float64 // mirrored from the ArkTS WebviewController zoom cache
}

func newWindowImpl(parent *WebviewWindow) *openharmonyWebviewWindow {
	return &openharmonyWebviewWindow{parent: parent, zoomLevel: 1.0}
}

func (w *openharmonyWebviewWindow) center() { ohBridgeVoid("center", "{}") }

func (w *openharmonyWebviewWindow) close() { ohBridgeVoid("close", "{}") }

func (w *openharmonyWebviewWindow) destroy() { w.parent.markAsDestroyed() }

func (w *openharmonyWebviewWindow) execJS(js string) {
	// Execute JavaScript via the WailsBridge.executeJavaScript bridge method.
	executeJavaScript(js)
}

func (w *openharmonyWebviewWindow) flash(_ bool) {}

func (w *openharmonyWebviewWindow) focus() { ohBridgeVoid("focus", "{}") }

func (w *openharmonyWebviewWindow) forceReload() { ohBridgeVoid("forceReload", "{}") }

func (w *openharmonyWebviewWindow) fullscreen() { ohBridgeVoid("fullscreen", "{}") }

func (w *openharmonyWebviewWindow) getScreen() (*Screen, error) {
	screens, err := getScreens()
	if err != nil {
		return nil, err
	}
	if len(screens) == 0 {
		return nil, fmt.Errorf("no screens available")
	}
	return screens[0], nil
}

func (w *openharmonyWebviewWindow) getZoom() float64 {
	return w.zoomLevel
}

func (w *openharmonyWebviewWindow) handleDragAndDropMessage(_ string) {}

func (w *openharmonyWebviewWindow) hasParent() bool { return false }

func (w *openharmonyWebviewWindow) height() int {
	_, h := w.size()
	return h
}

func (w *openharmonyWebviewWindow) hide() { ohBridgeVoid("hide", "{}") }

func (w *openharmonyWebviewWindow) isAlwaysOnTop() bool { return false }

func (w *openharmonyWebviewWindow) isCloseRequested() bool { return false }

func (w *openharmonyWebviewWindow) setCloseRequested(_ bool) {}

func (w *openharmonyWebviewWindow) isFocused() bool { return true }

func (w *openharmonyWebviewWindow) isFullscreen() bool { return true } // OHOS apps are typically fullscreen

func (w *openharmonyWebviewWindow) isMaximised() bool { return true }

func (w *openharmonyWebviewWindow) isMinimised() bool { return false }

func (w *openharmonyWebviewWindow) isNormal() bool { return false }

func (w *openharmonyWebviewWindow) isVisible() bool { return true }

func (w *openharmonyWebviewWindow) maximise() { ohBridgeVoid("maximise", "{}") }

func (w *openharmonyWebviewWindow) minimise() { ohBridgeVoid("minimise", "{}") }

func (w *openharmonyWebviewWindow) openContextMenu(_ *Menu, _ *ContextMenuData) {}

func (w *openharmonyWebviewWindow) openDevTools() {}

func (w *openharmonyWebviewWindow) print() error { return nil }

func (w *openharmonyWebviewWindow) reload() { ohBridgeVoid("reload", "{}") }

func (w *openharmonyWebviewWindow) relativePosition() (int, int) { return 0, 0 }

func (w *openharmonyWebviewWindow) resizable() bool { return false }

func (w *openharmonyWebviewWindow) restore() {}

func (w *openharmonyWebviewWindow) setAbsolutePosition(_ int, _ int) {}

func (w *openharmonyWebviewWindow) setAlwaysOnTop(value bool) {
	ohBridgeVoid("setAlwaysOnTop", string(mustJSON(map[string]any{"value": value})))
}

func (w *openharmonyWebviewWindow) setBackgroundColour(colour RGBA) {
	ohBridgeVoid("setBackgroundColour", string(mustJSON(map[string]any{
		"r": colour.Red,
		"g": colour.Green,
		"b": colour.Blue,
		"a": colour.Alpha,
	})))
}

func (w *openharmonyWebviewWindow) setEnabled(_ bool) {}

func (w *openharmonyWebviewWindow) setFrameless(_ bool) {}

func (w *openharmonyWebviewWindow) setFullscreenButtonState(_ ButtonState) {}

func (w *openharmonyWebviewWindow) setMaxSize(_ int, _ int) {}

func (w *openharmonyWebviewWindow) setMinSize(_ int, _ int) {}

func (w *openharmonyWebviewWindow) setRelativePosition(_ int, _ int) {}

func (w *openharmonyWebviewWindow) setResizable(value bool) {
	ohBridgeVoid("setResizable", string(mustJSON(map[string]any{"value": value})))
}

func (w *openharmonyWebviewWindow) setSize(width, height int) {
	ohBridgeVoid("setSize", string(mustJSON(map[string]any{"width": width, "height": height})))
}

func (w *openharmonyWebviewWindow) setTitle(title string) {
	ohBridgeVoid("setTitle", jsonStringField("title", title))
}

func (w *openharmonyWebviewWindow) setZoom(zoom float64) {
	w.zoomLevel = zoom
	ohBridgeVoid("setZoom", string(mustJSON(map[string]any{"zoom": zoom})))
}

func (w *openharmonyWebviewWindow) show() { ohBridgeVoid("show", "{}") }

func (w *openharmonyWebviewWindow) size() (int, int) {
	if s, err := w.getScreen(); err == nil && s != nil {
		return s.Size.Width, s.Size.Height
	}
	return 0, 0
}

func (w *openharmonyWebviewWindow) toggleDevTools() {}

func (w *openharmonyWebviewWindow) unfullscreen() { ohBridgeVoid("unfullscreen", "{}") }

func (w *openharmonyWebviewWindow) unmaximise() { ohBridgeVoid("unmaximise", "{}") }

func (w *openharmonyWebviewWindow) unminimise() { ohBridgeVoid("unminimise", "{}") }

func (w *openharmonyWebviewWindow) width() int {
	wi, _ := w.size()
	return wi
}

func (w *openharmonyWebviewWindow) zoom() {}

func (w *openharmonyWebviewWindow) zoomIn() {
	w.zoomLevel += 0.1
	ohBridgeVoid("zoomIn", "{}")
}

func (w *openharmonyWebviewWindow) zoomOut() {
	w.zoomLevel = math.Max(0.1, w.zoomLevel-0.1)
	ohBridgeVoid("zoomOut", "{}")
}

func (w *openharmonyWebviewWindow) zoomReset() {
	w.zoomLevel = 1.0
	ohBridgeVoid("zoomReset", "{}")
}

func (w *openharmonyWebviewWindow) setParent(_ *WebviewWindow) error { return nil }

func (w *openharmonyWebviewWindow) run() {}

func (w *openharmonyWebviewWindow) setIgnoreMouseEvents(_ bool) {}

func (w *openharmonyWebviewWindow) setOpacity(_ float32) {}

func (w *openharmonyWebviewWindow) setTheme(_ Theme) {}

func (w *openharmonyWebviewWindow) setPinned(_ bool) {}

func (w *openharmonyWebviewWindow) startResize(_ string) error { return nil }

func (w *openharmonyWebviewWindow) startDrag() error { return nil }

func (w *openharmonyWebviewWindow) enableDevTools() {}

func (w *openharmonyWebviewWindow) disableContextMenu() {}

func (w *openharmonyWebviewWindow) disableDefaultContextMenu() {}

func (w *openharmonyWebviewWindow) setShouldClose(_ func() bool) {}

func (w *openharmonyWebviewWindow) absolutePosition() (int, int) { return 0, 0 }

func (w *openharmonyWebviewWindow) startMove() {}

func (w *openharmonyWebviewWindow) windowMenu() *Menu { return nil }

func (w *openharmonyWebviewWindow) setWindowMenu(_ *Menu) {}

func (w *openharmonyWebviewWindow) isIgnoreMouseEvents() bool { return false }

func (w *openharmonyWebviewWindow) flashCancel() {}

func (w *openharmonyWebviewWindow) setFocusable(_ bool) {}

func (w *openharmonyWebviewWindow) bounds() Rect { return Rect{} }

func (w *openharmonyWebviewWindow) copy() { ohBridgeVoid("copy", "{}") }

func (w *openharmonyWebviewWindow) cut() { ohBridgeVoid("cut", "{}") }

func (w *openharmonyWebviewWindow) paste() { ohBridgeVoid("paste", "{}") }

func (w *openharmonyWebviewWindow) selectAll() { ohBridgeVoid("selectAll", "{}") }

func (w *openharmonyWebviewWindow) undo() { ohBridgeVoid("undo", "{}") }

func (w *openharmonyWebviewWindow) redo() { ohBridgeVoid("redo", "{}") }

func (w *openharmonyWebviewWindow) delete() { ohBridgeVoid("delete", "{}") }

func (w *openharmonyWebviewWindow) getBorderSizes() *LRTB { return &LRTB{} }

func (w *openharmonyWebviewWindow) handleKeyEvent(_ string) {}

func (w *openharmonyWebviewWindow) hideMenuBar() {}

func (w *openharmonyWebviewWindow) unhideMenuBar() {}

func (w *openharmonyWebviewWindow) toggleMenuBar() {}

func (w *openharmonyWebviewWindow) isMenuBarHidden() bool { return true }

func (w *openharmonyWebviewWindow) nativeWindow() unsafe.Pointer { return nil }

func (w *openharmonyWebviewWindow) attachModal(_ *WebviewWindow) {}

func (w *openharmonyWebviewWindow) on(eventID uint) { registerOpenharmonyListener(eventID) }

func (w *openharmonyWebviewWindow) position() (int, int) { return 0, 0 }

func (w *openharmonyWebviewWindow) physicalBounds() Rect { return Rect{} }

func (w *openharmonyWebviewWindow) setBounds(bounds Rect) {
	ohBridgeVoid("setBounds", string(mustJSON(map[string]any{
		"x":      bounds.X,
		"y":      bounds.Y,
		"width":  bounds.Width,
		"height": bounds.Height,
	})))
}

func (w *openharmonyWebviewWindow) setMinimiseButtonState(_ ButtonState) {}

func (w *openharmonyWebviewWindow) setMaximiseButtonState(_ ButtonState) {}

func (w *openharmonyWebviewWindow) setCloseButtonState(_ ButtonState) {}

func (w *openharmonyWebviewWindow) setContentProtection(_ bool) {}

func (w *openharmonyWebviewWindow) setNonClientHitTestRegions(_ []nonClientHitTestRegion) {}

func (w *openharmonyWebviewWindow) setHTML(html string) {
	ohBridgeVoid("setHTML", jsonStringField("html", html))
}

func (w *openharmonyWebviewWindow) setMenu(_ *Menu) {}

func (w *openharmonyWebviewWindow) setPhysicalBounds(_ Rect) {}

func (w *openharmonyWebviewWindow) setPosition(x, y int) {
	ohBridgeVoid("setPosition", string(mustJSON(map[string]any{"x": x, "y": y})))
}

func (w *openharmonyWebviewWindow) centerOnScreen(_ *Screen) {}

func (w *openharmonyWebviewWindow) setURL(url string) {
	ohBridgeVoid("setURL", jsonStringField("url", url))
}

func (w *openharmonyWebviewWindow) showMenuBar() {}

func (w *openharmonyWebviewWindow) snapAssist() {}
