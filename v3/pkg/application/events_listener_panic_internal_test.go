package application

import (
	"sync"
	"testing"
	"time"
)

func TestEventListenerKeepsDeliveringAfterPanic(t *testing.T) {
	var lock sync.Mutex
	var panics int
	prev := globalApplication
	globalApplication = &App{options: Options{PanicHandler: func(*PanicDetails) {
		lock.Lock()
		panics++
		lock.Unlock()
	}}}
	defer func() { globalApplication = prev }()

	processor := NewWailsEventProcessor(func(*CustomEvent) {})
	delivered := make(chan int, 1)
	processor.On("test", func(event *CustomEvent) {
		if event.Data.(int) == 1 {
			panic("listener failure")
		}
		delivered <- event.Data.(int)
	})

	_ = processor.Emit(&CustomEvent{Name: "test", Data: 1})
	_ = processor.Emit(&CustomEvent{Name: "test", Data: 2})

	select {
	case value := <-delivered:
		if value != 2 {
			t.Fatalf("delivered %d, want 2", value)
		}
	case <-time.After(time.Second):
		t.Fatal("a panicking handler stopped delivery of later events")
	}
	lock.Lock()
	defer lock.Unlock()
	if panics != 1 {
		t.Fatalf("panic handler called %d times, want 1", panics)
	}
}
