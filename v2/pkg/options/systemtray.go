package options

// SystemTray contains options for the system tray
type SystemTray struct {
	LightModeIcon *SystemTrayIcon
	DarkModeIcon  *SystemTrayIcon
	Title         string
	Tooltip       string
	StartHidden   bool
}

// SystemTrayIcon represents a system tray icon
type SystemTrayIcon struct {
	Data []byte
}
