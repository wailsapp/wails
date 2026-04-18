//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c -fblocks
#cgo LDFLAGS: -framework Foundation -framework Cocoa

#include <stdlib.h>
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

static char *g_wailsCapturedLaunchURL = NULL;

// stopAppEventLoop stops [NSApp run] by posting a synthetic event.
static void stopAppEventLoop(void) {
    [NSApp stop:nil];
    NSEvent *e = [NSEvent otherEventWithType:NSEventTypeApplicationDefined
                                   location:NSZeroPoint
                              modifierFlags:0
                                  timestamp:0
                               windowNumber:0
                                    context:nil
                                    subtype:0
                                      data1:0
                                      data2:0];
    [NSApp postEvent:e atStart:NO];
}

// _WailsURLCaptureHandler handles the kAEGetURL Apple Event and stops the run
// loop immediately so the second instance can proceed to notify the first.
@interface _WailsURLCaptureHandler : NSObject
@end

@implementation _WailsURLCaptureHandler
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
    NSString *urlStr = [[event paramDescriptorForKeyword:keyDirectObject] stringValue];
    if (urlStr && g_wailsCapturedLaunchURL == NULL) {
        g_wailsCapturedLaunchURL = strdup([urlStr UTF8String]);
    }
    stopAppEventLoop();
}
@end

// _WailsURLCaptureDelegate provides a safety-net timeout so the second
// instance never blocks indefinitely if no URL event arrives.
@interface _WailsURLCaptureDelegate : NSObject <NSApplicationDelegate>
@property (nonatomic, assign) double timeoutSeconds;
@end

@implementation _WailsURLCaptureDelegate
- (void)applicationDidFinishLaunching:(NSNotification *)note {
    double t = _timeoutSeconds;
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(t * NSEC_PER_SEC)),
                   dispatch_get_main_queue(), ^{
        stopAppEventLoop();
    });
}
@end

// CaptureLaunchURL briefly runs an NSApplication event loop so LaunchServices
// can deliver any pending kAEGetURL Apple Event (which it only does after the
// app has "finished launching").  The run loop is stopped as soon as the URL
// event is received, or after timeoutSeconds if no event arrives.
// Returns NULL on timeout.  Caller must free the returned string.
static char *CaptureLaunchURL(double timeoutSeconds) {
    g_wailsCapturedLaunchURL = NULL;

    NSApplication *app = [NSApplication sharedApplication];
    // Run without a dock icon — this is a short-lived second instance.
    [app setActivationPolicy:NSApplicationActivationPolicyProhibited];

    // Register the URL handler BEFORE [NSApp run] calls finishLaunching so we
    // catch the event the moment it is dispatched.
    NSAppleEventManager *mgr = [NSAppleEventManager sharedAppleEventManager];
    [mgr setEventHandler:[_WailsURLCaptureHandler class]
         andSelector:@selector(handleGetURLEvent:withReplyEvent:)
         forEventClass:kInternetEventClass
         andEventID:kAEGetURL];

    _WailsURLCaptureDelegate *delegate = [[_WailsURLCaptureDelegate alloc] init];
    delegate.timeoutSeconds = timeoutSeconds;
    [app setDelegate:delegate];

    // [NSApp run] calls finishLaunching (signalling to LaunchServices that this
    // process is ready) and then enters the event loop.  The run loop is stopped
    // by either the URL handler (early, on success) or the delegate timeout.
    [NSApp run];

    [app setDelegate:nil];
    [mgr removeEventHandlerForEventClass:kInternetEventClass andEventID:kAEGetURL];

    return g_wailsCapturedLaunchURL;
}
*/
import "C"
import "unsafe"

// captureLaunchURL briefly runs an NSApplication event loop so that
// LaunchServices can deliver any pending kAEGetURL Apple Event (e.g. when
// this process was force-launched via "open -n URL").
// Returns the URL string, or "" if none arrived within the timeout.
func captureLaunchURL() string {
	cURL := C.CaptureLaunchURL(C.double(0.3))
	if cURL == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cURL))
	return C.GoString(cURL)
}
