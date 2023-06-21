//go:build linux && cgo

package application

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>

typedef struct CallbackID
{
    unsigned int value;
} CallbackID;

extern void dispatchOnMainThreadCallback(unsigned int);

static gboolean dispatchCallback(gpointer data) {
    struct CallbackID *args = data;
    unsigned int cid = args->value;
    dispatchOnMainThreadCallback(cid);
    free(args);

    return G_SOURCE_REMOVE;
};

static void dispatchOnMainThread(unsigned int id) {
    CallbackID *args = malloc(sizeof(CallbackID));
    args->value = id;
    g_idle_add((GSourceFunc)dispatchCallback, (gpointer)args);
}

typedef struct WindowEvent {
    uint id;
    uint event;
} WindowEvent;

// exported below
void activateLinux(gpointer data);
extern void emit(WindowEvent* data);
void handleClick(void*);
extern gboolean onButtonEvent(GtkWidget *widget, GdkEventButton *event, gpointer user_data);
extern void onDragNDrop(
   void         *target,
   GdkDragContext* context,
   gint         x,
   gint         y,
   gpointer     seldata,
   guint        info,
   guint        time,
   gpointer     data);
extern void onProcessRequest(void *request, gpointer user_data);
// exported below (end)

static void signal_connect(GtkWidget *widget, char *event, void *cb, void* data) {
   // g_signal_connect is a macro and can't be called directly
   g_signal_connect(widget, event, cb, data);
}

static void* new_message_dialog(GtkWindow *parent, const gchar *msg, int dialogType, bool hasButtons) {
   // gtk_message_dialog_new is variadic!  Can't call from cgo directly
   GtkWidget *dialog;
   int buttonMask;

   // buttons will be added after creation
   buttonMask = GTK_BUTTONS_OK;
   if (hasButtons) {
       buttonMask = GTK_BUTTONS_NONE;
   }

   dialog = gtk_message_dialog_new(
       parent,
       GTK_DIALOG_MODAL | GTK_DIALOG_DESTROY_WITH_PARENT,
	   dialogType,
	   buttonMask,
       msg);

   // g_signal_connect_swapped (dialog,
   //                           "response",
   //                           G_CALLBACK (callback),
   //                           dialog);
   return dialog;
};

extern void messageDialogCB(gint button);

typedef struct Screen {
	const char* id;
	const char* name;
	int p_width;
	int p_height;
	int width;
	int height;
	int x;
	int y;
	int w_width;
	int w_height;
	int w_x;
	int w_y;
	float scale;
	double rotation;
	bool isPrimary;
} Screen;


static int GetNumScreens(){
    return 0;
}
*/
import "C"

type windowPointer *C.GtkWindow
type identifier C.uint
type pointer unsafe.Pointer
type GSList C.GSList
type GSListPointer *GSList

var (
	nilRadioGroup       GSListPointer = nil
	gtkSignalHandlers   map[*C.GtkWidget]C.gulong
	gtkSignalToMenuItem map[*C.GtkWidget]*MenuItem
)

func init() {
	fmt.Println("linux_cgo")

	gtkSignalHandlers = map[*C.GtkWidget]C.gulong{}
	gtkSignalToMenuItem = map[*C.GtkWidget]*MenuItem{}
}

// mainthread stuff
func dispatchOnMainThread(id uint) {
	C.dispatchOnMainThread(C.uint(id))
}

//export dispatchOnMainThreadCallback
func dispatchOnMainThreadCallback(callbackID C.uint) {
	executeOnMainThread(uint(callbackID))
}

//export activateLinux
func activateLinux(data pointer) {
	// NOOP: Callback for now
}

// implementation below
func appName() string {
	name := C.g_get_application_name()
	defer C.free(unsafe.Pointer(name))
	return C.GoString(name)
}

func appNew(name string) pointer {
	nameC := C.CString(fmt.Sprintf("org.wails.%s", name))
	defer C.free(unsafe.Pointer(nameC))
	return pointer(C.gtk_application_new(nameC, C.G_APPLICATION_DEFAULT_FLAGS))
}

