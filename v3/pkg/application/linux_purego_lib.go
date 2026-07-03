//go:build linux && purego && !gtk3 && !android && !server

package application

// Foundation for the CGO-free Linux backend.
//
// Instead of linking GTK4/WebKitGTK at build time via cgo, every library is
// dlopen(3)ed at runtime and each C function is bound with
// purego.RegisterLibFunc. This file owns:
//
//   - library loading (with per-distro soname fallbacks and an actionable
//     error message when a library or symbol is missing),
//   - the full function-variable registry for GLib/GObject/Gio/GTK4/
//     WebKitGTK-6.0/JavaScriptCore,
//   - small helpers shared by the rest of the backend (C-string conversion,
//     gboolean handling, GValue construction).
//
// Conventions:
//   - All C pointers are uintptr (the shared `pointer` type aliases uintptr).
//   - gboolean is int32 on the C side; helpers gbool/goBool convert.
//   - char* RETURN values that the caller must g_free are declared uintptr;
//     use goString()+gFree. Const char* returns are declared string (purego
//     copies them, nothing to free).
//   - Functions that only exist in newer library versions are registered with
//     registerOptional and must be nil-checked (or capability-checked via
//     haveSymbol) before use. Calling an unregistered function is a nil
//     dereference; calling a missing C symbol blindly would crash — there is
//     no compile-time SDK floor in a purego build, the runtime check IS the
//     guard.

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/ebitengine/purego"
)

type windowPointer uintptr
type identifier uint
type pointer uintptr

type GSList struct {
	data pointer
	next *GSList
}

type GSListPointer *GSList

const nilPointer pointer = 0

var nilRadioGroup GSListPointer = nil

// ----------------------------------------------------------------------------
// Library handles
// ----------------------------------------------------------------------------

var (
	libGLib    uintptr // libglib-2.0
	libGObject uintptr // libgobject-2.0
	libGio     uintptr // libgio-2.0
	libGtk     uintptr // libgtk-4
	libWebKit  uintptr // libwebkitgtk-6.0
	libJSC     uintptr // libjavascriptcoregtk-6.0
	libSoup    uintptr // libsoup-3.0
	libC       uintptr // libc (for sigaction)

	linuxLibsErr error
)

// soname candidates per library, most specific first. The unversioned name is
// tried last: it usually only exists when -dev packages are installed.
var linuxLibCandidates = []struct {
	handle  *uintptr
	install string // human hint: Debian/Ubuntu package (others named in error)
	names   []string
}{
	{&libGLib, "libglib2.0-0", []string{"libglib-2.0.so.0", "libglib-2.0.so"}},
	{&libGObject, "libglib2.0-0", []string{"libgobject-2.0.so.0", "libgobject-2.0.so"}},
	{&libGio, "libglib2.0-0", []string{"libgio-2.0.so.0", "libgio-2.0.so"}},
	{&libGtk, "libgtk-4-1", []string{"libgtk-4.so.1", "libgtk-4.so"}},
	{&libWebKit, "libwebkitgtk-6.0-4", []string{"libwebkitgtk-6.0.so.4", "libwebkitgtk-6.0.so"}},
	{&libJSC, "libjavascriptcoregtk-6.0-1", []string{"libjavascriptcoregtk-6.0.so.1", "libjavascriptcoregtk-6.0.so"}},
	{&libSoup, "libsoup-3.0-0", []string{"libsoup-3.0.so.0", "libsoup-3.0.so"}},
	{&libC, "libc6", []string{"libc.so.6", "libc.so"}},
}

