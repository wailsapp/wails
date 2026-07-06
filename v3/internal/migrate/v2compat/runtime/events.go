package runtime

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// eventArgs converts a v3 CustomEvent data payload into the variadic argument
// list expected by v2 event callbacks: nil becomes no arguments, a slice
// emitted from multiple arguments is spread, and any other value is passed
// as a single argument.
func eventArgs(data any) []interface{} {
	switch d := data.(type) {
	case nil:
		return nil
	case []interface{}:
		return d
	default:
		return []interface{}{d}
	}
}

// EventsOn mirrors the v2 runtime.EventsOn function.
// v3 equivalent: app.Event.On.
func EventsOn(_ context.Context, eventName string, callback func(optionalData ...interface{})) func() {
	a := app()
	if a == nil {
		return func() {}
	}
	return a.Event.On(eventName, func(event *application.CustomEvent) {
		callback(eventArgs(event.Data)...)
	})
}

// EventsOff mirrors the v2 runtime.EventsOff function.
// v3 equivalent: app.Event.Off.
func EventsOff(_ context.Context, eventName string, additionalEventNames ...string) {
	a := app()
	if a == nil {
		return
	}
	a.Event.Off(eventName)
	for _, name := range additionalEventNames {
		a.Event.Off(name)
	}
}

// EventsOnce mirrors the v2 runtime.EventsOnce function.
// v3 equivalent: app.Event.On combined with unregistering after the first event.
func EventsOnce(ctx context.Context, eventName string, callback func(optionalData ...interface{})) func() {
	return EventsOnMultiple(ctx, eventName, callback, 1)
}

// EventsOnMultiple mirrors the v2 runtime.EventsOnMultiple function. The
// callback is invoked at most counter times; a counter <= 0 means unlimited.
// v3 equivalent: app.Event.On combined with unregistering after counter events.
func EventsOnMultiple(_ context.Context, eventName string, callback func(optionalData ...interface{}), counter int) func() {
	a := app()
	if a == nil {
		return func() {}
	}
	if counter <= 0 {
		return a.Event.On(eventName, func(event *application.CustomEvent) {
			callback(eventArgs(event.Data)...)
		})
	}

	var remaining atomic.Int64
	remaining.Store(int64(counter))

	var lock sync.Mutex
	var unregister func()
	unregistered := false
	cancel := func() {
		lock.Lock()
		defer lock.Unlock()
		if unregistered {
			return
		}
		unregistered = true
		if unregister != nil {
			unregister()
		}
	}

	off := a.Event.On(eventName, func(event *application.CustomEvent) {
		if remaining.Add(-1) < 0 {
			return
		}
		callback(eventArgs(event.Data)...)
		if remaining.Load() <= 0 {
			cancel()
		}
	})

	lock.Lock()
	unregister = off
	if unregistered {
		// The counter was exhausted before registration completed.
		off()
	}
	lock.Unlock()

	return cancel
}

// EventsEmit mirrors the v2 runtime.EventsEmit function.
// v3 equivalent: app.Event.Emit.
func EventsEmit(_ context.Context, eventName string, optionalData ...interface{}) {
	a := app()
	if a == nil {
		return
	}
	a.Event.Emit(eventName, optionalData...)
}
