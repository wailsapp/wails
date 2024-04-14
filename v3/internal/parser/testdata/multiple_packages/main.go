package main

import (
	"embed"
	"log"
	"log/slog"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/parser/testdata/multiple_packages/other"
	otherother "github.com/wailsapp/wails/v3/internal/parser/testdata/multiple_packages/other/other"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type GreetService struct {
	ServiceName string
}

type Greeter GreetService

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + ", my name is " + g.ServiceName
}

func (g *GreetService) BuildInfo() (*debug.BuildInfo, bool) {
	return debug.ReadBuildInfo()
}

func (g *GreetService) UUID() uuid.UUID {
	return uuid.New()
}

//go:embed frontend/*
var assets embed.FS

func main() {
	tupel := lo.T2(0, 1)

	app := application.New(application.Options{
		Bind: []interface{}{
			log.Default(),
			&GreetService{
				ServiceName: "GreetService",
			},
			&Greeter{
				ServiceName: "Greeter",
			},
			&other.OtherService{},
			&otherother.OtherService{},
			&tupel,
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		LogLevel: slog.LevelWarn,
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		URL:             "/",
		DevToolsEnabled: true,
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
