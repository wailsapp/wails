//go:build !windows && !darwin

package userpaths

import (
	"os"
	"path/filepath"
	"strconv"
)

func platformDefaults(home string) Locations {
	p := personalDefaults(home)
	p.DataHome = filepath.Join(home, ".local", "share")
	p.ConfigHome = filepath.Join(home, ".config")
	p.StateHome = filepath.Join(home, ".local", "state")
	p.CacheHome = filepath.Join(home, ".cache")
	p.RuntimeDir = filepath.Join("/run/user", strconv.Itoa(os.Getuid()))
	p.DataDirs = []string{"/usr/local/share", "/usr/share"}
	p.ConfigDirs = []string{"/etc/xdg"}
	return p
}

func loadPersonalDirectories(p *Locations) {
	f, err := os.Open(filepath.Join(p.ConfigHome, "user-dirs.dirs"))
	if err != nil {
		return
	}
	defer f.Close()
	p.readPersonalDirectories(f, os.Getenv)
}

func (p *Locations) searchDirectories() {
	apps := []string{filepath.Join(p.DataHome, "applications"), filepath.Join(p.Home, ".local/share/applications"), "/usr/local/share/applications", "/usr/share/applications"}
	fonts := []string{filepath.Join(p.DataHome, "fonts"), filepath.Join(p.Home, ".fonts"), filepath.Join(p.Home, ".local/share/fonts"), "/usr/local/share/fonts", "/usr/share/fonts"}
	for _, base := range p.DataDirs {
		apps = append(apps, filepath.Join(base, "applications"))
		fonts = append(fonts, filepath.Join(base, "fonts"))
	}
	p.ApplicationDirs, p.FontDirs = uniqueAbsolute(apps), uniqueAbsolute(fonts)
}
