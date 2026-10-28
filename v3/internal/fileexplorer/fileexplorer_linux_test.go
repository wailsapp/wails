//go:build linux

package fileexplorer

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestExplorerBinArgsDesktopCommand(t *testing.T) {
	root := t.TempDir()
	applications := filepath.Join(root, "applications")
	if err := os.Mkdir(applications, 0755); err != nil {
		t.Fatal(err)
	}
	query := filepath.Join(root, "xdg-mime")
	if err := os.WriteFile(query, []byte("#!/bin/sh\nprintf 'wails-test.desktop\\n'\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", root)
	t.Setenv("XDG_DATA_HOME", root)
	file := filepath.Join(root, "selected.txt")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	uriFile := filepath.Join(root, "selected file #1 %20.txt")
	if err := os.WriteFile(uriFile, nil, 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		command  string
		path     string
		wantBin  string
		wantPath string
	}{
		{"missing/file", "", file, "xdg-open", root},
		{"missing/directory", "", root, "xdg-open", root},
		{"empty/file", "Exec=\n", file, "xdg-open", root},
		{"empty/directory", "Exec=\n", root, "xdg-open", root},
		{"whitespace/file", "Exec= \t \n", file, "xdg-open", root},
		{"whitespace/directory", "Exec= \t \n", root, "xdg-open", root},
		{"removed codes/file", "Exec=%d %D %n %N %v %m\n", file, "xdg-open", root},
		{"removed codes/directory", "Exec=%d %D %n %N %v %m\n", root, "xdg-open", root},
		{"valid/file", "Exec=wails-test-manager %f\n", file, "wails-test-manager", file},
		{"valid/directory", "Exec=wails-test-manager %f\n", root, "wails-test-manager", root},
		{"URI/plain", "Exec=wails-test-manager %U\n", file, "wails-test-manager", pathToURI(file)},
		{"URI/escaped-list", "Exec=wails-test-manager %U\n", uriFile, "wails-test-manager", pathToURI(uriFile)},
		{"URI/escaped-single", "Exec=wails-test-manager %u\n", uriFile, "wails-test-manager", pathToURI(uriFile)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.WriteFile(filepath.Join(applications, "wails-test.desktop"),
				[]byte("[Desktop Entry]\nName=Test\n"+tt.command), 0644); err != nil {
				t.Fatal(err)
			}
			bin, args, err := explorerBinArgs(tt.path, true)
			if err != nil {
				t.Fatal(err)
			}
			if bin != tt.wantBin || !slices.Equal(args, []string{tt.wantPath}) {
				t.Errorf("got %q %q, want %q [%q]", bin, args, tt.wantBin, tt.wantPath)
			}
		})
	}
}
