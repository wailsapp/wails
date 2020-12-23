package runtime

import (
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Events defines all events related operations
type Events interface {
	On(eventName string, callback func(optionalData ...interface{}))
	Once(eventName string, callback func(optionalData ...interface{}))
	OnMultiple(eventName string, callback func(optionalData ...interface{}), maxCallbacks int)
	Emit(eventName string, optionalData ...interface{})
	OnThemeChange(callback func(darkMode bool))
}

// event exposes the events interface
type event struct {
	bus *servicebus.ServiceBus
}

// newEvents creates a new Events struct
func newEvents(bus *servicebus.ServiceBus) Events {
	return &event{
		bus: bus,
	}
}

// On registers a listener for a particular event
func (r *event) On(eventName string, callback func(optionalData ...interface{})) {
	eventMessage := &message.OnEventMessage{
		Name:     eventName,
		Callback: callback,
		Counter:  -1,
	}
	r.bus.Publish("event:on", eventMessage)
}

// Once registers a listener for a particular event. After the first callback, the
// listener is deleted.
func (r *event) Once(eventName string, callback func(optionalData ...interface{})) {
	eventMessage := &message.OnEventMessage{
		Name:     eventName,
		Callback: callback,
		Counter:  1,
	}
	r.bus.Publish("event:on", eventMessage)
}

// OnMultiple registers a listener for a particular event, for a given maximum amount of callbacks.
// Once the callback has been run `maxCallbacks` times, the listener is deleted.
func (r *event) OnMultiple(eventName string, callback func(optionalData ...interface{}), maxCallbacks int) {
	eventMessage := &message.OnEventMessage{
		Name:     eventName,
		Callback: callback,
		Counter:  maxCallbacks,
	}
	r.bus.Publish("event:on", eventMessage)
}

// Emit pass through
func (r *event) Emit(eventName string, optionalData ...interface{}) {
	eventMessage := &message.EventMessage{
		Name: eventName,
		Data: optionalData,
	}

	r.bus.Publish("event:emit:from:g", eventMessage)
}

// OnThemeChange allows you to register callbacks when the system theme changes
// from light or dark.
func (r *event) OnThemeChange(callback func(darkMode bool)) {
	r.On("wails:system:themechange", func(data ...interface{}) {
		if len(data) != 1 {
			// TODO: Log error
			return
		}
		darkMode, ok := data[0].(bool)
		if !ok {
			// TODO: Log error
			return
		}
		callback(darkMode)
	})
}
