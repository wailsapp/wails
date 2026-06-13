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
// Strategy: use a lock to hold the debouncer's mutex while a timer is
// known to be mid-flight, then release it so the timer goroutine acquires
// it and sees a bumped generation.
//
// We use the internal debouncer directly (same package) to get precise
// control. The sequence is:
//  1. Call add(fn1) — starts a 2ms timer.
//  2. Sleep 3ms so the timer fires and its goroutine attempts to lock.
//  3. Meanwhile, hold d.mu externally so the goroutine blocks.
//  4. Call add(fn2) while holding the lock — bumps generation.
//  5. Release the lock: timer goroutine sees generation mismatch → returns.
// ---------------------------------------------------------------------------

func TestGenerationCounter_StaleCallbacksDiscarded(t *testing.T) {
	d := &debouncer{after: 2 * time.Millisecond}

	var fn1Called, fn2Called int64

	// Step 1: schedule fn1 with a 2ms timer.
	d.add(func() { atomic.AddInt64(&fn1Called, 1) })

	// Step 2: acquire the mutex BEFORE sleeping, so the timer goroutine will
	// block the instant it fires.
	d.mu.Lock()

	// Step 3: wait long enough for the timer to have fired (it will block on Lock).
	time.Sleep(10 * time.Millisecond)

	// Step 4: bump the generation by calling add directly while holding the lock.
	// add() tries to Lock too — we hold it, so we call the internals manually
	// to simulate what add() does under the lock.
	d.generation++ // stale gen ≠ new gen → the blocked goroutine will bail

	// Step 5: release; the blocked timer goroutine now acquires the lock, checks
	// generation mismatch, and returns without calling fn1.
	d.mu.Unlock()

	// Now schedule fn2 normally.
	d.add(func() { atomic.AddInt64(&fn2Called, 1) })

	time.Sleep(30 * time.Millisecond)

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
	debounced(func() {
		atomic.AddInt64(&count, 1)
	})

	time.Sleep(10 * time.Millisecond)

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
