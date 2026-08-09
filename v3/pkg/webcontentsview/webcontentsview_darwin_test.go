//go:build darwin && !server

package webcontentsview

import (
	"testing"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Dummy mock window that satisfies the interface
// and returns a nil NativeWindow so we can test the Attach nil-handling safely
// without spinning up the full NSApplication runloop in a headless test environment.
type mockWindow struct {
	application.Window
}

func (m *mockWindow) NativeWindow() unsafe.Pointer {
	return nil
}

func TestWebContentsView_APISurface(t *testing.T) {
	// We primarily want to ensure that the API surface compiles and functions
	// correctly at a struct level. Note: Full WKWebView instantiation without an NSApplication
	// runloop will crash on macOS, so we test the struct wiring here instead of the native allocations.

	options := WebContentsViewOptions{
		Name: "TestBrowser",
		URL:  "https://example.com",
		Bounds: application.Rect{
			X:      0,
			Y:      0,
			Width:  800,
			Height: 600,
		},
		WebPreferences: WebPreferences{
			DevTools:    Enabled,
			Javascript:  Enabled,
			WebSecurity: Disabled, // Disable CORS
			ZoomFactor:  1.2,
		},
	}

	// Native allocation is deferred to Attach after the host window exists, so
	// construction remains safe in a headless test without an NSApplication.
	view := NewWebContentsView(options)
	impl, ok := view.impl.(*macosWebContentsView)
	if !ok {
		t.Fatalf("implementation = %T, want *macosWebContentsView", view.impl)
	}
	if impl.nsView != nil {
		t.Fatal("WKWebView was created before a native window attached")
	}

	// 2. Test SetBounds
	view.SetBounds(application.Rect{X: 10, Y: 10, Width: 400, Height: 400})

	// 3. Test SetURL
	view.SetURL("https://google.com")

	// 4. Test ExecJS
	view.ExecJS("console.log('test');")

	// 5. Test Attach and Detach using a mock window
	win := &mockWindow{}
	view.Attach(win)
	view.Detach()

	t.Log("macOS WebContentsView API surface tests passed successfully.")
}
