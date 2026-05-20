package main

import (
	_ "embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type GreetService struct {
	SomeVariable int
	lowerCase    string
}

func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func NewGreetService() application.Service {
	return application.NewService(&GreetService{})
}
