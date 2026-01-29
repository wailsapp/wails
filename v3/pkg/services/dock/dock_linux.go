//go:build linux

package dock

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxDock struct{}

// New creates a new Dock Service.
// On Linux, this returns a stub implementation since dock icon visibility
// and badge functionality are not standardized across desktop environments.
func New() *DockService {
	return &DockService{
		impl: &linuxDock{},
	}
}

// NewWithOptions creates a new dock service with badge options.
// On Linux, this returns a stub implementation since badge functionality
// is not standardized across desktop environments. Options are ignored.
func NewWithOptions(options BadgeOptions) *DockService {
	return New()
}

func (l *linuxDock) Startup(ctx context.Context, options application.ServiceOptions) error {
	// No-op: Linux doesn't have standardized dock/badge support
	return nil
}

func (l *linuxDock) Shutdown() error {
	// No-op: Linux doesn't have standardized dock/badge support
	return nil
}

// HideAppIcon is a stub on Linux since dock icon visibility is not
// standardized across desktop environments.
func (l *linuxDock) HideAppIcon() {
	// No-op: Linux doesn't have standardized dock icon visibility support
}

// ShowAppIcon is a stub on Linux since dock icon visibility is not
// standardized across desktop environments.
func (l *linuxDock) ShowAppIcon() {
	// No-op: Linux doesn't have standardized dock icon visibility support
}

// SetBadge is a stub on Linux since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *linuxDock) SetBadge(label string) error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}

// SetCustomBadge is a stub on Linux since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *linuxDock) SetCustomBadge(label string, options BadgeOptions) error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}

// RemoveBadge is a stub on Linux since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *linuxDock) RemoveBadge() error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}

func (l *linuxDock) GetBadge() *string {
	// No-op: Linux doesn't have standardized badge support
	return nil
}