package main

import (
	"embed"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

var windowCount int32 = 1
var app *application.App

// WindowService allows the frontend to open new windows
type WindowService struct{}

func (s *WindowService) OpenNewWindow() int {
	count := atomic.AddInt32(&windowCount, 1)
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Broadcast Channel Demo",
		Width:  800,
		Height: 550,
		URL:    "/",
	})
	return int(count)
}

func main() {
	app = application.New(application.Options{
		Name:        "Broadcast Channel Demo",
		Description: "Cross-window communication via BroadcastChannel API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Services: []application.Service{
			application.NewService(&WindowService{}),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Broadcast Channel Demo",
		Width:  800,
		Height: 550,
		URL:    "/",
	})

	app.Run()
}
