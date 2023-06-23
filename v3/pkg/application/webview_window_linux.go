//go:build linux && !purego

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>


// exported below
extern gboolean buttonEvent(GtkWidget *widget, GdkEventButton *event, gpointer user_data);
   extern void processRequest(void *request, gpointer user_data);
extern void onDragNDrop(
   void         *target,
   GdkDragContext* context,
   gint         x,
   gint         y,
   gpointer     seldata,
   guint        info,
   guint        time,
   gpointer     data);
// exported below (end)

static void signal_connect(GtkWidget *widget, char *event, void *cb, void* data) {
   // g_signal_connect is a macro and can't be called directly
   g_signal_connect(widget, event, cb, data);
}
*/
import "C"

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var showDevTools = func(window unsafe.Pointer) {}

func gtkBool(input bool) C.gboolean {
	if input {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

type dragInfo struct {
	XRoot       int
	YRoot       int
	DragTime    int
	MouseButton uint
}

type linuxWebviewWindow struct {
	id          uint
	application unsafe.Pointer
	window      unsafe.Pointer
	webview     unsafe.Pointer
	parent      *WebviewWindow
	menubar     *C.GtkWidget
	vbox        *C.GtkWidget
	menu        *menu.Menu
	accels      *C.GtkAccelGroup
	lastWidth   int
	lastHeight  int
	drag        dragInfo
}

var (
	registered bool = false // avoid 'already registered message' about 'wails://'
)

//export buttonEvent
func buttonEvent(_ *C.GtkWidget, event *C.GdkEventButton, data unsafe.Pointer) C.gboolean {
	// Constants (defined here to be easier to use with )
	GdkButtonPress := C.GDK_BUTTON_PRESS     // 4
	Gdk2ButtonPress := C.GDK_2BUTTON_PRESS   // 5 for double-click
	GdkButtonRelease := C.GDK_BUTTON_RELEASE // 7

	windowId := uint(*((*C.uint)(data)))
	window := globalApplication.getWindowForID(windowId)
	if window == nil {
		return C.gboolean(0)
	}
	lw, ok := (window.impl).(*linuxWebviewWindow)
	if !ok {
		return C.gboolean(0)
	}

	if event == nil {
		return C.gboolean(0)
	}
	if event.button == 3 {
		return C.gboolean(0)
	}

	switch int(event._type) {
	case GdkButtonPress:
		lw.startDrag(uint(event.button), int(event.x_root), int(event.y_root))
	case Gdk2ButtonPress:
		fmt.Printf("%d - button %d - double-clicked\n", windowId, int(event.button))
	case GdkButtonRelease:
		lw.endDrag(uint(event.button), int(event.x_root), int(event.y_root))
	}

	return C.gboolean(0)
}

func (w *linuxWebviewWindow) startDrag(button uint, x, y int) {
	fmt.Println("startDrag ", button, x, y)
	w.drag.XRoot = x
	w.drag.YRoot = y
}

func (w *linuxWebviewWindow) endDrag(button uint, x, y int) {
	fmt.Println("endDrag", button, x, y)
}

//export onDragNDrop
func onDragNDrop(target unsafe.Pointer, context *C.GdkDragContext, x C.gint, y C.gint, seldata unsafe.Pointer, info C.guint, time C.guint, data unsafe.Pointer) {
	fmt.Println("target", target, info)
	var length C.gint
	selection := unsafe.Pointer(C.gtk_selection_data_get_data_with_length((*C.GtkSelectionData)(seldata), &length))
	extracted := C.g_uri_list_extract_uris((*C.char)(selection))
	defer C.g_strfreev(extracted)

	uris := unsafe.Slice(
		(**C.char)(unsafe.Pointer(extracted)),
		int(length))

	var filenames []string
	for _, uri := range uris {
		if uri == nil {
			break
		}
		filenames = append(filenames, strings.TrimPrefix(C.GoString(uri), "file://"))
	}
	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:  uint(*((*C.uint)(data))),
		filenames: filenames,
	}
	C.gtk_drag_finish(context, C.true, C.false, time)
}

//export processRequest
func processRequest(request unsafe.Pointer, data unsafe.Pointer) {
	windowId := uint(*((*C.uint)(data)))
	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(request),
		windowId:   windowId,
		windowName: globalApplication.getWindowForID(windowId).Name(),
	}
}

