//go:build linux && cgo

package application

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"

	"github.com/wailsapp/wails/v3/pkg/events"
)

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1  gdk-3.0

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

static void save_window_id(void *object, uint value)
{
    g_object_set_data((GObject *)object, "windowid", GUINT_TO_POINTER((guint)value));
}

static guint get_window_id(void *object)
{
    return GPOINTER_TO_UINT(g_object_get_data((GObject *)object, "windowid"));
}

// exported below
void activateLinux(gpointer data);
extern void emit(WindowEvent* data);
extern gboolean handleConfigureEvent(GtkWidget*, GdkEventConfigure*, uintptr_t);
extern gboolean handleDeleteEvent(GtkWidget*, GdkEvent*, uintptr_t);
extern gboolean handleFocusEvent(GtkWidget*, GdkEvent*, uintptr_t);
extern void handleLoadChanged(WebKitWebView*, WebKitLoadEvent, uintptr_t);
void handleClick(void*);
extern gboolean onButtonEvent(GtkWidget *widget, GdkEventButton *event, uintptr_t user_data);
extern gboolean onMenuButtonEvent(GtkWidget *widget, GdkEventButton *event, uintptr_t user_data);
extern void onUriList(char **extracted, gpointer data);
extern gboolean onKeyPressEvent (GtkWidget *widget, GdkEventKey *event, uintptr_t user_data);
extern void onProcessRequest(WebKitURISchemeRequest *request, uintptr_t user_data);
extern void sendMessageToBackend(WebKitUserContentManager *contentManager, WebKitJavascriptResult *result, void *data);
// exported below (end)

static void signal_connect(void *widget, char *event, void *cb, void* data) {
   // g_signal_connect is a macro and can't be called directly
   g_signal_connect(widget, event, cb, data);
}

static WebKitWebView* webkit_web_view(GtkWidget *webview) {
	return WEBKIT_WEB_VIEW(webview);
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
	float scaleFactor;
	double rotation;
	bool isPrimary;
} Screen;

// CREDIT: https://github.com/rainycape/magick
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

static void fix_signal(int signum) {
    struct sigaction st;

    if (sigaction(signum, NULL, &st) < 0) {
        goto fix_signal_error;
    }
    st.sa_flags |= SA_ONSTACK;
    if (sigaction(signum, &st,  NULL) < 0) {
        goto fix_signal_error;
    }
    return;
fix_signal_error:
        fprintf(stderr, "error fixing handler for signal %d, please "
                "report this issue to "
                "https://github.com/wailsapp/wails: %s\n",
                signum, strerror(errno));
}

static void install_signal_handlers() {
	#if defined(SIGCHLD)
		fix_signal(SIGCHLD);
	#endif
	#if defined(SIGHUP)
		fix_signal(SIGHUP);
	#endif
	#if defined(SIGINT)
		fix_signal(SIGINT);
	#endif
	#if defined(SIGQUIT)
		fix_signal(SIGQUIT);
	#endif
	#if defined(SIGABRT)
		fix_signal(SIGABRT);
	#endif
	#if defined(SIGFPE)
		fix_signal(SIGFPE);
	#endif
	#if defined(SIGTERM)
		fix_signal(SIGTERM);
	#endif
	#if defined(SIGBUS)
		fix_signal(SIGBUS);
	#endif
	#if defined(SIGSEGV)
		fix_signal(SIGSEGV);
	#endif
	#if defined(SIGXCPU)
		fix_signal(SIGXCPU);
	#endif
	#if defined(SIGXFSZ)
		fix_signal(SIGXFSZ);
	#endif
}

static int GetNumScreens(){
    return 0;
}

static void on_data_received(GtkWidget *widget, GdkDragContext *context, gint x, gint y,
                      GtkSelectionData *selection_data, guint target_type, guint time,
                      gpointer data)
{
    gint length = gtk_selection_data_get_length(selection_data);

    if (length < 0)
    {
        g_print("DnD failed!\n");
        gtk_drag_finish(context, FALSE, FALSE, time);
    }

    gchar *uri_data = (gchar *)gtk_selection_data_get_data(selection_data);
    gchar **uri_list = g_uri_list_extract_uris(uri_data);

    onUriList(uri_list, data);

    g_strfreev(uri_list);
    gtk_drag_finish(context, TRUE, TRUE, time);
}

// drag and drop tutorial: https://wiki.gnome.org/Newcomers/OldDragNDropTutorial
static void enableDND(GtkWidget *widget, gpointer data)
{
    GtkTargetEntry *target = gtk_target_entry_new("text/uri-list", 0, 0);
    gtk_drag_dest_set(widget, GTK_DEST_DEFAULT_MOTION | GTK_DEST_DEFAULT_HIGHLIGHT | GTK_DEST_DEFAULT_DROP, target, 1, GDK_ACTION_COPY);

    signal_connect(widget, "drag-data-received", on_data_received, data);
}
*/
import "C"

