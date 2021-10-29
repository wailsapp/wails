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
    [self.mainWindow makeKeyAndOrderFront:self];
    if (self.alwaysOnTop) {
        [self.mainWindow setLevel:NSStatusWindowLevel];
    }
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    [NSApp activateIgnoringOtherApps:YES];
}

@synthesize touchBar;

@end
