//go:build darwin && !ios

package dock

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GetBadge returns the badge label on the application icon.
func (d *darwinDock) GetBadge() *string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Badge
}

// RemoveBadge removes the badge label from the application icon.
func (d *darwinDock) RemoveBadge() error {
	return d.setBadge(nil)
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

func (d *darwinDock) Shutdown() error {
	return nil
}

func (d *darwinDock) Startup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

// New creates a new Dock Service.
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
