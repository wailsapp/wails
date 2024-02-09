//go:build darwin && (!production || devtools)

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

#include "webview_window_darwin.h"

@interface _WKInspector : NSObject
- (void)show;
- (void)detach;
@end

@interface WKWebView ()
- (_WKInspector *)_inspector;
@end

//void showDevTools(void *window) {
//    // get main window
//    WebviewWindow* nsWindow = (WebviewWindow*)window;
//    dispatch_async(dispatch_get_main_queue(), ^{
//        [nsWindow.webView._inspector show];
//    });
//}

void showDevTools(void *window) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
    dispatch_async(dispatch_get_main_queue(), ^{
		if (@available(macOS 12.0, *)) {
			WebviewWindow* nsWindow = (WebviewWindow*)window;

			@try {
				[nsWindow.webView._inspector show];
			} @catch (NSException *exception) {
				NSLog(@"Opening the inspector failed: %@", exception.reason);
				return;
			}
		} else {
			NSLog(@"Opening the inspector needs at least MacOS 12");
		}
    });
#endif
}

// Enable NSWindow devtools
void windowEnableDevTools(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	// Enable devtools in webview
	[window.webView.configuration.preferences setValue:@YES forKey:@"developerExtrasEnabled"];
}

*/
import "C"
import "unsafe"

func init() {
	showDevTools = func(window unsafe.Pointer) {
		C.showDevTools(window)
	}
}

func (w *macosWebviewWindow) enableDevTools() {
	C.windowEnableDevTools(w.nsWindow)
}
