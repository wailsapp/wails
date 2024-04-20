package main

import (
	_ "embed"
	"log"

	wails "github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct{}

// Greet does XYZ
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func main() {
	app := wails.New(wails.Options{
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
