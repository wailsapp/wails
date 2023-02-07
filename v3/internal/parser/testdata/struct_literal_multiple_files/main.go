package main

import (
	_ "embed"
	"log"

	"github.com/ciderapp/wails/v3/pkg/application"
	"github.com/ciderapp/wails/v3/pkg/options"
)

func main() {
	app := application.New(options.Application{
		Bind: []interface{}{
			&GreetService{},
			&OtherService{},
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
