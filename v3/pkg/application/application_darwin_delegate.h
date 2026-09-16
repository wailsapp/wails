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

// Declarations for Apple Event based custom URL handling and universal link
extern void HandleOpenURL(char*);

@interface CustomProtocolSchemeHandler : NSObject
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void StartCustomProtocolHandler(void);

#endif /* appdelegate_h */
