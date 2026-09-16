package application

// SpotlightManager indexes application content for system search.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type SpotlightManager struct {
	app *App
}

func newSpotlightManager(app *App) *SpotlightManager {
	return &SpotlightManager{app: app}
}
