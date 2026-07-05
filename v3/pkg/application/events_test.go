package application_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/matryer/is"
)

type mockNotifier struct {
	mu     sync.Mutex
	Events []*application.CustomEvent
}

// mu: dispatch now runs on a persistent worker (see runWindowDispatch), so it
// can race with a test's Reset() right after Emit(). Mutex makes that safe.
func (m *mockNotifier) dispatchEventToWindows(event *application.CustomEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Events = append(m.Events, event)
}

func (m *mockNotifier) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
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
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	wg.Wait()
	i.Equal(1, counter)

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter = 0
	_ = eventProcessor.Emit(&application.CustomEvent{
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
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	wg.Wait()
	i.Equal(1, counter)

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter = 0
	_ = eventProcessor.Emit(&application.CustomEvent{
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
	// atomic: 2 of the 3 Emit()s below can invoke this callback concurrently.
	var counter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)
	unregisterFn := eventProcessor.OnMultiple(eventName, func(event *application.CustomEvent) {
		// This is called in a goroutine
		counter.Add(1)
		wg.Done()
	}, 2)
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	wg.Wait()
	i.Equal(int32(2), counter.Load())

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter.Store(0)
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	i.Equal(int32(0), counter.Load())

}
