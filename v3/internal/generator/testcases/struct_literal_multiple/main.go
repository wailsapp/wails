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

type OtherService struct {
	t int
}

func (o *OtherService) Hello() {}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(&OtherService{}),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
