package application

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestHTTPTransportCleanupLifecycle(t *testing.T) {
	transport := NewHTTPTransport()
	t.Cleanup(func() { _ = transport.Stop() })

	for cycle := 0; cycle < 100; cycle++ {
		if err := transport.Start(context.Background(), nil); err != nil {
			t.Fatalf("cycle %d: Start: %v", cycle, err)
		}

		transport.cleanupMu.Lock()
		done := transport.cleanupDone
		transport.cleanupMu.Unlock()
		if done == nil {
			t.Fatalf("cycle %d: cleanup goroutine was not started", cycle)
		}
		if err := transport.Start(context.Background(), nil); err != nil {
			t.Fatalf("cycle %d: repeated Start: %v", cycle, err)
		}
		transport.cleanupMu.Lock()
		repeatedDone := transport.cleanupDone
		transport.cleanupMu.Unlock()
		if repeatedDone != done {
			t.Fatalf("cycle %d: repeated Start launched another cleanup goroutine", cycle)
		}

		// Stop is allowed to be called defensively from several shutdown paths.
		// It must close the channel only once and the call that owns the cleanup
		// goroutine must wait for it to finish before returning.
		var stops sync.WaitGroup
		stops.Add(4)
		for range 4 {
			go func() {
				defer stops.Done()
				if err := transport.Stop(); err != nil {
					t.Errorf("cycle %d: Stop: %v", cycle, err)
				}
			}()
		}
		stops.Wait()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("cycle %d: cleanup goroutine did not stop", cycle)
		}
	}
}
