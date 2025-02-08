package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Title is a title
type Title string

const (
	// Mister is a title
	Mister Title = "Mr"
	Miss   Title = "Miss"
	Ms     Title = "Ms"
	Mrs    Title = "Mrs"
	Dr     Title = "Dr"
)

// Age is an integer with some predefined values
type Age = int

const (
	NewBorn    Age = 0
	Teenager   Age = 12
	YoungAdult Age = 18

	// Oh no, some grey hair!
	MiddleAged Age = 50
	Mathusalem Age = 1000 // Unbelievable!
)

// GreetService is great
type GreetService struct {
	SomeVariable int
	lowerCase    string
	target       *Person
}

// Person represents a person
type Person struct {
	Title Title
	Name  string
	Age   Age
}

// Greet does XYZ
func (*GreetService) Greet(name string, title Title) string {
	return "Hello " + string(title) + " " + name
}

// NewPerson creates a new person
func (*GreetService) NewPerson(name string) *Person {
	return &Person{Name: name}
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
