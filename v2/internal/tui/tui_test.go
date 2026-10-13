package tui

import (
	"errors"
	"runtime"
	"testing"
	"time"
)

func withNoColour(t *testing.T) {
	t.Helper()
	old := ColourEnabled
	ColourEnabled = false
	t.Cleanup(func() { ColourEnabled = old })
}

func TestWithSpinnerNoColourSuccess(t *testing.T) {
	withNoColour(t)
	called := false
	if err := WithSpinner("work", func() error { called = true; return nil }); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called {
		t.Fatal("fn was not called")
	}
}

func TestWithSpinnerNoColourError(t *testing.T) {
	withNoColour(t)
	want := errors.New("boom")
	if err := WithSpinner("work", func() error { return want }); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestWithSpinnerNoColourPanicPropagates(t *testing.T) {
	withNoColour(t)
	before := runtime.NumGoroutine()

	finished := make(chan any, 1)
	go func() {
		defer func() { finished <- recover() }()
		_ = WithSpinner("work", func() error { panic("kaboom") })
	}()

	select {
	case r := <-finished:
		if r != "kaboom" {
			t.Fatalf("expected panic value %q, got %v", "kaboom", r)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("WithSpinner deadlocked on panic")
	}

	// No goroutines should be left behind on the non-colour path.
	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if after := runtime.NumGoroutine(); after > before {
		t.Fatalf("goroutine leak: before=%d after=%d", before, after)
	}
}
