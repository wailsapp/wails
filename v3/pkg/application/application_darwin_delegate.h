//go:build darwin

#ifndef appdelegate_h
#define appdelegate_h

#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate>
@property bool shouldTerminateWhenLastWindowClosed;
- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app;
@end

#endif
