package wails

type RuntimeWindow struct {
	renderer Renderer
}

func newRuntimeWindow(renderer Renderer) *RuntimeWindow {
	return &RuntimeWindow{
		renderer: renderer,
	}
}

func (r *RuntimeWindow) SetColour(colour string) error {
	return r.renderer.SetColour(colour)
}

func (r *RuntimeWindow) Fullscreen() {
	r.renderer.Fullscreen()
}

func (r *RuntimeWindow) UnFullscreen() {
	r.renderer.UnFullscreen()
}

func (r *RuntimeWindow) SetTitle(title string) {
	r.renderer.SetTitle(title)
}

func (r *RuntimeWindow) Close() {
	// TODO: Add shutdown mechanism
	r.renderer.Close()
}
