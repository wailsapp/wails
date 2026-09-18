package application

import (
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func newEventManagerTestApp() *App {
	return &App{
		applicationEventListeners: make(map[uint][]*EventListener),
		applicationEventHooks:     make(map[uint][]*eventHook),
	}
}

func TestHandleApplicationEventRunsHooksWithoutListeners(t *testing.T) {
	em := newEventManager(newEventManagerTestApp())

	fired := false
	em.RegisterApplicationEventHook(events.Common.ApplicationStarted, func(event *ApplicationEvent) {
		fired = true
	})

	em.handleApplicationEvent(newApplicationEvent(events.Common.ApplicationStarted))

	if !fired {
		t.Fatal("hook did not run for an event with no registered listeners")
	}
}

func TestHandleApplicationEventDeliversToListeners(t *testing.T) {
	em := newEventManager(newEventManagerTestApp())

	delivered := make(chan struct{})
	em.OnApplicationEvent(events.Common.ApplicationStarted, func(event *ApplicationEvent) {
		close(delivered)
	})

	em.handleApplicationEvent(newApplicationEvent(events.Common.ApplicationStarted))

	select {
	case <-delivered:
	case <-time.After(5 * time.Second):
		t.Fatal("listener was not invoked")
	}
}

func TestHandleApplicationEventHookCancelStopsListeners(t *testing.T) {
	em := newEventManager(newEventManagerTestApp())

	em.RegisterApplicationEventHook(events.Common.ApplicationStarted, func(event *ApplicationEvent) {
		event.Cancel()
	})
	invoked := make(chan struct{}, 1)
	em.OnApplicationEvent(events.Common.ApplicationStarted, func(event *ApplicationEvent) {
		invoked <- struct{}{}
	})

	em.handleApplicationEvent(newApplicationEvent(events.Common.ApplicationStarted))

	select {
	case <-invoked:
		t.Fatal("listener ran although a hook cancelled the event")
	case <-time.After(100 * time.Millisecond):
	}
}
