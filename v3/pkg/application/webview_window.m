//go:build darwin
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "webview_window.h"
#import "../events/events.h"
extern void processMessage(unsigned int, const char*);
extern void processURLRequest(unsigned int, void *);
extern bool hasListeners(unsigned int);
@implementation WebviewWindow
- (WebviewWindow*) initWithContentRect:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;
{
    self = [super initWithContentRect:contentRect styleMask:windowStyle backing:bufferingType defer:deferCreation];
    [self setAlphaValue:1.0];
    [self setBackgroundColor:[NSColor clearColor]];
    [self setOpaque:NO];
    [self setMovableByWindowBackground:YES];
    return self;
}
- (void)keyDown:(NSEvent *)event {
}
- (BOOL)canBecomeKeyWindow {
    return YES;
}
- (BOOL) canBecomeMainWindow {
    return YES;
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
@end
@implementation WebviewWindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    if( self.hideOnClose ) {
        [NSApp hide:nil];
        return false;
    }
    return true;
}
// Handle script messages from the external bridge
- (void)userContentController:(nonnull WKUserContentController *)userContentController didReceiveScriptMessage:(nonnull WKScriptMessage *)message {
    NSString *m = message.body;
    /*
    // TODO: Check for drag
    if ( [m isEqualToString:@"drag"] ) {
        if( [self IsFullScreen] ) {
            return;
        }
        if( self.mouseEvent != nil ) {
           [self.mainWindow performWindowDragWithEvent:self.mouseEvent];
        }
        return;
    }
    */
    const char *_m = [m UTF8String];
    processMessage(self.windowId, _m);
}
- (void)handleLeftMouseDown:(NSEvent *)event {
    self.leftMouseEvent = event;
    NSWindow *window = [event window];
    WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[window delegate];
    if( self.invisibleTitleBarHeight > 0 ) {
        NSPoint location = [event locationInWindow];
        NSRect frame = [window frame];
        if( location.y > frame.size.height - self.invisibleTitleBarHeight ) {
            [window performWindowDragWithEvent:event];
            return;
        }
    }
}
- (void)handleLeftMouseUp:(NSWindow *)window {
    self.leftMouseEvent = nil;
}
- (void)webView:(nonnull WKWebView *)webView startURLSchemeTask:(nonnull id<WKURLSchemeTask>)urlSchemeTask {
    processURLRequest(self.windowId, urlSchemeTask);
}
- (void)webView:(nonnull WKWebView *)webView stopURLSchemeTask:(nonnull id<WKURLSchemeTask>)urlSchemeTask {
    NSInputStream *stream = urlSchemeTask.request.HTTPBodyStream;
    if (stream) {
        NSStreamStatus status = stream.streamStatus;
        if (status != NSStreamStatusClosed && status != NSStreamStatusNotOpen) {
            [stream close];
        }
    }
}
// GENERATED EVENTS START
- (void)windowDidBecomeKey:(NSNotification *)notification {
    if( hasListeners(EventWindowDidBecomeKey) ) {
        processWindowEvent(self.windowId, EventWindowDidBecomeKey);
    }
}

- (void)windowDidBecomeMain:(NSNotification *)notification {
    if( hasListeners(EventWindowDidBecomeMain) ) {
        processWindowEvent(self.windowId, EventWindowDidBecomeMain);
    }
}

- (void)windowDidBeginSheet:(NSNotification *)notification {
    if( hasListeners(EventWindowDidBeginSheet) ) {
        processWindowEvent(self.windowId, EventWindowDidBeginSheet);
    }
}

- (void)windowDidChangeAlpha:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeAlpha) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeAlpha);
    }
}

- (void)windowDidChangeBackingLocation:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeBackingLocation) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeBackingLocation);
    }
}

- (void)windowDidChangeBackingProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeBackingProperties) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeBackingProperties);
    }
}

- (void)windowDidChangeCollectionBehavior:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeCollectionBehavior) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeCollectionBehavior);
    }
}

- (void)windowDidChangeEffectiveAppearance:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeEffectiveAppearance) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeEffectiveAppearance);
    }
}

- (void)windowDidChangeOcclusionState:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeOcclusionState) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeOcclusionState);
    }
}

- (void)windowDidChangeOrderingMode:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeOrderingMode) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeOrderingMode);
    }
}

- (void)windowDidChangeScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreen);
    }
}

- (void)windowDidChangeScreenParameters:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenParameters) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenParameters);
    }
}

- (void)windowDidChangeScreenProfile:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenProfile) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenProfile);
    }
}

- (void)windowDidChangeScreenSpace:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenSpace) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenSpace);
    }
}

- (void)windowDidChangeScreenSpaceProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenSpaceProperties) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenSpaceProperties);
    }
}

- (void)windowDidChangeSharingType:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeSharingType) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeSharingType);
    }
}

- (void)windowDidChangeSpace:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeSpace) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeSpace);
    }
}

- (void)windowDidChangeSpaceOrderingMode:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeSpaceOrderingMode) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeSpaceOrderingMode);
    }
}

- (void)windowDidChangeTitle:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeTitle) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeTitle);
    }
}

- (void)windowDidChangeToolbar:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeToolbar) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeToolbar);
    }
}

- (void)windowDidChangeVisibility:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeVisibility) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeVisibility);
    }
}

- (void)windowDidClose:(NSNotification *)notification {
    if( hasListeners(EventWindowDidClose) ) {
        processWindowEvent(self.windowId, EventWindowDidClose);
    }
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowDidDeminiaturize) ) {
        processWindowEvent(self.windowId, EventWindowDidDeminiaturize);
    }
}

- (void)windowDidEndSheet:(NSNotification *)notification {
    if( hasListeners(EventWindowDidEndSheet) ) {
        processWindowEvent(self.windowId, EventWindowDidEndSheet);
    }
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidEnterFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidEnterFullScreen);
    }
}

- (void)windowDidEnterVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowDidEnterVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowDidEnterVersionBrowser);
    }
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidExitFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidExitFullScreen);
    }
}

- (void)windowDidExitVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowDidExitVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowDidExitVersionBrowser);
    }
}

- (void)windowDidExpose:(NSNotification *)notification {
    if( hasListeners(EventWindowDidExpose) ) {
        processWindowEvent(self.windowId, EventWindowDidExpose);
    }
}

- (void)windowDidFocus:(NSNotification *)notification {
    if( hasListeners(EventWindowDidFocus) ) {
        processWindowEvent(self.windowId, EventWindowDidFocus);
    }
}

- (void)windowDidMiniaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowDidMiniaturize) ) {
        processWindowEvent(self.windowId, EventWindowDidMiniaturize);
    }
}

- (void)windowDidMove:(NSNotification *)notification {
    if( hasListeners(EventWindowDidMove) ) {
        processWindowEvent(self.windowId, EventWindowDidMove);
    }
}

- (void)windowDidOrderOffScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidOrderOffScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidOrderOffScreen);
    }
}

- (void)windowDidOrderOnScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidOrderOnScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidOrderOnScreen);
    }
}

- (void)windowDidResignKey:(NSNotification *)notification {
    if( hasListeners(EventWindowDidResignKey) ) {
        processWindowEvent(self.windowId, EventWindowDidResignKey);
    }
}

- (void)windowDidResignMain:(NSNotification *)notification {
    if( hasListeners(EventWindowDidResignMain) ) {
        processWindowEvent(self.windowId, EventWindowDidResignMain);
    }
}

- (void)windowDidResize:(NSNotification *)notification {
    if( hasListeners(EventWindowDidResize) ) {
        processWindowEvent(self.windowId, EventWindowDidResize);
    }
}

- (void)windowDidUnfocus:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUnfocus) ) {
        processWindowEvent(self.windowId, EventWindowDidUnfocus);
    }
}

- (void)windowDidUpdate:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdate) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdate);
    }
}

- (void)windowDidUpdateAlpha:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateAlpha) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateAlpha);
    }
}

- (void)windowDidUpdateCollectionBehavior:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateCollectionBehavior) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateCollectionBehavior);
    }
}

- (void)windowDidUpdateCollectionProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateCollectionProperties) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateCollectionProperties);
    }
}

- (void)windowDidUpdateShadow:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateShadow) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateShadow);
    }
}

- (void)windowDidUpdateTitle:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateTitle) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateTitle);
    }
}

- (void)windowDidUpdateToolbar:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateToolbar) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateToolbar);
    }
}

- (void)windowDidUpdateVisibility:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateVisibility) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateVisibility);
    }
}

- (void)windowWillBecomeKey:(NSNotification *)notification {
    if( hasListeners(EventWindowWillBecomeKey) ) {
        processWindowEvent(self.windowId, EventWindowWillBecomeKey);
    }
}

- (void)windowWillBecomeMain:(NSNotification *)notification {
    if( hasListeners(EventWindowWillBecomeMain) ) {
        processWindowEvent(self.windowId, EventWindowWillBecomeMain);
    }
}

- (void)windowWillBeginSheet:(NSNotification *)notification {
    if( hasListeners(EventWindowWillBeginSheet) ) {
        processWindowEvent(self.windowId, EventWindowWillBeginSheet);
    }
}

