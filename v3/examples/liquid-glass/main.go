package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed index.html
var indexHTML string

//go:embed wails-logo.png
var wailsLogo []byte

func main() {
	app := application.New(application.Options{
		Name:        "Wails Liquid Glass Demo",
		Description: "Demonstrates the native Liquid Glass effect on macOS",
	})

	// Check if running on macOS
	if runtime.GOOS != "darwin" {
		// Show dialog for non-macOS platforms
		app.Dialog.Info().
			SetTitle("macOS Only Demo").
			SetMessage("The Liquid Glass effect is a macOS-specific feature that uses native NSGlassEffectView (macOS 15.0+) or NSVisualEffectView.\n\nThis demo is not available on " + runtime.GOOS + ".").
			Show()
		fmt.Println("The Liquid Glass effect is a macOS-specific feature. This demo is not available on", runtime.GOOS)
		return
	}

	// Convert logo to base64 data URI
	logoDataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(wailsLogo)

	// Create different HTML for each window
	lightHTML := strings.Replace(indexHTML, "wails-logo.png", logoDataURI, 1)
	lightHTML = strings.Replace(lightHTML, "LIQUID GLASS", "Light Style", 1)

	darkHTML := strings.Replace(indexHTML, "wails-logo.png", logoDataURI, 1)
	darkHTML = strings.Replace(darkHTML, "LIQUID GLASS", "Dark Style", 1)

	vibrantHTML := strings.Replace(indexHTML, "wails-logo.png", logoDataURI, 1)
	vibrantHTML = strings.Replace(vibrantHTML, "LIQUID GLASS", "Vibrant Style", 1)

	tintedHTML := strings.Replace(indexHTML, "wails-logo.png", logoDataURI, 1)
	tintedHTML = strings.Replace(tintedHTML, "LIQUID GLASS", "Blue Tint", 1)

	sheetHTML := strings.Replace(indexHTML, "wails-logo.png", logoDataURI, 1)
	sheetHTML = strings.Replace(sheetHTML, "LIQUID GLASS", "Sheet Material", 1)

	hudHTML := strings.Replace(indexHTML, "wails-logo.png", logoDataURI, 1)
	hudHTML = strings.Replace(hudHTML, "LIQUID GLASS", "HUD Window", 1)

	contentHTML := strings.Replace(indexHTML, "wails-logo.png", logoDataURI, 1)
	contentHTML = strings.Replace(contentHTML, "LIQUID GLASS", "Content Background", 1)

	// Window 1: Light style with no tint
	window1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Light Glass",
		Width:             350,
		Height:            280,
		X:                 100,
		Y:                 100,
		Frameless:         true,
		EnableDragAndDrop: false,
		HTML:              lightHTML,
		InitialPosition:   application.WindowXY,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropLiquidGlass,
			InvisibleTitleBarHeight: 500,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleLight,
				Material:     application.NSVisualEffectMaterialAuto,
				CornerRadius: 20.0,
				TintColor:    nil,
			},
		},
	})

	// Window 2: Dark style
	window2 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Dark Glass",
		Width:             350,
		Height:            280,
		X:                 500,
		Y:                 100,
		Frameless:         true,
		EnableDragAndDrop: false,
		HTML:              darkHTML,
		InitialPosition:   application.WindowXY,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropLiquidGlass,
			InvisibleTitleBarHeight: 500,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleDark,
				Material:     application.NSVisualEffectMaterialAuto,
				CornerRadius: 20.0,
				TintColor:    nil,
			},
		},
	})

	// Window 3: Vibrant style
	window3 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Vibrant Glass",
		Width:             350,
		Height:            280,
		X:                 900,
		Y:                 100,
		Frameless:         true,
		EnableDragAndDrop: false,
		HTML:              vibrantHTML,
		InitialPosition:   application.WindowXY,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropLiquidGlass,
			InvisibleTitleBarHeight: 500,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleVibrant,
				Material:     application.NSVisualEffectMaterialAuto,
				CornerRadius: 20.0,
				TintColor:    nil,
			},
		},
	})

	// Window 4: Blue tinted glass
	window4 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Tinted Glass",
		Width:             350,
		Height:            280,
		X:                 300,
		Y:                 420,
		Frameless:         true,
		EnableDragAndDrop: false,
		HTML:              tintedHTML,
		InitialPosition:   application.WindowXY,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropLiquidGlass,
			InvisibleTitleBarHeight: 500,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleLight,
				Material:     application.NSVisualEffectMaterialAuto,
				CornerRadius: 25.0,                                                        // Different corner radius
				TintColor:    &application.RGBA{Red: 0, Green: 100, Blue: 200, Alpha: 50}, // Blue tint
			},
		},
	})

	// Window 5: Using specific NSVisualEffectMaterial
	window5 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Sheet Material",
		Width:             350,
		Height:            280,
		X:                 700,
		Y:                 420,
		Frameless:         true,
		EnableDragAndDrop: false,
		HTML:              sheetHTML,
		InitialPosition:   application.WindowXY,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropLiquidGlass,
			InvisibleTitleBarHeight: 500,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleAutomatic,   // Automatic style
				Material:     application.NSVisualEffectMaterialSheet, // Specific material
				CornerRadius: 15.0,                                    // Different corner radius
				TintColor:    nil,
			},
		},
	})

	// Window 6: HUD Window Material (very light, translucent)
	window6 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "HUD Window",
		Width:             350,
		Height:            280,
		X:                 100,
		Y:                 740,
		Frameless:         true,
		EnableDragAndDrop: false,
		HTML:              hudHTML,
		InitialPosition:   application.WindowXY,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropLiquidGlass,
			InvisibleTitleBarHeight: 500,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleAutomatic,
				Material:     application.NSVisualEffectMaterialHUDWindow, // HUD Window material - very light
				CornerRadius: 30.0,                                        // Larger corner radius
				TintColor:    nil,
			},
		},
	})

	// Window 7: Content Background Material
	window7 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Content Background",
		Width:             350,
		Height:            280,
		X:                 500,
		Y:                 740,
		Frameless:         true,
		EnableDragAndDrop: false,
		HTML:              contentHTML,
		InitialPosition:   application.WindowXY,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropLiquidGlass,
			InvisibleTitleBarHeight: 500,
			LiquidGlass: application.MacLiquidGlass{
				Style:        application.LiquidGlassStyleAutomatic,
				Material:     application.NSVisualEffectMaterialContentBackground,         // Content background
				CornerRadius: 10.0,                                                        // Smaller corner radius
				TintColor:    &application.RGBA{Red: 0, Green: 200, Blue: 100, Alpha: 30}, // Warm tint
			},
		},
	})

	// Show all windows
	window1.Show()
	window2.Show()
	window3.Show()
	window4.Show()
	window5.Show()
	window6.Show()
	window7.Show()

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
