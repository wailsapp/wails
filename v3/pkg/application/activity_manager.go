package application

// ActivityManager publishes and continues user activities for Handoff and universal links.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type ActivityManager struct {
	app *App
}

func newActivityManager(app *App) *ActivityManager {
	return &ActivityManager{app: app}
}
