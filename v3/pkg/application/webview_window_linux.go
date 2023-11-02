//go:build linux

package application

import (
	"fmt"
	"net/url"

	"github.com/wailsapp/wails/v3/pkg/events"
)

var showDevTools = func(window pointer) {}

type dragInfo struct {
	XRoot       int
	YRoot       int
	DragTime    int
	MouseButton uint
}

type linuxWebviewWindow struct {
	id          uint
	application pointer
	window      pointer
	webview     pointer
	parent      *WebviewWindow
	menubar     pointer
	vbox        pointer
	menu        *Menu
	accels      pointer
	lastWidth   int
	lastHeight  int
	drag        dragInfo
}

var (
	registered bool = false // avoid 'already registered message' about 'wails://'
)

func (w *linuxWebviewWindow) startDrag() error {
	return nil
}

func (w *linuxWebviewWindow) endDrag(button uint, x, y int) {

}

func (w *linuxWebviewWindow) enableDND() {
	windowEnableDND(w.parent.id, w.webview)
}

func (w *linuxWebviewWindow) connectSignals() {
	cb := func(e events.WindowEventType) {
		w.parent.emit(e)
	}
	windowSetupSignalHandlers(w.parent.id, w.window, w.webview, cb)
}

func (w *linuxWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu manually because we don't want a gtk_menu_bar
	// as the top-level item
	ctxMenu := &linuxMenu{
		menu: menu,
	}
	if menu.impl == nil {
		ctxMenu.update()
	}

	native := ctxMenu.menu.impl.(*linuxMenu).native
	contextMenuShow(w.window, native, data)
}

func (w *linuxWebviewWindow) getZoom() float64 {
	return windowZoom(w.webview)
}

func (w *linuxWebviewWindow) setZoom(zoom float64) {
	windowZoomSet(w.webview, zoom)
}

func (w *linuxWebviewWindow) setFrameless(frameless bool) {
	windowSetFrameless(w.window, frameless)
}

func (w *linuxWebviewWindow) getScreen() (*Screen, error) {
	mx, my, width, height, scale := windowGetCurrentMonitorGeometry(w.window)
	return &Screen{
		ID:        fmt.Sprintf("%d", w.id),            // A unique identifier for the display
		Name:      w.parent.Name(),                    // The name of the display
		Scale:     float32(scale),                     // The scale factor of the display
		X:         mx,                                 // The x-coordinate of the top-left corner of the rectangle
		Y:         my,                                 // The y-coordinate of the top-left corner of the rectangle
		Size:      Size{Width: width, Height: height}, // The size of the display
		Bounds:    Rect{},                             // The bounds of the display
		WorkArea:  Rect{},                             // The work area of the display
		IsPrimary: false,                              // Whether this is the primary display
		Rotation:  0.0,                                // The rotation of the display
	}, nil
}

func (w *linuxWebviewWindow) focus() {
	windowPresent(w.window)
}

func (w *linuxWebviewWindow) show() {
	windowShow(w.window)
}

func (w *linuxWebviewWindow) hide() {
	windowHide(w.window)
}

func (w *linuxWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}

func (w *linuxWebviewWindow) isVisible() bool {
	return windowIsVisible(w.window)
}

func (w *linuxWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	//	C.setFullscreenButtonEnabled(w.nsWindow, C.bool(enabled))
	fmt.Println("setFullscreenButtonEnabled - not implemented")
}

func (w *linuxWebviewWindow) disableSizeConstraints() {
	x, y, width, height, scale := windowGetCurrentMonitorGeometry(w.window)
	w.setMinMaxSize(x, y, width*scale, height*scale)
}

func (w *linuxWebviewWindow) unfullscreen() {
	windowUnfullscreen(w.window)
	w.unmaximise()
}

func (w *linuxWebviewWindow) fullscreen() {
	w.maximise()
	w.lastWidth, w.lastHeight = w.size()
	x, y, width, height, scale := windowGetCurrentMonitorGeometry(w.window)
	if x == -1 && y == -1 && width == -1 && height == -1 {
		return
	}
	w.setMinMaxSize(0, 0, width*scale, height*scale)
	w.setSize(width*scale, height*scale)
	windowFullscreen(w.window)
	w.setRelativePosition(0, 0)
}

func (w *linuxWebviewWindow) unminimise() {
	windowPresent(w.window)
}

func (w *linuxWebviewWindow) unmaximise() {
	windowUnmaximize(w.window)
}

func (w *linuxWebviewWindow) maximise() {
	windowMaximize(w.window)
}

func (w *linuxWebviewWindow) minimise() {
	windowMinimize(w.window)
}

func (w *linuxWebviewWindow) flash(enabled bool) {
	// Not supported on linux
}

