//go:build linux && purego && !gtk3 && !android && !server

package application

// Pure-Go port of linux_cgo.c: the GTK signal trampolines, main-thread
// dispatch, GAction-based menu machinery, GTK4 dialogs, drag-and-drop and the
// X11 window helpers. Every C callback becomes a package-level
// purego.NewCallback — a fixed, small set (purego's callback slots are capped
// process-wide and never freed, so nothing here creates callbacks per
// call/window/item).

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// ----------------------------------------------------------------------------
// Constants (values from the GTK4/GDK/GLib/WebKitGTK-6.0 headers; there is no
// compile step to import them from, so they are transcribed here)
// ----------------------------------------------------------------------------

const (
	gSourceRemove    = 0 // G_SOURCE_REMOVE (FALSE)
	gSourceContinue  = 1 // G_SOURCE_CONTINUE (TRUE)
	gPriorityDefault = 0

	gApplicationDefaultFlags = 0 // G_APPLICATION_DEFAULT_FLAGS

	gtkOrientationHorizontal = 0
	gtkOrientationVertical   = 1

	gtkAlignCenter = 3 // GtkAlign: FILL=0 START=1 END=2 CENTER=3

	gtkPosBottom = 3 // GtkPositionType: LEFT=0 RIGHT=1 TOP=2 BOTTOM=3

	gtkPhaseCapture = 1 // GtkPropagationPhase: NONE=0 CAPTURE=1 BUBBLE=2 TARGET=3

	gdkActionCopy = 1 << 0 // GdkDragAction

	gdkCurrentTime = 0 // GDK_CURRENT_TIME

	// GdkToplevelState (gdk/gdktoplevel.h)
	gdkToplevelStateMinimized = 1 << 0

	// GdkModifierType (gdk/gdkenums.h)
	gdkShiftMask   = 1 << 0
	gdkControlMask = 1 << 2
	gdkAltMask     = 1 << 3
	gdkSuperMask   = 1 << 26

	gdkKeyEscape = 0xff1b

	// WebKitLoadEvent (webkit/WebKitWebView.h)
	webkitLoadStarted    = 0
	webkitLoadRedirected = 1
	webkitLoadCommitted  = 2
	webkitLoadFinished   = 3

	// WebKitHardwareAccelerationPolicy (WebKitGTK 6.0: ON_DEMAND was removed,
	// leaving ALWAYS=0, NEVER=1)
	webkitHardwareAccelerationPolicyAlways = 0
	webkitHardwareAccelerationPolicyNever  = 1
)

// ----------------------------------------------------------------------------
// Main-thread dispatch (g_idle_add onto the GTK main loop)
// ----------------------------------------------------------------------------

var dispatchCallbackPtr = purego.NewCallback(func(data uintptr) uintptr {
	executeOnMainThread(uint(data))
	return gSourceRemove
})

func dispatchOnMainThread(id uint) {
	g_idle_add(dispatchCallbackPtr, uintptr(id))
}

// ----------------------------------------------------------------------------
// Signal handling (SA_ONSTACK fix, port of install_signal_handlers)
// ----------------------------------------------------------------------------

// glibc/musl struct sigaction on linux amd64/arm64:
//
//	0   sa_handler  (8)
//	8   sa_mask     (128)
//	136 sa_flags    (4)
//	144 sa_restorer (8)
const (
	sigactionSize     = 152
	sigactionFlagsOff = 136
	saOnStack         = 0x08000000
)

func fixSignal(sig syscall.Signal) {
	if libc_sigaction == nil {
		return
	}
	var st [sigactionSize + 8]byte
	p := uintptr(unsafe.Pointer(&st[0]))
	if libc_sigaction(int32(sig), 0, p) < 0 {
		return
	}
	flags := (*int32)(unsafe.Pointer(&st[sigactionFlagsOff]))
	*flags |= saOnStack
	libc_sigaction(int32(sig), p, 0)
}

// installSignalHandlers re-applies SA_ONSTACK to the signals Go cares about.
// GTK/WebKit install their own handlers without SA_ONSTACK; without this fix
// they run on goroutine stacks and crash the Go runtime.
//
// NOTE: SIGUSR1 is deliberately NOT fixed. WebKit's JavaScriptCore uses
// SIGUSR1 to suspend/resume threads for conservative GC stack scanning; once
// JSC owns that signal, forcing SA_ONSTACK onto its handler breaks GC thread
// synchronisation and freezes WebKit during idle collection. See issue #5527.
func installSignalHandlers() {
	for _, sig := range []syscall.Signal{
		syscall.SIGCHLD, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGABRT, syscall.SIGFPE, syscall.SIGTERM, syscall.SIGBUS,
		syscall.SIGSEGV, syscall.SIGXCPU, syscall.SIGXFSZ,
	} {
		fixSignal(sig)
	}
}

// WebKit's JSC lazily installs signal handlers without SA_ONSTACK when
// JavaScript first executes. This timer re-applies the fix every 50ms for the
// first 5 seconds, covering the JSC initialization window.
var signalFixRemaining int32

var signalFixTimeoutPtr = purego.NewCallback(func(data uintptr) uintptr {
	installSignalHandlers()
	signalFixRemaining--
	if signalFixRemaining <= 0 {
		return gSourceRemove
	}
	return gSourceContinue
})

func scheduleSignalHandlerFix() {
	signalFixRemaining = 100
	g_timeout_add_full(gPriorityDefault, 50, signalFixTimeoutPtr, 0, 0)
}

// ----------------------------------------------------------------------------
// Object data helpers (port of save_window_id & co)
// ----------------------------------------------------------------------------

func saveWindowID(object uintptr, id uint) {
	g_object_set_data(object, "windowid", uintptr(id))
}

func windowIDFromObject(object uintptr) uint {
	return uint(g_object_get_data(object, "windowid"))
}

func saveWebviewToContentManager(contentManager, webview uintptr) {
	g_object_set_data(contentManager, "webview", webview)
}

func getWebviewFromContentManager(contentManager uintptr) uintptr {
	return g_object_get_data(contentManager, "webview")
}

