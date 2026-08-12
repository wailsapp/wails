//go:build darwin && !ios && !server && privatemacapis

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#include "webview_window_darwin.h"

// Keep every private macOS API used by the normal window implementation in
// this opt-in translation unit so default builds cannot accidentally link it.
void wailsSetWebViewTransparent(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	[window.webView setValue:@NO forKey:@"drawsBackground"];
	[window.webView setValue:[NSColor clearColor] forKey:@"backgroundColor"];
}

void wailsSetWebViewBackgroundColour(void* nsWindow, int r, int g, int b, int alpha) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	NSColor* color = [NSColor colorWithRed:r/255.0
		green:g/255.0
		blue:b/255.0
		alpha:alpha/255.0];
	[window.webView setValue:color forKey:@"backgroundColor"];
}

void wailsConfigurePrivateLiquidGlass(void* nativeGlassView, int style, const char* groupID, double groupSpacing) {
	NSView* glassView = (NSView*)nativeGlassView;
	if ([glassView respondsToSelector:@selector(setStyle:)]) {
		// Preserve the legacy Wails mapping in opt-in builds. Values above 1
		// are undocumented NSGlassEffectView styles.
		NSInteger nativeStyle = style == 3 ? 1 : style;
		[glassView setValue:@(nativeStyle) forKey:@"style"];
	}
	if (groupID && groupID[0] != '\0') {
		NSString* identifier = [NSString stringWithUTF8String:groupID];
		if ([glassView respondsToSelector:@selector(setGroupIdentifier:)]) {
			[glassView performSelector:@selector(setGroupIdentifier:) withObject:identifier];
		} else if ([glassView respondsToSelector:@selector(setGroupName:)]) {
			[glassView performSelector:@selector(setGroupName:) withObject:identifier];
		}
	}
	if (groupSpacing > 0 && [glassView respondsToSelector:@selector(setGroupSpacing:)]) {
		[glassView setValue:@(groupSpacing) forKey:@"groupSpacing"];
	}
}
*/
import "C"