func dlopenFirst(names []string) (uintptr, []string) {
	var errs []string
	for _, name := range names {
		handle, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil && handle != 0 {
			return handle, nil
		}
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	return 0, errs
}

// missingSymbols accumulates required symbols that failed to resolve, so the
// user gets one complete report instead of the first failure.
var missingSymbols []string

func register(fptr any, lib uintptr, name string) {
	if lib == 0 {
		return // library load already failed; reported separately
	}
	sym, err := purego.Dlsym(lib, name)
	if err != nil || sym == 0 {
		missingSymbols = append(missingSymbols, name)
		return
	}
	purego.RegisterFunc(fptr, sym)
}

// registerOptional binds fptr only if the symbol exists in the loaded library
// version. Returns false (leaving the func nil) when absent — the caller must
// check before use.
func registerOptional(fptr any, lib uintptr, name string) bool {
	if lib == 0 {
		return false
	}
	sym, err := purego.Dlsym(lib, name)
	if err != nil || sym == 0 {
		return false
	}
	purego.RegisterFunc(fptr, sym)
	return true
}

func haveSymbol(lib uintptr, name string) bool {
	if lib == 0 {
		return false
	}
	sym, err := purego.Dlsym(lib, name)
	return err == nil && sym != 0
}

// loadLinuxLibraries dlopens everything and registers every function variable.
// Called from init() so that mainThreadId can be captured on the startup
// thread; a failure is recorded in linuxLibsErr and reported from appNew with
// an actionable message instead of panicking at import time.
func loadLinuxLibraries() error {
	var missingLibs []string
	for _, lib := range linuxLibCandidates {
		handle, errs := dlopenFirst(lib.names)
		if handle == 0 {
			detail := ""
			if len(errs) > 0 {
				detail = " (" + errs[len(errs)-1] + ")"
			}
			missingLibs = append(missingLibs,
				fmt.Sprintf("  %s [Debian/Ubuntu package: %s]%s", lib.names[0], lib.install, detail))
			continue
		}
		*lib.handle = handle
	}
	if len(missingLibs) > 0 {
		return fmt.Errorf("this application requires GTK4 and WebKitGTK 6.0 at runtime, "+
			"but the following libraries could not be loaded:\n%s\n"+
			"Install them via your distribution's package manager, e.g.\n"+
			"  Debian/Ubuntu: sudo apt install libgtk-4-1 libwebkitgtk-6.0-4\n"+
			"  Fedora:        sudo dnf install gtk4 webkitgtk6.0\n"+
			"  Arch:          sudo pacman -S gtk4 webkitgtk-6.0",
			strings.Join(missingLibs, "\n"))
	}

	registerGLibFuncs()
	registerGtkFuncs()
	registerWebKitFuncs()
	registerLibcFuncs()

	if len(missingSymbols) > 0 {
		return fmt.Errorf("the installed GTK4/WebKitGTK libraries are missing required symbols: %s\n"+
			"This usually means the libraries are older than the minimum supported versions "+
			"(GTK 4.10+, WebKitGTK 2.40+). Please upgrade them",
			strings.Join(missingSymbols, ", "))
	}
	return nil
}

// ----------------------------------------------------------------------------
// GLib / GObject / Gio
// ----------------------------------------------------------------------------

var (
	g_thread_self              func() uintptr
	g_idle_add                 func(uintptr, uintptr) uint32
	g_timeout_add_full         func(int32, uint32, uintptr, uintptr, uintptr) uint32
	g_free                     func(uintptr)
	g_malloc0                  func(uintptr) uintptr
	g_strdup                   func(string) uintptr
	g_get_application_name     func() uintptr // const, do NOT free
	g_set_prgname              func(string)
	g_main_context_default     func() uintptr
	g_main_context_iteration   func(uintptr, int32) int32
	g_main_context_invoke      func(uintptr, uintptr, uintptr)
	g_slist_length             func(uintptr) uint32
	g_bytes_new                func(uintptr, uintptr) uintptr
	g_bytes_unref              func(uintptr)
	g_variant_new_boolean      func(int32) uintptr
	g_variant_new_string       func(string) uintptr
	g_variant_get_boolean      func(uintptr) int32
	g_variant_get_string       func(uintptr, uintptr) uintptr // const return
	g_variant_unref            func(uintptr)
	g_variant_type_new         func(string) uintptr
	g_variant_type_free        func(uintptr)
	g_error_free               func(uintptr)
	g_error_new_literal        func(uint32, int32, string) uintptr
	g_quark_from_static_string func(uintptr) uint32

	g_object_ref                 func(uintptr) uintptr
	g_object_unref               func(uintptr)
	g_object_ref_sink            func(uintptr) uintptr
	g_object_set_data            func(uintptr, string, uintptr)
	g_object_get_data            func(uintptr, string) uintptr
	g_object_new_with_properties func(uintptr, uint32, uintptr, uintptr) uintptr
	g_signal_connect_data        func(uintptr, string, uintptr, uintptr, uintptr, uint32) uint64
	g_value_init                 func(uintptr, uintptr) uintptr
	g_value_set_object           func(uintptr, uintptr)
	g_value_set_string           func(uintptr, string)
	g_value_unset                func(uintptr)
	g_type_from_name             func(string) uintptr

	g_application_run     func(uintptr, int32, uintptr) int32
	g_application_hold    func(uintptr)
	g_application_release func(uintptr)
	g_application_quit    func(uintptr)

	g_simple_action_group_new    func() uintptr
	g_simple_action_new          func(string, uintptr) uintptr
	g_simple_action_new_stateful func(string, uintptr, uintptr) uintptr
	g_simple_action_set_enabled  func(uintptr, int32)
	g_simple_action_set_state    func(uintptr, uintptr)
	g_action_map_add_action      func(uintptr, uintptr)
	g_action_map_lookup_action   func(uintptr, string) uintptr
	g_action_get_state           func(uintptr) uintptr

	g_menu_new              func() uintptr
	g_menu_item_new         func(string, uintptr) uintptr
	g_menu_item_set_label   func(uintptr, string)
	g_menu_item_set_submenu func(uintptr, uintptr)
	// g_menu_item_set_action_and_target is variadic; the _value variant
	// takes the target as a GVariant instead.
	g_menu_item_set_action_and_target_value func(uintptr, string, uintptr)
	g_menu_append_item                      func(uintptr, uintptr)
	g_menu_append_section                   func(uintptr, uintptr, uintptr)
	g_menu_remove                           func(uintptr, int32)
	g_menu_remove_all                       func(uintptr)
	g_menu_insert_item                      func(uintptr, int32, uintptr)

	g_list_store_new         func(uintptr) uintptr
	g_list_store_append      func(uintptr, uintptr)
	g_list_model_get_n_items func(uintptr) uint32
	g_list_model_get_item    func(uintptr, uint32) uintptr

	g_file_new_for_path func(string) uintptr
	g_file_get_path     func(uintptr) uintptr // transfer full

	g_unix_input_stream_new func(int32, int32) uintptr
	g_input_stream_read_all func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) int32
	g_input_stream_close    func(uintptr, uintptr, uintptr) int32

	g_strsplit func(string, string, int32) uintptr
	g_strfreev func(uintptr)
	g_strstrip func(uintptr) uintptr
)

