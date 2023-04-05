//go:build darwin

// Package mac provides a set of functions to interact with the macOS platform.
package mac

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework ServiceManagement

#import <Foundation/Foundation.h>
#import <ServiceManagement/ServiceManagement.h>

// Get the bundle ID
char* getBundleID() {
	NSString *bundleID = [[NSBundle mainBundle] bundleIdentifier];
	return (char*)[bundleID UTF8String];
}
*/
import "C"

// GetBundleID returns the bundle ID of the application.
func GetBundleID() string {
	return C.GoString(C.getBundleID())
}
