//go:build windows && !server

package application

import (
	"testing"

	"golang.org/x/sys/windows/registry"
)

// newWindowsAutostartForTest returns a windowsAutostart that reads and writes
// under HKCU\Software\Wails\Tests\Autostart-<test-id>. The test creates the
// key, runs, and deletes it on cleanup so production HKCU\…\Run is never
// touched.
func newWindowsAutostartForTest(t *testing.T, name string) *windowsAutostart {
	t.Helper()
	// Per-test subkey to allow parallel runs without collisions.
	subKey := `Software\Wails\Tests\Autostart-` + t.Name()
	// Pre-create so OpenKey succeeds on first read; the impl's CreateKey will
	// be a no-op when the key already exists.
	k, _, err := registry.CreateKey(registry.CURRENT_USER, subKey, registry.SET_VALUE)
	if err != nil {
		t.Fatalf("create test subkey: %v", err)
	}
	_ = k.Close()
	t.Cleanup(func() {
		_ = registry.DeleteKey(registry.CURRENT_USER, subKey)
	})

	app := &App{options: Options{Name: name}}
	return &windowsAutostart{
		app:            app,
		registrySubKey: subKey,
	}
}

func TestWindowsAutostartRoundTrip(t *testing.T) {
	a := newWindowsAutostartForTest(t, "Test App")

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
	if st.Strategy != AutostartStrategyRegistryRun {
		t.Errorf("strategy=%q want %q", st.Strategy, AutostartStrategyRegistryRun)
	}

	if err := a.disable(); err != nil {
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

func TestWindowsAutostartCustomIdentifier(t *testing.T) {
	a := newWindowsAutostartForTest(t, "Test App")
	if err := a.enable(AutostartOptions{
		Identifier: "com.example.foo",
		Arguments:  []string{"--hidden"},
	}); err != nil {
		t.Fatalf("enable: %v", err)
	}

	// Read the registry value directly to confirm both the value name and
	// the quoted command are what we expect.
	k, err := registry.OpenKey(registry.CURRENT_USER, a.registrySubKey, registry.QUERY_VALUE)
	if err != nil {
		t.Fatalf("open key: %v", err)
	}
	defer k.Close()
	val, _, err := k.GetStringValue("com.example.foo")
	if err != nil {
		t.Fatalf("read value: %v", err)
	}
	// The value should contain --hidden as a separate token after the exe.
	if want := " --hidden"; len(val) < len(want) || val[len(val)-len(want):] != want {
		t.Errorf("registry value=%q, expected to end with %q", val, want)
	}
}

func TestWindowsAutostartIdentifierValidation(t *testing.T) {
	a := newWindowsAutostartForTest(t, "Test App")
	if err := a.enable(AutostartOptions{Identifier: "bad/identifier"}); err == nil {
		t.Error("expected error for invalid identifier")
	}
}

func TestWindowsAutostartDisableNoOp(t *testing.T) {
	a := newWindowsAutostartForTest(t, "Test App")
	if err := a.disable(); err != nil {
		t.Errorf("disable when not enabled should be nil, got %v", err)
	}
}

func TestParseWindowsCommandExe(t *testing.T) {
	cases := map[string]string{
		`C:\app\foo.exe`:                  `c:\app\foo.exe`,
		`"C:\Program Files\foo\foo.exe"`:  `c:\program files\foo\foo.exe`,
		`C:\app\foo.exe --flag`:           `c:\app\foo.exe`,
		`"C:\app\foo.exe" --flag`:         `c:\app\foo.exe`,
		``:                                ``,
		`  C:\app\foo.exe`:                `c:\app\foo.exe`,
	}
	for in, want := range cases {
		if got := parseWindowsCommandExe(in); got != want {
			t.Errorf("parseWindowsCommandExe(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestQuoteWindowsArg(t *testing.T) {
	cases := map[string]string{
		`foo`:               `foo`,
		``:                  `""`,
		`hello world`:       `"hello world"`,
		`C:\Program Files`:  `"C:\Program Files"`,
		`a"b`:               `"a\"b"`,
		`a\b`:               `a\b`,
		`a\"b`:              `"a\\\"b"`,
		`a\\"b`:             `"a\\\\\"b"`,
		`trailing\`:         `trailing\`,
		`with space\`:       `"with space\\"`,
	}
	for in, want := range cases {
		if got := quoteWindowsArg(in); got != want {
			t.Errorf("quoteWindowsArg(%q) = %q, want %q", in, got, want)
		}
	}
}
