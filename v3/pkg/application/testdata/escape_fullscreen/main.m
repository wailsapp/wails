#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include <assert.h>
#include <stdio.h>
#include <string.h>
#include "webview_window_darwin.h"

// Exercises windowDropsEscapeCommand, the decision WebviewWindow and
// WebviewPanel make in doCommandBySelector: for DisableEscapeExitsFullscreen.
//
//   --filter                 decide against probe windows; needs no display.
//   --fullscreen MODE        a real fullscreen window holding a WKWebView with a
//                            focused <textarea> receives Escape, and the page
//                            must see the keydown. MODE is one of:
//     --guarded              both overrides, option on: the window must stay
//                            fullscreen, and toggleFullScreen: must still leave.
//     --cancel-only          only the cancelOperation: override, option on: the
//                            window must leave, showing that override alone does
//                            not keep fullscreen.
//     --allow                both overrides, option off: the window must leave.
// Exits non-zero on any failure, and within 20 seconds whatever happens.

typedef NS_ENUM(int, Mode) { ModeGuarded, ModeCancelOnly, ModeAllow };
static Mode mode = ModeGuarded;

// A window as far as windowDropsEscapeCommand reads one.
@interface ProbeWindow : NSObject
@property BOOL disableEscapeExitsFullscreen;
@property NSWindowStyleMask styleMask;
@property (assign) id delegate;
@end
@implementation ProbeWindow
@end

@interface CancellingDelegate : NSObject
- (void)cancelOperation:(id)unused;
@end
@implementation CancellingDelegate
- (void)cancelOperation:(id)unused {}
@end

static BOOL drops(ProbeWindow *window, SEL selector) {
    return windowDropsEscapeCommand((NSWindow<WailsWebviewWindow> *)window, selector);
}

static void testFilter(void) {
    ProbeWindow *window = [ProbeWindow new];
    window.disableEscapeExitsFullscreen = YES;
    window.styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskFullScreen;
    assert(drops(window, @selector(cancelOperation:)));
    // Every other command still reaches NSWindow.
    assert(!drops(window, @selector(cancel:)));
    assert(!drops(window, @selector(insertNewline:)));
    assert(!drops(window, @selector(complete:)));
    // Only while fullscreen, and only when asked for.
    window.styleMask = NSWindowStyleMaskTitled;
    assert(!drops(window, @selector(cancelOperation:)));
    window.styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskFullScreen;
    window.disableEscapeExitsFullscreen = NO;
    assert(!drops(window, @selector(cancelOperation:)));
    // A delegate that handles the command still gets it.
    window.disableEscapeExitsFullscreen = YES;
    CancellingDelegate *delegate = [CancellingDelegate new];
    window.delegate = delegate;
    assert(!drops(window, @selector(cancelOperation:)));
}

// The method bodies below are copies of WebviewWindow's and WebviewPanel's;
// TestEscapeFullscreenHarnessMatchesWindowClasses fails if they diverge.

// The cancelOperation: override alone.
@interface CancelGuardWindow : NSWindow <WailsWebviewWindow>
@property (assign) WKWebView *webView;
@property BOOL disableEscapeExitsFullscreen;
@end
@implementation CancelGuardWindow
- (BOOL)canBecomeKeyWindow {
    return YES;
}
- (void)cancelOperation:(id)sender {
    if (self.disableEscapeExitsFullscreen &&
        (self.styleMask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen) {
        return;
    }
    [super cancelOperation:sender];
}
@end

// Both overrides, as WebviewWindow and WebviewPanel have them.
@interface GuardedWindow : CancelGuardWindow
@end
@implementation GuardedWindow
- (void)doCommandBySelector:(SEL)selector {
    if (windowDropsEscapeCommand(self, selector)) {
        return;
    }
    [super doCommandBySelector:selector];
}
@end

@interface Driver : NSObject <NSApplicationDelegate, NSWindowDelegate, WKScriptMessageHandler>
@property (retain) CancelGuardWindow *window;
@property BOOL pageSawEscape;
@property BOOL leaving;
@end

static BOOL isFullscreen(NSWindow *window) {
    return (window.styleMask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen;
}

__attribute__((noreturn)) static void fail(const char *message) {
    fprintf(stderr, "FAIL: %s\n", message);
    exit(1);
}

__attribute__((noreturn)) static void pass(const char *message) {
    printf("ok: %s\n", message);
    exit(0);
}

static void after(double seconds, dispatch_block_t block) {
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(seconds * NSEC_PER_SEC)), dispatch_get_main_queue(), block);
}

