//go:build !production

package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

var devMode = false

const (
	frontendDevServerRetryInterval = 500 * time.Millisecond
	frontendDevServerProbeTimeout  = 2 * time.Second
)

func waitForFrontendDevServer(ctx context.Context, client *http.Client, frontendURL string, retry func()) error {
	request, err := http.NewRequest(http.MethodGet, frontendURL, nil)
	if err != nil {
		return fmt.Errorf("invalid frontend dev server URL: %w", err)
	}

	for {
		response, err := client.Do(request.Clone(ctx))
		if err == nil {
			response.Body.Close()
			return nil
		}

		timer := time.NewTimer(frontendDevServerRetryInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			if retry != nil {
				retry()
			}
		}
	}
}

func (a *App) preRun() error {
	// Check for frontend server url
	frontendURL := assetserver.GetDevServerURL()
	if frontendURL != "" {
		devMode = true
		client := &http.Client{Timeout: frontendDevServerProbeTimeout}
		a.Logger.Info("Waiting for frontend dev server to start...", "url", frontendURL)
		retries := 0
		err := waitForFrontendDevServer(a.Context(), client, frontendURL, func() {
			retries++
			if retries%2 == 1 {
				a.Logger.Info("Retrying...")
			}
		})
		if err != nil {
			return fmt.Errorf("unable to connect to frontend server at FRONTEND_DEVSERVER_URL=%q: %w", frontendURL, err)
		}
		a.Logger.Info("Connected to frontend dev server!")
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