func (w *linuxWebviewWindow) on(eventID uint) {
	// Don't think this is correct!
	// GTK Events are strings
	fmt.Println("on()", eventID)
	//C.registerListener(C.uint(eventID))
}

func (w *linuxWebviewWindow) zoom() {
	w.zoomIn()
}

func (w *linuxWebviewWindow) windowZoom() {
	w.zoom() // FIXME> This should be removed
}

func (w *linuxWebviewWindow) close() {
	windowClose(w.window)
}

func (w *linuxWebviewWindow) zoomIn() {
	windowZoomIn(w.webview)
}

func (w *linuxWebviewWindow) zoomOut() {
	windowZoomOut(w.webview)
}

func (w *linuxWebviewWindow) zoomReset() {
	windowZoomSet(w.webview, 1.0)
}

func (w *linuxWebviewWindow) reload() {
	windowReload(w.webview, "wails://")
}

func (w *linuxWebviewWindow) forceReload() {
	w.reload()
}

func (w *linuxWebviewWindow) center() {
	x, y, width, height, _ := windowGetCurrentMonitorGeometry(w.window)
	if x == -1 && y == -1 && width == -1 && height == -1 {
		return
	}
	windowWidth, windowHeight := windowGetSize(w.window)

	newX := ((width - int(windowWidth)) / 2) + x
	newY := ((height - int(windowHeight)) / 2) + y

	// Place the window at the center of the monitor
	windowMove(w.window, newX, newY)
}

func (w *linuxWebviewWindow) isMinimised() bool {
	return windowIsMinimized(w.window)
}

func (w *linuxWebviewWindow) isMaximised() bool {
	return windowIsMaximized(w.window)
}

func (w *linuxWebviewWindow) isFocused() bool {
	return windowIsFocused(w.window)
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	return windowIsFullscreen(w.window)
}

func (w *linuxWebviewWindow) restore() {
	// restore window to normal size
	// FIXME: never called!  - remove from webviewImpl interface
}

func (w *linuxWebviewWindow) execJS(js string) {
	windowExecJS(w.webview, js)
}

func (w *linuxWebviewWindow) setURL(uri string) {
	if uri != "" {
		url, err := url.Parse(uri)
		if err == nil && url.Scheme == "" && url.Host == "" {
			// TODO handle this in a central location, the scheme and host might be platform dependant.
			url.Scheme = "wails"
			url.Host = "wails"
			uri = url.String()
		}
	}
	windowSetURL(w.webview, uri)
}

func (w *linuxWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	windowSetKeepAbove(w.window, alwaysOnTop)
}

func newWindowImpl(parent *WebviewWindow) *linuxWebviewWindow {
	//	(*C.struct__GtkWidget)(m.native)
	//var menubar *C.struct__GtkWidget
	return &linuxWebviewWindow{
		application: getNativeApplication().application,
		parent:      parent,
		//		menubar:     menubar,
	}
}

func (w *linuxWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		windowSetTitle(w.window, title)
	}
}

func (w *linuxWebviewWindow) setSize(width, height int) {
	windowResize(w.window, width, height)
}

func (w *linuxWebviewWindow) setMinMaxSize(minWidth, minHeight, maxWidth, maxHeight int) {
	if minWidth == 0 {
		minWidth = -1
	}
	if minHeight == 0 {
		minHeight = -1
	}
	if maxWidth == 0 {
		maxWidth = -1
	}
	if maxHeight == 0 {
		maxHeight = -1
	}
	windowSetGeometryHints(w.window, minWidth, minHeight, maxWidth, maxHeight)
}

func (w *linuxWebviewWindow) setMinSize(width, height int) {
	w.setMinMaxSize(width, height, w.parent.options.MaxWidth, w.parent.options.MaxHeight)
}

func (w *linuxWebviewWindow) setMaxSize(width, height int) {
	w.setMinMaxSize(w.parent.options.MinWidth, w.parent.options.MinHeight, width, height)
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	windowSetResizable(w.window, resizable)
}

func (w *linuxWebviewWindow) toggleDevTools() {
	showDevTools(w.webview)
}

func (w *linuxWebviewWindow) size() (int, int) {
	return windowGetSize(w.window)
}

func (w *linuxWebviewWindow) setRelativePosition(x, y int) {
	mx, my, _, _, _ := windowGetCurrentMonitorGeometry(w.window)
	windowMove(w.window, x+mx, y+my)
}

func (w *linuxWebviewWindow) width() int {
	width, _ := w.size()
	return width
}

func (w *linuxWebviewWindow) height() int {
	_, height := w.size()
	return height
}

