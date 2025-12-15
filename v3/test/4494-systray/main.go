package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

//go:embed assets/icon.png
var icon []byte

func main() {
	app := application.New(application.Options{
		Name:        "Systray Test (#4494)",
		Description: "Test for systray context menu hiding attached window",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	// Create the main window
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Systray Test (#4494)",
		Width:  400,
		Height: 300,
	})

	// Create a systray menu
	menu := app.NewMenu()
	menu.Add("Show Window").OnClick(func(ctx *application.Context) {
		log.Println("Show Window clicked")
		window.Show()
	})
	menu.Add("Hide Window").OnClick(func(ctx *application.Context) {
		log.Println("Hide Window clicked")
		window.Hide()
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		log.Println("Quit clicked")
		app.Quit()
	})

	// Create system tray with attached window
	// Issue #4494: Right-clicking to open context menu hides the attached window
	systemTray := app.SystemTray.New()
	systemTray.SetIcon(icon)
	systemTray.SetMenu(menu)
	systemTray.AttachWindow(window)

	log.Println("Starting application...")
	log.Println("TEST: Right-click the systray icon. The window should NOT hide when the context menu opens.")
	log.Println("BUG: On Linux (especially Wayland/KDE), the window hides when right-clicking the systray icon.")

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
