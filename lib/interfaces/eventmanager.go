package interfaces

import "github.com/wailsapp/wails/lib/messages"

// EventManager is the event manager interface
type EventManager interface {
	PushEvent(*messages.EventData)
	Emit(eventName string, optionalData ...interface{})
	OnMultiple(eventName string, callback func(...interface{}), counter uint)
	Once(eventName string, callback func(...interface{}))
	On(eventName string, callback func(...interface{}))
	Start(Renderer)
	Shutdown()
}
