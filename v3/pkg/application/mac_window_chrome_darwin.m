//go:build darwin && !ios && !server

#import "mac_window_chrome_darwin.h"
#import "mac_private_api_darwin.h"

void windowSetToolbarStyle(void* pointer, int style) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil || window.toolbar == nil) return;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) window.toolbarStyle = (NSWindowToolbarStyle)style;
#endif
}

void windowSetHideToolbarSeparator(void* pointer, bool hideSeparator) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil || window.toolbar == nil) return;
    window.toolbar.showsBaselineSeparator = !hideSeparator;
}

void windowSetShowToolbarWhenFullscreen(void* pointer, bool setting) {
    NSWindow* window = (NSWindow*)pointer;
    id delegate = window.delegate;
    if (delegate != nil && [delegate respondsToSelector:@selector(setShowToolbarWhenFullscreen:)]) {
        [delegate setValue:@(setting) forKey:@"showToolbarWhenFullscreen"];
    }
}

// The functions below apply MacWindow options to a plain NSWindow. Each one
// mirrors the WebviewWindow implementation of the same option in
// webview_window_darwin.go or webview_window_darwin.m; keep the two in step.

void windowChromeSetShadow(void* pointer, bool hasShadow) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil) return;
    window.hasShadow = hasShadow;
}

void windowChromeSetLevel(void* pointer, int levelCode) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil) return;
    switch (levelCode) {
    case 1: window.level = NSFloatingWindowLevel; break;
    case 2: window.level = NSTornOffMenuWindowLevel; break;
    case 3: window.level = NSModalPanelWindowLevel; break;
    case 4: window.level = NSMainMenuWindowLevel; break;
    case 5: window.level = NSStatusWindowLevel; break;
    case 6: window.level = NSPopUpMenuWindowLevel; break;
    case 7: window.level = NSScreenSaverWindowLevel; break;
    default: window.level = NSNormalWindowLevel; break;
    }
}

void windowChromeSetCollectionBehavior(void* pointer, int behavior) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil) return;
    if (behavior == 0) {
        window.collectionBehavior = NSWindowCollectionBehaviorFullScreenPrimary;
    } else {
        window.collectionBehavior = (NSWindowCollectionBehavior)behavior;
    }
}

void windowChromeSetTabbingMode(void* pointer, int mode) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil) return;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101200
    if (@available(macOS 10.12, *)) window.tabbingMode = (NSWindowTabbingMode)mode;
#endif
}

void windowChromeSetAppearanceByName(void* pointer, const char* appearanceName) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil || appearanceName == NULL || appearanceName[0] == '\0') return;
    window.appearance = [NSAppearance appearanceNamed:[NSString stringWithUTF8String:appearanceName]];
}

void windowChromeSetTransparentSurface(void* pointer) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil) return;
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
}

void windowChromeSetNormalSurface(void* pointer) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil) return;
    NSColor* background = window.backgroundColor;
    if (background == nil || background.alphaComponent <= 0.0) {
        background = [NSColor windowBackgroundColor];
    }
    window.backgroundColor = background;
    window.opaque = YES;
}

bool windowChromeLiquidGlassSupported(void) {
    if (@available(macOS 26.0, *)) {
        return NSClassFromString(@"NSGlassEffectView") != nil;
    }
    return false;
}

void* windowChromeCreateTranslucentView(void) {
    NSVisualEffectView* effectView = [[[NSVisualEffectView alloc] initWithFrame:NSZeroRect] autorelease];
    effectView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    effectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
    effectView.state = NSVisualEffectStateActive;
    return effectView;
}

// Style values match the Go MacLiquidGlassStyle constants (and the
// MacLiquidGlassStyle enum in webview_window_darwin.h, which is not included
// here because that header is excluded from wails_native builds).
enum {
    WailsChromeGlassStyleAutomatic = 0,
    WailsChromeGlassStyleLight = 1,
    WailsChromeGlassStyleDark = 2,
    WailsChromeGlassStyleVibrant = 3,
};

void* windowChromeCreateLiquidGlassView(int style, int material, double cornerRadius,
    int r, int g, int b, int a, const char* groupID, double groupSpacing) {
    if (!windowChromeLiquidGlassSupported()) return nil;
    NSView* glassView = nil;
    if (@available(macOS 26.0, *)) {
        Class glassClass = NSClassFromString(@"NSGlassEffectView");
        if (glassClass != nil) {
            glassView = [[[glassClass alloc] initWithFrame:NSZeroRect] autorelease];
            if (cornerRadius > 0 && [glassView respondsToSelector:@selector(setCornerRadius:)]) {
                [glassView setValue:@(cornerRadius) forKey:@"cornerRadius"];
            }
            if (a > 0 && [glassView respondsToSelector:@selector(setTintColor:)]) {
                NSColor* tint = [NSColor colorWithRed:r / 255.0 green:g / 255.0 blue:b / 255.0 alpha:a / 255.0];
                [glassView performSelector:@selector(setTintColor:) withObject:tint];
            }
            // Only regular and clear are documented styles, and grouping has no
            // public API, so both go through the private_mac_apis-guarded
            // implementations in mac_private_api_darwin.go / mac_public_api_darwin.go.
            wailsPrivateSetGlassStyle((void*)glassView, style);
            wailsPrivateSetGlassGrouping((void*)glassView, groupID, groupSpacing);
        }
    }
    if (glassView == nil) {
        // Same fallback material table as windowSetLiquidGlass.
        NSVisualEffectView* effectView = [[[NSVisualEffectView alloc] initWithFrame:NSZeroRect] autorelease];
        glassView = effectView;
        if (material >= 0) {
            effectView.material = (NSVisualEffectMaterial)material;
            effectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
        } else {
            switch (style) {
            case WailsChromeGlassStyleLight:
                if (@available(macOS 15.0, *)) effectView.material = NSVisualEffectMaterialUnderPageBackground;
                else if (@available(macOS 10.14, *)) effectView.material = NSVisualEffectMaterialHUDWindow;
                else effectView.material = NSVisualEffectMaterialLight;
                effectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
                break;
            case WailsChromeGlassStyleDark:
                if (@available(macOS 15.0, *)) effectView.material = NSVisualEffectMaterialHeaderView;
                else if (@available(macOS 10.14, *)) effectView.material = NSVisualEffectMaterialFullScreenUI;
                else effectView.material = NSVisualEffectMaterialDark;
                effectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
                break;
            case WailsChromeGlassStyleVibrant:
                if (@available(macOS 11.0, *)) effectView.material = NSVisualEffectMaterialHUDWindow;
                else if (@available(macOS 10.14, *)) effectView.material = NSVisualEffectMaterialSheet;
                else effectView.material = NSVisualEffectMaterialLight;
                effectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
                break;
            default:
                if (@available(macOS 10.14, *)) effectView.material = NSVisualEffectMaterialContentBackground;
                else effectView.material = NSVisualEffectMaterialAppearanceBased;
                effectView.blendingMode = NSVisualEffectBlendingModeWithinWindow;
                break;
            }
        }
        effectView.state = NSVisualEffectStateFollowsWindowActiveState;
        if (@available(macOS 10.12, *)) effectView.emphasized = NO;
        if (cornerRadius > 0) {
            effectView.wantsLayer = YES;
            effectView.layer.cornerRadius = cornerRadius;
            effectView.layer.masksToBounds = YES;
        }
    }
    glassView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    return glassView;
}
