package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService int

// Alpha sorts before Zeta, so the generator emits it first.
type Alpha struct {
	Name string
	Box  Epsilon[string]
}

type Epsilon[T any] struct {
	IntItems []int
	Zarray   []Zeta[T]
}

// Zeta is generic; []T and *T each need an element converter.
type Zeta[T any] struct {
	Items []T
	Ref   *Epsilon[T]
}

// Make a cycle.
func (GreetService) Greet() (_ Alpha) {
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

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
