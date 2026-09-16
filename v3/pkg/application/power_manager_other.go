//go:build !darwin || ios || server

package application

// newPowerImpl returns nil where no native power API is wired up;
// PowerManager then hands out no-op releases and a zero PowerState.
func newPowerImpl(app *App) powerImpl {
	return nil
}