func registerGLibFuncs() {
	register(&g_thread_self, libGLib, "g_thread_self")
	register(&g_idle_add, libGLib, "g_idle_add")
	register(&g_timeout_add_full, libGLib, "g_timeout_add_full")
	register(&g_free, libGLib, "g_free")
	register(&g_malloc0, libGLib, "g_malloc0")
	register(&g_strdup, libGLib, "g_strdup")
	register(&g_get_application_name, libGLib, "g_get_application_name")
	register(&g_set_prgname, libGLib, "g_set_prgname")
	register(&g_main_context_default, libGLib, "g_main_context_default")
	register(&g_main_context_iteration, libGLib, "g_main_context_iteration")
	register(&g_main_context_invoke, libGLib, "g_main_context_invoke")
	register(&g_slist_length, libGLib, "g_slist_length")
	register(&g_bytes_new, libGLib, "g_bytes_new")
	register(&g_bytes_unref, libGLib, "g_bytes_unref")
	register(&g_variant_new_boolean, libGLib, "g_variant_new_boolean")
	register(&g_variant_new_string, libGLib, "g_variant_new_string")
	register(&g_variant_get_boolean, libGLib, "g_variant_get_boolean")
	register(&g_variant_get_string, libGLib, "g_variant_get_string")
	register(&g_variant_unref, libGLib, "g_variant_unref")
	register(&g_variant_type_new, libGLib, "g_variant_type_new")
	register(&g_variant_type_free, libGLib, "g_variant_type_free")
	register(&g_error_free, libGLib, "g_error_free")
	register(&g_error_new_literal, libGLib, "g_error_new_literal")
	register(&g_quark_from_static_string, libGLib, "g_quark_from_static_string")
	register(&g_strsplit, libGLib, "g_strsplit")
	register(&g_strfreev, libGLib, "g_strfreev")
	register(&g_strstrip, libGLib, "g_strstrip")

	register(&g_object_ref, libGObject, "g_object_ref")
	register(&g_object_unref, libGObject, "g_object_unref")
	register(&g_object_ref_sink, libGObject, "g_object_ref_sink")
	register(&g_object_set_data, libGObject, "g_object_set_data")
	register(&g_object_get_data, libGObject, "g_object_get_data")
	register(&g_object_new_with_properties, libGObject, "g_object_new_with_properties")
	register(&g_signal_connect_data, libGObject, "g_signal_connect_data")
	register(&g_value_init, libGObject, "g_value_init")
	register(&g_value_set_object, libGObject, "g_value_set_object")
	register(&g_value_unset, libGObject, "g_value_unset")
	register(&g_type_from_name, libGObject, "g_type_from_name")

	register(&g_application_run, libGio, "g_application_run")
	register(&g_application_hold, libGio, "g_application_hold")
	register(&g_application_release, libGio, "g_application_release")
	register(&g_application_quit, libGio, "g_application_quit")
	register(&g_simple_action_group_new, libGio, "g_simple_action_group_new")
	register(&g_simple_action_new, libGio, "g_simple_action_new")
	register(&g_simple_action_new_stateful, libGio, "g_simple_action_new_stateful")
	register(&g_simple_action_set_enabled, libGio, "g_simple_action_set_enabled")
	register(&g_simple_action_set_state, libGio, "g_simple_action_set_state")
	register(&g_action_map_add_action, libGio, "g_action_map_add_action")
	register(&g_action_map_lookup_action, libGio, "g_action_map_lookup_action")
	register(&g_action_get_state, libGio, "g_action_get_state")
	register(&g_menu_new, libGio, "g_menu_new")
	register(&g_menu_item_new, libGio, "g_menu_item_new")
	register(&g_menu_item_set_label, libGio, "g_menu_item_set_label")
	register(&g_menu_item_set_submenu, libGio, "g_menu_item_set_submenu")
	register(&g_menu_item_set_action_and_target_value, libGio, "g_menu_item_set_action_and_target_value")
	register(&g_menu_append_item, libGio, "g_menu_append_item")
	register(&g_menu_append_section, libGio, "g_menu_append_section")
	register(&g_menu_remove, libGio, "g_menu_remove")
	register(&g_menu_remove_all, libGio, "g_menu_remove_all")
	register(&g_menu_insert_item, libGio, "g_menu_insert_item")
	register(&g_list_store_new, libGio, "g_list_store_new")
	register(&g_list_store_append, libGio, "g_list_store_append")
	register(&g_list_model_get_n_items, libGio, "g_list_model_get_n_items")
	register(&g_list_model_get_item, libGio, "g_list_model_get_item")
	register(&g_file_new_for_path, libGio, "g_file_new_for_path")
	register(&g_file_get_path, libGio, "g_file_get_path")
	register(&g_unix_input_stream_new, libGio, "g_unix_input_stream_new")
	register(&g_input_stream_read_all, libGio, "g_input_stream_read_all")
	register(&g_input_stream_close, libGio, "g_input_stream_close")
}

// ----------------------------------------------------------------------------
// GTK4 / GDK
// ----------------------------------------------------------------------------

