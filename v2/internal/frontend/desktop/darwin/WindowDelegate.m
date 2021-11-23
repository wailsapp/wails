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
    [sender orderOut:nil];
    if( self.hideOnClose == false ) {
        processMessage("Q");
    }
    return !self.hideOnClose;
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
