package wails

type RuntimeDialog struct {
	renderer Renderer
}

func newRuntimeDialog(renderer Renderer) *RuntimeDialog {
	return &RuntimeDialog{
		renderer: renderer,
	}
}

func (r *RuntimeDialog) SelectFile() string {
	return r.renderer.SelectFile()
}

func (r *RuntimeDialog) SelectDirectory() string {
	return r.renderer.SelectDirectory()
}

func (r *RuntimeDialog) SelectSaveFile() string {
	return r.renderer.SelectSaveFile()
}
