package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	runtime context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	//TODO: move to new runtime layout

	//b.runtime = runtime
	//runtime.Window.SetTitle("{{.ProjectName}}")
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (b *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
