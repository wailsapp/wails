package options

import (
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

// Default options for creating the App
var Default = &App{
	Title:    "My Wails App",
	Width:    1024,
	Height:   768,
	DevTools: false,
	RGBA:     0xFFFFFFFF,
	Mac: &mac.Options{
		TitleBar:                      mac.TitleBarDefault(),
		Appearance:                    mac.DefaultAppearance,
		WebviewIsTransparent:          false,
		WindowBackgroundIsTranslucent: false,
	},
	Logger:   logger.NewDefaultLogger(),
	LogLevel: logger.INFO,
}
