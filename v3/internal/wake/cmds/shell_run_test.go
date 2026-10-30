package cmds

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestShellCmdPreservesExecutionOptions(t *testing.T) {
	base := t.TempDir()
	script := `echo "$WAKE_SHELL_TEST_VALUE" > marker.txt && echo reported && echo diagnostic >&2`
	if runtime.GOOS == "windows" {
		script = `echo %WAKE_SHELL_TEST_VALUE%> marker.txt && echo reported && echo diagnostic 1>&2`
	}
	var stdout, stderr bytes.Buffer
	command := &ShellCmd{
		Cmd: script, Dir: base,
		Env:    append(os.Environ(), "WAKE_SHELL_TEST_VALUE=expected"),
		Output: Output{Stdout: &stdout, Stderr: &stderr},
	}
	if err := command.Run(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(base, "marker.txt"))
	if err != nil || strings.TrimSpace(string(data)) != "expected" {
		t.Errorf("shell output file = %q, %v", data, err)
	}
	if strings.TrimSpace(stdout.String()) != "reported" || strings.TrimSpace(stderr.String()) != "diagnostic" {
		t.Errorf("stdout = %q, stderr = %q", stdout.String(), stderr.String())
	}
}

func TestShellCmdPreservesExitStatus(t *testing.T) {
	err := (&ShellCmd{Cmd: "exit 7"}).Run()
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 7 {
		t.Errorf("shell error = %v, want exit status 7", err)
	}
}
