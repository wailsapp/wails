package browser

import (
	"os"
	"os/exec"
	"testing"
)

// TestOpen_Success covers the happy path via a stub that runs the test binary
// itself with no-match filter (exits 0, no side effects, cross-platform).
func TestOpen_Success(t *testing.T) {
	orig := openCmd
	openCmd = func(_ string) *exec.Cmd { return exec.Command(os.Args[0], "-test.run=^$") }
	t.Cleanup(func() { openCmd = orig })

	if err := OpenURL("https://example.com"); err != nil {
		t.Errorf("OpenURL unexpected error: %v", err)
	}
	if err := OpenFile("some/file"); err != nil {
		t.Errorf("OpenFile unexpected error: %v", err)
	}
}

// TestOpen_StartError exercises the cmd.Start() error path.
func TestOpen_StartError(t *testing.T) {
	t.Setenv("PATH", "/nonexistent_path_that_does_not_exist")
	if err := OpenURL("https://example.com"); err == nil {
		t.Error("expected error with empty PATH, got nil")
	}
}
