//go:build darwin

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
			DevTools:    application.Enabled,
			Javascript:  application.Enabled,
			WebSecurity: application.Disabled, // Disable CORS
			ZoomFactor:  1.2,
		},
	}

	// Because calling NewWebContentsView invokes C.createWebContentsView which 
	// traps without a runloop during go test, we will just manually instantiate
	// the Go wrapper to verify the methods.
	view := &WebContentsView{
		id:      1,
		options: options,
		impl:    &mockWebContentsViewImpl{}, // Mock the impl to bypass Objective-C in headless test
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

type mockWebContentsViewImpl struct{}

func (m *mockWebContentsViewImpl) setBounds(bounds application.Rect) {}
func (m *mockWebContentsViewImpl) setURL(url string) {}
func (m *mockWebContentsViewImpl) execJS(js string) {}
func (m *mockWebContentsViewImpl) attach(window application.Window) {}
func (m *mockWebContentsViewImpl) detach() {}
func (m *mockWebContentsViewImpl) nativeView() unsafe.Pointer { return nil }
