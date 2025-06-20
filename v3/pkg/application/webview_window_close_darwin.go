//go:build darwin

package application

/*
#include <stdbool.h>
*/
import "C"
import "sync/atomic"

//export windowShouldUnconditionallyClose
func windowShouldUnconditionallyClose(windowId C.uint) C.bool {
	window, _ := globalApplication.Windows.GetByID(uint(windowId))
	if window == nil {
		globalApplication.debug("windowShouldUnconditionallyClose: window not found", "windowId", windowId)
		return C.bool(false)
	}
	webviewWindow, ok := window.(*WebviewWindow)
	if !ok {
		globalApplication.debug("windowShouldUnconditionallyClose: window is not WebviewWindow", "windowId", windowId)
		return C.bool(false)
	}
	unconditionallyClose := atomic.LoadUint32(&webviewWindow.unconditionallyClose) != 0
	globalApplication.debug("windowShouldUnconditionallyClose check", "windowId", windowId, "unconditionallyClose", unconditionallyClose)
	return C.bool(unconditionallyClose)
}
