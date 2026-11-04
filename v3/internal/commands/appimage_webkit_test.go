package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRelocateWebKitHelpers(t *testing.T) {
	appDir := t.TempDir()
	helperDir := "/usr/lib/x86_64-linux-gnu/webkitgtk-6.0"
	if err := os.MkdirAll(filepath.Join(appDir, helperDir), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appDir, helperDir, "WebKitWebProcess"), nil, 0755); err != nil {
		t.Fatal(err)
	}
	libDir := filepath.Join(appDir, "usr", "lib")
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

	stray := "/usr/local/src/other/AppDir/usr/lib/x86_64-linux-gnu/webkitgtk-6.0"
	if err := relocateWebKitHelpers(appDir, []string{stray, helperDir}); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(lib)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte("\x00././webkitgtk-6.0/\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00/usr/lib/x86_64-linux-gnu/webkitgtk-6.0/injected-bundle/\x00/usr/share\x00")
	if !bytes.Equal(got, want) {
		t.Errorf("library not relocated:\n got %q\nwant %q", got, want)
	}
	// Windows has no permission bits to keep; AppImages are built on Linux.
	if runtime.GOOS != "windows" {
		if info, err := os.Stat(lib); err != nil {
			t.Errorf("stat library: %v", err)
		} else if info.Mode().Perm() != 0755 {
			t.Errorf("library mode changed: %v", info.Mode())
		}
	}
	if got, _ := os.ReadFile(other); !bytes.Equal(got, original) {
		t.Errorf("non-WebKit library was modified: %q", got)
	}
	if _, err := os.Stat(filepath.Join(libDir, "webkitgtk-6.0", "WebKitWebProcess")); err != nil {
		t.Errorf("helpers not moved beside the libraries: %v", err)
	}

	if err := relocateWebKitHelpers(appDir, []string{"/usr/libexec/webkitgtk-6.0"}); err == nil {
		t.Error("expected an error when no WebKit library references the helper directory")
	}
}