func appRun(app pointer) error {
	application := (*C.GApplication)(app)
	C.g_application_hold(application) // allows it to run without a window
	signal := C.CString("activate")
	C.g_signal_connect_data(C.gpointer(application), signal, C.GCallback(C.activateLinux), nil, nil, 0)
	status := C.g_application_run(application, 0, nil)
	C.g_application_release(application)
	C.g_object_unref(C.gpointer(app))

	var err error
	if status != 0 {
		err = fmt.Errorf("exit code: %d", status)
	}
	return err
}

func appDestroy(application pointer) {
	C.g_application_quit((*C.GApplication)(application))
}

func getCurrentWindowID(application pointer, windows map[windowPointer]uint) uint {
	// TODO: Add extra metadata to window and use it!
	window := (*C.GtkWindow)(C.gtk_application_get_active_window((*C.GtkApplication)(application)))
	if window == nil {
		return uint(1)
	}
	identifier, ok := windows[window]
	if ok {
		return identifier
	}
	// FIXME: Should we panic here if not found?
	return uint(1)
}

func getWindows(application pointer) []pointer {
	result := []pointer{}
	windows := C.gtk_application_get_windows((*C.GtkApplication)(application))
	for {
		result = append(result, pointer(windows.data))
		windows = windows.next
		if windows == nil {
			return result
		}
	}
}

func hideAllWindows(application pointer) {
	for _, window := range getWindows(application) {
		C.gtk_widget_hide((*C.GtkWidget)(window))
	}
}

func showAllWindows(application pointer) {
	for _, window := range getWindows(application) {
		C.gtk_widget_show_all((*C.GtkWidget)(window))
	}
}

// Menu
func menuAddSeparator(menu *Menu) {
	C.gtk_menu_shell_append(
		(*C.GtkMenuShell)((menu.impl).(*linuxMenu).native),
		C.gtk_separator_menu_item_new())
}

func menuAppend(parent *Menu, menu *MenuItem) {
	C.gtk_menu_shell_append(
		(*C.GtkMenuShell)((parent.impl).(*linuxMenu).native),
		(*C.GtkWidget)((menu.impl).(*linuxMenuItem).native),
	)
	/* gtk4
	C.gtk_menu_item_set_submenu(
		(*C.struct__GtkMenuItem)((menu.impl).(*linuxMenuItem).native),
		(*C.struct__GtkWidget)((parent.impl).(*linuxMenu).native),
	)
	*/
}

func menuBarNew() pointer {
	return pointer(C.gtk_menu_bar_new())
}

func menuNew() pointer {
	return pointer(C.gtk_menu_new())
}

func menuSetSubmenu(item *MenuItem, menu *Menu) {
	C.gtk_menu_item_set_submenu(
		(*C.GtkMenuItem)((item.impl).(*linuxMenuItem).native),
		(*C.GtkWidget)((menu.impl).(*linuxMenu).native))
}

func menuGetRadioGroup(item *linuxMenuItem) *GSList {
	return (*GSList)(C.gtk_radio_menu_item_get_group((*C.GtkRadioMenuItem)(item.native)))
}

//export handleClick
func handleClick(idPtr unsafe.Pointer) {
	id := (*C.GtkWidget)(idPtr)
	item, ok := gtkSignalToMenuItem[id]
	if !ok {
		return
	}

	switch item.itemType {
	case text, checkbox:
		menuItemClicked <- item.id
	case radio:
		menuItem := (item.impl).(*linuxMenuItem)
		if menuItem.isChecked() {
			menuItemClicked <- item.id
		}
	}
}

func attachMenuHandler(item *MenuItem) {
	signal := C.CString("activate")
	defer C.free(unsafe.Pointer(signal))

	impl := (item.impl).(*linuxMenuItem)
	widget := impl.native
	flags := C.GConnectFlags(0)
	handlerId := C.g_signal_connect_object(
		C.gpointer(widget),
		signal,
		C.GCallback(C.handleClick),
		C.gpointer(widget),
		flags)

	id := (*C.GtkWidget)(widget)
	gtkSignalToMenuItem[id] = item
	gtkSignalHandlers[id] = handlerId
	impl.handlerId = uint(handlerId)
}

