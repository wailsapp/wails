package application

import (
	"testing"
	"time"
)

// TestWindowEmitEventStillReachesHooksAndGoListeners pins down that targeting
// a window only narrows frontend delivery: hooks and Go listeners keep
// observing window-emitted events, and the event carries its target to the
// frontend transports.
func TestWindowEmitEventStillReachesHooksAndGoListeners(t *testing.T) {
	transported := make(chan *CustomEvent, 1)
	app := &App{}
	app.Event = newEventManager(app)
	app.customEventProcessor = NewWailsEventProcessor(func(event *CustomEvent) {
		transported <- event
	})

	prev := globalApplication
	globalApplication = app
	defer func() { globalApplication = prev }()

	window := &WebviewWindow{id: 7, options: WebviewWindowOptions{Name: "second"}}

	hooked := make(chan *CustomEvent, 1)
	app.customEventProcessor.RegisterHook("notification", func(event *CustomEvent) { hooked <- event })
	listened := make(chan *CustomEvent, 1)
	app.Event.On("notification", func(event *CustomEvent) { listened <- event })

	if cancelled := window.EmitEvent("notification", "Hello from Go!"); cancelled {
		t.Fatal("event was unexpectedly cancelled")
	}

	for name, ch := range map[string]chan *CustomEvent{
		"hook":      hooked,
		"listener":  listened,
		"transport": transported,
	} {
		select {
		case event := <-ch:
			if event.Target() != Window(window) {
				t.Errorf("%s: Target() = %v, want the emitting window", name, event.Target())
			}
			if event.Sender != "second" {
				t.Errorf("%s: Sender = %q, want %q", name, event.Sender, "second")
			}
		case <-time.After(time.Second):
			t.Fatalf("%s did not receive the window-emitted event", name)
		}
	}
}

func TestTargetIsNotSerialised(t *testing.T) {
	window := &WebviewWindow{id: 7}
	event := &CustomEvent{Name: "n", Data: 1, Sender: "s", target: window}
	if got, want := event.ToJSON(), `{"name":"n","data":1,"sender":"s"}`; got != want {
		t.Fatalf("ToJSON() = %s, want %s", got, want)
	}
}
