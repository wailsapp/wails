package application

import (
	"testing"
	"time"
)

// recordingIconTray records the icon setter calls that reach the platform tray.
type recordingIconTray struct {
	stubTrayImpl
	icons     chan string
	darkIcons chan string
}

func newRecordingIconTray() recordingIconTray {
	return recordingIconTray{icons: make(chan string, 1), darkIcons: make(chan string, 1)}
}

func (t recordingIconTray) setIcon(icon []byte)         { t.icons <- string(icon) }
func (t recordingIconTray) setDarkModeIcon(icon []byte) { t.darkIcons <- string(icon) }

// anyThreadIconTray is a platform tray whose icon setters lock their own state,
// as the Windows tray does.
type anyThreadIconTray struct{ recordingIconTray }

func (anyThreadIconTray) iconSettersRunOnAnyThread() {}

// trayOffMainThread returns a running tray whose main thread never runs queued
// work, so any InvokeSync from the calling goroutine waits forever.
func trayOffMainThread(t *testing.T, impl systemTrayImpl) (*SystemTray, *threadProbeApp) {
	t.Helper()
	probe := &threadProbeApp{}
	prev := globalApplication
	globalApplication = &App{impl: probe}
	t.Cleanup(func() { globalApplication = prev })

	tray := newSystemTray(1)
	tray.impl = impl
	return tray, probe
}

func returnsWithin(t *testing.T, name string, call func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		call()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s waited for the main thread", name)
	}
}

// A tray whose icon setters lock their own state must be updated on the caller's
// goroutine. A Windows Shell_NotifyIcon that Explorer or a security hook never
// answers then blocks only the caller, not the main thread and every window.
func TestSystemTraySetIconSkipsMainThreadWhenImplAllowsIt(t *testing.T) {
	impl := anyThreadIconTray{newRecordingIconTray()}
	tray, probe := trayOffMainThread(t, impl)

	returnsWithin(t, "SetIcon", func() { tray.SetIcon([]byte("light")) })
	returnsWithin(t, "SetDarkModeIcon", func() { tray.SetDarkModeIcon([]byte("dark")) })

	if got := <-impl.icons; got != "light" {
		t.Errorf("setIcon got %q, want %q", got, "light")
	}
	if got := <-impl.darkIcons; got != "dark" {
		t.Errorf("setDarkModeIcon got %q, want %q", got, "dark")
	}
	probe.mu.Lock()
	queued := len(probe.pending)
	probe.mu.Unlock()
	if queued != 0 {
		t.Errorf("%d calls were queued for the main thread, want 0", queued)
	}
}

// Platforms whose tray state belongs to the main thread keep the old dispatch.
func TestSystemTraySetIconRunsOnMainThreadByDefault(t *testing.T) {
	impl := newRecordingIconTray()
	tray, probe := trayOffMainThread(t, impl)

	done := make(chan struct{})
	go func() {
		tray.SetIcon([]byte("light"))
		close(done)
	}()

	deadline := time.Now().Add(2 * time.Second)
	for {
		probe.mu.Lock()
		queued := len(probe.pending)
		probe.mu.Unlock()
		if queued == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("SetIcon did not queue its work for the main thread")
		}
		time.Sleep(time.Millisecond)
	}
	select {
	case <-impl.icons:
		t.Fatal("setIcon ran before the main thread ran the queued work")
	default:
	}

	probe.runPending()
	<-done
	if got := <-impl.icons; got != "light" {
		t.Errorf("setIcon got %q, want %q", got, "light")
	}
}