// menuItem
func menuItemChecked(widget pointer) bool {
	if C.gtk_check_menu_item_get_active((*C.GtkCheckMenuItem)(widget)) == C.int(1) {
		return true
	}
	return false
}

func menuItemNew(label string) pointer {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	return pointer(C.gtk_menu_item_new_with_label(cLabel))
}

func menuCheckItemNew(label string) pointer {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	return pointer(C.gtk_check_menu_item_new_with_label(cLabel))
}

func menuItemSetChecked(widget pointer, checked bool) {
	value := C.int(0)
	if checked {
		value = C.int(1)
	}
	C.gtk_check_menu_item_set_active(
		(*C.GtkCheckMenuItem)(widget),
		value)
}

func menuItemSetDisabled(widget pointer, disabled bool) {
	value := C.int(1)
	if disabled {
		value = C.int(0)
	}
	C.gtk_widget_set_sensitive(
		(*C.GtkWidget)(widget),
		value)
}

func menuItemSetLabel(widget pointer, label string) {
	value := C.CString(label)
	C.gtk_menu_item_set_label(
		(*C.GtkMenuItem)(widget),
		value)
	C.free(unsafe.Pointer(value))
}

func menuItemSetToolTip(widget pointer, tooltip string) {
	value := C.CString(tooltip)
	C.gtk_widget_set_tooltip_text(
		(*C.GtkWidget)(widget),
		value)
	C.free(unsafe.Pointer(value))
}

func menuItemSignalBlock(widget pointer, handlerId uint, block bool) {
	if block {
		C.g_signal_handler_block(C.gpointer(widget), C.ulong(handlerId))
	} else {
		C.g_signal_handler_unblock(C.gpointer(widget), C.ulong(handlerId))
	}
}

func menuRadioItemNew(group *GSList, label string) pointer {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	return pointer(C.gtk_radio_menu_item_new_with_label((*C.GSList)(group), cLabel))
}

// screen related

func getScreenByIndex(display *C.struct__GdkDisplay, index int) *Screen {
	monitor := C.gdk_display_get_monitor(display, C.int(index))
	// TODO: Do we need to update Screen to contain current info?
	//	currentMonitor := C.gdk_display_get_monitor_at_window(display, window)

	var geometry C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &geometry)
	primary := false
	if C.gdk_monitor_is_primary(monitor) == 1 {
		primary = true
	}

	return &Screen{
		IsPrimary: primary,
		Scale:     1.0,
		X:         int(geometry.x),
		Y:         int(geometry.y),
		Size: Size{
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
	}
}

func getScreens(app pointer) ([]*Screen, error) {
	var screens []*Screen
	window := C.gtk_application_get_active_window((*C.GtkApplication)(app))
	display := C.gdk_window_get_display((*C.GdkWindow)(unsafe.Pointer(window)))
	count := C.gdk_display_get_n_monitors(display)
	for i := 0; i < int(count); i++ {
		screens = append(screens, getScreenByIndex(display, i))
	}
	return screens, nil
}

// widgets
func widgetSetVisible(widget pointer, hidden bool) {
	if hidden {
		C.gtk_widget_hide((*C.GtkWidget)(widget))
	} else {
		C.gtk_widget_show((*C.GtkWidget)(widget))
	}
}

// window related functions
func windowClose(window pointer) {
	C.gtk_window_close((*C.GtkWindow)(window))
}

func windowEnableDND(id uint, webview pointer) {
	dnd := C.CString("text/uri-list")
	defer C.free(unsafe.Pointer(dnd))
	targetentry := C.gtk_target_entry_new(dnd, 0, C.guint(id))
	defer C.gtk_target_entry_free(targetentry)
	C.gtk_drag_dest_set((*C.GtkWidget)(webview), C.GTK_DEST_DEFAULT_DROP, targetentry, 1, C.GDK_ACTION_COPY)
	event := C.CString("drag-data-received")
	defer C.free(unsafe.Pointer(event))
	windowId := C.uint(id)
	C.signal_connect((*C.GtkWidget)(unsafe.Pointer(webview)), event, C.onDragNDrop, unsafe.Pointer(C.gpointer(&windowId)))
}