func (w *linuxWebviewWindow) enableDND() {
	dnd := C.CString("text/uri-list")
	defer C.free(unsafe.Pointer(dnd))
	targetentry := C.gtk_target_entry_new(dnd, 0, C.guint(w.parent.id))
	defer C.gtk_target_entry_free(targetentry)
	C.gtk_drag_dest_set((*C.GtkWidget)(w.webview), C.GTK_DEST_DEFAULT_DROP, targetentry, 1, C.GDK_ACTION_COPY)
	event := C.CString("drag-data-received")
	defer C.free(unsafe.Pointer(event))
	id := C.uint(w.parent.id)
	C.signal_connect((*C.GtkWidget)(unsafe.Pointer(w.webview)), event, C.onDragNDrop, unsafe.Pointer(C.gpointer(&id)))
}

func (w *linuxWebviewWindow) newWebview(gpuPolicy int) unsafe.Pointer {
	manager := C.webkit_user_content_manager_new()
	external := C.CString("external")
	C.webkit_user_content_manager_register_script_message_handler(manager, external)

	C.free(unsafe.Pointer(external))
	webview := C.webkit_web_view_new_with_user_content_manager(manager)
	id := C.uint(w.parent.id)
	if !registered {
		wails := C.CString("wails")
		C.webkit_web_context_register_uri_scheme(
			C.webkit_web_context_get_default(),
			wails,
			C.WebKitURISchemeRequestCallback(C.processRequest),
			C.gpointer(&id),
			nil)
		registered = true
		C.free(unsafe.Pointer(wails))
	}
	settings := C.webkit_web_view_get_settings((*C.WebKitWebView)(unsafe.Pointer(webview)))
	wails_io := C.CString("wails.io")
	empty := C.CString("")
	defer C.free(unsafe.Pointer(wails_io))
	defer C.free(unsafe.Pointer(empty))
	C.webkit_settings_set_user_agent_with_application_details(settings, wails_io, empty)

	switch gpuPolicy {
	case 0:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
		break
	case 1:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND)
		break
	case 2:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER)
		break
	default:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND)
	}
	return unsafe.Pointer(webview)
}

func (w *linuxWebviewWindow) connectSignals() {
	event := C.CString("delete-event")
	defer C.free(unsafe.Pointer(event))

	// Window close handler

	if w.parent.options.HideOnClose {
		C.signal_connect((*C.GtkWidget)(w.window), event, C.gtk_widget_hide_on_delete, C.NULL)
	} else {

		//		C.signal_connect((*C.GtkWidget)(window), event, C.close_button_pressed, w.parent.id)
	}
	/*
		event = C.CString("load-changed")
		defer C.free(unsafe.Pointer(event))
		C.signal_connect(webview, event, C.webviewLoadChanged, unsafe.Pointer(&w.parent.id))
	*/
	id := C.uint(w.parent.id)
	event = C.CString("button-press-event")
	C.signal_connect((*C.GtkWidget)(unsafe.Pointer(w.webview)), event, C.buttonEvent, unsafe.Pointer(&id))
	C.free(unsafe.Pointer(event))
	event = C.CString("button-release-event")
	defer C.free(unsafe.Pointer(event))
	C.signal_connect((*C.GtkWidget)(unsafe.Pointer(w.webview)), event, C.buttonEvent, unsafe.Pointer(&id))
}

func (w *linuxWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	fmt.Println("linux.openContextMenu()")
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
	return float64(C.webkit_web_view_get_zoom_level((*C.WebKitWebView)(w.webview)))
}

func (w *linuxWebviewWindow) setZoom(zoom float64) {
	C.webkit_web_view_set_zoom_level((*C.WebKitWebView)(w.webview), C.double(zoom))
}

func (w *linuxWebviewWindow) setFrameless(frameless bool) {
	if frameless {
		C.gtk_window_set_decorated((*C.GtkWindow)(w.window), C.gboolean(0))
	} else {
		C.gtk_window_set_decorated((*C.GtkWindow)(w.window), C.gboolean(1))
		// TODO: Deal with transparency for the titlebar if possible
		//       Perhaps we just make it undecorated and add a menu bar inside?
	}
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

func (w *linuxWebviewWindow) show() {
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_widget_show_all((*C.GtkWidget)(w.window))
	})
}

