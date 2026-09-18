package dock

import (
	"context"
	"image/color"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type platformDock interface {
	// Lifecycle methods
	Startup(ctx context.Context, options application.ServiceOptions) error
	Shutdown() error

	// Dock icon visibility methods
	HideAppIcon()
	ShowAppIcon()

	// Badge methods
	SetBadge(label string) error
	SetCustomBadge(label string, options BadgeOptions) error
	RemoveBadge() error
	GetBadge() *string

	// Progress methods
	SetProgress(fraction float64) error
	ClearProgress() error
	GetProgress() *float64
}

// Service represents the dock service
type DockService struct {
	impl platformDock
}

// BadgeOptions represents options for customizing badge appearance
type BadgeOptions struct {
	TextColour       color.RGBA
	BackgroundColour color.RGBA
	FontName         string
	FontSize         int
	SmallFontSize    int
}

// ServiceName returns the name of the service.
func (d *DockService) ServiceName() string {
	return "github.com/wailsapp/wails/v3/pkg/services/dock"
}

// ServiceStartup is called when the service is loaded.
func (d *DockService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return d.impl.Startup(ctx, options)
}

// ServiceShutdown is called when the service is unloaded.
func (d *DockService) ServiceShutdown() error {
	return d.impl.Shutdown()
}

// HideAppIcon hides the app icon in the dock/taskbar.
func (d *DockService) HideAppIcon() {
	d.impl.HideAppIcon()
}

// ShowAppIcon shows the app icon in the dock/taskbar.
func (d *DockService) ShowAppIcon() {
	d.impl.ShowAppIcon()
}

// SetBadge sets the badge label on the application icon.
func (d *DockService) SetBadge(label string) error {
	return d.impl.SetBadge(label)
}

// SetCustomBadge sets the badge label on the application icon with custom options.
func (d *DockService) SetCustomBadge(label string, options BadgeOptions) error {
	return d.impl.SetCustomBadge(label, options)
}

// RemoveBadge removes the badge label from the application icon.
func (d *DockService) RemoveBadge() error {
	return d.impl.RemoveBadge()
}

// GetBadge returns the badge label on the application icon.
func (d *DockService) GetBadge() *string {
	return d.impl.GetBadge()
}

// SetProgress draws a progress bar over the application icon (macOS Dock
// tile). fraction is clamped to 0..1. The badge, if any, keeps showing.
func (d *DockService) SetProgress(fraction float64) error {
	if fraction < 0 {
		fraction = 0
	}
	if fraction > 1 {
		fraction = 1
	}
	return d.impl.SetProgress(fraction)
}

// ClearProgress removes the progress bar from the application icon.
func (d *DockService) ClearProgress() error {
	return d.impl.ClearProgress()
}

// GetProgress returns the current progress fraction, or nil when no
// progress bar is shown.
func (d *DockService) GetProgress() *float64 {
	return d.impl.GetProgress()
}
