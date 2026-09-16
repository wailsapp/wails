package application

// HapticsManager performs system haptic feedback where the platform supports it.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type HapticsManager struct {
	app *App
}

func newHapticsManager(app *App) *HapticsManager {
	return &HapticsManager{app: app}
}