func (w *linuxWebviewWindow) hide() {
	C.gtk_widget_hide((*C.GtkWidget)(w.window))
}

func (w *linuxWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}

func (w *linuxWebviewWindow) isVisible() bool {
	if C.gtk_widget_is_visible((*C.GtkWidget)(w.window)) == 1 {
		return true
	}
	return false
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
	fmt.Println("unfullscreen")
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_unfullscreen((*C.GtkWindow)(w.window))
		w.unmaximise()
	})
}

func (w *linuxWebviewWindow) fullscreen() {
	w.maximise()
	w.lastWidth, w.lastHeight = w.size()
	globalApplication.dispatchOnMainThread(func() {
		x, y, width, height, scale := w.getCurrentMonitorGeometry()
		if x == -1 && y == -1 && width == -1 && height == -1 {
			return
		}
		w.setMinMaxSize(0, 0, width*scale, height*scale)
		w.setSize(width*scale, height*scale)
		C.gtk_window_fullscreen((*C.GtkWindow)(w.window))
		w.setRelativePosition(0, 0)
	})
}

func (w *linuxWebviewWindow) setEnabled(enabled bool) {
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_widget_set_sensitive((*C.GtkWidget)(w.window), C.gboolean(enabled))
	})
}

func (w *linuxWebviewWindow) unminimise() {
	C.gtk_window_present((*C.GtkWindow)(w.window))
	// gtk_window_unminimize ((*C.GtkWindow)(w.window)) /// gtk4
}

func (w *linuxWebviewWindow) unmaximise() {
	C.gtk_window_unmaximize((*C.GtkWindow)(w.window))
}

func (w *linuxWebviewWindow) maximise() {
	C.gtk_window_maximize((*C.GtkWindow)(w.window))
}

func (w *linuxWebviewWindow) minimise() {
	C.gtk_window_iconify((*C.GtkWindow)(w.window))
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
	C.gtk_window_close((*C.GtkWindow)(w.window))
	if !w.parent.options.HideOnClose {
		globalApplication.deleteWindowByID(w.parent.id)
	}
}

func (w *linuxWebviewWindow) zoomIn() {
	lvl := C.webkit_web_view_get_zoom_level((*C.WebKitWebView)(w.webview))
	C.webkit_web_view_set_zoom_level((*C.WebKitWebView)(w.webview), lvl+0.5)
}

func (w *linuxWebviewWindow) zoomOut() {
	lvl := C.webkit_web_view_get_zoom_level((*C.WebKitWebView)(w.webview))
	C.webkit_web_view_set_zoom_level((*C.WebKitWebView)(w.webview), lvl-0.5)
}

func (w *linuxWebviewWindow) zoomReset() {
	C.webkit_web_view_set_zoom_level((*C.WebKitWebView)(w.webview), 0.0)
}

func (w *linuxWebviewWindow) reload() {
	// TODO: This should be a constant somewhere I feel
	uri := C.CString("wails://")
	C.webkit_web_view_load_uri((*C.WebKitWebView)(w.window), uri)
	C.free(unsafe.Pointer(uri))
}

func (w *linuxWebviewWindow) forceReload() {
	w.reload()
}

func (w linuxWebviewWindow) getCurrentMonitor() *C.GdkMonitor {
	// Get the monitor that the window is currently on
	display := C.gtk_widget_get_display((*C.GtkWidget)(w.window))
	gdk_window := C.gtk_widget_get_window((*C.GtkWidget)(w.window))
	if gdk_window == nil {
		return nil
	}
	return C.gdk_display_get_monitor_at_window(display, gdk_window)
}

func (w linuxWebviewWindow) getCurrentMonitorGeometry() (x int, y int, width int, height int, scale int) {
	monitor := w.getCurrentMonitor()
	if monitor == nil {
		return -1, -1, -1, -1, 1
	}
	var result C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &result)
	scale = int(C.gdk_monitor_get_scale_factor(monitor))
	return int(result.x), int(result.y), int(result.width), int(result.height), scale
}

