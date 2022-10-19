package fs

import (
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestRelativePath(t *testing.T) {

	i := is.New(t)

	cwd, err := os.Getwd()
	i.Equal(err, nil)

	// Check current directory
	actual := RelativePath(".")
	i.Equal(actual, cwd)

	// Check 2 parameters
	actual = RelativePath("..", "fs")
	i.Equal(actual, cwd)

	// Check 3 parameters including filename
	actual = RelativePath("..", "fs", "fs.go")
	expected := filepath.Join(cwd, "fs.go")
	i.Equal(actual, expected)

}

func Test_FindFileInParents(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (startDir string, configDir string)
		wantErr bool
	}{
		{
			name: "should error when no wails.json file is found in local or parent dirs",
			setup: func() (string, string) {
				tempDir := os.TempDir()
				testDir := lo.Must(os.MkdirTemp(tempDir, "projectPath"))
				_ = os.MkdirAll(testDir, 0755)
				return testDir, ""
			},
			wantErr: true,
		},
		{
			name: "should find wails.json in local path",
			setup: func() (string, string) {
				tempDir := os.TempDir()
				testDir := lo.Must(os.MkdirTemp(tempDir, "projectPath"))
				_ = os.MkdirAll(testDir, 0755)
				configFile := filepath.Join(testDir, "wails.json")
				_ = os.WriteFile(configFile, []byte("{}"), 0755)
				return testDir, configFile
			},
			wantErr: false,
		},
		{
			name: "should find wails.json in parent path",
			setup: func() (string, string) {
				tempDir := os.TempDir()
				testDir := lo.Must(os.MkdirTemp(tempDir, "projectPath"))
				_ = os.MkdirAll(testDir, 0755)
				parentDir := filepath.Dir(testDir)
				configFile := filepath.Join(parentDir, "wails.json")
				_ = os.WriteFile(configFile, []byte("{}"), 0755)
				return testDir, configFile
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, expectedPath := tt.setup()
			defer func() {
				if expectedPath != "" {
					_ = os.Remove(expectedPath)
				}
			}()
			got := FindFileInParents(path, "wails.json")
			if got != expectedPath {
				t.Errorf("FindFileInParents() got = %v, want %v", got, expectedPath)
			}
		})
	}
}
