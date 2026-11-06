//go:build !windows && !darwin

package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	doctorng "github.com/wailsapp/wails/v3/pkg/doctor-ng"
)

func TestCopyReportNotice(t *testing.T) {
	for _, exit := range []int{0, 7} {
		t.Run(fmt.Sprintf("exit_%d", exit), func(t *testing.T) {
			directory := t.TempDir()
			content := filepath.Join(directory, "report")
			t.Setenv("PATH", directory)
			t.Setenv("WAYLAND_DISPLAY", "test")
			t.Setenv("COPY_CONTENT", content)
			script := fmt.Sprintf("#!/bin/sh\n/bin/cat > \"$COPY_CONTENT\"\nexit %d\n", exit)
			if err := os.WriteFile(filepath.Join(directory, "wl-copy"), []byte(script), 0755); err != nil {
				t.Fatal(err)
			}
			model := Model{state: stateReport, report: &doctorng.Report{}}
			updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			if updated.(Model).copiedNotice != (exit == 0) {
				t.Fatal("copy notice does not reflect command success")
			}
			got, err := os.ReadFile(content)
			if err != nil || string(got) != model.generateClipboardText() {
				t.Fatalf("copied report = %q, err=%v", got, err)
			}
		})
	}
}
