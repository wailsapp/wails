package application

// PermissionsManager requests and reports system permissions such as camera, microphone, screen recording and location.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type PermissionsManager struct {
	app *App
}

func newPermissionsManager(app *App) *PermissionsManager {
	return &PermissionsManager{app: app}
}
