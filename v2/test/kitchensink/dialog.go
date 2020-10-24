package main

import (
	"fmt"

	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// Dialog struct
type Dialog struct {
	runtime *wails.Runtime
}

// WailsInit is called at application startup
func (l *Dialog) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	l.runtime = runtime
	return nil
}

// Open Dialog
func (l *Dialog) Open(options *options.OpenDialog) []string {
	fmt.Printf("%#v\n", options)
	// Perform your setup here
	return l.runtime.Dialog.Open(options)
}
