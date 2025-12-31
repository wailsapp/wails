//go:build ios

package dock

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type iosDock struct {
}

// New creates a new Dock Service.
// On iOS, this returns a stub implementation.
// iOS badge functionality will be implemented via native bridges.
func New() *DockService {
	return &DockService{
		impl: &iosDock{},
	}
}

// NewWithOptions creates a new dock service with badge options.
// On iOS, this returns a stub implementation. Options are ignored.
func NewWithOptions(options BadgeOptions) *DockService {
	return New()
}

func (d *iosDock) Startup(ctx context.Context, options application.ServiceOptions) error {
	// iOS dock/badge startup - implementation pending native bridge
	return nil
}

func (d *iosDock) Shutdown() error {
	// iOS dock/badge shutdown - implementation pending native bridge
	return nil
}

// HideAppIcon is a stub on iOS.
func (d *iosDock) HideAppIcon() {
	// No-op: iOS doesn't support hiding app icon
}

// ShowAppIcon is a stub on iOS.
func (d *iosDock) ShowAppIcon() {
	// No-op: iOS doesn't support showing/hiding app icon
}

// SetBadge sets the badge on the iOS app icon.
func (d *iosDock) SetBadge(label string) error {
	// iOS badge implementation would go here via native bridge
	return nil
}

// SetCustomBadge is a stub on iOS since iOS badges don't support custom styling.
func (d *iosDock) SetCustomBadge(label string, options BadgeOptions) error {
	// iOS doesn't support custom badge styling, fall back to standard badge
	return d.SetBadge(label)
}

// RemoveBadge removes the badge from the iOS app icon.
func (d *iosDock) RemoveBadge() error {
	// iOS badge removal would go here via native bridge
	return nil
}

// GetBadge retrieves the badge from the iOS app icon.
func (d *iosDock) GetBadge() *string {
	// iOS badge retrieval would go here via native bridge
	return nil
}
