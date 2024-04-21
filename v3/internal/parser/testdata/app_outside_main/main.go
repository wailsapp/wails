package main

import (
	"log"

	"github.com/wailsapp/wails/v3/internal/parser/testdata/app_outside_main/app"
)

func main() {
	app := app.NewApp()
	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
