package application_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/matryer/is"
)

type mockNotifier struct {
	mu     sync.Mutex
	Events []*application.CustomEvent
}

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
	var counter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(1)
	unregisterFn := eventProcessor.On(eventName, func(event *application.CustomEvent) {
		// This is called in a goroutine
		counter.Add(1)
		wg.Done()
	})
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	wg.Wait()
	i.Equal(1, int(counter.Load()))

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter.Store(0)
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	i.Equal(0, int(counter.Load()))

}

func Test_EventsOnce(t *testing.T) {
	i := is.New(t)
	notifier := &mockNotifier{}
	eventProcessor := application.NewWailsEventProcessor(notifier.dispatchEventToWindows)

	// Test OnApplicationEvent
	eventName := "test"
	var counter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(1)
	unregisterFn := eventProcessor.Once(eventName, func(event *application.CustomEvent) {
		// This is called in a goroutine
		counter.Add(1)
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
	i.Equal(1, int(counter.Load()))

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter.Store(0)
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	i.Equal(0, int(counter.Load()))

}
func Test_EventsOnMultiple(t *testing.T) {
	i := is.New(t)
	notifier := &mockNotifier{}
	eventProcessor := application.NewWailsEventProcessor(notifier.dispatchEventToWindows)

	// Test OnApplicationEvent
	eventName := "test"
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
	i.Equal(2, int(counter.Load()))

	// Unregister
	notifier.Reset()
	unregisterFn()
	counter.Store(0)
	_ = eventProcessor.Emit(&application.CustomEvent{
		Name: "test",
		Data: "test payload",
	})
	i.Equal(0, int(counter.Load()))

}

func TestEventProcessorDispatchesWindowEventsInEmitOrder(t *testing.T) {
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondStarted := make(chan struct{})
	dispatched := make(chan struct{}, 2)

	var lock sync.Mutex
	var order []int
	processor := application.NewWailsEventProcessor(func(event *application.CustomEvent) {
		value := event.Data.(int)
		switch value {
		case 1:
			close(firstStarted)
			<-releaseFirst
		case 2:
			close(secondStarted)
		}

		lock.Lock()
		order = append(order, value)
		lock.Unlock()
		dispatched <- struct{}{}
	})

	_ = processor.Emit(&application.CustomEvent{Name: "test", Data: 1})
	<-firstStarted
	_ = processor.Emit(&application.CustomEvent{Name: "test", Data: 2})

	select {
	case <-secondStarted:
		t.Fatal("second window event overtook the first")
	case <-time.After(100 * time.Millisecond):
	}

	close(releaseFirst)
	<-dispatched
	<-dispatched

	lock.Lock()
	defer lock.Unlock()
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("window event order = %v, want [1 2]", order)
	}
}