@implementation Driver
- (void)userContentController:(WKUserContentController *)controller didReceiveScriptMessage:(WKScriptMessage *)message {
    if ([[message.body description] isEqualToString:@"Escape"]) {
        self.pageSawEscape = YES;
    }
}

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
    WKWebViewConfiguration *config = [WKWebViewConfiguration new];
    [config.userContentController addScriptMessageHandler:self name:@"keydown"];
    NSRect frame = NSMakeRect(200, 200, 640, 400);
    Class windowClass = mode == ModeCancelOnly ? [CancelGuardWindow class] : [GuardedWindow class];
    self.window = [[windowClass alloc] initWithContentRect:frame
                                                 styleMask:NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskResizable
                                                   backing:NSBackingStoreBuffered
                                                     defer:NO];
    self.window.releasedWhenClosed = NO;
    self.window.collectionBehavior = NSWindowCollectionBehaviorFullScreenPrimary;
    self.window.disableEscapeExitsFullscreen = mode != ModeAllow;
    self.window.delegate = self;
    WKWebView *webView = [[WKWebView alloc] initWithFrame:frame configuration:config];
    self.window.webView = webView;
    self.window.contentView = webView;
    [webView loadHTMLString:@"<textarea autofocus>draft</textarea><script>"
        "addEventListener('load', () => document.querySelector('textarea').focus());"
        "addEventListener('keydown', e => { if (e.target === document.querySelector('textarea')) window.webkit.messageHandlers.keydown.postMessage(e.key); });"
        "</script>" baseURL:nil];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    [self.window makeKeyAndOrderFront:nil];
    [NSApp activateIgnoringOtherApps:YES];
    after(1.0, ^{
        [self.window makeFirstResponder:webView];
        [self.window toggleFullScreen:nil];
    });
    after(20.0, ^{
        fail("timed out");
    });
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
    after(1.0, ^{
        [self pressEscape];
    });
}

- (void)pressEscape {
    NSInteger number = self.window.windowNumber;
    NSTimeInterval now = [[NSProcessInfo processInfo] systemUptime];
    for (NSEventType type = NSEventTypeKeyDown; type <= NSEventTypeKeyUp; type++) {
        NSEvent *event = [NSEvent keyEventWithType:type
                                          location:NSZeroPoint
                                     modifierFlags:0
                                         timestamp:now
                                      windowNumber:number
                                           context:nil
                                        characters:@"\x1b"
                       charactersIgnoringModifiers:@"\x1b"
                                         isARepeat:NO
                                           keyCode:53];
        [NSApp postEvent:event atStart:NO];
    }
    after(2.0, ^{
        if (!self.pageSawEscape) {
            fail("the focused textarea did not receive the Escape keydown");
        }
        switch (mode) {
        case ModeCancelOnly:
            if (isFullscreen(self.window)) {
                fail("Escape kept fullscreen with only the cancelOperation: override");
            }
            pass("with only the cancelOperation: override, Escape left fullscreen");
        case ModeAllow:
            if (isFullscreen(self.window)) {
                fail("Escape did not leave fullscreen with the option off");
            }
            pass("with the option off, Escape left fullscreen");
        case ModeGuarded:
            if (!isFullscreen(self.window)) {
                fail("Escape left fullscreen");
            }
            self.leaving = YES;
            [self.window toggleFullScreen:nil];
        }
    });
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    if (self.leaving) {
        pass("with both overrides, Escape kept fullscreen and toggleFullScreen: left it");
    }
}
@end

int main(int argc, const char **argv) {
    @autoreleasepool {
        BOOL fullscreen = NO;
        for (int i = 1; i < argc; i++) {
            if (strcmp(argv[i], "--filter") == 0) {
                testFilter();
                return 0;
            }
            if (strcmp(argv[i], "--fullscreen") == 0) {
                fullscreen = YES;
            } else if (strcmp(argv[i], "--guarded") == 0) {
                mode = ModeGuarded;
            } else if (strcmp(argv[i], "--cancel-only") == 0) {
                mode = ModeCancelOnly;
            } else if (strcmp(argv[i], "--allow") == 0) {
                mode = ModeAllow;
            }
        }
        if (!fullscreen) {
            fprintf(stderr, "usage: %s --filter | --fullscreen [--guarded | --cancel-only | --allow]\n", argv[0]);
            return 2;
        }
        NSApplication *app = [NSApplication sharedApplication];
        Driver *driver = [Driver new];
        app.delegate = driver;
        [app run];
    }
    return 0;
}
