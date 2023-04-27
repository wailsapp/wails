//go:build linux && purego

package application

import (
	"fmt"
	"net/url"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var (
	registered bool = false // avoid 'already registered message'
)

const (
	// https://gitlab.gnome.org/GNOME/gtk/-/blob/gtk-3-24/gdk/gdkwindow.h#L121
	GDK_HINT_MIN_SIZE = 1 << 1
	GDK_HINT_MAX_SIZE = 1 << 2
	// https://gitlab.gnome.org/GNOME/gtk/-/blob/gtk-3-24/gdk/gdkevents.h#L512
	GDK_WINDOW_STATE_ICONIFIED  = 1 << 1
	GDK_WINDOW_STATE_MAXIMIZED  = 1 << 2
	GDK_WINDOW_STATE_FULLSCREEN = 1 << 4
)

type GdkGeometry struct {
	minWidth   int32
	minHeight  int32
	maxWidth   int32
	maxHeight  int32
	baseWidth  int32
	baseHeight int32
	widthInc   int32
	heightInc  int32
	padding    int32
	minAspect  float64
	maxAspect  float64
	GdkGravity int32
}

type linuxWebviewWindow struct {
	application                              uintptr
	window                                   uintptr
	webview                                  uintptr
	parent                                   *WebviewWindow
	menubar                                  uintptr
	vbox                                     uintptr
	menu                                     *menu.Menu
	accels                                   uintptr
	minWidth, minHeight, maxWidth, maxHeight int
}

func (w *linuxWebviewWindow) newWebview(gpuPolicy int) uintptr {
	var newContentMgr func() uintptr
	purego.RegisterLibFunc(
		&newContentMgr,
		webkit,
		"webkit_user_content_manager_new")
	var registerScriptMessageHandler func(uintptr, string)
	purego.RegisterLibFunc(&registerScriptMessageHandler, webkit, "webkit_user_content_manager_register_script_message_handler")
	var newWebview func(uintptr) uintptr
	purego.RegisterLibFunc(&newWebview, webkit, "webkit_web_view_new_with_user_content_manager")

	manager := newContentMgr()
	registerScriptMessageHandler(manager, "external")
	webview := newWebview(manager)
	if !registered {
		var registerUriScheme func(uintptr, string, uintptr, uintptr, uintptr)
		purego.RegisterLibFunc(&registerUriScheme, webkit, "webkit_web_context_register_uri_scheme")
		cb := purego.NewCallback(func(request uintptr) {
			processURLRequest(w.parent.id, request)
		})
		var defaultContext func() uintptr
		purego.RegisterLibFunc(&defaultContext, webkit, "webkit_web_context_get_default")
		registerUriScheme(defaultContext(), "wails", cb, 0, 0)
	}

	var g_signal_connect func(uintptr, string, uintptr, uintptr, bool, int) int
	purego.RegisterLibFunc(&g_signal_connect, gtk, "g_signal_connect_data")

	loadChanged := purego.NewCallback(func(window uintptr) {
		fmt.Println("loadChanged", window)
	})
	g_signal_connect(webview, "load-changed", loadChanged, 0, false, 0)

	if g_signal_connect(webview, "button-press-event", purego.NewCallback(w.buttonPress), 0, false, 0) == 0 {
		fmt.Println("failed to connect 'button-press-event")
	}
	if g_signal_connect(webview, "button-release-event", purego.NewCallback(w.buttonRelease), 0, false, 0) == 0 {
		fmt.Println("failed to connect 'button-release-event")
	}

	handleDelete := purego.NewCallback(func(uintptr) {
		w.close()
		if !w.parent.options.HideOnClose {
			fmt.Println("Need to do more!")
		}
	})
	g_signal_connect(w.window, "delete-event", handleDelete, 0, false, 0)

	var getSettings func(uintptr) uintptr
	purego.RegisterLibFunc(&getSettings, webkit, "webkit_web_view_get_settings")
	var setSettings func(uintptr, uintptr)
	purego.RegisterLibFunc(&setSettings, webkit, "webkit_web_view_set_settings")
	var setUserAgent func(uintptr, string, string)
	purego.RegisterLibFunc(&setUserAgent, webkit, "webkit_settings_set_user_agent_with_application_details")
	settings := getSettings(webview)
	setUserAgent(settings, "wails.io", "")

	var setHWAccel func(uintptr, int)
	purego.RegisterLibFunc(&setHWAccel, webkit, "webkit_settings_set_hardware_acceleration_policy")

	setHWAccel(settings, gpuPolicy)
	setSettings(webview, settings)

	return webview
}

func (w *linuxWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	fmt.Println("linux.openContextMenu()")
	//C.windowShowMenu(w.nsWindow, thisMenu.nsMenu, C.int(data.X), C.int(data.Y))
}

func (w *linuxWebviewWindow) getZoom() float64 {
	var getZoom func(uintptr) float32
	purego.RegisterLibFunc(&getZoom, webkit, "webkit_web_view_get_zoom_level")
	return float64(getZoom(w.webview))
}

func (w *linuxWebviewWindow) setZoom(zoom float64) {
	var setZoom func(uintptr, float64)
	purego.RegisterLibFunc(&setZoom, webkit, "webkit_web_view_set_zoom_level")
	setZoom(w.webview, zoom)
}

func (w *linuxWebviewWindow) setFrameless(frameless bool) {
	var setDecorated func(uintptr, int)
	purego.RegisterLibFunc(&setDecorated, gtk, "gtk_window_set_decorated")
	decorated := 1
	if frameless {
		decorated = 0
	}
	setDecorated(w.window, decorated)
	if !frameless {
		// TODO: Deal with transparency for the titlebar if possible
		//       Perhaps we just make it undecorated and add a menu bar inside?
	}
}

func (w *linuxWebviewWindow) getScreen() (*Screen, error) {
	return getScreenForWindow(w)
}

func (w *linuxWebviewWindow) show() {
	var widgetShow func(uintptr)
	purego.RegisterLibFunc(&widgetShow, gtk, "gtk_widget_show_all")
	globalApplication.dispatchOnMainThread(func() {
		widgetShow(w.window)
	})
}

func (w *linuxWebviewWindow) hide() {
	var widgetHide func(uintptr)
	purego.RegisterLibFunc(&widgetHide, gtk, "gtk_widget_hide")
	widgetHide(w.window)
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
	var unfullScreen func(uintptr)
	purego.RegisterLibFunc(&unfullScreen, gtk, "gtk_window_unfullscreen")

	globalApplication.dispatchOnMainThread(func() {
		unfullScreen(w.window)
		w.unmaximise()
	})
}

func (w *linuxWebviewWindow) fullscreen() {
	var fullScreen func(uintptr)
	purego.RegisterLibFunc(&fullScreen, gtk, "gtk_window_fullscreen")

	globalApplication.dispatchOnMainThread(func() {
		w.maximise()
		//		w.lastWidth, w.lastHeight = w.size() // do we need this?

		x, y, width, height, scale := w.getCurrentMonitorGeometry()
		if x == -1 && y == -1 && width == -1 && height == -1 {
			return
		}
		w.setMinMaxSize(0, 0, width*scale, height*scale)
		w.setSize(width*scale, height*scale)
		w.setPosition(0, 0)
		fullScreen(w.window)
	})
}

func (w *linuxWebviewWindow) unminimise() {
	var present func(uintptr)
	purego.RegisterLibFunc(&present, gtk, "gtk_window_present")
	present(w.window)
}

func (w *linuxWebviewWindow) unmaximise() {
	var unmaximize func(uintptr)
	purego.RegisterLibFunc(&unmaximize, gtk, "gtk_window_unmaximize")
	unmaximize(w.window)
}

func (w *linuxWebviewWindow) maximise() {
	var maximize func(uintptr)
	purego.RegisterLibFunc(&maximize, gtk, "gtk_window_maximize")
	maximize(w.window)
}

func (w *linuxWebviewWindow) minimise() {
	var iconify func(uintptr)
	purego.RegisterLibFunc(&iconify, gtk, "gtk_window_iconify")
	iconify(w.window)
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
	w.zoom()
}

func (w *linuxWebviewWindow) close() {
	var close func(uintptr)
	purego.RegisterLibFunc(&close, gtk, "gtk_window_close")
	close(w.window)
}

func (w *linuxWebviewWindow) zoomIn() {
	var getZoom func(uintptr) float32
	purego.RegisterLibFunc(&getZoom, webkit, "webkit_web_view_get_zoom_level")
	var setZoom func(uintptr, float32)
	purego.RegisterLibFunc(&setZoom, webkit, "webkit_web_view_set_zoom_level")
	lvl := getZoom(w.webview)
	setZoom(w.webview, lvl+0.5)
}

func (w *linuxWebviewWindow) zoomOut() {
	var getZoom func(uintptr) float32
	purego.RegisterLibFunc(&getZoom, webkit, "webkit_web_view_get_zoom_level")
	var setZoom func(uintptr, float32)
	purego.RegisterLibFunc(&setZoom, webkit, "webkit_web_view_set_zoom_level")
	lvl := getZoom(w.webview)
	setZoom(w.webview, lvl-0.5)
}

func (w *linuxWebviewWindow) zoomReset() {
	var setZoom func(uintptr, float32)
	purego.RegisterLibFunc(&setZoom, webkit, "webkit_web_view_set_zoom_level")
	setZoom(w.webview, 0.0)
}

func (w *linuxWebviewWindow) toggleDevTools() {
	var getSettings func(uintptr) uintptr
	purego.RegisterLibFunc(&getSettings, webkit, "webkit_web_view_get_settings")
	var isEnabled func(uintptr) bool
	purego.RegisterLibFunc(&isEnabled, webkit, "webkit_settings_get_enable_developer_extras")
	var enableDev func(uintptr, bool)
	purego.RegisterLibFunc(&enableDev, webkit, "webkit_settings_set_enable_developer_extras")
	settings := getSettings(w.webview)
	enabled := isEnabled(settings)
	enableDev(settings, !enabled)
}

func (w *linuxWebviewWindow) reload() {
	var reload func(uintptr)
	purego.RegisterLibFunc(&reload, webkit, "webkit_web_view_reload")
	reload(w.webview)
}

func (w *linuxWebviewWindow) forceReload() {
	var reload func(uintptr)
	purego.RegisterLibFunc(&reload, webkit, "webkit_web_view_reload_bypass_cache")
	reload(w.webview)
}

func (w linuxWebviewWindow) getCurrentMonitor() uintptr {
	var getDisplay func(uintptr) uintptr
	purego.RegisterLibFunc(&getDisplay, gtk, "gtk_widget_get_display")
	var getWindow func(uintptr) uintptr
	purego.RegisterLibFunc(&getWindow, gtk, "gtk_widget_get_window")
	var getMonitor func(uintptr, uintptr) uintptr
	purego.RegisterLibFunc(&getMonitor, gtk, "gdk_display_get_monitor_at_window")

	display := getDisplay(w.window)
	window := getWindow(w.window)
	if window == 0 {
		return 0
	}
	return getMonitor(display, window)
}

func (w linuxWebviewWindow) getCurrentMonitorGeometry() (x int, y int, width int, height int, scale int) {
	var getGeometry func(uintptr, uintptr)
	purego.RegisterLibFunc(&getGeometry, gtk, "gdk_monitor_get_geometry")
	var getScaleFactor func(uintptr) int
	purego.RegisterLibFunc(&getScaleFactor, gtk, "gdk_monitor_get_scale_factor")

	monitor := w.getCurrentMonitor()
	if monitor == 0 {
		return -1, -1, -1, -1, 1
	}
	result := struct {
		x      int32
		y      int32
		width  int32
		height int32
	}{}
	getGeometry(monitor, uintptr(unsafe.Pointer(&result)))
	scale = getScaleFactor(monitor)
	return int(result.x), int(result.y), int(result.width), int(result.height), scale
}

func (w *linuxWebviewWindow) center() {
	x, y, width, height, _ := w.getCurrentMonitorGeometry()
	if x == -1 && y == -1 && width == -1 && height == -1 {
		return
	}

	windowWidth, windowHeight := w.size()

	newX := ((width - int(windowWidth)) / 2) + x
	newY := ((height - int(windowHeight)) / 2) + y

	w.setPosition(newX, newY)
}

func (w *linuxWebviewWindow) isMinimised() bool {
	var getWindow func(uintptr) uintptr
	purego.RegisterLibFunc(&getWindow, gtk, "gtk_widget_get_window")
	var getWindowState func(uintptr) int
	purego.RegisterLibFunc(&getWindowState, gtk, "gdk_window_get_state")

	return w.syncMainThreadReturningBool(func() bool {
		state := getWindowState(getWindow(w.window))
		return state&GDK_WINDOW_STATE_ICONIFIED > 0
	})
}

func (w *linuxWebviewWindow) isMaximised() bool {
	var getWindow func(uintptr) uintptr
	purego.RegisterLibFunc(&getWindow, gtk, "gtk_widget_get_window")
	var getWindowState func(uintptr) int
	purego.RegisterLibFunc(&getWindowState, gtk, "gdk_window_get_state")

	return w.syncMainThreadReturningBool(func() bool {
		state := getWindowState(getWindow(w.window))
		return state&GDK_WINDOW_STATE_MAXIMIZED > 0 && state&GDK_WINDOW_STATE_FULLSCREEN == 0
	})
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	var getWindow func(uintptr) uintptr
	purego.RegisterLibFunc(&getWindow, gtk, "gtk_widget_get_window")
	var getWindowState func(uintptr) int
	purego.RegisterLibFunc(&getWindowState, gtk, "gdk_window_get_state")

	return w.syncMainThreadReturningBool(func() bool {
		state := getWindowState(getWindow(w.window))
		return state&GDK_WINDOW_STATE_FULLSCREEN > 0
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
	fmt.Println("restore")
}

func (w *linuxWebviewWindow) execJS(js string) {
	var evalJS func(uintptr, string, int, uintptr, string, uintptr, uintptr, uintptr)
	purego.RegisterLibFunc(&evalJS, webkit, "webkit_web_view_evaluate_javascript")
	evalJS(w.webview, js, len(js), 0, "", 0, 0, 0)
}

func (w *linuxWebviewWindow) setURL(uri string) {
	fmt.Println("setURL", uri)
	var loadUri func(uintptr, string)
	purego.RegisterLibFunc(&loadUri, webkit, "webkit_web_view_load_uri")

	url, err := url.Parse(uri)
	if url != nil && err == nil && url.Scheme == "" && url.Host == "" {
		// TODO handle this in a central location, the scheme and host might be platform dependant
		url.Scheme = "wails"
		url.Host = "wails"
		uri = url.String()
		loadUri(w.webview, uri)
	}
}

func (w *linuxWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	var keepAbove func(uintptr, bool)
	purego.RegisterLibFunc(&keepAbove, gtk, "gtk_window_set_keep_above")
	keepAbove(w.window, alwaysOnTop)
}

func newWindowImpl(parent *WebviewWindow) *linuxWebviewWindow {
	return &linuxWebviewWindow{
		application: (globalApplication.impl).(*linuxApp).application,
		parent:      parent,
	}
}

func (w *linuxWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		var setTitle func(uintptr, string)
		purego.RegisterLibFunc(&setTitle, gtk, "gtk_window_set_title")
		setTitle(w.window, title)
	}
}

func (w *linuxWebviewWindow) setSize(width, height int) {
	var setSize func(uintptr, int, int)
	purego.RegisterLibFunc(&setSize, gtk, "gtk_window_set_default_size")
	setSize(w.window, width, height)
}

func (w *linuxWebviewWindow) setMinMaxSize(minWidth, minHeight, maxWidth, maxHeight int) {
	fmt.Println("setMinMaxSize", minWidth, minHeight, maxWidth, maxHeight)
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
	size := GdkGeometry{
		minWidth:  int32(minWidth),
		minHeight: int32(minHeight),
		maxWidth:  int32(maxWidth),
		maxHeight: int32(maxHeight),
	}

	var setHints func(uintptr, uintptr, uintptr, int)
	purego.RegisterLibFunc(&setHints, gtk, "gtk_window_set_geometry_hints")
	setHints(w.window, 0, uintptr(unsafe.Pointer(&size)), GDK_HINT_MIN_SIZE|GDK_HINT_MAX_SIZE)
}

func (w *linuxWebviewWindow) setMinSize(width, height int) {
	w.setMinMaxSize(width, height, w.parent.options.MaxWidth, w.parent.options.MaxHeight)
}

func (w *linuxWebviewWindow) setMaxSize(width, height int) {
	w.setMinMaxSize(w.parent.options.MinWidth, w.parent.options.MinHeight, width, height)
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	var setResizable func(uintptr, int)
	purego.RegisterLibFunc(&setResizable, gtk, "gtk_window_set_resizable")
	globalApplication.dispatchOnMainThread(func() {
		if resizable {
			setResizable(w.window, 1)
		} else {
			setResizable(w.window, 0)
		}
	})
}

func (w *linuxWebviewWindow) size() (int, int) {
	var width, height int
	var windowGetSize func(uintptr, *int, *int)
	purego.RegisterLibFunc(&windowGetSize, gtk, "gtk_window_get_size")

	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		windowGetSize(w.window, &width, &height)
		wg.Done()
	})
	wg.Wait()
	return width, height
}

