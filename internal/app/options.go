package app

// Options for creating the App
type Options struct {
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
}

// mergeDefaults will set the minimum default values for an application
func (o *Options) mergeDefaults() {

	// Create a default title
	if len(o.Title) == 0 {
		o.Title = "My Wails App"
	}

	// Default width
	if o.Width == 0 {
		o.Width = 1024
	}

	// Default height
	if o.Height == 0 {
		o.Height = 768
	}
}
