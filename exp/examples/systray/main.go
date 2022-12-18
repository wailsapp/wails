package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/exp/pkg/options"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {
	app := application.New()
	app.SetActivationPolicy(options.ActivationPolicyAccessory)

	systemTray := app.NewSystemTray().SetIcon(application.DefaultMacTemplateIcon)

	myMenu := app.NewMenu()
	myMenu.Add("Hello World!").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().SetLabel("Clicked!")
	})
	subMenu := myMenu.AddSubmenu("Submenu")
	subMenu.Add("Submenu Item")
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
