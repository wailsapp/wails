// +build debug

package app

import (
	"flag"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"strings"
)

// Init initialises the application for a debug environment
func (a *App) Init() error {
	// Indicate debug mode
	a.debug = true

	if a.appType == "desktop" {
		// Enable dev tools
		a.options.DevTools = true
	}

	// Set log levels
	greeting := flag.String("loglevel", "debug", "Loglevel to use - Trace, Debug, Info, Warning, Error")
	flag.Parse()
	if len(*greeting) > 0 {
		switch strings.ToLower(*greeting) {
		case "trace":
			a.logger.SetLogLevel(logger.TRACE)
		case "info":
			a.logger.SetLogLevel(logger.INFO)
		case "warning":
			a.logger.SetLogLevel(logger.WARNING)
		case "error":
			a.logger.SetLogLevel(logger.ERROR)
		default:
			a.logger.SetLogLevel(logger.DEBUG)
		}
	}

	return nil
}