var (
	gtk_application_new                   func(string, uint32) uintptr
	gtk_application_get_active_window     func(uintptr) uintptr
	gtk_application_get_windows           func(uintptr) *GSList
	gtk_application_window_new            func(uintptr) uintptr
	gtk_application_set_accels_for_action func(uintptr, string, uintptr)
	gtk_accelerator_name                  func(uint32, uint32) uintptr // transfer full

	gtk_window_new                func() uintptr
	gtk_window_destroy            func(uintptr)
	gtk_window_present            func(uintptr)
	gtk_window_set_title          func(uintptr, string)
	gtk_window_set_default_size   func(uintptr, int32, int32)
	gtk_window_get_default_size   func(uintptr, uintptr, uintptr)
	gtk_window_set_resizable      func(uintptr, int32)
	gtk_window_set_decorated      func(uintptr, int32)
	gtk_window_set_child          func(uintptr, uintptr)
	gtk_window_set_titlebar       func(uintptr, uintptr)
	gtk_window_set_modal          func(uintptr, int32)
	gtk_window_set_transient_for  func(uintptr, uintptr)
	gtk_window_set_default_widget func(uintptr, uintptr)
	gtk_window_fullscreen         func(uintptr)
	gtk_window_unfullscreen       func(uintptr)
	gtk_window_maximize           func(uintptr)
	gtk_window_unmaximize         func(uintptr)
	gtk_window_minimize           func(uintptr)
	gtk_window_is_fullscreen      func(uintptr) int32
	gtk_window_is_maximized       func(uintptr) int32
	gtk_window_is_active          func(uintptr) int32

	gtk_widget_set_visible         func(uintptr, int32)
	gtk_widget_is_visible          func(uintptr) int32
	gtk_widget_set_sensitive       func(uintptr, int32)
	gtk_widget_set_opacity         func(uintptr, float64)
	gtk_widget_set_name            func(uintptr, string)
	gtk_widget_set_vexpand         func(uintptr, int32)
	gtk_widget_set_hexpand         func(uintptr, int32)
	gtk_widget_set_size_request    func(uintptr, int32, int32)
	gtk_widget_get_width           func(uintptr) int32
	gtk_widget_get_height          func(uintptr) int32
	gtk_widget_get_display         func(uintptr) uintptr
	gtk_widget_get_native          func(uintptr) uintptr
	gtk_widget_add_controller      func(uintptr, uintptr)
	gtk_widget_add_css_class       func(uintptr, string)
	gtk_widget_set_parent          func(uintptr, uintptr)
	gtk_widget_unparent            func(uintptr)
	gtk_widget_insert_action_group func(uintptr, string, uintptr)
	gtk_widget_set_halign          func(uintptr, int32)
	gtk_widget_set_margin_start    func(uintptr, int32)
	gtk_widget_set_margin_end      func(uintptr, int32)
	gtk_widget_set_margin_top      func(uintptr, int32)
	gtk_widget_set_margin_bottom   func(uintptr, int32)
	gtk_widget_set_tooltip_text    func(uintptr, string)
	gtk_widget_grab_focus          func(uintptr) int32
	gtk_widget_activate            func(uintptr) int32

	gtk_box_new     func(int32, int32) uintptr
	gtk_box_append  func(uintptr, uintptr)
	gtk_box_prepend func(uintptr, uintptr)

	gtk_native_get_surface func(uintptr) uintptr

	gtk_event_controller_focus_new             func() uintptr
	gtk_event_controller_key_new               func() uintptr
	gtk_event_controller_set_propagation_phase func(uintptr, int32)
	gtk_gesture_click_new                      func() uintptr
	gtk_gesture_single_set_button              func(uintptr, uint32)
	gtk_gesture_single_get_current_button      func(uintptr) uint32
	gtk_drop_controller_motion_new             func() uintptr
	gtk_drop_target_new                        func(uintptr, uint32) uintptr

	gtk_popover_menu_bar_new_from_model func(uintptr) uintptr
	gtk_popover_menu_new_from_model     func(uintptr) uintptr
	gtk_popover_set_has_arrow           func(uintptr, int32)
	gtk_popover_set_position            func(uintptr, int32)
	gtk_popover_set_pointing_to         func(uintptr, uintptr)
	gtk_popover_popup                   func(uintptr)
	gtk_header_bar_new                  func() uintptr
	gtk_header_bar_pack_end             func(uintptr, uintptr)
	gtk_menu_button_new                 func() uintptr
	gtk_menu_button_set_icon_name       func(uintptr, string)
	gtk_menu_button_set_menu_model      func(uintptr, uintptr)

	gtk_file_dialog_new                            func() uintptr
	gtk_file_dialog_set_title                      func(uintptr, string)
	gtk_file_dialog_set_filters                    func(uintptr, uintptr)
	gtk_file_dialog_set_initial_folder             func(uintptr, uintptr)
	gtk_file_dialog_set_accept_label               func(uintptr, string)
	gtk_file_dialog_open                           func(uintptr, uintptr, uintptr, uintptr, uintptr)
	gtk_file_dialog_open_finish                    func(uintptr, uintptr, uintptr) uintptr
	gtk_file_dialog_open_multiple                  func(uintptr, uintptr, uintptr, uintptr, uintptr)
	gtk_file_dialog_open_multiple_finish           func(uintptr, uintptr, uintptr) uintptr
	gtk_file_dialog_save                           func(uintptr, uintptr, uintptr, uintptr, uintptr)
	gtk_file_dialog_save_finish                    func(uintptr, uintptr, uintptr) uintptr
	gtk_file_dialog_select_folder                  func(uintptr, uintptr, uintptr, uintptr, uintptr)
	gtk_file_dialog_select_folder_finish           func(uintptr, uintptr, uintptr) uintptr
	gtk_file_dialog_select_multiple_folders        func(uintptr, uintptr, uintptr, uintptr, uintptr)
	gtk_file_dialog_select_multiple_folders_finish func(uintptr, uintptr, uintptr) uintptr
	gtk_file_filter_new                            func() uintptr
	gtk_file_filter_set_name                       func(uintptr, string)
	gtk_file_filter_add_pattern                    func(uintptr, string)
	gtk_file_filter_get_type                       func() uintptr

	gtk_button_new_with_label     func(string) uintptr
	gtk_label_new                 func(string) uintptr
	gtk_label_set_wrap            func(uintptr, int32)
	gtk_label_set_max_width_chars func(uintptr, int32)
	gtk_image_new_from_paintable  func(uintptr) uintptr
	gtk_image_new_from_icon_name  func(string) uintptr
	gtk_image_set_pixel_size      func(uintptr, int32)
	// gtk_accessible_update_property is variadic (unsupported by purego);
	// gtk_accessible_update_property_value takes arrays instead.
	gtk_accessible_update_property_value func(uintptr, int32, uintptr, uintptr)

	gdk_display_get_default            func() uintptr
	gdk_display_get_monitors           func(uintptr) uintptr
	gdk_display_get_clipboard          func(uintptr) uintptr
	gdk_display_get_default_seat       func(uintptr) uintptr
	gdk_display_get_monitor_at_surface func(uintptr, uintptr) uintptr
	gdk_seat_get_pointer               func(uintptr) uintptr
	gdk_surface_get_display            func(uintptr) uintptr
	gdk_toplevel_get_state             func(uintptr) uint32
	gdk_toplevel_begin_move            func(uintptr, uintptr, int32, float64, float64, uint32)
	gdk_toplevel_begin_resize          func(uintptr, int32, uintptr, int32, float64, float64, uint32)
	gdk_monitor_get_geometry           func(uintptr, uintptr)
	gdk_monitor_get_model              func(uintptr) string // const return
	gdk_monitor_get_scale_factor       func(uintptr) int32
	gdk_monitor_get_scale              func(uintptr) float64 // GTK 4.14+, optional
	gdk_clipboard_set_text             func(uintptr, string)
	gdk_clipboard_read_text_async      func(uintptr, uintptr, uintptr, uintptr)
	gdk_clipboard_read_text_finish     func(uintptr, uintptr, uintptr) uintptr
	gdk_unicode_to_keyval              func(uint32) uint32
	gdk_texture_new_from_bytes         func(uintptr, uintptr) uintptr
	gdk_texture_get_width              func(uintptr) int32
	gdk_file_list_get_type             func() uintptr
	gdk_content_formats_contain_gtype  func(uintptr, uintptr) int32
	gdk_drop_get_formats               func(uintptr) uintptr

	gtk_get_major_version func() uint32
	gtk_get_minor_version func() uint32
	gtk_get_micro_version func() uint32

	// X11 (resolved from GTK's already-loaded X11 backend, optional: absent
	// on Wayland-only systems)
	gdk_x11_display_get_xdisplay func(uintptr) uintptr
	gdk_x11_surface_get_xid      func(uintptr) uintptr
)

