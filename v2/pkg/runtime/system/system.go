// +build !experimental

package system

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Quit the application
func Quit(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	// Start shutdown of Wails
	bus.Publish("quit", "runtime.Quit()")
}
