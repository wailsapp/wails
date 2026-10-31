package cmds

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestGoVetRespectsPackageAndAnalyzer(t *testing.T) {
	for _, tc := range []struct {
		name, command string
		wantExit      int
	}{
		{"requested package", "go vet ./sub", 1},
		{"disabled analyzer", "go vet -printf=false ./sub", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, output, err := runOfflineGo(t, tc.command, map[string]string{
				"go.mod":     "module example.com/wake-vet-check\n\ngo 1.25.0\n",
				"main.go":    "package fixture\n",
				"sub/sub.go": "package sub\n\nimport \"fmt\"\n\nfunc check() { fmt.Printf(\"%d\", \"text\") }\n",
			})
			if tc.wantExit == 0 {
				if err != nil {
					t.Fatalf("disabled analyzer failed: %v; output: %s", err, output)
				}
				return
			}
			var exitError *exec.ExitError
			if !errors.As(err, &exitError) || exitError.ExitCode() != tc.wantExit {
				t.Fatalf("vet error = %v, want exit %d; output: %s", err, tc.wantExit, output)
			}
			if !strings.Contains(output, "Printf format %d") {
				t.Errorf("missing requested-package diagnostic: %s", output)
			}
		})
	}
}
