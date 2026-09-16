package application

// ServicesProviderManager registers handlers for the macOS Services menu so other applications can send this one text or files.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type ServicesProviderManager struct {
	app *App
}

func newServicesProviderManager(app *App) *ServicesProviderManager {
	return &ServicesProviderManager{app: app}
}
