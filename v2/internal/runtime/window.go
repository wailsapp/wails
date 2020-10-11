package runtime

import (
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Window defines all Window related operations
type Window interface {
	Close()
	Show()
	Hide()
	Maximise()
	Unmaximise()
	Minimise()
	Unminimise()
	SetTitle(title string)
	Fullscreen()
	UnFullscreen()
	SetColour(colour int)
}

// Window exposes the Windows interface
type window struct {
	bus *servicebus.ServiceBus
}

// newWindow creates a new window struct
func newWindow(bus *servicebus.ServiceBus) Window {
	return &window{
		bus: bus,
	}
}

// Close the Window
// DISCUSSION:
//   Should we even be doing this now we have a server build?
//   Runtime.Quit() makes more sense than closing a window...
func (w *window) Close() {
	w.bus.Publish("quit", "runtime.Close()")
}

// SetTitle sets the title of the window
func (w *window) SetTitle(title string) {
	w.bus.Publish("window:settitle", title)
}

// Fullscreen makes the window fullscreen
func (w *window) Fullscreen() {
	w.bus.Publish("window:fullscreen", "")
}

// UnFullscreen makes the window UnFullscreen
func (w *window) UnFullscreen() {
	w.bus.Publish("window:unfullscreen", "")
}

// SetColour sets the window colour to the given int
func (w *window) SetColour(colour int) {
	w.bus.Publish("window:setcolour", colour)
}

// Show shows the window if hidden
func (w *window) Show() {
	w.bus.Publish("window:show", "")
}

// Hide the window
func (w *window) Hide() {
	w.bus.Publish("window:hide", "")
}

// SetSize sets the size of the window
func (w *window) SetSize(width int, height int) {
	size := []int{width, height}
	w.bus.Publish("window:setsize", size)
}

// SetPosition sets the position of the window
func (w *window) SetPosition(x int, y int) {
	position := []int{x, y}
	w.bus.Publish("window:position", position)
}

// Maximise the window
func (w *window) Maximise() {
	w.bus.Publish("window:maximise", "")
}

// Unmaximise the window
func (w *window) Unmaximise() {
	w.bus.Publish("window:unmaximise", "")
}

// Minimise the window
func (w *window) Minimise() {
	w.bus.Publish("window:minimise", "")
}

// Unminimise the window
func (w *window) Unminimise() {
	w.bus.Publish("window:unminimise", "")
}