// ----------------------------------------------------------------------------
// Application activate
// ----------------------------------------------------------------------------

var activateLinuxPtr = purego.NewCallback(func(app, data uintptr) uintptr {
	nativeApp := getNativeApplication()
	nativeApp.markActivated()
	processApplicationEvent(uint(events.Linux.ApplicationStartup), nilPointer)
	return 0
})

func processApplicationEvent(eventID uint, _ pointer) {
	event := newApplicationEvent(events.ApplicationEventType(eventID))

	switch event.Id {
	case uint(events.Linux.SystemThemeChanged):
		isDark := globalApplication.Env.IsDarkMode()
		event.Context().setIsDarkMode(isDark)
	}
	applicationEvents <- event
}

func processWindowEvent(windowID uint, eventID uint) {
	windowEvents <- &windowEvent{
		WindowID: windowID,
		EventID:  eventID,
	}
}

// ----------------------------------------------------------------------------
// Window / webview signal trampolines
// (port of setupWindowEventControllers and the //export handlers)
// ----------------------------------------------------------------------------

var handleCloseRequestPtr = purego.NewCallback(func(window, data uintptr) uintptr {
	processWindowEvent(uint(data), uint(events.Linux.WindowDeleteEvent))
	return 1 // stop the default handler destroying the window
})

var handleNotifyStatePtr = purego.NewCallback(func(object, pspec, data uintptr) uintptr {
	windowId := uint(data)
	window, ok := globalApplication.Window.GetByID(windowId)
	if !ok || window == nil {
		return 0
	}
	lw := getLinuxWebviewWindow(window)
	if lw == nil {
		return 0
	}
	if lw.isMaximised() {
		processWindowEvent(windowId, uint(events.Linux.WindowDidResize))
	}
	if lw.isFullscreen() {
		processWindowEvent(windowId, uint(events.Linux.WindowDidResize))
	}
	return 0
})

var handleFocusEnterPtr = purego.NewCallback(func(controller, data uintptr) uintptr {
	processWindowEvent(uint(data), uint(events.Linux.WindowFocusIn))
	return 0
})

var handleFocusLeavePtr = purego.NewCallback(func(controller, data uintptr) uintptr {
	processWindowEvent(uint(data), uint(events.Linux.WindowFocusOut))
	return 0
})

var handleLoadChangedPtr = purego.NewCallback(func(wv uintptr, event int32, data uintptr) uintptr {
	switch event {
	case webkitLoadStarted:
		processWindowEvent(uint(data), uint(events.Linux.WindowLoadStarted))
	case webkitLoadRedirected:
		processWindowEvent(uint(data), uint(events.Linux.WindowLoadRedirected))
	case webkitLoadCommitted:
		processWindowEvent(uint(data), uint(events.Linux.WindowLoadCommitted))
	case webkitLoadFinished:
		// JSC is guaranteed to have initialised by page-load completion, so
		// re-apply SA_ONSTACK now to cover any handlers it installed during load.
		installSignalHandlers()
		processWindowEvent(uint(data), uint(events.Linux.WindowLoadFinished))
	}
	return 0
})

var handlePermissionRequestPtr = purego.NewCallback(func(wv, request, data uintptr) uintptr {
	// WebKitGTK denies any permission request nobody handles, so without this
	// getUserMedia always fails with NotAllowedError. Honour the window's
	// Permissions for camera/microphone; leave every other request to WebKit's
	// default handling (deny).
	if !gTypeInstanceIsA(request, webkit_user_media_permission_request_get_type()) {
		return 0
	}
	needAudio := webkit_user_media_permission_is_for_audio_device(request) != 0
	needVideo := webkit_user_media_permission_is_for_video_device(request) != 0
	if allowMediaCapture(uint(data), needAudio, needVideo) {
		webkit_permission_request_allow(request)
	} else {
		webkit_permission_request_deny(request)
	}
	return 1
})

var handleButtonPressedPtr = purego.NewCallback(func(gesture uintptr, nPress int32, x, y float64, data uintptr) uintptr {
	windowId := uint(data)
	window, ok := globalApplication.Window.GetByID(windowId)
	if !ok || window == nil {
		return 0
	}
	lw := getLinuxWebviewWindow(window)
	if lw == nil {
		return 0
	}
	button := gtk_gesture_single_get_current_button(gesture)
	lw.drag.MouseButton = uint(button)
	lw.drag.XRoot = int(x)
	lw.drag.YRoot = int(y)
	lw.drag.DragTime = uint32(gdkCurrentTime)
	return 0
})

var handleButtonReleasedPtr = purego.NewCallback(func(gesture uintptr, nPress int32, x, y float64, data uintptr) uintptr {
	windowId := uint(data)
	window, ok := globalApplication.Window.GetByID(windowId)
	if !ok || window == nil {
		return 0
	}
	lw := getLinuxWebviewWindow(window)
	if lw == nil {
		return 0
	}
	button := gtk_gesture_single_get_current_button(gesture)
	lw.endDrag(uint(button), int(x), int(y))
	return 0
})

var handleKeyPressedPtr = purego.NewCallback(func(controller uintptr, keyval, keycode uint32, state uint32, data uintptr) uintptr {
	windowID := uint(data)

	modifiers := uint(state)
	var acc accelerator

	if modifiers&gdkShiftMask != 0 {
		acc.Modifiers = append(acc.Modifiers, ShiftKey)
	}
	if modifiers&gdkControlMask != 0 {
		acc.Modifiers = append(acc.Modifiers, ControlKey)
	}
	if modifiers&gdkAltMask != 0 {
		acc.Modifiers = append(acc.Modifiers, OptionOrAltKey)
	}
	if modifiers&gdkSuperMask != 0 {
		acc.Modifiers = append(acc.Modifiers, SuperKey)
	}

	keyString, ok := VirtualKeyCodes[uint(keyval)]
	if !ok {
		return 0
	}
	acc.Key = keyString

	windowKeyEvents <- &windowKeyEvent{
		windowId:          windowID,
		acceleratorString: acc.String(),
	}

	return 0
})

