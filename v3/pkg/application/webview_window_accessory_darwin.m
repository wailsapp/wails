//go:build darwin && !ios && !server

#import "webview_window_accessory_darwin.h"

enum {
    WailsMacScrollEdgeEffectStyleAutomatic = 0,
    WailsMacScrollEdgeEffectStyleSoft = 1,
    WailsMacScrollEdgeEffectStyleHard = 2,
};

static bool isSupportedAccessoryViewController(id controller) {
    if (controller == nil) return false;
    if ([controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return true;
    Class splitAccessoryClass = NSClassFromString(@"NSSplitViewItemAccessoryViewController");
    return splitAccessoryClass != Nil && [controller isKindOfClass:splitAccessoryClass];
}

static int accessoryViewControllerKind(id controller) {
    if (controller == nil) return 0;
    if ([controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return 1;
    Class splitAccessoryClass = NSClassFromString(@"NSSplitViewItemAccessoryViewController");
    if (splitAccessoryClass != Nil && [controller isKindOfClass:splitAccessoryClass]) return 2;
    return 0;
}

int macAccessoryViewControllerKind(void* pointer) {
    if ([NSThread isMainThread]) return accessoryViewControllerKind((id)pointer);
    __block int result = 0;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessoryViewControllerKind((id)pointer);
    });
    return result;
}

static bool accessorySupportsScrollEdgeEffectStyle(id controller) {
    if (!isSupportedAccessoryViewController(controller)) return false;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260100
    if (@available(macOS 26.1, *)) return true;
#endif
    return false;
}

bool macAccessoryViewControllerSupportsScrollEdgeEffectStyle(void* pointer) {
    if ([NSThread isMainThread]) return accessorySupportsScrollEdgeEffectStyle((id)pointer);
    __block bool result = false;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessorySupportsScrollEdgeEffectStyle((id)pointer);
    });
    return result;
}

static int accessorySetScrollEdgeEffectStyle(id controller, int style) {
    if (!isSupportedAccessoryViewController(controller)) {
        return WailsMacAccessoryStyleInvalidController;
    }
    if (style < WailsMacScrollEdgeEffectStyleAutomatic ||
        style > WailsMacScrollEdgeEffectStyleHard) {
        return WailsMacAccessoryStyleInvalidController;
    }
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260100
    if (@available(macOS 26.1, *)) {
        NSScrollEdgeEffectStyle* nativeStyle = NSScrollEdgeEffectStyle.automaticStyle;
        if (style == WailsMacScrollEdgeEffectStyleSoft) {
            nativeStyle = NSScrollEdgeEffectStyle.softStyle;
        } else if (style == WailsMacScrollEdgeEffectStyleHard) {
            nativeStyle = NSScrollEdgeEffectStyle.hardStyle;
        }
        [(NSTitlebarAccessoryViewController*)controller setPreferredScrollEdgeEffectStyle:nativeStyle];
        return WailsMacAccessoryStyleApplied;
    }
#endif
    return style == WailsMacScrollEdgeEffectStyleAutomatic
        ? WailsMacAccessoryStyleApplied
        : WailsMacAccessoryStyleUnavailable;
}

int macAccessoryViewControllerSetScrollEdgeEffectStyle(void* pointer, int style) {
    if ([NSThread isMainThread]) return accessorySetScrollEdgeEffectStyle((id)pointer, style);
    __block int result = WailsMacAccessoryStyleInvalidController;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessorySetScrollEdgeEffectStyle((id)pointer, style);
    });
    return result;
}

static int accessoryScrollEdgeEffectStyle(id controller) {
    if (!isSupportedAccessoryViewController(controller)) {
        return WailsMacAccessoryStyleInvalidController;
    }
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260100
    if (@available(macOS 26.1, *)) {
        NSScrollEdgeEffectStyle* style =
            [(NSTitlebarAccessoryViewController*)controller preferredScrollEdgeEffectStyle];
        if ([style isEqual:NSScrollEdgeEffectStyle.softStyle]) {
            return WailsMacScrollEdgeEffectStyleSoft;
        }
        if ([style isEqual:NSScrollEdgeEffectStyle.hardStyle]) {
            return WailsMacScrollEdgeEffectStyleHard;
        }
        return WailsMacScrollEdgeEffectStyleAutomatic;
    }
#endif
    return WailsMacAccessoryStyleUnavailable;
}

int macAccessoryViewControllerScrollEdgeEffectStyle(void* pointer) {
    if ([NSThread isMainThread]) return accessoryScrollEdgeEffectStyle((id)pointer);
    __block int result = WailsMacAccessoryStyleInvalidController;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessoryScrollEdgeEffectStyle((id)pointer);
    });
    return result;
}
