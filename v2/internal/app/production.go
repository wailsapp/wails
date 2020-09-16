// +build !debug

package app

// Init initialises the application for a production environment
func (a *App) Init() error {
	println("Processing production cli options")
	return nil
}
