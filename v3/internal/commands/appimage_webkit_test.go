package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRelocateWebKitHelpers(t *testing.T) {
	libDir := t.TempDir()
	original := []byte("\x00/usr/lib/x86_64-linux-gnu/webkitgtk-6.0\x00/usr/lib/x86_64-linux-gnu/webkitgtk-6.0/injected-bundle/\x00/usr/share\x00")
	lib := filepath.Join(libDir, "libwebkitgtk-6.0.so.4")
	if err := os.WriteFile(lib, original, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(lib, 0755); err != nil { // independent of the umask
		t.Fatal(err)
	}
	other := filepath.Join(libDir, "libgtk-4.so.1")
	if err := os.WriteFile(other, original, 0644); err != nil {
		t.Fatal(err)
	}

	if err := relocateWebKitHelpers(libDir, []string{"/usr/lib/x86_64-linux-gnu/webkitgtk-6.0"}); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(lib)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte("\x00././/lib/x86_64-linux-gnu/webkitgtk-6.0\x00././/lib/x86_64-linux-gnu/webkitgtk-6.0/injected-bundle/\x00/usr/share\x00")
	if !bytes.Equal(got, want) {
		t.Errorf("library not relocated:\n got %q\nwant %q", got, want)
	}
	if info, err := os.Stat(lib); err != nil || info.Mode().Perm() != 0755 {
		t.Errorf("library mode changed: %v %v", info.Mode(), err)
	}
	if got, _ := os.ReadFile(other); !bytes.Equal(got, original) {
		t.Errorf("non-WebKit library was modified: %q", got)
	}

	if err := relocateWebKitHelpers(libDir, []string{"/usr/libexec/webkitgtk-6.0"}); err == nil {
		t.Error("expected an error when no WebKit library references the helper directory")
	}
}
