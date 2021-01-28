package runtime

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Window defines all Window related operations
type Window interface {
	Close()
	Center()
	Show()
	Hide()
	Maximise()
	Unmaximise()
	Minimise()
	Unminimise()
	SetTitle(title string)
	SetSize(width int, height int)
	SetMinSize(width int, height int)
	SetMaxSize(width int, height int)
	SetPosition(x int, y int)
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

// Center the window on the current screen
func (w *window) Center() {
	w.bus.Publish("window:center", "")
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
	message := fmt.Sprintf("window:size:%d:%d", width, height)
	w.bus.Publish(message, "")
}

// SetSize sets the size of the window
func (w *window) SetMinSize(width int, height int) {
	message := fmt.Sprintf("window:minsize:%d:%d", width, height)
	w.bus.Publish(message, "")
}

// SetSize sets the size of the window
func (w *window) SetMaxSize(width int, height int) {
	message := fmt.Sprintf("window:maxsize:%d:%d", width, height)
	w.bus.Publish(message, "")
}

// SetPosition sets the position of the window
func (w *window) SetPosition(x int, y int) {
	message := fmt.Sprintf("window:position:%d:%d", x, y)
	w.bus.Publish(message, "")
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
