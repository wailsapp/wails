package main

import (
	"context"
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

// Greet someone
func (*GreetService) GreetWithContext(ctx context.Context, name string) string {
	return "Hello " + name
}

func NewGreetService() application.Service {
	return application.NewService(&GreetService{})
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			NewGreetService(),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
