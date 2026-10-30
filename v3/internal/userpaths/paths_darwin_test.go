package userpaths

import (
	"path/filepath"
	"testing"
)

func TestDarwinDefaults(t *testing.T) {
	home := t.TempDir()
	p := platformDefaults(home)
	support := filepath.Join(home, "Library", "Application Support")
	if p.DataHome != support || p.ConfigHome != support || p.StateHome != support || p.RuntimeDir != support || p.CacheHome != filepath.Join(home, "Library", "Caches") || p.Videos != filepath.Join(home, "Movies") {
		t.Fatalf("incorrect macOS defaults: %+v", p)
	}
	if len(p.ApplicationDirs) != 1 || p.ApplicationDirs[0] != "/Applications" {
		t.Fatalf("application directories=%v", p.ApplicationDirs)
	}
}
