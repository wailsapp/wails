//go:build !windows

package updater

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// failCrossDevice makes rename report EXDEV for moves out of src, as when
// the staged artifact sits on a different filesystem from the target.
func failCrossDevice(t *testing.T, src string) {
	t.Helper()
	prev := rename
	rename = func(from, to string) error {
		if from == src {
			return &os.LinkError{Op: "rename", Old: from, New: to, Err: syscall.EXDEV}
		}
		return prev(from, to)
	}
	t.Cleanup(func() { rename = prev })
}

func TestReplaceTarget_CrossDevice_File(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app")
	newPath := filepath.Join(t.TempDir(), "app.new")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))
	failCrossDevice(t, newPath)

	if err := replaceTarget(target, newPath); err != nil {
		t.Fatalf("replaceTarget: %v", err)
	}
	if got := readFile(t, target); string(got) != "NEW" {
		t.Errorf("target contents: %q", got)
	}
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Errorf("source should be removed after copy: %v", err)
	}
	if _, err := os.Stat(target + ".wails-new"); !os.IsNotExist(err) {
		t.Errorf("intermediate copy left behind: %v", err)
	}
}

func TestReplaceTarget_CrossDevice_Directory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "App.app")
	newPath := filepath.Join(t.TempDir(), "App.app")
	writeFile(t, filepath.Join(target, "Contents", "MacOS", "App"), []byte("OLD"))
	writeFile(t, filepath.Join(newPath, "Contents", "MacOS", "App"), []byte("NEW"))
	failCrossDevice(t, newPath)

	if err := replaceTarget(target, newPath); err != nil {
		t.Fatalf("replaceTarget: %v", err)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "App")); string(got) != "NEW" {
		t.Errorf("bundle contents: %q", got)
	}
}

func TestReplaceTarget_CrossDevice_CopyFailureKeepsTarget(t *testing.T) {
	target := filepath.Join(t.TempDir(), "app")
	newPath := filepath.Join(t.TempDir(), "app.new")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))
	if err := os.Chmod(newPath, 0); err != nil {
		t.Fatal(err)
	}
	failCrossDevice(t, newPath)

	if err := replaceTarget(target, newPath); err == nil {
		t.Fatal("replaceTarget succeeded with an unreadable source")
	}
	if got := readFile(t, target); string(got) != "OLD" {
		t.Errorf("target contents after failed copy: %q", got)
	}
	if _, err := os.Stat(target + ".wails-old"); !os.IsNotExist(err) {
		t.Errorf("aside left behind: %v", err)
	}
}
