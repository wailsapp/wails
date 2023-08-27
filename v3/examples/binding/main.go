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
		Bind: []any{
			&GreetService{},
		},
		BindAliases: map[uint32]uint32{
			1: 1411160069,
			2: 4021313248,
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
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
