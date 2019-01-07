package wails

type Runtime struct {
	Events *RuntimeEvents
	Log    *RuntimeLog
	Dialog *RuntimeDialog
	Window *RuntimeWindow
}

func newRuntime(eventManager *eventManager, renderer Renderer) *Runtime {
	return &Runtime{
		Events: newRuntimeEvents(eventManager),
		Log:    newRuntimeLog(),
		Dialog: newRuntimeDialog(renderer),
		Window: newRuntimeWindow(renderer),
	}
}
