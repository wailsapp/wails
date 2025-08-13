package runtime

import (
	"context"
)

// EventsOn registers a listener for the given event name. It returns a function to cancel the listener
func EventsOn(ctx context.Context, eventName string, callback func(optionalData ...interface{})) func() {
	events := getEvents(ctx)
	return events.On(eventName, callback)
}

// EventsOff unregisters a listener for the given event name, optionally multiple listeners can be unregistered via `additionalEventNames`
func EventsOff(ctx context.Context, eventName string, additionalEventNames ...string) {
	events := getEvents(ctx)
	events.Off(eventName)

	if len(additionalEventNames) > 0 {
		for _, eventName := range additionalEventNames {
			events.Off(eventName)
		}
	}
}

// EventsOff unregisters a listener for the given event name, optionally multiple listeners can be unregistered via `additionalEventNames`
func EventsOffAll(ctx context.Context) {
	events := getEvents(ctx)
	events.OffAll()
}

// EventsOnce registers a listener for the given event name. After the first callback, the
// listener is deleted. It returns a function to cancel the listener
func EventsOnce(ctx context.Context, eventName string, callback func(optionalData ...interface{})) func() {
	events := getEvents(ctx)
	return events.Once(eventName, callback)
}

// EventsOnMultiple registers a listener for the given event name, that may be called a maximum of 'counter' times. It returns a function
// to cancel the listener
func EventsOnMultiple(ctx context.Context, eventName string, callback func(optionalData ...interface{}), counter int) func() {
	events := getEvents(ctx)
	return events.OnMultiple(eventName, callback, counter)
}

// EventsEmit pass through
func EventsEmit(ctx context.Context, eventName string, optionalData ...interface{}) {
	events := getEvents(ctx)
	events.Emit(eventName, optionalData...)
}
