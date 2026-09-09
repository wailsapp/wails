//go:build darwin && !ios && !server && private_mac_apis

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include <string.h>

#include "webview_window_darwin.h"
#include "mac_private_api_darwin.h"

// This file is the only place in Wails v3 that calls undocumented macOS APIs.
// It is compiled only with -tags private_mac_apis. Default builds select
// mac_public_api_darwin.go, whose private-only operations are no-ops.

// WKWebView has no public API for drawing without a background. This key is
// backed by a private property that WebKit has honoured since WebKit1.
void wailsPrivateSetWebviewTransparent(void *webView) {
	WKWebView *view = (WKWebView *)webView;
	if (view == nil) {
		return;
	}
	@try {
		[view setValue:@NO forKey:@"drawsBackground"];
	} @catch (NSException *exception) {
		NSLog(@"[Wails] Could not make the webview transparent: %@", exception.reason);
	}
}

void wailsPrivateSetWebviewBackgroundColour(void *webView, int r, int g, int b, int alpha) {
	WKWebView *view = (WKWebView *)webView;
	if (view == nil) {
		return;
	}
	NSColor *colour = [NSColor colorWithRed:r / 255.0 green:g / 255.0 blue:b / 255.0 alpha:alpha / 255.0];
	@try {
		[view setValue:colour forKey:@"backgroundColor"];
	} @catch (NSException *exception) {
		NSLog(@"[Wails] Could not set the webview background colour: %@", exception.reason);
	}
}

// Wails' historic mapping: vibrant borrows the light style, and light/dark are
// passed through as raw style values even though only regular (0) and clear (1)
// are documented.
void wailsPrivateSetGlassStyle(void *glassView, int style) {
	NSView *view = (NSView *)glassView;
	if (![view respondsToSelector:@selector(setStyle:)]) {
		return;
	}
	int nativeStyle = (style == LiquidGlassStyleVibrant) ? LiquidGlassStyleLight : style;
	@try {
		[view setValue:@(nativeStyle) forKey:@"style"];
	} @catch (NSException *exception) {
		NSLog(@"[Wails] Could not set the glass style: %@", exception.reason);
	}
}

void wailsPrivateSetGlassGrouping(void *glassView, const char *groupID, double groupSpacing) {
	NSView *view = (NSView *)glassView;
	if (groupID && strlen(groupID) > 0) {
		NSString *identifier = [NSString stringWithUTF8String:groupID];
		if ([view respondsToSelector:@selector(setGroupIdentifier:)]) {
			[view performSelector:@selector(setGroupIdentifier:) withObject:identifier];
		} else if ([view respondsToSelector:@selector(setGroupName:)]) {
			[view performSelector:@selector(setGroupName:) withObject:identifier];
		}
	}
	if (groupSpacing > 0 && [view respondsToSelector:@selector(setGroupSpacing:)]) {
		@try {
			[view setValue:@(groupSpacing) forKey:@"groupSpacing"];
		} @catch (NSException *exception) {
			NSLog(@"[Wails] Could not set the glass group spacing: %@", exception.reason);
		}
	}
}

// The devtools bridge defines this for development and production+devtools.
// Plain production builds must not contain private inspector selectors.
#ifdef WAILS_MAC_DEVTOOLS
@interface _WKInspector : NSObject
- (void)show;
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
	NSLog(@"Opening the inspector needs at least macOS 12");
#endif
}

void wailsPrivateEnableWebInspector(void *window) {
	NSWindow<WailsWebviewWindow> *nsWindow = (NSWindow<WailsWebviewWindow> *)window;

	@try {
		[nsWindow.webView.configuration.preferences setValue:@YES forKey:@"developerExtrasEnabled"];
	} @catch (NSException *exception) {
		NSLog(@"[Wails] Could not enable developer extras: %@", exception.reason);
	}
}
#endif
*/
import "C"
