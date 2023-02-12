//go:build darwin
#import "app_delegate.h"
#import "../events/events.h"
extern bool hasListeners(unsigned int);
@implementation AppDelegate
- (void)dealloc
{
    [super dealloc];
}
// Create the applicationShouldTerminateAfterLastWindowClosed: method
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)theApplication
{
    return self.shouldTerminateWhenLastWindowClosed;
}
// GENERATED EVENTS START
- (void)applicationDidBecomeActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidBecomeActive) ) {
        processApplicationEvent(EventApplicationDidBecomeActive);
    }
}

- (void)applicationDidChangeBackingProperties:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeBackingProperties) ) {
        processApplicationEvent(EventApplicationDidChangeBackingProperties);
    }
}

- (void)applicationDidChangeEffectiveAppearance:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeEffectiveAppearance) ) {
        processApplicationEvent(EventApplicationDidChangeEffectiveAppearance);
    }
}

- (void)applicationDidChangeIcon:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeIcon) ) {
        processApplicationEvent(EventApplicationDidChangeIcon);
    }
}

- (void)applicationDidChangeOcclusionState:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeOcclusionState) ) {
        processApplicationEvent(EventApplicationDidChangeOcclusionState);
    }
}

- (void)applicationDidChangeScreenParameters:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeScreenParameters) ) {
        processApplicationEvent(EventApplicationDidChangeScreenParameters);
    }
}

- (void)applicationDidChangeStatusBarFrame:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeStatusBarFrame) ) {
        processApplicationEvent(EventApplicationDidChangeStatusBarFrame);
    }
}

- (void)applicationDidChangeStatusBarOrientation:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidChangeStatusBarOrientation) ) {
        processApplicationEvent(EventApplicationDidChangeStatusBarOrientation);
    }
}

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidFinishLaunching) ) {
        processApplicationEvent(EventApplicationDidFinishLaunching);
    }
}

- (void)applicationDidHide:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidHide) ) {
        processApplicationEvent(EventApplicationDidHide);
    }
}

- (void)applicationDidResignActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidResignActive) ) {
        processApplicationEvent(EventApplicationDidResignActive);
    }
}

- (void)applicationDidUnhide:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidUnhide) ) {
        processApplicationEvent(EventApplicationDidUnhide);
    }
}

- (void)applicationDidUpdate:(NSNotification *)notification {
    if( hasListeners(EventApplicationDidUpdate) ) {
        processApplicationEvent(EventApplicationDidUpdate);
    }
}

- (void)applicationWillBecomeActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillBecomeActive) ) {
        processApplicationEvent(EventApplicationWillBecomeActive);
    }
}

- (void)applicationWillFinishLaunching:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillFinishLaunching) ) {
        processApplicationEvent(EventApplicationWillFinishLaunching);
    }
}

- (void)applicationWillHide:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillHide) ) {
        processApplicationEvent(EventApplicationWillHide);
    }
}

- (void)applicationWillResignActive:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillResignActive) ) {
        processApplicationEvent(EventApplicationWillResignActive);
    }
}

- (void)applicationWillTerminate:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillTerminate) ) {
        processApplicationEvent(EventApplicationWillTerminate);
    }
}

- (void)applicationWillUnhide:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillUnhide) ) {
        processApplicationEvent(EventApplicationWillUnhide);
    }
}

- (void)applicationWillUpdate:(NSNotification *)notification {
    if( hasListeners(EventApplicationWillUpdate) ) {
        processApplicationEvent(EventApplicationWillUpdate);
    }
}

// GENERATED EVENTS END
@end
