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
	data, err := os.ReadFile("linux_cgo.go")
	if err != nil {
		t.Skip("linux_cgo.go not available")
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
			// The body must contain the execCommand call, not be empty or commented
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
	data, err := os.ReadFile("linux_cgo.go")
	if err != nil {
		t.Skip("linux_cgo.go not available")
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
func TestLinuxOnKeyPressEventConsumesCtrlZ(t *testing.T) {
	data, err := os.ReadFile("linux_cgo.go")
	if err != nil {
		t.Skip("linux_cgo.go not available")
	}
	content := string(data)

	// Must contain a return-1 branch guarded by the Ctrl+Z accelerator string
	hasCtrlZ := strings.Contains(content, `"Ctrl+Z"`)
	hasReturn1 := strings.Contains(content, "return C.gboolean(1)")
	if !hasCtrlZ || !hasReturn1 {
		t.Error("onKeyPressEvent must return C.gboolean(1) for Ctrl+Z to prevent webkit2gtk's broken native undo from racing with document.execCommand")
	}
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
