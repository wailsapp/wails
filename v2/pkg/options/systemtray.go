package options

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// SystemTray contains options for the system tray
type SystemTray struct {
	LightModeIcon      *SystemTrayIcon
	DarkModeIcon       *SystemTrayIcon
	Title              string
	Tooltip            string
	StartHidden        bool
	Menu               *menu.Menu
	OnLeftClick        func()
	OnRightClick       func()
	OnLeftDoubleClick  func()
	OnRightDoubleClick func()
	OnMenuClose        func()
	OnMenuOpen         func()
}

// SystemTrayIcon represents a system tray icon
type SystemTrayIcon struct {
	Data []byte
}
