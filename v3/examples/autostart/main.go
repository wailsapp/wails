package main

import (
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Autostart Demo",
		Description: "Toggle whether the app launches on login",
		Assets:      application.AlphaAssets,
	})

	menu := app.NewMenu()

	menu.Add("Status").OnClick(func(_ *application.Context) {
		st, err := app.Autostart.Status()
		if err != nil {
			app.Dialog.Error().SetTitle("Status Error").SetMessage(err.Error()).Show()
			return
		}
		msg := fmt.Sprintf("Enabled: %v\nStrategy: %s\nPath: %s",
			st.Enabled, st.Strategy, st.Path)
		app.Dialog.Info().SetMessage(msg).Show()
	})

	menu.Add("Enable").OnClick(func(_ *application.Context) {
		if err := app.Autostart.Enable(); err != nil {
			app.Dialog.Error().SetTitle("Enable Failed").SetMessage(err.Error()).Show()
		}
	})

	menu.Add("Disable").OnClick(func(_ *application.Context) {
		if err := app.Autostart.Disable(); err != nil {
			app.Dialog.Error().SetTitle("Disable Failed").SetMessage(err.Error()).Show()
		}
	})

	menu.Add("Quit").OnClick(func(_ *application.Context) { app.Quit() })

	app.Menu.SetApplicationMenu(menu)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Autostart Demo",
		HTML:  `<h1>Autostart Demo</h1><p>Use the application menu to control autostart.</p>`,
	})
	_ = window

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
