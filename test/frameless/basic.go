package main

import (
	"fmt"

	wails "github.com/wailsapp/wails/v2"
)

// Basic application struct
type Basic struct {
	runtime *wails.Runtime
}

// newBasic creates a new Basic application struct
func newBasic() *Basic {
	return &Basic{}
}

// WailsInit is called at application startup
func (b *Basic) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	b.runtime = runtime
	return nil
}

// WailsShutdown is called at application termination
func (b *Basic) WailsShutdown() {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (b *Basic) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

// Close shuts down the application
func (b *Basic) Close() {
	b.runtime.Quit()
}
