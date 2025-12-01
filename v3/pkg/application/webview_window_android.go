//go:build android

package application

import "unsafe"

// androidWebviewWindow implements the webviewWindowImpl interface for Android
type androidWebviewWindow struct {
	windowID uint32 // Wails window ID for tracking
	parent   *WebviewWindow
}

func newWindowImpl(parent *WebviewWindow) *androidWebviewWindow {
	return &androidWebviewWindow{
		parent: parent,
	}
}

func (w *androidWebviewWindow) center() {}

func (w *androidWebviewWindow) close() {}

func (w *androidWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
}

func (w *androidWebviewWindow) execJS(js string) {
	// Execute JavaScript via JNI callback to Java's WailsBridge.executeJavaScript()
	androidLogf("debug", "execJS: %s", js)
	executeJavaScript(js)
}

func (w *androidWebviewWindow) flash(_ bool) {}

func (w *androidWebviewWindow) focus() {}

func (w *androidWebviewWindow) forceReload() {}

func (w *androidWebviewWindow) fullscreen() {}

func (w *androidWebviewWindow) getScreen() (*Screen, error) {
	// Android has a single "screen" (the device display)
	return &Screen{
		ID:        "main",
		Name:      "Main Display",
		IsPrimary: true,
	}, nil
}

func (w *androidWebviewWindow) getZoom() float64 {
	return 1.0
}

func (w *androidWebviewWindow) handleDragAndDropMessage(_ string) {}

func (w *androidWebviewWindow) hasParent() bool {
	return false
}

func (w *androidWebviewWindow) height() int {
	return 2400 // Default Android height (common flagship)
}

func (w *androidWebviewWindow) hide() {}

func (w *androidWebviewWindow) isAlwaysOnTop() bool {
	return false
}

func (w *androidWebviewWindow) isCloseRequested() bool {
	return false
}

func (w *androidWebviewWindow) setCloseRequested(_ bool) {}

func (w *androidWebviewWindow) isFocused() bool {
	return true
}

func (w *androidWebviewWindow) isFullscreen() bool {
	return true // Android apps are typically fullscreen
}

func (w *androidWebviewWindow) isMaximised() bool {
	return true
}

func (w *androidWebviewWindow) isMinimised() bool {
	return false
}

func (w *androidWebviewWindow) isNormal() bool {
	return false
}

func (w *androidWebviewWindow) isVisible() bool {
	return true
}

func (w *androidWebviewWindow) maximise() {}

func (w *androidWebviewWindow) minimise() {}

func (w *androidWebviewWindow) openContextMenu(_ *Menu, _ *ContextMenuData) {}

func (w *androidWebviewWindow) openDevTools() {}

func (w *androidWebviewWindow) print() error {
	return nil
}

func (w *androidWebviewWindow) reload() {}

func (w *androidWebviewWindow) relativePosition() (int, int) {
	return 0, 0
}

func (w *androidWebviewWindow) resizable() bool {
	return false
}

func (w *androidWebviewWindow) restore() {}

func (w *androidWebviewWindow) setAbsolutePosition(_ int, _ int) {}

func (w *androidWebviewWindow) setAlwaysOnTop(_ bool) {}

func (w *androidWebviewWindow) setBackgroundColour(col RGBA) {
	androidLogf("debug", "setBackgroundColour: rgba(%d,%d,%d,%d)", col.Red, col.Green, col.Blue, col.Alpha)
	AndroidSetWebViewBackgroundColor(col.Red, col.Green, col.Blue, col.Alpha)
}

func (w *androidWebviewWindow) setEnabled(_ bool) {}

func (w *androidWebviewWindow) setFrameless(_ bool) {}

func (w *androidWebviewWindow) setFullscreenButtonEnabled(_ bool) {}

func (w *androidWebviewWindow) setMaxSize(_ int, _ int) {}

func (w *androidWebviewWindow) setMinSize(_ int, _ int) {}

func (w *androidWebviewWindow) setRelativePosition(_ int, _ int) {}

func (w *androidWebviewWindow) setResizable(_ bool) {}

func (w *androidWebviewWindow) setSize(_ int, _ int) {}

func (w *androidWebviewWindow) setTitle(_ string) {}

func (w *androidWebviewWindow) setZoom(_ float64) {}

func (w *androidWebviewWindow) show() {}

func (w *androidWebviewWindow) size() (int, int) {
	return 1080, 2400 // Default Android size (common flagship)
}

func (w *androidWebviewWindow) toggleDevTools() {}

func (w *androidWebviewWindow) unfullscreen() {}

func (w *androidWebviewWindow) unmaximise() {}

func (w *androidWebviewWindow) unminimise() {}

func (w *androidWebviewWindow) width() int {
	return 1080 // Default Android width
}

func (w *androidWebviewWindow) zoom() {}

func (w *androidWebviewWindow) zoomIn() {}

func (w *androidWebviewWindow) zoomOut() {}

func (w *androidWebviewWindow) zoomReset() {}

