package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/exp/pkg/application"
	"github.com/wailsapp/wails/exp/pkg/options"
)

type GreetService struct {
	SomeVariable int
	lowerCase    string
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

type OtherService struct {
	t int
}

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
