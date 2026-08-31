//go:build darwin && !ios && !server && appstore

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"
#include "mac_private_api_darwin.h"

// Public-API replacements for every private macOS call in
// mac_private_api_darwin.go, selected by `-tags appstore`. Where macOS offers
// no public equivalent the function is a documented no-op: the Go API is
// unchanged, only the visual result differs.

// WKWebView cannot be made transparent with public API. underPageBackgroundColor
// (macOS 12+) is the closest documented control: it tints the area behind the
// page, so overscroll and unpainted regions blend with the window, but the page
// itself stays opaque.
void wailsPrivateSetWebviewTransparent(void *webView) {
	WKWebView *view = (WKWebView *)webView;
	if (view == nil) {
		return;
	}
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
	if (@available(macOS 12.0, *)) {
		view.underPageBackgroundColor = [NSColor clearColor];
	}
#endif
}

void wailsPrivateSetWebviewBackgroundColour(void *webView, int r, int g, int b, int alpha) {
	WKWebView *view = (WKWebView *)webView;
	if (view == nil) {
		return;
	}
	NSColor *colour = [NSColor colorWithRed:r / 255.0 green:g / 255.0 blue:b / 255.0 alpha:alpha / 255.0];
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
	if (@available(macOS 12.0, *)) {
		view.underPageBackgroundColor = colour;
		return;
	}
#endif
	view.wantsLayer = YES;
	view.layer.backgroundColor = colour.CGColor;
}

// NSGlassEffectView documents two styles: regular (0) and clear (1). Wails'
// light and dark styles have no native counterpart, so they are expressed with
// NSAppearance instead.
void wailsPrivateSetGlassStyle(void *glassView, int style) {
	NSView *view = (NSView *)glassView;
	if ([view respondsToSelector:@selector(setStyle:)]) {
		NSInteger nativeStyle = (style == LiquidGlassStyleVibrant) ? 1 : 0;
		[view setValue:@(nativeStyle) forKey:@"style"];
	}
	if (style == LiquidGlassStyleLight) {
		view.appearance = [NSAppearance appearanceNamed:NSAppearanceNameAqua];
	} else if (style == LiquidGlassStyleDark) {
		view.appearance = [NSAppearance appearanceNamed:NSAppearanceNameDarkAqua];
	} else {
		view.appearance = nil;
	}
}

// Cross-window glass grouping has no public AppKit equivalent. Each glass view
// stays independent.
void wailsPrivateSetGlassGrouping(void *glassView, const char *groupID, double groupSpacing) {
	(void)glassView;
	(void)groupID;
	(void)groupSpacing;
}

// WebKit's 60fps rendering cadence can only be changed through private feature
// flags, so an appstore build stays at 60fps on high-refresh displays.
void wailsPrivateSetPreferPageRenderingUpdatesNear60FPS(void *wkPreferences, bool enabled) {
	(void)wkPreferences;
	if (enabled) {
		// Requesting WebKit's default needs no action.
		return;
	}
	static dispatch_once_t onceToken;
	dispatch_once(&onceToken, ^{
		NSLog(@"[Wails] PreferPageRenderingUpdatesNear60FPS is ignored in appstore builds: unlocking the display refresh rate requires a private WebKit API");
	});
}
*/
import "C"
