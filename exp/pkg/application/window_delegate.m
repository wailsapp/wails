//go:build darwin
//
//  WindowDelegate.m
//  test
//
//  Created by Lea Anthony on 10/10/21.
//
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "window_delegate.h"
#import "../events/events.h"
extern void processMessage(unsigned int, const char*);
@implementation WindowDelegate
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

- (void) mouseDown:(NSEvent*)someEvent {
    NSLog(@"MOUSE DOWN!!!");
}

// GENERATED EVENTS START
- (void)windowDidBecomeKey:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidBecomeKey);
}

- (void)windowDidBecomeMain:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidBecomeMain);
}

- (void)windowDidBeginSheet:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidBeginSheet);
}

- (void)windowDidChangeAlpha:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeAlpha);
}

- (void)windowDidChangeBackingLocation:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeBackingLocation);
}

- (void)windowDidChangeBackingProperties:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeBackingProperties);
}

- (void)windowDidChangeCollectionBehavior:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeCollectionBehavior);
}

- (void)windowDidChangeEffectiveAppearance:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeEffectiveAppearance);
}

- (void)windowDidChangeOcclusionState:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeOcclusionState);
}

- (void)windowDidChangeOrderingMode:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeOrderingMode);
}

- (void)windowDidChangeScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeScreen);
}

- (void)windowDidChangeScreenParameters:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeScreenParameters);
}

- (void)windowDidChangeScreenProfile:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeScreenProfile);
}

- (void)windowDidChangeScreenSpace:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeScreenSpace);
}

- (void)windowDidChangeScreenSpaceProperties:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeScreenSpaceProperties);
}

- (void)windowDidChangeSharingType:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeSharingType);
}

- (void)windowDidChangeSpace:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeSpace);
}

- (void)windowDidChangeSpaceOrderingMode:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeSpaceOrderingMode);
}

- (void)windowDidChangeTitle:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeTitle);
}

- (void)windowDidChangeToolbar:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeToolbar);
}

- (void)windowDidChangeVisibility:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidChangeVisibility);
}

- (void)windowDidClose:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidClose);
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidDeminiaturize);
}

- (void)windowDidEndSheet:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidEndSheet);
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidEnterFullScreen);
}

- (void)windowDidEnterVersionBrowser:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidEnterVersionBrowser);
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidExitFullScreen);
}

- (void)windowDidExitVersionBrowser:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidExitVersionBrowser);
}

- (void)windowDidExpose:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidExpose);
}

- (void)windowDidFocus:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidFocus);
}

- (void)windowDidMiniaturize:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidMiniaturize);
}

- (void)windowDidMove:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidMove);
}

- (void)windowDidOrderOffScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidOrderOffScreen);
}

- (void)windowDidOrderOnScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidOrderOnScreen);
}

- (void)windowDidResignKey:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidResignKey);
}

- (void)windowDidResignMain:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidResignMain);
}

- (void)windowDidResize:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidResize);
}

- (void)windowDidUnfocus:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUnfocus);
}

- (void)windowDidUpdate:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdate);
}

- (void)windowDidUpdateAlpha:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdateAlpha);
}

- (void)windowDidUpdateCollectionBehavior:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdateCollectionBehavior);
}

- (void)windowDidUpdateCollectionProperties:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdateCollectionProperties);
}

- (void)windowDidUpdateShadow:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdateShadow);
}

- (void)windowDidUpdateTitle:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdateTitle);
}

- (void)windowDidUpdateToolbar:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdateToolbar);
}

- (void)windowDidUpdateVisibility:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowDidUpdateVisibility);
}

- (void)windowWillBecomeKey:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillBecomeKey);
}

- (void)windowWillBecomeMain:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillBecomeMain);
}

- (void)windowWillBeginSheet:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillBeginSheet);
}

- (void)windowWillChangeOrderingMode:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillChangeOrderingMode);
}

- (void)windowWillClose:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillClose);
}

- (void)windowWillDeminiaturize:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillDeminiaturize);
}

- (void)windowWillEnterFullScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillEnterFullScreen);
}

- (void)windowWillEnterVersionBrowser:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillEnterVersionBrowser);
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillExitFullScreen);
}

- (void)windowWillExitVersionBrowser:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillExitVersionBrowser);
}

- (void)windowWillFocus:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillFocus);
}

- (void)windowWillMiniaturize:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillMiniaturize);
}

- (void)windowWillMove:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillMove);
}

- (void)windowWillOrderOffScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillOrderOffScreen);
}

- (void)windowWillOrderOnScreen:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillOrderOnScreen);
}

- (void)windowWillResignMain:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillResignMain);
}

- (void)windowWillResize:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillResize);
}

- (void)windowWillUnfocus:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUnfocus);
}

- (void)windowWillUpdate:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdate);
}

- (void)windowWillUpdateAlpha:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdateAlpha);
}

- (void)windowWillUpdateCollectionBehavior:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdateCollectionBehavior);
}

- (void)windowWillUpdateCollectionProperties:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdateCollectionProperties);
}

- (void)windowWillUpdateShadow:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdateShadow);
}

- (void)windowWillUpdateTitle:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdateTitle);
}

- (void)windowWillUpdateToolbar:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdateToolbar);
}

- (void)windowWillUpdateVisibility:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUpdateVisibility);
}

- (void)windowWillUseStandardFrame:(NSNotification *)notification {
    processWindowEvent(self.windowId, EventWindowWillUseStandardFrame);
}

- (void)webView:(WKWebView *)webview didStartProvisionalNavigation:(WKNavigation *)navigation {
	processWindowEvent(self.windowId, EventWebViewDidStartProvisionalNavigation);
}

- (void)webView:(WKWebView *)webview didReceiveServerRedirectForProvisionalNavigation:(WKNavigation *)navigation {
	processWindowEvent(self.windowId, EventWebViewDidReceiveServerRedirectForProvisionalNavigation);
}

- (void)webView:(WKWebView *)webview didFinishNavigation:(WKNavigation *)navigation {
	processWindowEvent(self.windowId, EventWebViewDidFinishNavigation);
}

- (void)webView:(WKWebView *)webview didCommitNavigation:(WKNavigation *)navigation {
	processWindowEvent(self.windowId, EventWebViewDidCommitNavigation);
}

// GENERATED EVENTS END
@end
