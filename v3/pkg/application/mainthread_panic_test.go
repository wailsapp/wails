package application

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// dispatchWait bounds how long a test will wait for the caller of an
// InvokeSync* variant to come back. The bug under test wedges it forever, so
// the only alternative to a bound is a package-wide timeout that fails every
// other test in the run as collateral.
const dispatchWait = 5 * time.Second

// panicProbe records what the panic handler saw. The default handler ends in
// os.Exit(1), so a test that lets a panic reach it takes the whole binary with
// it: every test here has to install one of these.
type panicProbe struct {
	mu       sync.Mutex
	handled  int
	released chan struct{} // closed by the caller once its InvokeSync* returns

	// earlyRelease is set when the caller was already released at the moment
	// the panic reached the handler.
	earlyRelease bool
	// observeRelease makes the handler wait for the caller instead of only
	// glancing at it. Only the ordering test needs that.
	observeRelease bool
}

func (p *panicProbe) handle(*PanicDetails) {
	early := false
	if p.observeRelease {
		// A correctly ordered closure cannot release the caller from in here:
		// wg.Done() is deferred first, so it runs only after this returns, and
		// the wait always expires. A reversed closure releases the caller
		// before this runs, so the wait fires instead. Holding here is also
		// what keeps handled at 0 until the caller has read it.
		select {
		case <-p.released:
			early = true
		case <-time.After(250 * time.Millisecond):
		}
	} else {
		select {
		case <-p.released:
			early = true
		default:
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.handled++
	p.earlyRelease = p.earlyRelease || early
}

func (p *panicProbe) counts() (handled int, earlyRelease bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.handled, p.earlyRelease
}

// newPanicProbe installs a controllable platform layer plus a recording panic
// handler, and hands back the restore func for globalApplication.
func newPanicProbe(t *testing.T) (*threadProbeApp, *panicProbe, func()) {
	t.Helper()

	dispatcher := &threadProbeApp{}
	probe := &panicProbe{released: make(chan struct{})}

	prev := globalApplication
	globalApplication = &App{
		impl:    dispatcher,
		options: Options{PanicHandler: probe.handle},
	}

	// Off the main thread by default, so the work is queued rather than run
	// inline and the caller genuinely parks in wg.Wait().
	dispatcher.onMain.Store(false)

	return dispatcher, probe, func() { globalApplication = prev }
}

// runQueuedWork releases the UI thread once the caller has reached it, and
// fails rather than hangs if the work never arrives.
func runQueuedWork(t *testing.T, dispatcher *threadProbeApp) {
	t.Helper()

	deadline := time.Now().Add(dispatchWait)
	for dispatcher.pendingCount() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("the dispatched work never reached the main-thread queue")
		}
		time.Sleep(time.Millisecond)
	}
	dispatcher.runPending()
}

// waitForCaller fails instead of hanging when the caller is never released,
// which is the whole symptom being tested. Every caller here closes released
// from a defer, so returning from this also means the calling goroutine has
// finished and whatever it assigned is safe to read.
func waitForCaller(t *testing.T, probe *panicProbe, name string) {
	t.Helper()

	select {
	case <-probe.released:
	case <-time.After(dispatchWait):
		t.Fatalf("%s never returned", name)
	}
}

