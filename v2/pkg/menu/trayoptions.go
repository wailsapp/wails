package menu

type TrayType string

const (
	TrayIcon  TrayType = "icon"
	TrayLabel TrayType = "label"
)

// TrayOptions are the options
type TrayOptions struct {
	// Type is the type of tray item we want
	Type TrayType

	// Label is what is displayed initially when the type is TrayLabel
	Label string

	// Icon is the name of the tray icon we wish to display.
	// These are read up during build from <projectdir>/trayicons and
	// the filenames are used as IDs, minus the extension
	// EG: <projectdir>/trayicons/main.png can be referenced here with "main"
	Icon string

	// Menu is the initial menu we wish to use for the tray
	Menu *Menu
}
