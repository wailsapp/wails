package main

import (
	_ "embed"
	"log"

	"github.com/ciderapp/wails/v3/pkg/application"
	"github.com/ciderapp/wails/v3/pkg/options"
)

type GreetService struct {
	SomeVariable int
	lowerCase    string
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func main() {
	app := application.New(options.Application{
		Bind: []interface{}{
			&GreetService{},
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
