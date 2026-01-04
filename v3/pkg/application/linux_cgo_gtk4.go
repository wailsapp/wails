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

#include "linux_cgo_gtk4.h"
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
	C.set_app_menu_model(gmenu)
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

// Keyboard accelerator support for GTK4 menus

// namedKeysToGTK maps Wails key names to GDK keysym values
// These are X11 keysym values that GDK uses
var namedKeysToGTK = map[string]C.guint{
	"backspace": C.guint(0xff08),
	"tab":       C.guint(0xff09),
	"return":    C.guint(0xff0d),
	"enter":     C.guint(0xff0d),
	"escape":    C.guint(0xff1b),
	"left":      C.guint(0xff51),
	"right":     C.guint(0xff53),
	"up":        C.guint(0xff52),
	"down":      C.guint(0xff54),
	"space":     C.guint(0xff80),
	"delete":    C.guint(0xff9f),
	"home":      C.guint(0xff95),
	"end":       C.guint(0xff9c),
	"page up":   C.guint(0xff9a),
	"page down": C.guint(0xff9b),
	"f1":        C.guint(0xffbe),
	"f2":        C.guint(0xffbf),
	"f3":        C.guint(0xffc0),
	"f4":        C.guint(0xffc1),
	"f5":        C.guint(0xffc2),
	"f6":        C.guint(0xffc3),
	"f7":        C.guint(0xffc4),
	"f8":        C.guint(0xffc5),
	"f9":        C.guint(0xffc6),
	"f10":       C.guint(0xffc7),
	"f11":       C.guint(0xffc8),
	"f12":       C.guint(0xffc9),
	"f13":       C.guint(0xffca),
	"f14":       C.guint(0xffcb),
	"f15":       C.guint(0xffcc),
	"f16":       C.guint(0xffcd),
	"f17":       C.guint(0xffce),
	"f18":       C.guint(0xffcf),
	"f19":       C.guint(0xffd0),
	"f20":       C.guint(0xffd1),
	"f21":       C.guint(0xffd2),
	"f22":       C.guint(0xffd3),
	"f23":       C.guint(0xffd4),
	"f24":       C.guint(0xffd5),
	"f25":       C.guint(0xffd6),
	"f26":       C.guint(0xffd7),
	"f27":       C.guint(0xffd8),
	"f28":       C.guint(0xffd9),
	"f29":       C.guint(0xffda),
	"f30":       C.guint(0xffdb),
	"f31":       C.guint(0xffdc),
	"f32":       C.guint(0xffdd),
	"f33":       C.guint(0xffde),
	"f34":       C.guint(0xffdf),
	"f35":       C.guint(0xffe0),
	"numlock":   C.guint(0xff7f),
}

// parseKeyGTK converts a Wails key string to a GDK keysym value
func parseKeyGTK(key string) C.guint {
	// Check named keys first
	if result, found := namedKeysToGTK[key]; found {
		return result
	}
	// For single character keys, convert using gdk_unicode_to_keyval
	if len(key) != 1 {
		return C.guint(0)
	}
	keyval := rune(key[0])
	return C.gdk_unicode_to_keyval(C.guint(keyval))
}

// parseModifiersGTK converts Wails modifiers to GDK modifier type
func parseModifiersGTK(modifiers []modifier) C.GdkModifierType {
	var result C.GdkModifierType

	for _, mod := range modifiers {
		switch mod {
		case ShiftKey:
			result |= C.GDK_SHIFT_MASK
		case ControlKey, CmdOrCtrlKey:
			result |= C.GDK_CONTROL_MASK
		case OptionOrAltKey:
			result |= C.GDK_ALT_MASK
		case SuperKey:
			result |= C.GDK_SUPER_MASK
		}
	}
	return result
}

// acceleratorToGTK converts a Wails accelerator to GTK key/modifiers
func acceleratorToGTK(accel *accelerator) (C.guint, C.GdkModifierType) {
	key := parseKeyGTK(accel.Key)
	mods := parseModifiersGTK(accel.Modifiers)
	return key, mods
}

