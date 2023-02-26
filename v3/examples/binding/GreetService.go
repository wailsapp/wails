package main

type GreetService struct {
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}
