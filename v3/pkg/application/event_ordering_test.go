package application

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// threadProbeApp stands in for the platform layer so a test can control which
// thread it appears to be on and when queued work actually runs.
//
// Embedding platformApp supplies nil-bodied placeholders for everything the
// dispatch path never touches.
type threadProbeApp struct {
	platformApp

	onMain  atomic.Bool
	mu      sync.Mutex
	pending []uint
}

func (f *threadProbeApp) isOnMainThread() bool { return f.onMain.Load() }

// dispatchOnMainThread records the work instead of running it, so a test can
// hold the UI thread still and then release it.
func (f *threadProbeApp) dispatchOnMainThread(id uint) {
	f.mu.Lock()
	f.pending = append(f.pending, id)
	f.mu.Unlock()
}

// runPending executes the deferred work in the order it was dispatched, which
// is what the real main-thread queue does.
func (f *threadProbeApp) runPending() {
	for {
		f.mu.Lock()
		if len(f.pending) == 0 {
			f.mu.Unlock()
			return
		}
		id := f.pending[0]
		f.pending = f.pending[1:]
		f.mu.Unlock()

		mainThreadFunctionStoreLock.Lock()
		fn := mainThreadFunctionStore[id]
		delete(mainThreadFunctionStore, id)
		mainThreadFunctionStoreLock.Unlock()

		if fn == nil {
			continue
		}
		// Anything dispatched this way is, by definition, running on the UI
		// thread — so report that while it runs. Without this, an InvokeSync
		// inside the callback would try to hop to a thread it is already on
		// and wait forever.
		prev := f.onMain.Load()
		f.onMain.Store(true)
		fn()
		f.onMain.Store(prev)
	}
}

// stubWindowImpl records what actually reached the webview, in the order it
// got there. That order is the thing under test.
type stubWindowImpl struct {
	webviewWindowImpl
	mu   sync.Mutex
	seen []string
}

func (s *stubWindowImpl) execJS(js string) {
	s.mu.Lock()
	s.seen = append(s.seen, js)
	s.mu.Unlock()
}

func (s *stubWindowImpl) delivered() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.seen...)
}

func (f *threadProbeApp) pendingCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.pending)
}

// newOrderingProbe wires a window to a controllable platform layer. Leaving
// runtimeLoaded false makes ExecJS record into pendingJS instead of hopping to
// the UI thread, so pendingJS is the delivery order.
func newOrderingProbe(t *testing.T) (*threadProbeApp, *WebviewWindow, func()) {
	t.Helper()

	probe := &threadProbeApp{}
	prev := globalApplication
	globalApplication = &App{impl: probe}

	// runtimeLoaded must be true: with it false ExecJS short-circuits into
	// pendingJS and never touches the main-thread dispatch that this is all
	// about, which would make the test pass whether or not the bug is present.
	win := &WebviewWindow{
		id:            1,
		impl:          &stubWindowImpl{},
		runtimeLoaded: true,
	}

	return probe, win, func() { globalApplication = prev }
}

// A UI-thread emit must not overtake an event a goroutine queued earlier.
//
// This is the regression test for the inline fast path in
// App.dispatchOnMainThread: when the caller is already on the UI thread the
// work runs immediately, so before the queue existed a main-thread emit
// executed its eval while an earlier goroutine emit was still waiting to be
// dispatched. Measured at ~4.4% of events inverted under two concurrent
// emitters on macOS, Linux and Windows.
func TestEventFromMainThreadDoesNotOvertakeQueuedEvent(t *testing.T) {
	probe, win, restore := newOrderingProbe(t)
	defer restore()
	impl := win.impl.(*stubWindowImpl)

	// Emitted from a goroutine. Without the queue this parks in InvokeSync
	// waiting for the UI thread, so it runs in its own goroutine either way.
	probe.onMain.Store(false)
	firstReturned := make(chan struct{})
	go func() {
		defer close(firstReturned)
		win.DispatchWailsEvent(&CustomEvent{Name: "first"})
	}()

	// Wait until it has reached the UI-thread dispatch queue.
	deadline := time.Now().Add(5 * time.Second)
	for probe.pendingCount() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if probe.pendingCount() == 0 {
		t.Fatal("the goroutine emit never reached the main-thread dispatch queue")
	}

	// Now emit from the UI thread, while that one is still outstanding.
	// dispatchOnMainThread runs inline here, which is exactly how the later
	// event used to overtake the earlier one.
	probe.onMain.Store(true)
	win.DispatchWailsEvent(&CustomEvent{Name: "second"})

	// Release the UI thread.
	probe.runPending()

	select {
	case <-firstReturned:
	case <-time.After(5 * time.Second):
		t.Fatal("the goroutine emit never completed")
	}

	got := impl.delivered()
	if len(got) != 2 {
		t.Fatalf("delivered %d events, want 2: %v", len(got), got)
	}
	if !strings.Contains(got[0], `"first"`) || !strings.Contains(got[1], `"second"`) {
		t.Fatalf("the later main-thread event overtook the earlier queued one\n  1st delivered: %s\n  2nd delivered: %s", got[0], got[1])
	}
}