// Calloc handles alloc/dealloc of C data
type Calloc struct {
	pool []unsafe.Pointer
}

// NewCalloc creates a new allocator
func NewCalloc() Calloc {
	return Calloc{}
}

// String creates a new C string and retains a reference to it
func (c Calloc) String(in string) *C.char {
	result := C.CString(in)
	c.pool = append(c.pool, unsafe.Pointer(result))
	return result
}

// Free frees all allocated C memory
func (c Calloc) Free() {
	for _, str := range c.pool {
		C.free(str)
	}
	c.pool = []unsafe.Pointer{}
}

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

var registerURIScheme sync.Once

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
	processApplicationEvent(C.uint(events.Linux.ApplicationStartup), data)
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint, data pointer) {
	event := newApplicationEvent(events.ApplicationEventType(eventID))

	//if data != nil {
	//	dataCStrJSON := C.serializationNSDictionary(data)
	//	if dataCStrJSON != nil {
	//		defer C.free(unsafe.Pointer(dataCStrJSON))
	//
	//		dataJSON := C.GoString(dataCStrJSON)
	//		var result map[string]any
	//		err := json.Unmarshal([]byte(dataJSON), &result)
	//
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		event.Context().setData(result)
	//	}
	//}

	switch event.Id {
	case uint(events.Linux.SystemThemeChanged):
		isDark := globalApplication.IsDarkMode()
		event.Context().setIsDarkMode(isDark)
	}
	applicationEvents <- event
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
	C.install_signal_handlers()

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

func setProgramName(prgName string) {
	cPrgName := C.CString(prgName)
	defer C.free(unsafe.Pointer(cPrgName))
	C.g_set_prgname(cPrgName)
}

func appRun(app pointer) error {
	application := (*C.GApplication)(app)
	//TODO: Only set this if we configure it to do so
	C.g_application_hold(application) // allows it to run without a window

	signal := C.CString("activate")
	defer C.free(unsafe.Pointer(signal))
	C.signal_connect(unsafe.Pointer(application), signal, C.activateLinux, nil)
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

func (w *linuxWebviewWindow) contextMenuSignals(menu pointer) {
	c := NewCalloc()
	defer c.Free()
	winID := unsafe.Pointer(uintptr(C.uint(w.parent.ID())))
	C.signal_connect(unsafe.Pointer(menu), c.String("button-release-event"), C.onMenuButtonEvent, winID)
}

func (w *linuxWebviewWindow) contextMenuShow(menu pointer, data *ContextMenuData) {
	geometry := C.GdkRectangle{
		x: C.int(data.X),
		y: C.int(data.Y),
	}
	event := C.GdkEvent{}
	gdkWindow := C.gtk_widget_get_window(w.gtkWidget())
	C.gtk_menu_popup_at_rect(
		(*C.GtkMenu)(menu),
		gdkWindow,
		(*C.GdkRectangle)(&geometry),
		C.GDK_GRAVITY_NORTH_WEST,
		C.GDK_GRAVITY_NORTH_WEST,
		(*C.GdkEvent)(&event),
	)
	w.ctxMenuOpened = true
}

func (a *linuxApp) getCurrentWindowID() uint {
	// TODO: Add extra metadata to window and use it!
	window := (*C.GtkWindow)(C.gtk_application_get_active_window((*C.GtkApplication)(a.application)))
	if window == nil {
		return uint(1)
	}
	identifier, ok := a.windowMap[window]
	if ok {
		return identifier
	}
	// FIXME: Should we panic here if not found?
	return uint(1)
}

func (a *linuxApp) getWindows() []pointer {
	result := []pointer{}
	windows := C.gtk_application_get_windows((*C.GtkApplication)(a.application))
	for {
		result = append(result, pointer(windows.data))
		windows = windows.next
		if windows == nil {
			return result
		}
	}
}

func (a *linuxApp) hideAllWindows() {
	for _, window := range a.getWindows() {
		C.gtk_widget_hide((*C.GtkWidget)(window))
	}
}

func (a *linuxApp) showAllWindows() {
	for _, window := range a.getWindows() {
		C.gtk_window_present((*C.GtkWindow)(window))
	}
}

func (a *linuxApp) setIcon(icon []byte) {
	gbytes := C.g_bytes_new_static(C.gconstpointer(unsafe.Pointer(&icon[0])), C.ulong(len(icon)))
	stream := C.g_memory_input_stream_new_from_bytes(gbytes)
	var gerror *C.GError
	pixbuf := C.gdk_pixbuf_new_from_stream(stream, nil, &gerror)
	if gerror != nil {
		a.parent.error("failed to load application icon: %s", C.GoString(gerror.message))
		C.g_error_free(gerror)
		return
	}

	a.icon = pointer(pixbuf)
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

func menuItemDestroy(widget pointer) {
	C.gtk_widget_destroy((*C.GtkWidget)(widget))
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
	name := C.gdk_monitor_get_model(monitor)
	return &Screen{
		ID:          fmt.Sprintf("%d", index),
		Name:        C.GoString(name),
		IsPrimary:   primary,
		ScaleFactor: float32(C.gdk_monitor_get_scale_factor(monitor)),
		X:           int(geometry.x),
		Y:           int(geometry.y),
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
		PhysicalBounds: Rect{
			X:      int(geometry.x),
			Y:      int(geometry.y),
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
		WorkArea: Rect{
			X:      int(geometry.x),
			Y:      int(geometry.y),
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
		PhysicalWorkArea: Rect{
			X:      int(geometry.x),
			Y:      int(geometry.y),
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
		Rotation: 0.0,
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

func (w *linuxWebviewWindow) setEnabled(enabled bool) {
	var value C.int
	if enabled {
		value = C.int(1)
	}

	C.gtk_widget_set_sensitive(w.gtkWidget(), value)
}

func widgetSetVisible(widget pointer, hidden bool) {
	if hidden {
		C.gtk_widget_hide((*C.GtkWidget)(widget))
	} else {
		C.gtk_widget_show((*C.GtkWidget)(widget))
	}
}

func (w *linuxWebviewWindow) close() {
	C.gtk_widget_destroy(w.gtkWidget())
	getNativeApplication().unregisterWindow(windowPointer(w.window))
}

func (w *linuxWebviewWindow) enableDND() {
	C.gtk_drag_dest_unset((*C.GtkWidget)(w.webview))

	windowId := C.uint(w.parent.id)
	C.enableDND((*C.GtkWidget)(w.vbox), C.gpointer(&windowId))
}

func (w *linuxWebviewWindow) execJS(js string) {
	InvokeAsync(func() {
		value := C.CString(js)
		C.webkit_web_view_evaluate_javascript(w.webKitWebView(),
			value,
			C.long(len(js)),
			nil,
			C.CString(""),
			nil,
			nil,
			nil)
		C.free(unsafe.Pointer(value))
	})
}

func getMousePosition() (int, int, *Screen) {
	var x, y C.gint
	var screen *C.GdkScreen
	defaultDisplay := C.gdk_display_get_default()
	device := C.gdk_seat_get_pointer(C.gdk_display_get_default_seat(defaultDisplay))
	C.gdk_device_get_position(device, &screen, &x, &y)
	// Get Monitor for screen
	monitor := C.gdk_display_get_monitor_at_point(defaultDisplay, x, y)
	geometry := C.GdkRectangle{}
	C.gdk_monitor_get_geometry(monitor, &geometry)
	scaleFactor := int(C.gdk_monitor_get_scale_factor(monitor))
	return int(x), int(y), &Screen{
		ID:          fmt.Sprintf("%d", 0),                         // A unique identifier for the display
		Name:        C.GoString(C.gdk_monitor_get_model(monitor)), // The name of the display
		ScaleFactor: float32(scaleFactor),                         // The scale factor of the display
		X:           int(geometry.x),                              // The x-coordinate of the top-left corner of the rectangle
		Y:           int(geometry.y),                              // The y-coordinate of the top-left corner of the rectangle
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
		WorkArea: Rect{
			X:      int(geometry.x),
			Y:      int(geometry.y),
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
		IsPrimary: false,
		Rotation:  0.0,
	}
}

func (w *linuxWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
	// Free menu
	if w.gtkmenu != nil {
		C.gtk_widget_destroy((*C.GtkWidget)(w.gtkmenu))
		w.gtkmenu = nil
	}
	// Free window
	C.gtk_widget_destroy(w.gtkWidget())
}

func (w *linuxWebviewWindow) fullscreen() {
	w.maximise()
	//w.lastWidth, w.lastHeight = w.size()
	x, y, width, height, scaleFactor := w.getCurrentMonitorGeometry()
	if x == -1 && y == -1 && width == -1 && height == -1 {
		return
	}
	w.setMinMaxSize(0, 0, width*scaleFactor, height*scaleFactor)
	w.setSize(width*scaleFactor, height*scaleFactor)
	C.gtk_window_fullscreen(w.gtkWindow())
	w.setRelativePosition(0, 0)
}

func (w *linuxWebviewWindow) getCurrentMonitor() *C.GdkMonitor {
	// Get the monitor that the window is currently on
	display := C.gtk_widget_get_display(w.gtkWidget())
	gdkWindow := C.gtk_widget_get_window(w.gtkWidget())
	if gdkWindow == nil {
		return nil
	}
	return C.gdk_display_get_monitor_at_window(display, gdkWindow)
}

func (w *linuxWebviewWindow) getScreen() (*Screen, error) {
	// Get the current screen for the window
	monitor := w.getCurrentMonitor()
	name := C.gdk_monitor_get_model(monitor)
	mx, my, width, height, scaleFactor := w.getCurrentMonitorGeometry()
	return &Screen{
		ID:          fmt.Sprintf("%d", w.id), // A unique identifier for the display
		Name:        C.GoString(name),        // The name of the display
		ScaleFactor: float32(scaleFactor),    // The scale factor of the display
		X:           mx,                      // The x-coordinate of the top-left corner of the rectangle
		Y:           my,                      // The y-coordinate of the top-left corner of the rectangle
		Size: Size{
			Height: int(height),
			Width:  int(width),
		},
		Bounds: Rect{
			X:      int(mx),
			Y:      int(my),
			Height: int(height),
			Width:  int(width),
		},
		WorkArea: Rect{
			X:      int(mx),
			Y:      int(my),
			Height: int(height),
			Width:  int(width),
		},
		PhysicalBounds: Rect{
			X:      int(mx),
			Y:      int(my),
			Height: int(height),
			Width:  int(width),
		},
		PhysicalWorkArea: Rect{
			X:      int(mx),
			Y:      int(my),
			Height: int(height),
			Width:  int(width),
		},
		IsPrimary: false,
		Rotation:  0.0,
	}, nil
}

func (w *linuxWebviewWindow) getCurrentMonitorGeometry() (x int, y int, width int, height int, scaleFactor int) {
	monitor := w.getCurrentMonitor()
	if monitor == nil {
		// Best effort to find screen resolution of default monitor
		display := C.gdk_display_get_default()
		monitor = C.gdk_display_get_primary_monitor(display)
		if monitor == nil {
			return -1, -1, -1, -1, 1
		}
	}
	var result C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &result)
	scaleFactor = int(C.gdk_monitor_get_scale_factor(monitor))
	return int(result.x), int(result.y), int(result.width), int(result.height), scaleFactor
}

func (w *linuxWebviewWindow) size() (int, int) {
	var windowWidth C.int
	var windowHeight C.int
	C.gtk_window_get_size(w.gtkWindow(), &windowWidth, &windowHeight)
	return int(windowWidth), int(windowHeight)
}

func (w *linuxWebviewWindow) relativePosition() (int, int) {
	x, y := w.position()
	// The position must be relative to the screen it is on
	// We need to get the screen it is on
	monitor := w.getCurrentMonitor()
	geometry := C.GdkRectangle{}
	C.gdk_monitor_get_geometry(monitor, &geometry)
	x = x - int(geometry.x)
	y = y - int(geometry.y)

	// TODO: Scale based on DPI

	return x, y
}

func (w *linuxWebviewWindow) gtkWidget() *C.GtkWidget {
	return (*C.GtkWidget)(w.window)
}

func (w *linuxWebviewWindow) hide() {
	// save position
	w.lastX, w.lastY = w.position()
	C.gtk_widget_hide(w.gtkWidget())
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	gdkWindow := C.gtk_widget_get_window(w.gtkWidget())
	state := C.gdk_window_get_state(gdkWindow)
	return state&C.GDK_WINDOW_STATE_FULLSCREEN > 0
}

func (w *linuxWebviewWindow) isFocused() bool {
	// returns true if window is focused
	return C.gtk_window_has_toplevel_focus(w.gtkWindow()) == 1
}

func (w *linuxWebviewWindow) isMaximised() bool {
	gdkwindow := C.gtk_widget_get_window(w.gtkWidget())
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_MAXIMIZED > 0 && state&C.GDK_WINDOW_STATE_FULLSCREEN == 0
}

func (w *linuxWebviewWindow) isMinimised() bool {
	gdkwindow := C.gtk_widget_get_window(w.gtkWidget())
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_ICONIFIED > 0
}

func (w *linuxWebviewWindow) isVisible() bool {
	if C.gtk_widget_is_visible(w.gtkWidget()) == 1 {
		return true
	}
	return false
}

func (w *linuxWebviewWindow) maximise() {
	C.gtk_window_maximize(w.gtkWindow())
}

func (w *linuxWebviewWindow) minimise() {
	C.gtk_window_iconify(w.gtkWindow())
}

func windowNew(application pointer, menu pointer, windowId uint, gpuPolicy WebviewGpuPolicy) (window, webview, vbox pointer) {
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

func windowNewWebview(parentId uint, gpuPolicy WebviewGpuPolicy) pointer {
	c := NewCalloc()
	defer c.Free()
	manager := C.webkit_user_content_manager_new()
	C.webkit_user_content_manager_register_script_message_handler(manager, c.String("external"))
	webView := C.webkit_web_view_new_with_user_content_manager(manager)

	// attach window id to both the webview and contentmanager
	C.save_window_id(unsafe.Pointer(webView), C.uint(parentId))
	C.save_window_id(unsafe.Pointer(manager), C.uint(parentId))

	registerURIScheme.Do(func() {
		context := C.webkit_web_view_get_context(C.webkit_web_view(webView))
		C.webkit_web_context_register_uri_scheme(
			context,
			c.String("wails"),
			C.WebKitURISchemeRequestCallback(C.onProcessRequest),
			nil,
			nil)
	})
	settings := C.webkit_web_view_get_settings((*C.WebKitWebView)(unsafe.Pointer(webView)))
	C.webkit_settings_set_user_agent_with_application_details(settings, c.String("wails.io"), c.String(""))

	switch gpuPolicy {
	case WebviewGpuPolicyAlways:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
		break
	case WebviewGpuPolicyOnDemand:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND)
		break
	case WebviewGpuPolicyNever:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER)
		break
	default:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND)
	}
	return pointer(webView)
}

func (w *linuxWebviewWindow) present() {
	C.gtk_window_present(w.gtkWindow())
	// gtk_window_unminimize (w.gtkWindow()) /// gtk4
}

func (w *linuxWebviewWindow) setSize(width, height int) {
	C.gtk_window_resize(
		w.gtkWindow(),
		C.gint(width),
		C.gint(height))
}

func (w *linuxWebviewWindow) show() {
	if w.gtkWidget() == nil {
		return
	}
	C.gtk_widget_show_all(w.gtkWidget())
	//w.setPosition(w.lastX, w.lastY)
}

func windowIgnoreMouseEvents(window pointer, webview pointer, ignore bool) {
	var enable C.int
	if ignore {
		enable = 1
	}
	gdkWindow := (*C.GdkWindow)(window)
	C.gdk_window_set_pass_through(gdkWindow, enable)
	C.webkit_web_view_set_editable((*C.WebKitWebView)(webview), C.gboolean(enable))
}

func (w *linuxWebviewWindow) webKitWebView() *C.WebKitWebView {
	return (*C.WebKitWebView)(w.webview)
}

func (w *linuxWebviewWindow) setBorderless(borderless bool) {
	C.gtk_window_set_decorated(w.gtkWindow(), gtkBool(!borderless))
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	C.gtk_window_set_resizable(w.gtkWindow(), gtkBool(resizable))
}

func (w *linuxWebviewWindow) setDefaultSize(width int, height int) {
	C.gtk_window_set_default_size(w.gtkWindow(), C.gint(width), C.gint(height))
}

func (w *linuxWebviewWindow) setBackgroundColour(colour RGBA) {
	rgba := C.GdkRGBA{C.double(colour.Red) / 255.0, C.double(colour.Green) / 255.0, C.double(colour.Blue) / 255.0, C.double(colour.Alpha) / 255.0}
	C.webkit_web_view_set_background_color((*C.WebKitWebView)(w.webview), &rgba)

	colour.Alpha = 255
	cssStr := C.CString(fmt.Sprintf("#webview-box {background-color: rgba(%d, %d, %d, %1.1f);}", colour.Red, colour.Green, colour.Blue, float32(colour.Alpha)/255.0))
	provider := C.gtk_css_provider_new()
	C.gtk_style_context_add_provider(
		C.gtk_widget_get_style_context((*C.GtkWidget)(w.vbox)),
		(*C.GtkStyleProvider)(unsafe.Pointer(provider)),
		C.GTK_STYLE_PROVIDER_PRIORITY_USER)
	C.g_object_unref(C.gpointer(provider))
	C.gtk_css_provider_load_from_data(provider, cssStr, -1, nil)
	C.free(unsafe.Pointer(cssStr))
}

func getPrimaryScreen() (*Screen, error) {
	display := C.gdk_display_get_default()
	monitor := C.gdk_display_get_primary_monitor(display)
	geometry := C.GdkRectangle{}
	C.gdk_monitor_get_geometry(monitor, &geometry)
	scaleFactor := int(C.gdk_monitor_get_scale_factor(monitor))
	// get the name for the screen
	name := C.gdk_monitor_get_model(monitor)
	return &Screen{
		ID:        "0",
		Name:      C.GoString(name),
		IsPrimary: true,
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
		ScaleFactor: float32(scaleFactor),
	}, nil
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

func (w *linuxWebviewWindow) setFrameless(frameless bool) {
	C.gtk_window_set_decorated(w.gtkWindow(), gtkBool(!frameless))
	// TODO: Deal with transparency for the titlebar if possible when !frameless
	//       Perhaps we just make it undecorated and add a menu bar inside?
}

// TODO: confirm this is working properly
func (w *linuxWebviewWindow) setHTML(html string) {
	cHTML := C.CString(html)
	uri := C.CString("wails://")
	empty := C.CString("")
	defer C.free(unsafe.Pointer(cHTML))
	defer C.free(unsafe.Pointer(uri))
	defer C.free(unsafe.Pointer(empty))
	C.webkit_web_view_load_alternate_html(
		w.webKitWebView(),
		cHTML,
		uri,
		empty)
}

func (w *linuxWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	C.gtk_window_set_keep_above(w.gtkWindow(), gtkBool(alwaysOnTop))
}

func (w *linuxWebviewWindow) flash(_ bool) {
	// Not supported on Linux
}

func (w *linuxWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		cTitle := C.CString(title)
		C.gtk_window_set_title(w.gtkWindow(), cTitle)
		C.free(unsafe.Pointer(cTitle))
	}
}

func (w *linuxWebviewWindow) setIcon(icon pointer) {
	if icon != nil {
		C.gtk_window_set_icon(w.gtkWindow(), (*C.GdkPixbuf)(icon))
	}
}

func (w *linuxWebviewWindow) gtkWindow() *C.GtkWindow {
	return (*C.GtkWindow)(w.window)
}

func (w *linuxWebviewWindow) setTransparent() {
	screen := C.gtk_widget_get_screen(w.gtkWidget())
	visual := C.gdk_screen_get_rgba_visual(screen)

	if visual != nil && C.gdk_screen_is_composited(screen) == C.int(1) {
		C.gtk_widget_set_app_paintable(w.gtkWidget(), C.gboolean(1))
		C.gtk_widget_set_visual(w.gtkWidget(), visual)
	}
}

func (w *linuxWebviewWindow) setURL(uri string) {
	target := C.CString(uri)
	C.webkit_web_view_load_uri(w.webKitWebView(), target)
	C.free(unsafe.Pointer(target))
}

//export emit
func emit(we *C.WindowEvent) {
	window := globalApplication.getWindowForID(uint(we.id))
	if window != nil {
		windowEvents <- &windowEvent{
			WindowID: window.ID(),
			EventID:  uint(events.WindowEventType(we.event)),
		}
	}
}

//export handleConfigureEvent
func handleConfigureEvent(widget *C.GtkWidget, event *C.GdkEventConfigure, data C.uintptr_t) C.gboolean {
	window := globalApplication.getWindowForID(uint(data))
	if window != nil {
		lw, ok := window.(*WebviewWindow).impl.(*linuxWebviewWindow)
		if !ok {
			return C.gboolean(1)
		}
		if lw.lastX != int(event.x) || lw.lastY != int(event.y) {
			lw.moveDebouncer(func() {
				processWindowEvent(C.uint(data), C.uint(events.Linux.WindowDidMove))
			})
		}

		if lw.lastWidth != int(event.width) || lw.lastHeight != int(event.height) {
			lw.resizeDebouncer(func() {
				processWindowEvent(C.uint(data), C.uint(events.Linux.WindowDidResize))
			})
		}

		lw.lastX = int(event.x)
		lw.lastY = int(event.y)
		lw.lastWidth = int(event.width)
		lw.lastHeight = int(event.height)
	}

	return C.gboolean(0)
}

//export handleDeleteEvent
func handleDeleteEvent(widget *C.GtkWidget, event *C.GdkEvent, data C.uintptr_t) C.gboolean {
	processWindowEvent(C.uint(data), C.uint(events.Linux.WindowDeleteEvent))
	return C.gboolean(1)
}

//export handleFocusEvent
func handleFocusEvent(widget *C.GtkWidget, event *C.GdkEvent, data C.uintptr_t) C.gboolean {
	focusEvent := (*C.GdkEventFocus)(unsafe.Pointer(event))
	if focusEvent._type == C.GDK_FOCUS_CHANGE {
		if focusEvent.in == C.TRUE {
			processWindowEvent(C.uint(data), C.uint(events.Linux.WindowFocusIn))
		} else {
			processWindowEvent(C.uint(data), C.uint(events.Linux.WindowFocusOut))
		}
	}
	return C.gboolean(0)
}

//export handleLoadChanged
func handleLoadChanged(webview *C.WebKitWebView, event C.WebKitLoadEvent, data C.uintptr_t) {
	switch event {
	case C.WEBKIT_LOAD_FINISHED:
		processWindowEvent(C.uint(data), C.uint(events.Linux.WindowLoadChanged))
	}
}

func (w *linuxWebviewWindow) setupSignalHandlers(emit func(e events.WindowEventType)) {

	c := NewCalloc()
	defer c.Free()

	winID := unsafe.Pointer(uintptr(C.uint(w.parent.ID())))

	// Set up the window close event
	wv := unsafe.Pointer(w.webview)
	C.signal_connect(unsafe.Pointer(w.window), c.String("delete-event"), C.handleDeleteEvent, winID)
	C.signal_connect(unsafe.Pointer(w.window), c.String("focus-out-event"), C.handleFocusEvent, winID)
	C.signal_connect(wv, c.String("load-changed"), C.handleLoadChanged, winID)
	C.signal_connect(unsafe.Pointer(w.window), c.String("configure-event"), C.handleConfigureEvent, winID)

	contentManager := C.webkit_web_view_get_user_content_manager(w.webKitWebView())
	C.signal_connect(unsafe.Pointer(contentManager), c.String("script-message-received::external"), C.sendMessageToBackend, nil)
	C.signal_connect(wv, c.String("button-press-event"), C.onButtonEvent, winID)
	C.signal_connect(wv, c.String("button-release-event"), C.onButtonEvent, winID)
	C.signal_connect(wv, c.String("key-press-event"), C.onKeyPressEvent, winID)
}

func getMouseButtons() (bool, bool, bool) {
	var pointer *C.GdkDevice
	var state C.GdkModifierType
	pointer = C.gdk_seat_get_pointer(C.gdk_display_get_default_seat(C.gdk_display_get_default()))
	C.gdk_device_get_state(pointer, nil, nil, &state)
	return state&C.GDK_BUTTON1_MASK > 0, state&C.GDK_BUTTON2_MASK > 0, state&C.GDK_BUTTON3_MASK > 0
}

func openDevTools(webview pointer) {
	inspector := C.webkit_web_view_get_inspector((*C.WebKitWebView)(webview))
	C.webkit_web_inspector_show(inspector)
}

func (w *linuxWebviewWindow) startDrag() error {
	C.gtk_window_begin_move_drag(
		(*C.GtkWindow)(w.window),
		C.int(w.drag.MouseButton),
		C.int(w.drag.XRoot),
		C.int(w.drag.YRoot),
		C.uint32_t(w.drag.DragTime))
	return nil
}

func enableDevTools(webview pointer) {
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

func (w *linuxWebviewWindow) unfullscreen() {
	C.gtk_window_unfullscreen((*C.GtkWindow)(w.window))
	w.unmaximise()
}

func (w *linuxWebviewWindow) unmaximise() {
	C.gtk_window_unmaximize((*C.GtkWindow)(w.window))
}

func (w *linuxWebviewWindow) getZoom() float64 {
	return float64(C.webkit_web_view_get_zoom_level(w.webKitWebView()))
}

func (w *linuxWebviewWindow) zoomIn() {
	// FIXME: ZoomIn/Out is assumed to be incorrect!
	ZoomInFactor := 1.10
	w.setZoom(w.getZoom() * ZoomInFactor)
}

func (w *linuxWebviewWindow) zoomOut() {
	ZoomInFactor := -1.10
	w.setZoom(w.getZoom() * ZoomInFactor)
}

func (w *linuxWebviewWindow) zoomReset() {
	w.setZoom(1.0)
}

func (w *linuxWebviewWindow) reload() {
	uri := C.CString("wails://")
	C.webkit_web_view_load_uri(w.webKitWebView(), uri)
	C.free(unsafe.Pointer(uri))
}

func (w *linuxWebviewWindow) setZoom(zoom float64) {
	if zoom < 1 { // 1.0 is the smallest allowable
		zoom = 1
	}
	C.webkit_web_view_set_zoom_level(w.webKitWebView(), C.double(zoom))
}

func (w *linuxWebviewWindow) move(x, y int) {
	// Move the window to these coordinates
	C.gtk_window_move(w.gtkWindow(), C.int(x), C.int(y))
}

func (w *linuxWebviewWindow) position() (int, int) {
	var x C.int
	var y C.int
	C.gtk_window_get_position((*C.GtkWindow)(w.window), &x, &y)
	return int(x), int(y)
}

func (w *linuxWebviewWindow) ignoreMouse(ignore bool) {
	if ignore {
		C.gtk_widget_set_events((*C.GtkWidget)(unsafe.Pointer(w.window)), C.GDK_ENTER_NOTIFY_MASK|C.GDK_LEAVE_NOTIFY_MASK)
	} else {
		C.gtk_widget_set_events((*C.GtkWidget)(unsafe.Pointer(w.window)), C.GDK_ALL_EVENTS_MASK)
	}
}

// FIXME Change this to reflect mouse button!
//
//export onButtonEvent
func onButtonEvent(_ *C.GtkWidget, event *C.GdkEventButton, data C.uintptr_t) C.gboolean {
	// Constants (defined here to be easier to use with purego)
	GdkButtonPress := C.GDK_BUTTON_PRESS     // 4
	Gdk2ButtonPress := C.GDK_2BUTTON_PRESS   // 5 for double-click
	GdkButtonRelease := C.GDK_BUTTON_RELEASE // 7

	windowId := uint(C.uint(data))
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

//export onMenuButtonEvent
func onMenuButtonEvent(_ *C.GtkWidget, event *C.GdkEventButton, data C.uintptr_t) C.gboolean {
	// Constants (defined here to be easier to use with purego)
	GdkButtonRelease := C.GDK_BUTTON_RELEASE // 7

	windowId := uint(C.uint(data))
	window := globalApplication.getWindowForID(windowId)
	if window == nil {
		return C.gboolean(0)
	}
	lw, ok := (window.(*WebviewWindow).impl).(*linuxWebviewWindow)
	if !ok {
		return C.gboolean(0)
	}

	// prevent custom context menu from closing immediately
	if event.button == 3 && int(event._type) == GdkButtonRelease && lw.ctxMenuOpened {
		lw.ctxMenuOpened = false
		return C.gboolean(1)
	}

	return C.gboolean(0)
}

//export onUriList
func onUriList(extracted **C.char, data unsafe.Pointer) {
	// Credit: https://groups.google.com/g/golang-nuts/c/bI17Bpck8K4/m/DVDa7EMtDAAJ
	offset := unsafe.Sizeof(uintptr(0))
	filenames := []string{}
	for *extracted != nil {
		filenames = append(filenames, strings.TrimPrefix(C.GoString(*extracted), "file://"))
		extracted = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(extracted)) + offset))
	}

	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:  uint(*((*C.uint)(data))),
		filenames: filenames,
	}
}

var debounceTimer *time.Timer
var isDebouncing bool = false

//export onKeyPressEvent
func onKeyPressEvent(_ *C.GtkWidget, event *C.GdkEventKey, userData C.uintptr_t) C.gboolean {
	// Keypress re-emits if the key is pressed over a certain threshold so we need a debounce
	if isDebouncing {
		debounceTimer.Reset(50 * time.Millisecond)
		return C.gboolean(0)
	}

	// Start the debounce
	isDebouncing = true
	debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
		isDebouncing = false
	})

	windowID := uint(C.uint(userData))
	if accelerator, ok := getKeyboardState(event); ok {
		windowKeyEvents <- &windowKeyEvent{
			windowId:          windowID,
			acceleratorString: accelerator,
		}
	}
	return C.gboolean(0)
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
		return "", false
	}
	acc.Key = keyString
	return acc.String(), true
}

//export onProcessRequest
func onProcessRequest(request *C.WebKitURISchemeRequest, data C.uintptr_t) {
	webView := C.webkit_uri_scheme_request_get_web_view(request)
	windowId := uint(C.get_window_id(unsafe.Pointer(webView)))
	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(unsafe.Pointer(request)),
		windowId:   windowId,
		windowName: globalApplication.getWindowForID(windowId).Name(),
	}
}

