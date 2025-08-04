//go:build darwin

#ifndef appdelegate_h
#define appdelegate_h

#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSResponder <NSApplicationDelegate>
@property bool shouldTerminateWhenLastWindowClosed;
@property bool shuttingDown;
- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app;
@end

extern void HandleOpenFile(char *);

// Declarations for Apple Event based custom URL handling
extern void HandleCustomProtocol(char*);

@interface CustomProtocolSchemeHandler : NSObject
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void StartCustomProtocolHandler(void);

#endif /* appdelegate_h */
