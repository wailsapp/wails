package runtime

import "github.com/wailsapp/wails/lib/interfaces"

// Dialog exposes an interface to native dialogs
type Dialog struct {
	renderer interfaces.Renderer
}

// NewDialog creates a new Dialog struct
func NewDialog(renderer interfaces.Renderer) *Dialog {
	return &Dialog{
		renderer: renderer,
	}
}

// SelectFile prompts the user to select a file
func (r *Dialog) SelectFile(params ...string) string {
	filter := ""
	if len(params) > 0 {
		filter = params[0]
	}
	return r.renderer.SelectFile(filter)
}

// SelectDirectory prompts the user to select a directory
func (r *Dialog) SelectDirectory() string {
	return r.renderer.SelectDirectory()
}

// SelectSaveFile prompts the user to select a file for saving
func (r *Dialog) SelectSaveFile(params ...string) string {
	filter := ""
	if len(params) > 0 {
		filter = params[0]
	}
	return r.renderer.SelectSaveFile(filter)
}
