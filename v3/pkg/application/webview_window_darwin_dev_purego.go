//go:build darwin && purego && !ios && !server && (!production || devtools)

package application

// enableDevTools turns on the Web Inspector for the window's webview. As in the
// cgo windowEnableDevTools, macOS 13.3+ uses the public inspectable property;
// older systems fall back to wailsPrivateEnableWebInspector, which writes the
// private developerExtrasEnabled preference in private builds and does nothing
// in public builds.
func (w *macosWebviewWindow) enableDevTools() {
	runOnMain(func() {
		webView := w.webview()
		if macOSAtLeast(13, 3) && respondsTo(webView, "setInspectable:") {
			webView.send("setInspectable:", true)
			return
		}
		wailsPrivateEnableWebInspector(webView)
	})
}

// openDevTools shows the Web Inspector through wailsPrivateOpenWebInspector,
// which is a no-op in public builds.
func (w *macosWebviewWindow) openDevTools() {
	runOnMain(func() {
		wailsPrivateOpenWebInspector(w.webview())
	})
}
