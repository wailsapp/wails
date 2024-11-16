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

#endif
