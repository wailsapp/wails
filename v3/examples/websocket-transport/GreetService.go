package main

import (
	"fmt"
	"time"
)

// GreetService is a service that demonstrates bound methods over WebSocket transport
type GreetService struct {
	greetCount int
}

// Greet greets a person by name
func (g *GreetService) Greet(name string) string {
	g.greetCount++
	return fmt.Sprintf("Hello, %s! (Greeted %d times via WebSocket)", name, g.greetCount)
}

// GetTime returns the current server time
func (g *GreetService) GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Echo echoes back the input message
func (g *GreetService) Echo(message string) string {
	return "Echo: " + message
}

// Add adds two numbers together
func (g *GreetService) Add(a, b int) int {
	return a + b
}