// setupWindowEventControllers wires the GTK4-style event controllers for a
// window and its webview (port of the C function of the same name).
func setupWindowEventControllers(window, webview uintptr, winID uintptr) {
	// Close request (replaces delete-event)
	signalConnect(window, "close-request", handleCloseRequestPtr, winID)

	// Window state changes (maximize, fullscreen, etc)
	signalConnect(window, "notify::maximized", handleNotifyStatePtr, winID)
	signalConnect(window, "notify::fullscreened", handleNotifyStatePtr, winID)

	// Focus controller for window
	focusController := gtk_event_controller_focus_new()
	gtk_widget_add_controller(window, focusController)
	signalConnect(focusController, "enter", handleFocusEnterPtr, winID)
	signalConnect(focusController, "leave", handleFocusLeavePtr, winID)

	// Click gesture for webview (button press/release)
	clickGesture := gtk_gesture_click_new()
	gtk_gesture_single_set_button(clickGesture, 0) // listen to all buttons
	gtk_widget_add_controller(webview, clickGesture)
	signalConnect(clickGesture, "pressed", handleButtonPressedPtr, winID)
	signalConnect(clickGesture, "released", handleButtonReleasedPtr, winID)

	// Key controller for webview
	keyController := gtk_event_controller_key_new()
	gtk_widget_add_controller(webview, keyController)
	signalConnect(keyController, "key-pressed", handleKeyPressedPtr, winID)
}

// ----------------------------------------------------------------------------
// Asset scheme + script-message bridge
// ----------------------------------------------------------------------------

var onProcessRequestPtr = purego.NewCallback(func(request, data uintptr) uintptr {
	onProcessRequestGo(request)
	return 0
})

var sendMessageToBackendPtr = purego.NewCallback(func(contentManager, value, data uintptr) uintptr {
	// Get the windowID from the contentManager
	thisWindowID := windowIDFromObject(contentManager)

	webView := getWebviewFromContentManager(contentManager)
	var origin string
	if webView != 0 {
		currentURI := webkit_web_view_get_uri(webView)
		if currentURI != 0 {
			origin = goString(currentURI)
		}
	}

	// WebKitGTK 6.0: the JSCValue is passed directly
	msg := takeGString(jsc_value_to_string(value))
	windowMessageBuffer <- &windowMessage{
		windowId: thisWindowID,
		message:  msg,
		originInfo: &OriginInfo{
			Origin: origin,
		},
	}
	return 0
})

// ----------------------------------------------------------------------------
// Window drag / resize (GdkToplevel)
// ----------------------------------------------------------------------------

func toplevelForWindow(window uintptr) uintptr {
	native := gtk_widget_get_native(window)
	if native == 0 {
		return 0
	}
	// A GtkWindow's native surface is a GdkToplevel; gtk_native_get_surface
	// returns 0 before the window is realized.
	return gtk_native_get_surface(native)
}

func beginWindowDrag(window uintptr, button int32, x, y float64, timestamp uint32) {
	surface := toplevelForWindow(window)
	if surface == 0 {
		return
	}
	var device uintptr
	display := gdk_surface_get_display(surface)
	if seat := gdk_display_get_default_seat(display); seat != 0 {
		device = gdk_seat_get_pointer(seat)
	}
	gdk_toplevel_begin_move(surface, device, button, x, y, timestamp)
}

func beginWindowResize(window uintptr, edge int32, button int32, x, y float64, timestamp uint32) {
	surface := toplevelForWindow(window)
	if surface == 0 {
		return
	}
	var device uintptr
	display := gdk_surface_get_display(surface)
	if seat := gdk_display_get_default_seat(display); seat != 0 {
		device = gdk_seat_get_pointer(seat)
	}
	gdk_toplevel_begin_resize(surface, edge, device, button, x, y, timestamp)
}

// ----------------------------------------------------------------------------
// Drag and drop (GtkDropTarget + GtkDropControllerMotion)
// ----------------------------------------------------------------------------

var onDropAcceptPtr = purego.NewCallback(func(target, drop, data uintptr) uintptr {
	formats := gdk_drop_get_formats(drop)
	if gdk_content_formats_contain_gtype(formats, gdk_file_list_get_type()) != 0 {
		return 1
	}
	return 0
})

var onDropEnterPtr = purego.NewCallback(func(target uintptr, x, y float64, data uintptr) uintptr {
	onDropEnterGo(uint(data))
	return gdkActionCopy
})

var onDropLeavePtr = purego.NewCallback(func(target, data uintptr) uintptr {
	onDropLeaveGo(uint(data))
	return 0
})

var onDropMotionPtr = purego.NewCallback(func(target uintptr, x, y float64, data uintptr) uintptr {
	onDropMotionGo(int(x), int(y), uint(data))
	return gdkActionCopy
})

var onDropPtr = purego.NewCallback(func(target, value uintptr, x, y float64, data uintptr) uintptr {
	return uintptr(handleDrop(value, int(x), int(y), uint(data)))
})

var onMotionEnterPtr = purego.NewCallback(func(ctrl uintptr, x, y float64, data uintptr) uintptr {
	onDropEnterGo(uint(data))
	return 0
})

var onMotionLeavePtr = purego.NewCallback(func(ctrl, data uintptr) uintptr {
	onDropLeaveGo(uint(data))
	return 0
})

var onMotionMotionPtr = purego.NewCallback(func(ctrl uintptr, x, y float64, data uintptr) uintptr {
	onDropMotionGo(int(x), int(y), uint(data))
	return 0
})

func onDropEnterGo(windowId uint) {
	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return
	}
	if w, ok := targetWindow.(*WebviewWindow); ok {
		w.HandleDragEnter()
	}
}

func onDropLeaveGo(windowId uint) {
	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return
	}
	if w, ok := targetWindow.(*WebviewWindow); ok {
		w.HandleDragLeave()
	}
}

func onDropMotionGo(x, y int, windowId uint) {
	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return
	}
	if w, ok := targetWindow.(*WebviewWindow); ok {
		w.HandleDragOver(x, y)
	}
}

