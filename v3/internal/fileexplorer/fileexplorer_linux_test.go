//go:build linux

package fileexplorer

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func writeSearchDesktopFile(t *testing.T, dataDir string) string {
	t.Helper()
	dir := filepath.Join(dataDir, "applications")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "wails-xdg-search-test.desktop")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	return file
}

func TestFindDesktopFileXDGDataDirs(t *testing.T) {
	root := t.TempDir()
	t.Setenv("HOME", filepath.Join(root, "home"))
	user := filepath.Join(root, "user")
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	t.Setenv("XDG_DATA_HOME", user)
	t.Setenv("XDG_DATA_DIRS", first+string(os.PathListSeparator)+second)
	files := []string{
		writeSearchDesktopFile(t, user),
		writeSearchDesktopFile(t, first),
		writeSearchDesktopFile(t, second),
	}
	writeSearchDesktopFile(t, filepath.Join(root, "home", ".local", "share"))
	for _, want := range files {
		got, err := findDesktopFile(filepath.Base(want))
		if err != nil || got != want {
			t.Fatalf("findDesktopFile = %q, %v; want %q", got, err, want)
		}
		if err := os.Remove(want); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFindDesktopFileDefaultHomeIgnoresRelativePaths(t *testing.T) {
	for _, dataHome := range []string{"", "relative-user"} {
		t.Run("home="+dataHome, func(t *testing.T) {
			root := t.TempDir()
			t.Chdir(root)
			home := filepath.Join(root, "home")
			t.Setenv("HOME", home)
			t.Setenv("XDG_DATA_HOME", dataHome)
			t.Setenv("XDG_DATA_DIRS", filepath.Join(root, "system"))
			want := writeSearchDesktopFile(t, filepath.Join(home, ".local", "share"))
			writeSearchDesktopFile(t, ".")
			writeSearchDesktopFile(t, "relative-user")
			got, err := findDesktopFile(filepath.Base(want))
			if err != nil || got != want {
				t.Errorf("findDesktopFile = %q, %v; want %q", got, err, want)
			}
		})
	}
}

func TestFindDesktopFileIgnoresRelativeDataDirs(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	t.Setenv("HOME", filepath.Join(root, "home"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(root, "user"))
	system := filepath.Join(root, "system")
	t.Setenv("XDG_DATA_DIRS", "relative-system"+string(os.PathListSeparator)+system)
	writeSearchDesktopFile(t, "relative-system")
	want := writeSearchDesktopFile(t, system)
	got, err := findDesktopFile(filepath.Base(want))
	if err != nil || got != want {
		t.Errorf("findDesktopFile = %q, %v; want %q", got, err, want)
	}
}

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
	uriFile := filepath.Join(root, "selected file #1 %u %20.txt")
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
		{"path/spaced-single", "Exec=wails-test-manager %f\n", uriFile, "wails-test-manager", uriFile},
		{"path/spaced-list", "Exec=wails-test-manager %F\n", uriFile, "wails-test-manager", uriFile},
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
