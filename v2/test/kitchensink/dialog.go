package main

import (
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
	return l.runtime.Dialog.Open(options)
}

// Save Dialog
func (l *Dialog) Save(options *options.SaveDialog) string {
	return l.runtime.Dialog.Save(options)
}
