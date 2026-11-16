// Package userpaths resolves standard application and personal directories.
package userpaths

import (
	"os"
	"os/user"
	"path/filepath"
)

type Locations struct {
	Home            string
	DataHome        string
	ConfigHome      string
	StateHome       string
	CacheHome       string
	RuntimeDir      string
	Desktop         string
	Download        string
	Documents       string
	Music           string
	Pictures        string
	Videos          string
	Templates       string
	PublicShare     string
	DataDirs        []string
	ConfigDirs      []string
	FontDirs        []string
	ApplicationDirs []string
}

// Resolve takes a snapshot of the current user's directories without creating them.
func Resolve() Locations {
	home, _ := os.UserHomeDir()
	if home == "" {
		if current, err := user.Current(); err == nil {
			home = current.HomeDir
		}
	}
	p := platformDefaults(home)
	p.applyEnvironment(os.Getenv)
	loadPersonalDirectories(&p)
	p.searchDirectories()
	return p
}

func personalDefaults(home string) Locations {
	return Locations{
		Home:        home,
		Desktop:     filepath.Join(home, "Desktop"),
		Download:    filepath.Join(home, "Downloads"),
		Documents:   filepath.Join(home, "Documents"),
		Music:       filepath.Join(home, "Music"),
		Pictures:    filepath.Join(home, "Pictures"),
		Videos:      filepath.Join(home, "Videos"),
		Templates:   filepath.Join(home, "Templates"),
		PublicShare: filepath.Join(home, "Public"),
	}
}

type setting struct {
	key  string
	path *string
}

func (p *Locations) personalSettings() []setting {
	return []setting{
		{"XDG_DESKTOP_DIR", &p.Desktop},
		{"XDG_DOWNLOAD_DIR", &p.Download},
		{"XDG_DOCUMENTS_DIR", &p.Documents},
		{"XDG_MUSIC_DIR", &p.Music},
		{"XDG_PICTURES_DIR", &p.Pictures},
		{"XDG_VIDEOS_DIR", &p.Videos},
		{"XDG_TEMPLATES_DIR", &p.Templates},
		{"XDG_PUBLICSHARE_DIR", &p.PublicShare},
	}
}

func (p *Locations) applyEnvironment(getenv func(string) string) {
	fields := []setting{
		{"XDG_DATA_HOME", &p.DataHome},
		{"XDG_CONFIG_HOME", &p.ConfigHome},
		{"XDG_STATE_HOME", &p.StateHome},
		{"XDG_CACHE_HOME", &p.CacheHome},
		{"XDG_RUNTIME_DIR", &p.RuntimeDir},
	}
	fields = append(fields, p.personalSettings()...)
	for _, field := range fields {
		if value := getenv(field.key); filepath.IsAbs(value) {
			*field.path = value
		}
	}
	for _, list := range []struct {
		key    string
		values *[]string
	}{{"XDG_DATA_DIRS", &p.DataDirs},
		{"XDG_CONFIG_DIRS", &p.ConfigDirs}} {
		if value := getenv(list.key); value != "" {
			if paths := uniqueAbsolute(filepath.SplitList(value)); len(paths) > 0 {
				*list.values = paths
			}
		}
	}
}

func uniqueAbsolute(paths []string) []string {
	result := make([]string, 0, len(paths))
	seen := make(map[string]bool, len(paths))
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			continue
		}
		if !seen[path] {
			result = append(result, path)
			seen[path] = true
		}
	}
	return result
}
