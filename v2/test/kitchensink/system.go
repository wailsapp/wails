package main

import (
	wails "github.com/wailsapp/wails/v2"
)

// System struct
type System struct {
	runtime *wails.Runtime
}

// WailsInit is called at application startup
func (l *System) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	l.runtime = runtime
	return nil
}

// Platform will return the runtime platform value
func (l *System) Platform() string {
	// Perform your setup here
	return l.runtime.System.Platform()
}
