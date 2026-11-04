//go:build darwin && purego && !ios && !server

package application

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

// windowShouldUnconditionallyClose is the pure-Go entry point for the
// WindowShouldClose check.
func windowShouldUnconditionallyClose(windowId uint) bool {
	return checkUnconditionalClose(windowId)
}