func windowExecJS(webview pointer, js string) {
	value := C.CString(js)
	C.webkit_web_view_evaluate_javascript((*C.WebKitWebView)(webview),
		value,
		C.long(len(js)),
		nil,
		C.CString(""),
		nil,
		nil,
		nil)
	C.free(unsafe.Pointer(value))
}

func windowDestroy(window pointer) {
	// Should this truly 'destroy' ?
	C.gtk_window_close((*C.GtkWindow)(window))
	//C.gtk_widget_destroy((*C.GtkWidget)(window))
}

func windowFullscreen(window pointer) {
	C.gtk_window_fullscreen((*C.GtkWindow)(window))
}

func windowGetCurrentMonitor(window pointer) *C.GdkMonitor {
	// Get the monitor that the window is currently on
	display := C.gtk_widget_get_display((*C.GtkWidget)(window))
	gdk_window := C.gtk_widget_get_window((*C.GtkWidget)(window))
	if gdk_window == nil {
		return nil
	}
	return C.gdk_display_get_monitor_at_window(display, gdk_window)
}

func windowGetCurrentMonitorGeometry(window pointer) (x int, y int, width int, height int, scale int) {
	monitor := windowGetCurrentMonitor(window)
	if monitor == nil {
		return -1, -1, -1, -1, 1
	}
	var result C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &result)
	scale = int(C.gdk_monitor_get_scale_factor(monitor))
	return int(result.x), int(result.y), int(result.width), int(result.height), scale
}

func windowGetSize(window pointer) (int, int) {
	var windowWidth C.int
	var windowHeight C.int
	C.gtk_window_get_size((*C.GtkWindow)(window), &windowWidth, &windowHeight)
	return int(windowWidth), int(windowHeight)
}

func windowGetPosition(window pointer) (int, int) {
	var x C.int
	var y C.int
	C.gtk_window_get_position((*C.GtkWindow)(window), &x, &y)
	return int(x), int(y)
}

func windowHide(window pointer) {
	C.gtk_widget_hide((*C.GtkWidget)(window))
}

func windowIsFullscreen(window pointer) bool {
	gdkwindow := C.gtk_widget_get_window((*C.GtkWidget)(window))
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_FULLSCREEN > 0
}

func windowIsMaximized(window pointer) bool {
	gdkwindow := C.gtk_widget_get_window((*C.GtkWidget)(window))
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_MAXIMIZED > 0 && state&C.GDK_WINDOW_STATE_FULLSCREEN == 0
}

func windowIsMinimized(window pointer) bool {
	gdkwindow := C.gtk_widget_get_window((*C.GtkWidget)(window))
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_ICONIFIED > 0
}

func windowIsVisible(window pointer) bool {
	if C.gtk_widget_is_visible((*C.GtkWidget)(window)) == 1 {
		return true
	}
	return false
}

func windowMaximize(window pointer) {
	C.gtk_window_maximize((*C.GtkWindow)(window))
}

func windowMinimize(window pointer) {
	C.gtk_window_iconify((*C.GtkWindow)(window))
}

func windowNew(application pointer, menu pointer, windowId uint, gpuPolicy int) (window pointer, webview pointer) {
	window = pointer(C.gtk_application_window_new((*C.GtkApplication)(application)))
	C.g_object_ref_sink(C.gpointer(window))
	webview = windowNewWebview(windowId, gpuPolicy)
	vbox := pointer(C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0))
	C.gtk_container_add((*C.GtkContainer)(window), (*C.GtkWidget)(vbox))

	if menu != nil {
		C.gtk_box_pack_start((*C.GtkBox)(vbox), (*C.GtkWidget)(menu), 0, 0, 0)
	}
	C.gtk_box_pack_start((*C.GtkBox)(unsafe.Pointer(vbox)), (*C.GtkWidget)(webview), 1, 1, 0)
	return
}

