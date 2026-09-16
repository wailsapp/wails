//go:build darwin && !ios && !server

#import "native_window_darwin.h"
#import "mac_window_chrome_darwin.h"
#import <objc/runtime.h>

@interface WailsNativeWindowDelegate : NSObject <NSWindowDelegate>
@property unsigned int windowID;
@property BOOL hideOnClose;
// Set by windowSetShowToolbarWhenFullscreen (TitleBar.ShowToolbarWhenFullscreen)
// through the shared toolbar chrome helper.
@property BOOL showToolbarWhenFullscreen;
@end

@implementation WailsNativeWindowDelegate
- (BOOL)windowShouldClose:(NSWindow*)sender {
    if (!self.hideOnClose) return YES;
    [sender orderOut:nil];
    return NO;
}
- (void)windowWillClose:(NSNotification*)notification {
    processNativeWindowClosed(self.windowID);
}
// Same fullscreen toolbar policy as WebviewWindowDelegate.
- (NSApplicationPresentationOptions)window:(NSWindow*)window
    willUseFullScreenPresentationOptions:(NSApplicationPresentationOptions)proposedOptions {
    if (self.showToolbarWhenFullscreen) return proposedOptions;
    return proposedOptions | NSApplicationPresentationAutoHideToolbar;
}
@end

// WailsNativeWindow and WailsNativePanel are the WebView-free counterparts of
// WebviewWindow and WebviewPanel. A borderless NSWindow refuses key status by
// default, so both override the responder policies the same way the WebView
// classes do, and both honour DisableEscapeExitsFullscreen in cancelOperation:.
@interface WailsNativeWindow : NSWindow
@property BOOL disableEscapeExitsFullscreen;
@end

