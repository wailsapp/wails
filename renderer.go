package wails

// Renderer is an interface describing a Wails target to render the app to
type Renderer interface {
	Initialise(*AppConfig, *ipcManager, *eventManager) error
	Run() error

	// Binding
	NewBinding(bindingName string) error
	Callback(data string) error

	// Events
	NotifyEvent(eventData *eventData) error

	// Dialog Runtime
	SelectFile() string
	SelectDirectory() string
	SelectSaveFile() string

	// Window Runtime
	SetColour(string) error
	Fullscreen()
	UnFullscreen()
	SetTitle(title string)
	Close()
}
