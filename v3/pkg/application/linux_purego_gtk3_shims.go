//go:build linux && purego && gtk3 && !android && !server

package application

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var gtk3Libraries sync.Once
var gtk3LibsErr error

var (
	gtk3Malloc                     func(uintptr) unsafe.Pointer
	gtk3SetData                    func(unsafe.Pointer, string, uintptr)
	gtk3GetData                    func(unsafe.Pointer, string) uintptr
	gtk3Connect                    func(unsafe.Pointer, *_Ctype_char, unsafe.Pointer, uintptr, uintptr, uint32) uint64
	gtk3Idle                       func(uintptr, uintptr) uint32
	gtk3Timeout                    func(int32, uint32, uintptr, uintptr, uintptr) uint32
	gtk3IsA                        func(unsafe.Pointer, uint64) int32
	gtk3MediaType                  func() uint64
	gtk3MediaAudio, gtk3MediaVideo func(unsafe.Pointer) int32
	gtk3Atom                       func(string, int32) uintptr
	gtk3MessageNew                 func(*_Ctype_GtkWindow, uint32, int32, int32, string, *_Ctype_char) unsafe.Pointer
	gtk3ChooserNew                 func(*_Ctype_char, *_Ctype_GtkWindow, _Ctype_GtkFileChooserAction, *_Ctype_char, int32, *_Ctype_char, int32, uintptr) unsafe.Pointer
	// gtk3GThreadSelf is bound by bindGTK3Production. It is a shim rather than
	// a generated binding so that the init() copied from the cgo source can
	// run when the libraries failed to load: _Cfunc_g_thread_self then returns nil.
	gtk3GThreadSelf func() *_Ctype_GThread
)

// gtk3LibraryNames are tried in order. The WebKit2GTK 4.1 handle also resolves
// the GTK3 and GLib symbols bound below through its dependencies.
var gtk3LibraryNames = []string{"libwebkit2gtk-4.1.so.0"}

// gtk3Binder binds every GTK3 function variable against lib. Symbols that
// cannot be resolved are appended to missing and their variables left unset.
type gtk3Binder func(lib uintptr, missing *[]string)

// bindGTK3Symbol binds *fptr to the named symbol. A missing symbol is recorded
// in missing instead of panicking, which purego.RegisterLibFunc would do.
func bindGTK3Symbol(missing *[]string, lib uintptr, fptr any, name string) {
	sym, err := purego.Dlsym(lib, name)
	if err != nil || sym == 0 {
		*missing = append(*missing, name)
		return
	}
	purego.RegisterFunc(fptr, sym)
}

// bindGTK3Production binds the full GTK3 function table.
func bindGTK3Production(lib uintptr, missing *[]string) {
	registerGTK3Functions(lib, missing)
	for name, target := range map[string]any{
		"g_malloc": &gtk3Malloc, "g_object_set_data": &gtk3SetData, "g_object_get_data": &gtk3GetData,
		"g_signal_connect_data": &gtk3Connect, "g_idle_add": &gtk3Idle, "g_timeout_add_full": &gtk3Timeout,
		"g_type_check_instance_is_a":                       &gtk3IsA,
		"webkit_user_media_permission_request_get_type":    &gtk3MediaType,
		"webkit_user_media_permission_is_for_audio_device": &gtk3MediaAudio,
		"webkit_user_media_permission_is_for_video_device": &gtk3MediaVideo,
		"gdk_atom_intern":        &gtk3Atom,
		"gtk_message_dialog_new": &gtk3MessageNew, "gtk_file_chooser_dialog_new": &gtk3ChooserNew,
		"g_thread_self": &gtk3GThreadSelf,
	} {
		bindGTK3Symbol(missing, lib, target, name)
	}
	loadGTK3DragFunctions(lib, missing)
}

