package runtime

import "github.com/wailsapp/wails/lib/interfaces"

// Window exposes an interface for manipulating the window
type Window struct {
	renderer interfaces.Renderer
}

func newWindow(renderer interfaces.Renderer) *Window {
	return &Window{
		renderer: renderer,
	}
}

// SetColour sets the the window colour
func (r *Window) SetColour(colour string) error {
	return r.renderer.SetColour(colour)
}

// Fullscreen makes the window fullscreen
func (r *Window) Fullscreen() {
	r.renderer.Fullscreen()
}

// UnFullscreen attempts to restore the window to the size/position before fullscreen
func (r *Window) UnFullscreen() {
	r.renderer.UnFullscreen()
}

// SetTitle sets the the window title
func (r *Window) SetTitle(title string) {
	r.renderer.SetTitle(title)
}

// Close shuts down the window and therefore the app
func (r *Window) Close() {
	r.renderer.Close()
}
