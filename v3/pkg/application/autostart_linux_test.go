//go:build linux && !android && !server

package application

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newLinuxAutostartForTest(t *testing.T, name string) *linuxAutostart {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	app := &App{options: Options{Name: name}}
	return &linuxAutostart{app: app}
}

func TestLinuxAutostartRoundTrip(t *testing.T) {
	a := newLinuxAutostartForTest(t, "Test App")

	st, _ := a.status()
	if st.Enabled {
		t.Fatalf("expected disabled before enable, got %+v", st)
	}

	if err := a.enable(AutostartOptions{Arguments: []string{"--hidden"}}); err != nil {
		t.Fatalf("enable: %v", err)
	}

	st, err := a.status()
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !st.Enabled {
		t.Fatalf("expected enabled, got %+v", st)
	}
	if st.Strategy != AutostartStrategyXDGAutostart {
		t.Errorf("strategy=%q", st.Strategy)
	}
	if filepath.Base(st.Path) != "test-app.desktop" {
		t.Errorf("path=%q", st.Path)
	}

	data, _ := os.ReadFile(st.Path)
	body := string(data)
	for _, want := range []string{
		"[Desktop Entry]",
		"Type=Application",
		"Name=Test App",
		"Hidden=false",
		"X-GNOME-Autostart-enabled=true",
		"--hidden",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q in desktop file:\n%s", want, body)
		}
	}

	if err := a.disable(); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if st, _ := a.status(); st.Enabled {
		t.Error("still enabled after disable")
	}
}

func TestLinuxAutostartCustomIdentifier(t *testing.T) {
	a := newLinuxAutostartForTest(t, "Test App")
	if err := a.enable(AutostartOptions{Identifier: "my-custom-id"}); err != nil {
		t.Fatalf("enable: %v", err)
	}
	st, _ := a.status()
	if filepath.Base(st.Path) != "my-custom-id.desktop" {
		t.Errorf("path=%q want my-custom-id.desktop", st.Path)
	}
}

func TestLinuxQuoteExec(t *testing.T) {
	cases := map[string]string{
		"/usr/bin/foo":            "/usr/bin/foo",
		"/path with spaces/foo":   `"/path with spaces/foo"`,
		`/has"quote`:              `"/has\"quote"`,
		`/has\back`:               `"/has\\back"`,
	}
	for in, want := range cases {
		if got := quoteExec(in); got != want {
			t.Errorf("quoteExec(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDesktopExecPath(t *testing.T) {
	cases := map[string]string{
		"Exec=/usr/bin/foo\n":                   "/usr/bin/foo",
		"Exec=/usr/bin/foo --flag\n":            "/usr/bin/foo",
		`Exec="/path with spaces/foo" --flag` + "\n":          "/path with spaces/foo",
		`[Desktop Entry]` + "\nExec=/x/y\n":     "/x/y",
		"NoExec=/x\n":                          "",
	}
	for in, want := range cases {
		if got := desktopExecPath(in); got != want {
			t.Errorf("desktopExecPath(%q) = %q, want %q", in, got, want)
		}
	}
}