func (w *linuxWebviewWindow) setPosition(x, y int) {
	var windowMove func(uintptr, int, int)
	purego.RegisterLibFunc(&windowMove, gtk, "gtk_window_move")
	mx, my, _, _, _ := w.getCurrentMonitorGeometry()
	fmt.Println("setPosition", mx, my)
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

func (w *linuxWebviewWindow) buttonPress(widget uintptr, event uintptr, user_data uintptr) {
	GdkEventButton := (*byte)(unsafe.Pointer(event))
	fmt.Println("buttonpress", w.parent.id, widget, GdkEventButton, user_data)
}

func (w *linuxWebviewWindow) buttonRelease(widget uintptr, event uintptr, user_data uintptr) {
	GdkEventButton := (*byte)(unsafe.Pointer(event))
	fmt.Println("buttonrelease", w.parent.id, widget, GdkEventButton, user_data)
}

func (w *linuxWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}

	globalApplication.dispatchOnMainThread(func() {
		app := (globalApplication.impl).(*linuxApp)
		menu := app.applicationMenu
		var newWindow func(uintptr) uintptr
		purego.RegisterLibFunc(&newWindow, gtk, "gtk_application_window_new")
		var refSink func(uintptr)
		purego.RegisterLibFunc(&refSink, gtk, "g_object_ref_sink")
		var boxNew func(int, int) uintptr
		purego.RegisterLibFunc(&boxNew, gtk, "gtk_box_new")
		var containerAdd func(uintptr, uintptr)
		purego.RegisterLibFunc(&containerAdd, gtk, "gtk_container_add")
		var boxPackStart func(uintptr, uintptr, int, int, int)
		purego.RegisterLibFunc(&boxPackStart, gtk, "gtk_box_pack_start")

		var g_signal_connect func(uintptr, string, uintptr, uintptr, bool, int) int
		purego.RegisterLibFunc(&g_signal_connect, gtk, "g_signal_connect_data")

		w.window = newWindow(w.application)

		refSink(w.window)
		w.webview = w.newWebview(1)
		w.vbox = boxNew(1, 0)
		containerAdd(w.window, w.vbox)
		if menu != 0 {
			w.menubar = menu
			boxPackStart(w.vbox, menu, 0, 0, 0)
		}
		boxPackStart(w.vbox, w.webview, 1, 1, 0)

		w.setSize(w.parent.options.Width, w.parent.options.Height)
		w.setTitle(w.parent.options.Title)
		w.setAlwaysOnTop(w.parent.options.AlwaysOnTop)
		w.setResizable(!w.parent.options.DisableResize)
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
		w.setZoom(w.parent.options.Zoom)
		w.setBackgroundColour(w.parent.options.BackgroundColour)
		w.setFrameless(w.parent.options.Frameless)

		switch w.parent.options.StartState {
		case WindowStateMaximised:
			w.maximise()
		case WindowStateMinimised:
			w.minimise()
		case WindowStateFullscreen:
			w.fullscreen()

		}
		w.center()

		if w.parent.options.URL != "" {
			w.setURL(w.parent.options.URL)
		}
		// We need to wait for the HTML to load before we can execute the javascript
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
		if w.parent.options.Hidden == false {
			w.show()
			if w.parent.options.X != 0 || w.parent.options.Y != 0 {
				w.setPosition(w.parent.options.X, w.parent.options.Y)
			} else {
				fmt.Println("attempting to set in the center")
				w.center()
			}
		}
	})
}

