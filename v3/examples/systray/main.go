package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Systray Demo",
		Description: "A demo of the Systray API",
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	window := app.NewWebviewWindow().Hide()

	systemTray := app.NewSystemTray()
	if runtime.GOOS == "darwin" {
		systemTray.SetIcon(application.DefaultMacTemplateIcon)
	}

	myMenu := app.NewMenu()
	myMenu.Add("Hello World!").OnClick(func(ctx *application.Context) {
		app.InfoDialog().SetTitle("Hello World!").SetMessage("Hello World!").Show()
	})
	subMenu := myMenu.AddSubmenu("Submenu")
	subMenu.Add("Click me!").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().SetLabel("Clicked!")
	})
	myMenu.AddSeparator()
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systemTray.SetMenu(myMenu)

	systemTray.OnClick(func() {
		window.Show()
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
