//go:build darwin && !ios && !server

#ifndef WailsMacWindowChromeDarwin_h
#define WailsMacWindowChromeDarwin_h

#import <Cocoa/Cocoa.h>

// Shared NSWindow chrome operations. These deliberately accept NSWindow,
// rather than WebviewWindow, so toolbars work in both normal and native builds.
void windowSetToolbarStyle(void* nsWindow, int style);
void windowSetHideToolbarSeparator(void* nsWindow, bool hideSeparator);
void windowSetShowToolbarWhenFullscreen(void* nsWindow, bool setting);

// Plain NSWindow equivalents of the MacWindow option appliers that
// webview_window_darwin.go implements against WebviewWindow. They exist so a
// WebView-free NativeWindow (including wails_native builds, which exclude the
// WebviewWindow sources) applies each option with the same AppKit semantics.

// windowChromeSetShadow mirrors windowSetShadow.
void windowChromeSetShadow(void* nsWindow, bool hasShadow);
// windowChromeSetLevel takes a nativeMacWindowLevelCode value from
// native_window_mac_options.go, not a raw NSWindowLevel.
void windowChromeSetLevel(void* nsWindow, int levelCode);
// windowChromeSetCollectionBehavior mirrors windowSetCollectionBehavior:
// zero selects NSWindowCollectionBehaviorFullScreenPrimary.
void windowChromeSetCollectionBehavior(void* nsWindow, int behavior);
// windowChromeSetTabbingMode takes an NSWindowTabbingMode value.
void windowChromeSetTabbingMode(void* nsWindow, int mode);
// windowChromeSetAppearanceByName sets the NSAppearance by name. The caller
// keeps ownership of the string.
void windowChromeSetAppearanceByName(void* nsWindow, const char* appearanceName);
// windowChromeSetTransparentSurface mirrors windowSetTransparent;
// windowChromeSetNormalSurface mirrors windowSetNormalBackdrop.
void windowChromeSetTransparentSurface(void* nsWindow);
void windowChromeSetNormalSurface(void* nsWindow);
// windowChromeLiquidGlassSupported mirrors isLiquidGlassSupported.
bool windowChromeLiquidGlassSupported(void);
// windowChromeCreateTranslucentView mirrors the NSVisualEffectView created by
// windowSetTranslucent. windowChromeCreateLiquidGlassView mirrors the glass
// view (NSGlassEffectView, or the NSVisualEffectView fallback with the same
// style-to-material mapping) created by windowSetLiquidGlass. Both return an
// autoreleased NSView that the caller must insert into a view hierarchy; the
// Liquid Glass variant returns nil when glass is unsupported so the caller can
// fall back to the translucent view exactly as WebviewWindow does.
void* windowChromeCreateTranslucentView(void);
void* windowChromeCreateLiquidGlassView(int style, int material, double cornerRadius,
    int r, int g, int b, int a, const char* groupID, double groupSpacing);

#endif
