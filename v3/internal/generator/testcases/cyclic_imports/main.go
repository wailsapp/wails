package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService int

type StructA struct {
	B *structB
}

type structB struct {
	A *StructA
}

type StructC struct {
	D structD
}

type structD struct {
	E StructE
}

type StructE struct{}

// Make a cycle.
func (GreetService) MakeCycles() (_ StructA, _ StructC) {
	return
}

func NewGreetService() application.Service {
	return application.NewService(new(GreetService))
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
