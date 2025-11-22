package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is a service that demonstrates bound methods over WebSocket transport
type GreetService struct {
	mu         sync.Mutex
	greetCount int
	app        *application.App
}

// ServiceStartup is called when the service is initialized
func (g *GreetService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	g.app = application.Get()

	// Start a timer that emits events every second
	// This demonstrates automatic event forwarding to WebSocket transport
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				// Emit a timer event - automatically forwarded to WebSocket!
				g.app.Event.Emit("timer:tick", t.Format("15:04:05"))
			}
		}
	}()

	return nil
}

// Greet greets a person by name and emits an event
func (g *GreetService) Greet(name string) string {
	g.mu.Lock()
	g.greetCount++
	count := g.greetCount
	g.mu.Unlock()
	result := fmt.Sprintf("Hello, %s! (Greeted %d times via WebSocket)", name, count)

	// Emit an event to demonstrate event support over WebSocket
	// Events are automatically forwarded to the WebSocket transport!
	if g.app != nil {
		g.app.Event.Emit("greet:count", count)
	}

	return result
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
