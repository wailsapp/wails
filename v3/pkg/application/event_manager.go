package application

import (
	"slices"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// EventManager manages event-related operations
type EventManager struct {
	app *App
}

// newEventManager creates a new EventManager instance
func newEventManager(app *App) *EventManager {
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

// EmitEvent emits a custom event object (internal use)
func (em *EventManager) EmitEvent(event *CustomEvent) {
	em.app.customEventProcessor.Emit(event)
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
	eventID := uint(eventType)
	em.app.applicationEventListenersLock.Lock()
	defer em.app.applicationEventListenersLock.Unlock()
	listener := &EventListener{
		callback: callback,
	}
	em.app.applicationEventListeners[eventID] = append(em.app.applicationEventListeners[eventID], listener)
	if em.app.impl != nil {
		go func() {
			defer handlePanic()
			em.app.impl.on(eventID)
		}()
	}

	return func() {
		// lock the map
		em.app.applicationEventListenersLock.Lock()
		defer em.app.applicationEventListenersLock.Unlock()
		// Remove listener
		em.app.applicationEventListeners[eventID] = lo.Without(em.app.applicationEventListeners[eventID], listener)
	}
}

// RegisterApplicationEventHook registers an application event hook
func (em *EventManager) RegisterApplicationEventHook(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func() {
	eventID := uint(eventType)
	em.app.applicationEventHooksLock.Lock()
	defer em.app.applicationEventHooksLock.Unlock()
	thisHook := &eventHook{
		callback: callback,
	}
	em.app.applicationEventHooks[eventID] = append(em.app.applicationEventHooks[eventID], thisHook)

	return func() {
		em.app.applicationEventHooksLock.Lock()
		em.app.applicationEventHooks[eventID] = lo.Without(em.app.applicationEventHooks[eventID], thisHook)
		em.app.applicationEventHooksLock.Unlock()
	}
}

// Dispatch dispatches an event to listeners (internal use)
func (em *EventManager) dispatch(event *CustomEvent) {
	// Snapshot windows under RLock
	em.app.windowsLock.RLock()
	for _, window := range em.app.windows {
		if event.IsCancelled() {
			em.app.windowsLock.RUnlock()
			return
		}
		window.DispatchWailsEvent(event)
	}
	em.app.windowsLock.RUnlock()

	// Snapshot listeners under Lock
	em.app.wailsEventListenerLock.Lock()
	listeners := slices.Clone(em.app.wailsEventListeners)
	em.app.wailsEventListenerLock.Unlock()

	for _, listener := range listeners {
		if event.IsCancelled() {
			return
		}
		listener.DispatchWailsEvent(event)
	}
}

// HandleApplicationEvent handles application events (internal use)
func (em *EventManager) handleApplicationEvent(event *ApplicationEvent) {
	defer handlePanic()
	em.app.applicationEventListenersLock.RLock()
	listeners, ok := em.app.applicationEventListeners[event.Id]
	em.app.applicationEventListenersLock.RUnlock()
	if !ok {
		return
	}

	// Process Hooks
	em.app.applicationEventHooksLock.RLock()
	hooks, ok := em.app.applicationEventHooks[event.Id]
	em.app.applicationEventHooksLock.RUnlock()
	if ok {
		for _, thisHook := range hooks {
			thisHook.callback(event)
			if event.IsCancelled() {
				return
			}
		}
	}

	for _, listener := range listeners {
		go func() {
			if event.IsCancelled() {
				return
			}
			defer handlePanic()
			listener.callback(event)
		}()
	}
}
