package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
)

//go:embed assets/*
var assets embed.FS

type WindowService struct{}

func (s *WindowService) GeneratePanic() {
	s.call1()
}

func (s *WindowService) call1() {
	s.call2()
}

func (s *WindowService) call2() {
	panic("oh no! something went wrong deep in my service! :(")
}

// ==============================================

func main() {
	app := application.New(application.Options{
		Name:        "Panic Handler Demo",
		Description: "A demo of Handling Panics",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Services: []application.Service{
			application.NewService(&WindowService{}),
		},
		PanicHandler: func(panicDetails *application.PanicDetails) {
			fmt.Printf("*** There was a panic! ***\n")
			fmt.Printf("Time: %s\n", panicDetails.Time)
			fmt.Printf("Error: %s\n", panicDetails.Error)
			fmt.Printf("Stacktrace: %s\n", panicDetails.StackTrace)
			application.InfoDialog().SetMessage("There was a panic!").Show()
		},
	})

	app.NewWebviewWindow().
		SetTitle("WebviewWindow 1").
		Show()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
