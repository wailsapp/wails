//go:build !darwin || ios || server

package application

func newActivityImpl(app *App) activityImpl {
	return nil
}
