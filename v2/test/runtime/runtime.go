package main

import (
	"time"

	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// RuntimeTest to test the runtimes
type RuntimeTest struct {
	runtime *wails.Runtime
}

// WailsInit is an initialisation method
func (r *RuntimeTest) WailsInit(runtime *wails.Runtime) error {
	r.runtime = runtime
	println("Woohoo I'm here!")

	// Set title!
	// runtime.Window.SetTitle("My App!")

	r.runtime.Events.On("testevent", func(optionalParams ...interface{}) {
		println("Wooohoooo! I got called!")
	})
	r.runtime.Events.Once("testeventonce", func(optionalParams ...interface{}) {
		println("I only get called once!")
	})
	return nil
}

// WailsShutdown is called during shutdown
func (r *RuntimeTest) WailsShutdown() {
	// This is a test
	println("WOOOOOOOOOOOOOO WailsShutdown CALLED")
}

// EmitSingleEventNoParams will emit a single event with the given name and no params
func (r *RuntimeTest) EmitSingleEventNoParams(name string) {
	r.runtime.Events.Emit(name)
}

// EmitSingleEventWithOneParam will emit a single event with the given name
func (r *RuntimeTest) EmitSingleEventWithOneParam(name string) {
	r.runtime.Events.Emit(name, 1)
}

// RuntimeQuit will call the Runtime.Quit method
func (r *RuntimeTest) RuntimeQuit() {
	r.runtime.Quit()
}

// Fullscreen will call the Runtime.Fullscreen method
func (r *RuntimeTest) Fullscreen() {
	r.runtime.Window.Fullscreen()
}

// SetTitle will call the SetTitle method
func (r *RuntimeTest) SetTitle(title string) {
	r.runtime.Window.SetTitle(title)
}

// UnFullscreen will call the Runtime.UnFullscreen method
func (r *RuntimeTest) UnFullscreen() {
	r.runtime.Window.UnFullscreen()
}

// SetColour will call the Runtime.UnFullscreen method
func (r *RuntimeTest) SetColour(colour int) {
	r.runtime.Window.SetColour(colour)
}

// OpenFileDialog will call the Runtime.Dialog.OpenDialog method requesting File selection
func (r *RuntimeTest) OpenFileDialog(title string, filter string) []string {
	dialogOptions := &options.OpenDialog{
		Title:      title,
		Filters:    filter,
		AllowFiles: true,
	}
	return r.runtime.Dialog.Open(dialogOptions)
}

// OpenDirectoryDialog will call the Runtime.Dialog.OpenDialog method requesting File selection
func (r *RuntimeTest) OpenDirectoryDialog(title string, filter string) []string {
	dialogOptions := &options.OpenDialog{
		Title:            title,
		Filters:          filter,
		AllowDirectories: true,
	}
	return r.runtime.Dialog.Open(dialogOptions)
}

// OpenDialog will call the Runtime.Dialog.OpenDialog method requesting both Files and Directories
func (r *RuntimeTest) OpenDialog(title string, filter string) []string {
	dialogOptions := &options.OpenDialog{
		Title:            title,
		Filters:          filter,
		AllowDirectories: true,
		AllowFiles:       true,
	}
	return r.runtime.Dialog.Open(dialogOptions)
}

// OpenDialogMultiple will call the Runtime.Dialog.OpenDialog method allowing multiple selection
func (r *RuntimeTest) OpenDialogMultiple(title string, filter string) []string {
	dialogOptions := &options.OpenDialog{
		Title:            title,
		Filters:          filter,
		AllowDirectories: true,
		AllowFiles:       true,
		AllowMultiple:    true,
	}
	return r.runtime.Dialog.Open(dialogOptions)
}

// OpenDialogAllOptions will call the Runtime.Dialog.OpenDialog method allowing multiple selection
func (r *RuntimeTest) OpenDialogAllOptions(filter string, defaultDir string, defaultFilename string) []string {
	dialogOptions := &options.OpenDialog{
		DefaultDirectory:           defaultDir,
		DefaultFilename:            defaultFilename,
		Filters:                    filter,
		AllowFiles:                 true,
		AllowDirectories:           true,
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		TreatPackagesAsDirectories: true,
		ResolveAliases:             true,
	}
	return r.runtime.Dialog.Open(dialogOptions)
}

// SaveFileDialog will call the Runtime.Dialog.SaveDialog method requesting a File selection
func (r *RuntimeTest) SaveFileDialog(title string, filter string) string {
	dialogOptions := &options.SaveDialog{
		Title:   title,
		Filters: filter,
	}
	return r.runtime.Dialog.Save(dialogOptions)
}

// SaveDialogAllOptions will call the Runtime.Dialog.SaveDialog method allowing multiple selection
func (r *RuntimeTest) SaveDialogAllOptions(filter string, defaultDir string, defaultFilename string) string {
	dialogOptions := &options.SaveDialog{
		DefaultDirectory:           defaultDir,
		DefaultFilename:            defaultFilename,
		Filters:                    filter,
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		TreatPackagesAsDirectories: true,
	}
	return r.runtime.Dialog.Save(dialogOptions)
}

// HideWindow will call the Runtime.Window.Hide method and then call
// Runtime.Window.Show 3 seconds later.
func (r *RuntimeTest) HideWindow() {
	time.AfterFunc(3*time.Second, func() { r.runtime.Window.Show() })
	r.runtime.Window.Hide()
}

// Maximise the Window
func (r *RuntimeTest) Maximise() {
	r.runtime.Window.Maximise()
}

// Unmaximise the Window
func (r *RuntimeTest) Unmaximise() {
	r.runtime.Window.Unmaximise()
}

// Minimise the Window
func (r *RuntimeTest) Minimise() {
	r.runtime.Window.Minimise()
}

// Unminimise the Window
func (r *RuntimeTest) Unminimise() {
	r.runtime.Window.Unminimise()
}

// Check is system is running in dark mode
func (r *RuntimeTest) IsDarkMode() bool {
	return r.runtime.System.IsDarkMode()
}