// loadGTK3 dlopens the first loadable library in libNames and binds the
// function table with bind. It returns an actionable error, never panics, and
// leaves the function variables of unresolved symbols nil.
func loadGTK3(libNames []string, bind gtk3Binder) error {
	var lib uintptr
	var errs []string
	for _, name := range libNames {
		handle, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil && handle != 0 {
			lib = handle
			break
		}
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if lib == 0 {
		detail := ""
		if len(errs) > 0 {
			detail = " (" + errs[len(errs)-1] + ")"
		}
		return fmt.Errorf("this application requires GTK3 and WebKit2GTK 4.1 at runtime, "+
			"but the following library could not be loaded:\n"+
			"  %s [Debian/Ubuntu package: libwebkit2gtk-4.1-0]%s\n"+
			"Install them via your distribution's package manager, e.g.\n"+
			"  Debian/Ubuntu: sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0\n"+
			"  Fedora:        sudo dnf install gtk3 webkit2gtk4.1\n"+
			"  Arch:          sudo pacman -S gtk3 webkit2gtk-4.1",
			libNames[0], detail)
	}
	var missing []string
	bind(lib, &missing)
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("the installed GTK3/WebKit2GTK libraries are missing required symbols: %s\n"+
			"This usually means the installed libraries are older than this build expects. Please upgrade them",
			strings.Join(missing, ", "))
	}
	return nil
}

// ensureGTK3Libraries loads the GTK3 function table once. A failure is kept in
// gtk3LibsErr; it is never a panic. Use requireGTK3Libraries at startup.
func ensureGTK3Libraries() {
	gtk3Libraries.Do(func() {
		gtk3LibsErr = loadGTK3(gtk3LibraryNames, bindGTK3Production)
	})
}

// requireGTK3Libraries is inserted at the start of appNew (see gtk3.json
// "guards"). It reports a load failure the same way the GTK4 path does.
func requireGTK3Libraries() {
	ensureGTK3Libraries()
	if gtk3LibsErr != nil {
		Fatal("%v", gtk3LibsErr)
	}
}

func _Cfunc_g_thread_self() *_Ctype_GThread {
	if gtk3GThreadSelf == nil {
		return nil
	}
	return gtk3GThreadSelf()
}

func _Cfunc_CString(value string) *_Ctype_char {
	ensureGTK3Libraries()
	ptr := gtk3Malloc(uintptr(len(value) + 1))
	buf := unsafe.Slice((*byte)(ptr), len(value)+1)
	copy(buf, value)
	buf[len(value)] = 0
	return (*_Ctype_char)(ptr)
}

func _Cfunc_GoString(ptr *_Ctype_char) string {
	if ptr == nil {
		return ""
	}
	n := 0
	for *(*byte)(unsafe.Add(unsafe.Pointer(ptr), n)) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), n))
}

func _Cmacro_GDK_SELECTION_CLIPBOARD() uintptr { return gtk3Atom("CLIPBOARD", 0) }
func _Cmacro_GDK_SELECTION_PRIMARY() uintptr   { return gtk3Atom("PRIMARY", 0) }

func _Cfunc_get_compiled_gtk_major_version() _Ctype_guint { return _Cfunc_gtk_get_major_version() }
func _Cfunc_get_compiled_gtk_minor_version() _Ctype_guint { return _Cfunc_gtk_get_minor_version() }
func _Cfunc_get_compiled_gtk_micro_version() _Ctype_guint { return _Cfunc_gtk_get_micro_version() }
func _Cfunc_get_compiled_webkit_major_version() _Ctype_guint {
	return _Cfunc_webkit_get_major_version()
}
func _Cfunc_get_compiled_webkit_minor_version() _Ctype_guint {
	return _Cfunc_webkit_get_minor_version()
}
func _Cfunc_get_compiled_webkit_micro_version() _Ctype_guint {
	return _Cfunc_webkit_get_micro_version()
}

