//go:build darwin && !ios && !server && !purego

package application

/*
#include <stdbool.h>
*/
import "C"

//export windowIsHidden
func windowIsHidden(windowId C.uint) C.bool {
	window, _ := globalApplication.Window.GetByID(uint(windowId))
	if window == nil {
		return C.bool(false)
	}
	webviewWindow, ok := window.(*WebviewWindow)
	if !ok {
		return C.bool(false)
	}
	return C.bool(webviewWindow.options.Hidden)
}

// windowShouldUnconditionallyClose is the //export entry point for the
// WindowShouldClose check.
//
//export windowShouldUnconditionallyClose
func windowShouldUnconditionallyClose(windowId C.uint) C.bool {
	return C.bool(checkUnconditionalClose(uint(windowId)))
}
