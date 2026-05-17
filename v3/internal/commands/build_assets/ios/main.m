//go:build ios
// Minimal bootstrap: delegate comes from Go archive (WailsAppDelegate)
#import <UIKit/UIKit.h>
#include <stdio.h>

// External Go initialization function from the c-archive (declare before use)
extern void WailsIOSMain();

int main(int argc, char * argv[]) {
    @autoreleasepool {
        // Disable buffering so stdout/stderr from Go log.Printf flush immediately
        setvbuf(stdout, NULL, _IONBF, 0);
        setvbuf(stderr, NULL, _IONBF, 0);

        // Start Go runtime on a background queue to avoid blocking main thread/UI
        dispatch_async(dispatch_get_global_queue(QOS_CLASS_USER_INITIATED, 0), ^{
            WailsIOSMain();
        });

        // Run UIApplicationMain using WailsAppDelegate provided by the Go archive
        return UIApplicationMain(argc, argv, nil, @"WailsAppDelegate");
    }
}