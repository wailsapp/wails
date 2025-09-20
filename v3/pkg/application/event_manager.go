package application

import (
	"fmt"
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
	println("🟢 [newEventManager] Creating EventManager")
	em := &EventManager{
		app: app,
	}
	println("🟢 [newEventManager] EventManager created")
	return em
}

// Emit emits a custom event
func (em *EventManager) Emit(name string, data ...any) {
	// Use fmt.Printf to ensure output on iOS
	fmt.Printf("🔵 [EventManager.Emit] Emitting event: %s, Data count: %d\n", name, len(data))
	if em.app == nil {
		fmt.Println("🔴 [EventManager.Emit] App is nil!")
		return
	}
	if em.app.customEventProcessor == nil {
		fmt.Println("🔴 [EventManager.Emit] customEventProcessor is nil!")
		return
	}
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
	println("🔵 [EventManager.On] Registering listener for:", name)
	if em.app == nil {
		println("🔴 [EventManager.On] App is nil!")
		return func() {}
	}
	if em.app.customEventProcessor == nil {
		println("🔴 [EventManager.On] customEventProcessor is nil!")
		return func() {}
	}
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
	println("🔵 [EventManager.dispatch] Dispatching event to windows:", event.Name)

	// Snapshot windows under RLock
	em.app.windowsLock.RLock()
	windowCount := len(em.app.windows)
	println("🔵 [EventManager.dispatch] Window count:", windowCount)
	for i, window := range em.app.windows {
		if event.IsCancelled() {
			println("🟠 [EventManager.dispatch] Event cancelled, stopping dispatch")
			em.app.windowsLock.RUnlock()
			return
		}
		println("🔵 [EventManager.dispatch] Dispatching to window", i+1)
		window.DispatchWailsEvent(event)
	}
	em.app.windowsLock.RUnlock()

	// Snapshot listeners under Lock
	em.app.wailsEventListenerLock.Lock()
	listeners := slices.Clone(em.app.wailsEventListeners)
	em.app.wailsEventListenerLock.Unlock()

	println("🔵 [EventManager.dispatch] Wails event listener count:", len(listeners))

	for i, listener := range listeners {
		if event.IsCancelled() {
			println("🟠 [EventManager.dispatch] Event cancelled during listener dispatch")
			return
		}
		println("🔵 [EventManager.dispatch] Dispatching to Wails listener", i+1)
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
