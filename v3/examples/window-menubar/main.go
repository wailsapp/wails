package main

import (
	"embed"
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
)

//go:embed assets/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Window MenuBar Demo",
		Description: "A demo of menu bar toggling",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	// Create a menu
	menu := app.NewMenu()
	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	editMenu := menu.AddSubmenu("MenuBar")
	editMenu.Add("Hide MenuBar").OnClick(func(ctx *application.Context) {
		app.CurrentWindow().HideMenuBar()
	})

	helpMenu := menu.AddSubmenu("Help")
	helpMenu.Add("About").OnClick(func(ctx *application.Context) {
		app.CurrentWindow().SetURL("/about.html")
	})

	// Create window with menu
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Window MenuBar Demo",
		Width:  800,
		Height: 600,
		Windows: application.WindowsWindow{
			Menu: menu,
		},
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"F1": func(window *application.WebviewWindow) {
				window.ToggleMenuBar()
			},
			"F2": func(window *application.WebviewWindow) {
				window.ShowMenuBar()
			},
			"F3": func(window *application.WebviewWindow) {
				window.HideMenuBar()
			},
		},
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
