//go:build darwin

#import "app_delegate.h"
#import "../events/events.h"

@implementation AppDelegate

- (void)dealloc
{
    [super dealloc];
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    [NSApp setActivationPolicy:self.activationPolicy];
    [NSApp activateIgnoringOtherApps:YES];

    //callOnApplicationDidFinishLaunchingHandler();
    applicationEventHandler(EventApplicationDidFinishLaunching);

}


- (void)setApplicationActivationPolicy:(NSApplicationActivationPolicy)policy
{
    self.activationPolicy = policy;
}

// GENERATED EVENTS START
- (void)applicationDidBecomeActive:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidBecomeActive);
}

- (void)applicationDidChangeBackingProperties:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidChangeBackingProperties);
}

- (void)applicationDidChangeEffectiveAppearance:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidChangeEffectiveAppearance);
}

- (void)applicationDidChangeIcon:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidChangeIcon);
}

- (void)applicationDidChangeOcclusionState:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidChangeOcclusionState);
}

- (void)applicationDidChangeScreenParameters:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidChangeScreenParameters);
}

- (void)applicationDidChangeStatusBarFrame:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidChangeStatusBarFrame);
}

- (void)applicationDidChangeStatusBarOrientation:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidChangeStatusBarOrientation);
}

- (void)applicationDidHide:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidHide);
}

- (void)applicationDidResignActive:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidResignActive);
}

- (void)applicationDidUnhide:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidUnhide);
}

- (void)applicationDidUpdate:(NSNotification *)notification {
    applicationEventHandler(EventApplicationDidUpdate);
}

- (void)applicationWillBecomeActive:(NSNotification *)notification {
    applicationEventHandler(EventApplicationWillBecomeActive);
}

- (void)applicationWillFinishLaunching:(NSNotification *)notification {
    applicationEventHandler(EventApplicationWillFinishLaunching);
}

- (void)applicationWillHide:(NSNotification *)notification {
    applicationEventHandler(EventApplicationWillHide);
}

- (void)applicationWillResignActive:(NSNotification *)notification {
    applicationEventHandler(EventApplicationWillResignActive);
}

- (void)applicationWillTerminate:(NSNotification *)notification {
    applicationEventHandler(EventApplicationWillTerminate);
}

- (void)applicationWillUnhide:(NSNotification *)notification {
    applicationEventHandler(EventApplicationWillUnhide);
}

- (void)applicationWillUpdate:(NSNotification *)notification {
    applicationEventHandler(EventApplicationWillUpdate);
}

// GENERATED EVENTS END

@end



