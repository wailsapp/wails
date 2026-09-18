//go:build windows
// +build windows

package edge

import (
	"sync"
	"testing"
)

// TestFocusBeforeControllerReady guards the WM_SETFOCUS startup race from
// wailsapp/wails#5446: the WndProc can call Focus() re-entrantly before
// CreateCoreWebView2ControllerCompleted has finished configuring the
// controller. Focus() must be a safe no-op until IsReady() — previously it
// dereferenced the nil/partially-initialised controller, and any COM error
// routed through errorCallback, which exits the process.
func TestFocusBeforeControllerReady(t *testing.T) {
	e := NewChromium()
	e.SetErrorCallback(func(err error) {
		t.Fatalf("errorCallback invoked during pre-init Focus: %v", err)
	})

	if e.IsReady() {
		t.Fatal("new Chromium must not report ready before controller setup")
	}

	// Must not panic, must not invoke the error callback.
	e.Focus()
}

// TestFocusBeforeControllerReadyConcurrent hammers Focus/IsReady from many
// goroutines to give the race detector a chance to object — the WndProc and
// the controller-creation callback run on different stacks in production.
func TestFocusBeforeControllerReadyConcurrent(t *testing.T) {
	e := NewChromium()
	e.SetErrorCallback(func(err error) {
		t.Errorf("errorCallback invoked during pre-init Focus: %v", err)
	})

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				e.Focus()
				_ = e.IsReady()
			}
		}()
	}
	wg.Wait()
}
