//go:build darwin && !ios && !server && (!production || devtools) && privatemacapis

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"

@interface _WKInspector : NSObject
- (void)show;
@end

@interface WKWebView ()
- (_WKInspector *)_inspector;
@end

void openDevTools(void *window) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
	if (@available(macOS 12.0, *)) {
		dispatch_async(dispatch_get_main_queue(), ^{
			WebviewWindow* nsWindow = (WebviewWindow*)window;
			@try {
				[nsWindow.webView._inspector show];
			} @catch (NSException *exception) {
				NSLog(@"Opening the inspector failed: %@", exception.reason);
			}
		});
	}
#else
	NSLog(@"Opening the inspector needs at least macOS 12");
#endif
}

void windowEnableDevTools(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 130300
	if (@available(macOS 13.3, *)) {
		window.webView.inspectable = YES;
		return;
	}
#endif
	[window.webView.configuration.preferences setValue:@YES forKey:@"developerExtrasEnabled"];
}
*/
import "C"

func (w *macosWebviewWindow) openDevTools() {
	C.openDevTools(w.nsWindow)
}

func (w *macosWebviewWindow) enableDevTools() {
	C.windowEnableDevTools(w.nsWindow)
}
