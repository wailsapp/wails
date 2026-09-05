package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leaanthony/dmg/dmg"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestAddDMGFiles(t *testing.T) {
	dir := t.TempDir()
	resource := filepath.Join(dir, "install.command")
	if err := os.WriteFile(resource, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	opts := &dmg.Options{Files: map[string]string{}}
	if err := addDMGFiles(opts, "Install.command="+resource); err != nil {
		t.Fatalf("addDMGFiles() error = %v", err)
	}
	if got := opts.Files["Install.command"]; got != resource {
		t.Fatalf("resource path = %q, want %q", got, resource)
	}
}

func TestAddDMGFilesRejectsMalformedPair(t *testing.T) {
	if err := addDMGFiles(&dmg.Options{Files: map[string]string{}}, "install.command"); err == nil {
		t.Fatal("addDMGFiles() accepted malformed pair")
	}
}

func TestBuildDMGOptions(t *testing.T) {
	dir := t.TempDir()
	appPath := filepath.Join(dir, "Example.app")
	resource := filepath.Join(dir, "install.command")
	if err := os.WriteFile(resource, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	opts, err := buildDMGOptions(&flags.ToolPackage{
		ExecutableName:  "Example",
		BackgroundImage: "background.png",
		DmgVolumeIcon:   "volume.icns",
		DmgFileIcon:     "file.icns",
		DmgFiles:        "install.command=" + resource,
		DmgWindowWidth:  700,
		DmgWindowHeight: 500,
	}, appPath, filepath.Join(dir, "Example.dmg"))
	if err != nil {
		t.Fatalf("buildDMGOptions() error = %v", err)
	}

	if opts.VolumeName != "Example" {
		t.Errorf("VolumeName = %q, want Example", opts.VolumeName)
	}
	if opts.Window.Width != 700 || opts.Window.Height != 500 {
		t.Errorf("Window = %#v, want 700x500", opts.Window)
	}
	if opts.Background == nil || opts.Background.File != "background.png" {
		t.Errorf("Background = %#v, want background.png", opts.Background)
	}
	if opts.VolumeIcon != "volume.icns" || opts.FileIcon != "file.icns" {
		t.Errorf("icons = (%q, %q), want (volume.icns, file.icns)", opts.VolumeIcon, opts.FileIcon)
	}
	if opts.Files["Example.app"] != appPath || opts.Files["install.command"] != resource {
		t.Errorf("Files = %#v, want app and installer resource", opts.Files)
	}
	if got := opts.IconPositions["Example.app"]; got != (dmg.IconPosition{X: 196, Y: 224}) {
		t.Errorf("app icon position = %#v, want {196 224}", got)
	}
	if got := opts.IconPositions["Applications"]; got != (dmg.IconPosition{X: 504, Y: 224}) {
		t.Errorf("Applications icon position = %#v, want {504 224}", got)
	}
}

func TestBuildDMGOptionsDefaults(t *testing.T) {
	opts, err := buildDMGOptions(&flags.ToolPackage{ExecutableName: "Example"}, "Example.app", "Example.dmg")
	if err != nil {
		t.Fatalf("buildDMGOptions() error = %v", err)
	}
	if opts.Window.Width != 540 || opts.Window.Height != 380 {
		t.Errorf("Window = %#v, want 540x380", opts.Window)
	}
	if got := opts.IconPositions["Example.app"]; got != (dmg.IconPosition{X: 151, Y: 164}) {
		t.Errorf("app icon position = %#v, want {151 164}", got)
	}
	if got := opts.IconPositions["Applications"]; got != (dmg.IconPosition{X: 388, Y: 164}) {
		t.Errorf("Applications icon position = %#v, want {388 164}", got)
	}
}

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
