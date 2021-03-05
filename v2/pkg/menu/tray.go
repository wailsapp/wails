package menu

// TrayMenu are the options
type TrayMenu struct {

	// Label is the text we wish to display in the tray
	Label string

	// Icon is the name of the tray icon we wish to display.
	// These are read up during build from <projectdir>/trayicons and
	// the filenames are used as IDs, minus the extension
	// EG: <projectdir>/trayicons/main.png can be referenced here with "main"
	Icon string

	// Menu is the initial menu we wish to use for the tray
	Menu *Menu

	// OnOpen is called when the Menu is opened
	OnOpen func()

	// OnClose is called when the Menu is closed
	OnClose func()
}
