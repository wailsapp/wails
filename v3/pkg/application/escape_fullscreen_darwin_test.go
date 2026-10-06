//go:build darwin && !ios && !server

package application

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode"
)

const escapeFullscreenHarness = "testdata/escape_fullscreen/main.m"

// buildEscapeFullscreenProbe compiles testdata/escape_fullscreen/main.m, which
// drives windowDropsEscapeCommand from webview_window_darwin.h — the decision
// WebviewWindow and WebviewPanel make for DisableEscapeExitsFullscreen.
func buildEscapeFullscreenProbe(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("xcrun"); err != nil {
		t.Skip("requires Xcode command line tools")
	}
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(t.TempDir(), "escape-fullscreen")
	macNativeCommand(t, "xcrun", "clang", "-Werror", "-I", root,
		escapeFullscreenHarness, "-framework", "Cocoa", "-framework", "WebKit", "-o", executable)
	return executable
}

func TestEscapeFullscreenCommandFilter(t *testing.T) {
	macNativeCommand(t, buildEscapeFullscreenProbe(t), "--filter")
}

// objcMethodBody returns the body of the one method in source whose
// declaration matches signature, with all whitespace removed.
func objcMethodBody(t *testing.T, file, signature string) string {
	t.Helper()
	source, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	declaration := regexp.MustCompile(`(?m)^-\s*\(void\)\s*` + regexp.QuoteMeta(signature) + `\s*\{`)
	matches := declaration.FindAllIndex(source, -1)
	if len(matches) != 1 {
		t.Fatalf("%s: found %d definitions of %s, want 1", file, len(matches), signature)
	}
	start := matches[0][1]
	depth := 1
	for i := start; i < len(source); i++ {
		switch source[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return strings.Map(func(r rune) rune {
					if unicode.IsSpace(r) {
						return -1
					}
					return r
				}, string(source[start:i]))
			}
		}
	}
	t.Fatalf("%s: %s has no closing brace", file, signature)
	return ""
}

// TestEscapeFullscreenHarnessMatchesWindowClasses keeps the overrides that
// TestEscapeKeepsFullscreen puts on its windows identical to those of
// WebviewWindow and WebviewPanel.
func TestEscapeFullscreenHarnessMatchesWindowClasses(t *testing.T) {
	for _, signature := range []string{"cancelOperation:(id)sender", "doCommandBySelector:(SEL)selector"} {
		want := objcMethodBody(t, escapeFullscreenHarness, signature)
		for _, file := range []string{"webview_window_darwin.m", "webview_panel_darwin.m"} {
			if got := objcMethodBody(t, file, signature); got != want {
				t.Errorf("%s in %s differs from %s:\n%s\nwant\n%s", signature, file, escapeFullscreenHarness, got, want)
			}
		}
	}
}

// TestEscapeKeepsFullscreen presses Escape in a focused <textarea> of a real
// fullscreen window, where WebKit's command for the key reaches the window as
// doCommandBySelector:cancelOperation: rather than as cancelOperation:. It
// needs a logged-in display session and takes over the screen while it runs,
// so it is opt-in.
func TestEscapeKeepsFullscreen(t *testing.T) {
	if os.Getenv("WAILS_TEST_FULLSCREEN") == "" {
		t.Skip("takes the screen into fullscreen; set WAILS_TEST_FULLSCREEN=1 to run")
	}
	executable := buildEscapeFullscreenProbe(t)
	for _, tc := range []struct{ name, mode string }{
		{"cancelOperation override alone leaves fullscreen", "--cancel-only"},
		{"option on stays fullscreen", "--guarded"},
		{"option off leaves fullscreen", "--allow"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			output, err := exec.CommandContext(ctx, executable, "--fullscreen", tc.mode).CombinedOutput()
			if err != nil {
				t.Fatalf("%s: %v\n%s", tc.mode, err, output)
			}
			t.Log(strings.TrimSpace(string(output)))
		})
	}
}
