package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2"
)

// Basic application struct
type Basic struct {
	runtime *wails.Runtime
}

// NewBasic creates a new Basic application struct
func NewBasic() *Basic {
	return &Basic{}
}

// startup is called at application startup
func (b *Basic) startup(runtime *wails.Runtime) {
	// Perform your setup here
	b.runtime = runtime
	runtime.Window.SetTitle("{{.ProjectName}}")
}

// shutdown is called at application termination
func (b *Basic) shutdown() {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (b *Basic) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
