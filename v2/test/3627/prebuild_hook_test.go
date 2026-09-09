package test_3627

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPreBuildHookRunsInProjectDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wails-test-3627-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	hookScript := filepath.Join(tmpDir, "testhook.sh")
	outputFile := filepath.Join(tmpDir, "pwd_output.txt")

	scriptContent := "#!/bin/sh\npwd > " + outputFile + "\n"
	if err := os.WriteFile(hookScript, []byte(scriptContent), 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(hookScript)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	gotPwd := string(data)
	gotPwd = gotPwd[:len(gotPwd)-1]

	if gotPwd != tmpDir {
		t.Errorf("hook ran in %q, want %q", gotPwd, tmpDir)
	}
}
