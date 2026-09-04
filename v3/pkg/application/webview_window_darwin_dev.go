//go:build darwin && !ios && !server && (!production || devtools)

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"
#include "mac_private_api_darwin.h"

// The implementations live in mac_private_api_devtools_darwin.go and
// mac_private_api_devtools_appstore_darwin.go, so that the private WebKit
// inspector API stays out of `-tags appstore` builds.
*/
import "C"

func (w *macosWebviewWindow) openDevTools() {
	C.wailsPrivateOpenWebInspector(w.nsWindow)
}

func (w *macosWebviewWindow) enableDevTools() {
	C.wailsPrivateEnableWebInspector(w.nsWindow)
}
