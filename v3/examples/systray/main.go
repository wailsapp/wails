package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:        "Systray Demo",
		Description: "A demo of the Systray API",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systemTray := app.NewSystemTray()

	_ = app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:         500,
		Height:        500,
		Name:          "Systray Demo Window",
		Frameless:     true,
		AlwaysOnTop:   true,
		Hidden:        true,
		DisableResize: true,
		ShouldClose: func(window *application.WebviewWindow) bool {
			window.Hide()
			return false
		},
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"F12": func(window *application.WebviewWindow) {
				systemTray.OpenMenu()
			},
		},
	})

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	myMenu := app.NewMenu()
	myMenu.Add("Wails").SetBitmap(icons.WailsLogoBlackTransparent).SetEnabled(false)

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
	// Callbacks can be shared. This is useful for radio groups
	radioCallback := func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		menuItem.SetLabel(menuItem.Label() + "!")
	}

	// Radio groups are created implicitly by placing radio items next to each other in a menu
	myMenu.AddRadio("Radio 1", true).OnClick(radioCallback)
	myMenu.AddRadio("Radio 2", false).OnClick(radioCallback)
	myMenu.AddRadio("Radio 3", false).OnClick(radioCallback)

	myMenu.AddSeparator()
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systemTray.SetMenu(myMenu)

	//systemTray.AttachWindow(window).WindowOffset(5)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
