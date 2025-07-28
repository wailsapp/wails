package dock

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type platformDock interface {
	// Lifecycle methods
	Startup(ctx context.Context, options application.ServiceOptions) error
	Shutdown() error

	HideAppIcon()
	ShowAppIcon()
}

// Service represents the notifications service
type DockService struct {
	impl platformDock
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

// HideAppIcon hides the app icon in the macOS Dock.
func (d *DockService) HideAppIcon() {
	d.impl.HideAppIcon()
}

// ShowAppIcon shows the app icon in the macOS Dock.
func (d *DockService) ShowAppIcon() {
	d.impl.ShowAppIcon()
}
