//go:build linux && purego && !android

package webview

import "testing"

func TestFormatGErrorPreservesNativeMessage(t *testing.T) {
	ensureWebviewLibs()
	originalFree := g_error_free
	t.Cleanup(func() { g_error_free = originalFree })
	var freed uintptr
	g_error_free = func(pointer uintptr) {
		freed = pointer
		originalFree(pointer)
	}
	nativeError := g_error_new_literal(g_quark_from_string("wails-test"), 1, "failed at /a%20b")
	got := formatGError("stream read failed", nativeError).Error()
	if want := "stream read failed: failed at /a%20b"; got != want {
		t.Fatalf("error = %q; want %q", got, want)
	}
	if freed != nativeError {
		t.Fatal("native GError was not freed")
	}
}
