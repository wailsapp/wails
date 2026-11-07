//go:build windows && !server && !cef

package application

import (
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

// The edge tests verify that ShuttingDown revokes OLE registration. These
// exercise both application destruction paths against real HWNDs, ensuring
// they actually call ShuttingDown rather than relying on a later Close.
func TestWindowsWindowDestructionShutsDownChromium(t *testing.T) {
	for _, path := range []string{"WM_CLOSE", "destroy"} {
		t.Run(path, func(t *testing.T) {
			w := newTestMenuWindow(t)
			previous := globalApplication
			app := &App{options: Options{Windows: WindowsOptions{DisableQuitOnLastWindowClosed: true}}}
			app.impl = &windowsApp{parent: app, windowMap: map[w32.HWND]*windowsWebviewWindow{w.hwnd: w}}
			globalApplication = app
			t.Cleanup(func() { globalApplication = previous })

			w.chromium = edge.NewChromium()
			// Observe shutdown through its documented callback cleanup, without
			// requiring a WebView2 install or reaching into edge's private fields.
			w.chromium.MessageCallback = func(string, *edge.ICoreWebView2, *edge.ICoreWebView2WebMessageReceivedEventArgs) {}
			w.parent.options.Windows.WebView2CompositionHosting = true
			if path == "WM_CLOSE" {
				atomic.StoreUint32(&w.parent.unconditionallyClose, 1)
				w.WndProc(w32.WM_CLOSE, 0, 0)
			} else {
				w.destroy()
			}
			// Do not call Chromium.Close before either assertion: that would
			// hide a missing shutdown call in the application path under test.
			if w.chromium.MessageCallback != nil {
				t.Fatal("window destruction did not shut down Chromium")
			}
			if w32.IsWindow(w.hwnd) {
				t.Fatal("window destruction left the HWND alive")
			}
		})
	}
}
