//go:build darwin

package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

// PrintService provides print functionality to the frontend
type PrintService struct {
	app *application.App
}

func (p *PrintService) Print() error {
	if w := p.app.Window.Current(); w != nil {
		log.Println("PrintService.Print() called")
		return w.Print()
	}
	return nil
}

func main() {
	// Only run on macOS
	if runtime.GOOS != "darwin" {
		log.Fatal("This test is only for macOS")
	}

	printService := &PrintService{}

	app := application.New(application.Options{
		Name:        "Print Dialog Test",
		Description: "Test for macOS print dialog (Issue #4290)",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Services: []application.Service{
			application.NewService(printService),
		},
	})

	printService.app = app

	// Create application menu
	menu := app.NewMenu()

	// File menu
	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("Print...").
		SetAccelerator("CmdOrCtrl+P").
		OnClick(func(ctx *application.Context) {
			if w := app.Window.Current(); w != nil {
				log.Println("Attempting to print...")
				if err := w.Print(); err != nil {
					log.Printf("Print error: %v", err)
				} else {
					log.Println("Print completed (or dialog dismissed)")
				}
			}
		})
	fileMenu.AddSeparator()
	fileMenu.Add("Quit").
		SetAccelerator("CmdOrCtrl+Q").
		OnClick(func(ctx *application.Context) {
			app.Quit()
		})

	app.Menu.Set(menu)

	// Create main window
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Print Dialog Test - Issue #4290",
		Width:  800,
		Height: 600,
		URL:    "/index.html",
	})

	log.Println("Starting application. Use File > Print or Cmd+P to test print dialog.")

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
