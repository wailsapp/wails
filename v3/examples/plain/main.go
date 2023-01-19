package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/options"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(options.Application{
		Name:        "Plain",
		Description: "A demo of using raw HTML & CSS",
		Mac: options.Mac{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	// Create window
	app.NewWebviewWindowWithOptions(&options.WebviewWindow{
		Title: "Plain Bundle",
		CSS:   `body { background-color: rgba(255, 255, 255, 0); } .main { color: white; margin: 20%; }`,
		Mac: options.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                options.MacBackdropTranslucent,
			TitleBar:                options.TitleBarHiddenInset,
		},

		URL: "/",
		Assets: options.Assets{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<html><head><title>Plain Bundle</title></head><body><div class="main"><h1>Plain Bundle</h1><p>This is a plain bundle. It has no frontend code but this was Served by the AssetServer's Handler</p></div></body></html>`))
			}),
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
