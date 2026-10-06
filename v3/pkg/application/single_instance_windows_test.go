//go:build windows && !server

package application

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestWindowsLockReleaseFreesTheMutex(t *testing.T) {
	uniqueID := fmt.Sprintf("org.wails.test.release.p%d", os.Getpid())
	newLock := func() *windowsLock {
		return &windowsLock{manager: &singleInstanceManager{options: &SingleInstanceOptions{}}}
	}

	first := newLock()
	if err := first.acquire(uniqueID); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	second := newLock()
	if err := second.acquire(uniqueID); !errors.Is(err, alreadyRunningError) {
		t.Fatalf("while held: got %v, want alreadyRunningError", err)
	}

	first.release()

	third := newLock()
	if err := third.acquire(uniqueID); err != nil {
		t.Fatalf("after release: %v; want the mutex free", err)
	}
	third.release()
}
