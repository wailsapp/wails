package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
}

// Notify emits an event to the frontend
func (a *App) Notify(message string) {
	runtime.EventsEmit(a.ctx, "notify", message)
	runtime.WindowSetTitle(a.ctx, message)
}

// GreetService greets people
type GreetService struct{}

// Greet returns a greeting
func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