//export sendMessageToBackend
func sendMessageToBackend(contentManager *C.WebKitUserContentManager, result *C.WebKitJavascriptResult,
	data unsafe.Pointer) {

	// Get the windowID from the contentManager
	thisWindowID := uint(C.get_window_id(unsafe.Pointer(contentManager)))

	var msg string
	value := C.webkit_javascript_result_get_js_value(result)
	message := C.jsc_value_to_string(value)
	msg = C.GoString(message)
	defer C.g_free(C.gpointer(message))
	windowMessageBuffer <- &windowMessage{
		windowId: thisWindowID,
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
		f := C.gtk_file_filter_new()
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
		response := C.gtk_dialog_run((*C.GtkDialog)(fc))
		go func() {
			defer handlePanic()
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
			}
			close(selections)
		}()
	})
	C.gtk_widget_destroy((*C.GtkWidget)(unsafe.Pointer(fc)))
	return selections, nil
}

func runOpenFileDialog(dialog *OpenFileDialogStruct) (chan string, error) {
	var action int

	if dialog.canChooseDirectories {
		action = C.GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER
	} else {
		action = C.GTK_FILE_CHOOSER_ACTION_OPEN
	}

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
		action,
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

func (w *linuxWebviewWindow) cut() {
	//C.webkit_web_view_execute_editing_command(w.webview, C.WEBKIT_EDITING_COMMAND_CUT)
}

func (w *linuxWebviewWindow) paste() {
	//C.webkit_web_view_execute_editing_command(w.webview, C.WEBKIT_EDITING_COMMAND_PASTE)
}

func (w *linuxWebviewWindow) copy() {
	//C.webkit_web_view_execute_editing_command(w.webview, C.WEBKIT_EDITING_COMMAND_COPY)
}

func (w *linuxWebviewWindow) selectAll() {
	//C.webkit_web_view_execute_editing_command(w.webview, C.WEBKIT_EDITING_COMMAND_SELECT_ALL)
}

func (w *linuxWebviewWindow) undo() {
	//C.webkit_web_view_execute_editing_command(w.webview, C.WEBKIT_EDITING_COMMAND_UNDO)
}

func (w *linuxWebviewWindow) redo() {
}

func (w *linuxWebviewWindow) delete() {
}
