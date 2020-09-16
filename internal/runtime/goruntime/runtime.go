package goruntime

import "github.com/wailsapp/wails/v2/internal/servicebus"

// Runtime is a means for the user to interact with the application at runtime
type Runtime struct {
	Browser Browser
	Events  Events
	Window  Window
	Dialog  Dialog
	bus     *servicebus.ServiceBus
}

// New creates a new runtime
func New(serviceBus *servicebus.ServiceBus) *Runtime {
	return &Runtime{
		Browser: newBrowser(),
		Events:  newEvents(serviceBus),
		Window:  newWindow(serviceBus),
		Dialog:  newDialog(serviceBus),
		bus:     serviceBus,
	}
}

// Quit the application
func (r *Runtime) Quit() {
	r.bus.Publish("quit", "runtime.Quit()")
}
