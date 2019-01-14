package wails

// RuntimeWindow exposes an interface for manipulating the window
type RuntimeWindow struct {
	renderer Renderer
}

func newRuntimeWindow(renderer Renderer) *RuntimeWindow {
	return &RuntimeWindow{
		renderer: renderer,
	}
}

// SetColour sets the the window colour
func (r *RuntimeWindow) SetColour(colour string) error {
	return r.renderer.SetColour(colour)
}

// Fullscreen makes the window fullscreen
func (r *RuntimeWindow) Fullscreen() {
	r.renderer.Fullscreen()
}

// UnFullscreen attempts to restore the window to the size/position before fullscreen
func (r *RuntimeWindow) UnFullscreen() {
	r.renderer.UnFullscreen()
}

// SetTitle sets the the window title
func (r *RuntimeWindow) SetTitle(title string) {
	r.renderer.SetTitle(title)
}

// Close shuts down the window and therefore the app
func (r *RuntimeWindow) Close() {
	r.renderer.Close()
}