func _Cfunc_set_object_uint(obj *_Ctype_GObject, key *_Ctype_gchar, value _Ctype_guint) {
	gtk3SetData(unsafe.Pointer(obj), _Cfunc_GoString(key), uintptr(value))
}
func _Cfunc_get_object_uint(obj *_Ctype_GObject, key *_Ctype_gchar) _Ctype_guint {
	return _Ctype_guint(gtk3GetData(unsafe.Pointer(obj), _Cfunc_GoString(key)))
}
func _Cfunc_save_window_id(obj unsafe.Pointer, id _Ctype_uint) {
	gtk3SetData(obj, "windowid", uintptr(id))
}
func _Cfunc_get_window_id(obj unsafe.Pointer) _Ctype_guint {
	return _Ctype_guint(gtk3GetData(obj, "windowid"))
}
func _Cfunc_save_webview_to_content_manager(manager, webview unsafe.Pointer) {
	gtk3SetData(manager, "webview", uintptr(webview))
}
func _Cfunc_get_webview_from_content_manager(manager unsafe.Pointer) *_Ctype_WebKitWebView {
	return (*_Ctype_WebKitWebView)(unsafe.Pointer(gtk3GetData(manager, "webview")))
}
func _Cfunc_webkit_web_view(widget *_Ctype_GtkWidget) *_Ctype_WebKitWebView {
	return (*_Ctype_WebKitWebView)(unsafe.Pointer(widget))
}
func _Cfunc_is_user_media_permission_request(request *_Ctype_WebKitPermissionRequest) _Ctype_int {
	return _Ctype_int(gtk3IsA(unsafe.Pointer(request), gtk3MediaType()))
}
func _Cfunc_is_user_media_for_audio(request *_Ctype_WebKitPermissionRequest) _Ctype_int {
	return _Ctype_int(gtk3MediaAudio(unsafe.Pointer(request)))
}
func _Cfunc_is_user_media_for_video(request *_Ctype_WebKitPermissionRequest) _Ctype_int {
	return _Ctype_int(gtk3MediaVideo(unsafe.Pointer(request)))
}
func _Cfunc_signal_connect(obj unsafe.Pointer, name *_Ctype_char, callback unsafe.Pointer, data _Ctype_uintptr_t) {
	gtk3Connect(obj, name, callback, uintptr(data), 0, 0)
}

var gtk3Dispatch = purego.NewCallback(func(id uintptr) int32 { executeOnMainThread(uint(id)); return 0 })

func _Cfunc_dispatchOnMainThread(id _Ctype_uint) { gtk3Idle(gtk3Dispatch, uintptr(id)) }

// These two fixed prototypes use only integer/pointer arguments, including the
// variadic tail. No floating-point varargs or Go variadic calls cross the ABI.
func _Cfunc_new_message_dialog(parent *_Ctype_GtkWindow, message *_Ctype_gchar, typ _Ctype_int, buttons _Ctype_bool) unsafe.Pointer {
	mask := int32(1) // GTK_BUTTONS_OK
	if buttons {
		mask = 0
	}
	return gtk3MessageNew(parent, 3, int32(typ), mask, "%s", message)
}
func _Cfunc_gtkFileChooserDialogNew(title *_Ctype_char, parent *_Ctype_GtkWindow, action _Ctype_GtkFileChooserAction, cancel, accept *_Ctype_char) unsafe.Pointer {
	return gtk3ChooserNew(title, parent, action, cancel, -6, accept, -3, 0)
}

func _Cfunc_install_signal_handlers() { installSignalHandlers() }

var gtk3SignalTimer = purego.NewCallback(func(remaining uintptr) int32 {
	_Cfunc_install_signal_handlers()
	value := (*int32)(unsafe.Pointer(remaining))
	*value--
	if *value == 0 {
		return 0
	}
	return 1
})
var gtk3SignalTimerFree = purego.NewCallback(func(ptr unsafe.Pointer) { _Cfunc_free(ptr) })

func _Cfunc_schedule_signal_handler_fix() {
	ptr := gtk3Malloc(4)
	*(*int32)(ptr) = 100
	gtk3Timeout(0, 50, gtk3SignalTimer, uintptr(ptr), gtk3SignalTimerFree)
}
