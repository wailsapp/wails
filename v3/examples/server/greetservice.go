package main

// Person holds someone's most important attributes
type Person struct {
	// Name is the person's name
	Name string `json:"name"`
}

// GreetService is a service that greets people
type GreetService struct {
}

// Greet greets a person
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

// GreetPerson greets a person
func (*GreetService) GreetPerson(person Person) string {
	return "Hello " + person.Name
}
