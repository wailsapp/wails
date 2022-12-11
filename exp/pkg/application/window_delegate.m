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

// GENERATED EVENTS START
- (void)windowDidBecomeKey:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidBecomeKey);
}

- (void)windowDidBecomeMain:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidBecomeMain);
}

- (void)windowDidBeginSheet:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidBeginSheet);
}

- (void)windowDidChangeAlpha:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeAlpha);
}

- (void)windowDidChangeBackingLocation:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeBackingLocation);
}

- (void)windowDidChangeBackingProperties:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeBackingProperties);
}

- (void)windowDidChangeCollectionBehavior:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeCollectionBehavior);
}

- (void)windowDidChangeEffectiveAppearance:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeEffectiveAppearance);
}

- (void)windowDidChangeOcclusionState:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeOcclusionState);
}

- (void)windowDidChangeOrderingMode:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeOrderingMode);
}

- (void)windowDidChangeScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeScreen);
}

- (void)windowDidChangeScreenParameters:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeScreenParameters);
}

- (void)windowDidChangeScreenProfile:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeScreenProfile);
}

- (void)windowDidChangeScreenSpace:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeScreenSpace);
}

- (void)windowDidChangeScreenSpaceProperties:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeScreenSpaceProperties);
}

- (void)windowDidChangeSharingType:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeSharingType);
}

- (void)windowDidChangeSpace:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeSpace);
}

- (void)windowDidChangeSpaceOrderingMode:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeSpaceOrderingMode);
}

- (void)windowDidChangeTitle:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeTitle);
}

- (void)windowDidChangeToolbar:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeToolbar);
}

- (void)windowDidChangeVisibility:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidChangeVisibility);
}

- (void)windowDidClose:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidClose);
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidDeminiaturize);
}

- (void)windowDidEndSheet:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidEndSheet);
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidEnterFullScreen);
}

- (void)windowDidEnterVersionBrowser:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidEnterVersionBrowser);
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidExitFullScreen);
}

- (void)windowDidExitVersionBrowser:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidExitVersionBrowser);
}

- (void)windowDidExpose:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidExpose);
}

- (void)windowDidFocus:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidFocus);
}

- (void)windowDidMiniaturize:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidMiniaturize);
}

- (void)windowDidMove:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidMove);
}

- (void)windowDidOrderOffScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidOrderOffScreen);
}

- (void)windowDidOrderOnScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidOrderOnScreen);
}

- (void)windowDidResignKey:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidResignKey);
}

- (void)windowDidResignMain:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidResignMain);
}

- (void)windowDidResize:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidResize);
}

- (void)windowDidUnfocus:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUnfocus);
}

- (void)windowDidUpdate:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdate);
}

- (void)windowDidUpdateAlpha:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdateAlpha);
}

- (void)windowDidUpdateCollectionBehavior:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdateCollectionBehavior);
}

- (void)windowDidUpdateCollectionProperties:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdateCollectionProperties);
}

- (void)windowDidUpdateShadow:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdateShadow);
}

- (void)windowDidUpdateTitle:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdateTitle);
}

- (void)windowDidUpdateToolbar:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdateToolbar);
}

- (void)windowDidUpdateVisibility:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowDidUpdateVisibility);
}

- (void)windowWillBecomeKey:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillBecomeKey);
}

- (void)windowWillBecomeMain:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillBecomeMain);
}

- (void)windowWillBeginSheet:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillBeginSheet);
}

- (void)windowWillChangeOrderingMode:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillChangeOrderingMode);
}

- (void)windowWillClose:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillClose);
}

- (void)windowWillDeminiaturize:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillDeminiaturize);
}

- (void)windowWillEnterFullScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillEnterFullScreen);
}

- (void)windowWillEnterVersionBrowser:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillEnterVersionBrowser);
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillExitFullScreen);
}

- (void)windowWillExitVersionBrowser:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillExitVersionBrowser);
}

- (void)windowWillFocus:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillFocus);
}

- (void)windowWillMiniaturize:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillMiniaturize);
}

- (void)windowWillMove:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillMove);
}

- (void)windowWillOrderOffScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillOrderOffScreen);
}

- (void)windowWillOrderOnScreen:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillOrderOnScreen);
}

- (void)windowWillResignMain:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillResignMain);
}

- (void)windowWillResize:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillResize);
}

- (void)windowWillUnfocus:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUnfocus);
}

- (void)windowWillUpdate:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdate);
}

- (void)windowWillUpdateAlpha:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdateAlpha);
}

- (void)windowWillUpdateCollectionBehavior:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdateCollectionBehavior);
}

- (void)windowWillUpdateCollectionProperties:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdateCollectionProperties);
}

- (void)windowWillUpdateShadow:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdateShadow);
}

- (void)windowWillUpdateTitle:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdateTitle);
}

- (void)windowWillUpdateToolbar:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdateToolbar);
}

- (void)windowWillUpdateVisibility:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUpdateVisibility);
}

- (void)windowWillUseStandardFrame:(NSNotification *)notification {
    windowEventHandler(self.windowId, EventWindowWillUseStandardFrame);
}

// GENERATED EVENTS END

@end











