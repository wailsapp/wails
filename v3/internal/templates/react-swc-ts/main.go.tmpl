package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/options"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	app := application.New(options.Application{
		Name:        "{{.ProjectName}}",
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
			FS: assets,
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
