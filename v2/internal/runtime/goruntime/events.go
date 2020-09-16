package goruntime

import (
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Events defines all events related operations
type Events interface {
	On(eventName string, callback func(optionalData ...interface{}))
	Emit(eventName string, optionalData ...interface{})
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

// On pass through
func (r *event) On(eventName string, callback func(optionalData ...interface{})) {
	eventMessage := &message.OnEventMessage{
		Name:     eventName,
		Callback: callback,
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
