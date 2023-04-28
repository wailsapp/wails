package application

type BackdropType int32

const (
	Auto    BackdropType = 0
	None    BackdropType = 1
	Mica    BackdropType = 2
	Acrylic BackdropType = 3
	Tabbed  BackdropType = 4
)

type WindowsWindow struct {
	// Select the type of translucent backdrop. Requires Windows 11 22621 or later.
	BackdropType BackdropType
	// Disable the icon in the titlebar
	DisableIcon bool
}
