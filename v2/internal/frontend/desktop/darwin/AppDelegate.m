//
//  AppDelegate.m
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

#import "AppDelegate.h"

@implementation AppDelegate
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return NO;  
}
- (void)applicationWillFinishLaunching:(NSNotification *)aNotification {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    if (self.alwaysOnTop) {
        [self.mainWindow setLevel:NSStatusWindowLevel];
    }
    if ( !self.startHidden ) {
        [self.mainWindow makeKeyAndOrderFront:self];
    }
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    [NSApp activateIgnoringOtherApps:YES];
    if ( self.startFullscreen ) {
        NSWindowCollectionBehavior behaviour = [self.mainWindow collectionBehavior];
        behaviour |= NSWindowCollectionBehaviorFullScreenPrimary;
        [self.mainWindow setCollectionBehavior:behaviour];
        [self.mainWindow toggleFullScreen:nil];
    }
}

- (void)dealloc {
    [super dealloc];
}

@synthesize touchBar;

@end
