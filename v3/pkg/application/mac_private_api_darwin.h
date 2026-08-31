//go:build darwin && !ios && !server

#ifndef WailsMacPrivateAPI_h
#define WailsMacPrivateAPI_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include <stdbool.h>

// Every undocumented macOS API that Wails calls lives behind one of these
// functions. Two mutually exclusive implementations are compiled:
//
//   mac_private_api_darwin.go            default builds, real private APIs
//   mac_private_api_appstore_darwin.go   -tags appstore, public APIs or no-ops
//
// Nothing outside those two files may reference a private selector or key, so
// an `appstore` build contains no reference to them at all. This is enforced by
// TestPrivateMacAPIsAreIsolated in mac_private_api_test.go.

// Make the WKWebView itself draw no background, so the window (or a backdrop
// view behind it) shows through. macOS has no public equivalent.
void wailsPrivateSetWebviewTransparent(void *webView);

// Set the WKWebView's background colour. The appstore build falls back to the
// documented underPageBackgroundColor.
void wailsPrivateSetWebviewBackgroundColour(void *webView, int r, int g, int b, int alpha);

// Apply a Wails Liquid Glass style to an NSGlassEffectView. The default build
// keeps Wails' historic mapping, which passes undocumented style values; the
// appstore build uses only the two documented styles plus NSAppearance.
void wailsPrivateSetGlassStyle(void *glassView, int style);

// Group glass views across windows. There is no public AppKit equivalent, so
// the appstore build ignores it.
void wailsPrivateSetGlassGrouping(void *glassView, const char *groupID, double groupSpacing);

// Toggle WebKit's PreferPageRenderingUpdatesNear60FPSEnabled feature on a
// WKPreferences. Disabling it lets the webview render at the display's native
// refresh rate (120Hz on ProMotion) instead of being clamped to 60. There is no
// public API for this, so the appstore build is a no-op. See issue #6056.
void wailsPrivateSetPreferPageRenderingUpdatesNear60FPS(void *wkPreferences, bool enabled);

// Web Inspector. Implemented by the devtools-only variants:
//
//   mac_private_api_devtools_darwin.go
//   mac_private_api_devtools_appstore_darwin.go
void wailsPrivateOpenWebInspector(void *window);
void wailsPrivateEnableWebInspector(void *window);

#endif /* WailsMacPrivateAPI_h */
