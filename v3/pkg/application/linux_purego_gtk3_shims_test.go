//go:build linux && purego && gtk3 && !android && !server

package application

import (
	"strings"
	"testing"
	"unsafe"
)

// The GTK3 purego shims stand in for the cgo helpers in linux_cgo_gtk3.go. The
// tests below cover the pure helpers among them; the rest need a live GTK3
// window and are exercised by the GTK3 window tests.

func TestGTK3CStringRoundTrip(t *testing.T) {
	ensureGTK3Libraries()
	for _, s := range []string{"", "gtk3", "ünïcode ✓", "a longer string with spaces"} {
		p := _Cfunc_CString(s)
		if p == nil {
			t.Fatalf("_Cfunc_CString(%q) returned nil", s)
		}
		if got := _Cfunc_GoString(p); got != s {
			t.Errorf("_Cfunc_GoString(_Cfunc_CString(%q)) = %q", s, got)
		}
	}
	if got := _Cfunc_GoString(nil); got != "" {
		t.Fatalf("_Cfunc_GoString(nil) = %q", got)
	}
}

func TestGTK3CompiledVersionsMatchRuntime(t *testing.T) {
	ensureGTK3Libraries()
	if got, want := _Cfunc_get_compiled_gtk_major_version(), _Cfunc_gtk_get_major_version(); got != want {
		t.Fatalf("compiled GTK major %d != runtime %d", got, want)
	}
	if got, want := _Cfunc_get_compiled_webkit_major_version(), _Cfunc_webkit_get_major_version(); got != want {
		t.Fatalf("compiled WebKit major %d != runtime %d", got, want)
	}
	if _Cfunc_get_compiled_gtk_major_version() != 3 {
		t.Fatalf("GTK3 build loaded GTK major %d", _Cfunc_get_compiled_gtk_major_version())
	}
}

func TestGTK3SelectionAtomsAreDistinct(t *testing.T) {
	ensureGTK3Libraries()
	clipboard := _Cmacro_GDK_SELECTION_CLIPBOARD()
	primary := _Cmacro_GDK_SELECTION_PRIMARY()
	if clipboard == 0 || primary == 0 || clipboard == primary {
		t.Fatalf("selection atoms clipboard=%#x primary=%#x", clipboard, primary)
	}
}

func TestGTK3EnsureLibrariesIsIdempotent(t *testing.T) {
	ensureGTK3Libraries()
	ensureGTK3Libraries()
	if gtk3LibsErr != nil {
		t.Fatal(gtk3LibsErr)
	}
	if gtk3Idle == nil || gtk3SetData == nil {
		t.Fatal("GTK3 function table not bound")
	}
}

// mustNotPanic fails the test with the panic value instead of aborting the run.
func mustNotPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("loader panicked: %v", r)
		}
	}()
	fn()
}

func TestGTK3LoadMissingLibraryReturnsActionableError(t *testing.T) {
	const missingLib = "libwails-no-such-gtk3-test.so.0"
	var err error
	mustNotPanic(t, func() {
		err = loadGTK3([]string{missingLib}, func(uintptr, *[]string) {
			t.Error("binder must not run when no library loads")
		})
	})
	if err == nil {
		t.Fatal("loadGTK3 with a missing library returned nil")
	}
	for _, want := range []string{"GTK3 and WebKit2GTK 4.1", missingLib, "libwebkit2gtk-4.1-0", "apt install"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
}

func TestGTK3LoadMissingSymbolReturnsActionableError(t *testing.T) {
	var local func(uintptr) unsafe.Pointer
	const missingSym = "wails_no_such_gtk3_symbol_test"
	var err error
	mustNotPanic(t, func() {
		err = loadGTK3(gtk3LibraryNames, func(lib uintptr, missing *[]string) {
			bindGTK3Symbol(missing, lib, &local, missingSym)
		})
	})
	if err == nil {
		t.Fatal("loadGTK3 with a missing symbol returned nil")
	}
	if !strings.Contains(err.Error(), missingSym) {
		t.Errorf("error %q does not name the missing symbol", err)
	}
	if local != nil {
		t.Error("a missing symbol must leave its function variable unset")
	}
}

func TestGTK3LoadReportsEveryMissingSymbolSorted(t *testing.T) {
	var a func(uintptr) unsafe.Pointer
	var b func(uintptr)
	err := loadGTK3(gtk3LibraryNames, func(lib uintptr, missing *[]string) {
		bindGTK3Symbol(missing, lib, &a, "wails_zz_missing_test")
		bindGTK3Symbol(missing, lib, &b, "wails_aa_missing_test")
	})
	if err == nil {
		t.Fatal("expected an error for two missing symbols")
	}
	if !strings.Contains(err.Error(), "wails_aa_missing_test, wails_zz_missing_test") {
		t.Errorf("error %q does not list both symbols sorted", err)
	}
}

func TestGTK3LoadSucceedsWithPresentLibraryAndSymbol(t *testing.T) {
	var malloc func(uintptr) unsafe.Pointer
	var err error
	mustNotPanic(t, func() {
		err = loadGTK3(gtk3LibraryNames, func(lib uintptr, missing *[]string) {
			bindGTK3Symbol(missing, lib, &malloc, "g_malloc")
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if malloc == nil {
		t.Fatal("g_malloc was not bound")
	}
	if p := malloc(16); p == nil {
		t.Fatal("bound g_malloc returned nil")
	}
}

func TestGTK3ThreadSelfIsNilWhenLibrariesFailed(t *testing.T) {
	saved := gtk3GThreadSelf
	defer func() { gtk3GThreadSelf = saved }()
	gtk3GThreadSelf = nil
	if got := _Cfunc_g_thread_self(); got != nil {
		t.Fatalf("_Cfunc_g_thread_self() = %p with no library bound, want nil", got)
	}
}
