//go:build darwin && (dev || debug)

package darwin

// We are using private APIs here, make sure this is only included in a dev/debug build and not in a production build.
// Otherwise the binary might get rejected by the AppReview-Team when pushing it to the AppStore.

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "WailsContext.h"

@interface _WKInspector : NSObject
- (void)show;
- (void)detach;
@end

@interface WKWebView ()
- (_WKInspector *)_inspector;
@end

void showInspector(void *inctx) {
	if (@available(macOS 12.0, *)) {
		WailsContext *ctx = (__bridge WailsContext*) inctx;

		@try {
			[ctx.webview._inspector show];
		} @catch (NSException *exception) {
			NSLog(@"Opening the inspector failed: %@", exception.reason);
			return;
		}

		dispatch_time_t popTime = dispatch_time(DISPATCH_TIME_NOW, 1 * NSEC_PER_SEC);
		dispatch_after(popTime, dispatch_get_main_queue(), ^(void){
			// Detach must be deferred a little bit and is ignored directly after a show.
			@try {
				[ctx.webview._inspector detach];
			} @catch (NSException *exception) {
				NSLog(@"Detaching the inspector failed: %@", exception.reason);
			}
		});
	} else {
		NSLog(@"Opening the inspector needs at least MacOS 12");
	}
}
*/
import "C"
import (
	"unsafe"
)

func showInspector(context unsafe.Pointer) {
	C.showInspector(context)
}
