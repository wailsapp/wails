//go:build darwin && !ios && !server && !appstore

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include <string.h>

#include "webview_window_darwin.h"
#include "mac_private_api_darwin.h"

// This file is the only place in Wails that references undocumented macOS APIs.
// Build with `-tags appstore` to compile mac_private_api_appstore_darwin.go
// instead, which provides the same functions using public APIs or no-ops.

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

// WebKit exposes its feature flags - the same list Safari shows under
// Develop > Feature Flags - through +[WKPreferences _features]. Every feature
// is an object with a `key` matching the name in WebKit's
// UnifiedWebPreferences.yaml.
@interface WKPreferences (WailsPrivateFeatures)
+ (NSArray *)_features;
- (void)_setEnabled:(BOOL)value forFeature:(id)feature;
@end

void wailsPrivateSetPreferPageRenderingUpdatesNear60FPS(void *wkPreferences, bool enabled) {
	WKPreferences *preferences = (WKPreferences *)wkPreferences;
	if (preferences == nil) {
		return;
	}
	if (![[WKPreferences class] respondsToSelector:@selector(_features)] ||
		![preferences respondsToSelector:@selector(_setEnabled:forFeature:)]) {
		NSLog(@"[Wails] PreferPageRenderingUpdatesNear60FPS is unavailable: this WebKit does not expose feature flags");
		return;
	}

	@try {
		for (id feature in [WKPreferences _features]) {
			if (![feature respondsToSelector:@selector(valueForKey:)]) {
				continue;
			}
			NSString *key = [feature valueForKey:@"key"];
			if ([key isEqualToString:@"PreferPageRenderingUpdatesNear60FPSEnabled"]) {
				[preferences _setEnabled:(enabled ? YES : NO) forFeature:feature];
				return;
			}
		}
		NSLog(@"[Wails] PreferPageRenderingUpdatesNear60FPS is unavailable: WebKit no longer defines this feature");
	} @catch (NSException *exception) {
		NSLog(@"[Wails] Could not set PreferPageRenderingUpdatesNear60FPS: %@", exception.reason);
	}
}
*/
import "C"
