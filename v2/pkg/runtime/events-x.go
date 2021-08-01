// +build experimental

package runtime

import (
	"context"
)

// EventsOn registers a listener for the given event name
func EventsOn(ctx context.Context, eventName string, callback func(optionalData ...interface{})) {
	events := getEvents(ctx)
	events.On(eventName, callback)
}

// EventsOff unregisters a listener for the given event name
func EventsOff(ctx context.Context, eventName string) {
	events := getEvents(ctx)
	events.Off(eventName)
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
