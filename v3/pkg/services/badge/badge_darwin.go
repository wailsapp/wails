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

func New() *Service {
	return &Service{
		impl: &darwinBadge{},
	}
}

func (d *darwinBadge) Startup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

func (d *darwinBadge) Shutdown() error {
	return nil
}

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

func (d *darwinBadge) RemoveBadge() error {
	C.setBadge(nil)
	return nil
}
