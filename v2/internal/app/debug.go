// +build debug

package app

import (
	"flag"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
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
	loglevel := flag.String("loglevel", "debug", "Loglevel to use - Trace, Debug, Info, Warning, Error")
	flag.Parse()
	if len(*loglevel) > 0 {
		switch strings.ToLower(*loglevel) {
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
