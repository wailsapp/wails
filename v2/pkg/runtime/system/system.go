// +build !experimental

package system

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

type System struct{}

// Quit the application
func (s *System) Quit(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	// Start shutdown of Wails
	bus.Publish("quit", "runtime.Quit()")
}
