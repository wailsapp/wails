package options

// App contains options for creating the App
type App struct {
	Title         string
	Width         int
	Height        int
	DisableResize bool
	Fullscreen    bool
	Frameless     bool
	MinWidth      int
	MinHeight     int
	MaxWidth      int
	MaxHeight     int
	StartHidden   bool
	DevTools      bool
	Colour        int
	Mac           Mac
}

// MergeDefaults will set the minimum default values for an application
func (a *App) MergeDefaults() {

	// Create a default title
	if len(a.Title) == 0 {
		a.Title = "My Wails App"
	}

	// Default width
	if a.Width == 0 {
		a.Width = 1024
	}

	// Default height
	if a.Height == 0 {
		a.Height = 768
	}
}
