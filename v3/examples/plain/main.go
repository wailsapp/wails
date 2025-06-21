package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Plain",
		Description: "A demo of using raw HTML & CSS",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<html><head><title>Plain Bundle</title><script type="module" src="/wails/runtime.js"></script></head><body><div class="main"><h1>Plain Bundle</h1><p>This is a plain bundle. It has no frontend code but this was Served by the AssetServer's Handler.</p><br/><br/><p data-wml-event="clicked">Clicking this paragraph emits an event...<p></div></body></html>`))
			}),
		},
	})
	// Create window - Note: In future versions, window creation may return errors
	// that should be checked. For now, window creation is deferred until app.Run()
	app.Windows.NewWithOptions(application.WebviewWindowOptions{
		Title: "Plain Bundle",
		CSS:   `body { background-color: rgb(255, 255, 255); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif; user-select: none; -ms-user-select: none; -webkit-user-select: none; } .main { color: white; margin: 20%; }`,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
		URL: "/",
	})

	// Create second window with direct HTML content
	app.Windows.NewWithOptions(application.WebviewWindowOptions{
		Title: "HTML TEST",
		HTML:  "<h1>AWESOME!</h1>",
		CSS:   `body { background-color: rgb(255, 0, 0); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif; user-select: none; -ms-user-select: none; -webkit-user-select: none; } .main { color: white; margin: 20%; }`,
		JS:    `window.iamhere = function() { console.log("Hello World!"); }`,
	})

	// Store the cleanup function to remove event listener when needed
	removeClickHandler := app.Events.On("clicked", func(_ *application.CustomEvent) {
		println("clicked")
	})
	// Note: In a real application, you would call removeClickHandler() when appropriate
	_ = removeClickHandler // Acknowledge we're storing the cleanup function

	// Use context-aware goroutine for graceful shutdown
	go func() {
		// Use a ticker instead of sleep to allow for cancellation
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		select {
		case <-ticker.C:
			// Create window after delay - in production, you should handle potential errors
			app.Windows.NewWithOptions(application.WebviewWindowOptions{
				Title:  "Plain Bundle new Window from GoRoutine",
				Width:  500,
				Height: 500,
				Mac: application.MacWindow{
					Backdrop:                application.MacBackdropTranslucent,
					TitleBar:                application.MacTitleBarHiddenInsetUnified,
					InvisibleTitleBarHeight: 50,
				},
			})
		case <-app.Context().Done():
			// Application is shutting down, cancel the goroutine
			return
		}
	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
