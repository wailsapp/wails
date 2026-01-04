//go:build linux && cgo && !gtk3 && !android

package application

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"

	"github.com/wailsapp/wails/v3/pkg/events"
)

/*
#cgo linux pkg-config: gtk4 webkitgtk-6.0

#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>

// Use NON_UNIQUE to allow multiple instances of the application to run.
#define APPLICATION_DEFAULT_FLAGS G_APPLICATION_DEFAULT_FLAGS

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

static void save_webview_to_content_manager(void *contentManager, void *webview)
{
    g_object_set_data(G_OBJECT((WebKitUserContentManager *)contentManager), "webview", webview);
}

static WebKitWebView* get_webview_from_content_manager(void *contentManager)
{
	return WEBKIT_WEB_VIEW(g_object_get_data(G_OBJECT(contentManager), "webview"));
}

static guint get_window_id(void *object)
{
    return GPOINTER_TO_UINT(g_object_get_data((GObject *)object, "windowid"));
}

// exported below
void activateLinux(gpointer data);
extern void emit(WindowEvent* data);
extern gboolean handleCloseRequest(GtkWindow*, uintptr_t);
extern void handleNotifyState(GObject*, GParamSpec*, uintptr_t);
extern gboolean handleFocusEnter(GtkEventController*, uintptr_t);
extern gboolean handleFocusLeave(GtkEventController*, uintptr_t);
extern void handleLoadChanged(WebKitWebView*, WebKitLoadEvent, uintptr_t);
extern void handleButtonPressed(GtkGestureClick*, gint, gdouble, gdouble, uintptr_t);
extern void handleButtonReleased(GtkGestureClick*, gint, gdouble, gdouble, uintptr_t);
extern gboolean handleKeyPressed(GtkEventControllerKey*, guint, guint, GdkModifierType, uintptr_t);
void handleClick(void*);
extern void onProcessRequest(WebKitURISchemeRequest *request, uintptr_t user_data);
extern void sendMessageToBackend(WebKitUserContentManager *contentManager, void *result, void *data);
// exported below (end)

static void signal_connect(void *widget, char *event, void *cb, void* data) {
   // g_signal_connect is a macro and can't be called directly
   g_signal_connect(widget, event, cb, data);
}

static WebKitWebView* webkit_web_view(GtkWidget *webview) {
	return WEBKIT_WEB_VIEW(webview);
}

// GTK4: Window positioning is NO-OP on Wayland (documented limitation)
// These functions exist for API compatibility but may not have effect on Wayland

typedef struct Screen {
	const char* id;
	const char* name;
	int p_width;
	int p_height;
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

// GTK4 Menu System - uses GMenu/GAction instead of GtkMenu/GtkMenuItem
// Each menu item has an associated GSimpleAction in an action group

typedef struct MenuItemData {
    guint id;
    GSimpleAction *action;
} MenuItemData;

static GMenu *app_menu_model = NULL;
static GSimpleActionGroup *app_action_group = NULL;

extern void menuActionActivated(guint id);

static void on_action_activated(GSimpleAction *action, GVariant *parameter, gpointer user_data) {
    MenuItemData *data = (MenuItemData *)user_data;
    if (data != NULL) {
        menuActionActivated(data->id);
    }
}

static void init_app_action_group() {
    if (app_action_group == NULL) {
        app_action_group = g_simple_action_group_new();
    }
}

static GMenuItem* create_menu_item(const char *label, const char *action_name, guint item_id) {
    init_app_action_group();

    char full_action_name[256];
    snprintf(full_action_name, sizeof(full_action_name), "app.%s", action_name);

    GMenuItem *item = g_menu_item_new(label, full_action_name);

    GSimpleAction *action = g_simple_action_new(action_name, NULL);
    MenuItemData *data = g_new0(MenuItemData, 1);
    data->id = item_id;
    data->action = action;
    g_signal_connect(action, "activate", G_CALLBACK(on_action_activated), data);
    g_action_map_add_action(G_ACTION_MAP(app_action_group), G_ACTION(action));

    return item;
}

static GMenuItem* create_check_menu_item(const char *label, const char *action_name, guint item_id, gboolean initial_state) {
    init_app_action_group();

    char full_action_name[256];
    snprintf(full_action_name, sizeof(full_action_name), "app.%s", action_name);

    GMenuItem *item = g_menu_item_new(label, full_action_name);

    GSimpleAction *action = g_simple_action_new_stateful(action_name, NULL, g_variant_new_boolean(initial_state));
    MenuItemData *data = g_new0(MenuItemData, 1);
    data->id = item_id;
    data->action = action;
    g_signal_connect(action, "activate", G_CALLBACK(on_action_activated), data);
    g_action_map_add_action(G_ACTION_MAP(app_action_group), G_ACTION(action));

    return item;
}

static GtkWidget* create_menu_bar_from_model(GMenu *menu_model) {
    return gtk_popover_menu_bar_new_from_model(G_MENU_MODEL(menu_model));
}

static void attach_action_group_to_widget(GtkWidget *widget) {
    init_app_action_group();
    gtk_widget_insert_action_group(widget, "app", G_ACTION_GROUP(app_action_group));
}

static void set_action_enabled(const char *action_name, gboolean enabled) {
    if (app_action_group == NULL) return;
    GAction *action = g_action_map_lookup_action(G_ACTION_MAP(app_action_group), action_name);
    if (action != NULL && G_IS_SIMPLE_ACTION(action)) {
        g_simple_action_set_enabled(G_SIMPLE_ACTION(action), enabled);
    }
}

static void set_action_state(const char *action_name, gboolean state) {
    if (app_action_group == NULL) return;
    GAction *action = g_action_map_lookup_action(G_ACTION_MAP(app_action_group), action_name);
    if (action != NULL && G_IS_SIMPLE_ACTION(action)) {
        g_simple_action_set_state(G_SIMPLE_ACTION(action), g_variant_new_boolean(state));
    }
}

static gboolean get_action_state(const char *action_name) {
    if (app_action_group == NULL) return FALSE;
    GAction *action = g_action_map_lookup_action(G_ACTION_MAP(app_action_group), action_name);
    if (action != NULL) {
        GVariant *state = g_action_get_state(action);
        if (state != NULL) {
            gboolean result = g_variant_get_boolean(state);
            g_variant_unref(state);
            return result;
        }
    }
    return FALSE;
}

// GTK4 uses GtkEventController for events instead of direct signal handlers
static void setupWindowEventControllers(GtkWindow *window, GtkWidget *webview, uintptr_t winID) {
    // Close request (replaces delete-event)
    g_signal_connect(window, "close-request", G_CALLBACK(handleCloseRequest), (gpointer)winID);

    // Window state changes (maximize, fullscreen, etc)
    g_signal_connect(window, "notify::maximized", G_CALLBACK(handleNotifyState), (gpointer)winID);
    g_signal_connect(window, "notify::fullscreened", G_CALLBACK(handleNotifyState), (gpointer)winID);

    // Focus controller for window
    GtkEventController *focus_controller = gtk_event_controller_focus_new();
    gtk_widget_add_controller(GTK_WIDGET(window), focus_controller);
    g_signal_connect(focus_controller, "enter", G_CALLBACK(handleFocusEnter), (gpointer)winID);
    g_signal_connect(focus_controller, "leave", G_CALLBACK(handleFocusLeave), (gpointer)winID);

    // Click gesture for webview (button press/release)
    GtkGesture *click_gesture = gtk_gesture_click_new();
    gtk_gesture_single_set_button(GTK_GESTURE_SINGLE(click_gesture), 0); // Listen to all buttons
    gtk_widget_add_controller(webview, GTK_EVENT_CONTROLLER(click_gesture));
    g_signal_connect(click_gesture, "pressed", G_CALLBACK(handleButtonPressed), (gpointer)winID);
    g_signal_connect(click_gesture, "released", G_CALLBACK(handleButtonReleased), (gpointer)winID);

    // Key controller for webview
    GtkEventController *key_controller = gtk_event_controller_key_new();
    gtk_widget_add_controller(webview, key_controller);
    g_signal_connect(key_controller, "key-pressed", G_CALLBACK(handleKeyPressed), (gpointer)winID);
}

// GTK4 window drag using GdkToplevel
static void beginWindowDrag(GtkWindow *window, int button, double x, double y, guint32 timestamp) {
    GtkNative *native = gtk_widget_get_native(GTK_WIDGET(window));
    if (native == NULL) return;

    GdkSurface *surface = gtk_native_get_surface(native);
    if (surface == NULL || !GDK_IS_TOPLEVEL(surface)) return;

    GdkToplevel *toplevel = GDK_TOPLEVEL(surface);
    GdkDevice *device = NULL;
    GdkDisplay *display = gdk_surface_get_display(surface);
    GdkSeat *seat = gdk_display_get_default_seat(display);
    if (seat) {
        device = gdk_seat_get_pointer(seat);
    }

    gdk_toplevel_begin_move(toplevel, device, button, x, y, timestamp);
}

// GTK4 window resize using GdkToplevel
static void beginWindowResize(GtkWindow *window, GdkSurfaceEdge edge, int button, double x, double y, guint32 timestamp) {
    GtkNative *native = gtk_widget_get_native(GTK_WIDGET(window));
    if (native == NULL) return;

    GdkSurface *surface = gtk_native_get_surface(native);
    if (surface == NULL || !GDK_IS_TOPLEVEL(surface)) return;

    GdkToplevel *toplevel = GDK_TOPLEVEL(surface);
    GdkDevice *device = NULL;
    GdkDisplay *display = gdk_surface_get_display(surface);
    GdkSeat *seat = gdk_display_get_default_seat(display);
    if (seat) {
        device = gdk_seat_get_pointer(seat);
    }

    gdk_toplevel_begin_resize(toplevel, edge, device, button, x, y, timestamp);
}

// GTK4 drag-and-drop uses GtkDropTarget instead of GTK3's drag signals
extern void onDropEnter(uintptr_t);
extern void onDropLeave(uintptr_t);
extern void onDropMotion(gint, gint, uintptr_t);
extern void onDropFiles(char**, gint, gint, uintptr_t);

static GdkDragAction on_drop_enter(GtkDropTarget *target, gdouble x, gdouble y, gpointer data) {
    onDropEnter((uintptr_t)data);
    return GDK_ACTION_COPY;
}

static void on_drop_leave(GtkDropTarget *target, gpointer data) {
    onDropLeave((uintptr_t)data);
}

static GdkDragAction on_drop_motion(GtkDropTarget *target, gdouble x, gdouble y, gpointer data) {
    onDropMotion((gint)x, (gint)y, (uintptr_t)data);
    return GDK_ACTION_COPY;
}

static gboolean on_drop(GtkDropTarget *target, const GValue *value, gdouble x, gdouble y, gpointer data) {
    if (!G_VALUE_HOLDS(value, GDK_TYPE_FILE_LIST)) {
        return FALSE;
    }

    GSList *file_list = g_value_get_boxed(value);
    if (file_list == NULL) {
        return FALSE;
    }

    // Count files
    guint count = g_slist_length(file_list);
    if (count == 0) {
        return FALSE;
    }

    // Build array of file paths
    char **paths = g_new0(char*, count + 1);
    guint i = 0;
    for (GSList *l = file_list; l != NULL; l = l->next) {
        GFile *file = G_FILE(l->data);
        paths[i++] = g_file_get_path(file);
    }
    paths[count] = NULL;

    onDropFiles(paths, (gint)x, (gint)y, (uintptr_t)data);

    // Cleanup
    for (i = 0; i < count; i++) {
        g_free(paths[i]);
    }
    g_free(paths);

    return TRUE;
}

static void enableDND(GtkWidget *widget, gpointer data) {
    GtkDropTarget *target = gtk_drop_target_new(GDK_TYPE_FILE_LIST, GDK_ACTION_COPY);
    g_signal_connect(target, "enter", G_CALLBACK(on_drop_enter), data);
    g_signal_connect(target, "leave", G_CALLBACK(on_drop_leave), data);
    g_signal_connect(target, "motion", G_CALLBACK(on_drop_motion), data);
    g_signal_connect(target, "drop", G_CALLBACK(on_drop), data);
    gtk_widget_add_controller(widget, GTK_EVENT_CONTROLLER(target));
}

static void disableDND(GtkWidget *widget, gpointer data) {
    // In GTK4, we don't add a drop target to block drops
    // The default behavior is to not accept drops
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

// getLinuxWebviewWindow safely extracts a linuxWebviewWindow from a Window interface
func getLinuxWebviewWindow(window Window) *linuxWebviewWindow {
	if window == nil {
		return nil
	}

	webviewWindow, ok := window.(*WebviewWindow)
	if !ok {
		return nil
	}

	lw, ok := webviewWindow.impl.(*linuxWebviewWindow)
	if !ok {
		return nil
	}

	return lw
}

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

	switch event.Id {
	case uint(events.Linux.SystemThemeChanged):
		isDark := globalApplication.Env.IsDarkMode()
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
	C.g_application_hold(application)

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
	// GTK4: Context menus use GtkPopoverMenu, signals handled differently
	// TODO: Implement GTK4 context menu signal handling
}

func (w *linuxWebviewWindow) contextMenuShow(menu pointer, data *ContextMenuData) {
	// GTK4: Use GtkPopoverMenu instead of gtk_menu_popup_at_rect
	// TODO: Implement GTK4 context menu popup
}

func (a *linuxApp) getCurrentWindowID() uint {
	window := (*C.GtkWindow)(C.gtk_application_get_active_window((*C.GtkApplication)(a.application)))
	if window == nil {
		return uint(1)
	}
	identifier, ok := a.windowMap[window]
	if ok {
		return identifier
	}
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
		C.gtk_widget_set_visible((*C.GtkWidget)(window), C.gboolean(0))
	}
}

func (a *linuxApp) showAllWindows() {
	for _, window := range a.getWindows() {
		C.gtk_window_present((*C.GtkWindow)(window))
	}
}

func (a *linuxApp) setIcon(icon []byte) {
	// TODO: Implement GTK4 icon setting using GdkTexture
	gbytes := C.g_bytes_new_static(C.gconstpointer(unsafe.Pointer(&icon[0])), C.ulong(len(icon)))
	defer C.g_bytes_unref(gbytes)
}

// Clipboard - GTK4 uses GdkClipboard API
func clipboardGet() string {
	display := C.gdk_display_get_default()
	clip := C.gdk_display_get_clipboard(display)
	// GTK4: Async clipboard API - this is a simplified sync version
	// TODO: Implement proper async clipboard for GTK4
	_ = clip
	return ""
}

func clipboardSet(text string) {
	display := C.gdk_display_get_default()
	clip := C.gdk_display_get_clipboard(display)
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gdk_clipboard_set_text(clip, cText)
}

// Menu - GTK4 uses GMenu/GAction instead of GtkMenu

var menuItemActionCounter uint32 = 0
var menuItemActions = make(map[uint]string)

func generateActionName(itemId uint) string {
	menuItemActionCounter++
	name := fmt.Sprintf("action_%d", menuItemActionCounter)
	menuItemActions[itemId] = name
	return name
}

//export menuActionActivated
func menuActionActivated(id C.guint) {
	item, ok := gtkSignalToMenuItem[uint(id)]
	if !ok {
		return
	}
	switch item.itemType {
	case text:
		menuItemClicked <- item.id
	case checkbox:
		impl := item.impl.(*linuxMenuItem)
		currentState := impl.isChecked()
		impl.setChecked(!currentState)
		menuItemClicked <- item.id
	case radio:
		menuItem := item.impl.(*linuxMenuItem)
		if !menuItem.isChecked() {
			menuItem.setChecked(true)
			menuItemClicked <- item.id
		}
	}
}

func menuAddSeparator(menu *Menu) {
	if menu.impl == nil {
		return
	}
	impl := menu.impl.(*linuxMenu)
	if impl.native == nil {
		return
	}
	gmenu := (*C.GMenu)(impl.native)
	section := C.g_menu_new()
	C.g_menu_append_section(gmenu, nil, (*C.GMenuModel)(unsafe.Pointer(section)))
}

func menuAppend(parent *Menu, menu *MenuItem) {
	if parent.impl == nil || menu.impl == nil {
		return
	}
	parentImpl := parent.impl.(*linuxMenu)
	menuImpl := menu.impl.(*linuxMenuItem)
	if parentImpl.native == nil || menuImpl.native == nil {
		return
	}
	gmenu := (*C.GMenu)(parentImpl.native)
	gitem := (*C.GMenuItem)(menuImpl.native)
	C.g_menu_append_item(gmenu, gitem)
}

func menuBarNew() pointer {
	gmenu := C.g_menu_new()
	C.app_menu_model = gmenu
	return pointer(gmenu)
}

func menuNew() pointer {
	return pointer(C.g_menu_new())
}

func menuSetSubmenu(item *MenuItem, menu *Menu) {
	if item.impl == nil || menu.impl == nil {
		return
	}
	itemImpl := item.impl.(*linuxMenuItem)
	menuImpl := menu.impl.(*linuxMenu)
	if itemImpl.native == nil || menuImpl.native == nil {
		return
	}
	gitem := (*C.GMenuItem)(itemImpl.native)
	gmenu := (*C.GMenu)(menuImpl.native)
	C.g_menu_item_set_submenu(gitem, (*C.GMenuModel)(unsafe.Pointer(gmenu)))
}

func menuGetRadioGroup(item *linuxMenuItem) *GSList {
	return nil
}

//export handleClick
func handleClick(idPtr unsafe.Pointer) {
}

func attachMenuHandler(item *MenuItem) uint {
	gtkSignalToMenuItem[item.id] = item
	return item.id
}

func menuItemChecked(widget pointer) bool {
	if widget == nil {
		return false
	}
	cKey := C.CString("item_id")
	defer C.free(unsafe.Pointer(cKey))
	itemId := uint(uintptr(C.g_object_get_data((*C.GObject)(widget), cKey)))
	actionName, ok := menuItemActions[itemId]
	if !ok {
		return false
	}
	cName := C.CString(actionName)
	defer C.free(unsafe.Pointer(cName))
	return C.get_action_state(cName) != 0
}

func menuItemNew(label string, bitmap []byte) pointer {
	return nil
}

func menuItemNewWithId(label string, bitmap []byte, itemId uint) pointer {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	actionName := generateActionName(itemId)
	cAction := C.CString(actionName)
	defer C.free(unsafe.Pointer(cAction))

	gitem := C.create_menu_item(cLabel, cAction, C.guint(itemId))

	cKey := C.CString("item_id")
	defer C.free(unsafe.Pointer(cKey))
	C.g_object_set_data((*C.GObject)(unsafe.Pointer(gitem)), cKey, C.gpointer(uintptr(itemId)))
	return pointer(gitem)
}

func menuItemDestroy(widget pointer) {
	if widget != nil {
		C.g_object_unref(C.gpointer(widget))
	}
}

func menuItemAddProperties(menuItem *C.GtkWidget, label string, bitmap []byte) pointer {
	return nil
}

func menuCheckItemNew(label string, bitmap []byte) pointer {
	return nil
}

func menuCheckItemNewWithId(label string, bitmap []byte, itemId uint, checked bool) pointer {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	actionName := generateActionName(itemId)
	cAction := C.CString(actionName)
	defer C.free(unsafe.Pointer(cAction))

	initialState := C.gboolean(0)
	if checked {
		initialState = C.gboolean(1)
	}

	gitem := C.create_check_menu_item(cLabel, cAction, C.guint(itemId), initialState)

	cKey := C.CString("item_id")
	defer C.free(unsafe.Pointer(cKey))
	C.g_object_set_data((*C.GObject)(unsafe.Pointer(gitem)), cKey, C.gpointer(uintptr(itemId)))
	return pointer(gitem)
}

func menuItemSetChecked(widget pointer, checked bool) {
	if widget == nil {
		return
	}
	cKey := C.CString("item_id")
	defer C.free(unsafe.Pointer(cKey))
	itemId := uint(uintptr(C.g_object_get_data((*C.GObject)(widget), cKey)))
	actionName, ok := menuItemActions[itemId]
	if !ok {
		return
	}
	cName := C.CString(actionName)
	defer C.free(unsafe.Pointer(cName))
	state := C.gboolean(0)
	if checked {
		state = C.gboolean(1)
	}
	C.set_action_state(cName, state)
}

func menuItemSetDisabled(widget pointer, disabled bool) {
	if widget == nil {
		return
	}
	cKey := C.CString("item_id")
	defer C.free(unsafe.Pointer(cKey))
	itemId := uint(uintptr(C.g_object_get_data((*C.GObject)(widget), cKey)))
	actionName, ok := menuItemActions[itemId]
	if !ok {
		return
	}
	cName := C.CString(actionName)
	defer C.free(unsafe.Pointer(cName))
	enabled := C.gboolean(1)
	if disabled {
		enabled = C.gboolean(0)
	}
	C.set_action_enabled(cName, enabled)
}

func menuItemSetLabel(widget pointer, label string) {
	if widget == nil {
		return
	}
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	C.g_menu_item_set_label((*C.GMenuItem)(widget), cLabel)
}

func menuItemRemoveBitmap(widget pointer) {
}

func menuItemSetBitmap(widget pointer, bitmap []byte) {
}

func menuItemSetToolTip(widget pointer, tooltip string) {
}

func menuItemSignalBlock(widget pointer, handlerId uint, block bool) {
}

func menuRadioItemNew(group *GSList, label string) pointer {
	return nil
}

func menuRadioItemNewWithId(label string, itemId uint, checked bool) pointer {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	actionName := generateActionName(itemId)
	cAction := C.CString(actionName)
	defer C.free(unsafe.Pointer(cAction))

	initialState := C.gboolean(0)
	if checked {
		initialState = C.gboolean(1)
	}

	gitem := C.create_check_menu_item(cLabel, cAction, C.guint(itemId), initialState)

	cKey := C.CString("item_id")
	defer C.free(unsafe.Pointer(cKey))
	C.g_object_set_data((*C.GObject)(unsafe.Pointer(gitem)), cKey, C.gpointer(uintptr(itemId)))
	return pointer(gitem)
}

// screen related
func getScreenByIndex(display *C.GdkDisplay, index int) *Screen {
	monitors := C.gdk_display_get_monitors(display)
	monitor := (*C.GdkMonitor)(C.g_list_model_get_item(monitors, C.guint(index)))
	if monitor == nil {
		return nil
	}
	defer C.g_object_unref(C.gpointer(monitor))

	var geometry C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &geometry)
	name := C.gdk_monitor_get_model(monitor)
	return &Screen{
		ID:          fmt.Sprintf("%d", index),
		Name:        C.GoString(name),
		IsPrimary:   false, // GTK4 doesn't have gdk_monitor_is_primary
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
	display := C.gdk_display_get_default()
	monitors := C.gdk_display_get_monitors(display)
	count := C.g_list_model_get_n_items(monitors)
	for i := 0; i < int(count); i++ {
		screens = append(screens, getScreenByIndex(display, i))
	}
	return screens, nil
}

// widgets
func (w *linuxWebviewWindow) setEnabled(enabled bool) {
	C.gtk_widget_set_sensitive(w.gtkWidget(), C.gboolean(btoi(enabled)))
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func widgetSetVisible(widget pointer, hidden bool) {
	C.gtk_widget_set_visible((*C.GtkWidget)(widget), C.gboolean(btoi(!hidden)))
}

func (w *linuxWebviewWindow) close() {
	C.gtk_window_close(w.gtkWindow())
	getNativeApplication().unregisterWindow(windowPointer(w.window))
}

func (w *linuxWebviewWindow) enableDND() {
	winID := unsafe.Pointer(uintptr(w.parent.id))
	C.enableDND((*C.GtkWidget)(w.webview), C.gpointer(winID))
}

func (w *linuxWebviewWindow) disableDND() {
	winID := unsafe.Pointer(uintptr(w.parent.id))
	C.disableDND((*C.GtkWidget)(w.webview), C.gpointer(winID))
}

func (w *linuxWebviewWindow) execJS(js string) {
	InvokeAsync(func() {
		value := C.CString(js)
		defer C.free(unsafe.Pointer(value))
		// WebKitGTK 6.0 uses webkit_web_view_evaluate_javascript
		C.webkit_web_view_evaluate_javascript(w.webKitWebView(),
			value,
			C.gssize(len(js)),
			nil,
			nil,
			nil,
			nil,
			nil)
	})
}

// Preallocated buffer for drag-over JS calls
var dragOverJSBuffer = C.CString(strings.Repeat(" ", 64))
var emptyWorldName = C.CString("")

func (w *linuxWebviewWindow) execJSDragOver(x, y int) {
	buf := (*[64]byte)(unsafe.Pointer(dragOverJSBuffer))
	n := copy(buf[:], "window._wails.handleDragOver(")
	n += writeInt(buf[n:], x)
	buf[n] = ','
	n++
	n += writeInt(buf[n:], y)
	buf[n] = ')'
	n++
	buf[n] = 0

	C.webkit_web_view_evaluate_javascript(w.webKitWebView(),
		dragOverJSBuffer,
		C.gssize(n),
		nil,
		emptyWorldName,
		nil,
		nil,
		nil)
}

func writeInt(buf []byte, n int) int {
	if n < 0 {
		buf[0] = '-'
		return 1 + writeInt(buf[1:], -n)
	}
	if n == 0 {
		buf[0] = '0'
		return 1
	}
	tmp := n
	digits := 0
	for tmp > 0 {
		digits++
		tmp /= 10
	}
	for i := digits - 1; i >= 0; i-- {
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return digits
}

func getMousePosition() (int, int, *Screen) {
	// GTK4: Pointer position API is different
	// On Wayland, this may not work reliably
	display := C.gdk_display_get_default()
	seat := C.gdk_display_get_default_seat(display)
	device := C.gdk_seat_get_pointer(seat)
	_ = device
	// TODO: Implement GTK4 pointer position
	return 0, 0, nil
}

func (w *linuxWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
	if w.gtkmenu != nil {
		// GTK4: Different menu destruction
		w.gtkmenu = nil
	}
	C.gtk_window_destroy(w.gtkWindow())
}

func (w *linuxWebviewWindow) fullscreen() {
	C.gtk_window_fullscreen(w.gtkWindow())
}

func (w *linuxWebviewWindow) getCurrentMonitor() *C.GdkMonitor {
	display := C.gtk_widget_get_display(w.gtkWidget())
	surface := C.gtk_native_get_surface((*C.GtkNative)(unsafe.Pointer(w.gtkWindow())))
	if surface != nil {
		monitor := C.gdk_display_get_monitor_at_surface(display, surface)
		if monitor != nil {
			return monitor
		}
	}
	return nil
}

func (w *linuxWebviewWindow) getScreen() (*Screen, error) {
	monitor := w.getCurrentMonitor()
	if monitor == nil {
		return nil, fmt.Errorf("no monitor found")
	}
	name := C.gdk_monitor_get_model(monitor)
	var geometry C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &geometry)
	scaleFactor := int(C.gdk_monitor_get_scale_factor(monitor))
	return &Screen{
		ID:          fmt.Sprintf("%d", w.id),
		Name:        C.GoString(name),
		ScaleFactor: float32(scaleFactor),
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
		WorkArea: Rect{
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
		PhysicalWorkArea: Rect{
			X:      int(geometry.x),
			Y:      int(geometry.y),
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
		IsPrimary: false,
		Rotation:  0.0,
	}, nil
}

func (w *linuxWebviewWindow) getCurrentMonitorGeometry() (x int, y int, width int, height int, scaleFactor int) {
	monitor := w.getCurrentMonitor()
	if monitor == nil {
		return -1, -1, -1, -1, 1
	}
	var result C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &result)
	scaleFactor = int(C.gdk_monitor_get_scale_factor(monitor))
	return int(result.x), int(result.y), int(result.width), int(result.height), scaleFactor
}

func (w *linuxWebviewWindow) size() (int, int) {
	var width, height C.int
	C.gtk_window_get_default_size(w.gtkWindow(), &width, &height)
	if width <= 0 || height <= 0 {
		width = C.int(C.gtk_widget_get_width(w.gtkWidget()))
		height = C.int(C.gtk_widget_get_height(w.gtkWidget()))
	}
	return int(width), int(height)
}

func (w *linuxWebviewWindow) relativePosition() (int, int) {
	// GTK4/Wayland: Window positioning is not reliable
	// This is a documented limitation
	return 0, 0
}

func (w *linuxWebviewWindow) gtkWidget() *C.GtkWidget {
	return (*C.GtkWidget)(w.window)
}

func (w *linuxWebviewWindow) windowHide() {
	C.gtk_widget_set_visible(w.gtkWidget(), C.gboolean(0))
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	return C.gtk_window_is_fullscreen(w.gtkWindow()) != 0
}

func (w *linuxWebviewWindow) isFocused() bool {
	return C.gtk_window_is_active(w.gtkWindow()) != 0
}

func (w *linuxWebviewWindow) isMaximised() bool {
	return C.gtk_window_is_maximized(w.gtkWindow()) != 0 && !w.isFullscreen()
}

func (w *linuxWebviewWindow) isMinimised() bool {
	surface := C.gtk_native_get_surface((*C.GtkNative)(unsafe.Pointer(w.gtkWindow())))
	if surface == nil {
		return false
	}
	state := C.gdk_toplevel_get_state((*C.GdkToplevel)(unsafe.Pointer(surface)))
	return state&C.GDK_TOPLEVEL_STATE_MINIMIZED != 0
}

func (w *linuxWebviewWindow) isVisible() bool {
	return C.gtk_widget_is_visible(w.gtkWidget()) != 0
}

func (w *linuxWebviewWindow) maximise() {
	C.gtk_window_maximize(w.gtkWindow())
}

func (w *linuxWebviewWindow) minimise() {
	C.gtk_window_minimize(w.gtkWindow())
}

func windowNew(application pointer, menu pointer, windowId uint, gpuPolicy WebviewGpuPolicy) (window, webview, vbox pointer) {
	window = pointer(C.gtk_application_window_new((*C.GtkApplication)(application)))
	C.g_object_ref_sink(C.gpointer(window))

	C.attach_action_group_to_widget((*C.GtkWidget)(window))

	webview = windowNewWebview(windowId, gpuPolicy)
	vbox = pointer(C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0))
	name := C.CString("webview-box")
	defer C.free(unsafe.Pointer(name))
	C.gtk_widget_set_name((*C.GtkWidget)(vbox), name)

	C.gtk_window_set_child((*C.GtkWindow)(window), (*C.GtkWidget)(vbox))

	if menu != nil {
		menuBar := C.create_menu_bar_from_model((*C.GMenu)(menu))
		C.gtk_box_prepend((*C.GtkBox)(vbox), menuBar)
	}

	C.gtk_box_append((*C.GtkBox)(vbox), (*C.GtkWidget)(webview))
	C.gtk_widget_set_vexpand((*C.GtkWidget)(webview), C.gboolean(1))
	C.gtk_widget_set_hexpand((*C.GtkWidget)(webview), C.gboolean(1))
	return
}

func windowNewWebview(parentId uint, gpuPolicy WebviewGpuPolicy) pointer {
	c := NewCalloc()
	defer c.Free()
	manager := C.webkit_user_content_manager_new()
	// WebKitGTK 6.0: register_script_message_handler signature changed
	C.webkit_user_content_manager_register_script_message_handler(manager, c.String("external"), nil)

	// WebKitGTK 6.0: Create network session first
	networkSession := C.webkit_network_session_get_default()

	// Create web view with settings
	settings := C.webkit_settings_new()
	webView := C.webkit_web_view_new_with_user_content_manager(manager)

	C.save_webview_to_content_manager(unsafe.Pointer(manager), unsafe.Pointer(webView))
	C.save_window_id(unsafe.Pointer(webView), C.uint(parentId))
	C.save_window_id(unsafe.Pointer(manager), C.uint(parentId))

	// GPU policy
	switch gpuPolicy {
	case WebviewGpuPolicyNever:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER)
	case WebviewGpuPolicyAlways:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
	default:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND)
	}

	C.webkit_web_view_set_settings(C.webkit_web_view((*C.GtkWidget)(webView)), settings)

	// Register URI scheme handler
	registerURIScheme.Do(func() {
		webContext := C.webkit_web_view_get_context(C.webkit_web_view((*C.GtkWidget)(webView)))
		cScheme := C.CString(webview.Scheme)
		defer C.free(unsafe.Pointer(cScheme))
		C.webkit_web_context_register_uri_scheme(webContext, cScheme,
			(*[0]byte)(C.onProcessRequest), nil, nil)
	})

	_ = networkSession
	return pointer(webView)
}

func gtkBool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func (w *linuxWebviewWindow) gtkWindow() *C.GtkWindow {
	return (*C.GtkWindow)(w.window)
}

func (w *linuxWebviewWindow) webKitWebView() *C.WebKitWebView {
	return C.webkit_web_view((*C.GtkWidget)(w.webview))
}

func (w *linuxWebviewWindow) present() {
	C.gtk_window_present(w.gtkWindow())
}

func (w *linuxWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		cTitle := C.CString(title)
		C.gtk_window_set_title(w.gtkWindow(), cTitle)
		C.free(unsafe.Pointer(cTitle))
	}
}

func (w *linuxWebviewWindow) setSize(width, height int) {
	C.gtk_window_set_default_size(w.gtkWindow(), C.int(width), C.int(height))
}

func (w *linuxWebviewWindow) setDefaultSize(width int, height int) {
	C.gtk_window_set_default_size(w.gtkWindow(), C.int(width), C.int(height))
}

func windowSetGeometryHints(window pointer, minWidth, minHeight, maxWidth, maxHeight int) {
	w := (*C.GtkWidget)(window)
	if minWidth > 0 && minHeight > 0 {
		C.gtk_widget_set_size_request(w, C.int(minWidth), C.int(minHeight))
	}
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	C.gtk_window_set_resizable(w.gtkWindow(), gtkBool(resizable))
}

func (w *linuxWebviewWindow) move(x, y int) {
	// GTK4/Wayland: Window positioning is controlled by compositor
}

func (w *linuxWebviewWindow) position() (int, int) {
	// GTK4/Wayland: Cannot reliably get window position
	return 0, 0
}

func (w *linuxWebviewWindow) unfullscreen() {
	C.gtk_window_unfullscreen(w.gtkWindow())
	w.unmaximise()
}

func (w *linuxWebviewWindow) unmaximise() {
	C.gtk_window_unmaximize(w.gtkWindow())
}

func (w *linuxWebviewWindow) show() {
	C.gtk_widget_set_visible(w.gtkWidget(), gtkBool(true))
}

func (w *linuxWebviewWindow) windowShow() {
	if w.gtkWidget() == nil {
		return
	}
	C.gtk_widget_set_visible(w.gtkWidget(), gtkBool(true))
}

func (w *linuxWebviewWindow) hide() {
	C.gtk_widget_set_visible(w.gtkWidget(), gtkBool(false))
}

func (w *linuxWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	// GTK4: No direct equivalent - compositor-dependent
}

func (w *linuxWebviewWindow) setBorderless(borderless bool) {
	C.gtk_window_set_decorated(w.gtkWindow(), gtkBool(!borderless))
}

func (w *linuxWebviewWindow) setFrameless(frameless bool) {
	C.gtk_window_set_decorated(w.gtkWindow(), gtkBool(!frameless))
}

func (w *linuxWebviewWindow) setTransparent() {
	// GTK4: Transparency via CSS - different from GTK3
}

func (w *linuxWebviewWindow) setBackgroundColour(colour RGBA) {
	rgba := C.GdkRGBA{C.double(colour.Red) / 255.0, C.double(colour.Green) / 255.0, C.double(colour.Blue) / 255.0, C.double(colour.Alpha) / 255.0}
	C.webkit_web_view_set_background_color(w.webKitWebView(), &rgba)
}

func (w *linuxWebviewWindow) setIcon(icon pointer) {
	// GTK4: Window icons handled differently - no gtk_window_set_icon
}

func (w *linuxWebviewWindow) startDrag() error {
	C.beginWindowDrag(
		w.gtkWindow(),
		C.int(w.drag.MouseButton),
		C.double(w.drag.XRoot),
		C.double(w.drag.YRoot),
		C.guint32(w.drag.DragTime))
	return nil
}

func (w *linuxWebviewWindow) startResize(edge uint) error {
	C.beginWindowResize(
		w.gtkWindow(),
		C.GdkSurfaceEdge(edge),
		C.int(w.drag.MouseButton),
		C.double(w.drag.XRoot),
		C.double(w.drag.YRoot),
		C.guint32(w.drag.DragTime))
	return nil
}

func (w *linuxWebviewWindow) getZoom() float64 {
	return float64(C.webkit_web_view_get_zoom_level(w.webKitWebView()))
}

func (w *linuxWebviewWindow) setZoom(zoom float64) {
	if zoom < 1 {
		zoom = 1
	}
	C.webkit_web_view_set_zoom_level(w.webKitWebView(), C.gdouble(zoom))
}

func (w *linuxWebviewWindow) zoomIn() {
	w.setZoom(w.getZoom() * 1.10)
}

func (w *linuxWebviewWindow) zoomOut() {
	w.setZoom(w.getZoom() / 1.10)
}

func (w *linuxWebviewWindow) zoomReset() {
	w.setZoom(1.0)
}

func (w *linuxWebviewWindow) reload() {
	uri := C.CString("wails://")
	C.webkit_web_view_load_uri(w.webKitWebView(), uri)
	C.free(unsafe.Pointer(uri))
}

func (w *linuxWebviewWindow) setURL(uri string) {
	target := C.CString(uri)
	C.webkit_web_view_load_uri(w.webKitWebView(), target)
	C.free(unsafe.Pointer(target))
}

func (w *linuxWebviewWindow) setHTML(html string) {
	cHTML := C.CString(html)
	uri := C.CString("wails://")
	empty := C.CString("")
	defer C.free(unsafe.Pointer(cHTML))
	defer C.free(unsafe.Pointer(uri))
	defer C.free(unsafe.Pointer(empty))
	C.webkit_web_view_load_alternate_html(w.webKitWebView(), cHTML, uri, empty)
}

func (w *linuxWebviewWindow) flash(_ bool) {}

func (w *linuxWebviewWindow) ignoreMouse(ignore bool) {
	// GTK4: Input handling is different
}

func (w *linuxWebviewWindow) copy() {
	w.execJS("document.execCommand('copy')")
}

func (w *linuxWebviewWindow) cut() {
	w.execJS("document.execCommand('cut')")
}

func (w *linuxWebviewWindow) paste() {
	w.execJS("document.execCommand('paste')")
}

func (w *linuxWebviewWindow) delete() {
	w.execJS("document.execCommand('delete')")
}

func (w *linuxWebviewWindow) selectAll() {
	w.execJS("document.execCommand('selectAll')")
}

func (w *linuxWebviewWindow) undo() {
	w.execJS("document.execCommand('undo')")
}

func (w *linuxWebviewWindow) redo() {
	w.execJS("document.execCommand('redo')")
}

func (w *linuxWebviewWindow) setupSignalHandlers(emit func(e events.WindowEventType)) {
	c := NewCalloc()
	defer c.Free()

	winID := C.uintptr_t(w.parent.ID())

	C.setupWindowEventControllers(w.gtkWindow(), (*C.GtkWidget)(w.webview), winID)

	wv := unsafe.Pointer(w.webview)
	C.signal_connect(wv, c.String("load-changed"), C.handleLoadChanged, unsafe.Pointer(uintptr(winID)))

	contentManager := C.webkit_web_view_get_user_content_manager(w.webKitWebView())
	C.signal_connect(unsafe.Pointer(contentManager), c.String("script-message-received::external"), C.sendMessageToBackend, nil)
}

//export handleCloseRequest
func handleCloseRequest(window *C.GtkWindow, data C.uintptr_t) C.gboolean {
	processWindowEvent(C.uint(data), C.uint(events.Linux.WindowDeleteEvent))
	return C.gboolean(1)
}

//export handleNotifyState
func handleNotifyState(object *C.GObject, pspec *C.GParamSpec, data C.uintptr_t) {
	windowId := uint(data)
	window, ok := globalApplication.Window.GetByID(windowId)
	if !ok || window == nil {
		return
	}

	lw := getLinuxWebviewWindow(window)
	if lw == nil {
		return
	}

	if lw.isMaximised() {
		processWindowEvent(C.uint(data), C.uint(events.Linux.WindowDidResize))
	}
	if lw.isFullscreen() {
		processWindowEvent(C.uint(data), C.uint(events.Linux.WindowDidResize))
	}
}

//export handleFocusEnter
func handleFocusEnter(controller *C.GtkEventController, data C.uintptr_t) C.gboolean {
	processWindowEvent(C.uint(data), C.uint(events.Linux.WindowFocusIn))
	return C.gboolean(0)
}

//export handleFocusLeave
func handleFocusLeave(controller *C.GtkEventController, data C.uintptr_t) C.gboolean {
	processWindowEvent(C.uint(data), C.uint(events.Linux.WindowFocusOut))
	return C.gboolean(0)
}

//export handleButtonPressed
func handleButtonPressed(gesture *C.GtkGestureClick, nPress C.gint, x C.gdouble, y C.gdouble, data C.uintptr_t) {
	windowId := uint(data)
	window, ok := globalApplication.Window.GetByID(windowId)
	if !ok || window == nil {
		return
	}

	lw := getLinuxWebviewWindow(window)
	if lw == nil {
		return
	}

	button := C.gtk_gesture_single_get_current_button((*C.GtkGestureSingle)(unsafe.Pointer(gesture)))
	lw.drag.MouseButton = uint(button)
	lw.drag.XRoot = int(x)
	lw.drag.YRoot = int(y)
	lw.drag.DragTime = uint32(C.GDK_CURRENT_TIME)
}

//export handleButtonReleased
func handleButtonReleased(gesture *C.GtkGestureClick, nPress C.gint, x C.gdouble, y C.gdouble, data C.uintptr_t) {
	windowId := uint(data)
	window, ok := globalApplication.Window.GetByID(windowId)
	if !ok || window == nil {
		return
	}

	lw := getLinuxWebviewWindow(window)
	if lw == nil {
		return
	}

	button := C.gtk_gesture_single_get_current_button((*C.GtkGestureSingle)(unsafe.Pointer(gesture)))
	lw.endDrag(uint(button), int(x), int(y))
}

//export handleKeyPressed
func handleKeyPressed(controller *C.GtkEventControllerKey, keyval C.guint, keycode C.guint, state C.GdkModifierType, data C.uintptr_t) C.gboolean {
	windowID := uint(data)

	modifiers := uint(state)
	var acc accelerator

	if modifiers&C.GDK_SHIFT_MASK != 0 {
		acc.Modifiers = append(acc.Modifiers, ShiftKey)
	}
	if modifiers&C.GDK_CONTROL_MASK != 0 {
		acc.Modifiers = append(acc.Modifiers, ControlKey)
	}
	if modifiers&C.GDK_ALT_MASK != 0 {
		acc.Modifiers = append(acc.Modifiers, OptionOrAltKey)
	}
	if modifiers&C.GDK_SUPER_MASK != 0 {
		acc.Modifiers = append(acc.Modifiers, SuperKey)
	}

	keyString, ok := VirtualKeyCodes[uint(keyval)]
	if !ok {
		return C.gboolean(0)
	}
	acc.Key = keyString

	windowKeyEvents <- &windowKeyEvent{
		windowId:          windowID,
		acceleratorString: acc.String(),
	}

	return C.gboolean(0)
}

//export onDropEnter
func onDropEnter(data C.uintptr_t) {
	windowId := uint(data)
	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return
	}
	if w, ok := targetWindow.(*WebviewWindow); ok {
		w.HandleDragEnter()
	}
}

//export onDropLeave
func onDropLeave(data C.uintptr_t) {
	windowId := uint(data)
	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return
	}
	if w, ok := targetWindow.(*WebviewWindow); ok {
		w.HandleDragLeave()
	}
}

//export onDropMotion
func onDropMotion(x C.gint, y C.gint, data C.uintptr_t) {
	windowId := uint(data)
	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return
	}
	if w, ok := targetWindow.(*WebviewWindow); ok {
		w.HandleDragOver(int(x), int(y))
	}
}

//export onDropFiles
func onDropFiles(paths **C.char, x C.gint, y C.gint, data C.uintptr_t) {
	windowId := uint(data)
	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return
	}

	offset := unsafe.Sizeof(uintptr(0))
	var filenames []string
	for *paths != nil {
		filenames = append(filenames, C.GoString(*paths))
		paths = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(paths)) + offset))
	}

	targetWindow.InitiateFrontendDropProcessing(filenames, int(x), int(y))
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &windowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

//export onProcessRequest
func onProcessRequest(request *C.WebKitURISchemeRequest, data C.uintptr_t) {
	webView := C.webkit_uri_scheme_request_get_web_view(request)
	windowId := uint(C.get_window_id(unsafe.Pointer(webView)))
	webviewRequests <- &webViewAssetRequest{
		Request:  webview.NewRequest(unsafe.Pointer(request)),
		windowId: windowId,
		windowName: func() string {
			if window, ok := globalApplication.Window.GetByID(windowId); ok {
				return window.Name()
			}
			return ""
		}(),
	}
}

//export sendMessageToBackend
func sendMessageToBackend(contentManager *C.WebKitUserContentManager, result *C.WebKitJavascriptResult,
	data unsafe.Pointer) {

	// Get the windowID from the contentManager
	thisWindowID := uint(C.get_window_id(unsafe.Pointer(contentManager)))

	webView := C.get_webview_from_content_manager(unsafe.Pointer(contentManager))
	var origin string
	if webView != nil {
		currentUri := C.webkit_web_view_get_uri(webView)
		if currentUri != nil {
			uri := C.g_strdup(currentUri)
			defer C.g_free(C.gpointer(uri))
			origin = C.GoString(uri)
		}
	}

	var msg string
	value := C.webkit_javascript_result_get_js_value(result)
	message := C.jsc_value_to_string(value)
	msg = C.GoString(message)
	defer C.g_free(C.gpointer(message))
	windowMessageBuffer <- &windowMessage{
		windowId: thisWindowID,
		message:  msg,
		originInfo: &OriginInfo{
			Origin: origin,
		},
	}
}

var _ = time.Now
var _ = events.Linux
var _ = webview.Scheme
