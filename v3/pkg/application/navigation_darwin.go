//go:build darwin && !ios && !server

package application

/*
#include <stdbool.h>

void navigationDecide(void *decisionHandler, bool allow);
*/
import "C"

import "unsafe"

// windowInterceptsNavigation lets the WKNavigationDelegate in
// webview_window_darwin.m allow a navigation on the spot when the window has
// no OnBeforeNavigate, which is the common case.
//
//export windowInterceptsNavigation
func windowInterceptsNavigation(windowID C.uint) C.bool {
	return C.bool(interceptsNavigation(uint(windowID)))
}

// processNavigationAction receives a navigation WebKit is holding for the
// window's OnBeforeNavigate. decisionHandler is a copy of WebKit's decision
// block; navigationDecide calls it on the main thread and releases it.
//
//export processNavigationAction
func processNavigationAction(windowID C.uint, url *C.char, isMainFrame C.bool, decisionHandler unsafe.Pointer) {
	request := NavigationRequest{
		URL:         C.GoString(url),
		IsMainFrame: bool(isMainFrame),
	}
	resolveNavigation(uint(windowID), request, func(allow bool) {
		C.navigationDecide(decisionHandler, C.bool(allow))
	})
}
