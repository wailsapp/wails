//go:build darwin

package single_instance

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework AppKit -mmacosx-version-min=10.13

#import <AppKit/AppKit.h>

void activateApplicationWithProcessID(int pid) {
	NSRunningApplication *app = [NSRunningApplication runningApplicationWithProcessIdentifier:pid];
	if (app != nil) {
		[app unhide];
		[app activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
	}
}
*/
import "C"

func (p *Plugin) activeInstance(pid int) error {
	C.activateApplicationWithProcessID(C.int(pid))
	return nil
}