func registerGtkFuncs() {
	register(&gtk_application_new, libGtk, "gtk_application_new")
	register(&gtk_application_get_active_window, libGtk, "gtk_application_get_active_window")
	register(&gtk_application_get_windows, libGtk, "gtk_application_get_windows")
	register(&gtk_application_window_new, libGtk, "gtk_application_window_new")
	register(&gtk_application_set_accels_for_action, libGtk, "gtk_application_set_accels_for_action")
	register(&gtk_accelerator_name, libGtk, "gtk_accelerator_name")

	register(&gtk_window_new, libGtk, "gtk_window_new")
	register(&gtk_window_destroy, libGtk, "gtk_window_destroy")
	register(&gtk_window_present, libGtk, "gtk_window_present")
	register(&gtk_window_set_title, libGtk, "gtk_window_set_title")
	register(&gtk_window_set_default_size, libGtk, "gtk_window_set_default_size")
	register(&gtk_window_get_default_size, libGtk, "gtk_window_get_default_size")
	register(&gtk_window_set_resizable, libGtk, "gtk_window_set_resizable")
	register(&gtk_window_set_decorated, libGtk, "gtk_window_set_decorated")
	register(&gtk_window_set_child, libGtk, "gtk_window_set_child")
	register(&gtk_window_set_titlebar, libGtk, "gtk_window_set_titlebar")
	register(&gtk_window_set_modal, libGtk, "gtk_window_set_modal")
	register(&gtk_window_set_transient_for, libGtk, "gtk_window_set_transient_for")
	register(&gtk_window_set_default_widget, libGtk, "gtk_window_set_default_widget")
	register(&gtk_window_fullscreen, libGtk, "gtk_window_fullscreen")
	register(&gtk_window_unfullscreen, libGtk, "gtk_window_unfullscreen")
	register(&gtk_window_maximize, libGtk, "gtk_window_maximize")
	register(&gtk_window_unmaximize, libGtk, "gtk_window_unmaximize")
	register(&gtk_window_minimize, libGtk, "gtk_window_minimize")
	register(&gtk_window_is_fullscreen, libGtk, "gtk_window_is_fullscreen")
	register(&gtk_window_is_maximized, libGtk, "gtk_window_is_maximized")
	register(&gtk_window_is_active, libGtk, "gtk_window_is_active")

	register(&gtk_widget_set_visible, libGtk, "gtk_widget_set_visible")
	register(&gtk_widget_is_visible, libGtk, "gtk_widget_is_visible")
	register(&gtk_widget_set_sensitive, libGtk, "gtk_widget_set_sensitive")
	register(&gtk_widget_set_opacity, libGtk, "gtk_widget_set_opacity")
	register(&gtk_widget_set_name, libGtk, "gtk_widget_set_name")
	register(&gtk_widget_set_vexpand, libGtk, "gtk_widget_set_vexpand")
	register(&gtk_widget_set_hexpand, libGtk, "gtk_widget_set_hexpand")
	register(&gtk_widget_set_size_request, libGtk, "gtk_widget_set_size_request")
	register(&gtk_widget_get_width, libGtk, "gtk_widget_get_width")
	register(&gtk_widget_get_height, libGtk, "gtk_widget_get_height")
	register(&gtk_widget_get_display, libGtk, "gtk_widget_get_display")
	register(&gtk_widget_get_native, libGtk, "gtk_widget_get_native")
	register(&gtk_widget_add_controller, libGtk, "gtk_widget_add_controller")
	register(&gtk_widget_add_css_class, libGtk, "gtk_widget_add_css_class")
	register(&gtk_widget_set_parent, libGtk, "gtk_widget_set_parent")
	register(&gtk_widget_unparent, libGtk, "gtk_widget_unparent")
	register(&gtk_widget_insert_action_group, libGtk, "gtk_widget_insert_action_group")
	register(&gtk_widget_set_halign, libGtk, "gtk_widget_set_halign")
	register(&gtk_widget_set_margin_start, libGtk, "gtk_widget_set_margin_start")
	register(&gtk_widget_set_margin_end, libGtk, "gtk_widget_set_margin_end")
	register(&gtk_widget_set_margin_top, libGtk, "gtk_widget_set_margin_top")
	register(&gtk_widget_set_margin_bottom, libGtk, "gtk_widget_set_margin_bottom")
	register(&gtk_widget_set_tooltip_text, libGtk, "gtk_widget_set_tooltip_text")
	register(&gtk_widget_grab_focus, libGtk, "gtk_widget_grab_focus")
	register(&gtk_widget_activate, libGtk, "gtk_widget_activate")

	register(&gtk_box_new, libGtk, "gtk_box_new")
	register(&gtk_box_append, libGtk, "gtk_box_append")
	register(&gtk_box_prepend, libGtk, "gtk_box_prepend")
	register(&gtk_native_get_surface, libGtk, "gtk_native_get_surface")

	register(&gtk_event_controller_focus_new, libGtk, "gtk_event_controller_focus_new")
	register(&gtk_event_controller_key_new, libGtk, "gtk_event_controller_key_new")
	register(&gtk_event_controller_set_propagation_phase, libGtk, "gtk_event_controller_set_propagation_phase")
	register(&gtk_gesture_click_new, libGtk, "gtk_gesture_click_new")
	register(&gtk_gesture_single_set_button, libGtk, "gtk_gesture_single_set_button")
	register(&gtk_gesture_single_get_current_button, libGtk, "gtk_gesture_single_get_current_button")
	register(&gtk_drop_controller_motion_new, libGtk, "gtk_drop_controller_motion_new")
	register(&gtk_drop_target_new, libGtk, "gtk_drop_target_new")

	register(&gtk_popover_menu_bar_new_from_model, libGtk, "gtk_popover_menu_bar_new_from_model")
	register(&gtk_popover_menu_new_from_model, libGtk, "gtk_popover_menu_new_from_model")
	register(&gtk_popover_set_has_arrow, libGtk, "gtk_popover_set_has_arrow")
	register(&gtk_popover_set_position, libGtk, "gtk_popover_set_position")
	register(&gtk_popover_set_pointing_to, libGtk, "gtk_popover_set_pointing_to")
	register(&gtk_popover_popup, libGtk, "gtk_popover_popup")
	register(&gtk_header_bar_new, libGtk, "gtk_header_bar_new")
	register(&gtk_header_bar_pack_end, libGtk, "gtk_header_bar_pack_end")
	register(&gtk_menu_button_new, libGtk, "gtk_menu_button_new")
	register(&gtk_menu_button_set_icon_name, libGtk, "gtk_menu_button_set_icon_name")
	register(&gtk_menu_button_set_menu_model, libGtk, "gtk_menu_button_set_menu_model")

	register(&gtk_file_dialog_new, libGtk, "gtk_file_dialog_new")
	register(&gtk_file_dialog_set_title, libGtk, "gtk_file_dialog_set_title")
	register(&gtk_file_dialog_set_filters, libGtk, "gtk_file_dialog_set_filters")
	register(&gtk_file_dialog_set_initial_folder, libGtk, "gtk_file_dialog_set_initial_folder")
	register(&gtk_file_dialog_set_accept_label, libGtk, "gtk_file_dialog_set_accept_label")
	register(&gtk_file_dialog_open, libGtk, "gtk_file_dialog_open")
	register(&gtk_file_dialog_open_finish, libGtk, "gtk_file_dialog_open_finish")
	register(&gtk_file_dialog_open_multiple, libGtk, "gtk_file_dialog_open_multiple")
	register(&gtk_file_dialog_open_multiple_finish, libGtk, "gtk_file_dialog_open_multiple_finish")
	register(&gtk_file_dialog_save, libGtk, "gtk_file_dialog_save")
	register(&gtk_file_dialog_save_finish, libGtk, "gtk_file_dialog_save_finish")
	register(&gtk_file_dialog_select_folder, libGtk, "gtk_file_dialog_select_folder")
	register(&gtk_file_dialog_select_folder_finish, libGtk, "gtk_file_dialog_select_folder_finish")
	register(&gtk_file_dialog_select_multiple_folders, libGtk, "gtk_file_dialog_select_multiple_folders")
	register(&gtk_file_dialog_select_multiple_folders_finish, libGtk, "gtk_file_dialog_select_multiple_folders_finish")
	register(&gtk_file_filter_new, libGtk, "gtk_file_filter_new")
	register(&gtk_file_filter_set_name, libGtk, "gtk_file_filter_set_name")
	register(&gtk_file_filter_add_pattern, libGtk, "gtk_file_filter_add_pattern")
	register(&gtk_file_filter_get_type, libGtk, "gtk_file_filter_get_type")

	register(&gtk_button_new_with_label, libGtk, "gtk_button_new_with_label")
	register(&gtk_label_new, libGtk, "gtk_label_new")
	register(&gtk_label_set_wrap, libGtk, "gtk_label_set_wrap")
	register(&gtk_label_set_max_width_chars, libGtk, "gtk_label_set_max_width_chars")
	register(&gtk_image_new_from_paintable, libGtk, "gtk_image_new_from_paintable")
	register(&gtk_image_new_from_icon_name, libGtk, "gtk_image_new_from_icon_name")
	register(&gtk_image_set_pixel_size, libGtk, "gtk_image_set_pixel_size")
	register(&gtk_accessible_update_property_value, libGtk, "gtk_accessible_update_property_value")
	register(&g_value_set_string, libGObject, "g_value_set_string")

	register(&gdk_display_get_default, libGtk, "gdk_display_get_default")
	register(&gdk_display_get_monitors, libGtk, "gdk_display_get_monitors")
	register(&gdk_display_get_clipboard, libGtk, "gdk_display_get_clipboard")
	register(&gdk_display_get_default_seat, libGtk, "gdk_display_get_default_seat")
	register(&gdk_display_get_monitor_at_surface, libGtk, "gdk_display_get_monitor_at_surface")
	register(&gdk_seat_get_pointer, libGtk, "gdk_seat_get_pointer")
	register(&gdk_surface_get_display, libGtk, "gdk_surface_get_display")
	register(&gdk_toplevel_get_state, libGtk, "gdk_toplevel_get_state")
	register(&gdk_toplevel_begin_move, libGtk, "gdk_toplevel_begin_move")
	register(&gdk_toplevel_begin_resize, libGtk, "gdk_toplevel_begin_resize")
	register(&gdk_monitor_get_geometry, libGtk, "gdk_monitor_get_geometry")
	register(&gdk_monitor_get_model, libGtk, "gdk_monitor_get_model")
	register(&gdk_monitor_get_scale_factor, libGtk, "gdk_monitor_get_scale_factor")
	// GTK 4.14+: fractional scaling. Fall back to gdk_monitor_get_scale_factor
	// (integer) on older GTK4 — see monitorScale().
	registerOptional(&gdk_monitor_get_scale, libGtk, "gdk_monitor_get_scale")
	register(&gdk_clipboard_set_text, libGtk, "gdk_clipboard_set_text")
	register(&gdk_clipboard_read_text_async, libGtk, "gdk_clipboard_read_text_async")
	register(&gdk_clipboard_read_text_finish, libGtk, "gdk_clipboard_read_text_finish")
	register(&gdk_unicode_to_keyval, libGtk, "gdk_unicode_to_keyval")
	register(&gdk_texture_new_from_bytes, libGtk, "gdk_texture_new_from_bytes")
	register(&gdk_texture_get_width, libGtk, "gdk_texture_get_width")
	register(&gdk_file_list_get_type, libGtk, "gdk_file_list_get_type")
	register(&gdk_content_formats_contain_gtype, libGtk, "gdk_content_formats_contain_gtype")
	register(&gdk_drop_get_formats, libGtk, "gdk_drop_get_formats")

	register(&gtk_get_major_version, libGtk, "gtk_get_major_version")
	register(&gtk_get_minor_version, libGtk, "gtk_get_minor_version")
	register(&gtk_get_micro_version, libGtk, "gtk_get_micro_version")

	// The X11 backend symbols only exist when GTK was built with X11 support
	// (absent on Wayland-only builds) — optional, nil-checked at use.
	registerOptional(&gdk_x11_display_get_xdisplay, libGtk, "gdk_x11_display_get_xdisplay")
	registerOptional(&gdk_x11_surface_get_xid, libGtk, "gdk_x11_surface_get_xid")
}

