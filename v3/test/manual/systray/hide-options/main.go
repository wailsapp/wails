package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:        "Systray Hide Options",
		Description: "Tests HideOnEscape and HideOnFocusLost options",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systray := app.SystemTray.New()

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Width:           400,
		Height:          300,
		Name:            "Hide Options Test",
		Title:           "Press Escape or click outside to hide",
		Frameless:       true,
		AlwaysOnTop:     true,
		Hidden:          true,
		DisableResize:   true,
		HideOnEscape:    true,
		HideOnFocusLost: true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systray.AttachWindow(window).WindowOffset(5)

	log.Println("Hide options test started")
	log.Println("Expected behavior:")
	log.Println("  - Left-click systray: Toggle window")
	log.Println("  - Press Escape: Hide window (HideOnEscape)")
	log.Println("  - Click outside window: Hide window (HideOnFocusLost)")
	log.Println("")
	log.Println("NOTE: On focus-follows-mouse WMs (Hyprland, Sway, i3),")
	log.Println("      HideOnFocusLost is automatically disabled to prevent")
	log.Println("      immediate hiding when mouse moves away.")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
