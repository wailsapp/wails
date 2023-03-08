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
	})

	app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