- (void)windowWillChangeOrderingMode:(NSNotification *)notification {
    if( hasListeners(EventWindowWillChangeOrderingMode) ) {
        processWindowEvent(self.windowId, EventWindowWillChangeOrderingMode);
    }
}

- (void)windowWillClose:(NSNotification *)notification {
    if( hasListeners(EventWindowWillClose) ) {
        processWindowEvent(self.windowId, EventWindowWillClose);
    }
}

- (void)windowWillDeminiaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowWillDeminiaturize) ) {
        processWindowEvent(self.windowId, EventWindowWillDeminiaturize);
    }
}

- (void)windowWillEnterFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillEnterFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillEnterFullScreen);
    }
}

- (void)windowWillEnterVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowWillEnterVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowWillEnterVersionBrowser);
    }
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillExitFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillExitFullScreen);
    }
}

- (void)windowWillExitVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowWillExitVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowWillExitVersionBrowser);
    }
}

- (void)windowWillFocus:(NSNotification *)notification {
    if( hasListeners(EventWindowWillFocus) ) {
        processWindowEvent(self.windowId, EventWindowWillFocus);
    }
}

- (void)windowWillMiniaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowWillMiniaturize) ) {
        processWindowEvent(self.windowId, EventWindowWillMiniaturize);
    }
}

- (void)windowWillMove:(NSNotification *)notification {
    if( hasListeners(EventWindowWillMove) ) {
        processWindowEvent(self.windowId, EventWindowWillMove);
    }
}

- (void)windowWillOrderOffScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillOrderOffScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillOrderOffScreen);
    }
}

- (void)windowWillOrderOnScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillOrderOnScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillOrderOnScreen);
    }
}

- (void)windowWillResignMain:(NSNotification *)notification {
    if( hasListeners(EventWindowWillResignMain) ) {
        processWindowEvent(self.windowId, EventWindowWillResignMain);
    }
}

- (void)windowWillResize:(NSNotification *)notification {
    if( hasListeners(EventWindowWillResize) ) {
        processWindowEvent(self.windowId, EventWindowWillResize);
    }
}

- (void)windowWillUnfocus:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUnfocus) ) {
        processWindowEvent(self.windowId, EventWindowWillUnfocus);
    }
}

- (void)windowWillUpdate:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdate) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdate);
    }
}

- (void)windowWillUpdateAlpha:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateAlpha) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateAlpha);
    }
}

- (void)windowWillUpdateCollectionBehavior:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateCollectionBehavior) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateCollectionBehavior);
    }
}

- (void)windowWillUpdateCollectionProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateCollectionProperties) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateCollectionProperties);
    }
}

- (void)windowWillUpdateShadow:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateShadow) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateShadow);
    }
}

- (void)windowWillUpdateTitle:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateTitle) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateTitle);
    }
}

- (void)windowWillUpdateToolbar:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateToolbar) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateToolbar);
    }
}

- (void)windowWillUpdateVisibility:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateVisibility) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateVisibility);
    }
}

- (void)windowWillUseStandardFrame:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUseStandardFrame) ) {
        processWindowEvent(self.windowId, EventWindowWillUseStandardFrame);
    }
}

- (void)windowFileDraggingEntered:(NSNotification *)notification {
    if( hasListeners(EventWindowFileDraggingEntered) ) {
        processWindowEvent(self.windowId, EventWindowFileDraggingEntered);
    }
}

- (void)windowFileDraggingPerformed:(NSNotification *)notification {
    if( hasListeners(EventWindowFileDraggingPerformed) ) {
        processWindowEvent(self.windowId, EventWindowFileDraggingPerformed);
    }
}

- (void)windowFileDraggingExited:(NSNotification *)notification {
    if( hasListeners(EventWindowFileDraggingExited) ) {
        processWindowEvent(self.windowId, EventWindowFileDraggingExited);
    }
}

- (void)webView:(WKWebView *)webview didStartProvisionalNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidStartProvisionalNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidStartProvisionalNavigation);
    }
}

- (void)webView:(WKWebView *)webview didReceiveServerRedirectForProvisionalNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidReceiveServerRedirectForProvisionalNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidReceiveServerRedirectForProvisionalNavigation);
    }
}

- (void)webView:(WKWebView *)webview didFinishNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidFinishNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidFinishNavigation);
    }
}

- (void)webView:(WKWebView *)webview didCommitNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidCommitNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidCommitNavigation);
    }
}

// GENERATED EVENTS END
@end
