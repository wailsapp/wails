package application

import (
	"github.com/wailsapp/wails/v3/pkg/events"
)

// EventManager manages event-related operations
type EventManager struct {
	app *App
}

// NewEventManager creates a new EventManager instance
func NewEventManager(app *App) *EventManager {
	return &EventManager{
		app: app,
	}
}

// Emit emits a custom event
func (em *EventManager) Emit(name string, data ...any) {
	em.app.customEventProcessor.Emit(&CustomEvent{
		Name: name,
		Data: data,
	})
}

// On registers a listener for custom events
func (em *EventManager) On(name string, callback func(event *CustomEvent)) func() {
	return em.app.customEventProcessor.On(name, callback)
}

// Off removes all listeners for a custom event
func (em *EventManager) Off(name string) {
	em.app.customEventProcessor.Off(name)
}

// OnMultiple registers a listener for custom events that will be called N times
func (em *EventManager) OnMultiple(name string, callback func(event *CustomEvent), counter int) {
	em.app.customEventProcessor.OnMultiple(name, callback, counter)
}

// Reset removes all custom event listeners
func (em *EventManager) Reset() {
	em.app.customEventProcessor.OffAll()
}

// OnApplicationEvent registers a listener for application events
func (em *EventManager) OnApplicationEvent(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func() {
	return em.app.OnApplicationEvent(eventType, callback)
}

// RegisterApplicationEventHook registers an application event hook
func (em *EventManager) RegisterApplicationEventHook(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func() {
	return em.app.RegisterApplicationEventHook(eventType, callback)
}

// Dispatch dispatches an event to listeners (internal use)
func (em *EventManager) dispatch(event *CustomEvent) {
	em.app.dispatchEventToListeners(event)
}

// HandleApplicationEvent handles application events (internal use)
func (em *EventManager) handleApplicationEvent(event *ApplicationEvent) {
	em.app.handleApplicationEvent(event)
}
