//go:build linux

package application

import (
	"fmt"
	"net/url"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var showDevTools = func(window unsafe.Pointer) {}

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
	menu        *menu.Menu
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
	fmt.Println("endDrag", button, x, y)
}

func (w *linuxWebviewWindow) enableDND() {
	windowEnableDND(w.parent.id, w.webview)
}

func (w *linuxWebviewWindow) connectSignals() {
	windowSetupSignalHandlers(w.parent.id, w.window, w.webview, w.parent.options.HideOnClose)
}

func (w *linuxWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	fmt.Println("linux.openContextMenu() - not implemented")
	/*	void
		gtk_menu_popup_at_rect (
		  GtkMenu* menu,
		  GdkWindow* rect_window,
		  const GdkRectangle* rect,
		  GdkGravity rect_anchor,
		  GdkGravity menu_anchor,
		  const GdkEvent* trigger_event
		)
	*/
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
	globalApplication.dispatchOnMainThread(func() {
		windowPresent(w.window)
	})
}

func (w *linuxWebviewWindow) show() {
	globalApplication.dispatchOnMainThread(func() {
		windowShow(w.window)
	})
}

func (w *linuxWebviewWindow) hide() {
	globalApplication.dispatchOnMainThread(func() {
		windowHide(w.window)
	})
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
	fmt.Println("unfullscreen")
	globalApplication.dispatchOnMainThread(func() {
		windowUnfullscreen(w.window)
		w.unmaximise()
	})
}

func (w *linuxWebviewWindow) fullscreen() {
	w.maximise()
	w.lastWidth, w.lastHeight = w.size()
	globalApplication.dispatchOnMainThread(func() {
		x, y, width, height, scale := windowGetCurrentMonitorGeometry(w.window)
		if x == -1 && y == -1 && width == -1 && height == -1 {
			return
		}
		w.setMinMaxSize(0, 0, width*scale, height*scale)
		w.setSize(width*scale, height*scale)
		windowFullscreen(w.window)
		w.setPosition(0, 0)
	})
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
	if !w.parent.options.HideOnClose {
		globalApplication.deleteWindowByID(w.parent.id)
	}
}

func (w *linuxWebviewWindow) zoomIn() {
	windowZoomIn(w.webview)
}

func (w *linuxWebviewWindow) zoomOut() {
	windowZoomOut(w.webview)
}

func (w *linuxWebviewWindow) zoomReset() {
	windowZoomSet(w.webview, 0.0)
}

func (w *linuxWebviewWindow) reload() {
	windowReload(w.webview, "wails://")
}

func (w *linuxWebviewWindow) forceReload() {
	w.reload()
}

func (w *linuxWebviewWindow) center() {
	globalApplication.dispatchOnMainThread(func() {
		x, y, width, height, _ := windowGetCurrentMonitorGeometry(w.window)
		if x == -1 && y == -1 && width == -1 && height == -1 {
			return
		}
		windowWidth, windowHeight := windowGetSize(w.window)

		newX := ((width - int(windowWidth)) / 2) + x
		newY := ((height - int(windowHeight)) / 2) + y

		// Place the window at the center of the monitor
		windowMove(w.window, newX, newY)
	})
}

func (w *linuxWebviewWindow) isMinimised() bool {
	return windowIsMinimized(w.window)
}

func (w *linuxWebviewWindow) isMaximised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return windowIsMaximized(w.window)
	})
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return windowIsFullscreen(w.window)
	})
}

func (w *linuxWebviewWindow) syncMainThreadReturningBool(fn func() bool) bool {
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
	windowToggleDevTools(w.webview)
}

func (w *linuxWebviewWindow) size() (int, int) {
	/*	var width, height C.int
		var wg sync.WaitGroup
		wg.Add(1)
		globalApplication.dispatchOnMainThread(func() {

			C.gtk_window_get_size((*C.GtkWindow)(w.window), &width, &height)
			wg.Done()
		})
		wg.Wait()
		return int(width), int(height)
	*/
	// Does this need to be guarded?
	return windowGetSize(w.window)
}

func (w *linuxWebviewWindow) setPosition(x, y int) {
	mx, my, _, _, _ := windowGetCurrentMonitorGeometry(w.window)
	globalApplication.dispatchOnMainThread(func() {
		windowMove(w.window, x+mx, y+my)
	})
}

func (w *linuxWebviewWindow) width() int {
	width, _ := w.size()
	return width
}

func (w *linuxWebviewWindow) height() int {
	_, height := w.size()
	return height
}

func (w *linuxWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}

	app := getNativeApplication()
	menu := app.applicationMenu

	globalApplication.dispatchOnMainThread(func() {
		w.window, w.webview = windowNew(app.application, menu, w.parent.id, 1)
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
		w.setBackgroundColour(w.parent.options.BackgroundColour)
		w.setFrameless(w.parent.options.Frameless)

		if w.parent.options.X != 0 || w.parent.options.Y != 0 {
			w.setPosition(w.parent.options.X, w.parent.options.Y)
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
		w.parent.On(events.Mac.WebViewDidFinishNavigation, func(_ *WindowEventContext) {
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
				w.setPosition(w.parent.options.X, w.parent.options.Y)
			} else {
				w.center() // needs to be queued until after GTK starts up!
			}
		}
	})
}

func (w *linuxWebviewWindow) setTransparent() {
	windowSetTransparent(w.window)
}

func (w *linuxWebviewWindow) setBackgroundColour(colour RGBA) {
	if colour.Alpha != 0 {
		w.setTransparent()
	}
	windowSetBackgroundColour(w.webview, colour)
}

func (w *linuxWebviewWindow) position() (int, int) {
	var x, y int
	var wg sync.WaitGroup
	wg.Add(1)
	go globalApplication.dispatchOnMainThread(func() {
		x, y = windowGetPosition(w.window)
		wg.Done()
	})
	wg.Wait()
	return x, y
}

func (w *linuxWebviewWindow) destroy() {
	windowDestroy(w.window)
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
