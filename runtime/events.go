package runtime

import "github.com/wailsapp/wails/lib/interfaces"

// Events exposes the events interface
type Events struct {
	eventManager interfaces.EventManager
}

// NewEvents creates a new Events struct
func NewEvents(eventManager interfaces.EventManager) *Events {
	return &Events{
		eventManager: eventManager,
	}
}

// On pass through
func (r *Events) On(eventName string, callback func(optionalData ...interface{})) {
	r.eventManager.On(eventName, callback)
}

// Emit pass through
func (r *Events) Emit(eventName string, optionalData ...interface{}) {
	r.eventManager.Emit(eventName, optionalData...)
}
