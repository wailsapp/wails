//go:build !windows

package commands

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/github"
)

func TestUpdateCLIInstall(t *testing.T) {
	for _, tc := range []struct {
		name, exit, current, target string
		wantCode                    int
		force                       bool
	}{
		{"failure", "7", "v3.0.0", "v3.0.1", 7, true},
		{"success", "0", "v3.0.0", "v3.0.1", 0, true},
		{"stable metadata", "0", "v3.0.0+build", "v3.0.1", 0, false},
		{"stable metadata to prerelease", "0", "v3.0.0+build", "v3.0.0-beta.1", 0, false},
		{"stable to prerelease", "0", "v3.0.0", "v3.0.0-beta.1", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			logFile := filepath.Join(dir, "invocations")
			script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$WAILS_UPDATE_LOG\"\necho simulated installer output\nexit \"$WAILS_UPDATE_EXIT\"\n"
			if err := os.WriteFile(filepath.Join(dir, "go"), []byte(script), 0700); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			t.Setenv("WAILS_UPDATE_LOG", logFile)
			t.Setenv("WAILS_UPDATE_EXIT", tc.exit)
			target, err := github.NewSemanticVersion(tc.target)
			if err != nil {
				t.Fatal(err)
			}
			err = updateToVersion(target, tc.force, tc.current)
			if tc.wantCode == 0 {
				if err != nil {
					t.Errorf("successful installer returned error: %v", err)
				}
			} else {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) || exitErr.ExitCode() != tc.wantCode {
					t.Errorf("installer error=%v, want exit code %d", err, tc.wantCode)
				}
			}
			data, err := os.ReadFile(logFile)
			if err != nil {
				t.Fatal(err)
			}
			want := "install github.com/wailsapp/wails/v3/cmd/wails3@" + tc.target + "\n"
			if string(data) != want {
				t.Errorf("installer invocations=%q, want exactly %q", strings.TrimSpace(string(data)), strings.TrimSpace(want))
			}
		})
	}
}
