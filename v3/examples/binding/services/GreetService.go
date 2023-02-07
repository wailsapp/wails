package services

import (
	"github.com/ciderapp/wails/v3/examples/binding/models"
)

type GreetService struct {
	SomeVariable int
	lowercase    string
	Parent       *models.Person
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func (g *GreetService) GetPerson() *models.Person {
	return g.Parent
}

func (g *GreetService) SetPerson(person *models.Person) {
	g.Parent = person
}