// handleDrop extracts the GFile list out of the dropped GValue and forwards
// the paths (port of on_drop). Returns 1 when the drop was handled.
func handleDrop(value uintptr, x, y int, windowId uint) int32 {
	if g_type_check_value_holds(value, gdk_file_list_get_type()) == 0 {
		return 0
	}
	fileList := g_value_get_boxed(value)
	if fileList == 0 {
		return 0
	}
	count := g_slist_length(fileList)
	if count == 0 {
		return 0
	}

	targetWindow, ok := globalApplication.Window.GetByID(windowId)
	if !ok || targetWindow == nil {
		return 0
	}

	var filenames []string
	for l := (*GSList)(unsafe.Pointer(fileList)); l != nil; l = l.next {
		if path := g_file_get_path(uintptr(l.data)); path != 0 {
			filenames = append(filenames, takeGString(path))
		}
	}

	targetWindow.InitiateFrontendDropProcessing(filenames, x, y)
	return 1
}

func enableDNDGo(widget uintptr, winID uintptr) {
	motionCtrl := gtk_drop_controller_motion_new()
	gtk_event_controller_set_propagation_phase(motionCtrl, gtkPhaseCapture)
	signalConnect(motionCtrl, "enter", onMotionEnterPtr, winID)
	signalConnect(motionCtrl, "leave", onMotionLeavePtr, winID)
	signalConnect(motionCtrl, "motion", onMotionMotionPtr, winID)
	gtk_widget_add_controller(widget, motionCtrl)

	target := gtk_drop_target_new(gdk_file_list_get_type(), gdkActionCopy)
	gtk_event_controller_set_propagation_phase(target, gtkPhaseCapture)
	signalConnect(target, "accept", onDropAcceptPtr, winID)
	signalConnect(target, "enter", onDropEnterPtr, winID)
	signalConnect(target, "leave", onDropLeavePtr, winID)
	signalConnect(target, "motion", onDropMotionPtr, winID)
	signalConnect(target, "drop", onDropPtr, winID)
	gtk_widget_add_controller(widget, target)
}

// ----------------------------------------------------------------------------
// Menus: GSimpleActionGroup + GAction activation
// ----------------------------------------------------------------------------

var (
	appActionGroup     uintptr
	appActionGroupOnce sync.Once
	appMenuModel       uintptr
)

func initAppActionGroup() {
	appActionGroupOnce.Do(func() {
		appActionGroup = g_simple_action_group_new()
	})
}

// onActionActivated handles plain and checkbox menu actions. The menu item id
// is attached to the GSimpleAction as object data at creation time (the C
// implementation used a heap-allocated MenuItemData for the same purpose).
var onActionActivatedPtr = purego.NewCallback(func(action, parameter, data uintptr) uintptr {
	menuActionActivated(uint(data))
	return 0
})

// onRadioActionActivated switches the stateful string action to the activated
// target and fires the menu item encoded in the target string.
var onRadioActionActivatedPtr = purego.NewCallback(func(action, parameter, data uintptr) uintptr {
	target := goString(g_variant_get_string(parameter, 0))
	g_simple_action_set_state(action, g_variant_new_string(target))
	itemId := 0
	for i := 0; i < len(target); i++ {
		if target[i] < '0' || target[i] > '9' {
			break
		}
		itemId = itemId*10 + int(target[i]-'0')
	}
	menuActionActivated(uint(itemId))
	return 0
})

var cachedVariantTypeString uintptr

func variantTypeString() uintptr {
	if cachedVariantTypeString == 0 {
		cachedVariantTypeString = g_variant_type_new("s")
	}
	return cachedVariantTypeString
}

// gMenuItemNew wraps g_menu_item_new, whose action argument may be NULL —
// an empty Go string would create an item bound to an action literally named
// "", so the C string is built by hand.
func gMenuItemNew(label, action string) uintptr {
	var cAction uintptr
	if action != "" {
		cAction = cString(action)
		defer g_free(cAction)
	}
	return g_menu_item_new(label, cAction)
}

func createMenuItem(label, actionName string, itemId uint) uintptr {
	initAppActionGroup()

	item := gMenuItemNew(label, "app."+actionName)

	action := g_simple_action_new(actionName, 0)
	signalConnect(action, "activate", onActionActivatedPtr, uintptr(itemId))
	g_action_map_add_action(appActionGroup, action)
	return item
}

func createCheckMenuItem(label, actionName string, itemId uint, initialState bool) uintptr {
	initAppActionGroup()

	item := gMenuItemNew(label, "app."+actionName)

	action := g_simple_action_new_stateful(actionName, 0, g_variant_new_boolean(gbool(initialState)))
	signalConnect(action, "activate", onActionActivatedPtr, uintptr(itemId))
	g_action_map_add_action(appActionGroup, action)
	return item
}

func createRadioMenuItem(label, actionName, target, initialValue string, itemId uint) uintptr {
	initAppActionGroup()

	item := gMenuItemNew(label, "")
	g_menu_item_set_action_and_target_value(item, "app."+actionName, g_variant_new_string(target))

	if g_action_map_lookup_action(appActionGroup, actionName) == 0 {
		action := g_simple_action_new_stateful(actionName, variantTypeString(), g_variant_new_string(initialValue))
		signalConnect(action, "activate", onRadioActionActivatedPtr, uintptr(itemId))
		g_action_map_add_action(appActionGroup, action)
	}
	return item
}

func createMenuBarFromModel(menuModel uintptr) uintptr {
	return gtk_popover_menu_bar_new_from_model(menuModel)
}

func createHeaderBarWithMenu(menuModel uintptr) uintptr {
	headerBar := gtk_header_bar_new()

	menuButton := gtk_menu_button_new()
	gtk_menu_button_set_icon_name(menuButton, "open-menu-symbolic")
	gtk_menu_button_set_menu_model(menuButton, menuModel)
	gtk_widget_set_tooltip_text(menuButton, "Main Menu")
	accessibleLabel(menuButton, "Main Menu")

	gtk_header_bar_pack_end(headerBar, menuButton)
	return headerBar
}

