package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
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
func (l *Dialog) Open(options *dialog.OpenDialog) []string {
	return l.runtime.Dialog.Open(options)
}

// Save Dialog
func (l *Dialog) Save(options *dialog.SaveDialog) string {
	return l.runtime.Dialog.Save(options)
}

// Message Dialog
func (l *Dialog) Message(options *dialog.MessageDialog) string {
	return l.runtime.Dialog.Message(options)
}
