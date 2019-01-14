package wails

// RuntimeDialog exposes an interface to native dialogs
type RuntimeDialog struct {
	renderer Renderer
}

// newRuntimeDialog creates a new RuntimeDialog struct
func newRuntimeDialog(renderer Renderer) *RuntimeDialog {
	return &RuntimeDialog{
		renderer: renderer,
	}
}

// SelectFile prompts the user to select a file
func (r *RuntimeDialog) SelectFile() string {
	return r.renderer.SelectFile()
}

// SelectDirectory prompts the user to select a directory
func (r *RuntimeDialog) SelectDirectory() string {
	return r.renderer.SelectDirectory()
}

// SelectSaveFile prompts the user to select a file for saving
func (r *RuntimeDialog) SelectSaveFile() string {
	return r.renderer.SelectSaveFile()
}
