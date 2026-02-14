package runtime

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

// TraySetSystemTray sets the system tray menu
func TraySetSystemTray(ctx context.Context, trayMenu *menu.TrayMenu) {
	frontend := getFrontend(ctx)
	frontend.TraySetSystemTray(trayMenu)
}
