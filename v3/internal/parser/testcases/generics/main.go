package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService int

// A nice type Alias
type Alias = int

// A generic struct
type Person[T any] struct {
	Name         T
	AliasedField Alias
}

// Get someone
func (GreetService) Get(aliasValue Alias) Person[string] {
	return Person[string]{"hello", aliasValue}
}

// Get someone quite different
func (GreetService) GetButDifferent() Person[bool] {
	return Person[bool]{true, 13}
}

func NewGreetService() application.Service {
	return application.NewService(new(GreetService))
}

func main() {
	app := application.New(application.Options{
		Bind: []application.Service{
			NewGreetService(),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
