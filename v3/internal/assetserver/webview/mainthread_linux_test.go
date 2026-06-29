//go:build linux && cgo && !android

package webview

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// withWatchdog runs body and fails the test if it does not finish within d. A
// deadlock in the dispatch path would otherwise hang the whole test binary; this
// turns it into a clear failure instead.
func withWatchdog(t *testing.T, d time.Duration, name string, body func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		body()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("%s did not complete within %s (deadlock?)", name, d)
	}
}

// TestMainThreadDispatch exercises invokeOnMainSync and DisableMainThreadDispatch
// against a real GLib main loop, covering the shutdown race CodeRabbit flagged on
// #5668: the enabled-check must be serialized with scheduling so a worker either
// dispatches onto the live loop or runs inline, never blocking on a dead loop.
//
// The subtests run in order and share the process-global dispatch-enabled flag:
// DisableRaceDoesNotDeadlock flips it off permanently, and InlineAfterDisable
// depends on that, so the subtests are not independently runnable via -run.
func TestMainThreadDispatch(t *testing.T) {
	// Drive a real GLib main loop on a dedicated, locked OS thread.
	loopStopped := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		testRunMainLoop()
		close(loopStopped)
	}()
	testWaitLoopRunning()

	// While dispatch is enabled, many concurrent workers must each have their
	// callback run on the loop thread, and every call must complete.
	t.Run("ConcurrentDispatchOnLoopThread", func(t *testing.T) {
		const workers, perWorker = 64, 50
		var ran int64
		var offThread int64

		withWatchdog(t, 30*time.Second, "concurrent dispatch", func() {
			var wg sync.WaitGroup
			wg.Add(workers)
			for w := 0; w < workers; w++ {
				go func() {
					defer wg.Done()
					for i := 0; i < perWorker; i++ {
						invokeOnMainSync(func() {
							if !testOnLoopThread() {
								atomic.AddInt64(&offThread, 1)
							}
							atomic.AddInt64(&ran, 1)
						})
					}
				}()
			}
			wg.Wait()
		})

		if got, want := atomic.LoadInt64(&ran), int64(workers*perWorker); got != want {
			t.Fatalf("callbacks run = %d, want %d", got, want)
		}
		if n := atomic.LoadInt64(&offThread); n != 0 {
			t.Fatalf("%d callbacks ran off the loop thread while dispatch enabled", n)
		}
	})

	// Disabling must be safe to race with in-flight workers (CodeRabbit's
	// finding) and must not deadlock. The loop stays alive here, so any source
	// scheduled just before the flag flipped still drains on it; this isolates
	// the serialization of the check + scheduling from loop teardown.
	t.Run("DisableRaceDoesNotDeadlock", func(t *testing.T) {
		const workers = 64
		withWatchdog(t, 30*time.Second, "disable race", func() {
			var wg sync.WaitGroup
			wg.Add(workers + 1)
			for w := 0; w < workers; w++ {
				go func() {
					defer wg.Done()
					for i := 0; i < 100; i++ {
						invokeOnMainSync(func() {})
					}
				}()
			}
			go func() {
				defer wg.Done()
				// Flip the flag partway through the burst.
				time.Sleep(2 * time.Millisecond)
				DisableMainThreadDispatch()
			}()
			wg.Wait()
		})
	})

	// After disable, invokeOnMainSync must run the callback inline on the caller
	// (not the loop thread) and return without scheduling onto the loop. The loop
	// is still alive, proving the inline path keys off the flag, not loop state.
	t.Run("InlineAfterDisable", func(t *testing.T) {
		var ranInline bool
		withWatchdog(t, 10*time.Second, "inline after disable", func() {
			done := make(chan struct{})
			go func() {
				invokeOnMainSync(func() {
					ranInline = !testOnLoopThread()
				})
				close(done)
			}()
			<-done
		})
		if !ranInline {
			t.Fatal("callback did not run inline on the calling thread after disable")
		}
	})

	testQuitMainLoop()
	<-loopStopped
}

// TestWebkitRequestBodyZeroLengthRead guards the io.Reader contract fix from
// #5668: a zero-length read must return (0, nil) and must not dereference &p[0]
// or the stream. The nil stream here would crash if the early return regressed,
// so no GTK runtime is needed.
func TestWebkitRequestBodyZeroLengthRead(t *testing.T) {
	b := &webkitRequestBody{}
	for _, p := range [][]byte{nil, {}, make([]byte, 0)} {
		n, err := b.Read(p)
		if n != 0 || err != nil {
			t.Fatalf("Read(len=%d) = (%d, %v), want (0, nil)", len(p), n, err)
		}
	}
}
