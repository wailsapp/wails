package main

import "fmt"

var Greeting = "Hello"

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("%s %s!", Greeting, name)
}
