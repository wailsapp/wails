package main

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
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
	greetService := &GreetService{}
	app := application.New(application.Options{
		Bind: []interface{}{
			greetService,
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
