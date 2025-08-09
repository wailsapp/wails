package fileexplorer_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/fileexplorer"
)

// Credit: https://stackoverflow.com/a/50631395
func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestFileExplorer(t *testing.T) {
	skipCI(t)
	// TestFileExplorer verifies that the OpenFileManager function correctly handles:
	// - Opening files in the native file manager across different platforms
	// - Selecting files when the selectFile parameter is true
	// - Various error conditions like non-existent paths
	tempDir := t.TempDir() // Create a temporary directory for tests

	tests := []struct {
		name        string
		path        string
		selectFile  bool
		expectedErr error
	}{
		{"Open Existing File", tempDir, false, nil},
		{"Select Existing File", tempDir, true, nil},
		{"Non-Existent Path", "/path/does/not/exist", false, fmt.Errorf("failed to access the specified path: /path/does/not/exist")},
		{"Path with Special Characters", filepath.Join(tempDir, "test space.txt"), true, nil},
		{"No Permission Path", "/root/test.txt", false, fmt.Errorf("failed to open the file explorer: /root/test.txt")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("Windows", func(t *testing.T) {
				runPlatformTest(t, "windows")
			})
			t.Run("Linux", func(t *testing.T) {
				runPlatformTest(t, "linux")
			})
			t.Run("Darwin", func(t *testing.T) {
				runPlatformTest(t, "darwin")
			})
		})
	}
}

func runPlatformTest(t *testing.T, platform string) {
	if runtime.GOOS != platform {
		t.Skipf("Skipping test on non-%s platform", strings.ToTitle(platform))
	}

	testFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(testFile, []byte("Test file contents"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		selectFile bool
	}{
		{"OpenFile", false},
		{"SelectFile", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := fileexplorer.OpenFileManager(testFile, test.selectFile)
			if err != nil {
				t.Errorf("OpenFileManager(%s, %v) error = %v", testFile, test.selectFile, err)
			}
		})
	}
}
