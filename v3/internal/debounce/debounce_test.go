package debounce

import (
	"sync/atomic"
	"testing"
	"time"
)

const (
	// debounceAfter must be well above OS timer granularity (~15ms on Windows).
	debounceAfter = 100 * time.Millisecond
	// waitFor must exceed debounceAfter by enough margin for the callback to fire.
	waitFor = 300 * time.Millisecond
)

// ---------------------------------------------------------------------------
// Basic: single call fires after duration
// ---------------------------------------------------------------------------

func TestBasic_SingleCallFires(t *testing.T) {
	debounced := New(debounceAfter)

	var count int64
	debounced(func() {
		atomic.AddInt64(&count, 1)
	})

	time.Sleep(waitFor)

	if got := atomic.LoadInt64(&count); got != 1 {
		t.Errorf("expected 1 call, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// Debouncing: rapid calls, only last fires
// ---------------------------------------------------------------------------

func TestDebouncing_RapidCallsOnlyLastFires(t *testing.T) {
	debounced := New(debounceAfter)

	var count int64
	for i := 0; i < 10; i++ {
		debounced(func() {
			atomic.AddInt64(&count, 1)
		})
		// Sleep well below debounceAfter so each call resets the timer before it fires.
		// Must be > 0 to avoid a tight loop but << debounceAfter.
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(waitFor)

	if got := atomic.LoadInt64(&count); got != 1 {
		t.Errorf("expected exactly 1 call, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// Different functions: last function wins
// ---------------------------------------------------------------------------

func TestDifferentFunctions_LastWins(t *testing.T) {
	debounced := New(debounceAfter)

	var firstCalled, lastCalled int64

	debounced(func() {
		atomic.AddInt64(&firstCalled, 1)
	})
	time.Sleep(5 * time.Millisecond)
	debounced(func() {
		atomic.AddInt64(&lastCalled, 1)
	})

	time.Sleep(waitFor)

	if got := atomic.LoadInt64(&firstCalled); got != 0 {
		t.Errorf("first function should not have been called, got %d", got)
	}
	if got := atomic.LoadInt64(&lastCalled); got != 1 {
		t.Errorf("last function should have been called once, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// Generation counter: stale callbacks are discarded
//
// Hold the mutex while a callback carrying the current generation starts,
// advance the generation, then release it. The callback must observe that it
// is stale without relying on timer or scheduler timing.
// ---------------------------------------------------------------------------

func TestGenerationCounter_StaleCallbacksDiscarded(t *testing.T) {
	d := &debouncer{generation: 1}

	var fn1Called, fn2Called int64
	started := make(chan struct{})
	done := make(chan struct{})
	d.mu.Lock()
	go func() {
		close(started)
		d.fire(1, func() { atomic.AddInt64(&fn1Called, 1) })
		close(done)
	}()
	<-started
	d.generation++
	d.mu.Unlock()
	<-done

	d.fire(2, func() {
		atomic.AddInt64(&fn2Called, 1)
	})
	if got := atomic.LoadInt64(&fn1Called); got != 0 {
		t.Errorf("stale fn1 should not have been called, got %d", got)
	}
	if got := atomic.LoadInt64(&fn2Called); got != 1 {
		t.Errorf("fn2 should have been called once, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// Zero duration: callback fires (almost) immediately
// ---------------------------------------------------------------------------

func TestZeroDuration(t *testing.T) {
	debounced := New(0)

	var count int64
	fired := make(chan struct{})
	debounced(func() {
		atomic.AddInt64(&count, 1)
		close(fired)
	})

	// Wait for the callback rather than for a clock. Even a zero-delay timer
	// still goes through the scheduler, and giving it a fixed 10ms was the same
	// bet that made the stale-callback test flaky on CI.
	select {
	case <-fired:
	case <-time.After(5 * time.Second):
		t.Fatal("callback was never called with zero duration")
	}

	if got := atomic.LoadInt64(&count); got != 1 {
		t.Errorf("expected 1 call with zero duration, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// Timer re-use: verify that calling add() when timer is not nil hits Stop()
// ---------------------------------------------------------------------------

func TestTimerStop_CalledOnSubsequentAdd(t *testing.T) {
	debounced := New(debounceAfter)

	var count int64
	fn := func() { atomic.AddInt64(&count, 1) }

	// First call: timer is nil, so the nil-check branch is skipped.
	debounced(fn)
	// Second call: timer is non-nil, Stop() is called.
	debounced(fn)

	time.Sleep(waitFor)

	if got := atomic.LoadInt64(&count); got != 1 {
		t.Errorf("expected 1 call, got %d", got)
	}
}
