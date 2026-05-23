//go:build !production

package application

import (
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

var devMode = false

// preRun is called before the application starts. In dev mode it waits for the
// frontend dev server to become reachable, retrying up to 10 times with
// exponential backoff (250 ms base, capped at 5 s per attempt).
func (a *App) preRun() error {
	frontendURL := assetserver.GetDevServerURL()
	if frontendURL == "" {
		return nil
	}

	devMode = true
	a.Logger.Info("Waiting for frontend dev server to start...", "url", frontendURL)

	const (
		maxRetries = 10
		maxDelay   = 5 * time.Second
	)
	client := http.DefaultClient
	delay := 250 * time.Millisecond

	for i := range maxRetries {
		_, err := client.Get(frontendURL)
		if err == nil {
			a.Logger.Info("Connected to frontend dev server!")
			return nil
		}
		a.Logger.Info("Retrying...", "attempt", i+1, "next_delay", delay)
		time.Sleep(delay)
		if delay *= 2; delay > maxDelay {
			delay = maxDelay
		}
	}

	a.fatal("unable to connect to frontend server. Please check it is running - FRONTEND_DEVSERVER_URL='%s'", frontendURL)
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
