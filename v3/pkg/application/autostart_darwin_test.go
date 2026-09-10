//go:build darwin && !ios && !server

package application

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newDarwinAutostartForTest builds an impl rooted at a temp HOME with a known
// Options.Name. The returned impl always uses the LaunchAgent fallback path
// (running under `go test` we're not in a .app bundle, so SMAppService is
// skipped). Also stubs out launchctlBootstrap/Bootout — otherwise a successful
// bootstrap would respawn the test binary via launchd (RunAtLoad=true) and
// recurse.
func newDarwinAutostartForTest(t *testing.T, name string) *darwinAutostart {
	t.Helper()
	t.Setenv("HOME", t.TempDir())

	origBootstrap := launchctlBootstrap
	origBootout := launchctlBootout
	launchctlBootstrap = func(string) error { return nil }
	launchctlBootout = func(string) error { return nil }
	t.Cleanup(func() {
		launchctlBootstrap = origBootstrap
		launchctlBootout = origBootout
	})

	app := &App{options: Options{Name: name}}
	return &darwinAutostart{app: app}
}

func TestDarwinAutostartRoundTrip(t *testing.T) {
	a := newDarwinAutostartForTest(t, "Test App")

	st, err := a.status()
	if err != nil {
		t.Fatalf("status before enable: %v", err)
	}
	if st.Enabled {
		t.Fatalf("expected disabled before enable, got %+v", st)
	}

	if err := a.enable(AutostartOptions{}); err != nil {
		t.Fatalf("enable: %v", err)
	}

	st, err = a.status()
	if err != nil {
		t.Fatalf("status after enable: %v", err)
	}
	if !st.Enabled {
		t.Fatalf("expected enabled, got %+v", st)
	}
	if st.Strategy != AutostartStrategyLaunchAgent {
		t.Errorf("strategy=%q want %q", st.Strategy, AutostartStrategyLaunchAgent)
	}
	if !strings.HasSuffix(st.Path, ".plist") {
		t.Errorf("path should end in .plist: %q", st.Path)
	}

	data, err := os.ReadFile(st.Path)
	if err != nil {
		t.Fatalf("read plist: %v", err)
	}
	body := string(data)
	for _, want := range []string{
		"<key>Label</key>",
		"<key>ProgramArguments</key>",
		"<key>RunAtLoad</key>",
		"<true/>",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("plist missing %q\n---\n%s", want, body)
		}
	}

	if err := a.disable(AutostartOptions{}); err != nil {
		t.Fatalf("disable: %v", err)
	}

	st, err = a.status()
	if err != nil {
		t.Fatalf("status after disable: %v", err)
	}
	if st.Enabled {
		t.Fatalf("expected disabled after disable, got %+v", st)
	}
}

func TestDarwinAutostartCustomIdentifier(t *testing.T) {
	a := newDarwinAutostartForTest(t, "Test App")
	if err := a.enable(AutostartOptions{Identifier: "com.example.foo", Arguments: []string{"--hidden"}}); err != nil {
		t.Fatalf("enable: %v", err)
	}
	st, _ := a.status()
	if !st.Enabled {
		t.Fatal("expected enabled")
	}
	if base := filepath.Base(st.Path); base != "com.example.foo.plist" {
		t.Errorf("plist filename=%q want com.example.foo.plist", base)
	}
	data, _ := os.ReadFile(st.Path)
	body := string(data)
	if !strings.Contains(body, "<string>com.example.foo</string>") {
		t.Errorf("plist missing Label com.example.foo:\n%s", body)
	}
	if !strings.Contains(body, "<string>--hidden</string>") {
		t.Errorf("plist missing --hidden argument:\n%s", body)
	}
}

func TestDarwinAutostartIdentifierValidation(t *testing.T) {
	a := newDarwinAutostartForTest(t, "Test App")
	err := a.enable(AutostartOptions{Identifier: "bad/identifier"})
	if err == nil {
		t.Error("expected error for invalid identifier")
	}
}

func TestDarwinAutostartDisableNoOp(t *testing.T) {
	a := newDarwinAutostartForTest(t, "Test App")
	if err := a.disable(AutostartOptions{}); err != nil {
		t.Errorf("disable when not enabled should be nil, got %v", err)
	}
}

func TestPlistFirstProgramArg(t *testing.T) {
	body, err := launchAgentPlist("com.example.foo", "/path/to/exe", []string{"--flag"})
	if err != nil {
		t.Fatal(err)
	}
	if got := plistFirstProgramArg(body); got != "/path/to/exe" {
		t.Errorf("plistFirstProgramArg = %q, want /path/to/exe", got)
	}
	// Malformed input must return empty, not panic.
	if got := plistFirstProgramArg([]byte("not a plist")); got != "" {
		t.Errorf("malformed plist returned %q", got)
	}
}

func TestRunningFromAppBundle(t *testing.T) {
	// `go test` runs from a tempdir-built binary, not an .app bundle.
	if runningFromAppBundle() {
		t.Skip("test binary surprisingly looks like it's inside a .app — skipping")
	}
}

func TestDarwinAutostartDisableWithOptions(t *testing.T) {
	a := newDarwinAutostartForTest(t, "Test App")
	dir, err := a.launchAgentsDir()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(dir, "com.example.stale.plist")
	plist, err := launchAgentPlist("com.example.stale", "/nonexistent/old/path", []string{"--hidden"})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stale, plist, 0644); err != nil {
		t.Fatal(err)
	}

	// Discovery-based disable cannot see the stale entry.
	if err := a.disable(AutostartOptions{}); err != nil {
		t.Fatalf("discovery disable: %v", err)
	}
	if _, err := os.Stat(stale); err != nil {
		t.Fatalf("stale plist should still exist after discovery disable: %v", err)
	}

	if err := a.disable(AutostartOptions{Identifier: "com.example.stale"}); err != nil {
		t.Fatalf("disable with identifier: %v", err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatal("stale plist should be gone")
	}
	if err := a.disable(AutostartOptions{Identifier: "com.example.stale"}); err != nil {
		t.Fatalf("second disable should be nil, got %v", err)
	}
	if err := a.disable(AutostartOptions{Identifier: "bad/identifier"}); err == nil {
		t.Error("expected error for invalid identifier")
	}
}
