package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct {
	AnInterface
}

type AnInterface interface {
	// Comment 1.
	Method1()

	Method2() // Comment 2.

	// Comment 3a.
	Method3() // Comment 3b.

	interface {
		// Comment 4.
		Method4()
	}

	InterfaceAlias
}

type InterfaceAlias = interface {
	// Comment 5.
	Method5()
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