func (w *androidWebviewWindow) setParent(_ *WebviewWindow) error {
	return nil
}

func (w *androidWebviewWindow) run() {
	androidLogf("info", "androidWebviewWindow.run() called")
	// Android WebView is created and managed by the Java Activity
	// Just store the window ID for reference
	w.windowID = uint32(w.parent.ID())
}

func (w *androidWebviewWindow) setIgnoreMouseEvents(_ bool) {}

func (w *androidWebviewWindow) setOpacity(_ float32) {}

func (w *androidWebviewWindow) setTheme(_ Theme) {}

func (w *androidWebviewWindow) setPinned(_ bool) {}

func (w *androidWebviewWindow) startResize(_ string) error {
	return nil
}

func (w *androidWebviewWindow) startDrag() error {
	return nil
}

func (w *androidWebviewWindow) enableDevTools() {}

func (w *androidWebviewWindow) disableContextMenu() {}

func (w *androidWebviewWindow) disableDefaultContextMenu() {}

func (w *androidWebviewWindow) setShouldClose(_ func() bool) {}

func (w *androidWebviewWindow) absolutePosition() (int, int) {
	return 0, 0
}

func (w *androidWebviewWindow) startMove() {}

func (w *androidWebviewWindow) windowMenu() *Menu {
	return nil
}

func (w *androidWebviewWindow) setWindowMenu(_ *Menu) {}

func (w *androidWebviewWindow) isIgnoreMouseEvents() bool {
	return false
}

func (w *androidWebviewWindow) flashCancel() {}

func (w *androidWebviewWindow) setFocusable(_ bool) {}

func (w *androidWebviewWindow) bounds() Rect {
	return Rect{
		X:      0,
		Y:      0,
		Width:  1080,
		Height: 2400,
	}
}

func (w *androidWebviewWindow) copy() {
	// Android copy implementation
}

func (w *androidWebviewWindow) cut() {
	// Android cut implementation
}

func (w *androidWebviewWindow) paste() {
	// Android paste implementation
}

func (w *androidWebviewWindow) selectAll() {
	// Android select all implementation
}

func (w *androidWebviewWindow) undo() {
	// Android undo implementation
}

func (w *androidWebviewWindow) redo() {
	// Android redo implementation
}

func (w *androidWebviewWindow) delete() {
	// Android delete implementation
}

func (w *androidWebviewWindow) getBorderSizes() *LRTB {
	return &LRTB{}
}

func (w *androidWebviewWindow) handleKeyEvent(acceleratorString string) {
	// Android handle key event
}

func (w *androidWebviewWindow) hideMenuBar() {
	// Android doesn't have menu bar
}

func (w *androidWebviewWindow) unhideMenuBar() {
	// Android doesn't have menu bar
}

func (w *androidWebviewWindow) toggleMenuBar() {
	// Android doesn't have menu bar
}

func (w *androidWebviewWindow) isMenuBarHidden() bool {
	return true // Android doesn't have menu bar
}

func (w *androidWebviewWindow) nativeWindow() unsafe.Pointer {
	return nil
}

func (w *androidWebviewWindow) on(eventID uint) {
	// Android event handling
}

func (w *androidWebviewWindow) position() (int, int) {
	return 0, 0
}

func (w *androidWebviewWindow) physicalBounds() Rect {
	return Rect{
		X:      0,
		Y:      0,
		Width:  1080,
		Height: 2400,
	}
}

func (w *androidWebviewWindow) setBounds(bounds Rect) {
	// Android set bounds - not applicable on mobile
}

func (w *androidWebviewWindow) setMinimiseButtonState(_ ButtonState) {
	// Android doesn't have minimize buttons like desktop platforms
}

func (w *androidWebviewWindow) setMaximiseButtonState(_ ButtonState) {
	// Android doesn't have maximize buttons like desktop platforms
}

func (w *androidWebviewWindow) setCloseButtonState(_ ButtonState) {
	// Android doesn't have close buttons like desktop platforms
}

func (w *androidWebviewWindow) setContentProtection(_ bool) {
	// Android content protection - could be implemented with FLAG_SECURE
}

func (w *androidWebviewWindow) setHTML(html string) {
	// TODO: Implement via JNI
	androidLogf("debug", "setHTML called")
}

func (w *androidWebviewWindow) setMenu(_ *Menu) {
	// Android doesn't support window menus like desktop platforms
}

func (w *androidWebviewWindow) setPhysicalBounds(_ Rect) {
	// Android doesn't support arbitrary window bounds - apps are fullscreen
}

func (w *androidWebviewWindow) setPosition(_ int, _ int) {
	// Android doesn't support window positioning - apps are fullscreen
}

func (w *androidWebviewWindow) setURL(url string) {
	// TODO: Implement via JNI
	androidLogf("debug", "setURL: %s", url)
}

func (w *androidWebviewWindow) showMenuBar() {
	// Android doesn't have menu bars like desktop platforms
}

func (w *androidWebviewWindow) snapAssist() {
	// Android doesn't support window snap assist like Windows
}
