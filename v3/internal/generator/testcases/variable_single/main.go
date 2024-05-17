package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct {
	SomeVariable int
	lowerCase    string
}

// Greet someone
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func main() {
	greetService := application.NewService(&GreetService{})
	app := application.New(application.Options{
		Services: []application.Service{
			greetService,
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
