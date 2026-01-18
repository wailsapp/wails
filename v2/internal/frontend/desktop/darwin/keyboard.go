//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "WailsContext.h"

extern void processMessage(const char *message);

void setupEscapeKeyMonitor(void *inctx) {
	WailsContext *ctx = (__bridge WailsContext*) inctx;

	[NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskKeyDown handler:^NSEvent * _Nullable(NSEvent * _Nonnull event) {
		// Check if this is the Escape key (keyCode 53)
		if (event.keyCode == 53) {
			// Check if the window is in fullscreen mode
			if ([ctx IsFullScreen]) {
				// Dispatch escape key event to the WebView so JavaScript can handle it
				// This allows modal dialogs to close without exiting fullscreen
				NSString *jsCode = @"window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', keyCode: 27, which: 27, bubbles: true, cancelable: true }));";
				[ctx ExecJS:jsCode];

				// Consume the event to prevent macOS from exiting fullscreen
				return nil;
			}
		}
		// Allow normal processing for all other cases
		return event;
	}];
}
*/
import "C"
import (
	"unsafe"
)

func setupEscapeKeyMonitor(context unsafe.Pointer) {
	C.setupEscapeKeyMonitor(context)
}
