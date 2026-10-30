package cmds

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestGoTestRunsAllRequestedPackages(t *testing.T) {
	_, output, err := runOfflineGo(t, "go test -buildvcs=false -count 1 ./first ./second", map[string]string{
		"go.mod":                "module example.com/wake-test-check\n\ngo 1.25.0\n",
		"first/first_test.go":   "package first\n\nimport \"testing\"\n\nfunc TestFirst(t *testing.T) { t.Fatal(\"first-package-ran\") }\n",
		"second/second_test.go": "package second\n\nimport \"testing\"\n\nfunc TestSecond(t *testing.T) {}\n",
	})
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
		t.Errorf("test command error = %v, want exit 1; output: %s", err, output)
	}
	if !strings.Contains(output, "first-package-ran") {
		t.Errorf("first package was skipped: %s", output)
	}
}
