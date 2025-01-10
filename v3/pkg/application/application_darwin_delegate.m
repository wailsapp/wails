//go:build darwin
#import "application_darwin_delegate.h"
#import "../events/events_darwin.h"
extern bool hasListeners(unsigned int);
extern bool shouldQuitApplication();
extern void cleanup();
extern void handleSecondInstanceData(char * message);
@implementation AppDelegate
- (void)dealloc
{
    [super dealloc];
}
-(BOOL)application:(NSApplication *)sender openFile:(NSString *)filename
 {
    const char* utf8FileName = filename.UTF8String;
    HandleOpenFile((char*)utf8FileName);
    return YES;
 }
// Create the applicationShouldTerminateAfterLastWindowClosed: method
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)theApplication
{
    return self.shouldTerminateWhenLastWindowClosed;
}
- (void)themeChanged:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeTheme) ) {
        processApplicationEvent(EventApplicationDidChangeTheme, NULL);
    }
}
- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender {
    if( ! shouldQuitApplication() ) {
        return NSTerminateCancel;
    }
    if( !self.shuttingDown ) {
        self.shuttingDown = true;
        cleanup();
    }
    return NSTerminateNow;
}
- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app
{
    return YES;
}
- (BOOL)applicationShouldHandleReopen:(NSNotification *)notification
                    hasVisibleWindows:(BOOL)flag {
    if( hasListeners(EventApplicationShouldHandleReopen) ) {
        processApplicationEvent(EventApplicationShouldHandleReopen, @{@"hasVisibleWindows": @(flag)});
    }
    
    return TRUE;
}
- (void)handleSecondInstanceNotification:(NSNotification *)note;
{
   if (note.userInfo[@"message"] != nil) {
        NSString *message = note.userInfo[@"message"];
        const char* utf8Message = message.UTF8String;
        handleSecondInstanceData((char*)utf8Message);
    }
}

// GENERATED EVENTS START
- (void)applicationDidBecomeActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidBecomeActive) ) {
        processApplicationEvent(EventApplicationDidBecomeActive, NULL);
    }
}

- (void)applicationDidChangeBackingProperties:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeBackingProperties) ) {
        processApplicationEvent(EventApplicationDidChangeBackingProperties, NULL);
    }
}

- (void)applicationDidChangeEffectiveAppearance:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeEffectiveAppearance) ) {
        processApplicationEvent(EventApplicationDidChangeEffectiveAppearance, NULL);
    }
}

- (void)applicationDidChangeIcon:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeIcon) ) {
        processApplicationEvent(EventApplicationDidChangeIcon, NULL);
    }
}

- (void)applicationDidChangeOcclusionState:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeOcclusionState) ) {
        processApplicationEvent(EventApplicationDidChangeOcclusionState, NULL);
    }
}

- (void)applicationDidChangeScreenParameters:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeScreenParameters) ) {
        processApplicationEvent(EventApplicationDidChangeScreenParameters, NULL);
    }
}

- (void)applicationDidChangeStatusBarFrame:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeStatusBarFrame) ) {
        processApplicationEvent(EventApplicationDidChangeStatusBarFrame, NULL);
    }
}

- (void)applicationDidChangeStatusBarOrientation:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeStatusBarOrientation) ) {
        processApplicationEvent(EventApplicationDidChangeStatusBarOrientation, NULL);
    }
}

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidFinishLaunching) ) {
        processApplicationEvent(EventApplicationDidFinishLaunching, NULL);
    }
}

- (void)applicationDidHide:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidHide) ) {
        processApplicationEvent(EventApplicationDidHide, NULL);
    }
}

- (void)applicationDidResignActiveNotification:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidResignActiveNotification) ) {
        processApplicationEvent(EventApplicationDidResignActiveNotification, NULL);
    }
}

- (void)applicationDidUnhide:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidUnhide) ) {
        processApplicationEvent(EventApplicationDidUnhide, NULL);
    }
}

- (void)applicationDidUpdate:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidUpdate) ) {
        processApplicationEvent(EventApplicationDidUpdate, NULL);
    }
}

- (void)applicationWillBecomeActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillBecomeActive) ) {
        processApplicationEvent(EventApplicationWillBecomeActive, NULL);
    }
}

- (void)applicationWillFinishLaunching:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillFinishLaunching) ) {
        processApplicationEvent(EventApplicationWillFinishLaunching, NULL);
    }
}

- (void)applicationWillHide:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillHide) ) {
        processApplicationEvent(EventApplicationWillHide, NULL);
    }
}

- (void)applicationWillResignActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillResignActive) ) {
        processApplicationEvent(EventApplicationWillResignActive, NULL);
    }
}

- (void)applicationWillTerminate:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillTerminate) ) {
        processApplicationEvent(EventApplicationWillTerminate, NULL);
    }
}

- (void)applicationWillUnhide:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillUnhide) ) {
        processApplicationEvent(EventApplicationWillUnhide, NULL);
    }
}

- (void)applicationWillUpdate:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillUpdate) ) {
        processApplicationEvent(EventApplicationWillUpdate, NULL);
    }
}

// GENERATED EVENTS END
@end
