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
// openFile implemented as well just in case, but it's not called (at least we don't know how to call it)
-(BOOL)application:(NSApplication *)sender openFile:(NSString *)filename
{
   const char* utf8FileName = filename.UTF8String;
   HandleOpenFile((char*)utf8FileName);
   return YES;
}

// for some reasons it's triggered even when only one file is opened, instead of openFile.
-(void)application:(NSApplication *)sender openFiles:(NSArray<NSString *> *)filenames
{
   int count = [filenames count];
	 int i;
	 char **fileNamesArray = NULL;
	 fileNamesArray = (char**)realloc(fileNamesArray, i*sizeof(*fileNamesArray));
   for (i=0; i<count; i++) {
      NSString* fnm = [filenames objectAtIndex: i];
      const char* utf8FileName = fnm.UTF8String;
      fileNamesArray[i] = (char*)utf8FileName;
   }
   HandleOpenFiles(fileNamesArray, count);
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
}

- (void)dealloc {
    [super dealloc];
}

@synthesize touchBar;

@end
