package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "MenuWaylandTest",
		Description: "Test for window menu crash on Wayland",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	// Create a menu - this would crash on Wayland before the fix
	menu := app.NewMenu()

	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("New").OnClick(func(ctx *application.Context) {
		log.Println("New clicked")
	})
	fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
		log.Println("Open clicked")
	})
	fileMenu.AddSeparator()
	fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	editMenu := menu.AddSubmenu("Edit")
	editMenu.Add("Cut").OnClick(func(ctx *application.Context) {
		log.Println("Cut clicked")
	})
	editMenu.Add("Copy").OnClick(func(ctx *application.Context) {
		log.Println("Copy clicked")
	})
	editMenu.Add("Paste").OnClick(func(ctx *application.Context) {
		log.Println("Paste clicked")
	})

	helpMenu := menu.AddSubmenu("Help")
	helpMenu.Add("About").OnClick(func(ctx *application.Context) {
		log.Println("About clicked")
	})

	// Create window with menu attached via Linux options
	// This tests the fix for issue #4769
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Menu Wayland Test (#4769)",
		Width:  800,
		Height: 600,
		Linux: application.LinuxWindow{
			Menu: menu,
		},
	})

	log.Println("Starting application - if you see this on Wayland, the fix works!")
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
