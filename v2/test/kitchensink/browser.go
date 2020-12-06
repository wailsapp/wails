package main

import (
	"github.com/wailsapp/wails/v2"
)

// Browser struct
type Browser struct {
	runtime *wails.Runtime
}

// WailsInit is called at application startup
func (l *Browser) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	l.runtime = runtime
	return nil
}

// Open will open the default browser with the given target
func (l *Browser) Open(target string) error {
	// Perform your setup here
	return l.runtime.Browser.Open(target)
}
