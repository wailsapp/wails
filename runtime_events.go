package wails

// RuntimeEvents exposes the events interface
type RuntimeEvents struct {
	eventManager *eventManager
}

func newRuntimeEvents(eventManager *eventManager) *RuntimeEvents {
	return &RuntimeEvents{
		eventManager: eventManager,
	}
}

// On pass through
func (r *RuntimeEvents) On(eventName string, callback func(optionalData ...interface{})) {
	r.eventManager.On(eventName, callback)
}

// Emit pass through
func (r *RuntimeEvents) Emit(eventName string, optionalData ...interface{}) {
	r.eventManager.Emit(eventName, optionalData)
}
