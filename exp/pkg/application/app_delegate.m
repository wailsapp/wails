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
    systemEventHandler(EventApplicationDidFinishLaunching);

}

- (void)setApplicationActivationPolicy:(NSApplicationActivationPolicy)policy
{
    self.activationPolicy = policy;
}

- (void)applicationWillTerminate:(NSNotification *)aNotification
{
    systemEventHandler(EventApplicationWillTerminate);
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification
{
    systemEventHandler(EventApplicationDidBecomeActive);
}

- (void)applicationWillHide:(NSNotification *)aNotification
{
    systemEventHandler(EventApplicationWillHide);
}

- (void)applicationDidHide:(NSNotification *)aNotification
{
    systemEventHandler(EventApplicationDidHide);
}

- (void)applicationWillUnhide:(NSNotification *)aNotification
{
    systemEventHandler(EventApplicationWillUnhide);
}

- (void)applicationDidUnhide:(NSNotification *)aNotification
{
    systemEventHandler(EventApplicationDidUnhide);
}






@end