// ----------------------------------------------------------------------------
// WebKitGTK 6.0 / JavaScriptCore
// ----------------------------------------------------------------------------

var (
	webkit_web_view_get_type                                    func() uintptr
	webkit_web_view_get_context                                 func(uintptr) uintptr
	webkit_web_view_get_user_content_manager                    func(uintptr) uintptr
	webkit_web_view_get_settings                                func(uintptr) uintptr
	webkit_web_view_set_settings                                func(uintptr, uintptr)
	webkit_web_view_load_uri                                    func(uintptr, string)
	webkit_web_view_load_alternate_html                         func(uintptr, string, string, string)
	webkit_web_view_get_uri                                     func(uintptr) uintptr // const return
	webkit_web_view_evaluate_javascript                         func(uintptr, uintptr, int, uintptr, uintptr, uintptr, uintptr, uintptr)
	webkit_web_view_get_zoom_level                              func(uintptr) float64
	webkit_web_view_set_zoom_level                              func(uintptr, float64)
	webkit_web_view_set_background_color                        func(uintptr, uintptr)
	webkit_web_view_get_inspector                               func(uintptr) uintptr
	webkit_web_inspector_show                                   func(uintptr)
	webkit_settings_new                                         func() uintptr
	webkit_settings_set_hardware_acceleration_policy            func(uintptr, int32)
	webkit_settings_get_enable_developer_extras                 func(uintptr) int32
	webkit_settings_set_enable_developer_extras                 func(uintptr, int32)
	webkit_user_content_manager_new                             func() uintptr
	webkit_user_content_manager_register_script_message_handler func(uintptr, string, uintptr) int32
	webkit_web_context_register_uri_scheme                      func(uintptr, string, uintptr, uintptr, uintptr)
	webkit_network_session_get_default                          func() uintptr
	webkit_uri_scheme_request_get_web_view                      func(uintptr) uintptr
	webkit_permission_request_allow                             func(uintptr)
	webkit_permission_request_deny                              func(uintptr)
	webkit_user_media_permission_is_for_audio_device            func(uintptr) int32
	webkit_user_media_permission_is_for_video_device            func(uintptr) int32
	webkit_get_major_version                                    func() uint32
	webkit_get_minor_version                                    func() uint32
	webkit_get_micro_version                                    func() uint32

	jsc_value_to_string func(uintptr) uintptr // transfer full

	g_type_check_instance_is_a                    func(uintptr, uintptr) int32
	g_type_check_value_holds                      func(uintptr, uintptr) int32
	g_value_get_boxed                             func(uintptr) uintptr
	webkit_user_media_permission_request_get_type func() uintptr
)

