//go:build darwin && !ios

#ifndef appdelegate_h
#define appdelegate_h

#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSResponder <NSApplicationDelegate>
@property bool shouldTerminateWhenLastWindowClosed;
@property bool shuttingDown;
- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app;
@end

extern void HandleOpenFile(char *);

// Dock menu: returns the NSMenu* to show for a right-click on the Dock icon,
// or NULL for none (see menu_mac_extras_darwin.go)
extern void* HandleDockMenu(void);

// State restoration: MacOptions.SupportsSecureRestorableState (see
// webview_window_restoration_darwin.go)
extern bool HandleSupportsSecureRestorableState(void);

// Declarations for Apple Event based custom URL handling and universal link
extern void HandleOpenURL(char*);

// User activities (Handoff, Spotlight continuation, universal links). The
// Go exports live in activity_manager_darwin.go; wailsUserActivityJSON is
// implemented in activity_manager_darwin.m and returns a strdup'd string
// the caller frees.
extern bool activityWillContinue(char *activityType);
extern bool activityContinue(char *activityJSON);
extern void activityDidFail(char *activityType, char *message);
extern void activityDidUpdate(char *activityJSON);
extern char *wailsUserActivityJSON(NSUserActivity *activity);

@interface CustomProtocolSchemeHandler : NSObject
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void StartCustomProtocolHandler(void);

#endif /* appdelegate_h */
