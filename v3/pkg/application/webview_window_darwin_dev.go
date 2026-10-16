//go:build darwin && !ios && !server && (!production || devtools)

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c -DWAILS_MAC_DEVTOOLS
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "webview_window_darwin.h"
#include "mac_private_api_darwin.h"

void windowEnableDevTools(void *window) {
    NSWindow<WailsWebviewWindow> *nsWindow = (NSWindow<WailsWebviewWindow> *)window;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 130300
    if (@available(macOS 13.3, *)) {
        nsWindow.webView.inspectable = YES;
        return;
    }
#endif
    wailsPrivateEnableWebInspector(window);
}
*/
import "C"

func (w *macosWebviewWindow) openDevTools() {
	C.wailsPrivateOpenWebInspector(w.nsWindow)
}

func (w *macosWebviewWindow) enableDevTools() {
	C.windowEnableDevTools(w.nsWindow)
}
