//go:build darwin

package dock

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void hideDockIcon() {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    });
}

void showDockIcon() {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    });
}

static void setBadge(const char *label) {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSString *nsLabel = nil;
		if (label != NULL) {
			nsLabel = [NSString stringWithUTF8String:label];
		}
		[[NSApp dockTile] setBadgeLabel:nsLabel];
		[[NSApp dockTile] display];
    });
}
*/
import "C"
import (
	"context"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type darwinDock struct{}

// Creates a new Dock Service.
func New() *DockService {
	return &DockService{
		impl: &darwinDock{},
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
func (d *darwinDock) ShowAppIcon() {
	C.showDockIcon()
}

// SetBadge sets the badge label on the application icon.
func (d *darwinDock) SetBadge(label string) error {
    // Always pick a label (use “●” if empty), then allocate + free exactly once.
    value := label
    if value == "" {
        value = "●" // Default badge character
    }
    cLabel := C.CString(value)
    defer C.free(unsafe.Pointer(cLabel))

    C.setBadge(cLabel)
    return nil
}

// SetCustomBadge is not supported on macOS, SetBadge is called instead.
func (d *darwinDock) SetCustomBadge(label string, options BadgeOptions) error {
	return d.SetBadge(label)
}

// RemoveBadge removes the badge label from the application icon.
func (d *darwinDock) RemoveBadge() error {
	C.setBadge(nil)
	return nil
}
