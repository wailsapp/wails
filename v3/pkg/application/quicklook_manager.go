package application

// QuickLookManager previews files and generates thumbnails with Quick Look.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type QuickLookManager struct {
	app *App
}

func newQuickLookManager(app *App) *QuickLookManager {
	return &QuickLookManager{app: app}
}
