package main

type Person struct {
	name string
}

// GreetService is a service that greets people
type GreetService struct {
}

// Greet greets a person
//
//wails:methodID 1
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

// GreetPerson greets a person
//
//wails:methodID 2
func (*GreetService) GreetPerson(person Person) string {
	return "Hello " + person.name
}
