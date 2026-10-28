//go:build !windows

package git

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestHeadHashSeparatesDiagnostics(t *testing.T) {
	for _, exitCode := range []int{0, 1} {
		t.Run(strconv.Itoa(exitCode), func(t *testing.T) {
			bin := t.TempDir()
			script := "#!/bin/sh\nprintf 'warning: diagnostic output\\n' >&2\nprintf '0123456789abcdef\\n'\nexit " + strconv.Itoa(exitCode) + "\n"
			if err := os.WriteFile(filepath.Join(bin, "git"), []byte(script), 0755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", bin)
			hash, err := HeadHash(t.TempDir())
			if exitCode != 0 {
				if err == nil || !strings.Contains(err.Error(), "diagnostic output") || !strings.Contains(err.Error(), "0123456789abcdef") {
					t.Fatalf("failed Git diagnostics were lost: %v", err)
				}
				return
			}
			if err != nil || hash != "01234567" {
				t.Errorf("HeadHash = %q, %v; want 01234567, nil", hash, err)
			}
		})
	}
}
