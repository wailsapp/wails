//go:build darwin && !ios

package mac

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestResourceAPIsThroughSymlink(t *testing.T) {
	if os.Getenv("WAILS_MAC_RESOURCE_HELPER") == "1" {
		resources, err := ResourceFS()
		if err != nil {
			t.Fatalf("ResourceFS() error = %v", err)
		}
		data, err := fs.ReadFile(resources, "test.txt")
		if err != nil {
			t.Fatalf("fs.ReadFile(ResourceFS(), %q) error = %v", "test.txt", err)
		}
		if string(data) != "resource contents" {
			t.Errorf("fs.ReadFile(ResourceFS(), %q) = %q, want %q", "test.txt", data, "resource contents")
		}

		data, err = LoadResource("test.txt")
		if err != nil {
			t.Fatalf("LoadResource(%q) error = %v", "test.txt", err)
		}
		if string(data) != "resource contents" {
			t.Errorf("LoadResource(%q) = %q, want %q", "test.txt", data, "resource contents")
		}
		return
	}

	testBinary, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	bundlePath := filepath.Join(t.TempDir(), "Test.app")
	canonicalExecutable := filepath.Join(bundlePath, "Contents", "MacOS", "Test")
	resourcesPath := filepath.Join(bundlePath, "Contents", "Resources")
	if err := os.MkdirAll(resourcesPath, 0o755); err != nil {
		t.Fatalf("os.MkdirAll(%q) error = %v", resourcesPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(canonicalExecutable), 0o755); err != nil {
		t.Fatalf("os.MkdirAll(%q) error = %v", filepath.Dir(canonicalExecutable), err)
	}
	testBinaryData, err := os.ReadFile(testBinary)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", testBinary, err)
	}
	if err := os.WriteFile(canonicalExecutable, testBinaryData, 0o755); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", canonicalExecutable, err)
	}
	if err := os.WriteFile(filepath.Join(resourcesPath, "test.txt"), []byte("resource contents"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(test resource) error = %v", err)
	}

	launcher := filepath.Join(t.TempDir(), "launcher")
	if err := os.Symlink(canonicalExecutable, launcher); err != nil {
		t.Fatalf("os.Symlink(%q, %q) error = %v", canonicalExecutable, launcher, err)
	}

	cmd := exec.Command(launcher, "-test.run=^TestResourceAPIsThroughSymlink$")
	cmd.Env = append(os.Environ(), "WAILS_MAC_RESOURCE_HELPER=1")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("test binary launched through symlink failed: %v\n%s", err, output)
	}
}

func TestResourcePathForExecutable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		executable string
		want       string
		ok         bool
	}{
		{
			name:       "application bundle",
			executable: "/Applications/My App.app/Contents/MacOS/My App",
			want:       "/Applications/My App.app/Contents/Resources",
			ok:         true,
		},
		{
			name:       "nested application bundle",
			executable: "/Users/user/Desktop/My App.app/Contents/MacOS/My App",
			want:       "/Users/user/Desktop/My App.app/Contents/Resources",
			ok:         true,
		},
		{
			name:       "normalises executable path",
			executable: "/Applications/My App.app/Contents/MacOS/../MacOS/My App",
			want:       "/Applications/My App.app/Contents/Resources",
			ok:         true,
		},
		{
			name:       "executable outside application bundle",
			executable: "/usr/local/bin/my-app",
			ok:         false,
		},
		{
			name:       "missing Contents directory",
			executable: "/Applications/My App.app/MacOS/My App",
			ok:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := resourcePathForExecutable(tt.executable)
			if ok != tt.ok {
				t.Fatalf("resourcePathForExecutable(%q) returned ok=%t, want %t", tt.executable, ok, tt.ok)
			}
			if got != filepath.FromSlash(tt.want) {
				t.Errorf("resourcePathForExecutable(%q) = %q, want %q", tt.executable, got, filepath.FromSlash(tt.want))
			}
		})
	}
}
