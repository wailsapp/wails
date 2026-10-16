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
// private -[WKWebView _inspector] selector (macOS 12+). Both selectors are
// PRIVATE API and can vanish in any WebKit update; a missing selector raises
// an uncatchable NSException, so guard every send (cgo wraps this in both
// @available(macOS 12,*) and @try/@catch).
func (w *macosWebviewWindow) openDevTools() {
	runOnMain(func() {
		if !respondsTo(w.webview(), "_inspector") {
			return
		}
		inspector := w.webview().send("_inspector")
		if !inspector.isNil() && respondsTo(inspector, "show") {
			inspector.send("show")
		}
	})
}
