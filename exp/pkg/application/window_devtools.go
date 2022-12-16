//go:build darwin && !production

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

#include "window_delegate.h"

@interface _WKInspector : NSObject
- (void)show;
- (void)detach;
@end

@interface WKWebView ()
- (_WKInspector *)_inspector;
@end

void showDevTools(void *window) {
    // Get the window delegate
    WindowDelegate* delegate = (WindowDelegate*)[(NSWindow*)window delegate];
    dispatch_async(dispatch_get_main_queue(), ^{
        [delegate.webView._inspector show];
        //dispatch_time_t popTime = dispatch_time(DISPATCH_TIME_NOW, 1 * NSEC_PER_SEC);
        //dispatch_after(popTime, dispatch_get_main_queue(), ^(void){
        //    // Detach must be deferred a little bit and is ignored directly after a show.
        //    [delegate.webView._inspector detach];
        //});
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
