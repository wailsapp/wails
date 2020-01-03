package interfaces

import (
	"github.com/wailsapp/wails/lib/messages"
)

// Renderer is an interface describing a Wails target to render the app to
type Renderer interface {
	Initialise(AppConfig, IPCManager, EventManager) error
	Run() error
	EnableConsole()

	// Binding
	NewBinding(bindingName string) error

	// Events
	NotifyEvent(eventData *messages.EventData) error

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
