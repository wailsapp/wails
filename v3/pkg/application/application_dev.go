//go:build !production

package application

import (
	"errors"
	"github.com/wailsapp/wails/v3/internal/assetserver"
	"net/http"
	"time"
)

func (a *App) preRun() error {
	// Check for frontend server url
	frontendURL := assetserver.GetDevServerURL()
	if frontendURL != "" {
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
		a.Logger.Info("failed!")
		return errors.New("unable to connect to frontend server. Please check it is running")
	}
	return nil
}
