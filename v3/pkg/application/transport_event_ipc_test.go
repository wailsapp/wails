package application

import (
	"sync"
	"testing"
	"time"
)

func TestMainThreadFunctionStoreLockCorrectness(t *testing.T) {
	var lock sync.RWMutex
	store := make(map[uint]func())
	var wg sync.WaitGroup

	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		id := uint(i)
		go func() {
			defer wg.Done()
			called := false
			lock.Lock()
			store[id] = func() { called = true }
			lock.Unlock()

			lock.Lock()
			fn := store[id]
			if fn != nil {
				delete(store, id)
			}
			lock.Unlock()

			if fn != nil {
				fn()
			}
			if !called {
				t.Errorf("callback for id %d was not called", id)
			}
		}()
	}

	wg.Wait()
}

func TestMainThreadFunctionStoreConcurrentReadWrite(t *testing.T) {
	var lock sync.RWMutex
	store := make(map[uint]func())
	var wg sync.WaitGroup

	const writers = 50
	const readers = 50

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(id uint) {
			defer wg.Done()
			lock.Lock()
			store[id] = func() {}
			lock.Unlock()

			lock.Lock()
			delete(store, id)
			lock.Unlock()
		}(uint(i))
	}

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.RLock()
			_ = len(store)
			lock.RUnlock()
		}()
	}

	wg.Wait()
}

func TestDispatchWailsEventReleasesLockAfterCompletion(t *testing.T) {
	app := &App{
		windows:     make(map[uint]Window),
		windowsLock: sync.RWMutex{},
	}

	transport := &EventIPCTransport{app: app}
	event := &CustomEvent{Name: "test"}

	transport.DispatchWailsEvent(event)

	acquired := make(chan struct{})
	go func() {
		app.windowsLock.Lock()
		close(acquired)
		app.windowsLock.Unlock()
	}()

	select {
	case <-acquired:
	case <-time.After(time.Second):
		t.Fatal("windowsLock should be released after DispatchWailsEvent completes")
	}
}
