package main

import (
	"log"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {
	app := application.New()

	app.NewWindow()
	app.NewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
