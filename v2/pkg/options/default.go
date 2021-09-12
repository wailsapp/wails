package options

import (
	"github.com/wailsapp/wails/v2/pkg/logger"
)

// Default options for creating the App
var Default = &App{
	Width:    1024,
	Height:   768,
	Logger:   logger.NewDefaultLogger(),
	LogLevel: logger.INFO,
}
