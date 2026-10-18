//go:build darwin && !ios && !server

package application

/*
#include <stdbool.h>
*/
import "C"

// windowInterceptsNavigation lets the WKNavigationDelegate in
// webview_window_darwin.m allow a navigation on the spot when the window has
// no OnBeforeNavigate, which is the common case.
//
//export windowInterceptsNavigation
func windowInterceptsNavigation(windowID C.uint) C.bool {
	return C.bool(windowInterceptsNavigationByID(uint(windowID)))
}

// processNavigationAction answers a navigation WebKit is holding for the
// window's OnBeforeNavigate. Called on the main thread.
//
//export processNavigationAction
func processNavigationAction(windowID C.uint, url *C.char, isMainFrame C.bool) C.bool {
	return C.bool(decideNavigation(uint(windowID), NavigationRequest{
		URL:         C.GoString(url),
		IsMainFrame: bool(isMainFrame),
	}))
}