// accessibleLabel sets GTK_ACCESSIBLE_PROPERTY_LABEL on a widget.
// gtk_accessible_update_property is variadic, so use the array variant.
func accessibleLabel(widget uintptr, label string) {
	// GtkAccessibleProperty (gtkenums.h): AUTOCOMPLETE=0, DESCRIPTION,
	// HAS_POPUP, KEY_SHORTCUTS, LABEL=4, ...
	const gtkAccessiblePropertyLabel = 4

	var value gValue
	g_value_init(uintptr(unsafe.Pointer(&value)), g_type_from_name("gchararray"))
	g_value_set_string(uintptr(unsafe.Pointer(&value)), label)
	properties := []int32{gtkAccessiblePropertyLabel}
	gtk_accessible_update_property_value(widget, 1,
		uintptr(unsafe.Pointer(&properties[0])),
		uintptr(unsafe.Pointer(&value)))
	g_value_unset(uintptr(unsafe.Pointer(&value)))
}

func attachActionGroupToWidget(widget uintptr) {
	initAppActionGroup()
	gtk_widget_insert_action_group(widget, "app", appActionGroup)
}

func setActionAccelerator(app uintptr, actionName, accel string) {
	if app == 0 || accel == "" {
		return
	}
	cAccel := cString(accel)
	defer g_free(cAccel)
	accels := []uintptr{cAccel, 0}
	gtk_application_set_accels_for_action(app, "app."+actionName, uintptr(unsafe.Pointer(&accels[0])))
}

func setActionEnabled(actionName string, enabled bool) {
	if appActionGroup == 0 {
		return
	}
	action := g_action_map_lookup_action(appActionGroup, actionName)
	if action != 0 {
		g_simple_action_set_enabled(action, gbool(enabled))
	}
}

func setActionState(actionName string, state bool) {
	if appActionGroup == 0 {
		return
	}
	action := g_action_map_lookup_action(appActionGroup, actionName)
	if action != 0 {
		g_simple_action_set_state(action, g_variant_new_boolean(gbool(state)))
	}
}

func getActionState(actionName string) bool {
	if appActionGroup == 0 {
		return false
	}
	action := g_action_map_lookup_action(appActionGroup, actionName)
	if action == 0 {
		return false
	}
	state := g_action_get_state(action)
	if state == 0 {
		return false
	}
	result := g_variant_get_boolean(state) != 0
	g_variant_unref(state)
	return result
}

// Context menu

var onContextMenuClosedPtr = purego.NewCallback(func(popover, data uintptr) uintptr {
	// Unparent on the next main loop iteration so the popover finishes its
	// close animation/cleanup before being removed from the widget tree.
	g_idle_add(unparentWidgetPtr, popover)
	return 0
})

var unparentWidgetPtr = purego.NewCallback(func(widget uintptr) uintptr {
	gtk_widget_unparent(widget)
	return gSourceRemove
})

func showContextMenu(parent, menuModel uintptr, x, y int) {
	initAppActionGroup()

	popover := gtk_popover_menu_new_from_model(menuModel)
	gtk_widget_set_parent(popover, parent)
	gtk_popover_set_has_arrow(popover, 0)
	gtk_popover_set_position(popover, gtkPosBottom)

	// Ensure the menu actions resolve even if the parent's hierarchy does not
	// already expose the "app" action group.
	gtk_widget_insert_action_group(popover, "app", appActionGroup)

	rect := gdkRectangle{x: int32(x), y: int32(y), width: 1, height: 1}
	gtk_popover_set_pointing_to(popover, uintptr(unsafe.Pointer(&rect)))

	signalConnect(popover, "closed", onContextMenuClosedPtr, 0)

	gtk_popover_popup(popover)
}

// ----------------------------------------------------------------------------
// File dialogs (GtkFileDialog, async)
// ----------------------------------------------------------------------------

// The async finish callbacks receive the dialog request id via user_data —
// no allocation needed, unlike the C FileDialogData struct.

func finishSingleFile(finish func(uintptr, uintptr, uintptr) uintptr, source, res uintptr, requestID uint) {
	var gerr uintptr
	file := finish(source, res, uintptr(unsafe.Pointer(&gerr)))
	switch {
	case gerr != 0:
		g_error_free(gerr)
		fileDialogCallback(requestID, nil, true)
	case file != 0:
		path := takeGString(g_file_get_path(file))
		g_object_unref(file)
		fileDialogCallback(requestID, []string{path}, false)
	default:
		fileDialogCallback(requestID, nil, true)
	}
}

func finishMultipleFiles(finish func(uintptr, uintptr, uintptr) uintptr, source, res uintptr, requestID uint) {
	var gerr uintptr
	files := finish(source, res, uintptr(unsafe.Pointer(&gerr)))
	switch {
	case gerr != 0:
		g_error_free(gerr)
		fileDialogCallback(requestID, nil, true)
	case files != 0:
		n := g_list_model_get_n_items(files)
		paths := make([]string, 0, n)
		for i := uint32(0); i < n; i++ {
			file := g_list_model_get_item(files, i)
			if file == 0 {
				continue
			}
			if p := g_file_get_path(file); p != 0 {
				paths = append(paths, takeGString(p))
			}
			g_object_unref(file)
		}
		g_object_unref(files)
		fileDialogCallback(requestID, paths, false)
	default:
		fileDialogCallback(requestID, nil, true)
	}
}

var onFileDialogOpenFinishPtr = purego.NewCallback(func(source, res, data uintptr) uintptr {
	finishSingleFile(gtk_file_dialog_open_finish, source, res, uint(data))
	return 0
})

var onFileDialogOpenMultipleFinishPtr = purego.NewCallback(func(source, res, data uintptr) uintptr {
	finishMultipleFiles(gtk_file_dialog_open_multiple_finish, source, res, uint(data))
	return 0
})

var onFileDialogSelectFolderFinishPtr = purego.NewCallback(func(source, res, data uintptr) uintptr {
	finishSingleFile(gtk_file_dialog_select_folder_finish, source, res, uint(data))
	return 0
})

var onFileDialogSelectMultipleFoldersFinishPtr = purego.NewCallback(func(source, res, data uintptr) uintptr {
	finishMultipleFiles(gtk_file_dialog_select_multiple_folders_finish, source, res, uint(data))
	return 0
})

