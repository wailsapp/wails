package main

import (
	_ "embed"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/icons"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var counter int

func clickCount() int {
	counter++
	return counter
}

func main() {
	app := application.New(application.Options{
		Name:        "Systray Demo",
		Description: "A demo of the Systray API",
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:  500,
		Height: 800,
		//Frameless: true,
		Hidden: true,
	})

	systemTray := app.NewSystemTray()
	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	myMenu := app.NewMenu()
	myMenu.Add("Hello World!").OnClick(func(ctx *application.Context) {
		println("Hello World!")
		q := application.QuestionDialog().SetTitle("Ready?").SetMessage("Are you feeling ready?")
		q.AddButton("Yes").OnClick(func() {
			println("Awesome!")
		})
		q.AddButton("No").SetAsDefault().OnClick(func() {
			println("Boo!")
		})
		q.Show()
	})
	subMenu := myMenu.AddSubmenu("Submenu")
	subMenu.Add("Click me!").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().SetLabel("Clicked!")
	})
	myMenu.AddSeparator()
	myMenu.AddCheckbox("Checked", true).OnClick(func(ctx *application.Context) {
		println("Checked: ", ctx.ClickedMenuItem().Checked())
		application.InfoDialog().SetTitle("Hello World!").SetMessage("Hello World!").Show()
	})
	myMenu.Add("Enabled").OnClick(func(ctx *application.Context) {
		println("Click me!")
		ctx.ClickedMenuItem().SetLabel("Disabled!").SetEnabled(false)
	})
	myMenu.AddSeparator()
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	if runtime.GOOS != "darwin" {
		systemTray.SetMenu(myMenu)
	}

	showWindow := func() {
		window.SetTitle(fmt.Sprintf("Clicked %d times", clickCount()))
		err := systemTray.PositionWindow(window)
		if err != nil {
			application.InfoDialog().SetTitle("Error").SetMessage(err.Error()).Show()
			return
		}
		window.Show().Focus()
	}
	systemTray.OnClick(showWindow)

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
