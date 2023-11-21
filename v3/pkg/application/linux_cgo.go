//go:build linux && cgo

package application

import (
	"fmt"
	"regexp"
	"strings"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0 javascriptcoregtk-4.1

#include <JavaScriptCore/JavaScript.h>
#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>
#ifdef G_APPLICATION_DEFAULT_FLAGS
    #define APPLICATION_DEFAULT_FLAGS G_APPLICATION_DEFAULT_FLAGS
#else
    #define APPLICATION_DEFAULT_FLAGS G_APPLICATION_FLAGS_NONE
#endif

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
extern gboolean windowDeleteHandler(GtkWidget *widget, GdkEvent *event, gpointer user_data);
extern gboolean onKeyPressEvent (GtkWidget *widget, GdkEventKey *event, gpointer user_data);
extern void onProcessRequest(void *request, gpointer user_data);
extern void sendMessageToBackend(WebKitUserContentManager *contentManager, WebKitJavascriptResult *result, void *data);
void webviewLoadChanged(WebKitWebView *web_view, WebKitLoadEvent load_event, gpointer data);
// exported below (end)

static void signal_connect(void *widget, char *event, void *cb, void* data) {
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
       "%s",
       msg);

   // g_signal_connect_swapped (dialog,
   //                           "response",
   //                           G_CALLBACK (callback),
   //                           dialog);
   return dialog;
};

extern void messageDialogCB(gint button);

static void* gtkFileChooserDialogNew(char* title, GtkWindow* window, GtkFileChooserAction action, char* cancelLabel, char* acceptLabel) {
   // gtk_file_chooser_dialog_new is variadic!  Can't call from cgo directly
	return (GtkFileChooser*)gtk_file_chooser_dialog_new(
		title,
		window,
		action,
		cancelLabel,
		GTK_RESPONSE_CANCEL,
		acceptLabel,
		GTK_RESPONSE_ACCEPT,
		NULL);
}

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
	nilPointer    pointer       = nil
	nilRadioGroup GSListPointer = nil
)

var (
	gtkSignalToMenuItem map[uint]*MenuItem
	mainThreadId        *C.GThread
)

func init() {
	gtkSignalToMenuItem = map[uint]*MenuItem{}

	mainThreadId = C.g_thread_self()
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
	applicationEvents <- newApplicationEvent(events.Linux.ApplicationStarted)
}

func isOnMainThread() bool {
	threadId := C.g_thread_self()
	return threadId == mainThreadId
}

// implementation below
func appName() string {
	name := C.g_get_application_name()
	defer C.free(unsafe.Pointer(name))
	return C.GoString(name)
}

func appNew(name string) pointer {
	// prevent leading number
	if matched, _ := regexp.MatchString(`^\d+`, name); matched {
		name = fmt.Sprintf("_%s", name)
	}
	name = strings.Replace(name, "(", "_", -1)
	name = strings.Replace(name, ")", "_", -1)
	appId := fmt.Sprintf("org.wails.%s", name)
	nameC := C.CString(appId)
	defer C.free(unsafe.Pointer(nameC))
	return pointer(C.gtk_application_new(nameC, C.APPLICATION_DEFAULT_FLAGS))
}

func appRun(app pointer) error {
	application := (*C.GApplication)(app)
	C.g_application_hold(application) // allows it to run without a window

	signal := C.CString("activate")
	defer C.free(unsafe.Pointer(signal))
	C.g_signal_connect_data(
		C.gpointer(application),
		signal,
		C.GCallback(C.activateLinux),
		nil,
		nil,
		0)
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

func contextMenuShow(window pointer, menu pointer, data *ContextMenuData) {
	geometry := C.GdkRectangle{
		x: C.int(data.X),
		y: C.int(data.Y),
	}
	event := C.GdkEvent{}
	gdkWindow := C.gtk_widget_get_window((*C.GtkWidget)(window))
	C.gtk_menu_popup_at_rect(
		(*C.GtkMenu)(menu),
		gdkWindow,
		(*C.GdkRectangle)(&geometry),
		C.GDK_GRAVITY_NORTH_WEST,
		C.GDK_GRAVITY_NORTH_WEST,
		(*C.GdkEvent)(&event),
	)
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
		C.gtk_window_present((*C.GtkWindow)(window))
	}
}

