package options

import "github.com/wailsapp/wails/v2/pkg/options/mac"

// Default options for creating the App
var Default = &App{
	Title:    "My Wails App",
	Width:    1024,
	Height:   768,
	DevTools: true,
	Colour:   0xFFFFFFFF,
	Mac: &mac.Options{
		TitleBar: mac.TitleBarDefault(),
	},
}
