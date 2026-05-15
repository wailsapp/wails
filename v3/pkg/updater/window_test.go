package updater_test

import (
	"context"
	"crypto/sha256"
	"strings"
	"testing"
	"time"

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

func TestCheckAndInstall_BuiltinWindow_ClosesOnUpToDate(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"} // returns (nil, nil)
	u := newConfigured(t, host, p)

	if err := u.CheckAndInstall(context.Background()); err != nil {
		t.Fatalf("CheckAndInstall: %v", err)
	}
	host.mu.Lock()
	defer host.mu.Unlock()
	if host.window == nil {
		t.Fatal("window was never created")
	}
	if !host.window.closed {
		t.Errorf("window should auto-close when up to date")
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
	if host.openCalls[0].InitialHTML != custom {
		t.Errorf("custom HTML not used; got first 60 chars: %q", host.openCalls[0].InitialHTML[:60])
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
	rel := &updater.Release{Version: "2.0.0", Artifact: updater.Artifact{Filename: "app.bin"}}
	p := &fakeProvider{name: "p", rel: rel, body: []byte("x")}

	u := updater.New(host)
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
	}); err != nil {
		t.Fatal(err)
	}

	// Open a window without running the full flow.
	go func() {
		// Slight delay so the listeners get installed before we trigger.
		time.Sleep(10 * time.Millisecond)
		host.Emit(updater.EventUserCancel)
	}()
	// We can't run the full CheckAndInstall and then trigger cancel
	// reliably; instead we open the session directly via the public API.
	_ = u.CheckAndInstall(context.Background())
	// Cancel fires after success — should now have closed the window.
	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		host.mu.Lock()
		closed := host.window != nil && host.window.closed
		host.mu.Unlock()
		if closed {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	host.mu.Lock()
	defer host.mu.Unlock()
	if host.window == nil || !host.window.closed {
		t.Fatal("window should have closed after user:cancel")
	}
}

// Repeated CheckAndInstall calls (e.g. periodic timer + manual click) must
// tear down the previous session before opening a new one — otherwise the
// stale listeners and window leak. Regression guard for the bug where each
// call appended 5 fresh listeners under sessMu without closing the prior set.
func TestCheckAndInstall_RepeatedCalls_DoNotLeakListeners(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"} // returns (nil, nil) — fast up-to-date path
	u := newConfigured(t, host, p)

	for i := 0; i < 3; i++ {
		if err := u.CheckAndInstall(context.Background()); err != nil {
			t.Fatalf("CheckAndInstall #%d: %v", i, err)
		}
	}

	host.mu.Lock()
	defer host.mu.Unlock()
	// Each call opens then closes a window — but closeWindow tears down all
	// listeners registered in that session, so after the last close the host
	// should be back to zero listeners for each user-action event.
	for _, name := range []string{
		updater.EventUserInstall,
		updater.EventUserCancel,
		updater.EventUserSkip,
		updater.EventUserRemind,
		updater.EventUserRestart,
	} {
		if got := len(host.listeners[name]); got != 0 {
			t.Errorf("listeners for %q: got %d, want 0 (session not torn down)", name, got)
		}
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
