package application_test

import (
	"sync"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/matryer/is"
)

type mockNotifier struct {
	Events []*application.CustomEvent
}

func (m *mockNotifier) dispatchEventToWindows(event *application.CustomEvent) {
	m.Events = append(m.Events, event)
}

func (m *mockNotifier) Reset() {
	m.Events = []*application.CustomEvent{}
}

func Test_EventsOn(t *testing.T) {
	i := is.New(t)
	notifier := &mockNotifier{}
	eventProcessor := application.NewWailsEventProcessor(notifier.dispatchEventToWindows)

	// Test OnApplicationEvent
	eventName := "test"
	counter := 0
	var wg sync.WaitGroup
	wg.Add(1)
	unregisterFn := eventProcessor.On(eventName, func(event *application.CustomEvent) {
		// This is called in a goroutine
		counter++
		wg.Done()
	})
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	wg.Wait()
	i.Equal(1, counter)

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter = 0
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	i.Equal(0, counter)

}

func Test_EventsOnce(t *testing.T) {
	i := is.New(t)
	notifier := &mockNotifier{}
	eventProcessor := application.NewWailsEventProcessor(notifier.dispatchEventToWindows)

	// Test OnApplicationEvent
	eventName := "test"
	counter := 0
	var wg sync.WaitGroup
	wg.Add(1)
	unregisterFn := eventProcessor.Once(eventName, func(event *application.CustomEvent) {
		// This is called in a goroutine
		counter++
		wg.Done()
	})
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	wg.Wait()
	i.Equal(1, counter)

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter = 0
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	i.Equal(0, counter)

}
func Test_EventsOnMultiple(t *testing.T) {
	i := is.New(t)
	notifier := &mockNotifier{}
	eventProcessor := application.NewWailsEventProcessor(notifier.dispatchEventToWindows)

	// Test OnApplicationEvent
	eventName := "test"
	counter := 0
	var wg sync.WaitGroup
	wg.Add(2)
	unregisterFn := eventProcessor.OnMultiple(eventName, func(event *application.CustomEvent) {
		// This is called in a goroutine
		counter++
		wg.Done()
	}, 2)
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	wg.Wait()
	i.Equal(2, counter)

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter = 0
	eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	i.Equal(0, counter)

}
