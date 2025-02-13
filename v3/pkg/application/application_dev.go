//go:build !production

package application

import (
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

var devMode = false

func (a *App) preRun() error {
	// Check for frontend server url
	frontendURL := assetserver.GetDevServerURL()
	if frontendURL != "" {
		devMode = true
		// We want to check if the frontend server is running by trying to http get the url
		// and if it is not, we wait 500ms and try again for a maximum of 10 times. If it is
		// still not available, we return an error.
		// This is to allow the frontend server to start up before the backend server.
		client := http.Client{}
		a.Logger.Info("Waiting for frontend dev server to start...", "url", frontendURL)
		for i := 0; i < 10; i++ {
			_, err := client.Get(frontendURL)
			if err == nil {
				a.Logger.Info("Connected to frontend dev server!")
				return nil
			}
			// Wait 500ms
			time.Sleep(500 * time.Millisecond)
			if i%2 == 0 {
				a.Logger.Info("Retrying...")
			}
		}
		a.fatal("unable to connect to frontend server. Please check it is running - FRONTEND_DEVSERVER_URL='%s'", frontendURL)
	}
	return nil
}

func (a *App) postQuit() {
	if devMode {
		a.Logger.Info("The application has terminated, but the watcher is still running.")
		a.Logger.Info("To terminate the watcher, press CTRL+C")
	}
}

func (a *App) enableDevTools() {

}
