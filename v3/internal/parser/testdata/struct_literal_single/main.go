package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type GreetService struct {
	SomeVariable int
	lowerCase    string
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func main() {
	app := application.New(application.Options{
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
