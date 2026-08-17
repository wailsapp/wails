//go:build darwin && !ios && !server

#import "native_window_darwin.h"
#import <objc/runtime.h>

@interface WailsNativeWindowDelegate : NSObject <NSWindowDelegate>
@property unsigned int windowID;
@property BOOL hideOnClose;
// Kept for compatibility with the existing toolbar chrome helper, which
// stores this preference on the window delegate.
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
@end

static const void* WailsNativeWindowDelegateAssociationKey = &WailsNativeWindowDelegateAssociationKey;

void* nativeWindowCreate(unsigned int windowID, int width, int height, bool hideOnClose) {
    NSWindowStyleMask style = NSWindowStyleMaskTitled |
        NSWindowStyleMaskClosable |
        NSWindowStyleMaskMiniaturizable |
        NSWindowStyleMaskResizable;
    NSWindow* window = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, width, height)
        styleMask:style backing:NSBackingStoreBuffered defer:NO];
    if (window == nil) return NULL;
    window.releasedWhenClosed = NO;
    window.backgroundColor = [NSColor windowBackgroundColor];
    window.opaque = YES;

    WailsNativeWindowDelegate* delegate = [[WailsNativeWindowDelegate alloc] init];
    delegate.windowID = windowID;
    delegate.hideOnClose = hideOnClose;
    window.delegate = delegate;
    objc_setAssociatedObject(window, WailsNativeWindowDelegateAssociationKey,
        delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    [delegate release];
    return window;
}

void nativeWindowDestroy(void* nsWindow) {
    NSWindow* window = (NSWindow*)nsWindow;
    if (window == nil) return;
    window.delegate = nil;
    objc_setAssociatedObject(window, WailsNativeWindowDelegateAssociationKey,
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

void nativeWindowSetAlwaysOnTop(void* nsWindow, bool alwaysOnTop) {
    ((NSWindow*)nsWindow).level = alwaysOnTop ? NSFloatingWindowLevel : NSNormalWindowLevel;
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
