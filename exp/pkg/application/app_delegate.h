//go:build darwin

#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate>
@property NSApplicationActivationPolicy activationPolicy;
- (void)setApplicationActivationPolicy:(NSApplicationActivationPolicy)policy;
@end