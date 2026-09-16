package application

// AppleEventsManager registers handlers for Apple Events so scripts and other
// applications can drive this one on macOS.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type AppleEventsManager struct {
	app *App
}

func newAppleEventsManager(app *App) *AppleEventsManager {
	return &AppleEventsManager{app: app}
}