func (w *linuxWebviewWindow) center() {
	globalApplication.dispatchOnMainThread(func() {
		x, y, width, height, _ := w.getCurrentMonitorGeometry()
		if x == -1 && y == -1 && width == -1 && height == -1 {
			return
		}

		var windowWidth C.int
		var windowHeight C.int
		C.gtk_window_get_size((*C.GtkWindow)(w.window), &windowWidth, &windowHeight)

		newX := C.int(((width - int(windowWidth)) / 2) + x)
		newY := C.int(((height - int(windowHeight)) / 2) + y)

		// Place the window at the center of the monitor
		C.gtk_window_move((*C.GtkWindow)(w.window), newX, newY)
	})
}

func (w *linuxWebviewWindow) isMinimised() bool {
	gdkwindow := C.gtk_widget_get_window((*C.GtkWidget)(w.window))
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_ICONIFIED > 0
}

func (w *linuxWebviewWindow) isMaximised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		gdkwindow := C.gtk_widget_get_window((*C.GtkWidget)(w.window))
		state := C.gdk_window_get_state(gdkwindow)
		return state&C.GDK_WINDOW_STATE_MAXIMIZED > 0 && state&C.GDK_WINDOW_STATE_FULLSCREEN == 0
	})
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	return w.syncMainThreadReturningBool(func() bool {
		gdkwindow := C.gtk_widget_get_window((*C.GtkWidget)(w.window))
		state := C.gdk_window_get_state(gdkwindow)
		return state&C.GDK_WINDOW_STATE_FULLSCREEN > 0
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
	value := C.CString(js)
	C.webkit_web_view_evaluate_javascript((*C.WebKitWebView)(w.webview),
		value,
		C.long(len(js)),
		nil,
		C.CString(""),
		nil,
		nil,
		nil)
	C.free(unsafe.Pointer(value))
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
	target := C.CString(uri)
	C.webkit_web_view_load_uri((*C.WebKitWebView)(w.webview), target)
	C.free(unsafe.Pointer(target))
}

func (w *linuxWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	C.gtk_window_set_keep_above((*C.GtkWindow)(w.window), gtkBool(alwaysOnTop))
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
		cTitle := C.CString(title)
		C.gtk_window_set_title((*C.GtkWindow)(w.window), cTitle)
		C.free(unsafe.Pointer(cTitle))
	}
}

func (w *linuxWebviewWindow) setSize(width, height int) {
	C.gtk_window_resize((*C.GtkWindow)(w.window), C.gint(width), C.gint(height))
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
	size := C.GdkGeometry{
		min_width:  C.int(minWidth),
		min_height: C.int(minHeight),
		max_width:  C.int(maxWidth),
		max_height: C.int(maxHeight),
	}
	C.gtk_window_set_geometry_hints((*C.GtkWindow)(w.window), nil, &size, C.GDK_HINT_MAX_SIZE|C.GDK_HINT_MIN_SIZE)
}

func (w *linuxWebviewWindow) setMinSize(width, height int) {
	w.setMinMaxSize(width, height, w.parent.options.MaxWidth, w.parent.options.MaxHeight)
}

func (w *linuxWebviewWindow) setMaxSize(width, height int) {
	w.setMinMaxSize(w.parent.options.MinWidth, w.parent.options.MinHeight, width, height)
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	if resizable {
		C.gtk_window_set_resizable((*C.GtkWindow)(w.window), 1)
	} else {
		C.gtk_window_set_resizable((*C.GtkWindow)(w.window), 0)
	}
}

func (w *linuxWebviewWindow) toggleDevTools() {
	settings := C.webkit_web_view_get_settings((*C.WebKitWebView)(w.webview))
	enabled := C.webkit_settings_get_enable_developer_extras(settings)
	if enabled == C.int(0) {
		enabled = C.int(1)
	} else {
		enabled = C.int(0)
	}
	C.webkit_settings_set_enable_developer_extras(settings, enabled)
}

func (w *linuxWebviewWindow) size() (int, int) {
	var width, height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_get_size((*C.GtkWindow)(w.window), &width, &height)
		wg.Done()
	})
	wg.Wait()
	return int(width), int(height)
}

func (w *linuxWebviewWindow) setRelativePosition(x, y int) {
	mx, my, _, _, _ := w.getCurrentMonitorGeometry()
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_move((*C.GtkWindow)(w.window), C.int(x+mx), C.int(y+my))
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

func (w *linuxWebviewWindow) absolutePosition() (int, int) {
	var x, y C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_get_position((*C.GtkWindow)(w.window), &x, &y)
		wg.Done()
	})
	wg.Wait()
	return int(x), int(y)
}

