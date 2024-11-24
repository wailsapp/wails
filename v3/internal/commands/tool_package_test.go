package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestToolPackage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (*flags.ToolPackage, func())
		wantErr bool
		errMsg  string
	}{
		{
			name: "should fail with invalid format",
			setup: func() (*flags.ToolPackage, func()) {
				return &flags.ToolPackage{
					Format:         "invalid",
					ConfigPath:     "config.yaml",
					ExecutableName: "myapp",
				}, func() {}
			},
			wantErr: true,
			errMsg:  "unsupported package format",
		},
		{
			name: "should fail with missing config file",
			setup: func() (*flags.ToolPackage, func()) {
				return &flags.ToolPackage{
					Format:         "deb",
					ConfigPath:     "nonexistent.yaml",
					ExecutableName: "myapp",
				}, func() {}
			},
			wantErr: true,
			errMsg:  "config file not found",
		},
		{
			name: "should handle case-insensitive format (DEB)",
			setup: func() (*flags.ToolPackage, func()) {
				// Create a temporary config file
				dir := t.TempDir()
				configPath := filepath.Join(dir, "config.yaml")
				err := os.WriteFile(configPath, []byte("name: test"), 0644)
				if err != nil {
					t.Fatal(err)
				}

				// Create bin directory
				err = os.MkdirAll(filepath.Join(dir, "bin"), 0755)
				if err != nil {
					t.Fatal(err)
				}

				return &flags.ToolPackage{
					Format:         "DEB",
					ConfigPath:     configPath,
					ExecutableName: "myapp",
				}, func() {
					os.RemoveAll(filepath.Join(dir, "bin"))
				}
			},
			wantErr: false,
		},
		{
			name: "should handle case-insensitive format (RPM)",
			setup: func() (*flags.ToolPackage, func()) {
				// Create a temporary config file
				dir := t.TempDir()
				configPath := filepath.Join(dir, "config.yaml")
				err := os.WriteFile(configPath, []byte("name: test"), 0644)
				if err != nil {
					t.Fatal(err)
				}

				// Create bin directory
				err = os.MkdirAll(filepath.Join(dir, "bin"), 0755)
				if err != nil {
					t.Fatal(err)
				}

				return &flags.ToolPackage{
					Format:         "RPM",
					ConfigPath:     configPath,
					ExecutableName: "myapp",
				}, func() {
					os.RemoveAll(filepath.Join(dir, "bin"))
				}
			},
			wantErr: false,
		},
		{
			name: "should handle case-insensitive format (ARCHLINUX)",
			setup: func() (*flags.ToolPackage, func()) {
				// Create a temporary config file
				dir := t.TempDir()
				configPath := filepath.Join(dir, "config.yaml")
				err := os.WriteFile(configPath, []byte("name: test"), 0644)
				if err != nil {
					t.Fatal(err)
				}

				// Create bin directory
				err = os.MkdirAll(filepath.Join(dir, "bin"), 0755)
				if err != nil {
					t.Fatal(err)
				}

				return &flags.ToolPackage{
					Format:         "ARCHLINUX",
					ConfigPath:     configPath,
					ExecutableName: "myapp",
				}, func() {
					os.RemoveAll(filepath.Join(dir, "bin"))
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options, cleanup := tt.setup()
			defer cleanup()

			err := ToolPackage(options)

			if (err != nil) != tt.wantErr {
				t.Errorf("ToolPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ToolPackage() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}
