package main

import (
	"strconv"

	"github.com/wailsapp/wails/v3/examples/binding/data"
)

// GreetService is a service that greets people
type GreetService struct {
}

// Greet greets a person
func (*GreetService) Greet(name string, counts ...int) string {
	times := " "

	for index, count := range counts {
		if index > 0 {
			times += ", "
		}
		times += strconv.Itoa(count)
	}

	if len(counts) > 0 {
		times += " times "
	}

	return "Hello" + times + name
}

// GreetPerson greets a person
func (srv *GreetService) GreetPerson(person data.Person) string {
	return srv.Greet(person.Name, person.Counts...)
}