func (w *linuxWebviewWindow) setTransparent() {
	var getScreen func(uintptr) uintptr
	purego.RegisterLibFunc(&getScreen, gtk, "gtk_widget_get_screen")
	var getVisual func(uintptr) uintptr
	purego.RegisterLibFunc(&getVisual, gtk, "gdk_screen_get_rgba_visual")
	var isComposited func(uintptr) int
	purego.RegisterLibFunc(&isComposited, gtk, "gdk_screen_is_composited")
	var setPaintable func(uintptr, int)
	purego.RegisterLibFunc(&setPaintable, gtk, "gtk_widget_set_app_paintable")
	var setVisual func(uintptr, uintptr)
	purego.RegisterLibFunc(&setVisual, gtk, "gtk_widget_set_visual")

	screen := getScreen(w.window)
	visual := getVisual(screen)
	if visual == 0 {
		return
	}
	if isComposited(screen) == 1 {
		setPaintable(w.window, 1)
		setVisual(w.window, visual)
	}
}

func (w *linuxWebviewWindow) setBackgroundColour(colour *RGBA) {
	if colour == nil {
		return
	}

	if colour.Alpha != 0 {
		w.setTransparent()
	}

	var rgbaParse func(uintptr, string) bool
	purego.RegisterLibFunc(&rgbaParse, gtk, "gdk_rgba_parse")
	var setBackgroundColor func(uintptr, uintptr)
	purego.RegisterLibFunc(&setBackgroundColor, webkit, "webkit_web_view_set_background_color")

	rgba := make([]byte, 4*8) // C.sizeof_GdkRGBA == 32
	pointer := uintptr(unsafe.Pointer(&rgba[0]))
	if !rgbaParse(
		pointer,
		fmt.Sprintf("rgba(%v,%v,%v,%v)",
			colour.Red,
			colour.Green,
			colour.Blue,
			float32(colour.Alpha)/255.0,
		)) {
		return
	}
	setBackgroundColor(w.webview, pointer)
}

func (w *linuxWebviewWindow) position() (int, int) {
	var getPosition func(uintptr, *int, *int) bool
	purego.RegisterLibFunc(&getPosition, gtk, "gtk_window_get_position")

	var x, y int
	var wg sync.WaitGroup
	wg.Add(1)
	go globalApplication.dispatchOnMainThread(func() {
		getPosition(w.window, &x, &y)
		wg.Done()
	})
	wg.Wait()
	return x, y
}

func (w *linuxWebviewWindow) destroy() {
	var close func(uintptr)
	purego.RegisterLibFunc(&close, gtk, "gtk_window_close")
	go globalApplication.dispatchOnMainThread(func() {
		close(w.window)
	})
}

func (w *linuxWebviewWindow) setHTML(html string) {
	fmt.Println("setHTML")
	var loadHTML func(uintptr, string, string, *string)
	purego.RegisterLibFunc(&loadHTML, webkit, "webkit_web_view_load_alternate_html")
	go globalApplication.dispatchOnMainThread(func() {
		loadHTML(w.webview, html, "wails://", nil)
	})
}
