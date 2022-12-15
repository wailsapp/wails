package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/exp/pkg/events"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {
	app := application.New()

	// Create window
	myWindow := app.NewWindow()
	myWindow.On(events.Mac.WindowDidBecomeMain, func() {
		println("Window did become main")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
