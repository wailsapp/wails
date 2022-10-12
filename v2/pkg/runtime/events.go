package runtime

import (
	"context"
)

// EventsOn registers a listener for the given event name
func EventsOn(ctx context.Context, eventName string, callback func(optionalData ...interface{})) {
	events := getEvents(ctx)
	events.On(eventName, callback)
}

// EventsOff unregisters a listener for the given event name, optionally multiple listeneres can be unregistered via `additionalEventNames`
func EventsOff(ctx context.Context, eventName string, additionalEventNames ...string) {
	events := getEvents(ctx)
	events.Off(eventName)

	if len(additionalEventNames) > 0 {
		for _, eventName := range additionalEventNames {
			events.Off(eventName)
		}
	}
}

// EventsOnce registers a listener for the given event name. After the first callback, the
// listener is deleted.
func EventsOnce(ctx context.Context, eventName string, callback func(optionalData ...interface{})) {
	events := getEvents(ctx)
	events.Once(eventName, callback)
}

// EventsOnMultiple registers a listener for the given event name, that may be called a maximum of 'counter' times
func EventsOnMultiple(ctx context.Context, eventName string, callback func(optionalData ...interface{}), counter int) {
	events := getEvents(ctx)
	events.OnMultiple(eventName, callback, counter)
}

// EventsEmit pass through
func EventsEmit(ctx context.Context, eventName string, optionalData ...interface{}) {
	events := getEvents(ctx)
	events.Emit(eventName, optionalData...)
}
