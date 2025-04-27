//go:build darwin

package badge

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

static void setBadge(const char *label) {
    NSString *nsLabel = nil;
	if (label != NULL) {
		nsLabel = [NSString stringWithUTF8String:label];
	}
	[[NSApp dockTile] setBadgeLabel:nsLabel];
	[[NSApp dockTile] display];
}
*/
import "C"
import (
	"context"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type darwinBadge struct{}

// Creates a new Badge Service.
func New() *Service {
	return &Service{
		impl: &darwinBadge{},
	}
}

// NewWithOptions creates a new badge service with the given options.
// Currently, options are not available on macOS and are ignored.
// (Windows-specific)
func NewWithOptions(options Options) *Service {
	return New()
}

func (d *darwinBadge) Startup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

func (d *darwinBadge) Shutdown() error {
	return nil
}

// SetBadge sets the badge label on the application icon.
func (d *darwinBadge) SetBadge(label string) error {
	var cLabel *C.char
	if label != "" {
		cLabel = C.CString(label)
		defer C.free(unsafe.Pointer(cLabel))
	} else {
		cLabel = C.CString("●") // Default badge character
	}
	C.setBadge(cLabel)
	return nil
}

// SetCustomBadge is not supported on macOS, SetBadge is called instead.
// (Windows-specific)
func (d *darwinBadge) SetCustomBadge(label string, options Options) error {
	return d.SetBadge(label)
}

// RemoveBadge removes the badge label from the application icon.
func (d *darwinBadge) RemoveBadge() error {
	C.setBadge(nil)
	return nil
}
