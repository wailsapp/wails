//go:build linux && cgo && !android && !server

package application

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestX11GlobalShortcutEndToEnd exercises the real X11 backend: it grabs a
// shortcut on the root window and uses xdotool to synthesize the key press,
// verifying the callback fires. It self-skips unless an X display and xdotool
// are available, so it is a no-op in headless CI. Run it under Xvfb:
//
//	Xvfb :99 & DISPLAY=:99 go test -tags gtk3 -run TestX11GlobalShortcutEndToEnd ./pkg/application/
func TestX11GlobalShortcutEndToEnd(t *testing.T) {
	if os.Getenv("DISPLAY") == "" {
		t.Skip("no X display; skipping X11 end-to-end test")
	}
	xdotool, err := exec.LookPath("xdotool")
	if err != nil {
		t.Skip("xdotool not available; skipping X11 end-to-end test")
	}

	m := newGlobalShortcutManager(&App{})
	fired := make(chan struct{}, 1)
	parsed, err := parseAccelerator("Ctrl+Shift+A")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	m.byID[1] = &globalShortcut{
		id:          1,
		accelerator: parsed.String(),
		parsed:      parsed,
		callback:    func() { fired <- struct{}{} },
	}

	impl := newX11GlobalShortcuts(m)
	if x, ok := impl.(*x11GlobalShortcuts); ok && x.startErr != nil {
		t.Skipf("cannot open X display: %v", x.startErr)
	}
	defer impl.unregisterAll()

	if err := impl.register(1, parsed); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Give the grab a moment to settle, then synthesize the key.
	time.Sleep(200 * time.Millisecond)
	if out, err := exec.Command(xdotool, "key", "--clearmodifiers", "ctrl+shift+a").CombinedOutput(); err != nil {
		t.Fatalf("xdotool failed: %v (%s)", err, out)
	}

	select {
	case <-fired:
		// success
	case <-time.After(3 * time.Second):
		t.Fatal("callback did not fire after synthesized key press")
	}

	// Verify unregister stops delivery.
	if err := impl.unregister(1); err != nil {
		t.Fatalf("unregister failed: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	_ = exec.Command(xdotool, "key", "--clearmodifiers", "ctrl+shift+a").Run()
	select {
	case <-fired:
		t.Fatal("callback fired after unregister")
	case <-time.After(1 * time.Second):
		// expected: no fire
	}
}
