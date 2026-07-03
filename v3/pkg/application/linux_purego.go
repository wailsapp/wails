//go:build linux && purego && !gtk3 && !android && !server

package application

// CGO-free port of linux_cgo.go: the GTK4/WebKitGTK-6.0 backend driven through
// purego instead of cgo. The function surface (names, signatures, behaviour)
// mirrors the cgo shim exactly so the shared Linux files compile unchanged
// against either backend.
//
// Known cgo bugs are FIXED here rather than ported 1:1 — each fix is recorded
// in BUGS_FOUND.md.

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

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
	// BUGS_FOUND #4: the cgo backend accesses this map from the GTK main
	// thread (menuActionActivated) while attachMenuHandler writes it from the
	// menu-processing goroutine, with no synchronisation. Guarded here.
	gtkSignalToMenuItem     = map[uint]*MenuItem{}
	gtkSignalToMenuItemLock sync.RWMutex

	mainThreadId uintptr
)

var (
	registerURIScheme sync.Once
	fixSignalHandlers sync.Once
)

func init() {
	linuxLibsErr = loadLinuxLibraries()
	if linuxLibsErr == nil {
		// Package init runs on the process's main thread, the same thread
		// that will run the GTK main loop (matching the cgo init()).
		mainThreadId = g_thread_self()
	}
}

func isOnMainThread() bool {
	return g_thread_self() == mainThreadId
}

// implementation below

func appName() string {
	// BUGS_FOUND #1: g_get_application_name returns a string owned by GLib;
	// the cgo backend free()s it (undefined behaviour). Copy, don't free.
	return goString(g_get_application_name())
}

func appNew(name string) pointer {
	if linuxLibsErr != nil {
		// The GUI libraries are dlopen'ed at runtime; fail with the full
		// actionable report (which libraries/symbols, what to install)
		// instead of a nil-pointer crash on the first GTK call.
		Fatal("%v", linuxLibsErr)
	}

	installSignalHandlers()

	appId := fmt.Sprintf("org.wails.%s", name)
	return pointer(gtk_application_new(appId, gApplicationDefaultFlags))
}

func setProgramName(prgName string) {
	g_set_prgname(prgName)
}

func appRun(app pointer) error {
	application := uintptr(app)
	g_application_hold(application)

	signalConnect(application, "activate", activateLinuxPtr, 0)
	status := g_application_run(application, 0, 0)
	// The GTK main loop has stopped. Tell the asset-server webview layer to stop
	// marshalling WebKit calls onto it, so any request still being completed on a
	// worker goroutine runs inline instead of blocking on a loop that is gone.
	// See #5631.
	webview.DisableMainThreadDispatch()
	g_application_release(application)
	g_object_unref(application)

	var err error
	if status != 0 {
		err = fmt.Errorf("exit code: %d", status)
	}
	return err
}

func appDestroy(application pointer) {
	g_application_quit(uintptr(application))
}

func (w *linuxWebviewWindow) contextMenuSignals(menu pointer) {
	// GTK4: GtkPopoverMenu items are wired through the "app" GAction group,
	// which is attached to the window in windowNew. The popover's "closed"
	// signal (used for cleanup) is connected in showContextMenu, so there
	// is nothing to wire up here.
}

func (w *linuxWebviewWindow) contextMenuShow(menu pointer, data *ContextMenuData) {
	// GTK4: present the GMenu model as a GtkPopoverMenu anchored to the
	// webview at the click coordinates (which are relative to the webview).
	showContextMenu(uintptr(w.webview), uintptr(menu), data.X, data.Y)
}

func (a *linuxApp) getCurrentWindowID() uint {
	window := gtk_application_get_active_window(uintptr(a.application))
	if window == 0 {
		return uint(1)
	}
	a.windowMapLock.Lock()
	identifier, ok := a.windowMap[windowPointer(window)]
	a.windowMapLock.Unlock()
	if ok {
		return identifier
	}
	return uint(1)
}

func (a *linuxApp) getWindows() []pointer {
	result := []pointer{}
	// BUGS_FOUND #2: the cgo backend dereferences the returned GList head
	// unconditionally; when no window exists the list is NULL and it crashes.
	windows := gtk_application_get_windows(uintptr(a.application))
	for windows != nil {
		result = append(result, windows.data)
		windows = windows.next
	}
	return result
}

func (a *linuxApp) hideAllWindows() {
	for _, window := range a.getWindows() {
		gtk_widget_set_visible(uintptr(window), 0)
	}
}

func (a *linuxApp) showAllWindows() {
	for _, window := range a.getWindows() {
		gtk_window_present(uintptr(window))
	}
}

func (a *linuxApp) setIcon(icon []byte) {
	// GTK4 removed per-window icon APIs. The application icon is determined by
	// the .desktop file's Icon= field at the desktop-integration level.
	// No programmatic equivalent exists for setting icons from bytes in GTK4.
}

func clipboardGet() string {
	return clipboardGetTextSync()
}

func clipboardSet(text string) {
	display := gdk_display_get_default()
	clip := gdk_display_get_clipboard(display)
	gdk_clipboard_set_text(clip, text)
}

// Menu - GTK4 uses GMenu/GAction instead of GtkMenu