func registerWebKitFuncs() {
	register(&webkit_web_view_get_type, libWebKit, "webkit_web_view_get_type")
	register(&webkit_web_view_get_context, libWebKit, "webkit_web_view_get_context")
	register(&webkit_web_view_get_user_content_manager, libWebKit, "webkit_web_view_get_user_content_manager")
	register(&webkit_web_view_get_settings, libWebKit, "webkit_web_view_get_settings")
	register(&webkit_web_view_set_settings, libWebKit, "webkit_web_view_set_settings")
	register(&webkit_web_view_load_uri, libWebKit, "webkit_web_view_load_uri")
	register(&webkit_web_view_load_alternate_html, libWebKit, "webkit_web_view_load_alternate_html")
	register(&webkit_web_view_get_uri, libWebKit, "webkit_web_view_get_uri")
	register(&webkit_web_view_evaluate_javascript, libWebKit, "webkit_web_view_evaluate_javascript")
	register(&webkit_web_view_get_zoom_level, libWebKit, "webkit_web_view_get_zoom_level")
	register(&webkit_web_view_set_zoom_level, libWebKit, "webkit_web_view_set_zoom_level")
	register(&webkit_web_view_set_background_color, libWebKit, "webkit_web_view_set_background_color")
	register(&webkit_web_view_get_inspector, libWebKit, "webkit_web_view_get_inspector")
	register(&webkit_web_inspector_show, libWebKit, "webkit_web_inspector_show")
	register(&webkit_settings_new, libWebKit, "webkit_settings_new")
	register(&webkit_settings_set_hardware_acceleration_policy, libWebKit, "webkit_settings_set_hardware_acceleration_policy")
	register(&webkit_settings_get_enable_developer_extras, libWebKit, "webkit_settings_get_enable_developer_extras")
	register(&webkit_settings_set_enable_developer_extras, libWebKit, "webkit_settings_set_enable_developer_extras")
	register(&webkit_user_content_manager_new, libWebKit, "webkit_user_content_manager_new")
	register(&webkit_user_content_manager_register_script_message_handler, libWebKit, "webkit_user_content_manager_register_script_message_handler")
	register(&webkit_web_context_register_uri_scheme, libWebKit, "webkit_web_context_register_uri_scheme")
	register(&webkit_network_session_get_default, libWebKit, "webkit_network_session_get_default")
	register(&webkit_uri_scheme_request_get_web_view, libWebKit, "webkit_uri_scheme_request_get_web_view")
	register(&webkit_permission_request_allow, libWebKit, "webkit_permission_request_allow")
	register(&webkit_permission_request_deny, libWebKit, "webkit_permission_request_deny")
	register(&webkit_user_media_permission_is_for_audio_device, libWebKit, "webkit_user_media_permission_is_for_audio_device")
	register(&webkit_user_media_permission_is_for_video_device, libWebKit, "webkit_user_media_permission_is_for_video_device")
	register(&webkit_user_media_permission_request_get_type, libWebKit, "webkit_user_media_permission_request_get_type")
	register(&webkit_get_major_version, libWebKit, "webkit_get_major_version")
	register(&webkit_get_minor_version, libWebKit, "webkit_get_minor_version")
	register(&webkit_get_micro_version, libWebKit, "webkit_get_micro_version")

	register(&jsc_value_to_string, libJSC, "jsc_value_to_string")

	register(&g_type_check_instance_is_a, libGObject, "g_type_check_instance_is_a")
	register(&g_type_check_value_holds, libGObject, "g_type_check_value_holds")
	register(&g_value_get_boxed, libGObject, "g_value_get_boxed")
}

