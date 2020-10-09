// +build debug

package app

// Init initialises the application for a debug environment
func (a *App) Init() error {
	a.debug = true
	return nil
}
