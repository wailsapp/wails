package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed index.html
var indexHTML string

func main() {
	app := application.New(application.Options{
		Name:        "Liquid Glass Demo",
		Description: "Demonstrates the Liquid Glass effect on macOS",
	})

	// Create main window with simple liquid glass
	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Liquid Glass - Simple",
		Width:  400,
		Height: 300,
		X:      100,
		Y:      100,
		HTML:   indexHTML,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropLiquidGlass,
		},
	})

	// Create second window with advanced liquid glass configuration  
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Liquid Glass - Advanced",
		Width:  400,
		Height: 300,
		X:      520,
		Y:      100,
		HTML:   indexHTML,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropLiquidGlass,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleVibrant,
				CornerRadius: 16.0,
				TintColor:    &application.RGBA{0, 122, 255, 50}, // Blue tint with transparency
				GroupID:      "main-group",
				GroupSpacing: 8.0,
			},
		},
	}).Show()

	// Create third window with different style
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Liquid Glass - Dark",
		Width:  400,
		Height: 300,
		X:      310,
		Y:      420,
		HTML:   indexHTML,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropLiquidGlass,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleDark,
				CornerRadius: 20.0,
				TintColor:    &application.RGBA{255, 0, 255, 30}, // Magenta tint
				GroupID:      "secondary-group",
			},
		},
	}).Show()

	// Show the main window
	mainWindow.Show()

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}