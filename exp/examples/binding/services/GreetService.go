package services

import "github.com/wailsapp/wails/exp/examples/binding/models"

type GreetService struct {
	SomeVariable int
	lowercase    string
	Parent       *models.Person
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}
