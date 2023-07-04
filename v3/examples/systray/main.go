package main

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
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

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:       500,
		Height:      800,
		Frameless:   true,
		AlwaysOnTop: true,
		Hidden:      true,
		ShouldClose: func(window *application.WebviewWindow) bool {
			window.Hide()
			return false
		},
	})

	window.On(events.Common.WindowLostFocus, func(ctx *application.WindowEventContext) {
		window.Hide()
	})

	app.On(events.Mac.ApplicationDidResignActiveNotification, func() {
		window.Hide()
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
		if window.IsVisible() {
			window.Hide()
			return
		}
		err := systemTray.PositionWindow(window, 5)
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