// ----------------------------------------------------------------------------
// libc (signal-handler SA_ONSTACK fix, see linux_purego.go)
// ----------------------------------------------------------------------------

var libc_sigaction func(int32, uintptr, uintptr) int32

func registerLibcFuncs() {
	register(&libc_sigaction, libC, "sigaction")
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func gbool(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// goString copies a NUL-terminated C string. The pointer is not freed.
func goString(c uintptr) string {
	if c == 0 {
		return ""
	}
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&c))
	n := 0
	for *(*byte)(unsafe.Add(ptr, n)) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(ptr), n))
}

// takeGString copies a transfer-full C string and g_free()s it.
func takeGString(c uintptr) string {
	if c == 0 {
		return ""
	}
	s := goString(c)
	g_free(c)
	return s
}

// cString allocates a NUL-terminated copy of s in C memory (g_strdup). Free
// with g_free. Only needed when a C string must outlive the call or sit
// inside an array/struct — plain string arguments are marshalled by purego.
func cString(s string) uintptr {
	return g_strdup(s)
}

// gValue mirrors GValue: a GType followed by 2 machine words of payload.
type gValue struct {
	gtype uintptr
	data  [2]uint64
}

// gObjectNewWithObjectProperty is the purego replacement for
// g_object_new(type, prop, obj, NULL): WebKitGTK 6.0 removed
// webkit_web_view_new_with_user_content_manager, and g_object_new is
// variadic, so use the non-variadic g_object_new_with_properties.
func gObjectNewWithObjectProperty(gtype uintptr, property string, obj uintptr) uintptr {
	cProp := cString(property)
	defer g_free(cProp)

	var value gValue
	g_value_init(uintptr(unsafe.Pointer(&value)), gtype_Object())
	g_value_set_object(uintptr(unsafe.Pointer(&value)), obj)

	names := []uintptr{cProp}
	result := g_object_new_with_properties(gtype, 1,
		uintptr(unsafe.Pointer(&names[0])),
		uintptr(unsafe.Pointer(&value)))
	g_value_unset(uintptr(unsafe.Pointer(&value)))
	return result
}

// gtype_Object returns G_TYPE_OBJECT. Fundamental types have fixed IDs
// (G_TYPE_OBJECT = 20 << 2), but resolve by name to stay ABI-agnostic.
var cachedGTypeObject uintptr

func gtype_Object() uintptr {
	if cachedGTypeObject == 0 {
		cachedGTypeObject = g_type_from_name("GObject")
	}
	return cachedGTypeObject
}

// gTypeInstanceIsA reports whether a GTypeInstance is of (a subtype of) the
// given GType — the purego equivalent of the WEBKIT_IS_* checking macros.
func gTypeInstanceIsA(instance, gtype uintptr) bool {
	if instance == 0 || gtype == 0 {
		return false
	}
	return g_type_check_instance_is_a(instance, gtype) != 0
}

// signalConnect wires a GObject signal to a purego callback:
// g_signal_connect(obj, sig, cb, data) is a macro over g_signal_connect_data.
func signalConnect(obj uintptr, sig string, cb uintptr, data uintptr) {
	g_signal_connect_data(obj, sig, cb, data, 0, 0)
}

// GdkRectangle (all gint)
type gdkRectangle struct {
	x, y, width, height int32
}

// GdkRGBA (all float)
type gdkRGBA struct {
	red, green, blue, alpha float32
}