func (w *linuxWebviewWindow) setAbsolutePosition(x int, y int) {
	// Set the window's absolute position
	windowMove(w.window, x, y)
}

func (w *linuxWebviewWindow) absolutePosition() (int, int) {
	var x, y int
	x, y = windowGetAbsolutePosition(w.window)
	return x, y
}

func (w *linuxWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}

	app := getNativeApplication()

	menu := app.getApplicationMenu()
	w.window, w.webview, w.vbox = windowNew(app.application, menu, w.parent.id, 1)
	app.registerWindow(w.window, w.parent.id) // record our mapping
	w.connectSignals()
	if w.parent.options.EnableDragAndDrop {
		w.enableDND()
	}
	w.setTitle(w.parent.options.Title)
	w.setAlwaysOnTop(w.parent.options.AlwaysOnTop)
	w.setResizable(!w.parent.options.DisableResize)
	// only set min/max size if actually set
	if w.parent.options.MinWidth != 0 &&
		w.parent.options.MinHeight != 0 &&
		w.parent.options.MaxWidth != 0 &&
		w.parent.options.MaxHeight != 0 {
		w.setMinMaxSize(
			w.parent.options.MinWidth,
			w.parent.options.MinHeight,
			w.parent.options.MaxWidth,
			w.parent.options.MaxHeight,
		)
	}
	w.setSize(w.parent.options.Width, w.parent.options.Height)
	w.setZoom(w.parent.options.Zoom)
	if w.parent.options.BackgroundType != BackgroundTypeSolid {
		w.setTransparent()
		w.setBackgroundColour(w.parent.options.BackgroundColour)
	}

	w.setFrameless(w.parent.options.Frameless)

	if w.parent.options.X != 0 || w.parent.options.Y != 0 {
		w.setRelativePosition(w.parent.options.X, w.parent.options.Y)
	} else {
		fmt.Println("attempting to set in the center")
		w.center()
	}
	switch w.parent.options.StartState {
	case WindowStateMaximised:
		w.maximise()
	case WindowStateMinimised:
		w.minimise()
	case WindowStateFullscreen:
		w.fullscreen()
	}

	if w.parent.options.URL != "" {
		w.setURL(w.parent.options.URL)
	}
	// We need to wait for the HTML to load before we can execute the javascript
	// FIXME: What event is this?  DomReady?
	w.parent.On(events.Mac.WebViewDidFinishNavigation, func(_ *WindowEvent) {
		if w.parent.options.JS != "" {
			w.execJS(w.parent.options.JS)
		}
		if w.parent.options.CSS != "" {
			js := fmt.Sprintf("(function() { var style = document.createElement('style'); style.appendChild(document.createTextNode('%s')); document.head.appendChild(style); })();", w.parent.options.CSS)
			w.execJS(js)
		}
	})
	if w.parent.options.HTML != "" {
		w.setHTML(w.parent.options.HTML)
	}
	if !w.parent.options.Hidden {
		w.show()
		if w.parent.options.X != 0 || w.parent.options.Y != 0 {
			w.setRelativePosition(w.parent.options.X, w.parent.options.Y)
		} else {
			w.center() // needs to be queued until after GTK starts up!
		}
	}
	if w.parent.options.DevToolsEnabled {
		w.toggleDevTools()
	}
}

func (w *linuxWebviewWindow) setTransparent() {
	windowSetTransparent(w.window)
}

func (w *linuxWebviewWindow) setBackgroundColour(colour RGBA) {
	windowSetBackgroundColour(w.vbox, w.webview, colour)
}

func (w *linuxWebviewWindow) relativePosition() (int, int) {
	var x, y int
	x, y = windowGetRelativePosition(w.window)
	return x, y
}

func (w *linuxWebviewWindow) destroy() {
	windowDestroy(w.window)
}

func (w *linuxWebviewWindow) setEnabled(enabled bool) {
	widgetSetSensitive(w.window, enabled)
}

func (w *linuxWebviewWindow) setHTML(html string) {
	windowSetHTML(w.webview, html)
}

func (w *linuxWebviewWindow) startResize(border string) error {
	// FIXME: what do we need to do here?
	return nil
}

func (w *linuxWebviewWindow) nativeWindowHandle() uintptr {
	return uintptr(w.window)
}

func (w *linuxWebviewWindow) print() error {
	w.execJS("window.print();")
	return nil
}

func (w *linuxWebviewWindow) handleKeyEvent(acceleratorString string) {
	// Parse acceleratorString
	// accelerator, err := parseAccelerator(acceleratorString)
	// if err != nil {
	// 	globalApplication.error("unable to parse accelerator: %s", err.Error())
	// 	return
	// }
	w.parent.processKeyBinding(acceleratorString)
}
