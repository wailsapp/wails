package main

import (
	"context"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is a service that greets people
type GreetService struct {
}

// Greet greets a person
func (*GreetService) Greet(win application.Window, name string) string {
	return "Hello " + name + " on " + win.Name()
}

// GreetWithCtx greets a person
func (*GreetService) GreetWithCtx(ctx context.Context, name string) string {
	win := ctx.Value("window").(application.Window)
	return "[ctx] Hello " + name + " on " + win.Name()
}

// GreetWithCtx greets a person
func (*GreetService) GreetWithBoth(_ context.Context, win application.Window, name string) string {
	return "[ctx+win] Hello " + name + " on " + win.Name()
}
