package application

// LifecycleManager controls process termination behaviour such as sudden and automatic termination.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type LifecycleManager struct {
	app *App
}

func newLifecycleManager(app *App) *LifecycleManager {
	return &LifecycleManager{app: app}
}
