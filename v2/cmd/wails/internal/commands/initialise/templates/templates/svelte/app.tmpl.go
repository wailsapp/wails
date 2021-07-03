package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2"
)

// App application struct
type App struct {
	runtime *wails.Runtime
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(runtime *wails.Runtime) {
	// Perform your setup here
	b.runtime = runtime
}

// shutdown is called at application termination
func (b *App) shutdown() {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (b *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
