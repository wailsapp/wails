//go:build darwin

package application

/*
#include <stdbool.h>
*/
import "C"

//export windowShouldUnconditionallyClose
func windowShouldUnconditionallyClose(windowId C.uint) C.bool {
	window, _ := globalApplication.Window.GetByID(uint(windowId))
	if window == nil {
		globalApplication.debug("windowShouldUnconditionallyClose: window not found", "windowId", windowId)
		return C.bool(false)
	}
	unconditionallyClose := window.shouldUnconditionallyClose()
	globalApplication.debug("windowShouldUnconditionallyClose check", "windowId", windowId, "unconditionallyClose", unconditionallyClose)
	return C.bool(unconditionallyClose)
}
