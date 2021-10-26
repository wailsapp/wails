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
//
//- (void) CreateMenu {
//    [NSApplication sharedApplication];
//    menubar = [[NSMenu new] autorelease];
//    id appMenuItem = [[NSMenuItem new] autorelease];
//    [menubar addItem:appMenuItem];
//    [NSApp setMainMenu:menubar];
//    id appMenu = [[NSMenu new] autorelease];
//    id appName = [[NSProcessInfo processInfo] processName];
//    id quitTitle = [@"Quit " stringByAppendingString:appName];
//    id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
//        action:@selector(terminate:) keyEquivalent:@"q"]
//              autorelease];
//    [appMenu addItem:quitMenuItem];
//    [appMenuItem setSubmenu:appMenu];
//}
//
//- (void) dealloc {
//    [super dealloc];
//    window = nil;
//    menubar = nil;
//}

@synthesize touchBar;

@end
