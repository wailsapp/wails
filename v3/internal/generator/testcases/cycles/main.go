package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService int

type Cyclic []map[string]Alias

type Alias = *Cyclic

type GenericCyclic[T any] []struct {
	X *GenericCyclic[T]
	Y []T
}

// Make a cycle.
func (GreetService) MakeCycles() (_ Cyclic, _ GenericCyclic[GenericCyclic[int]]) {
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