var onFileDialogSaveFinishPtr = purego.NewCallback(func(source, res, data uintptr) uintptr {
	finishSingleFile(gtk_file_dialog_save_finish, source, res, uint(data))
	return 0
})

func showOpenFileDialog(parent, dialog uintptr, requestID uint, allowMultiple, isFolder bool) {
	data := uintptr(requestID)
	switch {
	case isFolder && allowMultiple:
		gtk_file_dialog_select_multiple_folders(dialog, parent, 0, onFileDialogSelectMultipleFoldersFinishPtr, data)
	case isFolder:
		gtk_file_dialog_select_folder(dialog, parent, 0, onFileDialogSelectFolderFinishPtr, data)
	case allowMultiple:
		gtk_file_dialog_open_multiple(dialog, parent, 0, onFileDialogOpenMultipleFinishPtr, data)
	default:
		gtk_file_dialog_open(dialog, parent, 0, onFileDialogOpenFinishPtr, data)
	}
}

func showSaveFileDialog(parent, dialog uintptr, requestID uint) {
	gtk_file_dialog_save(dialog, parent, 0, onFileDialogSaveFinishPtr, uintptr(requestID))
}

func addFileFilter(dialog, filters uintptr, name, pattern string) {
	filter := gtk_file_filter_new()
	gtk_file_filter_set_name(filter, name)
	for _, p := range splitAndTrim(pattern, ";") {
		if p != "" {
			gtk_file_filter_add_pattern(filter, p)
		}
	}
	g_list_store_append(filters, filter)
	g_object_unref(filter)
}

// ----------------------------------------------------------------------------
// Message dialogs (custom GtkWindow-based, port of show_message_dialog)
// ----------------------------------------------------------------------------

type messageDialogState struct {
	dialog       uintptr
	requestID    uint
	cancelButton int
	buttons      []uintptr
}

var (
	messageDialogsLock sync.Mutex
	messageDialogs     = map[uintptr]*messageDialogState{} // keyed by handle id
	messageDialogNext  uintptr
)

func storeMessageDialog(s *messageDialogState) uintptr {
	messageDialogsLock.Lock()
	defer messageDialogsLock.Unlock()
	messageDialogNext++
	messageDialogs[messageDialogNext] = s
	return messageDialogNext
}

func loadMessageDialog(handle uintptr) *messageDialogState {
	messageDialogsLock.Lock()
	defer messageDialogsLock.Unlock()
	return messageDialogs[handle]
}

func dropMessageDialog(handle uintptr) {
	messageDialogsLock.Lock()
	defer messageDialogsLock.Unlock()
	delete(messageDialogs, handle)
}

var onMessageDialogButtonClickedPtr = purego.NewCallback(func(button, data uintptr) uintptr {
	state := loadMessageDialog(data)
	if state == nil {
		return 0
	}
	index := int(g_object_get_data(button, "button-index"))
	alertDialogCallback(state.requestID, index)
	gtk_window_destroy(state.dialog)
	dropMessageDialog(data)
	return 0
})

var onMessageDialogClosePtr = purego.NewCallback(func(window, data uintptr) uintptr {
	state := loadMessageDialog(data)
	if state == nil {
		return 0
	}
	result := -1
	if state.cancelButton >= 0 {
		result = state.cancelButton
	}
	alertDialogCallback(state.requestID, result)
	dropMessageDialog(data)
	return 0 // FALSE: allow the default close handling to destroy the window
})

var onMessageDialogKeyPressedPtr = purego.NewCallback(func(controller uintptr, keyval, keycode uint32, state uint32, data uintptr) uintptr {
	dlgState := loadMessageDialog(data)
	if dlgState == nil {
		return 0
	}
	if keyval == gdkKeyEscape && dlgState.cancelButton >= 0 && dlgState.cancelButton < len(dlgState.buttons) {
		gtk_widget_activate(dlgState.buttons[dlgState.cancelButton])
		return 1
	}
	return 0
})

func showMessageDialog(parent uintptr, heading, body, iconName string, iconData []byte,
	buttons []string, defaultButton, cancelButton, destructiveButton int, requestID uint) {

	dialog := gtk_window_new()
	gtk_window_set_modal(dialog, 1)
	gtk_window_set_resizable(dialog, 0)
	gtk_window_set_decorated(dialog, 1)
	gtk_widget_add_css_class(dialog, "message")
	gtk_widget_set_size_request(dialog, 300, -1)

	if parent != 0 {
		gtk_window_set_transient_for(dialog, parent)
	}

	state := &messageDialogState{
		dialog:       dialog,
		requestID:    requestID,
		cancelButton: cancelButton,
	}
	handle := storeMessageDialog(state)

	content := gtk_box_new(gtkOrientationVertical, 12)
	gtk_widget_set_margin_start(content, 24)
	gtk_widget_set_margin_end(content, 24)
	gtk_widget_set_margin_top(content, 24)
	gtk_widget_set_margin_bottom(content, 24)

	const symbolicIconSize = 32
	var iconWidget uintptr
	if len(iconData) > 0 {
		bytes := g_bytes_new(uintptr(unsafe.Pointer(&iconData[0])), uintptr(len(iconData)))
		texture := gdk_texture_new_from_bytes(bytes, 0)
		g_bytes_unref(bytes)
		if texture != 0 {
			texSize := gdk_texture_get_width(texture)
			image := gtk_image_new_from_paintable(texture)
			gtk_image_set_pixel_size(image, texSize)
			iconWidget = image
			g_object_unref(texture)
		}
	} else if iconName != "" {
		iconWidget = gtk_image_new_from_icon_name(iconName)
		gtk_image_set_pixel_size(iconWidget, symbolicIconSize)
	}

	if iconWidget != 0 {
		gtk_widget_set_halign(iconWidget, gtkAlignCenter)
		gtk_widget_set_margin_bottom(iconWidget, 12)
		gtk_box_append(content, iconWidget)
	}

	if heading != "" {
		headingLabel := gtk_label_new(heading)
		gtk_widget_add_css_class(headingLabel, "title-2")
		gtk_widget_set_halign(headingLabel, gtkAlignCenter)
		gtk_label_set_wrap(headingLabel, 1)
		gtk_label_set_max_width_chars(headingLabel, 50)
		gtk_box_append(content, headingLabel)
	}

	if body != "" {
		bodyLabel := gtk_label_new(body)
		gtk_widget_set_halign(bodyLabel, gtkAlignCenter)
		gtk_label_set_wrap(bodyLabel, 1)
		gtk_label_set_max_width_chars(bodyLabel, 50)
		gtk_widget_add_css_class(bodyLabel, "dim-label")
		gtk_box_append(content, bodyLabel)
	}

	if len(buttons) > 0 {
		buttonBox := gtk_box_new(gtkOrientationHorizontal, 8)
		gtk_widget_set_halign(buttonBox, gtkAlignCenter)
		gtk_widget_set_margin_top(buttonBox, 12)

		for i, label := range buttons {
			btn := gtk_button_new_with_label(label)
			g_object_set_data(btn, "button-index", uintptr(i))
			signalConnect(btn, "clicked", onMessageDialogButtonClickedPtr, handle)
			state.buttons = append(state.buttons, btn)

			if i == defaultButton {
				gtk_widget_add_css_class(btn, "suggested-action")
				gtk_widget_add_css_class(btn, "default")
			}
			if i == destructiveButton {
				gtk_widget_add_css_class(btn, "destructive-action")
			}

			gtk_box_append(buttonBox, btn)
		}

		gtk_box_append(content, buttonBox)
	}

	gtk_window_set_child(dialog, content)

	keyController := gtk_event_controller_key_new()
	signalConnect(keyController, "key-pressed", onMessageDialogKeyPressedPtr, handle)
	gtk_widget_add_controller(dialog, keyController)

	signalConnect(dialog, "close-request", onMessageDialogClosePtr, handle)

	gtk_window_present(dialog)

	if defaultButton >= 0 && defaultButton < len(state.buttons) {
		gtk_window_set_default_widget(dialog, state.buttons[defaultButton])
		gtk_widget_grab_focus(state.buttons[defaultButton])
	}
}

