//go:build darwin && purego && !ios && !server && (!production || devtools)

package application

// enableDevTools turns on the Web Inspector for the window's webview by setting
// the private developerExtrasEnabled preference (KVC), matching the cgo path.
func (w *macosWebviewWindow) enableDevTools() {
	runOnMain(func() {
		prefs := w.webview().send("configuration").send("preferences")
		prefs.send("setValue:forKey:", nsNumberBool(true), nsString("developerExtrasEnabled"))
	})
}

// openDevTools shows the Web Inspector. The inspector is reached through the
// private -[WKWebView _inspector] selector (macOS 12+).
func (w *macosWebviewWindow) openDevTools() {
	runOnMain(func() {
		inspector := w.webview().send("_inspector")
		if !inspector.isNil() {
			inspector.send("show")
		}
	})
}
