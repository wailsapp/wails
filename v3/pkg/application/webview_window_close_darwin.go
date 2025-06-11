//go:build darwin

package application

/*
#include <stdbool.h>
*/
import "C"

//export windowShouldUnconditionallyClose
func windowShouldUnconditionallyClose(windowId uint) bool {
	window := globalApplication.getWindowForID(windowId)
	if window == nil {
		globalApplication.debug("windowShouldUnconditionallyClose: window not found", "windowId", windowId)
		return false
	}
	webviewWindow, ok := window.(*WebviewWindow)
	if !ok {
		globalApplication.debug("windowShouldUnconditionallyClose: window is not WebviewWindow", "windowId", windowId)
		return false
	}
	globalApplication.debug("windowShouldUnconditionallyClose check", "windowId", windowId, "unconditionallyClose", webviewWindow.unconditionallyClose)
	return webviewWindow.unconditionallyClose
}
