//go:build darwin && !ios && !server
#import "application_darwin_delegate.h"
#import "../events/events_darwin.h"
#import <CoreServices/CoreServices.h> // For Apple Event constants
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
- (void)workspaceWillSleep:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillSleep) ) {
        processApplicationEvent(EventApplicationWillSleep, NULL);
    }
}
- (void)workspaceDidWake:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidWake) ) {
        processApplicationEvent(EventApplicationDidWake, NULL);
    }
}
- (void)workspaceScreensDidSleep:(NSNotification *)notification {
    if( hasListeners(EventApplicationScreensDidSleep) ) {
        processApplicationEvent(EventApplicationScreensDidSleep, NULL);
    }
}
- (void)workspaceScreensDidWake:(NSNotification *)notification {
    if( hasListeners(EventApplicationScreensDidWake) ) {
        processApplicationEvent(EventApplicationScreensDidWake, NULL);
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
// State restoration: answers MacOptions.SupportsSecureRestorableState. The
// method is always implemented so macOS 14 never logs the secure coding
// warning (see webview_window_restoration_darwin.go).
- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app
{
    return HandleSupportsSecureRestorableState() ? YES : NO;
}
- (BOOL)applicationShouldHandleReopen:(NSNotification *)notification
                    hasVisibleWindows:(BOOL)flag { // Changed from NSApplication to NSNotification
    if( hasListeners(EventApplicationShouldHandleReopen) ) {
        processApplicationEvent(EventApplicationShouldHandleReopen, @{@"hasVisibleWindows": @(flag)});
    }
    return TRUE;
}
- (void)handleSecondInstanceNotification:(NSNotification *)note;
{
   if (note.object != nil) {
        NSString *message = (NSString *)note.object;
        const char* utf8Message = message.UTF8String;
        handleSecondInstanceData((char*)utf8Message);
    }
}
// User activities (Handoff, Spotlight continuation and universal links).
// The Go side lives in activity_manager_darwin.go; universal links
// (NSUserActivityTypeBrowsingWeb) are also routed to HandleOpenURL there so
// they share the custom URL scheme path.
- (BOOL)application:(NSApplication *)application willContinueUserActivityWithType:(NSString *)userActivityType {
    return activityWillContinue((char *)[userActivityType UTF8String]);
}
- (BOOL)application:(NSApplication *)application continueUserActivity:(NSUserActivity *)userActivity restorationHandler:(void (^)(NSArray<id<NSUserActivityRestoring>> * _Nullable))restorationHandler {
    char *json = wailsUserActivityJSON(userActivity);
    BOOL handled = activityContinue(json);
    free(json);
    return handled;
}
- (void)application:(NSApplication *)application didFailToContinueUserActivityWithType:(NSString *)userActivityType error:(NSError *)error {
    NSString *message = error.localizedDescription ?: @"";
    activityDidFail((char *)[userActivityType UTF8String], (char *)[message UTF8String]);
}
- (void)application:(NSApplication *)application didUpdateUserActivity:(NSUserActivity *)userActivity {
    char *json = wailsUserActivityJSON(userActivity);
    activityDidUpdate(json);
    free(json);
}
- (NSMenu *)applicationDockMenu:(NSApplication *)sender {
    return (NSMenu *)HandleDockMenu();
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

- (void)applicationDidResignActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidResignActive) ) {
        processApplicationEvent(EventApplicationDidResignActive, NULL);
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
// Implementation for Apple Event based custom URL handling
@implementation CustomProtocolSchemeHandler
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
   NSString *urlStr = [[event paramDescriptorForKeyword:keyDirectObject] stringValue];
   if (urlStr) {
       HandleOpenURL((char*)[urlStr UTF8String]);
   }
}
@end
void StartCustomProtocolHandler(void) {
    NSAppleEventManager *appleEventManager = [NSAppleEventManager sharedAppleEventManager];
    [appleEventManager setEventHandler:[CustomProtocolSchemeHandler class]
        andSelector:@selector(handleGetURLEvent:withReplyEvent:)
        forEventClass:kInternetEventClass
        andEventID: kAEGetURL];
}
