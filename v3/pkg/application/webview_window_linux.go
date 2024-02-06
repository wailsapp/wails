//go:build linux

package application

import "C"
import (
	"fmt"
	"math"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/capabilities"
	"github.com/wailsapp/wails/v3/internal/runtime"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type dragInfo struct {
	XRoot       int
	YRoot       int
	DragTime    uint32
	MouseButton uint
}

type linuxWebviewWindow struct {
	id           uint
	application  pointer
	window       pointer
	webview      pointer
	parent       *WebviewWindow
	menubar      pointer
	vbox         pointer
	menu         *Menu
	accels       pointer
	lastWidth    int
	lastHeight   int
	drag         dragInfo
	lastX, lastY int
	gtkmenu      pointer
}

var (
	registered bool = false // avoid 'already registered message' about 'wails://'
)

func (w *linuxWebviewWindow) endDrag(button uint, x, y int) {
	w.drag.XRoot = 0.0
	w.drag.YRoot = 0.0
	w.drag.DragTime = 0
}

func (w *linuxWebviewWindow) connectSignals() {
	cb := func(e events.WindowEventType) {
		w.parent.emit(e)
	}
	w.setupSignalHandlers(cb)
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
	w.contextMenuShow(native, data)
}

func (w *linuxWebviewWindow) getScreen() (*Screen, error) {
	mx, my, width, height, scale := w.getCurrentMonitorGeometry()
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
	w.present()
}

func (w *linuxWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}

func (w *linuxWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	//	C.setFullscreenButtonEnabled(w.nsWindow, C.bool(enabled))
	fmt.Println("setFullscreenButtonEnabled - not implemented")
}

func (w *linuxWebviewWindow) disableSizeConstraints() {
	x, y, width, height, scale := w.getCurrentMonitorGeometry()
	w.setMinMaxSize(x, y, width*scale, height*scale)
}

func (w *linuxWebviewWindow) unfullscreen() {
	windowUnfullscreen(w.window)
	w.unmaximise()
}

func (w *linuxWebviewWindow) fullscreen() {
	w.maximise()
	//w.lastWidth, w.lastHeight = w.size()
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
	w.present()
}

func (w *linuxWebviewWindow) on(eventID uint) {
	// TODO: Test register/unregister listener for linux events
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
	getNativeApplication().unregisterWindow(windowPointer(w.window))
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
	x, y, width, height, _ := w.getCurrentMonitorGeometry()
	if x == -1 && y == -1 && width == -1 && height == -1 {
		return
	}
	windowWidth, windowHeight := w.size()

	newX := ((width - windowWidth) / 2) + x
	newY := ((height - windowHeight) / 2) + y

	// Place the window at the center of the monitor
	w.move(newX, newY)
}

func (w *linuxWebviewWindow) restore() {
	// restore window to normal size
	// FIXME: never called!  - remove from webviewImpl interface
}

func newWindowImpl(parent *WebviewWindow) *linuxWebviewWindow {
	//	(*C.struct__GtkWidget)(m.native)
	//var menubar *C.struct__GtkWidget
	result := &linuxWebviewWindow{
		application: getNativeApplication().application,
		parent:      parent,
		//		menubar:     menubar,
	}
	return result
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
		maxWidth = math.MaxInt
	}
	if maxHeight == 0 {
		maxHeight = math.MaxInt
	}
	windowSetGeometryHints(w.window, minWidth, minHeight, maxWidth, maxHeight)
}

func (w *linuxWebviewWindow) setMinSize(width, height int) {
	w.setMinMaxSize(width, height, w.parent.options.MaxWidth, w.parent.options.MaxHeight)
}

func (w *linuxWebviewWindow) setMaxSize(width, height int) {
	w.setMinMaxSize(w.parent.options.MinWidth, w.parent.options.MinHeight, width, height)
}

func (w *linuxWebviewWindow) setRelativePosition(x, y int) {
	mx, my, _, _, _ := w.getCurrentMonitorGeometry()
	w.move(x+mx, y+my)
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
	w.move(x, y)
}

func (w *linuxWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}

	// Register the capabilities
	globalApplication.capabilities = capabilities.NewCapabilities()

	app := getNativeApplication()

	var menu = w.menu
	if menu == nil && globalApplication.ApplicationMenu != nil {
		menu = globalApplication.ApplicationMenu.clone()
	}
	if menu != nil {
		InvokeSync(func() {
			menu.Update()
		})
		w.menu = menu
		w.gtkmenu = (menu.impl).(*linuxMenu).native
	}

	w.window, w.webview, w.vbox = windowNew(app.application, w.gtkmenu, w.parent.id, w.parent.options.Linux.WebviewGpuPolicy)
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
		w.center()
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

	//if w.parent.options.IgnoreMouseEvents {
	//	windowIgnoreMouseEvents(w.window, w.webview, true)
	//}

	startURL, err := assetserver.GetStartURL(w.parent.options.URL)
	if err != nil {
		globalApplication.fatal(err.Error())
	}

	w.setURL(startURL)
	w.parent.On(events.Linux.WindowLoadChanged, func(_ *WindowEvent) {
		if w.parent.options.JS != "" {
			w.execJS(w.parent.options.JS)
		}
		if w.parent.options.CSS != "" {
			js := fmt.Sprintf("(function() { var style = document.createElement('style'); style.appendChild(document.createTextNode('%s')); document.head.appendChild(style); })();", w.parent.options.CSS)
			w.execJS(js)
		}
	})
	w.parent.On(events.Linux.WindowDeleteEvent, func(e *WindowEvent) {
		w.parent.emit(events.Common.WindowClosing)
	})
	w.parent.RegisterHook(events.Linux.WindowLoadChanged, func(e *WindowEvent) {
		w.execJS(runtime.Core())
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
	if w.parent.options.DevToolsEnabled || globalApplication.isDebugMode {
		w.enableDevTools()
		if w.parent.options.OpenInspectorOnStartup {
			w.openDevTools()
		}
	}
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