// BUGS_FOUND #5: in the cgo backend menuItemActions and menuItemCounters are
// plain maps written during menu construction and read from GTK callbacks
// without any locking (only menuItemIds had a mutex). One mutex guards all
// three here.
var (
	menuItemActionCounter uint32
	menuItemActions       = make(map[uint]string)
	menuItemIds           = make(map[pointer]uint)
	menuItemCounters      = make(map[pointer]int)
	menuItemsLock         sync.RWMutex
)

func generateActionName(itemId uint) string {
	menuItemsLock.Lock()
	defer menuItemsLock.Unlock()
	menuItemActionCounter++
	name := fmt.Sprintf("action_%d", menuItemActionCounter)
	menuItemActions[itemId] = name
	return name
}

func lookupActionName(itemId uint) (string, bool) {
	menuItemsLock.RLock()
	defer menuItemsLock.RUnlock()
	name, ok := menuItemActions[itemId]
	return name, ok
}

func menuActionActivated(id uint) {
	gtkSignalToMenuItemLock.RLock()
	item, ok := gtkSignalToMenuItem[id]
	gtkSignalToMenuItemLock.RUnlock()
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

func menuNewSection() pointer {
	return pointer(g_menu_new())
}

func menuAppendSection(menu *Menu, section pointer) {
	if menu.impl == nil {
		return
	}
	impl := menu.impl.(*linuxMenu)
	if impl.native == nilPointer {
		return
	}
	g_menu_append_section(uintptr(impl.native), 0, uintptr(section))
}

func menuAppendItemToSection(section pointer, item *MenuItem) {
	if item.impl == nil {
		return
	}
	menuImpl := item.impl.(*linuxMenuItem)
	if menuImpl.native == nilPointer {
		return
	}

	menuImpl.parentMenu = section
	menuImpl.isHidden = item.hidden

	if !item.hidden {
		g_menu_append_item(uintptr(section), uintptr(menuImpl.native))
	}
}

func menuAppend(parent *Menu, menu *MenuItem, hidden bool) {
	if parent.impl == nil || menu.impl == nil {
		return
	}
	parentImpl := parent.impl.(*linuxMenu)
	menuImpl := menu.impl.(*linuxMenuItem)
	if parentImpl.native == nilPointer || menuImpl.native == nilPointer {
		return
	}

	menuImpl.parentMenu = parentImpl.native
	menuImpl.isHidden = hidden

	menuItemsLock.Lock()
	menuImpl.menuIndex = menuItemCounters[parentImpl.native]
	menuItemCounters[parentImpl.native]++
	menuItemsLock.Unlock()

	if !hidden {
		g_menu_append_item(uintptr(parentImpl.native), uintptr(menuImpl.native))
	}
}

// menuClear removes every item from the menu's native GMenu so that it can be
// rebuilt from scratch on Menu.Update() (#5464). The per-menu append counter is
// reset too, so rebuilt items get fresh 0-based positions (menuIndex is used by
// menu_remove_item for hide/show). This mirrors the GTK3/cgo menuClear.
func menuClear(menu *Menu) {
	if menu.impl == nil {
		return
	}
	impl := menu.impl.(*linuxMenu)
	if impl.native == nilPointer {
		return
	}
	g_menu_remove_all(uintptr(impl.native))
	menuItemsLock.Lock()
	delete(menuItemCounters, impl.native)
	menuItemsLock.Unlock()
}

func menuBarNew() pointer {
	gmenu := g_menu_new()
	appMenuModel = gmenu
	return pointer(gmenu)
}

func menuNew() pointer {
	return pointer(g_menu_new())
}

func menuSetSubmenu(item *MenuItem, menu *Menu) {
	if item.impl == nil || menu.impl == nil {
		return
	}
	itemImpl := item.impl.(*linuxMenuItem)
	menuImpl := menu.impl.(*linuxMenu)
	if itemImpl.native == nilPointer || menuImpl.native == nilPointer {
		return
	}
	g_menu_item_set_submenu(uintptr(itemImpl.native), uintptr(menuImpl.native))
}

func menuGetRadioGroup(item *linuxMenuItem) *GSList {
	return nil
}

func attachMenuHandler(item *MenuItem) uint {
	gtkSignalToMenuItemLock.Lock()
	gtkSignalToMenuItem[item.id] = item
	gtkSignalToMenuItemLock.Unlock()
	return item.id
}

func menuItemChecked(widget pointer) bool {
	if widget == nilPointer {
		return false
	}
	menuItemsLock.RLock()
	itemId, exists := menuItemIds[widget]
	menuItemsLock.RUnlock()
	if !exists {
		return false
	}
	actionName, ok := lookupActionName(itemId)
	if !ok {
		return false
	}
	return getActionState(actionName)
}

func menuItemNew(label string, bitmap []byte) pointer {
	return nilPointer
}

func menuItemNewWithId(label string, bitmap []byte, itemId uint) pointer {
	actionName := generateActionName(itemId)
	gitem := createMenuItem(label, actionName, itemId)

	menuItemsLock.Lock()
	menuItemIds[pointer(gitem)] = itemId
	menuItemsLock.Unlock()
	return pointer(gitem)
}

func menuItemDestroy(widget pointer) {
	if widget != nilPointer {
		g_object_unref(uintptr(widget))
	}
}

func menuItemSetHidden(item *linuxMenuItem, hidden bool) {
	if item.parentMenu == nilPointer {
		return
	}
	if hidden {
		g_menu_remove(uintptr(item.parentMenu), int32(item.menuIndex))
	} else {
		g_menu_insert_item(uintptr(item.parentMenu), int32(item.menuIndex), uintptr(item.native))
	}
}

func menuCheckItemNew(label string, bitmap []byte) pointer {
	return nilPointer
}

func menuCheckItemNewWithId(label string, bitmap []byte, itemId uint, checked bool) pointer {
	actionName := generateActionName(itemId)
	gitem := createCheckMenuItem(label, actionName, itemId, checked)

	menuItemsLock.Lock()
	menuItemIds[pointer(gitem)] = itemId
	menuItemsLock.Unlock()
	return pointer(gitem)
}

func menuItemSetChecked(widget pointer, checked bool) {
	if widget == nilPointer {
		return
	}
	menuItemsLock.RLock()
	itemId, exists := menuItemIds[widget]
	menuItemsLock.RUnlock()
	if !exists {
		return
	}
	actionName, ok := lookupActionName(itemId)
	if !ok {
		return
	}
	setActionState(actionName, checked)
}

func menuItemSetDisabled(widget pointer, disabled bool) {
	if widget == nilPointer {
		return
	}
	menuItemsLock.RLock()
	itemId, exists := menuItemIds[widget]
	menuItemsLock.RUnlock()
	if !exists {
		return
	}
	actionName, ok := lookupActionName(itemId)
	if !ok {
		return
	}
	setActionEnabled(actionName, !disabled)
}

func menuItemSetLabel(widget pointer, label string) {
	if widget == nilPointer {
		return
	}
	g_menu_item_set_label(uintptr(widget), label)
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
	return nilPointer
}

func menuRadioItemNewWithId(label string, itemId uint, checked bool) pointer {
	actionName := generateActionName(itemId)
	gitem := createCheckMenuItem(label, actionName, itemId, checked)

	menuItemsLock.Lock()
	menuItemIds[pointer(gitem)] = itemId
	menuItemsLock.Unlock()
	return pointer(gitem)
}

func menuRadioItemNewWithGroup(label string, itemId uint, groupId uint, checkedId uint) pointer {
	actionName := fmt.Sprintf("radio_group_%d", groupId)
	targetValue := fmt.Sprintf("%d", itemId)
	initialValue := fmt.Sprintf("%d", checkedId)

	gitem := createRadioMenuItem(label, actionName, targetValue, initialValue, itemId)

	menuItemsLock.Lock()
	menuItemIds[pointer(gitem)] = itemId
	menuItemsLock.Unlock()
	return pointer(gitem)
}

// Keyboard accelerator support for GTK4 menus

// namedKeysToGTK maps Wails key names to GDK keysym values
// These are X11 keysym values that GDK uses
var namedKeysToGTK = map[string]uint32{
	"backspace": 0xff08,
	"tab":       0xff09,
	"return":    0xff0d,
	"enter":     0xff0d,
	"escape":    0xff1b,
	"left":      0xff51,
	"right":     0xff53,
	"up":        0xff52,
	"down":      0xff54,
	"space":     0xff80,
	"delete":    0xff9f,
	"home":      0xff95,
	"end":       0xff9c,
	"page up":   0xff9a,
	"page down": 0xff9b,
	"f1":        0xffbe,
	"f2":        0xffbf,
	"f3":        0xffc0,
	"f4":        0xffc1,
	"f5":        0xffc2,
	"f6":        0xffc3,
	"f7":        0xffc4,
	"f8":        0xffc5,
	"f9":        0xffc6,
	"f10":       0xffc7,
	"f11":       0xffc8,
	"f12":       0xffc9,
	"f13":       0xffca,
	"f14":       0xffcb,
	"f15":       0xffcc,
	"f16":       0xffcd,
	"f17":       0xffce,
	"f18":       0xffcf,
	"f19":       0xffd0,
	"f20":       0xffd1,
	"f21":       0xffd2,
	"f22":       0xffd3,
	"f23":       0xffd4,
	"f24":       0xffd5,
	"f25":       0xffd6,
	"f26":       0xffd7,
	"f27":       0xffd8,
	"f28":       0xffd9,
	"f29":       0xffda,
	"f30":       0xffdb,
	"f31":       0xffdc,
	"f32":       0xffdd,
	"f33":       0xffde,
	"f34":       0xffdf,
	"f35":       0xffe0,
	"numlock":   0xff7f,
}

// parseKeyGTK converts a Wails key string to a GDK keysym value
func parseKeyGTK(key string) uint32 {
	// Check named keys first
	if result, found := namedKeysToGTK[key]; found {
		return result
	}
	// For single character keys, convert using gdk_unicode_to_keyval
	if len(key) != 1 {
		return 0
	}
	return gdk_unicode_to_keyval(uint32(key[0]))
}

// parseModifiersGTK converts Wails modifiers to GDK modifier type
func parseModifiersGTK(modifiers []modifier) uint32 {
	var result uint32

	for _, mod := range modifiers {
		switch mod {
		case ShiftKey:
			result |= gdkShiftMask
		case ControlKey, CmdOrCtrlKey:
			result |= gdkControlMask
		case OptionOrAltKey:
			result |= gdkAltMask
		case SuperKey:
			result |= gdkSuperMask
		}
	}
	return result
}

// acceleratorToGTK converts a Wails accelerator to GTK key/modifiers
func acceleratorToGTK(accel *accelerator) (uint32, uint32) {
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
	actionName, ok := lookupActionName(itemId)
	if !ok {
		return
	}

	// Get the GtkApplication pointer
	app := getNativeApplication()
	if app == nil || app.application == nilPointer {
		return
	}

	// Convert accelerator to GTK format
	key, mods := acceleratorToGTK(accel)
	if key == 0 {
		return
	}

	// Build accelerator string using GTK's function
	accelString := gtk_accelerator_name(key, mods)
	if accelString == 0 {
		return
	}
	setActionAccelerator(uintptr(app.application), actionName, takeGString(accelString))
}

// screen related

// monitorScale returns the monitor's scale factor. gdk_monitor_get_scale
// (GTK 4.14+) reports fractional scaling; older GTK4 only has the integer
// gdk_monitor_get_scale_factor. This runtime fallback replaces the cgo
// build's compile-time GTK version floor.
func monitorScale(monitor uintptr) float64 {
	if gdk_monitor_get_scale != nil {
		return gdk_monitor_get_scale(monitor)
	}
	return float64(gdk_monitor_get_scale_factor(monitor))
}

func monitorGeometry(monitor uintptr) gdkRectangle {
	var geometry gdkRectangle
	gdk_monitor_get_geometry(monitor, uintptr(unsafe.Pointer(&geometry)))
	return geometry
}

// buildScreen assembles a Screen from a monitor's logical geometry and scale.
// GTK4's gdk_monitor_get_geometry returns logical (DIP) coordinates;
// PhysicalBounds needs physical pixel dimensions for proper DPI scaling.
func buildScreen(id string, monitor uintptr, isPrimary bool) *Screen {
	geometry := monitorGeometry(monitor)
	scaleFactor := monitorScale(monitor)
	name := gdk_monitor_get_model(monitor)

	x := int(geometry.x)
	y := int(geometry.y)
	width := int(geometry.width)
	height := int(geometry.height)

	physical := Rect{
		X:      int(float64(x) * scaleFactor),
		Y:      int(float64(y) * scaleFactor),
		Height: int(float64(height) * scaleFactor),
		Width:  int(float64(width) * scaleFactor),
	}

	return &Screen{
		ID:          id,
		Name:        name,
		IsPrimary:   isPrimary,
		ScaleFactor: float32(scaleFactor),
		X:           x,
		Y:           y,
		Size: Size{
			Height: height,
			Width:  width,
		},
		Bounds: Rect{
			X:      x,
			Y:      y,
			Height: height,
			Width:  width,
		},
		WorkArea: Rect{
			X:      x,
			Y:      y,
			Height: height,
			Width:  width,
		},
		PhysicalBounds:   physical,
		PhysicalWorkArea: physical,
		Rotation:         0.0,
	}
}

func getScreenByIndex(display uintptr, index int) *Screen {
	monitors := gdk_display_get_monitors(display)
	monitor := g_list_model_get_item(monitors, uint32(index))
	if monitor == 0 {
		return nil
	}
	defer g_object_unref(monitor)
	return buildScreen(fmt.Sprintf("%d", index), monitor, index == 0)
}

func getScreens(app pointer) ([]*Screen, error) {
	var screens []*Screen
	display := gdk_display_get_default()
	monitors := gdk_display_get_monitors(display)
	count := g_list_model_get_n_items(monitors)
	for i := 0; i < int(count); i++ {
		screens = append(screens, getScreenByIndex(display, i))
	}
	return screens, nil
}

// widgets
func (w *linuxWebviewWindow) setEnabled(enabled bool) {
	gtk_widget_set_sensitive(uintptr(w.window), gbool(enabled))
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func widgetSetVisible(widget pointer, hidden bool) {
	gtk_widget_set_visible(uintptr(widget), gbool(!hidden))
}

func (w *linuxWebviewWindow) close() {
	gtk_window_destroy(uintptr(w.window))
	getNativeApplication().unregisterWindow(windowPointer(w.window))
}

func (w *linuxWebviewWindow) enableDND() {
	enableDNDGo(uintptr(w.webview), uintptr(w.parent.id))
}

func (w *linuxWebviewWindow) disableDND() {
	// Mirrors the cgo backend: disabling DND after enabling is not
	// implemented for GTK4.
}

func (w *linuxWebviewWindow) execJS(js string) {
	InvokeAsync(func() {
		script := cString(js)
		defer g_free(script)
		// WebKitGTK 6.0 uses webkit_web_view_evaluate_javascript
		webkit_web_view_evaluate_javascript(w.webKitWebView(), script, len(js), 0, 0, 0, 0, 0)
	})
}

// Preallocated buffer for drag-over JS calls, matching the cgo backend's
// allocation-free hot path (drop-motion events fire at pointer-move rate).
var dragOverJSBuffer uintptr
var dragOverJSOnce sync.Once

func (w *linuxWebviewWindow) execJSDragOver(x, y int) {
	dragOverJSOnce.Do(func() {
		dragOverJSBuffer = g_malloc0(64)
	})
	buf := unsafe.Slice((*byte)(unsafe.Pointer(dragOverJSBuffer)), 64)
	n := copy(buf, "window._wails.handleDragOver(")
	n += writeInt(buf[n:], x)
	buf[n] = ','
	n++
	n += writeInt(buf[n:], y)
	buf[n] = ')'
	n++
	buf[n] = 0

	webkit_web_view_evaluate_javascript(w.webKitWebView(), dragOverJSBuffer, n, 0, 0, 0, 0, 0)
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
	display := gdk_display_get_default()
	if display == 0 {
		return 0, 0, nil
	}

	monitors := gdk_display_get_monitors(display)
	if monitors == 0 {
		return 0, 0, nil
	}

	n := g_list_model_get_n_items(monitors)
	if n == 0 {
		return 0, 0, nil
	}

	var primaryMonitor uintptr
	for i := uint32(0); i < n; i++ {
		mon := g_list_model_get_item(monitors, i)
		if mon != 0 {
			primaryMonitor = mon
			break
		}
	}

	if primaryMonitor == 0 {
		return 0, 0, nil
	}
	defer g_object_unref(primaryMonitor)

	screen := buildScreen("0", primaryMonitor, true)

	centerX := screen.X + screen.Size.Width/2
	centerY := screen.Y + screen.Size.Height/2

	return centerX, centerY, screen
}

func (w *linuxWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
	if w.gtkmenu != nilPointer {
		// GTK4: Different menu destruction
		w.gtkmenu = nilPointer
	}
	gtk_window_destroy(uintptr(w.window))
}

func (w *linuxWebviewWindow) fullscreen() {
	gtk_window_fullscreen(uintptr(w.window))
}

func (w *linuxWebviewWindow) getCurrentMonitor() uintptr {
	display := gtk_widget_get_display(uintptr(w.window))
	surface := gtk_native_get_surface(uintptr(w.window))
	if surface != 0 {
		monitor := gdk_display_get_monitor_at_surface(display, surface)
		if monitor != 0 {
			return monitor
		}
	}
	return 0
}

func (w *linuxWebviewWindow) getScreen() (*Screen, error) {
	monitor := w.getCurrentMonitor()
	if monitor == 0 {
		return nil, fmt.Errorf("no monitor found")
	}
	screen := buildScreen(fmt.Sprintf("%d", w.id), monitor, false)
	return screen, nil
}

func (w *linuxWebviewWindow) getCurrentMonitorGeometry() (x int, y int, width int, height int, scaleFactor float64) {
	monitor := w.getCurrentMonitor()
	if monitor == 0 {
		return -1, -1, -1, -1, 1
	}
	geometry := monitorGeometry(monitor)
	scaleFactor = monitorScale(monitor)
	return int(geometry.x), int(geometry.y), int(geometry.width), int(geometry.height), scaleFactor
}

func (w *linuxWebviewWindow) size() (int, int) {
	var width, height int32
	gtk_window_get_default_size(uintptr(w.window),
		uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	if width <= 0 || height <= 0 {
		width = gtk_widget_get_width(uintptr(w.window))
		height = gtk_widget_get_height(uintptr(w.window))
	}
	return int(width), int(height)
}

func (w *linuxWebviewWindow) relativePosition() (int, int) {
	x, y := w.position()
	monitor := w.getCurrentMonitor()
	if monitor == 0 {
		return x, y
	}
	geometry := monitorGeometry(monitor)
	return x - int(geometry.x), y - int(geometry.y)
}

func (w *linuxWebviewWindow) windowHide() {
	gtk_widget_set_visible(uintptr(w.window), 0)
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	return gtk_window_is_fullscreen(uintptr(w.window)) != 0
}

func (w *linuxWebviewWindow) isFocused() bool {
	return gtk_window_is_active(uintptr(w.window)) != 0
}

func (w *linuxWebviewWindow) isMaximised() bool {
	return gtk_window_is_maximized(uintptr(w.window)) != 0 && !w.isFullscreen()
}

func (w *linuxWebviewWindow) isMinimised() bool {
	surface := gtk_native_get_surface(uintptr(w.window))
	if surface == 0 {
		return false
	}
	state := gdk_toplevel_get_state(surface)
	return state&gdkToplevelStateMinimized != 0
}

func (w *linuxWebviewWindow) isVisible() bool {
	return gtk_widget_is_visible(uintptr(w.window)) != 0
}

func (w *linuxWebviewWindow) maximise() {
	gtk_window_maximize(uintptr(w.window))
}

func (w *linuxWebviewWindow) minimise() {
	gtk_window_minimize(uintptr(w.window))
}

func windowNew(application pointer, menu pointer, menuStyle LinuxMenuStyle, windowId uint, gpuPolicy WebviewGpuPolicy) (window, webview, vbox pointer) {
	window = pointer(gtk_application_window_new(uintptr(application)))
	g_object_ref_sink(uintptr(window))

	attachActionGroupToWidget(uintptr(window))

	webview = windowNewWebview(windowId, gpuPolicy)
	vbox = pointer(gtk_box_new(gtkOrientationVertical, 0))
	gtk_widget_set_name(uintptr(vbox), "webview-box")

	gtk_window_set_child(uintptr(window), uintptr(vbox))

	if menu != nilPointer {
		switch menuStyle {
		case LinuxMenuStylePrimaryMenu:
			headerBar := createHeaderBarWithMenu(uintptr(menu))
			gtk_window_set_titlebar(uintptr(window), headerBar)
		default:
			menuBar := createMenuBarFromModel(uintptr(menu))
			gtk_box_prepend(uintptr(vbox), menuBar)
		}
	}

	gtk_box_append(uintptr(vbox), uintptr(webview))
	gtk_widget_set_vexpand(uintptr(webview), 1)
	gtk_widget_set_hexpand(uintptr(webview), 1)
	return
}

func windowNewWebview(parentId uint, gpuPolicy WebviewGpuPolicy) pointer {
	manager := webkit_user_content_manager_new()
	// WebKitGTK 6.0: register_script_message_handler(manager, name, world_name)
	webkit_user_content_manager_register_script_message_handler(manager, "external", 0)

	// Create web view with settings
	settings := webkit_settings_new()
	// WebKitGTK 6.0 removed webkit_web_view_new_with_user_content_manager;
	// user-content-manager is a construct-only property, so build the view
	// with g_object_new_with_properties (g_object_new is variadic).
	webView := gObjectNewWithObjectProperty(webkit_web_view_get_type(), "user-content-manager", manager)

	saveWebviewToContentManager(manager, webView)
	saveWindowID(webView, parentId)
	saveWindowID(manager, parentId)

	// GPU policy
	// WebKitGTK 6.0: WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND was removed
	// Only ALWAYS and NEVER are available
	switch gpuPolicy {
	case WebviewGpuPolicyNever:
		webkit_settings_set_hardware_acceleration_policy(settings, webkitHardwareAccelerationPolicyNever)
	case WebviewGpuPolicyAlways:
		webkit_settings_set_hardware_acceleration_policy(settings, webkitHardwareAccelerationPolicyAlways)
	default:
		// Default to ALWAYS (was ON_DEMAND in older WebKitGTK)
		webkit_settings_set_hardware_acceleration_policy(settings, webkitHardwareAccelerationPolicyAlways)
	}

	webkit_web_view_set_settings(webView, settings)

	// Register URI scheme handler
	registerURIScheme.Do(func() {
		webContext := webkit_web_view_get_context(webView)
		webkit_web_context_register_uri_scheme(webContext, "wails", onProcessRequestPtr, 0, 0)
	})

	// Start the periodic signal-handler fix now that a WebView exists and JSC
	// can actually initialise. Anchoring to first webview creation (not appNew)
	// ensures the 5s window covers the JSC lazy-init race window.
	fixSignalHandlers.Do(func() {
		installSignalHandlers()
		scheduleSignalHandlerFix()
	})

	return pointer(webView)
}

func (w *linuxWebviewWindow) webKitWebView() uintptr {
	// The webview widget IS the WebKitWebView instance (WEBKIT_WEB_VIEW is
	// just a checked cast in C).
	return uintptr(w.webview)
}

func (w *linuxWebviewWindow) present() {
	gtk_window_present(uintptr(w.window))
}

func (w *linuxWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		gtk_window_set_title(uintptr(w.window), title)
	}
}

func (w *linuxWebviewWindow) setSize(width, height int) {
	gtk_window_set_default_size(uintptr(w.window), int32(width), int32(height))
}

func (w *linuxWebviewWindow) setDefaultSize(width int, height int) {
	gtk_window_set_default_size(uintptr(w.window), int32(width), int32(height))
}

func windowSetGeometryHints(window pointer, minWidth, minHeight, maxWidth, maxHeight int) {
	if minWidth > 0 && minHeight > 0 {
		gtk_widget_set_size_request(uintptr(window), int32(minWidth), int32(minHeight))
	}
	if maxWidth > 0 || maxHeight > 0 {
		windowSetMaxSize(uintptr(window), maxWidth, maxHeight)
	}
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	gtk_window_set_resizable(uintptr(w.window), gbool(resizable))
	w.execJS(fmt.Sprintf("if(window._wails&&window._wails.setResizable)window._wails.setResizable(%v);", resizable))
}

func (w *linuxWebviewWindow) move(x, y int) {
	// The GDK_IS_X11_DISPLAY-equivalent check inside handles X11 vs Wayland
	// correctly, including XWayland and GDK_BACKEND=x11 scenarios.
	windowMoveX11(uintptr(w.window), x, y)
}

func (w *linuxWebviewWindow) position() (int, int) {
	// Returns 0,0 on non-X11 displays, matching the cgo backend.
	return windowGetPositionX11(uintptr(w.window))
}

func (w *linuxWebviewWindow) unfullscreen() {
	gtk_window_unfullscreen(uintptr(w.window))
	w.unmaximise()
}

func (w *linuxWebviewWindow) unmaximise() {
	gtk_window_unmaximize(uintptr(w.window))
}

func (w *linuxWebviewWindow) windowShow() {
	if w.window == nilPointer {
		return
	}
	gtk_window_present(uintptr(w.window))
	// Re-apply always-on-top state now that the surface exists.
	windowApplyPendingAlwaysOnTop(uintptr(w.window))
}

func (w *linuxWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	// X11 only: uses _NET_WM_STATE_ABOVE. No-op on Wayland (no standard protocol).
	windowSetAlwaysOnTop(uintptr(w.window), alwaysOnTop)
}

func (w *linuxWebviewWindow) setBorderless(borderless bool) {
	gtk_window_set_decorated(uintptr(w.window), gbool(!borderless))
}

func (w *linuxWebviewWindow) setFrameless(frameless bool) {
	gtk_window_set_decorated(uintptr(w.window), gbool(!frameless))
	w.execJS(fmt.Sprintf("if(window._wails&&window._wails.flags)window._wails.flags.frameless=%v;", frameless))
}

func (w *linuxWebviewWindow) setTransparent() {
	// GTK4: Transparency via CSS - different from GTK3
}

func (w *linuxWebviewWindow) setBackgroundColour(colour RGBA) {
	rgba := gdkRGBA{
		red:   float32(colour.Red) / 255.0,
		green: float32(colour.Green) / 255.0,
		blue:  float32(colour.Blue) / 255.0,
		alpha: float32(colour.Alpha) / 255.0,
	}
	webkit_web_view_set_background_color(w.webKitWebView(), uintptr(unsafe.Pointer(&rgba)))
}

func (w *linuxWebviewWindow) setIcon(icon pointer) {
	// GTK4 removed gtk_window_set_icon. Window icons are set via the
	// application's .desktop file at the desktop-integration level.
}

func (w *linuxWebviewWindow) startDrag() error {
	beginWindowDrag(uintptr(w.window),
		int32(w.drag.MouseButton),
		float64(w.drag.XRoot),
		float64(w.drag.YRoot),
		w.drag.DragTime)
	return nil
}

// gdkSurfaceEdgeForBorder maps the border strings sent by the Wails runtime
// (as injected by drag.ts — "n-resize", "ne-resize", etc.) to the
// corresponding GdkSurfaceEdge value expected by gdk_toplevel_begin_resize.
// GdkSurfaceEdge values (gdk/gdkenums.h): NORTH_WEST=0, NORTH=1, NORTH_EAST=2,
// WEST=3, EAST=4, SOUTH_WEST=5, SOUTH=6, SOUTH_EAST=7.
var gdkSurfaceEdgeForBorder = map[string]int32{
	"nw-resize": 0,
	"n-resize":  1,
	"ne-resize": 2,
	"w-resize":  3,
	"e-resize":  4,
	"sw-resize": 5,
	"s-resize":  6,
	"se-resize": 7,
}

func (w *linuxWebviewWindow) startResize(border string) error {
	edge, ok := gdkSurfaceEdgeForBorder[border]
	if !ok {
		return fmt.Errorf("unknown resize border: %q", border)
	}
	// Drag state (mouse button, root coords, timestamp) was captured by
	// the click gesture in the GTK4 controller and stored on w.drag.
	beginWindowResize(uintptr(w.window), edge,
		int32(w.drag.MouseButton),
		float64(w.drag.XRoot),
		float64(w.drag.YRoot),
		w.drag.DragTime)
	return nil
}

func (w *linuxWebviewWindow) getZoom() float64 {
	return webkit_web_view_get_zoom_level(w.webKitWebView())
}

func (w *linuxWebviewWindow) setZoom(zoom float64) {
	if zoom < 1 {
		zoom = 1
	}
	webkit_web_view_set_zoom_level(w.webKitWebView(), zoom)
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
	webkit_web_view_load_uri(w.webKitWebView(), "wails://")
}

func (w *linuxWebviewWindow) setURL(uri string) {
	webkit_web_view_load_uri(w.webKitWebView(), uri)
}

func (w *linuxWebviewWindow) setHTML(html string) {
	webkit_web_view_load_alternate_html(w.webKitWebView(), html, "wails://", "")
}

func (w *linuxWebviewWindow) flash(_ bool) {}

func (w *linuxWebviewWindow) setOpacity(opacity float64) {
	gtk_widget_set_opacity(uintptr(w.window), opacity)
}

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
	winID := uintptr(w.parent.ID())

	setupWindowEventControllers(uintptr(w.window), uintptr(w.webview), winID)

	signalConnect(uintptr(w.webview), "load-changed", handleLoadChangedPtr, winID)
	signalConnect(uintptr(w.webview), "permission-request", handlePermissionRequestPtr, winID)

	contentManager := webkit_web_view_get_user_content_manager(w.webKitWebView())
	signalConnect(contentManager, "script-message-received::external", sendMessageToBackendPtr, 0)
}

// onProcessRequestGo forwards a WebKitURISchemeRequest to the asset server
// (called from the registered URI scheme trampoline on the main thread).
func onProcessRequestGo(request uintptr) {
	webView := webkit_uri_scheme_request_get_web_view(request)
	windowId := windowIDFromObject(webView)
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

// ============================================================================
// GTK4 Dialog System
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

// fileDialogCallback delivers the chosen paths to the waiting dialog channel.
//
// BUGS_FOUND #6: the cgo backend sends each path inline on the GTK main
// thread into a channel with a fixed buffer of 100 — selecting more files
// than that deadlocks the main loop if the consumer isn't already draining.
// The results are handed off to a goroutine here, so the main thread never
// blocks regardless of selection size or consumer behaviour.
func fileDialogCallback(requestID uint, files []string, cancelled bool) {
	dialogRequestMutex.Lock()
	ch, ok := fileDialogCallbacks[requestID]
	if ok {
		delete(fileDialogCallbacks, requestID)
	}
	dialogRequestMutex.Unlock()

	if !ok {
		return
	}

	if cancelled {
		close(ch)
		return
	}

	go func() {
		defer handlePanic()
		for _, file := range files {
			ch <- file
		}
		close(ch)
	}()
}

func alertDialogCallback(requestID uint, buttonIndex int) {
	dialogRequestMutex.Lock()
	ch, ok := alertDialogCallbacks[requestID]
	if ok {
		delete(alertDialogCallbacks, requestID)
	}
	dialogRequestMutex.Unlock()

	if !ok {
		return
	}

	ch <- buttonIndex
	close(ch)
}

func runChooserDialog(window pointer, allowMultiple, createFolders, showHidden bool, currentFolder, title string, action int, acceptLabel string, filters []FileFilter) (chan string, error) {
	requestID := nextDialogRequestID()
	resultChan := make(chan string, 100)

	dialogRequestMutex.Lock()
	fileDialogCallbacks[requestID] = resultChan
	dialogRequestMutex.Unlock()

	InvokeAsync(func() {
		dialog := gtk_file_dialog_new()
		gtk_file_dialog_set_title(dialog, title)

		// Create filter list if we have filters
		if len(filters) > 0 {
			filterStore := g_list_store_new(gtk_file_filter_get_type())
			defer g_object_unref(filterStore)

			for _, filter := range filters {
				addFileFilter(dialog, filterStore, filter.DisplayName, filter.Pattern)
			}
			gtk_file_dialog_set_filters(dialog, filterStore)
		}

		if currentFolder != "" {
			file := g_file_new_for_path(currentFolder)
			gtk_file_dialog_set_initial_folder(dialog, file)
			g_object_unref(file)
		}

		if acceptLabel != "" {
			gtk_file_dialog_set_accept_label(dialog, acceptLabel)
		}

		isFolder := action == 2
		isSave := action == 1

		if isSave {
			showSaveFileDialog(uintptr(window), dialog, requestID)
		} else {
			showOpenFileDialog(uintptr(window), dialog, requestID, allowMultiple, isFolder)
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
			window = pointer(uintptr(nativeWindow))
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
			window = pointer(uintptr(nativeWindow))
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

func dialogTypeToIconName(dialogType DialogType) string {
	switch dialogType {
	case InfoDialogType:
		return "dialog-information-symbolic"
	case WarningDialogType:
		return "dialog-warning-symbolic"
	case ErrorDialogType:
		return "dialog-error-symbolic"
	case QuestionDialogType:
		return "dialog-question-symbolic"
	default:
		return ""
	}
}

func runQuestionDialog(parent pointer, options *MessageDialog) int {
	requestID := nextDialogRequestID()
	resultChan := make(chan int, 1)

	dialogRequestMutex.Lock()
	alertDialogCallbacks[requestID] = resultChan
	dialogRequestMutex.Unlock()

	InvokeAsync(func() {
		var iconName string
		var iconData []byte
		if len(options.Icon) > 0 {
			iconData = options.Icon
		} else {
			iconName = dialogTypeToIconName(options.DialogType)
		}

		buttons := options.Buttons
		if len(buttons) == 0 {
			buttons = []*Button{{Label: "OK", IsDefault: true}}
		}

		buttonLabels := make([]string, len(buttons))
		for i, btn := range buttons {
			buttonLabels[i] = btn.Label
		}

		defaultButton := -1
		cancelButton := -1
		destructiveButton := -1
		for i, btn := range buttons {
			if btn.IsDefault {
				defaultButton = i
			}
			if btn.IsCancel {
				cancelButton = i
			}
		}

		if options.DialogType == ErrorDialogType || options.DialogType == WarningDialogType {
			if defaultButton >= 0 && !buttons[defaultButton].IsCancel {
				destructiveButton = defaultButton
				defaultButton = -1
			}
		}

		showMessageDialog(uintptr(parent), options.Title, options.Message,
			iconName, iconData, buttonLabels,
			defaultButton, cancelButton, destructiveButton, requestID)
	})

	// Wait for result
	result := <-resultChan
	return result
}

func getPrimaryScreen() (*Screen, error) {
	display := gdk_display_get_default()
	monitors := gdk_display_get_monitors(display)
	if monitors == 0 {
		return nil, fmt.Errorf("no monitors found")
	}
	count := g_list_model_get_n_items(monitors)
	if count == 0 {
		return nil, fmt.Errorf("no monitors found")
	}
	monitor := g_list_model_get_item(monitors, 0)
	if monitor == 0 {
		return nil, fmt.Errorf("failed to get primary monitor")
	}
	defer g_object_unref(monitor)

	return buildScreen("0", monitor, true), nil
}

func openDevTools(wv pointer) {
	inspector := webkit_web_view_get_inspector(uintptr(wv))
	webkit_web_inspector_show(inspector)
}

func enableDevTools(wv pointer) {
	settings := webkit_web_view_get_settings(uintptr(wv))
	enabled := webkit_settings_get_enable_developer_extras(settings)
	if enabled == 0 {
		webkit_settings_set_enable_developer_extras(settings, 1)
	} else {
		webkit_settings_set_enable_developer_extras(settings, 0)
	}
}

// splitAndTrim splits s on sep and trims surrounding whitespace from each part.
func splitAndTrim(s, sep string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i:i+1] == sep {
			part := s[start:i]
			// trim spaces/tabs
			for len(part) > 0 && (part[0] == ' ' || part[0] == '\t') {
				part = part[1:]
			}
			for len(part) > 0 && (part[len(part)-1] == ' ' || part[len(part)-1] == '\t') {
				part = part[:len(part)-1]
			}
			out = append(out, part)
			start = i + 1
		}
	}
	return out
}
