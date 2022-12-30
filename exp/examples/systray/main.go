package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/exp/pkg/options"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {
	app := application.New(options.Application{
		Name:        "Systray Demo",
		Description: "A demo of the Systray API",
		Mac: options.Mac{
			ActivationPolicy: options.ActivationPolicyAccessory,
		},
	})

	systemTray := app.NewSystemTray()
	if runtime.GOOS == "darwin" {
		systemTray.SetIcon(application.DefaultMacTemplateIcon)
	}

	myMenu := app.NewMenu()
	myMenu.Add("Hello World!").OnClick(func(ctx *application.Context) {
		app.NewInfoDialog().SetTitle("Hello World!").SetMessage("Hello World!").Show()
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

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
