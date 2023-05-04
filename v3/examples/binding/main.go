package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Bind: []interface{}{
			&GreetService{},
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		URL: "/",
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
