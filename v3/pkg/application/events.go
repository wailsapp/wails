package application

import (
	"encoding/json"
	"sync"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type ApplicationEvent struct {
	Id        uint
	ctx       *ApplicationEventContext
	cancelled bool
	lock      sync.RWMutex
}

func (w *ApplicationEvent) Context() *ApplicationEventContext {
	return w.ctx
}

func newApplicationEvent(id events.ApplicationEventType) *ApplicationEvent {
	return &ApplicationEvent{
		Id:  uint(id),
		ctx: newApplicationEventContext(),
	}
}

func (w *ApplicationEvent) Cancel() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.cancelled = true
}

func (w *ApplicationEvent) IsCancelled() bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.cancelled
}

var applicationEvents = make(chan *ApplicationEvent, 5)

type windowEvent struct {
	WindowID uint
	EventID  uint
}

var windowEvents = make(chan *windowEvent, 5)

var menuItemClicked = make(chan uint, 5)

type CustomEvent struct {
	Name      string `json:"name"`
	Data      any    `json:"data"`
	Sender    string `json:"sender"` // Name of the window sending the event, or "" if sent from application
	cancelled bool
	lock      sync.RWMutex
}

func (e *CustomEvent) Cancel() {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.cancelled = true
}

func (e *CustomEvent) IsCancelled() bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.cancelled
}

func (e *CustomEvent) ToJSON() string {
	marshal, err := json.Marshal(&e)
	if err != nil {
		// TODO: Fatal error? log?
		return ""
	}
	return string(marshal)
}

// WailsEventListener is an interface that can be implemented to listen for Wails events
// It is used by the RegisterListener method of the Application.
type WailsEventListener interface {
	DispatchWailsEvent(event *CustomEvent)
}

type hook struct {
	callback func(*CustomEvent)
}

// eventListener holds a callback function which is invoked when
// the event listened for is emitted. It has a counter which indicates
// how the total number of events it is interested in. A value of zero
// means it does not expire (default).
type eventListener struct {
	callback func(*CustomEvent) // Function to call with emitted event data
	counter  int                // The number of times this callback may be called. -1 = infinite
	delete   bool               // Flag to indicate that this listener should be deleted
}

// EventProcessor handles custom events
type EventProcessor struct {
	// Go event listeners
	listeners              map[string][]*eventListener
	notifyLock             sync.RWMutex
	dispatchEventToWindows func(*CustomEvent)
	hooks                  map[string][]*hook
	hookLock               sync.RWMutex
}

func NewWailsEventProcessor(dispatchEventToWindows func(*CustomEvent)) *EventProcessor {
	return &EventProcessor{
		listeners:              make(map[string][]*eventListener),
		dispatchEventToWindows: dispatchEventToWindows,
		hooks:                  make(map[string][]*hook),
	}
}

// On is the equivalent of Javascript's `addEventListener`
func (e *EventProcessor) On(eventName string, callback func(event *CustomEvent)) func() {
	return e.registerListener(eventName, callback, -1)
}

// OnMultiple is the same as `OnApplicationEvent` but will unregister after `count` events
func (e *EventProcessor) OnMultiple(eventName string, callback func(event *CustomEvent), counter int) func() {
	return e.registerListener(eventName, callback, counter)
}

// Once is the same as `OnApplicationEvent` but will unregister after the first event
func (e *EventProcessor) Once(eventName string, callback func(event *CustomEvent)) func() {
	return e.registerListener(eventName, callback, 1)
}

// Emit sends an event to all listeners
func (e *EventProcessor) Emit(thisEvent *CustomEvent) {
	if thisEvent == nil {
		return
	}

	// If we have any hooks, run them first and check if the event was cancelled
	if e.hooks != nil {
		if hooks, ok := e.hooks[thisEvent.Name]; ok {
			for _, thisHook := range hooks {
				thisHook.callback(thisEvent)
				if thisEvent.IsCancelled() {
					return
				}
			}
		}
	}

	go func() {
		defer handlePanic()
		e.dispatchEventToListeners(thisEvent)
	}()
	go func() {
		defer handlePanic()
		e.dispatchEventToWindows(thisEvent)
	}()
}

func (e *EventProcessor) Off(eventName string) {
	e.unRegisterListener(eventName)
}

func (e *EventProcessor) OffAll() {
	e.notifyLock.Lock()
	defer e.notifyLock.Unlock()
	e.listeners = make(map[string][]*eventListener)
}

// registerListener provides a means of subscribing to events of type "eventName"
func (e *EventProcessor) registerListener(eventName string, callback func(*CustomEvent), counter int) func() {
	// Create new eventListener
	thisListener := &eventListener{
		callback: callback,
		counter:  counter,
		delete:   false,
	}
	e.notifyLock.Lock()
	// Append the new listener to the listeners slice
	e.listeners[eventName] = append(e.listeners[eventName], thisListener)
	e.notifyLock.Unlock()
	return func() {
		e.notifyLock.Lock()
		defer e.notifyLock.Unlock()

		if _, ok := e.listeners[eventName]; !ok {
			return
		}
		e.listeners[eventName] = lo.Filter(e.listeners[eventName], func(l *eventListener, i int) bool {
			return l != thisListener
		})
	}
}

// RegisterHook provides a means of registering methods to be called before emitting the event
func (e *EventProcessor) RegisterHook(eventName string, callback func(*CustomEvent)) func() {
	// Create new hook
	thisHook := &hook{
		callback: callback,
	}
	e.hookLock.Lock()
	// Append the new listener to the listeners slice
	e.hooks[eventName] = append(e.hooks[eventName], thisHook)
	e.hookLock.Unlock()
	return func() {
		e.hookLock.Lock()
		defer e.hookLock.Unlock()

		if _, ok := e.hooks[eventName]; !ok {
			return
		}
		e.hooks[eventName] = lo.Filter(e.hooks[eventName], func(l *hook, i int) bool {
			return l != thisHook
		})
	}
}

// unRegisterListener provides a means of unsubscribing to events of type "eventName"
func (e *EventProcessor) unRegisterListener(eventName string) {
	e.notifyLock.Lock()
	defer e.notifyLock.Unlock()
	delete(e.listeners, eventName)
}

// dispatchEventToListeners calls all registered listeners event name
func (e *EventProcessor) dispatchEventToListeners(event *CustomEvent) {

	e.notifyLock.Lock()
	defer e.notifyLock.Unlock()

	listeners := e.listeners[event.Name]
	if listeners == nil {
		return
	}

	// We have a dirty flag to indicate that there are items to delete
	itemsToDelete := false

	// Callback in goroutine
	for _, listener := range listeners {
		if listener.counter > 0 {
			listener.counter--
		}
		go func() {
			if event.IsCancelled() {
				return
			}
			defer handlePanic()
			listener.callback(event)
		}()

		if listener.counter == 0 {
			listener.delete = true
			itemsToDelete = true
		}
	}

	// Do we have items to delete?
	if itemsToDelete == true {
		e.listeners[event.Name] = lo.Filter(listeners, func(l *eventListener, i int) bool {
			return l.delete == false
		})
	}
}
