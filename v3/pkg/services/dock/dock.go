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
	return "github.com/wailsapp/wails/v3/services/dock"
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
	application.InvokeSync(func() {
		d.impl.HideAppIcon()
	})
}

// ShowAppIcon shows the app icon in the dock/taskbar.
func (d *DockService) ShowAppIcon() {
	application.InvokeSync(func() {
		d.impl.ShowAppIcon()
	})
}

// SetBadge sets the badge label on the application icon.
// This method ensures the badge call is made on the main thread to avoid crashes.
func (d *DockService) SetBadge(label string) error {
	return application.InvokeSyncWithError(func() error {
		return d.impl.SetBadge(label)
	})
}

// SetCustomBadge sets the badge label on the application icon with custom options.
// This method ensures the badge call is made on the main thread to avoid crashes.
func (d *DockService) SetCustomBadge(label string, options BadgeOptions) error {
	return application.InvokeSyncWithError(func() error {
		return d.impl.SetCustomBadge(label, options)
	})
}

// RemoveBadge removes the badge label from the application icon.
// This method ensures the badge call is made on the main thread to avoid crashes.
func (d *DockService) RemoveBadge() error {
	return application.InvokeSyncWithError(func() error {
		return d.impl.RemoveBadge()
	})
}
