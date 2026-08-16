//go:build darwin && !ios && !server && (!production || devtools) && noprivateapis

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"

void openDevTools(void *window) {
	(void)window;
	NSLog(@"Programmatically opening the Web Inspector is disabled by the noprivateapis build tag. Use Safari's Develop menu to inspect this window.");
}

// Enable NSWindow devtools
void windowEnableDevTools(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 130300
	if (@available(macOS 13.3, *)) {
		window.webView.inspectable = YES;
		return;
	}
#endif
	NSLog(@"Web Inspector requires macOS 13.3 when built with noprivateapis.");
}

*/
import "C"

func (w *macosWebviewWindow) openDevTools() {
	C.openDevTools(w.nsWindow)
}

func (w *macosWebviewWindow) enableDevTools() {
	C.windowEnableDevTools(w.nsWindow)
}
