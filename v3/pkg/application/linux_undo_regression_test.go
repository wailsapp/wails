package application

// Regression tests for https://github.com/wailsapp/wails/issues/4723
// Ctrl+Z (undo) not working in HTML <input> elements on Linux (webkit2gtk).

import (
	"os"
	"strings"
	"testing"
)

// TestLinuxCGOUndoIsImplemented verifies that undo() in the GTK3 CGO path
// calls document.execCommand('undo') rather than being a no-op.
// Prior to the fix the body was an empty commented-out stub.
func TestLinuxCGOUndoIsImplemented(t *testing.T) {
	data, err := os.ReadFile("linux_cgo_gtk3.go")
	if err != nil {
		t.Skip("linux_cgo_gtk3.go not available")
	}
	content := string(data)
	lines := strings.Split(content, "\n")

	inUndo := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "func (w *linuxWebviewWindow) undo()") {
			inUndo = true
			continue
		}
		if inUndo {
			if strings.HasPrefix(trimmed, "func ") {
				break
			}
			// Skip comment lines — we want live code, not dead comments
			if strings.HasPrefix(trimmed, "//") {
				continue
			}
			if strings.Contains(trimmed, "execJS") && strings.Contains(trimmed, "execCommand") && strings.Contains(trimmed, "undo") {
				return // pass
			}
		}
	}
	t.Error("linuxWebviewWindow.undo() must call execJS(\"document.execCommand('undo')\"); found no-op or commented implementation")
}

// TestLinuxCGORedoIsImplemented verifies that redo() in the GTK3 CGO path
// calls document.execCommand('redo').
func TestLinuxCGORedoIsImplemented(t *testing.T) {
	data, err := os.ReadFile("linux_cgo_gtk3.go")
	if err != nil {
		t.Skip("linux_cgo_gtk3.go not available")
	}
	content := string(data)
	lines := strings.Split(content, "\n")

	inRedo := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "func (w *linuxWebviewWindow) redo()") {
			inRedo = true
			continue
		}
		if inRedo {
			if strings.HasPrefix(trimmed, "func ") {
				break
			}
			// Skip comment lines — we want live code, not dead comments
			if strings.HasPrefix(trimmed, "//") {
				continue
			}
			if strings.Contains(trimmed, "execJS") && strings.Contains(trimmed, "execCommand") && strings.Contains(trimmed, "redo") {
				return // pass
			}
		}
	}
	t.Error("linuxWebviewWindow.redo() must call execJS(\"document.execCommand('redo')\"); found no-op implementation")
}

// TestLinuxOnKeyPressEventConsumesCtrlZ verifies that onKeyPressEvent returns 1
// (consumed) for Ctrl+Z so that webkit2gtk's unreliable native undo handling
// for <input> elements does not race with Wails' explicit execCommand call.
//
// The test checks structural proximity: the return C.gboolean(1) must appear
// within 3 lines of the guard condition (i.e. inside the same if-block, not
// coincidentally later in the function).
func TestLinuxOnKeyPressEventConsumesCtrlZ(t *testing.T) {
	data, err := os.ReadFile("linux_cgo_gtk3.go")
	if err != nil {
		t.Skip("linux_cgo_gtk3.go not available")
	}
	lines := strings.Split(string(data), "\n")

	inFn := false
	guardLineIdx := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "func onKeyPressEvent(") {
			inFn = true
			continue
		}
		if !inFn {
			continue
		}
		if strings.HasPrefix(trimmed, "func ") {
			break
		}
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		// Guard: a non-comment line that checks both accelerators with ||
		if strings.Contains(trimmed, `"Ctrl+Z"`) &&
			strings.Contains(trimmed, `"Ctrl+Shift+Z"`) &&
			strings.Contains(trimmed, "||") {
			guardLineIdx = i
			continue
		}
		// Return must be within the same if-block (≤3 lines after the guard).
		if guardLineIdx >= 0 && i-guardLineIdx <= 3 {
			if strings.Contains(trimmed, "return C.gboolean(1)") {
				return // pass
			}
		} else if guardLineIdx >= 0 && i-guardLineIdx > 3 {
			guardLineIdx = -1 // guard too far away — reset
		}
	}
	t.Error("onKeyPressEvent must return C.gboolean(1) immediately inside the Ctrl+Z||Ctrl+Shift+Z if-block to prevent webkit2gtk's broken native undo from racing with document.execCommand")
}

// TestLinuxHandleKeyEventFallbackForCtrlZ verifies that handleKeyEvent in the
// shared Linux file calls undo() when no binding is registered for Ctrl+Z.
func TestLinuxHandleKeyEventFallbackForCtrlZ(t *testing.T) {
	data, err := os.ReadFile("webview_window_linux.go")
	if err != nil {
		t.Skip("webview_window_linux.go not available")
	}
	content := string(data)
	lines := strings.Split(content, "\n")

	inHandleKeyEvent := false
	hasCtrlZCase := false
	hasUndoCall := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "func (w *linuxWebviewWindow) handleKeyEvent(") {
			inHandleKeyEvent = true
			continue
		}
		if inHandleKeyEvent {
			if strings.HasPrefix(trimmed, "func ") {
				break
			}
			if strings.HasPrefix(trimmed, "//") {
				continue
			}
			if strings.Contains(trimmed, `"Ctrl+Z"`) {
				hasCtrlZCase = true
			}
			if hasCtrlZCase && strings.Contains(trimmed, "w.undo()") {
				hasUndoCall = true
			}
		}
	}

	if !hasCtrlZCase || !hasUndoCall {
		t.Error("handleKeyEvent must have a fallback case for \"Ctrl+Z\" that calls w.undo() when no binding is registered")
	}
}
