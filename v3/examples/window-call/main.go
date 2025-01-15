package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"math/rand"
	"runtime"
	"strconv"
)

//go:embed assets/*
var assets embed.FS

type WindowService struct{}

func (s *WindowService) RandomTitle(ctx context.Context) {
	callingWindow := ctx.Value(application.WindowKey).(application.Window)
	title := "Random Title " + strconv.Itoa(rand.Intn(1000))
	callingWindow.SetTitle(title)
}

// ==============================================

func main() {
	app := application.New(application.Options{
		Name:        "Window call Demo",
		Description: "A demo of the WebviewWindow API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Services: []application.Service{
			application.NewService(&WindowService{}),
		},
	})

	app.NewWebviewWindow().
		SetTitle("WebviewWindow 1").
		Show()

	// Create a custom menu
	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	windowCounter := 1

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("New")

	myMenu.Add("New WebviewWindow").
		SetAccelerator("CmdOrCtrl+N").
		OnClick(func(ctx *application.Context) {
			app.NewWebviewWindow().
				SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
				SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
				Show()
			windowCounter++
		})

	app.SetMenu(menu)
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