func windowNewWebview(parentId uint, gpuPolicy int) pointer {
	manager := C.webkit_user_content_manager_new()
	external := C.CString("external")
	C.webkit_user_content_manager_register_script_message_handler(manager, external)
	C.free(unsafe.Pointer(external))
	webview := C.webkit_web_view_new_with_user_content_manager(manager)
	id := C.uint(parentId)
	if !registered {
		wails := C.CString("wails")
		C.webkit_web_context_register_uri_scheme(
			C.webkit_web_context_get_default(),
			wails,
			C.WebKitURISchemeRequestCallback(C.onProcessRequest),
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
	return pointer(webview)
}

func windowPresent(window pointer) {
	C.gtk_window_present((*C.GtkWindow)(window))
	// gtk_window_unminimize ((*C.GtkWindow)(w.window)) /// gtk4
}

func windowReload(webview pointer, address string) {
	uri := C.CString(address)
	C.webkit_web_view_load_uri((*C.WebKitWebView)(webview), uri)
	C.free(unsafe.Pointer(uri))
}

func windowResize(window pointer, width, height int) {
	C.gtk_window_resize(
		(*C.GtkWindow)(window),
		C.gint(width),
		C.gint(height))
}

func windowShow(window pointer) {
	C.gtk_widget_show_all((*C.GtkWidget)(window))
}

func windowSetBackgroundColour(webview pointer, colour RGBA) {
	rgba := C.GdkRGBA{C.double(colour.Red) / 255.0, C.double(colour.Green) / 255.0, C.double(colour.Blue) / 255.0, C.double(colour.Alpha) / 255.0}
	C.webkit_web_view_set_background_color((*C.WebKitWebView)(webview), &rgba)
}

func windowSetGeometryHints(window pointer, minWidth, minHeight, maxWidth, maxHeight int) {
	size := C.GdkGeometry{
		min_width:  C.int(minWidth),
		min_height: C.int(minHeight),
		max_width:  C.int(maxWidth),
		max_height: C.int(maxHeight),
	}
	C.gtk_window_set_geometry_hints((*C.GtkWindow)(window), nil, &size, C.GDK_HINT_MAX_SIZE|C.GDK_HINT_MIN_SIZE)
}

func windowSetFrameless(window pointer, frameless bool) {
	C.gtk_window_set_decorated((*C.GtkWindow)(window), gtkBool(!frameless))
	// TODO: Deal with transparency for the titlebar if possible when !frameless
	//       Perhaps we just make it undecorated and add a menu bar inside?
}

// TODO: confirm this is working properly
func windowSetHTML(webview pointer, html string) {
	cHTML := C.CString(html)
	uri := C.CString("wails://")
	empty := C.CString("")
	defer C.free(unsafe.Pointer(cHTML))
	defer C.free(unsafe.Pointer(uri))
	defer C.free(unsafe.Pointer(empty))
	C.webkit_web_view_load_alternate_html(
		(*C.WebKitWebView)(webview),
		cHTML,
		uri,
		empty)
}

func windowSetKeepAbove(window pointer, alwaysOnTop bool) {
	C.gtk_window_set_keep_above((*C.GtkWindow)(window), gtkBool(alwaysOnTop))
}

func windowSetResizable(window pointer, resizable bool) {
	C.gtk_window_set_resizable((*C.GtkWindow)(window), gtkBool(resizable))
}

func windowSetTitle(window pointer, title string) {
	cTitle := C.CString(title)
	C.gtk_window_set_title((*C.GtkWindow)(window), cTitle)
	C.free(unsafe.Pointer(cTitle))
}

func windowSetTransparent(window pointer) {
	screen := C.gtk_widget_get_screen((*C.GtkWidget)(window))
	visual := C.gdk_screen_get_rgba_visual(screen)

	if visual != nil && C.gdk_screen_is_composited(screen) == C.int(1) {
		C.gtk_widget_set_app_paintable((*C.GtkWidget)(window), C.gboolean(1))
		C.gtk_widget_set_visual((*C.GtkWidget)(window), visual)
	}
}

func windowSetURL(webview pointer, uri string) {
	target := C.CString(uri)
	C.webkit_web_view_load_uri((*C.WebKitWebView)(webview), target)
	C.free(unsafe.Pointer(target))
}

//export emit
func emit(we *C.WindowEvent) {
	window := globalApplication.getWindowForID(uint(we.id))
	if window != nil {
		window.emit(events.WindowEventType(we.event))
	}
}

func windowSetupSignalHandlers(windowId uint, window, webview pointer, emit func(e events.WindowEventType)) {
	event := C.CString("delete-event")
	defer C.free(unsafe.Pointer(event))
	wEvent := C.WindowEvent{
		id:    C.uint(windowId),
		event: C.uint(events.Common.WindowClosing),
	}
	C.signal_connect((*C.GtkWidget)(window), event, C.emit, unsafe.Pointer(&wEvent))

	/*
		event = C.CString("load-changed")
		defer C.free(unsafe.Pointer(event))
		C.signal_connect(webview, event, C.webviewLoadChanged, unsafe.Pointer(&w.parent.id))
	*/
	id := C.uint(windowId)
	event = C.CString("button-press-event")
	C.signal_connect((*C.GtkWidget)(unsafe.Pointer(webview)), event, C.onButtonEvent, unsafe.Pointer(&id))
	C.free(unsafe.Pointer(event))
	event = C.CString("button-release-event")
	defer C.free(unsafe.Pointer(event))
	C.signal_connect((*C.GtkWidget)(unsafe.Pointer(webview)), event, C.onButtonEvent, unsafe.Pointer(&id))
}

func windowToggleDevTools(webview pointer) {
	settings := C.webkit_web_view_get_settings((*C.WebKitWebView)(webview))
	enabled := C.webkit_settings_get_enable_developer_extras(settings)
	switch enabled {
	case C.int(0):
		enabled = C.int(1)
	case C.int(1):
		enabled = C.int(0)
	}
	C.webkit_settings_set_enable_developer_extras(settings, enabled)
}

func windowUnfullscreen(window pointer) {
	C.gtk_window_unfullscreen((*C.GtkWindow)(window))
}

func windowUnmaximize(window pointer) {
	C.gtk_window_unmaximize((*C.GtkWindow)(window))
}

func windowZoom(webview pointer) float64 {
	return float64(C.webkit_web_view_get_zoom_level((*C.WebKitWebView)(webview)))
}

// FIXME: ZoomIn/Out is assumed to be incorrect!
func windowZoomIn(webview pointer) {
	ZoomInFactor := 1.10
	windowZoomSet(webview, windowZoom(webview)*ZoomInFactor)
}
func windowZoomOut(webview pointer) {
	ZoomOutFactor := -1.10
	windowZoomSet(webview, windowZoom(webview)*ZoomOutFactor)
}

func windowZoomSet(webview pointer, zoom float64) {
	C.webkit_web_view_set_zoom_level((*C.WebKitWebView)(webview), C.double(zoom))
}

func windowMove(window pointer, x, y int) {
	C.gtk_window_move((*C.GtkWindow)(window), C.int(x), C.int(y))
}

//export onButtonEvent
func onButtonEvent(_ *C.GtkWidget, event *C.GdkEventButton, data unsafe.Pointer) C.gboolean {
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
		lw.startDrag() //uint(event.button), int(event.x_root), int(event.y_root))
	case Gdk2ButtonPress:
		fmt.Printf("%d - button %d - double-clicked\n", windowId, int(event.button))
	case GdkButtonRelease:
		lw.endDrag(uint(event.button), int(event.x_root), int(event.y_root))
	}

	return C.gboolean(0)
}

//export onDragNDrop
func onDragNDrop(target unsafe.Pointer, context *C.GdkDragContext, x C.gint, y C.gint, seldata unsafe.Pointer, info C.guint, time C.guint, data unsafe.Pointer) {
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

//export onProcessRequest
func onProcessRequest(request unsafe.Pointer, data unsafe.Pointer) {
	windowId := uint(*((*C.uint)(data)))
	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(request),
		windowId:   windowId,
		windowName: globalApplication.getWindowForID(windowId).Name(),
	}
}

func gtkBool(input bool) C.gboolean {
	if input {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

// dialog related

func setWindowIcon(window pointer, icon []byte) {
	fmt.Println("setWindowIcon", len(icon))
	loader := C.gdk_pixbuf_loader_new()
	if loader == nil {
		return
	}
	written := C.gdk_pixbuf_loader_write(
		loader,
		(*C.uchar)(&icon[0]),
		C.ulong(len(icon)),
		nil)
	if written == 0 {
		fmt.Println("failed to write icon")
		return
	}
	C.gdk_pixbuf_loader_close(loader, nil)
	pixbuf := C.gdk_pixbuf_loader_get_pixbuf(loader)
	if pixbuf != nil {
		fmt.Println("gtk_window_set_icon", window)
		C.gtk_window_set_icon((*C.GtkWindow)(window), pixbuf)
	}
	C.g_object_unref(C.gpointer(loader))
}

//export messageDialogCB
func messageDialogCB(button C.int) {
	fmt.Println("messageDialogCB", button)

}

func runQuestionDialog(parent pointer, options *MessageDialog) int {
	cMsg := C.CString(options.Message)
	cTitle := C.CString(options.Title)
	defer C.free(unsafe.Pointer(cMsg))
	defer C.free(unsafe.Pointer(cTitle))
	hasButtons := false
	if len(options.Buttons) > 0 {
		hasButtons = true
	}

	dType, ok := map[DialogType]C.int{
		InfoDialog:     C.GTK_MESSAGE_INFO,
		QuestionDialog: C.GTK_MESSAGE_QUESTION,
		WarningDialog:  C.GTK_MESSAGE_WARNING,
	}[options.DialogType]
	if !ok {
		// FIXME: Add logging here!
		dType = C.GTK_MESSAGE_INFO
	}

	dialog := C.new_message_dialog((*C.GtkWindow)(parent), cMsg, dType, C.bool(hasButtons))
	if options.Title != "" {
		C.gtk_window_set_title(
			(*C.GtkWindow)(unsafe.Pointer(dialog)),
			cTitle)
	}

	if img, err := pngToImage(options.Icon); err == nil {
		gbytes := C.g_bytes_new_static(
			C.gconstpointer(unsafe.Pointer(&img.Pix[0])),
			C.ulong(len(img.Pix)))
		defer C.g_bytes_unref(gbytes)
		pixBuf := C.gdk_pixbuf_new_from_bytes(
			gbytes,
			C.GDK_COLORSPACE_RGB,
			1, // has_alpha
			8,
			C.int(img.Bounds().Dx()),
			C.int(img.Bounds().Dy()),
			C.int(img.Stride),
		)
		image := C.gtk_image_new_from_pixbuf(pixBuf)
		C.gtk_widget_set_visible((*C.GtkWidget)(image), C.gboolean(1))
		contentArea := C.gtk_dialog_get_content_area((*C.GtkDialog)(dialog))
		C.gtk_container_add(
			(*C.GtkContainer)(unsafe.Pointer(contentArea)),
			(*C.GtkWidget)(image))
	}
	for i, button := range options.Buttons {
		cLabel := C.CString(button.Label)
		defer C.free(unsafe.Pointer(cLabel))
		index := C.int(i)
		C.gtk_dialog_add_button(
			(*C.GtkDialog)(dialog), cLabel, index)
		if button.IsDefault {
			C.gtk_dialog_set_default_response((*C.GtkDialog)(dialog), index)
		}
	}

	defer C.gtk_widget_destroy((*C.GtkWidget)(dialog))
	return int(C.gtk_dialog_run((*C.GtkDialog)(unsafe.Pointer(dialog))))
}

//export openFileDialogCallbackEnd
func openFileDialogCallbackEnd(cid C.uint) {
	id := uint(cid)
	channel, ok := openFileResponses[id]
	if ok {
		close(channel)
		delete(openFileResponses, id)
		freeDialogID(id)
	} else {
		panic("No channel found for open file dialog")
	}
}

//export openFileDialogCallback
func openFileDialogCallback(cid C.uint, cpath *C.char) {
	path := C.GoString(cpath)
	id := uint(cid)
	channel, ok := openFileResponses[id]
	if ok {
		channel <- path
	} else {
		panic("No channel found for open file dialog")
	}
}

//export saveFileDialogCallback
func saveFileDialogCallback(cid C.uint, cpath *C.char) {
	// Covert the path to a string
	path := C.GoString(cpath)
	id := uint(cid)
	// put response on channel
	channel, ok := saveFileResponses[id]
	if ok {
		channel <- path
		close(channel)
		delete(saveFileResponses, id)
		freeDialogID(id)

	} else {
		panic("No channel found for save file dialog")
	}
}
