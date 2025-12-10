//go:build ios

package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed test-assets/*
var assets embed.FS

type TestApp struct{}

func (a *TestApp) Greet(name string) string {
	return fmt.Sprintf("Hello %s from iOS build test!", name)
}

func main() {
	app := application.New(application.Options{
		Name:        "iOS Build Test",
		Description: "Testing iOS build system",
		Assets: application.AssetOptions{
			FS: assets,
		},
		Services: []application.Service{
			application.NewService(&TestApp{}),
		},
		LogLevel: application.LogLevelDebug,
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}