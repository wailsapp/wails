//go:build linux

package badge

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxBadge struct{}

// New creates a new Badge Service.
// On Linux, this returns a no-op implementation since most desktop environments
// don't have standardized dock badge functionality.
func New() *Service {
	return &Service{
		impl: &linuxBadge{},
	}
}

// NewWithOptions creates a new badge service with the given options.
// On Linux, this returns a no-op implementation since most desktop environments
// don't have standardized dock badge functionality. Options are ignored.
func NewWithOptions(options Options) *Service {
	return New()
}

func (l *linuxBadge) Startup(ctx context.Context, options application.ServiceOptions) error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}

func (l *linuxBadge) Shutdown() error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}

// SetBadge is a no-op on Linux since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *linuxBadge) SetBadge(label string) error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}

// SetCustomBadge is a no-op on Linux since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *linuxBadge) SetCustomBadge(label string, options Options) error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}

// RemoveBadge is a no-op on Linux since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *linuxBadge) RemoveBadge() error {
	// No-op: Linux doesn't have standardized badge support
	return nil
}