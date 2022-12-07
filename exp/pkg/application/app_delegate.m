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



@end