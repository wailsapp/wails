//go:build linux && purego

package application

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type windowPointer uintptr
type identifier uint
type pointer uintptr

// type GSList uintptr
type GSList struct {
	data pointer
	next *GSList
}

type GSListPointer *GSList

const (
	nilPointer pointer = 0
)

const (
	GSourceRemove int = 0

	// https://gitlab.gnome.org/GNOME/gtk/-/blob/gtk-3-24/gdk/gdkwindow.h#L121
	GdkHintMinSize = 1 << 1
	GdkHintMaxSize = 1 << 2
	// https://gitlab.gnome.org/GNOME/gtk/-/blob/gtk-3-24/gdk/gdkevents.h#L512
	GdkWindowStateIconified  = 1 << 1
	GdkWindowStateMaximized  = 1 << 2
	GdkWindowStateFullscreen = 1 << 4

	// https://gitlab.gnome.org/GNOME/gtk/-/blob/gtk-3-24/gtk/gtkmessagedialog.h#L87
	GtkButtonsNone     int = 0
	GtkButtonsOk           = 1
	GtkButtonsClose        = 2
	GtkButtonsCancel       = 3
	GtkButtonsYesNo        = 4
	GtkButtonsOkCancel     = 5

	// https://gitlab.gnome.org/GNOME/gtk/-/blob/gtk-3-24/gtk/gtkdialog.h#L36
	GtkDialogModal             = 1 << 0
	GtkDialogDestroyWithParent = 1 << 1
	GtkDialogUseHeaderBar      = 1 << 2 // actions in header bar instead of action area

	GtkOrientationVertical = 1

	// enum GtkMessageType
	GtkMessageInfo     = 0
	GtkMessageWarning  = 1
	GtkMessageQuestion = 2
	GtkMessageError    = 3
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

var (
	nilRadioGroup       GSListPointer         = nil
	gtkSignalHandlers   map[pointer]uint      = map[pointer]uint{}
	gtkSignalToMenuItem map[pointer]*MenuItem = map[pointer]*MenuItem{}
	mainThreadId        uint64
)

const (
	// TODO: map distro => so filename - with fallback?
	gtk3    = "libgtk-3.so.0"
	gtk4    = "libgtk-4.so.1"
	webkit4 = "libwebkit2gtk-4.1.so.0"
)

var (
	gtk        uintptr
	gtkVersion int
	webkit     uintptr

	// function references
	gApplicationHold      func(pointer)
	gApplicationQuit      func(pointer)
	gApplicationName      func() string
	gApplicationRelease   func(pointer)
	gApplicationRun       func(pointer, int, []string) int
	gBytesNewStatic       func(uintptr, int) uintptr
	gBytesUnref           func(uintptr)
	gFree                 func(pointer)
	gIdleAdd              func(uintptr)
	gObjectRefSink        func(pointer)
	gObjectUnref          func(pointer)
	gSignalConnectData    func(pointer, string, uintptr, pointer, bool, int) int
	gSignalConnectObject  func(pointer, string, pointer, pointer, int) uint
	gSignalHandlerBlock   func(pointer, uint)
	gSignalHandlerUnblock func(pointer, uint)
	gThreadSelf           func() uint64

	// gdk functions
	gdkDisplayGetMonitor         func(pointer, int) pointer
	gdkDisplayGetMonitorAtWindow func(pointer, pointer) pointer
	gdkDisplayGetNMonitors       func(pointer) int
	gdkMonitorGetGeometry        func(pointer, pointer) pointer
	gdkMonitorGetScaleFactor     func(pointer) int
	gdkMonitorIsPrimary          func(pointer) int
	gdkPixbufNewFromBytes        func(uintptr, int, int, int, int, int, int) pointer
	gdkRgbaParse                 func(pointer, string) bool
	gdkScreenGetRgbaVisual       func(pointer) pointer
	gdkScreenIsComposited        func(pointer) int
	gdkWindowGetState            func(pointer) int
	gdkWindowGetDisplay          func(pointer) pointer

	// gtk functions
	gtkApplicationNew               func(string, uint) pointer
	gtkApplicationGetActiveWindow   func(pointer) pointer
	gtkApplicationGetWindows        func(pointer) *GList
	gtkApplicationWindowNew         func(pointer) pointer
	gtkBoxNew                       func(int, int) pointer
	gtkBoxPackStart                 func(pointer, pointer, int, int, int)
	gtkCheckMenuItemGetActive       func(pointer) int
	gtkCheckMenuItemNewWithLabel    func(string) pointer
	gtkCheckMenuItemSetActive       func(pointer, int)
	gtkContainerAdd                 func(pointer, pointer)
	gtkCSSProviderLoadFromData      func(pointer, string, int, pointer)
	gtkCSSProviderNew               func() pointer
	gtkDialogAddButton              func(pointer, string, int)
	gtkDialogGetContentArea         func(pointer) pointer
	gtkDialogRun                    func(pointer) int
	gtkDialogSetDefaultResponse     func(pointer, int)
	gtkDragDestSet                  func(pointer, uint, pointer, uint, uint)
	gtkFileChooserAddFilter         func(pointer, pointer)
	gtkFileChooserDialogNew         func(string, pointer, int, string, int, string, int, pointer) pointer
	gtkFileChooserGetFilenames      func(pointer) *GSList
	gtkFileChooserSetAction         func(pointer, int)
	gtkFileChooserSetCreateFolders  func(pointer, bool)
	gtkFileChooserSetCurrentFolder  func(pointer, string)
	gtkFileChooserSetSelectMultiple func(pointer, bool)
	gtkFileChooserSetShowHidden     func(pointer, bool)
	gtkFileFilterAddPattern         func(pointer, string)
	gtkFileFilterNew                func() pointer
	gtkFileFilterSetName            func(pointer, string)
	gtkImageNewFromPixbuf           func(pointer) pointer
	gtkMenuBarNew                   func() pointer
	gtkMenuItemNewWithLabel         func(string) pointer
	gtkMenuItemSetLabel             func(pointer, string)
	gtkMenuItemSetSubmenu           func(pointer, pointer)
	gtkMenuNew                      func() pointer
	gtkMenuShellAppend              func(pointer, pointer)
	gtkMessageDialogNew             func(pointer, int, int, int, string) pointer
	gtkRadioMenuItemGetGroup        func(pointer) GSListPointer
	gtkRadioMenuItemNewWithLabel    func(GSListPointer, string) pointer
	gtkSeparatorMenuItemNew         func() pointer
	gtkStyleContextAddProvider      func(pointer, pointer, int)
	gtkTargetEntryFree              func(pointer)
	gtkTargetEntryNew               func(string, int, uint) pointer
	gtkWidgetDestroy                func(pointer)
	gtkWidgetGetDisplay             func(pointer) pointer
	gtkWidgetGetScreen              func(pointer) pointer
	gtkWidgetGetStyleContext        func(pointer) pointer
	gtkWidgetGetWindow              func(pointer) pointer
	gtkWidgetHide                   func(pointer)
	gtkWidgetIsVisible              func(pointer) bool
	gtkWidgetShow                   func(pointer)
	gtkWidgetShowAll                func(pointer)
	gtkWidgetSetAppPaintable        func(pointer, int)
	gtkWidgetSetName                func(pointer, string)
	gtkWidgetSetSensitive           func(pointer, int)
	gtkWidgetSetToolTipText         func(pointer, string)
	gtkWidgetSetVisual              func(pointer, pointer)
	gtkWindowClose                  func(pointer)
	gtkWindowFullScreen             func(pointer)
	gtkWindowGetPosition            func(pointer, *int, *int) bool
	gtkWindowGetSize                func(pointer, *int, *int)
	gtkWindowHasToplevelFocus       func(pointer) int
	gtkWindowKeepAbove              func(pointer, bool)
	gtkWindowMaximize               func(pointer)
	gtkWindowMinimize               func(pointer)
	gtkWindowMove                   func(pointer, int, int)
	gtkWindowPresent                func(pointer)
	gtkWindowResize                 func(pointer, int, int)
	gtkWindowSetDecorated           func(pointer, int)
	gtkWindowSetGeometryHints       func(pointer, pointer, pointer, int)
	gtkWindowSetKeepAbove           func(pointer, bool)
	gtkWindowSetResizable           func(pointer, bool)
	gtkWindowSetTitle               func(pointer, string)
	gtkWindowUnfullscreen           func(pointer)
	gtkWindowUnmaximize             func(pointer)

	// webkit
	webkitNewWithUserContentManager                      func(pointer) pointer
	webkitRegisterUriScheme                              func(pointer, string, pointer, int, int)
	webkitSettingsGetEnableDeveloperExtras               func(pointer) bool
	webkitSettingsSetHardwareAccelerationPolicy          func(pointer, int)
	webkitSettingsSetEnableDeveloperExtras               func(pointer, bool)
	webkitSettingsSetUserAgentWithApplicationDetails     func(pointer, string, string)
	webkitUserContentManagerNew                          func() pointer
	webkitUserContentManagerRegisterScriptMessageHandler func(pointer, string)
	webkitWebContextGetDefault                           func() pointer
	webkitWebViewEvaluateJS                              func(pointer, string, int, pointer, string, pointer, pointer, pointer)
	webkitWebViewGetSettings                             func(pointer) pointer
	webkitWebViewGetZoom                                 func(pointer) float64
	webkitWebViewLoadAlternateHTML                       func(pointer, string, string, *string)
	webkitWebViewLoadUri                                 func(pointer, string)
	webkitWebViewSetBackgroundColor                      func(pointer, pointer)
	webkitWebViewSetSettings                             func(pointer, pointer)
	webkitWebViewSetZoomLevel                            func(pointer, float64)
)

func init() {
	// needed for GTK4 to function
	_ = os.Setenv("GDK_BACKEND", "x11")
	var err error

	// gtk, err = purego.Dlopen(gtk4, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	// if err == nil {
	// 	version = 4
	// 	return
	// }

	// log.Println("Failed to open GTK4: Falling back to GTK3")

	gtk, err = purego.Dlopen(gtk3, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	gtkVersion = 3

	webkit, err = purego.Dlopen(webkit4, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	// Function registration
	// GLib
	purego.RegisterLibFunc(&gApplicationHold, gtk, "g_application_hold")
	purego.RegisterLibFunc(&gApplicationName, gtk, "g_get_application_name")
	purego.RegisterLibFunc(&gApplicationQuit, gtk, "g_application_quit")
	purego.RegisterLibFunc(&gApplicationRelease, gtk, "g_application_release")
	purego.RegisterLibFunc(&gApplicationRun, gtk, "g_application_run")
	purego.RegisterLibFunc(&gBytesNewStatic, gtk, "g_bytes_new_static")
	purego.RegisterLibFunc(&gBytesUnref, gtk, "g_bytes_unref")
	purego.RegisterLibFunc(&gFree, gtk, "g_free")
	purego.RegisterLibFunc(&gIdleAdd, gtk, "g_idle_add")
	purego.RegisterLibFunc(&gObjectRefSink, gtk, "g_object_ref_sink")
	purego.RegisterLibFunc(&gObjectUnref, gtk, "g_object_unref")
	purego.RegisterLibFunc(&gSignalConnectData, gtk, "g_signal_connect_data")
	purego.RegisterLibFunc(&gSignalConnectObject, gtk, "g_signal_connect_object")
	purego.RegisterLibFunc(&gSignalHandlerBlock, gtk, "g_signal_handler_block")
	purego.RegisterLibFunc(&gSignalHandlerUnblock, gtk, "g_signal_handler_unblock")
	purego.RegisterLibFunc(&gThreadSelf, gtk, "g_thread_self")

	// GDK
	purego.RegisterLibFunc(&gdkDisplayGetMonitor, gtk, "gdk_display_get_monitor")
	purego.RegisterLibFunc(&gdkDisplayGetMonitorAtWindow, gtk, "gdk_display_get_monitor_at_window")
	purego.RegisterLibFunc(&gdkDisplayGetNMonitors, gtk, "gdk_display_get_n_monitors")
	purego.RegisterLibFunc(&gdkMonitorGetGeometry, gtk, "gdk_monitor_get_geometry")
	purego.RegisterLibFunc(&gdkMonitorGetScaleFactor, gtk, "gdk_monitor_get_scale_factor")
	purego.RegisterLibFunc(&gdkMonitorIsPrimary, gtk, "gdk_monitor_is_primary")
	purego.RegisterLibFunc(&gdkPixbufNewFromBytes, gtk, "gdk_pixbuf_new_from_bytes")
	purego.RegisterLibFunc(&gdkRgbaParse, gtk, "gdk_rgba_parse")
	purego.RegisterLibFunc(&gdkScreenGetRgbaVisual, gtk, "gdk_screen_get_rgba_visual")
	purego.RegisterLibFunc(&gdkScreenIsComposited, gtk, "gdk_screen_is_composited")
	purego.RegisterLibFunc(&gdkWindowGetDisplay, gtk, "gdk_window_get_display")
	purego.RegisterLibFunc(&gdkWindowGetState, gtk, "gdk_window_get_state")

	// GTK3
	purego.RegisterLibFunc(&gtkApplicationNew, gtk, "gtk_application_new")
	purego.RegisterLibFunc(&gtkApplicationGetActiveWindow, gtk, "gtk_application_get_active_window")
	purego.RegisterLibFunc(&gtkApplicationGetWindows, gtk, "gtk_application_get_windows")
	purego.RegisterLibFunc(&gtkApplicationWindowNew, gtk, "gtk_application_window_new")
	purego.RegisterLibFunc(&gtkBoxNew, gtk, "gtk_box_new")
	purego.RegisterLibFunc(&gtkBoxPackStart, gtk, "gtk_box_pack_start")
	purego.RegisterLibFunc(&gtkCheckMenuItemGetActive, gtk, "gtk_check_menu_item_get_active")
	purego.RegisterLibFunc(&gtkCheckMenuItemNewWithLabel, gtk, "gtk_check_menu_item_new_with_label")
	purego.RegisterLibFunc(&gtkCheckMenuItemSetActive, gtk, "gtk_check_menu_item_set_active")
	purego.RegisterLibFunc(&gtkContainerAdd, gtk, "gtk_container_add")
	purego.RegisterLibFunc(&gtkCSSProviderLoadFromData, gtk, "gtk_css_provider_load_from_data")
	purego.RegisterLibFunc(&gtkDialogAddButton, gtk, "gtk_dialog_add_button")
	purego.RegisterLibFunc(&gtkDialogGetContentArea, gtk, "gtk_dialog_get_content_area")
	purego.RegisterLibFunc(&gtkDialogRun, gtk, "gtk_dialog_run")
	purego.RegisterLibFunc(&gtkDialogSetDefaultResponse, gtk, "gtk_dialog_set_default_response")
	purego.RegisterLibFunc(&gtkDragDestSet, gtk, "gtk_drag_dest_set")
	purego.RegisterLibFunc(&gtkFileChooserAddFilter, gtk, "gtk_file_chooser_add_filter")
	purego.RegisterLibFunc(&gtkFileChooserDialogNew, gtk, "gtk_file_chooser_dialog_new")
	purego.RegisterLibFunc(&gtkFileChooserGetFilenames, gtk, "gtk_file_chooser_get_filenames")
	purego.RegisterLibFunc(&gtkFileChooserSetAction, gtk, "gtk_file_chooser_set_action")
	purego.RegisterLibFunc(&gtkFileChooserSetCreateFolders, gtk, "gtk_file_chooser_set_create_folders")
	purego.RegisterLibFunc(&gtkFileChooserSetCurrentFolder, gtk, "gtk_file_chooser_set_current_folder")
	purego.RegisterLibFunc(&gtkFileChooserSetSelectMultiple, gtk, "gtk_file_chooser_set_select_multiple")
	purego.RegisterLibFunc(&gtkFileChooserSetShowHidden, gtk, "gtk_file_chooser_set_show_hidden")
	purego.RegisterLibFunc(&gtkFileFilterAddPattern, gtk, "gtk_file_filter_add_pattern")
	purego.RegisterLibFunc(&gtkFileFilterNew, gtk, "gtk_file_filter_new")
	purego.RegisterLibFunc(&gtkFileFilterSetName, gtk, "gtk_file_filter_set_name")
	purego.RegisterLibFunc(&gtkImageNewFromPixbuf, gtk, "gtk_image_new_from_pixbuf")
	purego.RegisterLibFunc(&gtkMenuItemSetLabel, gtk, "gtk_menu_item_set_label")
	purego.RegisterLibFunc(&gtkMenuBarNew, gtk, "gtk_menu_bar_new")
	purego.RegisterLibFunc(&gtkMenuItemNewWithLabel, gtk, "gtk_menu_item_new_with_label")
	purego.RegisterLibFunc(&gtkMenuItemSetSubmenu, gtk, "gtk_menu_item_set_submenu")
	purego.RegisterLibFunc(&gtkMenuNew, gtk, "gtk_menu_new")
	purego.RegisterLibFunc(&gtkMenuShellAppend, gtk, "gtk_menu_shell_append")
	purego.RegisterLibFunc(&gtkMessageDialogNew, gtk, "gtk_message_dialog_new")
	purego.RegisterLibFunc(&gtkRadioMenuItemGetGroup, gtk, "gtk_radio_menu_item_get_group")
	purego.RegisterLibFunc(&gtkRadioMenuItemNewWithLabel, gtk, "gtk_radio_menu_item_new_with_label")
	purego.RegisterLibFunc(&gtkSeparatorMenuItemNew, gtk, "gtk_separator_menu_item_new")
	purego.RegisterLibFunc(&gtkStyleContextAddProvider, gtk, "gtk_style_context_add_provider")
	purego.RegisterLibFunc(&gtkTargetEntryFree, gtk, "gtk_target_entry_free")
	purego.RegisterLibFunc(&gtkTargetEntryNew, gtk, "gtk_target_entry_new")
	purego.RegisterLibFunc(&gtkWidgetDestroy, gtk, "gtk_widget_destroy")
	purego.RegisterLibFunc(&gtkWidgetGetDisplay, gtk, "gtk_widget_get_display")
	purego.RegisterLibFunc(&gtkWidgetGetScreen, gtk, "gtk_widget_get_screen")
	purego.RegisterLibFunc(&gtkWidgetGetStyleContext, gtk, "gtk_widget_get_style_context")
	purego.RegisterLibFunc(&gtkWidgetGetWindow, gtk, "gtk_widget_get_window")
	purego.RegisterLibFunc(&gtkWidgetHide, gtk, "gtk_widget_hide")
	purego.RegisterLibFunc(&gtkWidgetIsVisible, gtk, "gtk_widget_is_visible")
	purego.RegisterLibFunc(&gtkWidgetSetAppPaintable, gtk, "gtk_widget_set_app_paintable")
	purego.RegisterLibFunc(&gtkWidgetSetName, gtk, "gtk_widget_set_name")
	purego.RegisterLibFunc(&gtkWidgetSetSensitive, gtk, "gtk_widget_set_sensitive")
	purego.RegisterLibFunc(&gtkWidgetSetToolTipText, gtk, "gtk_widget_set_tooltip_text")
	purego.RegisterLibFunc(&gtkWidgetSetVisual, gtk, "gtk_widget_set_visual")
	purego.RegisterLibFunc(&gtkWidgetShow, gtk, "gtk_widget_show")
	purego.RegisterLibFunc(&gtkWidgetShowAll, gtk, "gtk_widget_show_all")
	purego.RegisterLibFunc(&gtkWindowFullScreen, gtk, "gtk_window_fullscreen")
	purego.RegisterLibFunc(&gtkWindowClose, gtk, "gtk_window_close")
	purego.RegisterLibFunc(&gtkWindowGetPosition, gtk, "gtk_window_get_position")
	purego.RegisterLibFunc(&gtkWindowGetSize, gtk, "gtk_window_get_size")
	purego.RegisterLibFunc(&gtkWindowMaximize, gtk, "gtk_window_maximize")
	purego.RegisterLibFunc(&gtkWindowMove, gtk, "gtk_window_move")
	purego.RegisterLibFunc(&gtkWindowPresent, gtk, "gtk_window_present")
	//purego.RegisterLibFunc(&gtkWindowPresent, gtk, "gtk_window_unminimize")  // gtk4
	purego.RegisterLibFunc(&gtkWindowHasToplevelFocus, gtk, "gtk_window_has_toplevel_focus")
	purego.RegisterLibFunc(&gtkWindowMinimize, gtk, "gtk_window_iconify") // gtk3
	//	purego.RegisterLibFunc(&gtkWindowMinimize, gtk, "gtk_window_minimize") // gtk4
	purego.RegisterLibFunc(&gtkWindowResize, gtk, "gtk_window_resize")
	purego.RegisterLibFunc(&gtkWindowSetGeometryHints, gtk, "gtk_window_set_geometry_hints")
	purego.RegisterLibFunc(&gtkWindowSetDecorated, gtk, "gtk_window_set_decorated")
	purego.RegisterLibFunc(&gtkWindowKeepAbove, gtk, "gtk_window_set_keep_above")
	purego.RegisterLibFunc(&gtkWindowSetResizable, gtk, "gtk_window_set_resizable")
	purego.RegisterLibFunc(&gtkWindowSetTitle, gtk, "gtk_window_set_title")
	purego.RegisterLibFunc(&gtkWindowUnfullscreen, gtk, "gtk_window_unfullscreen")
	purego.RegisterLibFunc(&gtkWindowUnmaximize, gtk, "gtk_window_unmaximize")

	// webkit
	purego.RegisterLibFunc(&webkitNewWithUserContentManager, webkit, "webkit_web_view_new_with_user_content_manager")
	purego.RegisterLibFunc(&webkitRegisterUriScheme, webkit, "webkit_web_context_register_uri_scheme")
	purego.RegisterLibFunc(&webkitSettingsGetEnableDeveloperExtras, webkit, "webkit_settings_get_enable_developer_extras")
	purego.RegisterLibFunc(&webkitSettingsSetEnableDeveloperExtras, webkit, "webkit_settings_set_enable_developer_extras")
	purego.RegisterLibFunc(&webkitSettingsSetHardwareAccelerationPolicy, webkit, "webkit_settings_set_hardware_acceleration_policy")
	purego.RegisterLibFunc(&webkitSettingsSetUserAgentWithApplicationDetails, webkit, "webkit_settings_set_user_agent_with_application_details")
	purego.RegisterLibFunc(&webkitUserContentManagerNew, webkit, "webkit_user_content_manager_new")
	purego.RegisterLibFunc(&webkitUserContentManagerRegisterScriptMessageHandler, webkit, "webkit_user_content_manager_register_script_message_handler")
	purego.RegisterLibFunc(&webkitWebContextGetDefault, webkit, "webkit_web_context_get_default")
	purego.RegisterLibFunc(&webkitWebViewEvaluateJS, webkit, "webkit_web_view_evaluate_javascript")
	purego.RegisterLibFunc(&webkitWebViewGetSettings, webkit, "webkit_web_view_get_settings")
	purego.RegisterLibFunc(&webkitWebViewGetZoom, webkit, "webkit_web_view_get_zoom_level")
	purego.RegisterLibFunc(&webkitWebViewLoadAlternateHTML, webkit, "webkit_web_view_load_alternate_html")
	purego.RegisterLibFunc(&webkitWebViewLoadUri, webkit, "webkit_web_view_load_uri")
	purego.RegisterLibFunc(&webkitWebViewSetBackgroundColor, webkit, "webkit_web_view_set_background_color")
	purego.RegisterLibFunc(&webkitWebViewSetSettings, webkit, "webkit_web_view_set_settings")
	purego.RegisterLibFunc(&webkitWebViewSetZoomLevel, webkit, "webkit_web_view_set_zoom_level")
}

// mainthread stuff
func dispatchOnMainThread(id uint) {
	gIdleAdd(purego.NewCallback(func(pointer) int {
		executeOnMainThread(id)
		return GSourceRemove
	}))
}

// implementation below
func appName() string {
	return gApplicationName()
}

func appNew(name string) pointer {
	GApplicationDefaultFlags := uint(0)

	name = strings.ToLower(name)
	if name == "" {
		name = "undefined"
	}
	identifier := fmt.Sprintf("org.wails.%s", strings.Replace(name, " ", "-", -1))

	return pointer(gtkApplicationNew(identifier, GApplicationDefaultFlags))
}

func appRun(application pointer) error {
	mainThreadId = gThreadSelf()
	fmt.Println("linux_purego: appRun threadID", mainThreadId)

	app := pointer(application)
	activate := func() {
		// TODO: Do we care?
		fmt.Println("linux.activated!")
		gApplicationHold(app) // allow running without a window
	}
	gSignalConnectData(
		application,
		"activate",
		purego.NewCallback(activate),
		app,
		false,
		0)

	status := gApplicationRun(app, 0, nil)
	gApplicationRelease(app)
	gObjectUnref(app)

	var err error
	if status != 0 {
		err = fmt.Errorf("exit code: %d", status)
	}
	return err
}

func appDestroy(application pointer) {
	gApplicationQuit(pointer(application))
}

func getCurrentWindowID(application pointer, windows map[windowPointer]uint) uint {
	// TODO: Add extra metadata to window and use it!
	window := gtkApplicationGetActiveWindow(pointer(application))
	identifier, ok := windows[windowPointer(window)]
	if ok {
		return identifier
	}
	// FIXME: Should we panic here if not found?
	return uint(1)
}

type GList struct {
	data pointer
	next *GList
	prev *GList
}

func getWindows(application pointer) []pointer {
	result := []pointer{}
	windows := gtkApplicationGetWindows(pointer(application))
	// FIXME: Need to make a struct here to deal with response data
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
		gtkWidgetHide(window)
	}
}

func showAllWindows(application pointer) {
	for _, window := range getWindows(application) {
		gtkWidgetShowAll(window)
	}
}

// Menu
func menuAddSeparator(menu *Menu) {
	gtkMenuShellAppend(
		pointer((menu.impl).(*linuxMenu).native),
		gtkSeparatorMenuItemNew())
}

func menuAppend(parent *Menu, menu *MenuItem) {
	// TODO: override this with the GTK4 version if needed - possibly rename to imply it's an alias
	gtkMenuShellAppend(
		pointer((parent.impl).(*linuxMenu).native),
		pointer((menu.impl).(*linuxMenuItem).native))

	/* gtk4
	C.gtk_menu_item_set_submenu(
		(*C.struct__GtkMenuItem)((menu.impl).(*linuxMenuItem).native),
		(*C.struct__GtkWidget)((parent.impl).(*linuxMenu).native),
	)
	*/
}

func menuBarNew() pointer {
	return pointer(gtkMenuBarNew())
}

func menuNew() pointer {
	return pointer(gtkMenuNew())
}

func menuSetSubmenu(item *MenuItem, menu *Menu) {
	// FIXME: How is this different than `menuAppend` above?
	gtkMenuItemSetSubmenu(
		pointer((item.impl).(*linuxMenuItem).native),
		pointer((menu.impl).(*linuxMenu).native))
}

func menuGetRadioGroup(item *linuxMenuItem) *GSList {
	return (*GSList)(gtkRadioMenuItemGetGroup(pointer(item.native)))
}

func attachMenuHandler(item *MenuItem) {
	handleClick := func() {
		item := item
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

	impl := (item.impl).(*linuxMenuItem)
	widget := impl.native
	flags := 0
	handlerId := gSignalConnectObject(
		pointer(widget),
		"activate",
		pointer(purego.NewCallback(handleClick)),
		pointer(widget),
		flags)
	impl.handlerId = uint(handlerId)
}

// menuItem
func menuItemChecked(widget pointer) bool {
	if gtkCheckMenuItemGetActive(widget) == 1 {
		return true
	}
	return false
}

func menuItemNew(label string) pointer {
	return pointer(gtkMenuItemNewWithLabel(label))
}

func menuCheckItemNew(label string) pointer {
	return pointer(gtkCheckMenuItemNewWithLabel(label))
}

func menuItemSetChecked(widget pointer, checked bool) {
	value := 0
	if checked {
		value = 1
	}
	gtkCheckMenuItemSetActive(pointer(widget), value)
}

func menuItemSetDisabled(widget pointer, disabled bool) {
	value := 1
	if disabled {
		value = 0
	}
	gtkWidgetSetSensitive(widget, value)
}

func menuItemSetLabel(widget pointer, label string) {
	gtkMenuItemSetLabel(
		pointer(widget),
		label)
}

func menuItemSetToolTip(widget pointer, tooltip string) {
	gtkWidgetSetToolTipText(
		pointer(widget),
		tooltip)
}

func menuItemSignalBlock(widget pointer, handlerId uint, block bool) {
	if block {
		gSignalHandlerBlock(widget, handlerId)
	} else {
		gSignalHandlerUnblock(widget, handlerId)
	}
}

func menuRadioItemNew(group GSListPointer, label string) pointer {
	return pointer(gtkRadioMenuItemNewWithLabel(group, label))
}

// screen related

func getScreenByIndex(display pointer, index int) *Screen {
	monitor := gdkDisplayGetMonitor(display, index)
	// TODO: Do we need to update Screen to contain current info?
	//	currentMonitor := C.gdk_display_get_monitor_at_window(display, window)

	geometry := struct {
		x      int32
		y      int32
		width  int32
		height int32
	}{}
	result := pointer(unsafe.Pointer(&geometry))
	gdkMonitorGetGeometry(monitor, result)

	primary := false
	if gdkMonitorIsPrimary(monitor) == 1 {
		primary = true
	}

	return &Screen{
		IsPrimary:   primary,
		ScaleFactor: 1.0,
		X:           int(geometry.x),
		Y:           int(geometry.y),
		Size: Size{
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
	}
}

func getScreens(app pointer) ([]*Screen, error) {
	var screens []*Screen
	window := gtkApplicationGetActiveWindow(app)
	display := gdkWindowGetDisplay(window)
	count := gdkDisplayGetNMonitors(display)
	for i := 0; i < int(count); i++ {
		screens = append(screens, getScreenByIndex(display, i))
	}
	return screens, nil
}

// widgets
func widgetSetSensitive(widget pointer, enabled bool) {
	value := 0
	if enabled {
		value = 1
	}
	gtkWidgetSetSensitive(widget, value)
}

func widgetSetVisible(widget pointer, hidden bool) {
	if hidden {
		gtkWidgetHide(widget)
	} else {
		gtkWidgetShow(widget)
	}
}

// window related functions
func windowClose(window pointer) {
	gtkWindowClose(window)
}

func windowEnableDND(id uint, webview pointer) {
	targetentry := gtkTargetEntryNew("text/uri-list", 0, id)
	defer gtkTargetEntryFree(targetentry)

	GtkDestDefaultDrop := uint(0)
	GdkActionCopy := uint(0) //?
	gtkDragDestSet(webview, GtkDestDefaultDrop, targetentry, 1, GdkActionCopy)

	// FIXME: enable and process
	/*	gSignalConnectData(webview,
		"drag-data-received",
		purego.NewCallback(onDragNDrop),
		0,
		false,
		0)*/
}

func windowExecJS(webview pointer, js string) {
	webkitWebViewEvaluateJS(
		webview,
		js,
		len(js),
		0,
		"",
		0,
		0,
		0)
}

func windowDestroy(window pointer) {
	// Should this truly 'destroy' ?
	gtkWindowClose(window)
}

func windowFullscreen(window pointer) {
	gtkWindowFullScreen(window)
}

func windowGetPosition(window pointer) (int, int) {
	var x, y int
	gtkWindowGetPosition(window, &x, &y)
	return x, y
}

func windowGetCurrentMonitor(window pointer) pointer {
	// Get the monitor that the window is currently on
	display := gtkWidgetGetDisplay(window)
	window = gtkWidgetGetWindow(window)
	if window == 0 {
		return 0
	}
	return gdkDisplayGetMonitorAtWindow(display, window)
}

func windowGetCurrentMonitorGeometry(window pointer) (x int, y int, width int, height int, scaleFactor int) {
	monitor := windowGetCurrentMonitor(window)
	if monitor == 0 {
		return -1, -1, -1, -1, 1
	}

	result := struct {
		x      int32
		y      int32
		width  int32
		height int32
	}{}
	gdkMonitorGetGeometry(monitor, pointer(unsafe.Pointer(&result)))
	return int(result.x), int(result.y), int(result.width), int(result.height), gdkMonitorGetScaleFactor(monitor)
}

func windowGetRelativePosition(window pointer) (int, int) {
	absX, absY := windowGetPosition(window)
	x, y, _, _, _ := windowGetCurrentMonitorGeometry(window)

	relX := absX - x
	relY := absY - y

	// TODO: Scale based on DPI
	return relX, relY
}

func windowGetSize(window pointer) (int, int) {
	// TODO: dispatchOnMainThread?
	var width, height int
	gtkWindowGetSize(window, &width, &height)
	return width, height
}

func windowGetPosition(window pointer) (int, int) {
	// TODO: dispatchOnMainThread?
	var x, y int
	gtkWindowGetPosition(window, &x, &y)
	return x, y
}

func windowHide(window pointer) {
	gtkWidgetHide(window)
}

func windowIsFocused(window pointer) bool {
	return gtkWindowHasToplevelFocus(window) == 1
}

func windowIsFullscreen(window pointer) bool {
	gdkwindow := gtkWidgetGetWindow(window)
	state := gdkWindowGetState(gdkwindow)
	return state&GdkWindowStateFullscreen > 0
}

func windowIsMaximized(window pointer) bool {
	gdkwindow := gtkWidgetGetWindow(window)
	state := gdkWindowGetState(gdkwindow)
	return state&GdkWindowStateMaximized > 0 && state&GdkWindowStateFullscreen == 0
}

func windowIsMinimized(window pointer) bool {
	gdkwindow := gtkWidgetGetWindow(window)
	state := gdkWindowGetState(gdkwindow)
	return state&GdkWindowStateIconified > 0
}

func windowIsVisible(window pointer) bool {
	// TODO: validate this works..  (used a `bool` in the registration)
	return gtkWidgetIsVisible(window)
}

func windowMaximize(window pointer) {
	gtkWindowMaximize(window)
}

func windowMinimize(window pointer) {
	gtkWindowMinimize(window)
}

func windowNew(application pointer, menu pointer, windowId uint, gpuPolicy int) (pointer, pointer, pointer) {
	window := gtkApplicationWindowNew(application)
	gObjectRefSink(window)
	webview := windowNewWebview(windowId, gpuPolicy)
	vbox := gtkBoxNew(GtkOrientationVertical, 0)
	gtkContainerAdd(window, vbox)
	gtkWidgetSetName(vbox, "webview-box")

	if menu != 0 {
		gtkBoxPackStart(vbox, menu, 0, 0, 0)
	}
	gtkBoxPackStart(vbox, webview, 1, 1, 0)
	return pointer(window), pointer(webview), pointer(vbox)
}

func windowNewWebview(parentId uint, gpuPolicy int) pointer {
	manager := webkitUserContentManagerNew()
	webkitUserContentManagerRegisterScriptMessageHandler(manager, "external")
	wv := webkitNewWithUserContentManager(manager)
	if !registered {
		webkitRegisterUriScheme(
			webkitWebContextGetDefault(),
			"wails",
			pointer(purego.NewCallback(func(request uintptr) {
				webviewRequests <- &webViewAssetRequest{
					Request:    webview.NewRequest(request),
					windowId:   parentId,
					windowName: globalApplication.getWindowForID(parentId).Name(),
				}
			})),
			0,
			0,
		)
		registered = true
	}

	settings := webkitWebViewGetSettings(wv)
	webkitSettingsSetUserAgentWithApplicationDetails(
		settings,
		"wails.io",
		"")
	webkitSettingsSetHardwareAccelerationPolicy(settings, gpuPolicy)
	webkitWebViewSetSettings(wv, settings)
	return wv
}

func windowPresent(window pointer) {
	gtkWindowPresent(pointer(window))
}

func windowReload(webview pointer, address string) {
	webkitWebViewLoadUri(pointer(webview), address)
}

func windowResize(window pointer, width, height int) {
	gtkWindowResize(window, width, height)
}

func windowShow(window pointer) {
	gtkWidgetShowAll(pointer(window))
}

func windowSetBackgroundColour(vbox, webview pointer, colour RGBA) {
	const GtkStyleProviderPriorityUser = 800

	// FIXME: Use a struct!
	rgba := make([]byte, 4*8) // C.sizeof_GdkRGBA == 32
	rgbaPointer := pointer(unsafe.Pointer(&rgba[0]))
	if !gdkRgbaParse(
		rgbaPointer,
		fmt.Sprintf("rgba(%v,%v,%v,%v)",
			colour.Red,
			colour.Green,
			colour.Blue,
			float32(colour.Alpha)/255.0,
		)) {
		return
	}
	webkitWebViewSetBackgroundColor(pointer(webview), rgbaPointer)

	colour.Alpha = 255
	css := fmt.Sprintf("#webview-box {background-color: rgba(%d, %d, %d, %1.1f);}", colour.Red, colour.Green, colour.Blue, float32(colour.Alpha)/255.0)
	provider := gtkCSSProviderNew()
	defer gObjectUnref(provider)
	gtkStyleContextAddProvider(
		gtkWidgetGetStyleContext(vbox),
		provider,
		GtkStyleProviderPriorityUser,
	)
	gtkCSSProviderLoadFromData(provider, css, -1, 0)
}

func windowSetGeometryHints(window pointer, minWidth, minHeight, maxWidth, maxHeight int) {
	size := GdkGeometry{
		minWidth:  int32(minWidth),
		minHeight: int32(minHeight),
		maxWidth:  int32(maxWidth),
		maxHeight: int32(maxHeight),
	}
	gtkWindowSetGeometryHints(
		pointer(window),
		pointer(0),
		pointer(unsafe.Pointer(&size)),
		GdkHintMinSize|GdkHintMaxSize)
}

func windowSetFrameless(window pointer, frameless bool) {
	decorated := 1
	if frameless {
		decorated = 0
	}
	gtkWindowSetDecorated(pointer(window), decorated)

	// TODO: Deal with transparency for the titlebar if possible when !frameless
	//       Perhaps we just make it undecorated and add a menu bar inside?
}

// TODO: confirm this is working properly
func windowSetHTML(webview pointer, html string) {
	webkitWebViewLoadAlternateHTML(webview, html, "wails://", nil)
}

func windowSetKeepAbove(window pointer, alwaysOnTop bool) {
	gtkWindowKeepAbove(window, alwaysOnTop)
}

func windowSetResizable(window pointer, resizable bool) {
	// FIXME: Does this work?
	gtkWindowSetResizable(
		pointer(window),
		resizable,
	)
}

func windowSetTitle(window pointer, title string) {
	gtkWindowSetTitle(pointer(window), title)
}

func windowSetTransparent(window pointer) {
	screen := gtkWidgetGetScreen(pointer(window))
	visual := gdkScreenGetRgbaVisual(screen)
	if visual == 0 {
		return
	}
	if gdkScreenIsComposited(screen) == 1 {
		gtkWidgetSetAppPaintable(pointer(window), 1)
		gtkWidgetSetVisual(pointer(window), visual)
	}
}

func windowSetURL(webview pointer, uri string) {
	webkitWebViewLoadUri(webview, uri)
}

func windowSetupSignalHandlers(windowId uint, window, webview pointer, emit func(e events.WindowEventType)) {
	handleDelete := purego.NewCallback(func(pointer) {
		emit(events.Common.WindowClosing)
	})
	gSignalConnectData(window, "delete-event", handleDelete, 0, false, 0)

	/*
		event = C.CString("load-changed")
		defer C.free(unsafe.Pointer(event))
		C.signal_connect(webview, event, C.webviewLoadChanged, unsafe.Pointer(&w.parent.id))
	*/

	// TODO: Handle mouse button / drag events
	/*	id := C.uint(windowId)
		event = C.CString("button-press-event")
		C.signal_connect((*C.GtkWidget)(unsafe.Pointer(webview)), event, C.onButtonEvent, unsafe.Pointer(&id))
		C.free(unsafe.Pointer(event))
		event = C.CString("button-release-event")
		defer C.free(unsafe.Pointer(event))
		C.signal_connect((*C.GtkWidget)(unsafe.Pointer(webview)), event, onButtonEvent, unsafe.Pointer(&id))
	*/
}

func windowOpenDevTools(webview pointer) {
	settings := webkitWebViewGetSettings(pointer(webview))
	webkitSettingsSetEnableDeveloperExtras(
		settings,
		!webkitSettingsGetEnableDeveloperExtras(settings))
}

func windowUnfullscreen(window pointer) {
	gtkWindowUnfullscreen(window)
}

func windowUnmaximize(window pointer) {
	gtkWindowUnmaximize(window)
}

func windowZoom(webview pointer) float64 {
	return webkitWebViewGetZoom(webview)
}

func windowZoomIn(webview pointer) {
	ZoomInFactor := 1.10
	windowZoomSet(webview, windowZoom(webview)*ZoomInFactor)
}
func windowZoomOut(webview pointer) {
	ZoomOutFactor := -1.10
	windowZoomSet(webview, windowZoom(webview)*ZoomOutFactor)
}

func windowZoomSet(webview pointer, zoom float64) {
	if zoom < 1.0 { // 1.0 is the smallest allowable
		zoom = 1.0
	}
	webkitWebViewSetZoomLevel(webview, zoom)
}

func windowMove(window pointer, x, y int) {
	gtkWindowMove(window, x, y)
}

func runChooserDialog(window pointer, allowMultiple, createFolders, showHidden bool, currentFolder, title string, action int, acceptLabel string, filters []FileFilter) ([]string, error) {
	GtkResponseCancel := 0
	GtkResponseAccept := 1

	fc := gtkFileChooserDialogNew(
		title,
		window,
		action,
		"_Cancel",
		GtkResponseCancel,
		acceptLabel,
		GtkResponseAccept,
		0)

	gtkFileChooserSetAction(fc, action)

	gtkFilters := []pointer{}
	for _, filter := range filters {
		f := gtkFileFilterNew()
		gtkFileFilterSetName(f, filter.DisplayName)
		gtkFileFilterAddPattern(f, filter.Pattern)
		gtkFileChooserAddFilter(fc, f)
		gtkFilters = append(gtkFilters, f)
	}
	gtkFileChooserSetSelectMultiple(fc, allowMultiple)
	gtkFileChooserSetCreateFolders(fc, createFolders)
	gtkFileChooserSetShowHidden(fc, showHidden)

	if currentFolder != "" {
		gtkFileChooserSetCurrentFolder(fc, currentFolder)
	}

	buildStringAndFree := func(s pointer) string {
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
		gFree(s) // so we don't have to iterate a second time
		return string(bytes)
	}

	response := gtkDialogRun(fc)
	selections := []string{}
	if response == GtkResponseAccept {
		filenames := gtkFileChooserGetFilenames(fc)
		iter := filenames
		count := 0
		for {
			selections = append(selections, buildStringAndFree(iter.data))
			iter = iter.next
			if iter == nil || count == 1024 {
				break
			}
			count++
		}
	}
	defer gtkWidgetDestroy(fc)
	return selections, nil
}

// dialog related
func runOpenFileDialog(dialog *OpenFileDialogStruct) ([]string, error) {
	const GtkFileChooserActionOpen = 0
	const GtkFileChooserActionSelectFolder = 2

	var action int

	if dialog.canChooseDirectories {
		action = GtkFileChooserActionSelectFolder
	} else {
		action = GtkFileChooserActionOpen
	}

	window := pointer(0)
	if dialog.window != nil {
		window = (dialog.window.(*WebviewWindow).impl).(*linuxWebviewWindow).window
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
	dType, ok := map[DialogType]int{
		InfoDialogType:     GtkMessageInfo,
		WarningDialogType:  GtkMessageWarning,
		QuestionDialogType: GtkMessageQuestion,
	}[options.DialogType]
	if !ok {
		// FIXME: Add logging here!
		dType = GtkMessageInfo
	}
	buttonMask := GtkButtonsOk
	if len(options.Buttons) > 0 {
		buttonMask = GtkButtonsNone
	}

	dialog := gtkMessageDialogNew(
		pointer(parent),
		GtkDialogModal|GtkDialogDestroyWithParent,
		dType,
		buttonMask,
		options.Message)

	if options.Title != "" {
		gtkWindowSetTitle(dialog, options.Title)
	}

	GdkColorspaceRGB := 0

	if img, err := pngToImage(options.Icon); err == nil {
		gbytes := gBytesNewStatic(uintptr(unsafe.Pointer(&img.Pix[0])), len(img.Pix))

		defer gBytesUnref(gbytes)
		pixBuf := gdkPixbufNewFromBytes(
			gbytes,
			GdkColorspaceRGB,
			1, // has_alpha
			8,
			img.Bounds().Dx(),
			img.Bounds().Dy(),
			img.Stride,
		)
		image := gtkImageNewFromPixbuf(pixBuf)
		widgetSetVisible(image, false)
		contentArea := gtkDialogGetContentArea(dialog)
		gtkContainerAdd(contentArea, image)
	}

	for i, button := range options.Buttons {
		gtkDialogAddButton(
			dialog,
			button.Label,
			i,
		)
		if button.IsDefault {
			gtkDialogSetDefaultResponse(dialog, i)
		}
	}
	defer gtkWidgetDestroy(dialog)
	return gtkDialogRun(dialog)
}

func runSaveFileDialog(dialog *SaveFileDialogStruct) (string, error) {
	const GtkFileChooserActionSave = 1

	window := pointer(0)
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
		GtkFileChooserActionSave,
		buttonText,
		dialog.filters)

	if err != nil || len(results) == 0 {
		return "", err
	}

	return results[0], nil
}

func isOnMainThread() bool {
	return mainThreadId == gThreadSelf()
}
