package main

import (
	_ "embed"
	"log"
	"net/http"

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
				w.Write([]byte(`<html><head><title>Plain Bundle</title></head><body><div class="main"><h1>Plain Bundle</h1><p>This is a plain bundle. It has no frontend code but this was Served by the AssetServer's Handler.</p><br/><br/><p data-wml-event="clicked">Clicking this paragraph emits an event...<p></div></body></html>`))
			}),
		},
	})
	// Create window
	app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title: "Plain Bundle",
		CSS:   `body { background-color: rgba(255, 255, 255, 0); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif; user-select: none; -ms-user-select: none; -webkit-user-select: none; } .main { color: white; margin: 20%; }`,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
		URL: "/",
	})

	app.Events.On("clicked", func(_ *application.WailsEvent) {
		println("clicked")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
