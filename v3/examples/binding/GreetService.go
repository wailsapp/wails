package main

type Person struct {
	name string
}

type GreetService struct {
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func (*GreetService) GreetPerson(person Person) string {
	return "Hello " + person.name
}
