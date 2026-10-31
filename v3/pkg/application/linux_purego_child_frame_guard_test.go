//go:build linux && purego && !gtk3 && !android && !server

package application

import (
	"os"
	"strings"
	"testing"
)

// TestPuregoWebviewRegistersChildFrameGuard checks that the GTK4 purego
// windowNewWebview registers the child-frame guard exactly as linux_cgo.go
// does: a user script over linuxChildFrameGuardJS injected into every frame at
// document start (the same scope as the Blob body shim next to it), added to
// the content manager before the "external" message handler is registered.
//
// The registration is checked at the source level because the GTK4 bindings
// expose no way to read a content manager's scripts back: WebKitGTK 6.0 does
// not declare webkit_user_content_manager_get_scripts() or
// webkit_user_script_get_source() in WebKitUserContent.h, and reading WebKit's
// internal structs would pin private layouts.
func TestPuregoWebviewRegistersChildFrameGuard(t *testing.T) {
	source, err := os.ReadFile("linux_purego.go")
	if err != nil {
		t.Fatalf("read linux_purego.go: %v", err)
	}
	fn := functionSource(t, string(source), "func windowNewWebview(")

	// The cgo GTK4 backend, the reference implementation, uses the same calls.
	cgoSource, err := os.ReadFile("linux_cgo.go")
	if err != nil {
		t.Fatalf("read linux_cgo.go: %v", err)
	}
	cgoFn := functionSource(t, string(cgoSource), "func windowNewWebview(")
	if !strings.Contains(cgoFn, "C.webkit_user_content_manager_add_script(manager, childFrameGuard)") {
		t.Fatal("linux_cgo.go no longer registers childFrameGuard; this test's reference is stale")
	}

	// WEBKIT_USER_CONTENT_INJECT_ALL_FRAMES and
	// WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_START are both 0 in WebKitGTK 6.0.
	for _, want := range []string{
		"childFrameGuard := webkit_user_script_new(linuxChildFrameGuardJS, 0, 0, 0, 0)",
		"webkit_user_content_manager_add_script(manager, childFrameGuard)",
		"webkit_user_script_unref(childFrameGuard)",
	} {
		if !strings.Contains(fn, want) {
			t.Errorf("windowNewWebview does not contain %q", want)
		}
	}

	// Order: guard script before the message handler that it wraps, so no
	// message of a document that starts loading can bypass it.
	guard := strings.Index(fn, "childFrameGuard := webkit_user_script_new(")
	handler := strings.Index(fn, `register_script_message_handler(manager, "external"`)
	switch {
	case guard < 0 || handler < 0:
		t.Fatal("guard script or message handler registration missing")
	case guard > handler:
		t.Error("child-frame guard is registered after the external message handler")
	}
}

// TestPuregoChildFrameGuardMatchesCgo pins the purego registration to the cgo
// one, keyed by function and message, so a future diverging edit to either
// backend fails this test rather than silently changing child-frame behaviour.
func TestPuregoChildFrameGuardMatchesCgo(t *testing.T) {
	puregoSource, err := os.ReadFile("linux_purego.go")
	if err != nil {
		t.Fatalf("read linux_purego.go: %v", err)
	}
	cgoSource, err := os.ReadFile("linux_cgo.go")
	if err != nil {
		t.Fatalf("read linux_cgo.go: %v", err)
	}
	puregoFn := functionSource(t, string(puregoSource), "func windowNewWebview(")
	cgoFn := functionSource(t, string(cgoSource), "func windowNewWebview(")

	for _, call := range []string{
		"webkit_user_content_manager_new()",
		"webkit_user_script_new(linuxBlobBodyFetchShimJS, 0, 0, 0, 0)",
		"webkit_user_script_new(linuxChildFrameGuardJS, 0, 0, 0, 0)",
		`register_script_message_handler(manager, "external"`,
	} {
		if !strings.Contains(puregoFn, call) {
			t.Errorf("purego windowNewWebview is missing %q", call)
		}
	}
	for _, call := range []string{
		"linuxBlobBodyFetchShimJS",
		"linuxChildFrameGuardJS",
		"WEBKIT_USER_CONTENT_INJECT_ALL_FRAMES",
		"WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_START",
	} {
		if !strings.Contains(cgoFn, call) {
			t.Errorf("cgo windowNewWebview is missing %q; update the purego port", call)
		}
	}
}

// functionSource returns the text of the named top-level function in src.
func functionSource(t *testing.T, src, signature string) string {
	t.Helper()
	start := strings.Index(src, signature)
	if start < 0 {
		t.Fatalf("%q not found", signature)
	}
	rest := src[start:]
	if end := strings.Index(rest, "\n}\n"); end >= 0 {
		return rest[:end]
	}
	return rest
}
