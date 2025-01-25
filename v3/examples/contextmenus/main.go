package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Context Menu Demo",
		Description: "A demo of the Context Menu API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Context Menu Demo",
		Width:  1024,
		Height: 1024,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	contextMenu := application.NewContextMenu("test")
	clickMe := contextMenu.Add("Click to set Menuitem label to Context Data")
	contextDataMenuItem := contextMenu.Add("Current context data: No Context Data")
	clickMe.OnClick(func(data *application.Context) {
		app.Logger.Info("Context menu", "context data", data.ContextMenuData())
		contextDataMenuItem.SetLabel("Current context data: " + data.ContextMenuData())
		contextMenu.Update()
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
