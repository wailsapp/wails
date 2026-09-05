//go:build windows && !server

package application

import (
	"os"
	"strings"
	"testing"
)

// Regression test for https://github.com/wailsapp/wails/issues/6054.
func TestModalParentIsReleasedBeforeWindowDestruction(t *testing.T) {
	sourceBytes, err := os.ReadFile("webview_window_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	source := string(sourceBytes)

	assertCallOrder(
		t,
		source,
		"func (w *windowsWebviewWindow) destroy()",
		"func (w *windowsWebviewWindow) reload()",
		"w.releaseModalParent(w.isFocused())",
		"w32.DestroyWindow(w.hwnd)",
	)
	assertCallOrder(
		t,
		source,
		"case w32.WM_CLOSE:",
		"case w32.WM_SETCURSOR:",
		"w.releaseModalParent(restoreParentActivation)",
		"w32.DefWindowProc(w.hwnd, w32.WM_CLOSE, 0, 0)",
	)
}

func TestModalParentActivationIsConditional(t *testing.T) {
	sourceBytes, err := os.ReadFile("webview_window_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	assertCallOrder(
		t,
		string(sourceBytes),
		"func (w *windowsWebviewWindow) releaseModalParent(activate bool)",
		"func (w *windowsWebviewWindow) reload()",
		"if activate {",
		"w32.SetActiveWindow(w.parentHWND)",
	)
}

func assertCallOrder(t *testing.T, source, start, end, first, second string) {
	t.Helper()
	startIndex := strings.Index(source, start)
	if startIndex < 0 {
		t.Fatalf("start marker %q not found", start)
	}
	section := source[startIndex:]
	endIndex := strings.Index(section, end)
	if endIndex < 0 {
		t.Fatalf("end marker %q not found after %q", end, start)
	}
	section = section[:endIndex]
	firstIndex := strings.Index(section, first)
	secondIndex := strings.Index(section, second)
	if firstIndex < 0 || secondIndex < 0 {
		t.Fatalf("expected %q and %q between %q and %q", first, second, start, end)
	}
	if firstIndex > secondIndex {
		t.Fatalf("%q must occur before %q between %q and %q", first, second, start, end)
	}
}
