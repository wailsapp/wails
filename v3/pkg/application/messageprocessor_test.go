package application

import (
	"context"
	"strconv"
	"sync"
	"testing"
)

// registerCall is a test helper that registers a single in-flight call under
// the given window ID using the same internal bookkeeping that processCallMethod
// uses. It returns the call's context so the caller can observe cancellation.
func registerCall(m *MessageProcessor, windowID uint, callID string) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	m.l.Lock()
	defer m.l.Unlock()
	m.runningCalls[callID] = cancel
	if m.windowCalls[windowID] == nil {
		m.windowCalls[windowID] = make(map[string]bool)
	}
	m.windowCalls[windowID][callID] = true
	return ctx
}

func TestMessageProcessor_CancelWindowCalls_CancelsAllForWindow(t *testing.T) {
	m := NewMessageProcessor(nil)
	ctxA := registerCall(m, 7, "call-a")
	ctxB := registerCall(m, 7, "call-b")

	m.CancelWindowCalls(7)

	select {
	case <-ctxA.Done():
	default:
		t.Fatal("expected ctxA to be cancelled")
	}
	select {
	case <-ctxB.Done():
	default:
		t.Fatal("expected ctxB to be cancelled")
	}

	m.l.Lock()
	defer m.l.Unlock()
	if _, ok := m.runningCalls["call-a"]; ok {
		t.Error("runningCalls should not contain call-a after cancel")
	}
	if _, ok := m.runningCalls["call-b"]; ok {
		t.Error("runningCalls should not contain call-b after cancel")
	}
	if _, ok := m.windowCalls[7]; ok {
		t.Error("windowCalls[7] should be removed after cancel")
	}
}

func TestMessageProcessor_CancelWindowCalls_PreservesOtherWindows(t *testing.T) {
	m := NewMessageProcessor(nil)
	ctxW1 := registerCall(m, 1, "w1-call")
	ctxW2 := registerCall(m, 2, "w2-call")

	m.CancelWindowCalls(1)

	select {
	case <-ctxW1.Done():
	default:
		t.Fatal("expected window 1 call to be cancelled")
	}
	select {
	case <-ctxW2.Done():
		t.Fatal("window 2 call must not be cancelled when only window 1 is closing")
	default:
	}

	m.l.Lock()
	defer m.l.Unlock()
	if _, ok := m.runningCalls["w2-call"]; !ok {
		t.Error("runningCalls should still contain w2-call")
	}
	if _, ok := m.windowCalls[2]; !ok {
		t.Error("windowCalls[2] should still be tracked")
	}
	if _, ok := m.windowCalls[1]; ok {
		t.Error("windowCalls[1] should be removed")
	}
}

func TestMessageProcessor_CancelWindowCalls_UnknownWindowIsNoOp(t *testing.T) {
	m := NewMessageProcessor(nil)

	// Empty processor: cancelling an unknown window must not panic.
	m.CancelWindowCalls(99)

	// With one call registered under a different window, cancelling an
	// unknown window must not disturb it.
	ctx := registerCall(m, 1, "keep-me")
	m.CancelWindowCalls(99)

	select {
	case <-ctx.Done():
		t.Fatal("registered call must not be cancelled when an unrelated window closes")
	default:
	}

	m.l.Lock()
	defer m.l.Unlock()
	if _, ok := m.runningCalls["keep-me"]; !ok {
		t.Error("registered call should still be present")
	}
}

// If a callID is recorded in windowCalls but has already been removed from
// runningCalls (the call completed and self-removed concurrently with the
// window destroy), CancelWindowCalls must not panic and must still clean up
// the windowCalls bookkeeping. This guards the
// `if cancel, ok := m.runningCalls[callID]; ok` branch.
func TestMessageProcessor_CancelWindowCalls_StaleEntryCleansUp(t *testing.T) {
	m := NewMessageProcessor(nil)
	m.l.Lock()
	m.windowCalls[5] = map[string]bool{"stale-id": true}
	m.l.Unlock()

	m.CancelWindowCalls(5)

	m.l.Lock()
	defer m.l.Unlock()
	if _, ok := m.windowCalls[5]; ok {
		t.Error("stale windowCalls[5] entry should still be cleared")
	}
}

// Ensures CancelWindowCalls does not deadlock or corrupt internal maps under
// concurrent registration. Run with `go test -race`.
func TestMessageProcessor_CancelWindowCalls_ConcurrentSafe(t *testing.T) {
	m := NewMessageProcessor(nil)
	const windows = 4
	const callsPerWindow = 25

	var wg sync.WaitGroup
	for w := uint(1); w <= windows; w++ {
		wg.Add(1)
		go func(windowID uint) {
			defer wg.Done()
			for i := 0; i < callsPerWindow; i++ {
				registerCall(m, windowID, callIDFor(windowID, i))
			}
		}(w)
	}
	wg.Wait()

	// Cancel every window concurrently.
	var cancelWG sync.WaitGroup
	for w := uint(1); w <= windows; w++ {
		cancelWG.Add(1)
		go func(windowID uint) {
			defer cancelWG.Done()
			m.CancelWindowCalls(windowID)
		}(w)
	}
	cancelWG.Wait()

	m.l.Lock()
	defer m.l.Unlock()
	if len(m.runningCalls) != 0 {
		t.Errorf("expected runningCalls empty, got %d entries", len(m.runningCalls))
	}
	if len(m.windowCalls) != 0 {
		t.Errorf("expected windowCalls empty, got %d entries", len(m.windowCalls))
	}
}

func callIDFor(windowID uint, i int) string {
	return "w" + strconv.FormatUint(uint64(windowID), 10) + "-" + strconv.Itoa(i)
}
