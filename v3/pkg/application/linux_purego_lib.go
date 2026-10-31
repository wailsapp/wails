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

// Calloc owns native strings until Free is called.
type Calloc struct{ pool []uintptr }

func NewCalloc() *Calloc { return &Calloc{} }

func (c *Calloc) String(value string) *int8 {
	ptr := cString(value)
	c.pool = append(c.pool, ptr)
	return (*int8)(unsafe.Pointer(ptr))
}

func (c *Calloc) Free() {
	for _, ptr := range c.pool {
		g_free(ptr)
	}
	c.pool = nil
}

// Library handles

var (
	libGLib    uintptr // libglib-2.0
	libGObject uintptr // libgobject-2.0
	libGio     uintptr // libgio-2.0
	libGtk     uintptr // libgtk-4
	libWebKit  uintptr // libwebkitgtk-6.0
	libJSC     uintptr // libjavascriptcoregtk-6.0
	libSoup    uintptr // libsoup-3.0

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

	registerGeneratedNativeFunctions()

	if len(missingSymbols) > 0 {
		return fmt.Errorf("the installed GTK4/WebKitGTK libraries are missing required symbols: %s\n"+
			"This usually means the libraries are older than the minimum supported versions "+
			"(GTK 4.8+, WebKitGTK 2.40+). Please upgrade them",
			strings.Join(missingSymbols, ", "))
	}
	return nil
}

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
	gtype uint64
	data  [2]uint64
}

// gObjectNewWithObjectProperty is the purego replacement for
// g_object_new(type, prop, obj, NULL): WebKitGTK 6.0 removed
// webkit_web_view_new_with_user_content_manager, and g_object_new is
// variadic, so use the non-variadic g_object_new_with_properties.
func gObjectNewWithObjectProperty(gtype uint64, property string, obj uintptr) uintptr {
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
var cachedGTypeObject uint64

func gtype_Object() uint64 {
	if cachedGTypeObject == 0 {
		cachedGTypeObject = g_type_from_name("GObject")
	}
	return cachedGTypeObject
}

// gTypeInstanceIsA reports whether a GTypeInstance is of (a subtype of) the
// given GType — the purego equivalent of the WEBKIT_IS_* checking macros.
func gTypeInstanceIsA(instance uintptr, gtype uint64) bool {
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
