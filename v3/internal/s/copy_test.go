package s

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFileIdentity(t *testing.T) {
	// COPY exits on failure, so run it in a child process.
	if os.Getenv("WAILS_TEST_COPY_PROCESS") == "1" {
		COPY(os.Getenv("WAILS_TEST_COPY_SOURCE"), os.Getenv("WAILS_TEST_COPY_TARGET"))
		return
	}
	for _, name := range []string{"same_path", "same_directory", "symlink", "hardlink", "new_file", "existing_file", "directory"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			source := filepath.Join(dir, "source")
			const content = "source contents must survive"
			if err := os.WriteFile(source, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(dir, "target")
			outputPath := target
			wantExit := 1
			var err error
			switch name {
			case "same_path":
				target = source
			case "same_directory":
				target = dir
			case "symlink":
				err = os.Symlink(source, target)
			case "hardlink":
				err = os.Link(source, target)
			case "new_file":
				wantExit = 0
			case "existing_file":
				wantExit = 0
				err = os.WriteFile(target, []byte("old destination content longer than source contents"), 0600)
			case "directory":
				wantExit = 0
				err = os.Mkdir(target, 0700)
				outputPath = filepath.Join(target, filepath.Base(source))
			}
			if err != nil {
				if name == "symlink" || name == "hardlink" {
					t.Skipf("cannot create %s on this filesystem: %v", name, err)
				}
				t.Fatal(err)
			}
			cmd := exec.Command(os.Args[0], "-test.run=^TestCopyFileIdentity$")
			cmd.Env = append(os.Environ(), "WAILS_TEST_COPY_PROCESS=1", "WAILS_TEST_COPY_SOURCE="+source, "WAILS_TEST_COPY_TARGET="+target)
			output, err := cmd.CombinedOutput()
			if cmd.ProcessState == nil {
				t.Fatalf("start copy process: %v", err)
			}
			if got := cmd.ProcessState.ExitCode(); got != wantExit {
				t.Errorf("exit=%d, want %d: %s", got, wantExit, output)
			}
			if wantExit == 1 && !strings.Contains(string(output), "same file") {
				t.Errorf("missing same-file diagnostic: %s", output)
			}
			data, err := os.ReadFile(source)
			if err != nil || string(data) != content {
				t.Errorf("source=%q, error=%v", data, err)
			}
			if wantExit == 0 {
				data, err = os.ReadFile(outputPath)
				if err != nil || string(data) != content {
					t.Errorf("destination=%q, error=%v", data, err)
				}
			}
		})
	}
}
