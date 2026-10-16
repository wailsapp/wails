//go:build darwin && purego && !ios && !server

package application

// windowShouldUnconditionallyClose reports whether the window may close without
// firing the WindowShouldClose event (Go-native port of the cgo //export).
func windowShouldUnconditionallyClose(windowId uint) bool {
	window, _ := globalApplication.Window.GetByID(windowId)
	if window == nil {
		globalApplication.debug("windowShouldUnconditionallyClose: window not found", "windowId", windowId)
		return false
	}
	unconditionallyClose := window.shouldUnconditionallyClose()
	globalApplication.debug("windowShouldUnconditionallyClose check", "windowId", windowId, "unconditionallyClose", unconditionallyClose)
	return unconditionallyClose
}

// windowIsHidden reports whether the window was configured hidden.
func windowIsHidden(windowId uint) bool {
	window, _ := globalApplication.Window.GetByID(windowId)
	if window == nil {
		return false
	}
	webviewWindow, ok := window.(*WebviewWindow)
	if !ok {
		return false
	}
	return webviewWindow.options.Hidden
}
