package application

// PowerManager prevents sleep and reports power and thermal state.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type PowerManager struct {
	app *App
}

func newPowerManager(app *App) *PowerManager {
	return &PowerManager{app: app}
}
