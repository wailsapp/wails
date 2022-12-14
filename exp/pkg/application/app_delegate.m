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
    processApplicationEvent(EventApplicationDidFinishLaunching);

}


- (void)setApplicationActivationPolicy:(NSApplicationActivationPolicy)policy
{
    self.activationPolicy = policy;
}

// GENERATED EVENTS START
- (void)applicationDidBecomeActive:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidBecomeActive);
}

- (void)applicationDidChangeBackingProperties:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidChangeBackingProperties);
}

- (void)applicationDidChangeEffectiveAppearance:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidChangeEffectiveAppearance);
}

- (void)applicationDidChangeIcon:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidChangeIcon);
}

- (void)applicationDidChangeOcclusionState:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidChangeOcclusionState);
}

- (void)applicationDidChangeScreenParameters:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidChangeScreenParameters);
}

- (void)applicationDidChangeStatusBarFrame:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidChangeStatusBarFrame);
}

- (void)applicationDidChangeStatusBarOrientation:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidChangeStatusBarOrientation);
}

- (void)applicationDidHide:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidHide);
}

- (void)applicationDidResignActive:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidResignActive);
}

- (void)applicationDidUnhide:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidUnhide);
}

- (void)applicationDidUpdate:(NSNotification *)notification {
    processApplicationEvent(EventApplicationDidUpdate);
}

- (void)applicationWillBecomeActive:(NSNotification *)notification {
    processApplicationEvent(EventApplicationWillBecomeActive);
}

- (void)applicationWillFinishLaunching:(NSNotification *)notification {
    processApplicationEvent(EventApplicationWillFinishLaunching);
}

- (void)applicationWillHide:(NSNotification *)notification {
    processApplicationEvent(EventApplicationWillHide);
}

- (void)applicationWillResignActive:(NSNotification *)notification {
    processApplicationEvent(EventApplicationWillResignActive);
}

- (void)applicationWillTerminate:(NSNotification *)notification {
    processApplicationEvent(EventApplicationWillTerminate);
}

- (void)applicationWillUnhide:(NSNotification *)notification {
    processApplicationEvent(EventApplicationWillUnhide);
}

- (void)applicationWillUpdate:(NSNotification *)notification {
    processApplicationEvent(EventApplicationWillUpdate);
}

// GENERATED EVENTS END

@end



