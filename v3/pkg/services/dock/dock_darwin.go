//go:build darwin

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
	"context"
	"fmt"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type darwinDock struct {
	mu    sync.RWMutex
	Badge *string
}

// Creates a new Dock Service.
func New() *DockService {
	return &DockService{
		impl: &darwinDock{
			Badge: nil,
		},
	}
}

// NewWithOptions creates a new dock service with badge options.
// Currently, options are not available on macOS and are ignored.
func NewWithOptions(options BadgeOptions) *DockService {
	return New()
}

func (d *darwinDock) Startup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

func (d *darwinDock) Shutdown() error {
	return nil
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

// SetBadge sets the badge label on the application icon.
// Available default badge labels:
// Single space " " empty badge
// Empty string "" dot "●" indeterminate badge
func (d *darwinDock) SetBadge(label string) error {
	// Always pick a label (use "●" if empty), then allocate + free exactly once.
	if label == "" {
		label = "●" // Default badge character
	}
	return d.setBadge(&label)
}

// SetCustomBadge is not supported on macOS, SetBadge is called instead.
func (d *darwinDock) SetCustomBadge(label string, options BadgeOptions) error {
	return d.SetBadge(label)
}

// RemoveBadge removes the badge label from the application icon.
func (d *darwinDock) RemoveBadge() error {
	return d.setBadge(nil)
}

// GetBadge returns the badge label on the application icon.
func (d *darwinDock) GetBadge() *string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Badge
}
