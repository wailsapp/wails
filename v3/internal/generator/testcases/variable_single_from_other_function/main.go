package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/internal/generator/testcases/variable_single_from_other_function/services"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct {
	SomeVariable int
	lowerCase    string
	target       *Person
}

// Person is a person!
// They have a name and an address
type Person struct {
	Name    string
	Address *services.Address
}

// Greet does XYZ
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

// NewPerson creates a new person
func (*GreetService) NewPerson(name string) *Person {
	return &Person{Name: name}
}

func main() {
	otherService := services.NewOtherService()
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&GreetService{}),
			otherService,
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
