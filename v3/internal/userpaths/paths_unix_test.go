//go:build !windows && !darwin

package userpaths

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestUnixSearchOrder(t *testing.T) {
	home := t.TempDir()
	p := platformDefaults(home)
	if p.DataHome != filepath.Join(home, ".local/share") || p.ConfigHome != filepath.Join(home, ".config") || p.StateHome != filepath.Join(home, ".local/state") || p.CacheHome != filepath.Join(home, ".cache") {
		t.Fatalf("incorrect default locations: %+v", p)
	}
	custom := filepath.Join(home, "custom")
	extra := filepath.Join(home, "extra")
	p.DataHome = custom
	p.DataDirs = []string{extra, "/usr/share", extra}
	p.searchDirectories()
	want := []string{filepath.Join(custom, "applications"), filepath.Join(home, ".local/share/applications"), "/usr/local/share/applications", "/usr/share/applications", filepath.Join(extra, "applications")}
	if !reflect.DeepEqual(p.ApplicationDirs, want) {
		t.Fatalf("application search order=%v, want %v", p.ApplicationDirs, want)
	}
	if p.FontDirs[0] != filepath.Join(custom, "fonts") || p.FontDirs[1] != filepath.Join(home, ".fonts") {
		t.Fatalf("font priority=%v", p.FontDirs)
	}
}

func TestLoadPersonalDirectories(t *testing.T) {
	defaults := personalDefaults("")
	for _, field := range defaults.personalSettings() {
		t.Setenv(field.key, "")
	}
	home, config := t.TempDir(), t.TempDir()
	p := personalDefaults(home)
	p.ConfigHome = config
	if err := os.WriteFile(filepath.Join(config, "user-dirs.dirs"), []byte("XDG_DESKTOP_DIR=\"$HOME/Workspace\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	loadPersonalDirectories(&p)
	if p.Desktop != filepath.Join(home, "Workspace") {
		t.Fatalf("desktop=%q", p.Desktop)
	}
}
