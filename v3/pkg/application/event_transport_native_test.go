//go:build !server

package application

import (
	"fmt"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type readinessWindowImpl struct {
	stubWindowImpl
	resizable bool
}

func (impl *readinessWindowImpl) setResizable(value bool) { impl.resizable = value }

func TestNativeDefaultTransportPreservesPreReadyEvents(t *testing.T) {
	a := newEventTransportTestApp(t, Options{})
	probe := &threadProbeApp{}
	probe.onMain.Store(true)
	a.impl = probe
	impl := &readinessWindowImpl{}
	window := &WebviewWindow{id: 1, impl: impl}
	a.Window.Add(window)

	first := &CustomEvent{Name: "before-ready-first", Data: "first payload"}
	second := &CustomEvent{Name: "before-ready-second", Data: "second payload"}
	a.Event.dispatch(first)
	a.Event.dispatch(second)

	window.pendingJSMutex.Lock()
	pending := append([]string(nil), window.pendingJS...)
	ready := window.runtimeLoaded
	window.pendingJSMutex.Unlock()
	if ready || len(pending) != 2 {
		t.Fatalf("native pre-ready state: ready=%t, pending=%d; want false and two", ready, len(pending))
	}
	if len(impl.delivered()) != 0 {
		t.Fatal("native transport executed JavaScript before runtime readiness")
	}

	// Exercise the real runtime-ready message handler, not a manual queue flush.
	// No platform event loop is running in this test, so isolate its notification.
	previousEvents := windowEvents
	windowEvents = make(chan *windowEvent, 1)
	defer func() { windowEvents = previousEvents }()
	window.HandleMessage("wails:runtime:ready")
	select {
	case event := <-windowEvents:
		if event.WindowID != window.ID() || event.EventID != uint(events.Common.WindowRuntimeReady) {
			t.Fatalf("runtime-ready notification = %#v", event)
		}
	default:
		t.Fatal("runtime-ready handler did not emit its window notification")
	}
	if !window.runtimeLoaded || len(window.pendingJS) != 0 || !impl.resizable {
		t.Fatal("runtime-ready handler did not complete readiness and clear pending JavaScript")
	}
	third := &CustomEvent{Name: "after-ready", Data: "third payload"}
	a.Event.dispatch(third)
	delivered := impl.delivered()
	if len(delivered) != 3 {
		t.Fatalf("native runtime received %d events, want three exactly once", len(delivered))
	}
	for index, event := range []*CustomEvent{first, second, third} {
		if want := fmt.Sprintf(inlineEventJS, event.ToJSON()); delivered[index] != want {
			t.Fatalf("native event %d changed payload or order across readiness", index)
		}
	}
}
