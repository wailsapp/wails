// +build !experimental

package runtime

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// WindowSetTitle sets the title of the window
func WindowSetTitle(ctx context.Context, title string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:settitle", title)
}

// WindowFullscreen makes the window fullscreen
func WindowFullscreen(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:fullscreen", "")
}

// WindowUnFullscreen makes the window UnFullscreen
func WindowUnFullscreen(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:unfullscreen", "")
}

// WindowCenter the window on the current screen
func WindowCenter(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:center", "")
}

// WindowShow shows the window if hidden
func WindowShow(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:show", "")
}

// WindowHide the window
func WindowHide(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:hide", "")
}

// WindowSetSize sets the size of the window
func WindowSetSize(ctx context.Context, width int, height int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:size:%d:%d", width, height)
	bus.Publish(message, "")
}

// WindowSetSize sets the size of the window
func WindowSetMinSize(ctx context.Context, width int, height int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:minsize:%d:%d", width, height)
	bus.Publish(message, "")
}

// WindowSetSize sets the size of the window
func WindowSetMaxSize(ctx context.Context, width int, height int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:maxsize:%d:%d", width, height)
	bus.Publish(message, "")
}

// WindowSetPosition sets the position of the window
func WindowSetPosition(ctx context.Context, x int, y int) {
	bus := servicebus.ExtractBus(ctx)
	message := fmt.Sprintf("window:position:%d:%d", x, y)
	bus.Publish(message, "")
}

// WindowMaximise the window
func WindowMaximise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:maximise", "")
}

// WindowUnmaximise the window
func WindowUnmaximise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:unmaximise", "")
}

// WindowMinimise the window
func WindowMinimise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:minimise", "")
}

// WindowUnminimise the window
func WindowUnminimise(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("window:unminimise", "")
}
