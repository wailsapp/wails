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
extern void handleLoadChanged(WebKitWebView*, WebKitLoadEvent, uintptr_t);
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

// GTK4 drag-and-drop uses GtkDropTarget instead of GTK3's drag signals
// This is a significant API change - stub for now

static void enableDND(GtkWidget *widget, gpointer data) {
    // TODO: Implement GTK4 drag-and-drop with GtkDropTarget
}

static void disableDND(GtkWidget *widget, gpointer data) {
    // TODO: Implement GTK4 drag-and-drop blocking
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
func menuAddSeparator(menu *Menu) {
	// GTK4: GMenu separators are sections, not items
	// TODO: Implement GTK4 menu separators
}

func menuAppend(parent *Menu, menu *MenuItem) {
	// GTK4: Use g_menu_append_item
	// TODO: Implement GTK4 menu append
}

func menuBarNew() pointer {
	// GTK4: Use GtkPopoverMenuBar
	return nil
}

func menuNew() pointer {
	// GTK4: Use GMenu
	return pointer(C.g_menu_new())
}

func menuSetSubmenu(item *MenuItem, menu *Menu) {
	// GTK4: Use g_menu_item_set_submenu
	// TODO: Implement GTK4 submenu
}

func menuGetRadioGroup(item *linuxMenuItem) *GSList {
	// GTK4: Radio groups work via GAction state
	return nil
}

//export handleClick
func handleClick(idPtr unsafe.Pointer) {
	// GTK4: GAction activation callback
	// TODO: Implement GTK4 menu click handling
}

func attachMenuHandler(item *MenuItem) uint {
	// GTK4: Use g_signal_connect on GSimpleAction
	// TODO: Implement GTK4 menu handler
	return 0
}

// menuItem - GTK4 uses GMenuItem
func menuItemChecked(widget pointer) bool {
	// GTK4: Check GAction state
	return false
}

func menuItemNew(label string, bitmap []byte) pointer {
	// GTK4: Use g_menu_item_new
	return nil
}

func menuItemDestroy(widget pointer) {
	// GTK4: GMenuItem is reference counted
}

func menuItemAddProperties(menuItem *C.GtkWidget, label string, bitmap []byte) pointer {
	// GTK4: Different API for menu items
	return nil
}

func menuCheckItemNew(label string, bitmap []byte) pointer {
	// GTK4: Use GMenuItem with stateful GAction
	return nil
}

func menuItemSetChecked(widget pointer, checked bool) {
	// GTK4: Set GAction state
}

func menuItemSetDisabled(widget pointer, disabled bool) {
	// GTK4: Set GAction enabled state
}

func menuItemSetLabel(widget pointer, label string) {
	// GTK4: Use g_menu_item_set_label
}

func menuItemRemoveBitmap(widget pointer) {
	// GTK4: Different icon handling
}

func menuItemSetBitmap(widget pointer, bitmap []byte) {
	// GTK4: Use g_menu_item_set_icon
}

func menuItemSetToolTip(widget pointer, tooltip string) {
	// GTK4: Tooltips on menu items
}

func menuItemSignalBlock(widget pointer, handlerId uint, block bool) {
	if block {
		C.g_signal_handler_block(C.gpointer(widget), C.ulong(handlerId))
	} else {
		C.g_signal_handler_unblock(C.gpointer(widget), C.ulong(handlerId))
	}
}

func menuRadioItemNew(group *GSList, label string) pointer {
	// GTK4: Use GMenuItem with radio action
	return nil
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
	return C.gtk_window_get_default_size(w.gtkWindow(), nil, nil), 0
	// GTK4: gtk_window_get_size is deprecated, use gtk_window_get_default_size
	// TODO: Fix this to return proper values
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
	// GTK4: There's no direct API for this on Wayland
	// The window state is managed by the compositor
	return false
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
	webview = windowNewWebview(windowId, gpuPolicy)
	vbox = pointer(C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0))
	name := C.CString("webview-box")
	defer C.free(unsafe.Pointer(name))
	C.gtk_widget_set_name((*C.GtkWidget)(vbox), name)

	// GTK4: Use gtk_window_set_child instead of gtk_container_add
	C.gtk_window_set_child((*C.GtkWindow)(window), (*C.GtkWidget)(vbox))
	if menu != nil {
		C.gtk_box_prepend((*C.GtkBox)(vbox), (*C.GtkWidget)(menu))
	}
	C.gtk_box_append((*C.GtkBox)(vbox), (*C.GtkWidget)(webview))
	// GTK4: Set expand properties for webview
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
	// GTK4: GdkGeometry and gtk_window_set_geometry_hints are removed
	// Size constraints must be handled differently in GTK4
	// TODO: Implement via GtkConstraint or size request
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
	// TODO: GTK4 drag via GtkGestureDrag
	return nil
}

func (w *linuxWebviewWindow) startResize(edge uint) error {
	// TODO: GTK4 resize via gtk_window_begin_resize
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

	winID := unsafe.Pointer(uintptr(C.uint(w.parent.ID())))

	wv := unsafe.Pointer(w.webview)
	C.signal_connect(wv, c.String("load-changed"), C.handleLoadChanged, winID)

	contentManager := C.webkit_web_view_get_user_content_manager(w.webKitWebView())
	C.signal_connect(unsafe.Pointer(contentManager), c.String("script-message-received::external"), C.sendMessageToBackend, nil)
}

var _ = time.Now
var _ = events.Linux
var _ = webview.Scheme
