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
-(BOOL)application:(NSApplication *)sender openFile:(NSString *)filename
{
   const char* utf8FileName = filename.UTF8String;
   HandleOpenFile((char*)utf8FileName);
   return YES;
}

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

    if ( self.singleInstanceLockEnabled ) {
      [[NSDistributedNotificationCenter defaultCenter] addObserver:self
          selector:@selector(handleSecondInstanceNotification:) name:self.singleInstanceUniqueId object:nil];
    }
}

void SendDataToFirstInstance(char * singleInstanceUniqueId, char * message) {
    [[NSDistributedNotificationCenter defaultCenter]
        postNotificationName:[NSString stringWithUTF8String:singleInstanceUniqueId]
        object:nil
        userInfo:@{@"message": [NSString stringWithUTF8String:message]}
        deliverImmediately:YES];
}

- (void)handleSecondInstanceNotification:(NSNotification *)note;
{
    if (note.userInfo[@"message"] != nil) {
        NSString *message = note.userInfo[@"message"];
        const char* utf8Message = message.UTF8String;
        HandleSecondInstanceData((char*)utf8Message);
    }
}

- (void)dealloc {
    [super dealloc];
}

@synthesize touchBar;

@end
