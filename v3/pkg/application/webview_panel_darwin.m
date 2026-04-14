//go:build darwin && !ios
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "webview_panel_darwin.h"

@implementation WebviewPanel
- (WebviewPanel*) initWithContentRect:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;
{
    self = [super initWithContentRect:contentRect styleMask:windowStyle backing:bufferingType defer:deferCreation];
    [self setAlphaValue:1.0];
    [self setBackgroundColor:[NSColor clearColor]];
    [self setOpaque:NO];
    [self setMovableByWindowBackground:YES];
    return self;
}
// Override sendEvent to intercept key events BEFORE WKWebView consumes them
- (void)sendEvent:(NSEvent *)event {
    if (event.type == NSEventTypeKeyDown) {
        [self keyDown:event];
    }
    [super sendEvent:event];
}
- (void)keyDown:(NSEvent *)event {
    WebviewWindowDelegate *delegate = (WebviewWindowDelegate*)self.delegate;
    dispatchKeyDownEvent(event, delegate.windowId);
}
- (BOOL)canBecomeKeyWindow {
    return YES;
}
- (BOOL) canBecomeMainWindow {
    return NO;  // Panels typically don't become main window
}
- (BOOL) acceptsFirstResponder {
    return YES;
}
- (BOOL) becomeFirstResponder {
    return YES;
}
- (BOOL) resignFirstResponder {
    return YES;
}
- (void) setDelegate:(id<NSWindowDelegate>) delegate {
    [delegate retain];
    [super setDelegate: delegate];
    if ([delegate isKindOfClass:[WebviewWindowDelegate class]]) {
        [self registerForDraggedTypes:@[NSFilenamesPboardType]];
    }
}
- (void) dealloc {
    [self.webView.configuration.userContentController removeScriptMessageHandlerForName:@"external"];
    if (self.delegate) {
        [self.delegate release];
    }
    [super dealloc];
}
@end
