package application

import (
	"os"
	"strings"
	"testing"
)

func TestLinuxPuregoOpenFileDialogUsesComputedAction(t *testing.T) {
	data, err := os.ReadFile("linux_purego.go")
	if err != nil {
		t.Skip("linux_purego.go not available")
	}
	content := string(data)

	lines := strings.Split(content, "\n")
	inRunOpenFileDialog := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "func runOpenFileDialog") {
			inRunOpenFileDialog = true
		}
		if inRunOpenFileDialog && strings.HasPrefix(trimmed, "func ") && !strings.Contains(trimmed, "runOpenFileDialog") {
			break
		}
		if inRunOpenFileDialog && strings.Contains(trimmed, "GtkFileChooserActionOpen") && strings.Contains(trimmed, "runChooserDialog") {
			t.Errorf("runOpenFileDialog should pass computed 'action' variable to runChooserDialog, not hardcoded GtkFileChooserActionOpen")
		}
	}
}