@implementation WailsNativeWindow
- (BOOL)canBecomeKeyWindow { return YES; }
- (BOOL)canBecomeMainWindow { return YES; }
- (BOOL)acceptsFirstResponder { return YES; }
- (void)cancelOperation:(id)sender {
    if (self.disableEscapeExitsFullscreen &&
        (self.styleMask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen) {
        return;
    }
    [super cancelOperation:sender];
}
@end

@interface WailsNativePanel : NSPanel
@property BOOL disableEscapeExitsFullscreen;
@end

@implementation WailsNativePanel
- (BOOL)canBecomeKeyWindow { return YES; }
- (BOOL)canBecomeMainWindow { return NO; }
- (BOOL)acceptsFirstResponder { return YES; }
- (void)cancelOperation:(id)sender {
    if (self.disableEscapeExitsFullscreen &&
        (self.styleMask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen) {
        return;
    }
    [super cancelOperation:sender];
}
@end

static const void* WailsNativeWindowDelegateAssociationKey = &WailsNativeWindowDelegateAssociationKey;
static const void* WailsNativeWindowBackdropAssociationKey = &WailsNativeWindowBackdropAssociationKey;

// Same style mask rules as windowStyleMask in webview_window_darwin.go, with
// TitleBar.Hide standing in for WebviewWindowOptions.Frameless.
static NSWindowStyleMask nativeWindowStyleMask(WailsNativeWindowConfig config) {
    NSWindowStyleMask style = NSWindowStyleMaskTitled |
        NSWindowStyleMaskClosable |
        NSWindowStyleMaskMiniaturizable |
        NSWindowStyleMaskResizable;
    if (config.frameless && config.borderless) {
        style = NSWindowStyleMaskBorderless | NSWindowStyleMaskResizable | NSWindowStyleMaskMiniaturizable;
    } else if (config.frameless) {
        style |= NSWindowStyleMaskFullSizeContentView;
    }
    if (config.isPanel) {
        if (config.nonActivating) style |= NSWindowStyleMaskNonactivatingPanel;
        if (config.utilityWindow) style |= NSWindowStyleMaskUtilityWindow;
    }
    return style;
}

void* nativeWindowCreate(WailsNativeWindowConfig config) {
    NSWindowStyleMask style = nativeWindowStyleMask(config);
    NSRect contentRect = NSMakeRect(0, 0, config.width, config.height);
    NSWindow* window = nil;
    if (config.isPanel) {
        WailsNativePanel* panel = [[WailsNativePanel alloc] initWithContentRect:contentRect
            styleMask:style backing:NSBackingStoreBuffered defer:NO];
        if (panel == nil) return NULL;
        // NSPanel defaults differ from NSWindow; match WebviewPanel.
        panel.hidesOnDeactivate = NO;
        panel.floatingPanel = config.floatingPanel;
        panel.becomesKeyOnlyIfNeeded = config.becomesKeyOnlyIfNeeded;
        panel.disableEscapeExitsFullscreen = config.disableEscapeExitsFullscreen;
        window = panel;
    } else {
        WailsNativeWindow* plain = [[WailsNativeWindow alloc] initWithContentRect:contentRect
            styleMask:style backing:NSBackingStoreBuffered defer:NO];
        if (plain == nil) return NULL;
        plain.disableEscapeExitsFullscreen = config.disableEscapeExitsFullscreen;
        window = plain;
    }
    window.releasedWhenClosed = NO;
    window.backgroundColor = [NSColor windowBackgroundColor];
    window.opaque = YES;
    if (config.frameless) {
        // WebviewWindow is movable by its background in every configuration;
        // a native window only needs it once the titlebar is gone.
        window.movableByWindowBackground = YES;
        if (!config.borderless) {
            // AppKit's own rounded frame without a visible titlebar: hide the
            // title and the window buttons exactly like a frameless WebviewWindow.
            window.titlebarAppearsTransparent = YES;
            window.titleVisibility = NSWindowTitleHidden;
            [window standardWindowButton:NSWindowCloseButton].hidden = YES;
            [window standardWindowButton:NSWindowMiniaturizeButton].hidden = YES;
            [window standardWindowButton:NSWindowZoomButton].hidden = YES;
        }
    }

    WailsNativeWindowDelegate* delegate = [[WailsNativeWindowDelegate alloc] init];
    delegate.windowID = config.windowID;
    delegate.hideOnClose = config.hideOnClose;
    window.delegate = delegate;
    objc_setAssociatedObject(window, WailsNativeWindowDelegateAssociationKey,
        delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    [delegate release];
    return window;
}

// The primary pane of a native window is an NSScrollView hosting an
// NSTextView. Both draw an opaque background by default, which would hide any
// backdrop behind them; a WebviewWindow makes its WKWebView transparent for the
// same reason. Sidebar, inspector, and content-list panes keep their own
// material, as they do in a split WebviewWindow.
static void nativeWindowMakeTextEditorsTransparent(NSView* host) {
    for (NSView* subview in host.subviews) {
        if ([subview isKindOfClass:[NSScrollView class]]) {
            NSScrollView* scroll = (NSScrollView*)subview;
            if ([scroll.documentView isKindOfClass:[NSTextView class]]) {
                scroll.drawsBackground = NO;
                ((NSTextView*)scroll.documentView).drawsBackground = NO;
            }
            continue;
        }
        nativeWindowMakeTextEditorsTransparent(subview);
    }
}

// Inserts the backdrop beneath every pane. The content view of an installed
// native window is the NSSplitViewController's NSSplitView; it must not treat
// the backdrop as a pane, so subview arrangement is limited to its existing
// arranged subviews before the backdrop is added.
static void nativeWindowInsertBackdrop(NSWindow* window, NSView* backdrop) {
    NSView* host = window.contentView;
    if (host == nil || backdrop == nil) return;
    if ([host isKindOfClass:[NSSplitView class]]) {
        NSSplitView* split = (NSSplitView*)host;
        if (split.arrangesAllSubviews) {
            NSArray<NSView*>* panes = [split.arrangedSubviews copy];
            split.arrangesAllSubviews = NO;
            for (NSView* pane in panes) {
                if (![split.arrangedSubviews containsObject:pane]) [split addArrangedSubview:pane];
            }
            [panes release];
        }
    }
    backdrop.frame = host.bounds;
    backdrop.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    [host addSubview:backdrop positioned:NSWindowBelow relativeTo:nil];
    objc_setAssociatedObject(window, WailsNativeWindowBackdropAssociationKey,
        backdrop, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
}

void nativeWindowApplyContentChrome(void* nsWindow, bool belowToolbar,
    double cornerRadius, WailsNativeBackdropConfig backdrop) {
    NSWindow* window = (NSWindow*)nsWindow;
    if (window == nil) return;
    NSView* contentView = window.contentView;

    // ContentLayout. The split installer always enables full-size content;
    // below-toolbar layout removes it again so the whole layout, including a
    // full-height sidebar, starts beneath the titlebar and toolbar.
    if (belowToolbar) {
        NSRect frame = window.frame;
        window.styleMask &= ~NSWindowStyleMaskFullSizeContentView;
        [window setFrame:frame display:YES];
    }

    // CornerRadius on a borderless window: mask the content the way
    // createNativeWindow masks the WebviewWindow content view.
    if (cornerRadius > 0 && contentView != nil) {
        contentView.wantsLayer = YES;
        contentView.layer.cornerRadius = cornerRadius;
        contentView.layer.masksToBounds = YES;
    }

    // Backdrop. Mirrors the switch in macosWebviewWindow.run.
    NSView* backdropView = nil;
    switch (backdrop.backdrop) {
    case 1: // MacBackdropTransparent
        windowChromeSetTransparentSurface(window);
        nativeWindowMakeTextEditorsTransparent(contentView);
        break;
    case 2: // MacBackdropTranslucent
        windowChromeSetTransparentSurface(window);
        backdropView = (NSView*)windowChromeCreateTranslucentView();
        nativeWindowMakeTextEditorsTransparent(contentView);
        break;
    case 3: // MacBackdropLiquidGlass
        backdropView = (NSView*)windowChromeCreateLiquidGlassView(backdrop.glassStyle,
            backdrop.glassMaterial, backdrop.glassCornerRadius,
            backdrop.tintR, backdrop.tintG, backdrop.tintB, backdrop.tintA,
            backdrop.groupID, backdrop.groupSpacing);
        if (backdropView == nil) backdropView = (NSView*)windowChromeCreateTranslucentView();
        windowChromeSetTransparentSurface(window);
        nativeWindowMakeTextEditorsTransparent(contentView);
        break;
    default: // MacBackdropNormal
        if (cornerRadius > 0) {
            // A rounded borderless window needs a clear surface for the mask to
            // show, so the opaque window background moves into the content.
            windowChromeSetTransparentSurface(window);
            NSBox* fill = [[[NSBox alloc] initWithFrame:NSZeroRect] autorelease];
            fill.boxType = NSBoxCustom;
            fill.borderWidth = 0;
            fill.fillColor = [NSColor windowBackgroundColor];
            fill.titlePosition = NSNoTitle;
            backdropView = fill;
        } else {
            windowChromeSetNormalSurface(window);
        }
        break;
    }
    if (backdropView != nil) nativeWindowInsertBackdrop(window, backdropView);
    [window invalidateShadow];
}

void nativeWindowDestroy(void* nsWindow) {
    NSWindow* window = (NSWindow*)nsWindow;
    if (window == nil) return;
    window.delegate = nil;
    objc_setAssociatedObject(window, WailsNativeWindowDelegateAssociationKey,
        nil, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(window, WailsNativeWindowBackdropAssociationKey,
        nil, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    [window orderOut:nil];
    [window close];
    [window release];
}

void nativeWindowSetTitle(void* nsWindow, const char* title) {
    ((NSWindow*)nsWindow).title = title == NULL ? @"" : [NSString stringWithUTF8String:title];
}

void nativeWindowSetResizable(void* nsWindow, bool resizable) {
    NSWindow* window = (NSWindow*)nsWindow;
    if (resizable) window.styleMask |= NSWindowStyleMaskResizable;
    else window.styleMask &= ~NSWindowStyleMaskResizable;
}

void nativeWindowSetMinSize(void* nsWindow, int width, int height) {
    ((NSWindow*)nsWindow).contentMinSize = NSMakeSize(MAX(0, width), MAX(0, height));
}

void nativeWindowSetMaxSize(void* nsWindow, int width, int height) {
    CGFloat maxWidth = width > 0 ? width : CGFLOAT_MAX;
    CGFloat maxHeight = height > 0 ? height : CGFLOAT_MAX;
    ((NSWindow*)nsWindow).contentMaxSize = NSMakeSize(maxWidth, maxHeight);
}

void nativeWindowConfigureTitlebar(void* nsWindow, bool appearsTransparent,
    bool fullSizeContent, bool hideTitle, bool hideToolbarSeparator, int toolbarStyle) {
    NSWindow* window = (NSWindow*)nsWindow;
    window.titlebarAppearsTransparent = appearsTransparent;
    window.titleVisibility = hideTitle ? NSWindowTitleHidden : NSWindowTitleVisible;
    if (fullSizeContent) window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    else window.styleMask &= ~NSWindowStyleMaskFullSizeContentView;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) window.toolbarStyle = (NSWindowToolbarStyle)toolbarStyle;
#endif
    if ([window respondsToSelector:@selector(setTitlebarSeparatorStyle:)]) {
        [window setValue:@(hideToolbarSeparator ? 0 : 1) forKey:@"titlebarSeparatorStyle"];
    }
}

void nativeWindowShow(void* nsWindow) {
    [(NSWindow*)nsWindow makeKeyAndOrderFront:nil];
    [NSApp activateIgnoringOtherApps:YES];
}
void nativeWindowHide(void* nsWindow) { [(NSWindow*)nsWindow orderOut:nil]; }
void nativeWindowFocus(void* nsWindow) { [(NSWindow*)nsWindow makeKeyAndOrderFront:nil]; }
bool nativeWindowIsVisible(void* nsWindow) { return ((NSWindow*)nsWindow).visible; }
void nativeWindowCenter(void* nsWindow) { [(NSWindow*)nsWindow center]; }
void nativeWindowSetPosition(void* nsWindow, int x, int y) {
    NSWindow* window = (NSWindow*)nsWindow;
    NSRect frame = window.frame;
    NSScreen* screen = window.screen ?: [NSScreen mainScreen];
    CGFloat top = NSMaxY(screen.visibleFrame);
    [window setFrameOrigin:NSMakePoint(x, top - y - frame.size.height)];
}
