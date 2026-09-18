//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>
#import "webview_panel_darwin.h"

extern bool processWindowKeyEquivalent(unsigned int, const char*);
extern void processWindowKeyDownEvent(unsigned int, const char*);

@implementation WebviewPanel

- (instancetype)initWithContentRect:(NSRect)contentRect
                           styleMask:(NSWindowStyleMask)styleMask
                             backing:(NSBackingStoreType)bufferingType
                               defer:(BOOL)deferCreation {
    self = [super initWithContentRect:contentRect
                            styleMask:styleMask
                              backing:bufferingType
                                defer:deferCreation];
    if (self != nil) {
        self.alphaValue = 1.0;
        self.backgroundColor = [NSColor clearColor];
        self.opaque = NO;
        self.movableByWindowBackground = YES;

        // NSPanel defaults differ from NSWindow. Wails windows are expected to
        // remain visible when another app activates and to release on Close.
        self.hidesOnDeactivate = NO;
        self.releasedWhenClosed = YES;
    }
    return self;
}

// A non-activating NSPanel can leave WKWebView as first responder, bypassing
// NSWindow's keyDown:. Check Wails bindings before normal AppKit delivery and
// consume only a binding that was actually handled. Unbound web input follows
// the normal responder chain exactly once.
- (void)sendEvent:(NSEvent*)event {
    if (event.type == NSEventTypeKeyDown) {
        WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)self.delegate;
        NSString* accelerator = acceleratorStringFromKeyEvent(event);
        if (delegate != nil && accelerator.length > 0 &&
            processWindowKeyEquivalent(delegate.windowId, accelerator.UTF8String)) {
            return;
        }
    }
    [super sendEvent:event];
}

- (void)keyDown:(NSEvent*)event {
    WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)self.delegate;
    NSString* accelerator = acceleratorStringFromKeyEvent(event);
    if (delegate != nil && accelerator.length > 0) {
        processWindowKeyDownEvent(delegate.windowId, accelerator.UTF8String);
    }
}

- (BOOL)performKeyEquivalent:(NSEvent*)event {
    if (dispatchKeyEquivalent(event, self)) {
        return YES;
    }
    return [super performKeyEquivalent:event];
}

- (BOOL)canBecomeKeyWindow {
    return YES;
}

- (BOOL)canBecomeMainWindow {
    return NO;
}

- (BOOL)acceptsFirstResponder {
    return YES;
}

- (BOOL)becomeFirstResponder {
    return YES;
}

- (BOOL)resignFirstResponder {
    return YES;
}

- (void)cancelOperation:(id)sender {
    if (self.disableEscapeExitsFullscreen &&
        (self.styleMask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen) {
        return;
    }
    [super cancelOperation:sender];
}

- (void)setDelegate:(id<NSWindowDelegate>)delegate {
    id<NSWindowDelegate> previousDelegate = [super delegate];
    if (previousDelegate == delegate) {
        return;
    }
    [delegate retain];
    [super setDelegate:delegate];
    [previousDelegate release];
    if ([delegate isKindOfClass:[WebviewWindowDelegate class]]) {
        [self registerForDraggedTypes:@[NSFilenamesPboardType]];
    }
}

- (void)dealloc {
    [self.webView.configuration.userContentController removeScriptMessageHandlerForName:@"external"];
    id<NSWindowDelegate> retainedDelegate = [super delegate];
    [super setDelegate:nil];
    [retainedDelegate release];
    [super dealloc];
}

@end
