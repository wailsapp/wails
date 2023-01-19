package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/examples/binding/services"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/options"
)

type localStruct struct{}

func main() {
	app := application.New(options.Application{
		Bind: []interface{}{
			&localStruct{},
			&services.GreetService{},
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
