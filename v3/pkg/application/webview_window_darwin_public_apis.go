//go:build darwin && !ios && !server && noprivateapis

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit -framework QuartzCore

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"

void wailsSetWebViewTransparent(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
	if (@available(macOS 12.0, *)) {
		// This controls the documented colour below the page but does not make
		// WKWebView itself transparent. AppKit has no public API for that.
		window.webView.underPageBackgroundColor = [NSColor clearColor];
	}
#endif
}

void wailsSetWebViewBackgroundColour(void* nsWindow, int r, int g, int b, int alpha) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	NSColor* color = [NSColor colorWithRed:r/255.0
		green:g/255.0
		blue:b/255.0
		alpha:alpha/255.0];
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
	if (@available(macOS 12.0, *)) {
		window.webView.underPageBackgroundColor = color;
		return;
	}
#endif
	window.webView.wantsLayer = YES;
	window.webView.layer.backgroundColor = color.CGColor;
}

void wailsConfigurePrivateLiquidGlass(void* glassView, int style, const char* groupID, double groupSpacing) {
	(void)glassView;
	(void)style;
	(void)groupID;
	(void)groupSpacing;
}
*/
import "C"
