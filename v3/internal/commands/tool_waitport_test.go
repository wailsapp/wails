package commands

import (
	"testing"
	"time"
)

func TestWaitForPortRetriesUntilReady(t *testing.T) {
	attempts := 0
	ready := waitForPort(func() bool {
		attempts++
		return attempts == 3
	}, time.Second)

	if !ready {
		t.Fatal("waitForPort reported that the port was unavailable")
	}
	if attempts != 3 {
		t.Fatalf("attempt count = %d, want 3", attempts)
	}
}

func TestWaitForPortDoesNotRetryWithoutTimeout(t *testing.T) {
	attempts := 0
	ready := waitForPort(func() bool {
		attempts++
		return false
	}, 0)

	if ready {
		t.Fatal("waitForPort reported that the port was available")
	}
	if attempts != 1 {
		t.Fatalf("attempt count = %d, want 1", attempts)
	}
}
