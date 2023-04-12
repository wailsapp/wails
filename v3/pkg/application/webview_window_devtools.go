//go:build darwin && !production

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

#include "webview_window.h"

@interface _WKInspector : NSObject
- (void)show;
- (void)detach;
@end

@interface WKWebView ()
- (_WKInspector *)_inspector;
@end

void showDevTools(void *window) {
    // get main window
    WebviewWindow* nsWindow = (WebviewWindow*)window;
    dispatch_async(dispatch_get_main_queue(), ^{
        [nsWindow.webView._inspector show];
    });
}

*/
import "C"
import "unsafe"

func init() {
	showDevTools = func(window unsafe.Pointer) {
		C.showDevTools(window)
	}
}