// A recovered panic in the callback must release the caller of every
// InvokeSync* variant. The closures used to call wg.Done() as their last
// statement, below `defer handlePanic()`: handlePanic recovers and returns
// normally, so the panic is swallowed, the closure unwinds cleanly, wg.Done()
// is skipped and wg.Wait() blocks for the life of the process.
//
// The variants that assign a result hand back the zero value in that case.
// That is a deliberate consequence of releasing the caller at all: before this
// there was no value to observe because there was no return.
func TestInvokeSyncReleasesCallerWhenCallbackPanics(t *testing.T) {
	tests := []struct {
		name string
		// call reports a wrong result rather than asserting, so nothing in the
		// goroutine touches *testing.T: on unfixed code that goroutine is still
		// parked when the test gives up on it.
		call func() error
	}{
		{
			name: "InvokeSync",
			call: func() error {
				InvokeSync(func() { panic("boom") })
				return nil
			},
		},
		{
			name: "InvokeSyncWithResult",
			call: func() error {
				if got := InvokeSyncWithResult(func() string { panic("boom") }); got != "" {
					return fmt.Errorf("result = %q, want the zero value", got)
				}
				return nil
			},
		},
		{
			name: "InvokeSyncWithError",
			call: func() error {
				if err := InvokeSyncWithError(func() error { panic("boom") }); err != nil {
					return fmt.Errorf("error = %v, want the zero value", err)
				}
				return nil
			},
		},
		{
			name: "InvokeSyncWithResultAndError",
			call: func() error {
				got, err := InvokeSyncWithResultAndError(func() (string, error) { panic("boom") })
				if got != "" || err != nil {
					return fmt.Errorf("result = %q, error = %v, want the zero values", got, err)
				}
				return nil
			},
		},
		{
			name: "InvokeSyncWithResultAndOther",
			call: func() error {
				got, other := InvokeSyncWithResultAndOther(func() (string, int) { panic("boom") })
				if got != "" || other != 0 {
					return fmt.Errorf("result = %q, other = %d, want the zero values", got, other)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher, probe, restore := newPanicProbe(t)
			defer restore()

			// The call has to run somewhere other than here: on unfixed code it
			// never returns, and this goroutine is the one that has to release
			// the UI thread and then declare the failure.
			var callErr error
			go func() {
				defer close(probe.released)
				callErr = tt.call()
			}()

			runQueuedWork(t, dispatcher)
			waitForCaller(t, probe, tt.name)

			if callErr != nil {
				t.Error(callErr)
			}
			if handled, _ := probe.counts(); handled != 1 {
				t.Errorf("panic handler ran %d times, want 1: the panic must still be recovered and reported", handled)
			}
		})
	}
}

// When the caller is already on the main thread, dispatchOnMainThread runs the
// closure inline — so the callback and wg.Wait() sit on the same goroutine and
// an unreleased WaitGroup freezes the main thread itself, with nobody left to
// release it. That is the shape behind a panic handler that can no longer draw
// its dialog.
func TestInvokeSyncReleasesInlineCallerWhenCallbackPanics(t *testing.T) {
	dispatcher, probe, restore := newPanicProbe(t)
	defer restore()
	dispatcher.onMain.Store(true)

	go func() {
		defer close(probe.released)
		InvokeSync(func() { panic("boom") })
	}()

	// Nothing to drain: the closure runs on the calling goroutine.
	waitForCaller(t, probe, "InvokeSync")

	if handled, _ := probe.counts(); handled != 1 {
		t.Errorf("panic handler ran %d times, want 1", handled)
	}
	if queued := dispatcher.pendingCount(); queued != 0 {
		t.Errorf("%d items were queued; the inline path was not the one exercised", queued)
	}
}

// The caller must be released only after the panic has been processed.
//
// Deferred calls run LIFO, so `defer wg.Done()` has to come first and
// `defer handlePanic()` second. Reversing them also releases the caller, and
// also passes the test above — but it releases it while the panic is still
// being handled, which is what a panic handler that wants to draw a dialog on
// the main thread races against.
func TestInvokeSyncReleasesCallerAfterPanicIsProcessed(t *testing.T) {
	dispatcher, probe, restore := newPanicProbe(t)
	defer restore()
	probe.observeRelease = true

	// Read from the caller's own goroutine the instant it is released. With
	// the release deferred first, the handler has finished by then and the
	// wg.Done()/wg.Wait() pair publishes that to this goroutine, so a correct
	// closure always reports 1 here — no timing assumption involved. The
	// reversed closure releases the caller before the handler has even been
	// entered, and reports 0.
	var handledAtRelease int

	go func() {
		defer close(probe.released)
		InvokeSync(func() { panic("boom") })
		handledAtRelease, _ = probe.counts()
	}()

	runQueuedWork(t, dispatcher)
	waitForCaller(t, probe, "InvokeSync")

	handled, earlyRelease := probe.counts()
	if handled != 1 {
		t.Fatalf("panic handler ran %d times, want 1", handled)
	}
	if handledAtRelease != 1 {
		t.Errorf("the panic handler had not finished when the caller was released (saw %d handled): wg.Done() must be deferred before handlePanic()", handledAtRelease)
	}
	if earlyRelease {
		t.Error("the caller was released while the panic was still being processed: wg.Done() must be deferred before handlePanic()")
	}
}

// The variants must keep behaving normally when nothing panics: this is the
// control for the tests above, which would all pass against a closure that
// released the caller unconditionally and never ran the callback.
//
// One call per subtest, because runPending drains until empty and reports the
// UI thread as current while it does: a second call issued by a caller it just
// released would either be swept up by the same drain or run inline, and
// either way nothing would be left to release.
func TestInvokeSyncPassesThroughWhenCallbackDoesNotPanic(t *testing.T) {
	t.Run("InvokeSync", func(t *testing.T) {
		dispatcher, probe, restore := newPanicProbe(t)
		defer restore()

		var ran bool
		go func() {
			defer close(probe.released)
			InvokeSync(func() { ran = true })
		}()

		runQueuedWork(t, dispatcher)
		waitForCaller(t, probe, "InvokeSync")

		if !ran {
			t.Error("the callback never ran")
		}
		if handled, _ := probe.counts(); handled != 0 {
			t.Errorf("panic handler ran %d times, want 0", handled)
		}
	})

	t.Run("InvokeSyncWithResultAndError", func(t *testing.T) {
		dispatcher, probe, restore := newPanicProbe(t)
		defer restore()

		wantErr := errors.New("callback error")

		var (
			gotRes string
			gotErr error
		)
		go func() {
			defer close(probe.released)
			gotRes, gotErr = InvokeSyncWithResultAndError(func() (string, error) {
				return "value", wantErr
			})
		}()

		runQueuedWork(t, dispatcher)
		waitForCaller(t, probe, "InvokeSyncWithResultAndError")

		if gotRes != "value" || !errors.Is(gotErr, wantErr) {
			t.Errorf("result = %q, error = %v; want %q, %v", gotRes, gotErr, "value", wantErr)
		}
		if handled, _ := probe.counts(); handled != 0 {
			t.Errorf("panic handler ran %d times, want 0", handled)
		}
	})
}
