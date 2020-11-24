package main

import (
	"fmt"

	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
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

	// Setup Menu Listeners
	l.runtime.Menu.On("hello", func(m *menu.MenuItem) {
		fmt.Printf("The '%s' menu was clicked\n", m.Label)
	})
	l.runtime.Menu.On("checkbox-menu", func(m *menu.MenuItem) {
		fmt.Printf("The '%s' menu was clicked\n", m.Label)
		fmt.Printf("It is now %v\n", m.Checked)
		// m.Checked = false
		// runtime.Menu.Update()
	})

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
