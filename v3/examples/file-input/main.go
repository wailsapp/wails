package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

type FileService struct{}

func (f *FileService) OpenDirectoryOnly() string {
	result, err := application.Get().Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("Select Directory").
		PromptForSingleSelection()
	if err != nil {
		return "Error: " + err.Error()
	}
	if result == "" {
		return "Cancelled"
	}
	return result
}

func (f *FileService) OpenFilteredFile() string {
	result, err := application.Get().Dialog.OpenFile().
		SetTitle("Select Text File").
		AddFilter("Text Files", "*.txt").
		PromptForSingleSelection()
	if err != nil {
		return "Error: " + err.Error()
	}
	if result == "" {
		return "Cancelled"
	}
	return result
}

func main() {
	app := application.New(application.Options{
		Name:        "File Input Test",
		Description: "Test for HTML file input (#4862)",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Services: []application.Service{
			application.NewService(&FileService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "File Input Test",
		Width:  700,
		Height: 500,
		URL:    "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
