package wails

// Runtime is the Wails Runtime Interface, given to a user who has defined the WailsInit method
type Runtime struct {
	Events     *RuntimeEvents
	Log        *RuntimeLog
	Dialog     *RuntimeDialog
	Window     *RuntimeWindow
	FileSystem *RuntimeFileSystem
}

func newRuntime(eventManager *eventManager, renderer Renderer) *Runtime {
	return &Runtime{
		Events:     newRuntimeEvents(eventManager),
		Log:        newRuntimeLog(),
		Dialog:     newRuntimeDialog(renderer),
		Window:     newRuntimeWindow(renderer),
		FileSystem: newRuntimeFileSystem(),
	}
}
