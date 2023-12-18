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
        [self.mainWindow setLevel:NSFloatingWindowLevel];
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
    // we pass message in object because otherwise sandboxing will prevent us from sending it https://developer.apple.com/forums/thread/129437
    NSString * myString = [NSString stringWithUTF8String:message];
    [[NSDistributedNotificationCenter defaultCenter]
        postNotificationName:[NSString stringWithUTF8String:singleInstanceUniqueId]
        object:(__bridge const void *)(myString)
        userInfo:nil
        deliverImmediately:YES];
}

char* GetMacOsNativeTempDir() {
    NSString *tempDir = NSTemporaryDirectory();
    char *copy = strdup([tempDir UTF8String]);

    return copy;
}

- (void)handleSecondInstanceNotification:(NSNotification *)note;
{
    if (note.object != nil) {
        NSString * message = (__bridge NSString *)note.object;
        const char* utf8Message = message.UTF8String;
        HandleSecondInstanceData((char*)utf8Message);
    }
}

- (void)dealloc {
    [super dealloc];
}

@synthesize touchBar;

@end
