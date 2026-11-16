package userpaths

import (
	"path/filepath"
)

func platformDefaults(home string) Locations {
	p := personalDefaults(home)
	support := filepath.Join(home, "Library", "Application Support")
	p.DataHome, p.ConfigHome, p.StateHome = support, support, support
	p.CacheHome = filepath.Join(home, "Library", "Caches")
	p.RuntimeDir = support
	p.Videos = filepath.Join(home, "Movies")
	p.DataDirs = []string{"/Library/Application Support", filepath.Join(home, ".local/share")}
	p.ConfigDirs = []string{filepath.Join(home, "Library/Preferences"), "/Library/Application Support", "/Library/Preferences", filepath.Join(home, ".config")}
	p.ApplicationDirs = []string{"/Applications"}
	p.FontDirs = []string{filepath.Join(home, "Library/Fonts"), "/Library/Fonts", "/System/Library/Fonts", "/Network/Library/Fonts"}
	return p
}
func loadPersonalDirectories(*Locations) {}
func (*Locations) searchDirectories()    {}