// ----------------------------------------------------------------------------
// Clipboard (GTK4 async API driven to completion on the main context)
// ----------------------------------------------------------------------------

// Unlike the C implementation, which used static globals (racy if two reads
// ever overlapped), each read gets its own state keyed by handle.
type clipboardRead struct {
	done bool
	text string
}

var (
	clipboardReadsLock sync.Mutex
	clipboardReads     = map[uintptr]*clipboardRead{}
	clipboardReadNext  uintptr
)

var onClipboardReadFinishPtr = purego.NewCallback(func(source, res, data uintptr) uintptr {
	var gerr uintptr
	text := gdk_clipboard_read_text_finish(source, res, uintptr(unsafe.Pointer(&gerr)))
	clipboardReadsLock.Lock()
	state := clipboardReads[data]
	if state != nil {
		if gerr != 0 {
			state.text = ""
		} else if text != 0 {
			state.text = goString(text)
		}
		state.done = true
	}
	clipboardReadsLock.Unlock()
	if gerr != 0 {
		g_error_free(gerr)
	}
	if text != 0 {
		g_free(text)
	}
	return 0
})

// clipboardGetTextSync reads the clipboard, iterating the default main
// context until the async result lands (this is called on the main thread,
// mirroring the cgo implementation).
func clipboardGetTextSync() string {
	display := gdk_display_get_default()
	if display == 0 {
		return ""
	}
	clipboard := gdk_display_get_clipboard(display)

	state := &clipboardRead{}
	clipboardReadsLock.Lock()
	clipboardReadNext++
	handle := clipboardReadNext
	clipboardReads[handle] = state
	clipboardReadsLock.Unlock()

	gdk_clipboard_read_text_async(clipboard, 0, onClipboardReadFinishPtr, handle)

	ctx := g_main_context_default()
	for {
		clipboardReadsLock.Lock()
		done := state.done
		clipboardReadsLock.Unlock()
		if done {
			break
		}
		g_main_context_iteration(ctx, 1)
	}

	clipboardReadsLock.Lock()
	delete(clipboardReads, handle)
	clipboardReadsLock.Unlock()
	return state.text
}

// ----------------------------------------------------------------------------
// Window max-size enforcement
// ----------------------------------------------------------------------------

var onWindowSizeChangedPtr = purego.NewCallback(func(object, pspec, data uintptr) uintptr {
	window := object

	// Don't clamp during fullscreen or maximize - these should bypass max size
	// constraints, matching V2 behaviour where geometry hints are suspended.
	if gtk_window_is_fullscreen(window) != 0 || gtk_window_is_maximized(window) != 0 {
		return 0
	}

	maxW := int32(g_object_get_data(window, "wails-max-width"))
	maxH := int32(g_object_get_data(window, "wails-max-height"))
	if maxW <= 0 && maxH <= 0 {
		return 0
	}

	w := gtk_widget_get_width(window)
	h := gtk_widget_get_height(window)

	needsClamp := false
	if maxW > 0 && w > maxW {
		w = maxW
		needsClamp = true
	}
	if maxH > 0 && h > maxH {
		h = maxH
		needsClamp = true
	}
	if needsClamp {
		gtk_window_set_default_size(window, w, h)
	}
	return 0
})

func windowSetMaxSize(window uintptr, maxWidth, maxHeight int) {
	g_object_set_data(window, "wails-max-width", uintptr(maxWidth))
	g_object_set_data(window, "wails-max-height", uintptr(maxHeight))

	if g_object_get_data(window, "wails-max-size-connected") == 0 {
		signalConnect(window, "notify::default-width", onWindowSizeChangedPtr, 0)
		signalConnect(window, "notify::default-height", onWindowSizeChangedPtr, 0)
		g_object_set_data(window, "wails-max-size-connected", 1)
	}
}

// ----------------------------------------------------------------------------
// X11 window helpers (position / always-on-top)
// ----------------------------------------------------------------------------

