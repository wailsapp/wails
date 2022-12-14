package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/exp/pkg/application"
)

//go:embed macos_template_icon.png
var macosIcon []byte

func main() {
	app := application.New()
	systemTray := app.NewSystemTray().SetIcon(macosIcon)

	myMenu := app.NewMenu()
	myMenu.Add("Hello World!").OnClick(func(ctx *application.Context) {
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
