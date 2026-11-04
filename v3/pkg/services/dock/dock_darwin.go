//go:build darwin && !ios && !purego

package dock

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void hideDockIcon() {
    dispatch_sync(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    });
}

void showDockIcon() {
    dispatch_sync(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    });
}

bool setBadge(const char *label) {
    __block bool success = false;
    dispatch_sync(dispatch_get_main_queue(), ^{
        // Ensure the app is in Regular activation policy (dock icon visible)
        NSApplicationActivationPolicy currentPolicy = [NSApp activationPolicy];
        if (currentPolicy != NSApplicationActivationPolicyRegular) {
            success = false;
            return;
        }

        NSString *nsLabel = nil;
		if (label != NULL) {
			nsLabel = [NSString stringWithUTF8String:label];
		}
		[[NSApp dockTile] setBadgeLabel:nsLabel];
		[[NSApp dockTile] display];
		success = true;
    });
    return success;
}
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

type darwinDock struct {
	mu    sync.RWMutex
	Badge *string
}

// HideAppIcon hides the app icon in the macOS Dock.
func (d *darwinDock) HideAppIcon() {
	C.hideDockIcon()
}

// ShowAppIcon shows the app icon in the macOS Dock.
// Note: After showing the dock icon, you may need to call SetBadge again
// to reapply any previously set badge, as changing activation policies clears the badge.
func (d *darwinDock) ShowAppIcon() {
	C.showDockIcon()
}

// setBadge handles the C call and updates the internal badge state with locking.
func (d *darwinDock) setBadge(label *string) error {
	var cLabel *C.char
	if label != nil {
		cLabel = C.CString(*label)
		defer C.free(unsafe.Pointer(cLabel))
	}

	success := C.setBadge(cLabel)
	if !success {
		return fmt.Errorf("failed to set badge")
	}

	d.mu.Lock()
	d.Badge = label
	d.mu.Unlock()

	return nil
}
