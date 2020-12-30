package main

import (
	"github.com/wailsapp/wails/v2"
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

// Message Dialog
func (l *Dialog) Message(options *options.MessageDialog) string {
	return l.runtime.Dialog.Message(options)
}

// Message Dialog
func (l *Dialog) Test() string {
	return l.runtime.Dialog.Message(&options.MessageDialog{
		Type:    options.InfoDialog,
		Title:   " ",
		Message: "I am a longer message but these days, can't be too long!",
		// Buttons are declared in the order they should be appear in
		Buttons:       []string{"test", "Cancel", "OK"},
		DefaultButton: "OK",
		CancelButton:  "Cancel",
		//Icon:          "wails",
	})
}
