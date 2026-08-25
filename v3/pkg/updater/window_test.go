package updater_test

import (
	"context"
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
)

func TestCheckAndInstall_BuiltinWindow_OpensClosesOnSuccess(t *testing.T) {
	host := &fakeHost{}
	body := []byte("payload")
	digest := sha256.Sum256(body)
	rel := &updater.Release{
		Version:      "2.0.0",
		Notes:        "test release",
		Artifact:     updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{DigestAlgo: "sha256", Digest: digest[:]},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	if err := u.CheckAndInstall(context.Background()); err != nil {
		t.Fatalf("CheckAndInstall: %v", err)
	}
	host.mu.Lock()
	defer host.mu.Unlock()
	if len(host.openCalls) != 1 {
		t.Fatalf("expected 1 OpenWindow call, got %d", len(host.openCalls))
	}
	opts := host.openCalls[0]
	if opts.Title == "" || opts.Width == 0 || opts.Height == 0 {
		t.Errorf("default window options missing defaults: %+v", opts)
	}
	if !strings.Contains(opts.InitialHTML, "wails-updater") {
		t.Errorf("default HTML doesn't contain expected marker class")
	}
	if host.window == nil {
		t.Fatal("window was never created")
	}
	// Window stays open after install because Restart is the user's next
	// action; the framework does not auto-close on "ready."
	if host.window.closed {
		t.Errorf("window should remain open in ready state")
	}
}

// CheckAndInstall must keep the window open when the user is already on
// the latest version — closing immediately produced a flicker on every
// manual check. The user dismisses the window via the Close button (which
// fires EventUserCancel and goes through the regular closeWindow path).
func TestCheckAndInstall_BuiltinWindow_StaysOpenWhenUpToDate(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"} // returns (nil, nil)
	u := newConfigured(t, host, p)

	if err := u.CheckAndInstall(context.Background()); err != nil {
		t.Fatalf("CheckAndInstall: %v", err)
	}
	host.mu.Lock()
	w := host.window
	host.mu.Unlock()
	if w == nil {
		t.Fatal("window was never created")
	}
	if w.closed {
		t.Errorf("window should stay open in up-to-date state until user dismisses")
	}

	// User dismisses by firing the cancel event the Close button emits.
	host.Emit(updater.EventUserCancel)
	// closeWindow runs synchronously in the EventUserCancel listener so the
	// window's Close() method has been called by the time Emit returns in the
	// fakeHost's synchronous fan-out.
	host.mu.Lock()
	defer host.mu.Unlock()
	if !w.closed {
		t.Errorf("window should close after EventUserCancel; closed=%v", w.closed)
	}
}

func TestCheckAndInstall_WindowNone_NeverOpens(t *testing.T) {
	host := &fakeHost{}
	rel := &updater.Release{Version: "2.0.0", Artifact: updater.Artifact{Filename: "app.bin", Size: 1}}
	p := &fakeProvider{name: "p", rel: rel, body: []byte("x")}

	u := updater.New(host)
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
		Window:         updater.WindowNone,
	}); err != nil {
		t.Fatal(err)
	}
	if err := u.CheckAndInstall(context.Background()); err != nil {
		t.Fatalf("CheckAndInstall: %v", err)
	}
	if len(host.openCalls) != 0 {
		t.Errorf("WindowNone must never open a window, got %d calls", len(host.openCalls))
	}
}

func TestCheckAndInstall_BuiltinWindow_OverrideHTML(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"}
	u := updater.New(host)
	custom := "<html><body data-marker=\"custom\"></body></html>"
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
		Window:         &updater.BuiltinWindow{HTML: custom},
	}); err != nil {
		t.Fatal(err)
	}
	_ = u.CheckAndInstall(context.Background())
	host.mu.Lock()
	defer host.mu.Unlock()
	if len(host.openCalls) != 1 {
		t.Fatalf("expected 1 OpenWindow, got %d", len(host.openCalls))
	}
	got := host.openCalls[0].InitialHTML
	if got != custom {
		preview := got
		if len(preview) > 60 {
			preview = preview[:60]
		}
		t.Errorf("custom HTML not used; got preview: %q", preview)
	}
}

func TestCheckAndInstall_BuiltinWindow_AppendsCSS(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"}
	u := updater.New(host)
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
		Window: &updater.BuiltinWindow{
			CSS: "body { background: red; }",
		},
	}); err != nil {
		t.Fatal(err)
	}
	_ = u.CheckAndInstall(context.Background())
	host.mu.Lock()
	defer host.mu.Unlock()
	html := host.openCalls[0].InitialHTML
	if !strings.Contains(html, "wails-updater") {
		t.Errorf("default template lost when only CSS overridden")
	}
	if !strings.Contains(html, "background: red") {
		t.Errorf("CSS override not appended")
	}
	// Override styles must sit inside <head> so the resulting document is
	// still well-formed (a final <style> after </html> trips HTML linters
	// and some CSP-strict configurations).
	endHead := strings.Index(html, "</head>")
	endHTML := strings.LastIndex(html, "</html>")
	override := strings.Index(html, "background: red")
	if endHead < 0 || endHTML < 0 || override < 0 {
		t.Fatalf("missing markers: head=%d html=%d override=%d", endHead, endHTML, override)
	}
	if override > endHead {
		t.Errorf("CSS override emitted after </head>: head=%d override=%d", endHead, override)
	}
	if override > endHTML {
		t.Errorf("CSS override emitted after </html>: html=%d override=%d", endHTML, override)
	}
}