// Clipboard
func clipboardGet() string {
	clip := C.gtk_clipboard_get(C.GDK_SELECTION_CLIPBOARD)
	text := C.gtk_clipboard_wait_for_text(clip)
	return C.GoString(text)
}

func clipboardSet(text string) {
	cText := C.CString(text)
	clip := C.gtk_clipboard_get(C.GDK_SELECTION_CLIPBOARD)
	C.gtk_clipboard_set_text(clip, cText, -1)

	clip = C.gtk_clipboard_get(C.GDK_SELECTION_PRIMARY)
	C.gtk_clipboard_set_text(clip, cText, -1)
	C.free(unsafe.Pointer(cText))
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
	ident := C.CString("id")
	defer C.free(unsafe.Pointer(ident))
	value := C.g_object_get_data((*C.GObject)(idPtr), ident)
	id := uint(*(*C.uint)(value))
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

func attachMenuHandler(item *MenuItem) uint {
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

	id := C.uint(item.id)
	ident := C.CString("id")
	defer C.free(unsafe.Pointer(ident))
	C.g_object_set_data(
		(*C.GObject)(widget),
		ident,
		C.gpointer(&id),
	)

	gtkSignalToMenuItem[item.id] = item
	return uint(handlerId)
}

// menuItem
func menuItemChecked(widget pointer) bool {
	if C.gtk_check_menu_item_get_active((*C.GtkCheckMenuItem)(widget)) == C.int(1) {
		return true
	}
	return false
}

func menuItemNew(label string, bitmap []byte) pointer {
	return menuItemAddProperties(C.gtk_menu_item_new(), label, bitmap)
}

func menuItemAddProperties(menuItem *C.GtkWidget, label string, bitmap []byte) pointer {
	/*
		   // FIXME: Support accelerator configuration
		   activate := C.CString("activate")
			defer C.free(unsafe.Pointer(activate))
			accelGroup := C.gtk_accel_group_new()
			C.gtk_widget_add_accelerator(menuItem, activate, accelGroup,
				C.GDK_KEY_m, C.GDK_CONTROL_MASK, C.GTK_ACCEL_VISIBLE)
	*/
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	lbl := unsafe.Pointer(C.gtk_accel_label_new(cLabel))
	C.gtk_label_set_use_underline((*C.GtkLabel)(lbl), 1)
	C.gtk_label_set_xalign((*C.GtkLabel)(lbl), 0.0)
	C.gtk_accel_label_set_accel_widget(
		(*C.GtkAccelLabel)(lbl),
		(*C.GtkWidget)(unsafe.Pointer(menuItem)))

	box := C.gtk_box_new(C.GTK_ORIENTATION_HORIZONTAL, 6)
	if img, err := pngToImage(bitmap); err == nil {
		gbytes := C.g_bytes_new_static(C.gconstpointer(unsafe.Pointer(&img.Pix[0])),
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
		C.gtk_container_add(
			(*C.GtkContainer)(unsafe.Pointer(box)),
			(*C.GtkWidget)(unsafe.Pointer(image)))
	}

	C.gtk_box_pack_end(
		(*C.GtkBox)(unsafe.Pointer(box)),
		(*C.GtkWidget)(lbl), 1, 1, 0)
	C.gtk_container_add(
		(*C.GtkContainer)(unsafe.Pointer(menuItem)),
		(*C.GtkWidget)(unsafe.Pointer(box)))
	C.gtk_widget_show_all(menuItem)
	return pointer(menuItem)
}

func menuCheckItemNew(label string, bitmap []byte) pointer {
	return menuItemAddProperties(C.gtk_check_menu_item_new(), label, bitmap)
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

func menuItemRemoveBitmap(widget pointer) {
	box := C.gtk_bin_get_child((*C.GtkBin)(widget))
	if box == nil {
		return
	}

	children := C.gtk_container_get_children((*C.GtkContainer)(unsafe.Pointer(box)))
	defer C.g_list_free(children)
	count := int(C.g_list_length(children))
	if count == 2 {
		C.gtk_container_remove((*C.GtkContainer)(unsafe.Pointer(box)),
			(*C.GtkWidget)(children.data))
	}
}

func menuItemSetBitmap(widget pointer, bitmap []byte) {
	menuItemRemoveBitmap(widget)
	box := C.gtk_bin_get_child((*C.GtkBin)(widget))
	if img, err := pngToImage(bitmap); err == nil {
		gbytes := C.g_bytes_new_static(C.gconstpointer(unsafe.Pointer(&img.Pix[0])),
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
		C.gtk_container_add(
			(*C.GtkContainer)(unsafe.Pointer(box)),
			(*C.GtkWidget)(unsafe.Pointer(image)))
	}

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
		ID:        fmt.Sprintf("%d", index),
		Name:      fmt.Sprintf("Screen %d", index),
		IsPrimary: primary,
		Scale:     float32(C.gdk_monitor_get_scale_factor(monitor)),
		X:         int(geometry.x),
		Y:         int(geometry.y),
		Size: Size{
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
		Bounds: Rect{
			X:      int(geometry.x),
			Y:      int(geometry.y),
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
	}
}

func getScreens(app pointer) ([]*Screen, error) {
	var screens []*Screen
	window := C.gtk_application_get_active_window((*C.GtkApplication)(app))
	gdkWindow := C.gtk_widget_get_window((*C.GtkWidget)(unsafe.Pointer(window)))
	display := C.gdk_window_get_display(gdkWindow)
	count := C.gdk_display_get_n_monitors(display)
	for i := 0; i < int(count); i++ {
		screens = append(screens, getScreenByIndex(display, i))
	}
	return screens, nil
}

// widgets
func widgetSetSensitive(widget pointer, enabled bool) {
	value := C.int(0)
	if enabled {
		value = C.int(1)
	}

	C.gtk_widget_set_sensitive((*C.GtkWidget)(widget), value)
}

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
	C.signal_connect(unsafe.Pointer(webview), event, C.onDragNDrop, unsafe.Pointer(C.gpointer(&windowId)))
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

func windowGetAbsolutePosition(window pointer) (int, int) {
	var x C.int
	var y C.int
	C.gtk_window_get_position((*C.GtkWindow)(window), &x, &y)
	return int(x), int(y)
}

func windowGetSize(window pointer) (int, int) {
	var windowWidth C.int
	var windowHeight C.int
	C.gtk_window_get_size((*C.GtkWindow)(window), &windowWidth, &windowHeight)
	return int(windowWidth), int(windowHeight)
}

func windowGetRelativePosition(window pointer) (int, int) {
	x, y := windowGetAbsolutePosition(window)
	// The position must be relative to the screen it is on
	// We need to get the screen it is on
	monitor := windowGetCurrentMonitor(window)
	geometry := C.GdkRectangle{}
	C.gdk_monitor_get_geometry(monitor, &geometry)
	x = x - int(geometry.x)
	y = y - int(geometry.y)

	// TODO: Scale based on DPI

	return x, y
}

func windowHide(window pointer) {
	C.gtk_widget_hide((*C.GtkWidget)(window))
}

func windowIsFullscreen(window pointer) bool {
	gdkwindow := C.gtk_widget_get_window((*C.GtkWidget)(window))
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_FULLSCREEN > 0
}

func windowIsFocused(window pointer) bool {
	// returns true if window is focused
	return C.gtk_window_has_toplevel_focus((*C.GtkWindow)(window)) == 1
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

func windowNew(application pointer, menu pointer, windowId uint, gpuPolicy int) (window, webview, vbox pointer) {
	window = pointer(C.gtk_application_window_new((*C.GtkApplication)(application)))
	C.g_object_ref_sink(C.gpointer(window))
	webview = windowNewWebview(windowId, gpuPolicy)
	vbox = pointer(C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0))
	name := C.CString("webview-box")
	defer C.free(unsafe.Pointer(name))
	C.gtk_widget_set_name((*C.GtkWidget)(vbox), name)

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

func windowSetBackgroundColour(vbox, webview pointer, colour RGBA) {
	rgba := C.GdkRGBA{C.double(colour.Red) / 255.0, C.double(colour.Green) / 255.0, C.double(colour.Blue) / 255.0, C.double(colour.Alpha) / 255.0}
	C.webkit_web_view_set_background_color((*C.WebKitWebView)(webview), &rgba)

	colour.Alpha = 255
	cssStr := C.CString(fmt.Sprintf("#webview-box {background-color: rgba(%d, %d, %d, %1.1f);}", colour.Red, colour.Green, colour.Blue, float32(colour.Alpha)/255.0))
	provider := C.gtk_css_provider_new()
	C.gtk_style_context_add_provider(
		C.gtk_widget_get_style_context((*C.GtkWidget)(vbox)),
		(*C.GtkStyleProvider)(unsafe.Pointer(provider)),
		C.GTK_STYLE_PROVIDER_PRIORITY_USER)
	C.g_object_unref(C.gpointer(provider))
	C.gtk_css_provider_load_from_data(provider, cssStr, -1, nil)
	C.free(unsafe.Pointer(cssStr))
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

//export windowDeleteHandler
func windowDeleteHandler(window *C.GtkWidget, event *C.GdkEvent, data unsafe.Pointer) C.gboolean {
	id := uint(*((*C.uint)(data)))
	windowEvents <- &windowEvent{
		WindowID: id,
		EventID:  uint(events.WindowEventType(events.Common.WindowClosing)),
	}
	// stop propagation
	return C.gboolean(1)
}

//export webviewLoadChanged
func webviewLoadChanged(web_view *C.WebKitWebView, load_event C.WebKitLoadEvent, data unsafe.Pointer) {
	if load_event == C.WEBKIT_LOAD_FINISHED {
		//applicationEvents <- newApplicationEvent(events.Linux.ApplicationStarted)
	}
}

func windowSetupSignalHandlers(windowId uint, window, webview pointer, emit func(e events.WindowEventType)) {
	event := C.CString("delete-event")
	defer C.free(unsafe.Pointer(event))
	C.signal_connect(unsafe.Pointer(window), event, C.windowDeleteHandler, unsafe.Pointer(&windowId))

	contentManager := C.webkit_web_view_get_user_content_manager((*C.WebKitWebView)(webview))
	event = C.CString("script-message-received::external")
	defer C.free(unsafe.Pointer(event))
	C.signal_connect(unsafe.Pointer(contentManager), event, C.sendMessageToBackend, nil)

	event = C.CString("load-changed")
	defer C.free(unsafe.Pointer(event))
	C.signal_connect(unsafe.Pointer(webview), event, C.webviewLoadChanged, unsafe.Pointer(&windowId))

	id := C.uint(windowId)
	event = C.CString("button-press-event")
	C.signal_connect(unsafe.Pointer(webview), event, C.onButtonEvent, unsafe.Pointer(&id))
	C.free(unsafe.Pointer(event))
	event = C.CString("button-release-event")
	defer C.free(unsafe.Pointer(event))
	C.signal_connect(unsafe.Pointer(webview), event, C.onButtonEvent, unsafe.Pointer(&id))

	event = C.CString("key-press-event")
	defer C.free(unsafe.Pointer(event))
	C.signal_connect(unsafe.Pointer(webview), event, C.onKeyPressEvent, unsafe.Pointer(&id))
}

func windowShowDevTools(webview pointer) {
	inspector := C.webkit_web_view_get_inspector((*C.WebKitWebView)(webview))
	C.webkit_web_inspector_show(inspector)
}

func windowStartDrag(window pointer, button uint, xroot int, yroot int, dragTime uint32) {
	C.gtk_window_begin_move_drag(
		(*C.GtkWindow)(window),
		C.int(button),
		C.int(xroot),
		C.int(yroot),
		C.uint32_t(dragTime))
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
	if zoom < 1 { // 1.0 is the smallest allowable
		zoom = 1
	}
	C.webkit_web_view_set_zoom_level((*C.WebKitWebView)(webview), C.double(zoom))
}

func windowMove(window pointer, x, y int) {
	C.gtk_window_move((*C.GtkWindow)(window), C.int(x), C.int(y))
}

// FIXME Change this to reflect mouse button!
//
//export onButtonEvent
func onButtonEvent(_ *C.GtkWidget, event *C.GdkEventButton, data unsafe.Pointer) C.gboolean {
	// Constants (defined here to be easier to use with purego)
	GdkButtonPress := C.GDK_BUTTON_PRESS     // 4
	Gdk2ButtonPress := C.GDK_2BUTTON_PRESS   // 5 for double-click
	GdkButtonRelease := C.GDK_BUTTON_RELEASE // 7

	windowId := uint(*((*C.uint)(data)))
	window := globalApplication.getWindowForID(windowId)
	if window == nil {
		return C.gboolean(0)
	}
	lw, ok := (window.(*WebviewWindow).impl).(*linuxWebviewWindow)
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
		lw.drag.MouseButton = uint(event.button)
		lw.drag.XRoot = int(event.x_root)
		lw.drag.YRoot = int(event.y_root)
		lw.drag.DragTime = uint32(event.time)
	case Gdk2ButtonPress:
		// do we need something here?
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

//export onKeyPressEvent
func onKeyPressEvent(widget *C.GtkWidget, event *C.GdkEventKey, userData unsafe.Pointer) C.gboolean {
	windowID := uint(*((*C.uint)(userData)))
	accelerator, ok := getKeyboardState(event)
	if !ok {
		return C.gboolean(1)
	}
	windowKeyEvents <- &windowKeyEvent{
		windowId:          windowID,
		acceleratorString: accelerator,
	}
	return C.gboolean(1)
}

func getKeyboardState(event *C.GdkEventKey) (string, bool) {
	modifiers := uint(event.state) & C.GDK_MODIFIER_MASK
	keyCode := uint(event.keyval)

	var acc accelerator
	// Check Accelerators
	if modifiers&(C.GDK_SHIFT_MASK) != 0 {
		acc.Modifiers = append(acc.Modifiers, ShiftKey)
	}
	if modifiers&(C.GDK_CONTROL_MASK) != 0 {
		acc.Modifiers = append(acc.Modifiers, ControlKey)
	}
	if modifiers&(C.GDK_MOD1_MASK) != 0 {
		acc.Modifiers = append(acc.Modifiers, OptionOrAltKey)
	}
	if modifiers&(C.GDK_SUPER_MASK) != 0 {
		acc.Modifiers = append(acc.Modifiers, SuperKey)
	}
	keyString, ok := VirtualKeyCodes[keyCode]
	if !ok {
		fmt.Println("Error Could not find key code: ", keyCode)
		return "", false
	}
	acc.Key = keyString
	return acc.String(), true
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

//export sendMessageToBackend
func sendMessageToBackend(contentManager *C.WebKitUserContentManager, result *C.WebKitJavascriptResult,
	data unsafe.Pointer) {

	var msg string
	value := C.webkit_javascript_result_get_js_value(result)
	message := C.jsc_value_to_string(value)
	msg = C.GoString(message)
	defer C.g_free(C.gpointer(message))
	windowMessageBuffer <- &windowMessage{
		windowId: uint(windowID),
		message:  msg,
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
		return
	}
	C.gdk_pixbuf_loader_close(loader, nil)
	pixbuf := C.gdk_pixbuf_loader_get_pixbuf(loader)
	if pixbuf != nil {
		C.gtk_window_set_icon((*C.GtkWindow)(window), pixbuf)
	}
	C.g_object_unref(C.gpointer(loader))
}

//export messageDialogCB
func messageDialogCB(button C.int) {
	fmt.Println("messageDialogCB", button)
}

func runChooserDialog(window pointer, allowMultiple, createFolders, showHidden bool, currentFolder, title string, action int, acceptLabel string, filters []FileFilter) (chan string, error) {
	titleStr := C.CString(title)
	defer C.free(unsafe.Pointer(titleStr))
	cancelStr := C.CString("_Cancel")
	defer C.free(unsafe.Pointer(cancelStr))
	acceptLabelStr := C.CString(acceptLabel)
	defer C.free(unsafe.Pointer(acceptLabelStr))

	fc := C.gtkFileChooserDialogNew(
		titleStr,
		(*C.GtkWindow)(window),
		C.GtkFileChooserAction(action),
		cancelStr,
		acceptLabelStr)

	C.gtk_file_chooser_set_action((*C.GtkFileChooser)(fc), C.GtkFileChooserAction(action))

	gtkFilters := []*C.GtkFileFilter{}
	for _, filter := range filters {
		f := (C.gtk_file_filter_new())
		displayStr := C.CString(filter.DisplayName)
		C.gtk_file_filter_set_name(f, displayStr)
		C.free(unsafe.Pointer(displayStr))
		patternStr := C.CString(filter.Pattern)
		C.gtk_file_filter_add_pattern(f, patternStr)
		C.free(unsafe.Pointer(patternStr))
		C.gtk_file_chooser_add_filter((*C.GtkFileChooser)(fc), f)
		gtkFilters = append(gtkFilters, f)
	}
	C.gtk_file_chooser_set_select_multiple(
		(*C.GtkFileChooser)(fc),
		gtkBool(allowMultiple))
	C.gtk_file_chooser_set_create_folders(
		(*C.GtkFileChooser)(fc),
		gtkBool(createFolders))
	C.gtk_file_chooser_set_show_hidden(
		(*C.GtkFileChooser)(fc),
		gtkBool(showHidden))

	if currentFolder != "" {
		path := C.CString(currentFolder)
		C.gtk_file_chooser_set_current_folder(
			(*C.GtkFileChooser)(fc),
			path)
		C.free(unsafe.Pointer(path))
	}

	// FIXME: This should be consolidated - duplicate exists in linux_purego.go
	buildStringAndFree := func(s C.gpointer) string {
		bytes := []byte{}
		p := unsafe.Pointer(s)
		for {
			val := *(*byte)(p)
			if val == 0 { // this is the null terminator
				break
			}
			bytes = append(bytes, val)
			p = unsafe.Add(p, 1)
		}
		C.g_free(s) // so we don't have to iterate a second time
		return string(bytes)
	}

	selections := make(chan string)
	// run this on the gtk thread
	InvokeAsync(func() {
		go func() {
			response := C.gtk_dialog_run((*C.GtkDialog)(fc))
			if response == C.GTK_RESPONSE_ACCEPT {
				filenames := C.gtk_file_chooser_get_filenames((*C.GtkFileChooser)(fc))
				iter := filenames
				count := 0
				for {
					selections <- buildStringAndFree(C.gpointer(iter.data))
					iter = iter.next
					if iter == nil || count == 1024 {
						break
					}
					count++
				}
				close(selections)
				C.gtk_widget_destroy((*C.GtkWidget)(unsafe.Pointer(fc)))
			}
		}()
	})
	return selections, nil
}

func runOpenFileDialog(dialog *OpenFileDialogStruct) (chan string, error) {
	const GtkFileChooserActionOpen = C.GTK_FILE_CHOOSER_ACTION_OPEN

	window := nilPointer
	if dialog.window != nil {
		window = (dialog.window.impl).(*linuxWebviewWindow).window
	}

	buttonText := dialog.buttonText
	if buttonText == "" {
		buttonText = "_Open"
	}

	return runChooserDialog(
		window,
		dialog.allowsMultipleSelection,
		dialog.canCreateDirectories,
		dialog.showHiddenFiles,
		dialog.directory,
		dialog.title,
		GtkFileChooserActionOpen,
		buttonText,
		dialog.filters)
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
		InfoDialogType: C.GTK_MESSAGE_INFO,
		//		ErrorDialogType:
		QuestionDialogType: C.GTK_MESSAGE_QUESTION,
		WarningDialogType:  C.GTK_MESSAGE_WARNING,
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

func runSaveFileDialog(dialog *SaveFileDialogStruct) (chan string, error) {
	window := nilPointer
	buttonText := dialog.buttonText
	if buttonText == "" {
		buttonText = "_Save"
	}
	results, err := runChooserDialog(
		window,
		false, // multiple selection
		dialog.canCreateDirectories,
		dialog.showHiddenFiles,
		dialog.directory,
		dialog.title,
		C.GTK_FILE_CHOOSER_ACTION_SAVE,
		buttonText,
		dialog.filters)

	return results, err
}