func (w *linuxWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}

	app := getNativeApplication()
	menu := app.applicationMenu

	globalApplication.dispatchOnMainThread(func() {
		w.window = unsafe.Pointer(C.gtk_application_window_new((*C.GtkApplication)(w.application)))
		app.registerWindow((*C.GtkWindow)(w.window), w.parent.id) // record our mapping
		C.g_object_ref_sink(C.gpointer(w.window))
		w.webview = w.newWebview(1)
		w.connectSignals()
		if w.parent.options.EnableDragAndDrop {
			w.enableDND()
		}
		w.vbox = C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
		C.gtk_container_add((*C.GtkContainer)(w.window), w.vbox)
		if menu != nil {
			C.gtk_box_pack_start((*C.GtkBox)(unsafe.Pointer(w.vbox)), (*C.GtkWidget)(menu), 0, 0, 0)
		}
		C.gtk_box_pack_start((*C.GtkBox)(unsafe.Pointer(w.vbox)), (*C.GtkWidget)(w.webview), 1, 1, 0)

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
		w.parent.On(events.Mac.WebViewDidFinishNavigation, func(_ *WindowEventContext) {
			if w.parent.options.JS != "" {
				w.execJS(w.parent.options.JS)
			}
			if w.parent.options.CSS != "" {
				js := fmt.Sprintf("(function() { var style = document.createElement('style'); style.appendChild(document.createTextNode('%s')); document.head.appendChild(style); })();", w.parent.options.CSS)
				fmt.Println(js)
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
				fmt.Println("attempting to set in the center")
				w.center()
			}
		}
	})
}

func (w *linuxWebviewWindow) setTransparent() {
	screen := C.gtk_widget_get_screen((*C.GtkWidget)(w.window))
	visual := C.gdk_screen_get_rgba_visual(screen)

	if visual != nil && C.gdk_screen_is_composited(screen) == C.int(1) {
		C.gtk_widget_set_app_paintable((*C.GtkWidget)(w.window), C.gboolean(1))
		C.gtk_widget_set_visual((*C.GtkWidget)(w.window), visual)
	}
}

func (w *linuxWebviewWindow) setBackgroundColour(colour RGBA) {
	if colour.Alpha != 0 {
		w.setTransparent()
	}
	rgba := C.GdkRGBA{C.double(colour.Red) / 255.0, C.double(colour.Green) / 255.0, C.double(colour.Blue) / 255.0, C.double(colour.Alpha) / 255.0}
	fmt.Println(unsafe.Pointer(&rgba))
	C.webkit_web_view_set_background_color((*C.WebKitWebView)(w.webview), &rgba)
}

func (w *linuxWebviewWindow) relativePosition() (int, int) {
	var x, y C.int
	var wg sync.WaitGroup
	wg.Add(1)
	go globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_get_position((*C.GtkWindow)(w.window), &x, &y)

		// The position must be relative to the screen it is on
		// We need to get the screen it is on
		screen := C.gtk_widget_get_screen((*C.GtkWidget)(w.window))
		monitor := C.gdk_screen_get_monitor_at_window(screen, (*C.GdkWindow)(w.window))
		geometry := C.GdkRectangle{}
		C.gdk_screen_get_monitor_geometry(screen, monitor, &geometry)
		x = x - geometry.x
		y = y - geometry.y

		// TODO: Scale based on DPI

		wg.Done()
	})
	wg.Wait()
	return int(x), int(y)
}

func (w *linuxWebviewWindow) destroy() {
	C.gtk_window_close((*C.GtkWindow)(w.window))
}

func (w *linuxWebviewWindow) setHTML(html string) {
	cHTML := C.CString(html)
	uri := C.CString("wails://")
	empty := C.CString("")
	defer C.free(unsafe.Pointer(cHTML))
	defer C.free(unsafe.Pointer(uri))
	defer C.free(unsafe.Pointer(empty))
	C.webkit_web_view_load_alternate_html(
		(*C.WebKitWebView)(w.webview),
		cHTML,
		uri,
		empty)
}

func (w *linuxWebviewWindow) nativeWindowHandle() uintptr {
	return uintptr(w.window)
}