func TestCheckAndInstall_BuiltinWindow_OverrideOptions(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"}
	u := updater.New(host)
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
		Window: &updater.BuiltinWindow{
			Options: updater.WindowOptions{
				Title:     "Custom",
				Width:     900,
				Height:    600,
				Frameless: true,
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
	_ = u.CheckAndInstall(context.Background())
	host.mu.Lock()
	defer host.mu.Unlock()
	opts := host.openCalls[0]
	if opts.Title != "Custom" || opts.Width != 900 || opts.Height != 600 || !opts.Frameless {
		t.Errorf("overrides not applied: %+v", opts)
	}
}

func TestUserCancel_ClosesWindow(t *testing.T) {
	host := &fakeHost{}
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: 1},
	}
	p := &fakeProvider{name: "p", rel: rel, body: []byte("x")}

	u := updater.New(host)
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
	}); err != nil {
		t.Fatal(err)
	}

	// Drive the full Check + DownloadAndInstall flow synchronously. Listeners
	// are guaranteed registered by the time CheckAndInstall returns, and the
	// builtin window stays open in StateReady so we can observe the cancel
	// tearing it down.
	if err := u.CheckAndInstall(context.Background()); err != nil {
		t.Fatalf("CheckAndInstall: %v", err)
	}
	host.mu.Lock()
	if host.window == nil || host.window.closed {
		host.mu.Unlock()
		t.Fatal("window should be open and not yet closed after a ready install")
	}
	host.mu.Unlock()

	// Now fire cancel. fakeHost.Emit dispatches listener callbacks
	// synchronously on the calling goroutine, so the close has already
	// happened by the time Emit returns — no polling required.
	host.Emit(updater.EventUserCancel)

	host.mu.Lock()
	defer host.mu.Unlock()
	if host.window == nil || !host.window.closed {
		t.Fatal("window should have closed after user:cancel")
	}
}

// Repeated CheckAndInstall calls (e.g. periodic timer + manual click) must
// tear down the previous session before opening a new one — otherwise each
// invocation appends another 5 fresh listeners and the count grows
// unboundedly. Regression guard.
//
// After my fix to keep the window open on up-to-date, each call's listeners
// stay alive until the NEXT call's open-session tears them down. After N
// calls the count should equal exactly the per-session listener count
// (currently 5 user-action listeners), never N * 5.
func TestCheckAndInstall_RepeatedCalls_DoNotLeakListeners(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"} // returns (nil, nil) — fast up-to-date path
	u := newConfigured(t, host, p)

	for i := 0; i < 5; i++ {
		if err := u.CheckAndInstall(context.Background()); err != nil {
			t.Fatalf("CheckAndInstall #%d: %v", i, err)
		}
	}

	host.mu.Lock()
	defer host.mu.Unlock()
	// After 5 calls the only listeners that should still exist are the 5
	// from the most recent session. A leak would show as ~25 (or 5 per call
	// × 5 calls) and a regression of the close-before-open ordering would
	// likewise grow.
	for _, name := range []string{
		updater.EventUserInstall,
		updater.EventUserCancel,
		updater.EventUserSkip,
		updater.EventUserRemind,
		updater.EventUserRestart,
	} {
		if got := len(host.listeners[name]); got != 1 {
			t.Errorf("listeners for %q: got %d, want 1 (only the most recent session's set)", name, got)
		}
	}
}

// The default window template emits updater:window:ready on load so it can
// rehydrate from the current updater state when it opens. The Updater must
// reply with a snapshot event matching the current State.
func TestWindowReady_ReplaysCurrentStateSnapshot(t *testing.T) {
	host := &fakeHost{}
	rel := &updater.Release{Version: "2.0.0", Artifact: updater.Artifact{Filename: "app.bin", Size: 1}}
	p := &fakeProvider{name: "p", rel: rel, body: []byte("x")}
	u := newConfigured(t, host, p)

	// Advance to StateAvailable.
	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if u.State() != updater.StateAvailable {
		t.Fatalf("state: %s", u.State())
	}

	// Open a session (this is what CheckAndInstall does internally; we use
	// the same paths via Check above).
	_ = u.CheckAndInstall(context.Background())

	// Drop the noisy event history accumulated so far so the snapshot
	// re-emit we trigger next stands out.
	host.mu.Lock()
	host.events = nil
	host.mu.Unlock()

	// The window has just finished loading and asks for state.
	host.Emit(updater.EventWindowReady)

	host.mu.Lock()
	defer host.mu.Unlock()
	saw := map[string]bool{}
	for _, e := range host.events {
		saw[e.Name] = true
	}
	if !saw[updater.EventUpdateReady] && !saw[updater.EventUpdateAvailable] {
		t.Errorf("expected snapshot replay (update-available or update-ready); got events %v", host.events)
	}
}

func TestUserSkip_RecordsVersionAndShortCircuitsCheck(t *testing.T) {
	host := &fakeHost{}
	rel := &updater.Release{Version: "2.0.0", Artifact: updater.Artifact{Filename: "app.bin", Size: 1}}
	p := &fakeProvider{name: "p", rel: rel, body: []byte("x")}
	u := newConfigured(t, host, p)

	// First Check finds the release, marks it pending.
	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	// Skip it.
	u.SkipVersion("2.0.0")
	if u.SkippedVersion() != "2.0.0" {
		t.Errorf("SkippedVersion: %q", u.SkippedVersion())
	}
	// Subsequent Check should treat 2.0.0 as up-to-date.
	got, err := u.Check(context.Background())
	if err != nil || got != nil {
		t.Fatalf("expected (nil, nil) after skip, got %+v %v", got, err)
	}
	if u.State() != updater.StateUpToDate {
		t.Errorf("state: %s", u.State())
	}
}
