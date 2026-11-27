//go:build windows && !server

package application

import (
	"errors"
	"testing"
	"time"
)

func TestRetryScreenCache(t *testing.T) {
	errTransient := errors.New("transient")

	t.Run("returns nil without retrying when the first attempt succeeds", func(t *testing.T) {
		calls := 0
		err := retryScreenCache(3, time.Millisecond, func() error {
			calls++
			return nil
		})
		if err != nil || calls != 1 {
			t.Fatalf("got err=%v after %d calls, want nil after 1", err, calls)
		}
	})

	t.Run("recovers from a transient failure", func(t *testing.T) {
		calls := 0
		err := retryScreenCache(3, time.Millisecond, func() error {
			calls++
			if calls < 3 {
				return errTransient
			}
			return nil
		})
		if err != nil || calls != 3 {
			t.Fatalf("got err=%v after %d calls, want nil after 3", err, calls)
		}
	})

	t.Run("returns the last error once every attempt has failed", func(t *testing.T) {
		calls := 0
		err := retryScreenCache(3, time.Millisecond, func() error {
			calls++
			return errTransient
		})
		if !errors.Is(err, errTransient) || calls != 3 {
			t.Fatalf("got err=%v after %d calls, want %v after 3", err, calls, errTransient)
		}
	})
}