// Events emitted in sequence from one goroutine must arrive in that sequence.
func TestEventQueuePreservesSingleProducerOrder(t *testing.T) {
	probe, win, restore := newOrderingProbe(t)
	defer restore()

	probe.onMain.Store(false)

	// More events than the queue holds, so the drain has to run concurrently
	// with the producer; emitting them all first would (correctly) block on the
	// bound with nobody there to relieve it.
	const n = eventQueueCapacity * 4

	stop := make(chan struct{})
	drained := make(chan struct{})
	go func() {
		defer close(drained)
		for {
			probe.runPending()
			select {
			case <-stop:
				probe.runPending() // final sweep
				return
			default:
			}
		}
	}()

	for i := 0; i < n; i++ {
		win.enqueueEventJS(fmt.Sprintf("e%d", i))
	}
	close(stop)
	<-drained

	deadline := time.Now().Add(10 * time.Second)
	for {
		delivered := len(win.impl.(*stubWindowImpl).delivered())
		if delivered == n || time.Now().After(deadline) {
			break
		}
		probe.runPending()
	}

	got := win.impl.(*stubWindowImpl).delivered()

	if len(got) != n {
		t.Fatalf("delivered %d events, want %d", len(got), n)
	}
	for i := 0; i < n; i++ {
		if want := fmt.Sprintf("e%d", i); got[i] != want {
			t.Fatalf("event %d = %q, want %q", i, got[i], want)
		}
	}
}

// A UI-thread emitter must never block, however many events it emits: it is
// the drainer, so blocking it could not be relieved by anyone. Emitting far
// more than the queue holds must still complete, in order.
//
// In practice it never even approaches the bound, because the drain scheduled
// by InvokeAsync runs inline when already on the UI thread — so a main-thread
// emitter enqueues and immediately drains. The bound exists for goroutine
// emitters racing a busy UI thread.
func TestEventQueueDoesNotBlockMainThreadWhenFull(t *testing.T) {
	probe, win, restore := newOrderingProbe(t)
	defer restore()

	probe.onMain.Store(true)

	const n = eventQueueCapacity * 3

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < n; i++ {
			win.enqueueEventJS(fmt.Sprintf("e%d", i))
		}
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("a main-thread emitter blocked on a full queue; this deadlocks in a real app")
	}

	got := win.impl.(*stubWindowImpl).delivered()

	if len(got) != n {
		t.Fatalf("delivered %d events, want %d", len(got), n)
	}
	for i := 0; i < n; i++ {
		if want := fmt.Sprintf("e%d", i); got[i] != want {
			t.Fatalf("event %d = %q, want %q", i, got[i], want)
		}
	}

	// Draining inline is what keeps the UI thread from ever waiting.
	if high := win.eventQueueHighWater(); high > eventQueueCapacity {
		t.Errorf("high water = %d; a main-thread emitter should drain inline, not accumulate", high)
	}
}

// Destroying a window must release an emitter waiting for queue space rather
// than leaving the goroutine parked forever.
func TestEventQueueCloseReleasesBlockedProducer(t *testing.T) {
	probe, win, restore := newOrderingProbe(t)
	defer restore()

	probe.onMain.Store(false)

	// Fill to capacity so the next append has to wait.
	for i := 0; i < eventQueueCapacity; i++ {
		win.enqueueEventJS(fmt.Sprintf("e%d", i))
	}

	// Wait until the producer is genuinely parked before closing, otherwise
	// close can win the race and the test passes without ever exercising the
	// Broadcast that is the thing under test.
	blocked := make(chan struct{})
	go func() {
		defer close(blocked)
		win.enqueueEventJS("waits for space")
	}()

	deadline := time.Now().Add(5 * time.Second)
	for {
		win.eventQueueMu.Lock()
		waiting := len(win.eventQueue) >= eventQueueCapacity
		win.eventQueueMu.Unlock()

		select {
		case <-blocked:
			t.Fatal("the producer returned without waiting; the queue was not full")
		default:
		}
		if waiting && time.Now().After(deadline.Add(-4900*time.Millisecond)) {
			break // queue is full and the producer has had a chance to park
		}
		if time.Now().After(deadline) {
			t.Fatal("producer never reached the full queue")
		}
		time.Sleep(2 * time.Millisecond)
	}

	win.closeEventQueue()

	select {
	case <-blocked:
	case <-time.After(10 * time.Second):
		t.Fatal("closing the queue did not wake the blocked emitter")
	}
}

// A payload parked just as the window goes away must not be stranded. The
// dispatcher can pass the destroyed check, markAsDestroyed can then close the
// queue and run dropWindow, and only then does put succeed — so dropWindow
// cannot see it and the queue refuses the event that would have fetched it.
func TestOrphanedPayloadIsReclaimedWhenQueueRefuses(t *testing.T) {
	probe, win, restore := newOrderingProbe(t)
	defer restore()
	probe.onMain.Store(false)

	store := newEventPayloadStore()
	globalApplication.eventPayloads = store

	// Simulate the interleaving: the queue is already closed by the time the
	// dispatcher reaches it.
	win.closeEventQueue()

	big := make([]byte, maxInlineEventPayload+1)
	for i := range big {
		big[i] = 'x'
	}
	win.DispatchWailsEvent(&CustomEvent{Name: "big", Data: string(big)})

	store.mu.Lock()
	parked, bytes := len(store.items), store.bytes
	store.mu.Unlock()

	if parked != 0 || bytes != 0 {
		t.Errorf("payload left stranded in the store: %d entries, %d bytes", parked, bytes)
	}
}
