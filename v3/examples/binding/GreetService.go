package main

import "github.com/wailsapp/wails/v3/examples/binding/data"

// GreetService is a service that greets people
type GreetService struct {
}

// Greet greets a person
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

// GreetPerson greets a person
func (*GreetService) GreetPerson(person data.Person) string {
	return "Hello " + person.Name
}
