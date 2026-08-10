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

	// Step 3: give the timer a chance to fire and block on Lock. This is a nudge
	// rather than a requirement: if it has not fired yet, the generation bump
	// below still discards it when it does, so either ordering exercises the
	// same guarantee.
	time.Sleep(10 * time.Millisecond)

	// Step 4: bump the generation while holding the lock, simulating what add()
	// does, so the pending callback is stale.
	d.generation++

	// Step 5: release; the blocked timer goroutine acquires the lock, sees the
	// generation mismatch, and returns without calling fn1.
	d.mu.Unlock()

	// Now schedule fn2 normally, and wait for it rather than for a clock. The
	// previous version slept 30ms and asserted fn2 had run, which is a bet that
	// a 2ms timer gets scheduled within 30ms — false often enough on a loaded
	// CI runner to fail the build.
	fired := make(chan struct{})
	d.add(func() {
		atomic.AddInt64(&fn2Called, 1)
		close(fired)
	})

	select {
	case <-fired:
	case <-time.After(5 * time.Second):
		t.Fatal("fn2 was never called")
	}

	// fn1's timer was armed earlier and with the same delay, so by the time fn2
	// has run a stale fn1 would have run too. Checking here is deterministic
	// where a fixed sleep was not.
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