// setMenuItemAccelerator sets the keyboard accelerator for a menu item
// This uses gtk_application_set_accels_for_action to register the shortcut
func setMenuItemAccelerator(itemId uint, accel *accelerator) {
	if accel == nil {
		return
	}

	// Look up the action name for this menu item
	actionName, ok := menuItemActions[itemId]
	if !ok {
		return
	}

	// Get the GtkApplication pointer
	app := getNativeApplication()
	if app == nil || app.application == nil {
		return
	}

	// Convert accelerator to GTK format
	key, mods := acceleratorToGTK(accel)
	if key == 0 {
		return
	}

	// Build accelerator string using GTK's function
	accelString := C.build_accelerator_string(key, mods)
	if accelString == nil {
		return
	}
	defer C.g_free(C.gpointer(accelString))

	// Set the accelerator on the application
	cActionName := C.CString(actionName)
	defer C.free(unsafe.Pointer(cActionName))
	C.set_action_accelerator((*C.GtkApplication)(app.application), cActionName, accelString)
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
	// WebKitGTK 6.0: webkit_web_view_new_with_user_content_manager() was removed
	// Use create_webview_with_user_content_manager() helper instead
	webView := C.create_webview_with_user_content_manager(manager)

	C.save_webview_to_content_manager(unsafe.Pointer(manager), unsafe.Pointer(webView))
	C.save_window_id(unsafe.Pointer(webView), C.uint(parentId))
	C.save_window_id(unsafe.Pointer(manager), C.uint(parentId))

	// GPU policy
	// WebKitGTK 6.0: WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND was removed
	// Only ALWAYS and NEVER are available
	switch gpuPolicy {
	case WebviewGpuPolicyNever:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER)
	case WebviewGpuPolicyAlways:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
	default:
		// Default to ALWAYS (was ON_DEMAND in older WebKitGTK)
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
	}

	C.webkit_web_view_set_settings(C.webkit_web_view((*C.GtkWidget)(webView)), settings)

	// Register URI scheme handler
	registerURIScheme.Do(func() {
		webContext := C.webkit_web_view_get_context(C.webkit_web_view((*C.GtkWidget)(webView)))
		cScheme := C.CString("wails")
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

func (w *linuxWebviewWindow) windowShow() {
	if w.gtkWidget() == nil {
		return
	}
	C.gtk_widget_set_visible(w.gtkWidget(), gtkBool(true))
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
	rgba := C.GdkRGBA{C.float(colour.Red) / 255.0, C.float(colour.Green) / 255.0, C.float(colour.Blue) / 255.0, C.float(colour.Alpha) / 255.0}
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

// startResize is handled by webview_window_linux.go
// GTK4-specific resize using beginWindowResize can be added via a helper function

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

//export handleLoadChanged
func handleLoadChanged(wv *C.WebKitWebView, event C.WebKitLoadEvent, data C.uintptr_t) {
	switch event {
	case C.WEBKIT_LOAD_STARTED:
		processWindowEvent(C.uint(data), C.uint(events.Linux.WindowLoadStarted))
	case C.WEBKIT_LOAD_REDIRECTED:
		processWindowEvent(C.uint(data), C.uint(events.Linux.WindowLoadRedirected))
	case C.WEBKIT_LOAD_COMMITTED:
		processWindowEvent(C.uint(data), C.uint(events.Linux.WindowLoadCommitted))
	case C.WEBKIT_LOAD_FINISHED:
		processWindowEvent(C.uint(data), C.uint(events.Linux.WindowLoadFinished))
	}
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

// WebKitGTK 6.0: callback now receives JSCValue directly instead of WebKitJavascriptResult
//
//export sendMessageToBackend
func sendMessageToBackend(contentManager *C.WebKitUserContentManager, value *C.JSCValue,
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

	// WebKitGTK 6.0: JSCValue is passed directly, no need for webkit_javascript_result_get_js_value
	message := C.jsc_value_to_string(value)
	msg := C.GoString(message)
	defer C.g_free(C.gpointer(message))
	windowMessageBuffer <- &windowMessage{
		windowId: thisWindowID,
		message:  msg,
		originInfo: &OriginInfo{
			Origin: origin,
		},
	}
}

// ============================================================================
// GTK4 Dialog System - Go wrapper functions
// ============================================================================

// Dialog request tracking
var (
	dialogRequestCounter uint32
	dialogRequestMutex   sync.Mutex
	fileDialogCallbacks  = make(map[uint]chan string)
	alertDialogCallbacks = make(map[uint]chan int)
)

func nextDialogRequestID() uint {
	dialogRequestMutex.Lock()
	defer dialogRequestMutex.Unlock()
	dialogRequestCounter++
	return uint(dialogRequestCounter)
}

//export fileDialogCallback
func fileDialogCallback(requestID C.uint, files **C.char, count C.int, cancelled C.gboolean) {
	dialogRequestMutex.Lock()
	ch, ok := fileDialogCallbacks[uint(requestID)]
	if ok {
		delete(fileDialogCallbacks, uint(requestID))
	}
	dialogRequestMutex.Unlock()

	if !ok {
		return
	}

	if cancelled != 0 {
		close(ch)
		return
	}

	// Convert C string array to Go strings
	if count > 0 && files != nil {
		slice := unsafe.Slice(files, int(count))
		for _, cstr := range slice {
			if cstr != nil {
				ch <- C.GoString(cstr)
			}
		}
	}
	close(ch)
}

//export alertDialogCallback
func alertDialogCallback(requestID C.uint, buttonIndex C.int) {
	dialogRequestMutex.Lock()
	ch, ok := alertDialogCallbacks[uint(requestID)]
	if ok {
		delete(alertDialogCallbacks, uint(requestID))
	}
	dialogRequestMutex.Unlock()

	if !ok {
		return
	}

	ch <- int(buttonIndex)
	close(ch)
}

func runChooserDialog(window pointer, allowMultiple, createFolders, showHidden bool, currentFolder, title string, action int, acceptLabel string, filters []FileFilter) (chan string, error) {
	requestID := nextDialogRequestID()
	resultChan := make(chan string, 100)

	dialogRequestMutex.Lock()
	fileDialogCallbacks[requestID] = resultChan
	dialogRequestMutex.Unlock()

	InvokeAsync(func() {
		cTitle := C.CString(title)
		defer C.free(unsafe.Pointer(cTitle))

		dialog := C.create_file_dialog(cTitle)
		defer C.g_object_unref(C.gpointer(dialog))

		// Create filter list if we have filters
		if len(filters) > 0 {
			filterStore := C.g_list_store_new(C.gtk_file_filter_get_type())
			defer C.g_object_unref(C.gpointer(filterStore))

			for _, filter := range filters {
				cName := C.CString(filter.DisplayName)
				cPattern := C.CString(filter.Pattern)
				C.add_file_filter(dialog, filterStore, cName, cPattern)
				C.free(unsafe.Pointer(cName))
				C.free(unsafe.Pointer(cPattern))
			}
			C.set_file_dialog_filters(dialog, filterStore)
		}

		// Set initial folder if provided
		if currentFolder != "" {
			cFolder := C.CString(currentFolder)
			file := C.g_file_new_for_path(cFolder)
			C.gtk_file_dialog_set_initial_folder(dialog, file)
			C.g_object_unref(C.gpointer(file))
			C.free(unsafe.Pointer(cFolder))
		}

		var parent *C.GtkWindow
		if window != nil {
			parent = (*C.GtkWindow)(window)
		}

		// Determine dialog type based on action
		// GTK_FILE_CHOOSER_ACTION_OPEN = 0
		// GTK_FILE_CHOOSER_ACTION_SAVE = 1
		// GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER = 2
		isFolder := action == 2
		isSave := action == 1

		if isSave {
			C.show_save_file_dialog(parent, dialog, C.uint(requestID))
		} else {
			C.show_open_file_dialog(parent, dialog, C.uint(requestID), gtkBool(allowMultiple), gtkBool(isFolder))
		}
	})

	return resultChan, nil
}

func runOpenFileDialog(dialog *OpenFileDialogStruct) (chan string, error) {
	var action int

	if dialog.canChooseDirectories {
		action = 2 // GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER
	} else {
		action = 0 // GTK_FILE_CHOOSER_ACTION_OPEN
	}

	window := nilPointer
	if dialog.window != nil {
		nativeWindow := dialog.window.NativeWindow()
		if nativeWindow != nil {
			window = pointer(nativeWindow)
		}
	}

	buttonText := dialog.buttonText
	if buttonText == "" {
		buttonText = "_Open"
	}

	return runChooserDialog(
		window,
		dialog.allowsMultipleSelection,
		false, // createFolders not applicable for open
		dialog.showHiddenFiles,
		dialog.directory,
		dialog.title,
		action,
		buttonText,
		dialog.filters,
	)
}

func runSaveFileDialog(dialog *SaveFileDialogStruct) (chan string, error) {
	window := nilPointer
	if dialog.window != nil {
		nativeWindow := dialog.window.NativeWindow()
		if nativeWindow != nil {
			window = pointer(nativeWindow)
		}
	}

	buttonText := dialog.buttonText
	if buttonText == "" {
		buttonText = "_Save"
	}

	return runChooserDialog(
		window,
		false,
		dialog.canCreateDirectories,
		dialog.showHiddenFiles,
		dialog.directory,
		dialog.title,
		1, // GTK_FILE_CHOOSER_ACTION_SAVE
		buttonText,
		dialog.filters,
	)
}

func runQuestionDialog(parent pointer, options *MessageDialog) int {
	requestID := nextDialogRequestID()
	resultChan := make(chan int, 1)

	dialogRequestMutex.Lock()
	alertDialogCallbacks[requestID] = resultChan
	dialogRequestMutex.Unlock()

	InvokeAsync(func() {
		cMessage := C.CString(options.Message)
		defer C.free(unsafe.Pointer(cMessage))

		var cDetail *C.char
		if options.Message != "" {
			cDetail = C.CString(options.Message)
			defer C.free(unsafe.Pointer(cDetail))
		}

		// Build button labels
		buttonLabels := make([]*C.char, len(options.Buttons)+1)
		for i, btn := range options.Buttons {
			buttonLabels[i] = C.CString(btn.Label)
		}
		buttonLabels[len(options.Buttons)] = nil // NULL terminator

		defer func() {
			for _, label := range buttonLabels[:len(options.Buttons)] {
				C.free(unsafe.Pointer(label))
			}
		}()

		// Find default and cancel button indices
		defaultButton := 0
		cancelButton := -1
		for i, btn := range options.Buttons {
			if btn.IsDefault {
				defaultButton = i
			}
			if btn.IsCancel {
				cancelButton = i
			}
		}

		var parentWindow *C.GtkWindow
		if parent != nil {
			parentWindow = (*C.GtkWindow)(parent)
		}

		C.show_alert_dialog(
			parentWindow,
			cMessage,
			cDetail,
			(**C.char)(unsafe.Pointer(&buttonLabels[0])),
			C.int(len(options.Buttons)),
			C.int(defaultButton),
			C.int(cancelButton),
			C.uint(requestID),
		)
	})

	// Wait for result
	result := <-resultChan
	return result
}

func getPrimaryScreen() (*Screen, error) {
	display := C.gdk_display_get_default()
	monitors := C.gdk_display_get_monitors(display)
	if monitors == nil {
		return nil, fmt.Errorf("no monitors found")
	}
	count := C.g_list_model_get_n_items(monitors)
	if count == 0 {
		return nil, fmt.Errorf("no monitors found")
	}
	monitor := (*C.GdkMonitor)(C.g_list_model_get_item(monitors, 0))
	if monitor == nil {
		return nil, fmt.Errorf("failed to get primary monitor")
	}
	defer C.g_object_unref(C.gpointer(monitor))

	var geometry C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &geometry)
	scaleFactor := int(C.gdk_monitor_get_scale_factor(monitor))
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
		Rotation: 0.0,
	}, nil
}

func openDevTools(wv pointer) {
	inspector := C.webkit_web_view_get_inspector((*C.WebKitWebView)(wv))
	C.webkit_web_inspector_show(inspector)
}

func enableDevTools(wv pointer) {
	settings := C.webkit_web_view_get_settings((*C.WebKitWebView)(wv))
	enabled := C.webkit_settings_get_enable_developer_extras(settings)
	if enabled == 0 {
		C.webkit_settings_set_enable_developer_extras(settings, C.gboolean(1))
	} else {
		C.webkit_settings_set_enable_developer_extras(settings, C.gboolean(0))
	}
}

var _ = time.Now
var _ = events.Linux
var _ = strings.TrimSpace
