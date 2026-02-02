package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckDirectorySafety(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		force   bool
		setup   func() string // returns path to use, may create files
		wantErr bool
		errMsg  string // substring to check in error message
	}{
		{
			name:    "empty target dir string - should be safe",
			force:   false,
			setup:   func() string { return "" },
			wantErr: false,
		},
		{
			name:  "non-existent directory - should be safe",
			force: false,
			setup: func() string {
				return filepath.Join(tempDir, "nonexistent")
			},
			wantErr: false,
		},
		{
			name:  "empty existing directory - should be safe",
			force: false,
			setup: func() string {
				dir := filepath.Join(tempDir, "empty_dir")
				os.Mkdir(dir, 0755)
				return dir
			},
			wantErr: false,
		},
		{
			name:  "non-empty directory with force flag - should be safe",
			force: true,
			setup: func() string {
				dir := filepath.Join(tempDir, "nonempty_force")
				os.Mkdir(dir, 0755)
				os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644)
				return dir
			},
			wantErr: false,
		},
		{
			name:  "non-empty directory without force - should return error",
			force: false,
			setup: func() string {
				dir := filepath.Join(tempDir, "nonempty_no_force")
				os.Mkdir(dir, 0755)
				os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644)
				return dir
			},
			wantErr: true,
			errMsg:  "Use -f to force",
		},
		{
			name:  "non-empty directory with .git folder - should return error",
			force: false,
			setup: func() string {
				dir := filepath.Join(tempDir, "with_git")
				os.Mkdir(dir, 0755)
				os.Mkdir(filepath.Join(dir, ".git"), 0755)
				os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("[core]"), 0644)
				return dir
			},
			wantErr: true,
			errMsg:  "not empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDir := tt.setup()

			err := CheckDirectorySafety(targetDir, tt.force)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckDirectorySafety() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check error message contains expected substring
			if tt.wantErr && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("CheckDirectorySafety() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}
