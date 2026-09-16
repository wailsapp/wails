//go:build !darwin || ios || server

package application

func newServicesProviderImpl(app *App) servicesProviderImpl {
	return nil
}
