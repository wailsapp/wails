//go:build windows && !server

package webcontentsview

import (
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestWindowsWebContentsViewDefersControllerOperations(t *testing.T) {
	view := NewWebContentsView(WebContentsViewOptions{
		URL:    "https://example.com",
		Bounds: application.Rect{X: 10, Y: 20, Width: 640, Height: 480},
	})
	impl, ok := view.impl.(*windowsWebContentsView)
	if !ok {
		t.Fatalf("implementation = %T, want *windowsWebContentsView", view.impl)
	}
	if impl.chromium == nil {
		t.Fatal("WebView2 wrapper was not initialized")
	}
	if impl.chromium.IsReady() {
		t.Fatal("WebView2 controller unexpectedly initialized without a host window")
	}

	// These calls are valid before Attach: requested URL/bounds stay on the
	// parent view, and scripts wait until the native controller is ready.
	view.SetBounds(application.Rect{X: 30, Y: 40, Width: 800, Height: 600})
	view.SetURL("https://wails.io")
	view.ExecJS("window.__wailsProbe = true")
	view.GoBack()

	if got := view.GetURL(); got != "https://wails.io" {
		t.Fatalf("GetURL() = %q, want requested URL", got)
	}
	if got := len(impl.pendingJS); got != 1 {
		t.Fatalf("queued scripts = %d, want 1", got)
	}
}

func TestWindowsWebContentsViewUsesStableIsolatedUserDataPath(t *testing.T) {
	newPath := func(name string) string {
		view := NewWebContentsView(WebContentsViewOptions{Name: name})
		impl, ok := view.impl.(*windowsWebContentsView)
		if !ok {
			t.Fatalf("implementation = %T, want *windowsWebContentsView", view.impl)
		}
		return impl.chromium.DataPath
	}

	first := newPath("agent-browser")
	if first == "" {
		t.Fatal("child WebView2 user-data path must not fall back to the host default")
	}
	if filepath.Clean(first) != first {
		t.Fatalf("user-data path is not clean: %q", first)
	}
	if again := newPath("agent-browser"); again != first {
		t.Fatalf("same view name has unstable user-data paths: first=%q second=%q", first, again)
	}
	if other := newPath("agent-browser-2"); other == first {
		t.Fatalf("different views share user-data path %q", first)
	}
}
