package runtime

import "github.com/wailsapp/wails/lib/interfaces"

// Runtime is the Wails Runtime Interface, given to a user who has defined the WailsInit method
type Runtime struct {
	Events     *Events
	Log        *Log
	Dialog     *Dialog
	Window     *Window
	Browser    *Browser
	FileSystem *FileSystem
}

// NewRuntime creates a new Runtime struct
func NewRuntime(eventManager interfaces.EventManager, renderer interfaces.Renderer) *Runtime {
	return &Runtime{
		Events:     newEvents(eventManager),
		Log:        newLog(),
		Dialog:     newDialog(renderer),
		Window:     newWindow(renderer),
		Browser:    NewBrowser(),
		FileSystem: newFileSystem(),
	}
}
