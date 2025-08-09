package menu

// TrayMenu are the options
type TrayMenu struct {
	// Label is the text we wish to display in the tray
	Label string

	// Image is the name of the tray icon we wish to display.
	// These are read up during build from <projectdir>/trayicons and
	// the filenames are used as IDs, minus the extension
	// EG: <projectdir>/trayicons/main.png can be referenced here with "main"
	// If the image is not a filename, it will be treated as base64 image data
	Image string

	// MacTemplateImage indicates that on a Mac, this image is a template image
	MacTemplateImage bool

	// Text Colour
	RGBA string

	// Font
	FontSize int
	FontName string

	// Tooltip
	Tooltip string

	// Callback function when menu clicked
	// Click Callback `json:"-"`

	// Disabled makes the item unselectable
	Disabled bool

	// Menu is the initial menu we wish to use for the tray
	Menu *Menu

	// OnOpen is called when the Menu is opened
	OnOpen func()

	// OnClose is called when the Menu is closed
	OnClose func()
}
