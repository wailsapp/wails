//go:build darwin && !ios && !server && (!production || devtools) && !appstore

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"
#include "mac_private_api_darwin.h"

// Web Inspector support for default (non-appstore) builds. macOS has no public
// API for opening the inspector programmatically, so this file uses WebKit's
// private _inspector. See mac_private_api_devtools_appstore_darwin.go for the
// public-only variant.

@interface _WKInspector : NSObject
- (void)show;
- (void)detach;
@end

@interface WKWebView (WailsPrivateInspector)
- (_WKInspector *)_inspector;
@end

void wailsPrivateOpenWebInspector(void *window) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
	if (@available(macOS 12.0, *)) {
		dispatch_async(dispatch_get_main_queue(), ^{
			NSWindow<WailsWebviewWindow> *nsWindow = (NSWindow<WailsWebviewWindow> *)window;
			@try {
				[nsWindow.webView._inspector show];
			} @catch (NSException *exception) {
				NSLog(@"Opening the inspector failed: %@", exception.reason);
				return;
			}
		});
	}
#else
	NSLog(@"Opening the inspector needs at least MacOS 12");
#endif
}

void wailsPrivateEnableWebInspector(void *window) {
	NSWindow<WailsWebviewWindow> *nsWindow = (NSWindow<WailsWebviewWindow> *)window;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 130300
	if (@available(macOS 13.3, *)) {
		nsWindow.webView.inspectable = YES;
	}
#endif
	// Older WebKit versions only honour the private preference.
	@try {
		[nsWindow.webView.configuration.preferences setValue:@YES forKey:@"developerExtrasEnabled"];
	} @catch (NSException *exception) {
		NSLog(@"[Wails] Could not enable developer extras: %@", exception.reason);
	}
}
*/
import "C"
