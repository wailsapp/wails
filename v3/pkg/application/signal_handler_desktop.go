//go:build !ios

package application

import (
	"os"

	"github.com/wailsapp/wails/v3/internal/signal"
)

// setupSignalHandler sets up signal handling for desktop platforms
func (a *App) setupSignalHandler(options Options) {
	if !options.DisableDefaultSignalHandler {
		a.signalHandler = signal.NewSignalHandler(a.Quit)
		a.signalHandler.Logger = a.Logger
		a.signalHandler.ExitMessage = func(sig os.Signal) string {
			return "Quitting application..."
		}
	}
}
