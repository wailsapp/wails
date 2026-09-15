//go:build !production

package application

import "github.com/wailsapp/wails/v3/internal/assetserver"

func (a *App) postQuit() {
	if assetserver.GetDevServerURL() == "" {
		return
	}

	a.Logger.Info("The application has terminated, but the watcher is still running.")
	a.Logger.Info("To terminate the watcher, press CTRL+C")
}
