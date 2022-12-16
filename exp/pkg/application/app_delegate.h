//go:build darwin

#ifndef appdelegate_h
#define appdelegate_h

#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate>
@property NSApplicationActivationPolicy activationPolicy;
- (void)setApplicationActivationPolicy:(NSApplicationActivationPolicy)policy;
@end

#endif
