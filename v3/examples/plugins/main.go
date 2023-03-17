package main

import (
	"embed"
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"math/rand"
)

//go:embed assets/*
var assets embed.FS

type RandomNumberPlugin struct{}

func (r *RandomNumberPlugin) Name() string {
	return "Random Number Plugin"
}

func (r *RandomNumberPlugin) Init(_ *application.App) error {
	return nil
}

func (r *RandomNumberPlugin) Call(args []any) (any, error) {
	return rand.Intn(100), nil
}

func main() {

	app := application.New(application.Options{
		Name:        "Plugin Demo",
		Description: "A demo of the plugins API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Plugins: map[string]application.Plugin{
			"random": &RandomNumberPlugin{},
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	window := app.NewWebviewWindow()
	window.ToggleDevTools()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
