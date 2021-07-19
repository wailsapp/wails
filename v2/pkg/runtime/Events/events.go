package events

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// On registers a listener for a particular event
func On(ctx context.Context, eventName string, callback func(optionalData ...interface{})) {

	bus := servicebus.ExtractBus(ctx)

	eventMessage := &message.OnEventMessage{
		Name:     eventName,
		Callback: callback,
		Counter:  -1,
	}
	bus.Publish("event:on", eventMessage)
}

// Once registers a listener for a particular event. After the first callback, the
// listener is deleted.
func Once(ctx context.Context, eventName string, callback func(optionalData ...interface{})) {
	bus := servicebus.ExtractBus(ctx)
	eventMessage := &message.OnEventMessage{
		Name:     eventName,
		Callback: callback,
		Counter:  1,
	}
	bus.Publish("event:on", eventMessage)
}

// Emit pass through
func Emit(ctx context.Context, eventName string, optionalData ...interface{}) {
	bus := servicebus.ExtractBus(ctx)
	eventMessage := &message.EventMessage{
		Name: eventName,
		Data: optionalData,
	}

	bus.Publish("event:emit:from:g", eventMessage)
}
