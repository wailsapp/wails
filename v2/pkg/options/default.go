package options

import (
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Default options for creating the App
var Default = &App{
	Width:              1024,
	Height:             768,
	Logger:             logger.NewDefaultLogger(),
	LogLevel:           logger.INFO,
	LogLevelProduction: logger.ERROR,
	CSSDragProperty:    "--wails-draggable",
	CSSDragValue:       "drag",
}

var defaultMacMenu = menu.NewMenuFromItems(
	menu.AppMenu(),
	menu.EditMenu(),
)
