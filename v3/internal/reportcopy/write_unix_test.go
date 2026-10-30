//go:build !windows && !darwin

package reportcopy

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteDesktopReport(t *testing.T) {
	for _, tc := range []struct {
		name, wayland string
		tools         []string
		arguments     string
		exit          int
		wantError     bool
	}{
		{name: "wayland", wayland: "test", tools: []string{"wl-copy", "xclip"}, arguments: "wl-copy\n"},
		{name: "xclip", tools: []string{"xclip", "xsel"}, arguments: "xclip\n-selection\nclipboard\n-in\n"},
		{name: "xsel", tools: []string{"xsel"}, arguments: "xsel\n--clipboard\n--input\n"},
		{name: "missing", wantError: true},
		{name: "missing-wayland-tool", wayland: "test", tools: []string{"xclip"}, wantError: true},
		{name: "failed-command", tools: []string{"xclip"}, exit: 7, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			directory := t.TempDir()
			arguments, content := filepath.Join(directory, "arguments"), filepath.Join(directory, "content")
			t.Setenv("PATH", directory)
			t.Setenv("WAYLAND_DISPLAY", tc.wayland)
			t.Setenv("COPY_ARGUMENTS", arguments)
			t.Setenv("COPY_CONTENT", content)
			for _, name := range tc.tools {
				script := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' \"${0##*/}\" \"$@\" > \"$COPY_ARGUMENTS\"\n/bin/cat > \"$COPY_CONTENT\"\nexit %d\n", tc.exit)
				if err := os.WriteFile(filepath.Join(directory, name), []byte(script), 0755); err != nil {
					t.Fatal(err)
				}
			}
			text := "Wails report: 世界\nquotes '\"; $(literal)\n"
			if err := Write(text); (err != nil) != tc.wantError {
				t.Fatalf("Write() error = %v, want error=%v", err, tc.wantError)
			}
			if tc.wantError {
				return
			}
			got, err := os.ReadFile(arguments)
			if err != nil || string(got) != tc.arguments {
				t.Fatalf("arguments=%q, err=%v; want %q", got, err, tc.arguments)
			}
			got, err = os.ReadFile(content)
			if err != nil || string(got) != text {
				t.Fatalf("report text=%q, err=%v; want %q", got, err, text)
			}
		})
	}
}
