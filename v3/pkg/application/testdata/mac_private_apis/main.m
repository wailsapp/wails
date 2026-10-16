#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include <assert.h>
#include <math.h>
#include "mac_private_api_darwin.h"
#include "webview_window_darwin.h"

// Records the actual native messages sent by each bridge implementation.
@interface GlassProbe : NSView
@property NSInteger style;
@property (copy) NSString *groupIdentifier;
@property double groupSpacing;
@end
@implementation GlassProbe
@end

@interface LegacyGlassProbe : NSView
@property (copy) NSString *groupName;
@end
@implementation LegacyGlassProbe
@end

#ifdef WAILS_MAC_DEVTOOLS
void windowEnableDevTools(void *window);
@interface WindowProbe : NSObject
@property (strong) WKWebView *webView;
@end
@implementation WindowProbe
@end
#if EXPECT_PRIVATE
@interface InspectorProbe : NSObject
@property BOOL shown;
- (void)show;
@end
@implementation InspectorProbe
- (void)show { self.shown = YES; }
@end
@interface InspectorViewProbe : NSObject
@property (strong) InspectorProbe *inspector;
- (id)_inspector;
@end
@implementation InspectorViewProbe
- (id)_inspector { return self.inspector; }
@end
#endif
#endif

int main(void) {
    @autoreleasepool {
        [NSApplication sharedApplication];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyProhibited];
        WKWebView *view = [[WKWebView alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
        NSNumber *originalBackground = [view valueForKey:@"drawsBackground"];
        wailsPrivateSetWebviewTransparent((void *)view);
#if EXPECT_PRIVATE
        assert(![[view valueForKey:@"drawsBackground"] boolValue]);
#else
        assert([[view valueForKey:@"drawsBackground"] isEqual:originalBackground]);
        // No-op functions must not even message their receiver.
        NSObject *unsupported = [NSObject new];
        wailsPrivateSetWebviewTransparent((void *)unsupported);
        wailsPrivateSetGlassGrouping((void *)unsupported, "group", 8);
#ifdef WAILS_MAC_DEVTOOLS
        wailsPrivateOpenWebInspector((void *)unsupported);
        wailsPrivateEnableWebInspector((void *)unsupported);
#endif
#endif
        wailsPrivateSetWebviewBackgroundColour((void *)view, 32, 64, 128, 255);
#if EXPECT_PRIVATE
        NSColor *colour = [view valueForKey:@"backgroundColor"];
#else
        NSColor *colour;
        if (@available(macOS 12.0, *)) {
            colour = view.underPageBackgroundColor;
        } else {
            colour = [NSColor colorWithCGColor:view.layer.backgroundColor];
        }
#endif
        colour = [colour colorUsingColorSpace:NSColorSpace.deviceRGBColorSpace];
        assert(fabs(colour.redComponent - 32.0/255.0) < 0.001);
        assert(fabs(colour.greenComponent - 64.0/255.0) < 0.001);
        assert(fabs(colour.blueComponent - 128.0/255.0) < 0.001);
        GlassProbe *glass = [GlassProbe new];
        for (int style = 0; style <= 3; style++) {
            wailsPrivateSetGlassStyle((void *)glass, style);
#if EXPECT_PRIVATE
            assert(glass.style == (style == 3 ? 1 : style));
#else
            assert(glass.style == (style == 3 ? 1 : 0));
            if (style == 1) assert([glass.appearance.name isEqual:NSAppearanceNameAqua]);
            if (style == 2) assert([glass.appearance.name isEqual:NSAppearanceNameDarkAqua]);
            if (style == 0 || style == 3) assert(glass.appearance == nil);
#endif
        }
        wailsPrivateSetGlassGrouping((void *)glass, "shared", 12);
#if EXPECT_PRIVATE
        assert([glass.groupIdentifier isEqual:@"shared"]);
        assert(glass.groupSpacing == 12);
        LegacyGlassProbe *legacy = [LegacyGlassProbe new];
        wailsPrivateSetGlassGrouping((void *)legacy, "legacy", 12);
        assert([legacy.groupName isEqual:@"legacy"]);
#else
        assert(glass.groupIdentifier == nil);
        assert(glass.groupSpacing == 0);
#endif
        wailsPrivateSetGlassGrouping((void *)glass, NULL, 0);
        wailsPrivateSetGlassStyle(NULL, 0);
        wailsPrivateSetGlassGrouping(NULL, NULL, 0);
        wailsPrivateSetWebviewTransparent(NULL);
        wailsPrivateSetWebviewBackgroundColour(NULL, 0, 0, 0, 0);
#ifdef WAILS_MAC_DEVTOOLS
        WindowProbe *window = [WindowProbe new];
        window.webView = view;
        if (@available(macOS 13.3, *)) {
            assert(!view.inspectable);
            windowEnableDevTools((void *)window);
            assert(view.inspectable);
        }
#if EXPECT_PRIVATE
        wailsPrivateEnableWebInspector((void *)window);
        assert([[view.configuration.preferences valueForKey:@"developerExtrasEnabled"] boolValue]);
        if (@available(macOS 12.0, *)) {
            InspectorViewProbe *inspectorView = [InspectorViewProbe new];
            inspectorView.inspector = [InspectorProbe new];
            window.webView = (WKWebView *)inspectorView;
            wailsPrivateOpenWebInspector((void *)window);
            NSDate *deadline = [NSDate dateWithTimeIntervalSinceNow:2];
            while (!inspectorView.inspector.shown && [deadline timeIntervalSinceNow] > 0) {
                [[NSRunLoop mainRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.01]];
            }
            assert(inspectorView.inspector.shown);
        }
#endif
#endif
    }
    return 0;
}
