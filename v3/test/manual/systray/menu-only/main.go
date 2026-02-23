package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:        "Systray Menu Only",
		Description: "Tests systray with menu only (no attached window)",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systray := app.SystemTray.New()

	menu := app.Menu.New()
	menu.Add("Action 1").OnClick(func(ctx *application.Context) {
		log.Println("Action 1 clicked")
	})
	menu.Add("Action 2").OnClick(func(ctx *application.Context) {
		log.Println("Action 2 clicked")
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systray.SetMenu(menu)

	log.Println("Menu-only test started")
	log.Println("Expected behavior:")
	log.Println("  - Left-click: Nothing (no window)")
	log.Println("  - Right-click: Show menu")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
