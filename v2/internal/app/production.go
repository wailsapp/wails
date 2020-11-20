// +build !debug

package app

import "github.com/wailsapp/wails/v2/pkg/logger"

// Init initialises the application for a production environment
func (a *App) Init() error {
	a.logger.SetLogLevel(logger.ERROR)
	return nil
}
