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

func TestEventProcessorDeliversToGoListenerInEmitOrder(t *testing.T) {
	processor := application.NewWailsEventProcessor(func(*application.CustomEvent) {})
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondStarted := make(chan struct{})
	done := make(chan struct{}, 2)

	var lock sync.Mutex
	var order []int
	processor.On("test", func(event *application.CustomEvent) {
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
		done <- struct{}{}
	})

	_ = processor.Emit(&application.CustomEvent{Name: "test", Data: 1})
	<-firstStarted
	_ = processor.Emit(&application.CustomEvent{Name: "test", Data: 2})

	select {
	case <-secondStarted:
		t.Fatal("second event overtook the first in the same listener")
	case <-time.After(100 * time.Millisecond):
	}

	close(releaseFirst)
	<-done
	<-done

	lock.Lock()
	defer lock.Unlock()
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("listener order = %v, want [1 2]", order)
	}
}

func TestEventProcessorGoListenerSeesBurstInOrderOneAtATime(t *testing.T) {
	const count = 500
	processor := application.NewWailsEventProcessor(func(*application.CustomEvent) {})

	var wg sync.WaitGroup
	wg.Add(count)
	var running, maxRunning atomic.Int32
	var lock sync.Mutex
	var order []int
	processor.On("burst", func(event *application.CustomEvent) {
		now := running.Add(1)
		for {
			highest := maxRunning.Load()
			if now <= highest || maxRunning.CompareAndSwap(highest, now) {
				break
			}
		}
		lock.Lock()
		order = append(order, event.Data.(int))
		lock.Unlock()
		running.Add(-1)
		wg.Done()
	})

	for value := range count {
		_ = processor.Emit(&application.CustomEvent{Name: "burst", Data: value})
	}
	wg.Wait()

	if got := maxRunning.Load(); got != 1 {
		t.Fatalf("listener ran %d times concurrently, want 1", got)
	}
	lock.Lock()
	defer lock.Unlock()
	for index, value := range order {
		if value != index {
			t.Fatalf("event %d delivered at position %d", value, index)
		}
	}
}

func TestEventProcessorSlowListenerDoesNotDelayOtherListeners(t *testing.T) {
	processor := application.NewWailsEventProcessor(func(*application.CustomEvent) {})
	releaseSlow := make(chan struct{})
	defer close(releaseSlow)
	fastCalled := make(chan struct{}, 1)

	processor.On("test", func(*application.CustomEvent) { <-releaseSlow })
	processor.On("test", func(*application.CustomEvent) { fastCalled <- struct{}{} })

	_ = processor.Emit(&application.CustomEvent{Name: "test"})

	select {
	case <-fastCalled:
	case <-time.After(time.Second):
		t.Fatal("a blocked listener delayed another listener of the same event")
	}
}

func TestEventProcessorEmitDoesNotWaitForListeners(t *testing.T) {
	processor := application.NewWailsEventProcessor(func(*application.CustomEvent) {})
	releaseListener := make(chan struct{})
	defer close(releaseListener)
	processor.On("test", func(*application.CustomEvent) { <-releaseListener })

	emitted := make(chan struct{})
	go func() {
		for range 100 {
			_ = processor.Emit(&application.CustomEvent{Name: "test"})
		}
		close(emitted)
	}()

	select {
	case <-emitted:
	case <-time.After(time.Second):
		t.Fatal("Emit blocked on a busy listener")
	}
}

func TestEventProcessorOnceRunsOnceUnderBurst(t *testing.T) {
	processor := application.NewWailsEventProcessor(func(*application.CustomEvent) {})
	var calls atomic.Int32
	processor.Once("test", func(*application.CustomEvent) { calls.Add(1) })

	var emitters sync.WaitGroup
	for range 8 {
		emitters.Add(1)
		go func() {
			defer emitters.Done()
			for range 50 {
				_ = processor.Emit(&application.CustomEvent{Name: "test"})
			}
		}()
	}
	emitters.Wait()
	time.Sleep(50 * time.Millisecond)

	if got := calls.Load(); got != 1 {
		t.Fatalf("Once listener ran %d times, want 1", got)
	}
}

func TestEventProcessorDropsQueuedEventsAfterUnsubscribe(t *testing.T) {
	unsubscribers := map[string]func(processor *application.EventProcessor, off func()){
		"returned function": func(_ *application.EventProcessor, off func()) { off() },
		"Off":               func(processor *application.EventProcessor, _ func()) { processor.Off("test") },
		"OffAll":            func(processor *application.EventProcessor, _ func()) { processor.OffAll() },
	}
	for name, unsubscribe := range unsubscribers {
		t.Run(name, func(t *testing.T) {
			processor := application.NewWailsEventProcessor(func(*application.CustomEvent) {})
			firstStarted := make(chan struct{})
			releaseFirst := make(chan struct{})
			var lock sync.Mutex
			var delivered []int
			off := processor.On("test", func(event *application.CustomEvent) {
				value := event.Data.(int)
				if value == 1 {
					close(firstStarted)
					<-releaseFirst
				}
				lock.Lock()
				delivered = append(delivered, value)
				lock.Unlock()
			})

			_ = processor.Emit(&application.CustomEvent{Name: "test", Data: 1})
			<-firstStarted
			_ = processor.Emit(&application.CustomEvent{Name: "test", Data: 2})
			_ = processor.Emit(&application.CustomEvent{Name: "test", Data: 3})
			unsubscribe(processor, off)
			close(releaseFirst)
			time.Sleep(50 * time.Millisecond)

			lock.Lock()
			defer lock.Unlock()
			if len(delivered) != 1 || delivered[0] != 1 {
				t.Fatalf("delivered %v after unsubscribing, want only the event already running [1]", delivered)
			}
		})
	}
}

func TestEventProcessorRegisterHookWhileEmitting(t *testing.T) {
	processor := application.NewWailsEventProcessor(func(*application.CustomEvent) {})
	stop := make(chan struct{})
	var emitter sync.WaitGroup
	emitter.Add(1)
	go func() {
		defer emitter.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_ = processor.Emit(&application.CustomEvent{Name: "test"})
			}
		}
	}()

	for range 200 {
		remove := processor.RegisterHook("test", func(*application.CustomEvent) {})
		remove()
	}
	close(stop)
	emitter.Wait()
}
