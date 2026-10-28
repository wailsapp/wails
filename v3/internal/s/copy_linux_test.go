package s

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

func TestCopyClosesDestination(t *testing.T) {
	// Finalizers must not hide leaked descriptors during the check.
	previousGC := debug.SetGCPercent(-1)
	t.Cleanup(func() { debug.SetGCPercent(previousGC) })
	dir := t.TempDir()
	source, target := filepath.Join(dir, "source"), filepath.Join(dir, "target")
	if err := os.WriteFile(source, []byte("copied content"), 0600); err != nil {
		t.Fatal(err)
	}
	for range 32 {
		COPY(source, target)
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "copied content" {
		t.Fatalf("copied data=%q, error=%v", data, err)
	}
	fds, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		t.Fatal(err)
	}
	leaked := 0
	for _, fd := range fds {
		name, err := os.Readlink(filepath.Join("/proc/self/fd", fd.Name()))
		if err == nil && name == target {
			leaked++
		}
	}
	if leaked != 0 {
		t.Errorf("COPY left %d destination descriptors open after 32 copies", leaked)
	}
}
