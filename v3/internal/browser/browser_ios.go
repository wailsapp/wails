//go:build ios

package browser

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework UIKit

#include <stdlib.h>

// Objective-C function to open URL using UIApplication
static int iosOpenURL(const char* urlStr) {
    @autoreleasepool {
        NSString *urlString = [NSString stringWithUTF8String:urlStr];
        NSURL *url = [NSURL URLWithString:urlString];
        if (url == nil) {
            return 0;
        }

        // Use UIApplication to open URL
        // This must be called on the main thread
        __block BOOL success = NO;
        dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);

        dispatch_async(dispatch_get_main_queue(), ^{
            if (@available(iOS 10.0, *)) {
                [[UIApplication sharedApplication] openURL:url options:@{} completionHandler:^(BOOL opened) {
                    success = opened;
                    dispatch_semaphore_signal(semaphore);
                }];
            } else {
                success = [[UIApplication sharedApplication] openURL:url];
                dispatch_semaphore_signal(semaphore);
            }
        });

        // Wait for completion with timeout
        dispatch_semaphore_wait(semaphore, dispatch_time(DISPATCH_TIME_NOW, 5 * NSEC_PER_SEC));
        return success ? 1 : 0;
    }
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// openURL opens a URL using iOS's UIApplication openURL.
func openURL(url string) error {
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))

	result := C.iosOpenURL(cURL)
	if result == 0 {
		return fmt.Errorf("failed to open URL: %s", url)
	}
	return nil
}

// openFile opens a file URL using iOS's UIApplication.
// On iOS, this typically opens the file with the appropriate app.
func openFile(path string) error {
	// On iOS, we use file:// URLs to open files
	fileURL := "file://" + path
	return openURL(fileURL)
}
