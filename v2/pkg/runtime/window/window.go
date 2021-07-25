// +build !experimental

package window

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/internal/servicebus"
)

type Window struct{}

// SetTitle sets the title of the window
func (w *Window) SetTitle(ctx context.Context, title string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:settitle", title)
}

// Fullscreen makes the window fullscreen
func (w *Window) Fullscreen(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:fullscreen", "")
}

// UnFullscreen makes the window UnFullscreen
func (w *Window) UnFullscreen(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:unfullscreen", "")
}

// Center the window on the current screen
func (w *Window) Center(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:center", "")
}

// Show shows the window if hidden
func (w *Window) Show(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:show", "")
}

// Hide the window
func (w *Window) Hide(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:hide", "")
}

// SetSize sets the size of the window
func (w *Window) SetSize(ctx context.Context, width int, height int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:size:%d:%d", width, height)
	bus.Publish(message, "")
}

// SetSize sets the size of the window
func (w *Window) SetMinSize(ctx context.Context, width int, height int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:minsize:%d:%d", width, height)
	bus.Publish(message, "")
}

// SetSize sets the size of the window
func (w *Window) SetMaxSize(ctx context.Context, width int, height int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:maxsize:%d:%d", width, height)
	bus.Publish(message, "")
}

// SetPosition sets the position of the window
func (w *Window) SetPosition(ctx context.Context, x int, y int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:position:%d:%d", x, y)
	bus.Publish(message, "")
}

// Maximise the window
func (w *Window) Maximise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:maximise", "")
}

// Unmaximise the window
func (w *Window) Unmaximise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:unmaximise", "")
}

// Minimise the window
func (w *Window) Minimise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:minimise", "")
}

// Unminimise the window
func (w *Window) Unminimise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:unminimise", "")
}
