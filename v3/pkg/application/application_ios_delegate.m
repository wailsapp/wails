//go:build ios
#import "application_ios_delegate.h"
#import "../events/events_ios.h"
#import "application_ios.h"
extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);
extern bool hasListeners(unsigned int);
@implementation WailsAppDelegate
- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    // Set global appDelegate reference and bring up a window if needed
    appDelegate = self;
    if (self.window == nil) {
        self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
        self.window.backgroundColor = [UIColor whiteColor];
        UIViewController *rootVC = [[UIViewController alloc] init];
        rootVC.view.backgroundColor = [UIColor whiteColor];
        self.window.rootViewController = rootVC;
        [self.window makeKeyAndVisible];
    }
    // Apply app-wide background colour if configured
    unsigned char r = 255, g = 255, b = 255, a = 255;
    if (ios_get_app_background_color(&r, &g, &b, &a)) {
        CGFloat fr = ((CGFloat)r) / 255.0;
        CGFloat fg = ((CGFloat)g) / 255.0;
        CGFloat fb = ((CGFloat)b) / 255.0;
        CGFloat fa = ((CGFloat)a) / 255.0;
        UIColor *color = [UIColor colorWithRed:fr green:fg blue:fb alpha:fa];
        self.window.backgroundColor = color;
        self.window.rootViewController.view.backgroundColor = color;
    }
    if (!self.viewControllers) {
        self.viewControllers = [NSMutableArray array];
    }
    if (hasListeners(EventApplicationDidFinishLaunching)) {
        processApplicationEvent(EventApplicationDidFinishLaunching, NULL);
    }
    return YES;
}
// GENERATED EVENTS START
- (void)applicationDidBecomeActive:(UIApplication *)application {
    if( hasListeners(EventApplicationDidBecomeActive) ) {
        processApplicationEvent(EventApplicationDidBecomeActive, NULL);
    }
}

- (void)applicationDidEnterBackground:(UIApplication *)application {
    if( hasListeners(EventApplicationDidEnterBackground) ) {
        processApplicationEvent(EventApplicationDidEnterBackground, NULL);
    }
}

- (void)applicationDidFinishLaunching:(UIApplication *)application {
    if( hasListeners(EventApplicationDidFinishLaunching) ) {
        processApplicationEvent(EventApplicationDidFinishLaunching, NULL);
    }
}

- (void)applicationDidReceiveMemoryWarning:(UIApplication *)application {
    if( hasListeners(EventApplicationDidReceiveMemoryWarning) ) {
        processApplicationEvent(EventApplicationDidReceiveMemoryWarning, NULL);
    }
}

- (void)applicationWillEnterForeground:(UIApplication *)application {
    if( hasListeners(EventApplicationWillEnterForeground) ) {
        processApplicationEvent(EventApplicationWillEnterForeground, NULL);
    }
}

- (void)applicationWillResignActive:(UIApplication *)application {
    if( hasListeners(EventApplicationWillResignActive) ) {
        processApplicationEvent(EventApplicationWillResignActive, NULL);
    }
}

- (void)applicationWillTerminate:(UIApplication *)application {
    if( hasListeners(EventApplicationWillTerminate) ) {
        processApplicationEvent(EventApplicationWillTerminate, NULL);
    }
}

// GENERATED EVENTS END
@end
