//go:build !darwin || ios || server

package application

// newLifecycleImpl returns nil where the platform has no termination
// control; LifecycleManager then hands out no-op releases.
func newLifecycleImpl(app *App) lifecycleImpl {
	return nil
}
