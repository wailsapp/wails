// +build debug

package app

// Init initialises the application for a debug environment
func (a *App) Init() error {
	// Indicate debug mode
	a.debug = true
	// Enable dev tools
	a.options.DevTools = true
	return nil
}
