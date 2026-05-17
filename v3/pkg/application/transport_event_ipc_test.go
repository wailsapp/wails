package application

import (
	"sync"
	"testing"
	"time"
)

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
