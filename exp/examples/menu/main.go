package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {
	app := application.New()

	myMenu := app.NewMenu()
	myMenu.AddRole(application.AppMenu)
	myMenu.AddRole(application.EditMenu)

	app.SetMenu(myMenu)

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
