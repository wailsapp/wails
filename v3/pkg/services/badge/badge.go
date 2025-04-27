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
	RemoveBadge() error
}

// Service represents the notifications service
type Service struct {
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
func (b *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/badge"
}

// ServiceStartup is called when the service is loaded.
func (b *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return b.impl.Startup(ctx, options)
}

// ServiceShutdown is called when the service is unloaded.
func (b *Service) ServiceShutdown() error {
	return b.impl.Shutdown()
}

// SetBadge sets the badge label on the application icon.
func (b *Service) SetBadge(label string) error {
	return b.impl.SetBadge(label)
}

// RemoveBadge removes the badge label from the application icon.
func (b *Service) RemoveBadge() error {
	return b.impl.RemoveBadge()
}
