//go:build darwin
//
//  WindowDelegate.m
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "WindowDelegate.h"
#import "message.h"
#import "WailsContext.h"

@implementation WindowDelegate
- (BOOL)windowShouldClose:(WailsWindow *)sender {
    if( self.hideOnClose == 1 ) {
        [NSApp hide:nil];
        return false;
    }
    if( self.hideOnClose == 2 ) {
        [sender orderOut:nil];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        return false;
    }
    processMessage("Q");
    return false;
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    [self.ctx.mainWindow applyWindowConstraints];
}

- (void)windowWillEnterFullScreen:(NSNotification *)notification {
    [self.ctx.mainWindow disableWindowConstraints];
}

- (NSApplicationPresentationOptions)window:(WailsWindow *)window willUseFullScreenPresentationOptions:(NSApplicationPresentationOptions)proposedOptions {
    return NSApplicationPresentationAutoHideToolbar | NSApplicationPresentationAutoHideMenuBar | NSApplicationPresentationFullScreen;
}


@end