// Xlib functions are resolved with dlsym(RTLD_DEFAULT): they come from GTK4's
// already-loaded X11 backend, avoiding a hard libX11 dependency. On
// Wayland-only systems they stay nil and every helper is a no-op — matching
// the cgo implementation.
var (
	x11Once               sync.Once
	xMoveWindow           func(uintptr, uintptr, int32, int32) int32
	xFlush                func(uintptr) int32
	xTranslateCoordinates func(uintptr, uintptr, uintptr, int32, int32, uintptr, uintptr, uintptr) int32
	xSendEvent            func(uintptr, uintptr, int32, int64, uintptr) int32
	xInternAtom           func(uintptr, string, int32) uintptr
	xDefaultRootWindow    func(uintptr) uintptr
	gdkX11DisplayType     uintptr
)

func resolveX11Funcs() {
	x11Once.Do(func() {
		reg := func(fptr any, name string) {
			sym, err := purego.Dlsym(purego.RTLD_DEFAULT, name)
			if err == nil && sym != 0 {
				purego.RegisterFunc(fptr, sym)
			}
		}
		reg(&xMoveWindow, "XMoveWindow")
		reg(&xFlush, "XFlush")
		reg(&xTranslateCoordinates, "XTranslateCoordinates")
		reg(&xSendEvent, "XSendEvent")
		reg(&xInternAtom, "XInternAtom")
		reg(&xDefaultRootWindow, "XDefaultRootWindow")
	})
}

// isX11Display reports whether the display is backed by X11 (the purego
// equivalent of GDK_IS_X11_DISPLAY, checked via the GObject type system —
// the GdkX11Display type only registers when the X11 backend is in use).
func isX11Display(display uintptr) bool {
	if gdkX11DisplayType == 0 {
		gdkX11DisplayType = g_type_from_name("GdkX11Display")
	}
	return gdkX11DisplayType != 0 && gTypeInstanceIsA(display, gdkX11DisplayType)
}

func x11WindowForGtkWindow(window uintptr) (xdisplay, xwindow uintptr, ok bool) {
	surface := toplevelForWindow(window)
	if surface == 0 {
		return 0, 0, false
	}
	display := gdk_surface_get_display(surface)
	if !isX11Display(display) {
		return 0, 0, false
	}
	if gdk_x11_display_get_xdisplay == nil || gdk_x11_surface_get_xid == nil {
		return 0, 0, false
	}
	resolveX11Funcs()
	return gdk_x11_display_get_xdisplay(display), gdk_x11_surface_get_xid(surface), true
}

func windowMoveX11(window uintptr, x, y int) {
	xdisplay, xwindow, ok := x11WindowForGtkWindow(window)
	if !ok || xMoveWindow == nil {
		return
	}
	xMoveWindow(xdisplay, xwindow, int32(x), int32(y))
	if xFlush != nil {
		xFlush(xdisplay)
	}
}

func windowGetPositionX11(window uintptr) (int, int) {
	xdisplay, xwindow, ok := x11WindowForGtkWindow(window)
	if !ok || xTranslateCoordinates == nil || xDefaultRootWindow == nil {
		return 0, 0
	}
	root := xDefaultRootWindow(xdisplay)
	var absX, absY int32
	var child uintptr
	if xTranslateCoordinates(xdisplay, xwindow, root, 0, 0,
		uintptr(unsafe.Pointer(&absX)), uintptr(unsafe.Pointer(&absY)),
		uintptr(unsafe.Pointer(&child))) != 0 {
		return int(absX), int(absY)
	}
	return 0, 0
}

const (
	substructureNotifyMask   = 1 << 19
	substructureRedirectMask = 1 << 20
)

// xClientMessageEvent mirrors XEvent's XClientMessage member on 64-bit Linux.
// XEvent itself is a 192-byte union; pad accordingly so Xlib can copy it.
type xClientMessageEvent struct {
	typ         int32
	_           int32
	serial      uint64
	sendEvent   int32
	_           int32
	display     uintptr
	window      uintptr
	messageType uintptr
	format      int32
	_           int32
	dataL       [5]int64
	_           [96]byte // pad to sizeof(XEvent) == 192
}

func windowSendAlwaysOnTopX11(window uintptr, alwaysOnTop bool) {
	xdisplay, xwindow, ok := x11WindowForGtkWindow(window)
	if !ok || xSendEvent == nil || xInternAtom == nil || xDefaultRootWindow == nil {
		return
	}

	netWmState := xInternAtom(xdisplay, "_NET_WM_STATE", 0)
	netWmStateAbove := xInternAtom(xdisplay, "_NET_WM_STATE_ABOVE", 0)
	root := xDefaultRootWindow(xdisplay)

	const clientMessage = 33 // X11 ClientMessage event type
	xev := xClientMessageEvent{
		typ:         clientMessage,
		display:     xdisplay,
		window:      xwindow,
		messageType: netWmState,
		format:      32,
	}
	if alwaysOnTop {
		xev.dataL[0] = 1 // _NET_WM_STATE_ADD
	}
	xev.dataL[1] = int64(netWmStateAbove)
	xev.dataL[3] = 1 // source: normal application

	xSendEvent(xdisplay, root, 0, substructureRedirectMask|substructureNotifyMask,
		uintptr(unsafe.Pointer(&xev)))
	if xFlush != nil {
		xFlush(xdisplay)
	}
}

func windowSetAlwaysOnTop(window uintptr, alwaysOnTop bool) {
	// Store the desired state so windowShow can re-apply it if the surface
	// doesn't exist yet. Use 1=true, 2=false as sentinels (0 means never set).
	sentinel := uintptr(2)
	if alwaysOnTop {
		sentinel = 1
	}
	g_object_set_data(window, "wails-always-on-top", sentinel)
	windowSendAlwaysOnTopX11(window, alwaysOnTop)
}

// windowApplyPendingAlwaysOnTop applies a previously-set always-on-top state
// once the window surface exists. Called from windowShow after present.
func windowApplyPendingAlwaysOnTop(window uintptr) {
	stored := g_object_get_data(window, "wails-always-on-top")
	if stored == 0 {
		return // never been set
	}
	windowSendAlwaysOnTopX11(window, stored == 1)
}
