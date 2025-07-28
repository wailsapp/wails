//go:build darwin

package dock

// #cgo CFLAGS: -x objective-c
// #import <AppKit/AppKit.h>
//
// void hideDockIcon() {
//     [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
// }
//
// void showDockIcon() {
//     [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
// }
import "C"
import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type darwinDock struct{}

// Creates a new Dock Service.
func New() *DockService {
	return &DockService{
		impl: &darwinDock{},
	}
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
