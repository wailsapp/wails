package main

import (
	"embed"
	_ "embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Context Menu Demo",
		Description: "A demo of the Context Menu API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	mainWindow := app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title: "Context Menu Demo",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title: "Context Menu Demo",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	contextMenu := app.NewMenu()
	contextMenu.Add("Click Me").OnClick(func(data *application.Context) {
		fmt.Printf("Context menu data: %+v\n", data.ContextMenuData())
	})

	globalContextMenu := app.NewMenu()
	globalContextMenu.Add("Default context menu item").OnClick(func(data *application.Context) {
		fmt.Printf("Context menu data: %+v\n", data.ContextMenuData())
	})

	// Registering the menu with a window will make it available to that window only
	mainWindow.RegisterContextMenu("test", contextMenu)

	// Registering the menu with the app will make it available to all windows
	app.RegisterContextMenu("test", globalContextMenu)

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
