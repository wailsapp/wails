//go:build ios
#import "application_ios_delegate.h"
#import "../events/events_ios.h"
#import "application_ios.h"
extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);
extern bool hasListeners(unsigned int);
extern void iosApplicationDidLaunch(void);
// WailsIOSMain (app's generated main_ios.go) runs the user's main()/app.Run().
// The delegate starts it AFTER UIKit has launched (see below).
extern void WailsIOSMain(void);
// Registers the UNUserNotificationCenter delegate so local notifications are
// shown while the app is in the foreground (and taps are handled). Apple
// requires this be set before launch finishes, hence the call below.
extern void ios_notifications_init(void);

@implementation WailsAppDelegate
- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    // Set global appDelegate reference and bring up a window if needed
    appDelegate = self;
    if (self.window == nil) {
        // Start the window with the launch-screen colour (a "LaunchBackground"
        // colour asset, also referenced by UILaunchScreen) so there's no white
        // flash between the launch screen and the first WebView paint. The Go
        // options set the colour too, but that happens after this delegate runs,
        // so it can't colour the initial window. Falls back to white if the
        // asset isn't present.
        UIColor *launchBG = [UIColor colorNamed:@"LaunchBackground"] ?: [UIColor whiteColor];
        self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
        self.window.backgroundColor = launchBG;
        UIViewController *rootVC = [[UIViewController alloc] init];
        rootVC.view.backgroundColor = launchBG;
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
    // Register the notification-center delegate before launch finishes so local
    // notifications appear while the app is foregrounded (otherwise iOS delivers
    // them silently and no banner is shown).
    ios_notifications_init();
    // Unconditional launch signal for the Go runtime. platformRun waits on
    // this and emits ApplicationDidFinishLaunching from the Go side once the
    // event listeners are wired up - emitting it from here would race the Go
    // runtime's startup and the event could be dropped.
    // Start the Go runtime NOW — only after UIKit has delivered the launch and
    // the window exists. Starting Go earlier (concurrently with UIApplicationMain)
    // intermittently corrupts the FrontBoard launch handshake on a physical
    // device, so this method never fires (blank cold launch / 0x8BADF00D). Run it
    // on a background thread so app.Run()'s blocking loop never touches the main
    // thread. WailsIOSMain -> user main() -> app.Run(); the window's run() then
    // creates the WebView (appDelegate/window are already set above).
    dispatch_async(dispatch_get_global_queue(QOS_CLASS_USER_INITIATED, 0), ^{
        WailsIOSMain();
    });
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
