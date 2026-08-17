//go:build wails_native

package application

import "fmt"

type frontendState struct{}

func configureFrontend(*App, Options) {}

func startFrontendEventLoops(*App) {}

// Native applications may still use ServiceStartup and ServiceShutdown for
// lifecycle management. JavaScript bindings are deliberately absent.
func bindFrontendService(*App, Service) error { return nil }

func attachServiceRoute(_ *App, service Service) error {
	if service.options.Route != "" {
		return fmt.Errorf("service %q declares HTTP route %q, which is unavailable in a wails_native build", getServiceName(service), service.options.Route)
	}
	return nil
}
