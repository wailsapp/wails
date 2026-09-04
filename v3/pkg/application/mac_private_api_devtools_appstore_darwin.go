//go:build darwin && !ios && !server && (!production || devtools) && appstore

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"
#include "mac_private_api_darwin.h"

// Public-only Web Inspector support, selected by `-tags appstore`.

// WebKit has no public API for opening the inspector from code. The window
// stays inspectable, so Safari's Develop menu can still attach to it.
void wailsPrivateOpenWebInspector(void *window) {
	(void)window;
	NSLog(@"[Wails] Opening the Web Inspector from code needs a private WebKit API and is disabled in appstore builds. Attach from Safari's Develop menu instead.");
}

void wailsPrivateEnableWebInspector(void *window) {
	NSWindow<WailsWebviewWindow> *nsWindow = (NSWindow<WailsWebviewWindow> *)window;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 130300
	if (@available(macOS 13.3, *)) {
		nsWindow.webView.inspectable = YES;
		return;
	}
#endif
	NSLog(@"[Wails] The Web Inspector needs macOS 13.3 or later in appstore builds.");
}
*/
import "C"
