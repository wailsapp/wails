package interfaces

import (
	"github.com/wailsapp/wails/lib/messages"
)

// Renderer is an interface describing a Wails target to render the app to
type Renderer interface {
	Initialise(AppConfig, IPCManager, EventManager) error
	Run() error

	// Binding
	NewBinding(bindingName string) error

	// Events
	NotifyEvent(eventData *messages.EventData) error

	// Dialog Runtime
	SelectFile(title string, filter string) string
	SelectDirectory() string
	SelectSaveFile(title string, filter string) string

	// Window Runtime
	SetColour(string) error

	SetMinSize(width, height int)
	SetMaxSize(width, height int)

	Fullscreen()
	UnFullscreen()
	SetTitle(title string)
	Close()
}
