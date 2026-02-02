package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDirectorySafety(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		targetDir      string
		ciMode         bool
		force          bool
		setup          func() string // returns path to use, may create files
		wantResult     DirectorySafetyResult
		wantErr        bool
		wantErrType    string
	}{
		{
			name:       "empty target dir string - should be safe",
			targetDir:  "",
			ciMode:     false,
			force:      false,
			setup:      func() string { return "" },
			wantResult: DirectorySafe,
			wantErr:    false,
		},
		{
			name:      "non-existent directory - should be safe",
			targetDir: "",
			ciMode:    false,
			force:     false,
			setup: func() string {
				return filepath.Join(tempDir, "nonexistent")
			},
			wantResult: DirectorySafe,
			wantErr:    false,
		},
		{
			name:      "empty existing directory - should be safe",
			targetDir: "",
			ciMode:    false,
			force:     false,
			setup: func() string {
				dir := filepath.Join(tempDir, "empty_dir")
				os.Mkdir(dir, 0755)
				return dir
			},
			wantResult: DirectorySafe,
			wantErr:    false,
		},
		{
			name:      "non-empty directory with force flag - should be safe",
			targetDir: "",
			ciMode:    false,
			force:     true,
			setup: func() string {
				dir := filepath.Join(tempDir, "nonempty_force")
				os.Mkdir(dir, 0755)
				os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644)
				return dir
			},
			wantResult: DirectorySafe,
			wantErr:    false,
		},
		{
			name:      "non-empty directory in CI mode - should return error",
			targetDir: "",
			ciMode:    true,
			force:     false,
			setup: func() string {
				dir := filepath.Join(tempDir, "nonempty_ci")
				os.Mkdir(dir, 0755)
				os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644)
				return dir
			},
			wantResult:  DirectorySafe,
			wantErr:     true,
			wantErrType: "*main.DirectorySafetyError",
		},
		{
			name:      "non-empty directory in interactive mode - should need confirmation",
			targetDir: "",
			ciMode:    false,
			force:     false,
			setup: func() string {
				dir := filepath.Join(tempDir, "nonempty_interactive")
				os.Mkdir(dir, 0755)
				os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644)
				return dir
			},
			wantResult: DirectoryNeedsConfirmation,
			wantErr:    false,
		},
		{
			name:      "non-empty directory with .git folder - should need confirmation",
			targetDir: "",
			ciMode:    false,
			force:     false,
			setup: func() string {
				dir := filepath.Join(tempDir, "with_git")
				os.Mkdir(dir, 0755)
				os.Mkdir(filepath.Join(dir, ".git"), 0755)
				os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("[core]"), 0644)
				return dir
			},
			wantResult: DirectoryNeedsConfirmation,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDir := tt.setup()
			if tt.targetDir != "" {
				targetDir = tt.targetDir
			}

			result, err := CheckDirectorySafety(targetDir, tt.ciMode, tt.force)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckDirectorySafety() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check error type if expected
			if tt.wantErr && tt.wantErrType != "" {
				_, ok := err.(*DirectorySafetyError)
				if !ok {
					t.Errorf("CheckDirectorySafety() error type = %T, want %s", err, tt.wantErrType)
				}
			}

			// Check result
			if result != tt.wantResult {
				t.Errorf("CheckDirectorySafety() result = %v, want %v", result, tt.wantResult)
			}
		})
	}
}

func TestDirectorySafetyError_Error(t *testing.T) {
	err := &DirectorySafetyError{TargetDir: "/some/path"}
	expected := "target directory '/some/path' is not empty. Aborting to prevent data loss. Use an empty directory or remove existing files first"
	if err.Error() != expected {
		t.Errorf("DirectorySafetyError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestGetAbsoluteTargetDir(t *testing.T) {
	tests := []struct {
		name      string
		targetDir string
		wantEmpty bool
		wantErr   bool
	}{
		{
			name:      "empty string",
			targetDir: "",
			wantEmpty: true,
			wantErr:   false,
		},
		{
			name:      "relative path",
			targetDir: "relative/path",
			wantEmpty: false,
			wantErr:   false,
		},
		{
			name:      "absolute path",
			targetDir: "/absolute/path",
			wantEmpty: false,
			wantErr:   false,
		},
		{
			name:      "current directory",
			targetDir: ".",
			wantEmpty: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetAbsoluteTargetDir(tt.targetDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAbsoluteTargetDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantEmpty && result != "" {
				t.Errorf("GetAbsoluteTargetDir() = %v, want empty string", result)
			}

			if !tt.wantEmpty {
				if result == "" {
					t.Error("GetAbsoluteTargetDir() returned empty string, want non-empty")
				}
				if !filepath.IsAbs(result) {
					t.Errorf("GetAbsoluteTargetDir() = %v, want absolute path", result)
				}
			}
		})
	}
}
