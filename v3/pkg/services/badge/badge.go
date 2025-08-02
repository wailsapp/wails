package badge

import (
	"context"
	"image/color"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type platformBadge interface {
	// Lifecycle methods
	Startup(ctx context.Context, options application.ServiceOptions) error
	Shutdown() error

	SetBadge(label string) error
	SetCustomBadge(label string, options Options) error
	RemoveBadge() error
}

// Service represents the badge service
type BadgeService struct {
	impl platformBadge
}

type Options struct {
	TextColour       color.RGBA
	BackgroundColour color.RGBA
	FontName         string
	FontSize         int
	SmallFontSize    int
}

// ServiceName returns the name of the service.
func (b *BadgeService) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/badge"
}

// ServiceStartup is called when the service is loaded.
func (b *BadgeService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return b.impl.Startup(ctx, options)
}

// ServiceShutdown is called when the service is unloaded.
func (b *BadgeService) ServiceShutdown() error {
	return b.impl.Shutdown()
}

// SetBadge sets the badge label on the application icon.
func (b *BadgeService) SetBadge(label string) error {
	return b.impl.SetBadge(label)
}

func (b *BadgeService) SetCustomBadge(label string, options Options) error {
	return b.impl.SetCustomBadge(label, options)
}

// RemoveBadge removes the badge label from the application icon.
func (b *BadgeService) RemoveBadge() error {
	return b.impl.RemoveBadge()
}
