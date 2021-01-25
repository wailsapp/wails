package main

import (
	"github.com/wailsapp/wails/v2"
	"log"
)

func main() {

	// Create application with options
	app, err := wails.CreateApp("{{.ProjectName}}", 1024, 768)
	if err != nil {
		log.Fatal(err)
	}

	app.Bind(newBasic())

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
