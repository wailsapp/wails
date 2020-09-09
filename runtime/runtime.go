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
	Store      *StoreProvider
}

// NewRuntime creates a new Runtime struct
func NewRuntime(eventManager interfaces.EventManager, renderer interfaces.Renderer) *Runtime {
	result := &Runtime{
		Events:     NewEvents(eventManager),
		Log:        NewLog(),
		Dialog:     NewDialog(renderer),
		Window:     NewWindow(renderer),
		Browser:    NewBrowser(),
		FileSystem: NewFileSystem(),
	}
	// We need a reference to itself
	result.Store = NewStoreProvider(result)
	return result
}
