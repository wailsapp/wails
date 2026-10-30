package userpaths

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestEnvironmentOverrides(t *testing.T) {
	home := t.TempDir()
	p := personalDefaults(home)
	p.DataHome = filepath.Join(home, "default-data")
	p.DataDirs = []string{home}
	data := filepath.Join(home, "data")
	desktop := filepath.Join(home, "desktop")
	env := map[string]string{"XDG_DATA_HOME": data, "XDG_DESKTOP_DIR": desktop, "XDG_CACHE_HOME": "relative"}
	p.applyEnvironment(func(key string) string { return env[key] })
	if p.DataHome != data || p.Desktop != desktop || p.CacheHome != "" {
		t.Fatalf("invalid overrides: %+v", p)
	}
}

func TestPersonalDirectoryConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("user-dirs.dirs is read only on Unix")
	}
	home := t.TempDir()
	p := personalDefaults(home)
	config := `# personal directories
XDG_DESKTOP_DIR="$HOME/Desk Test"
XDG_DOWNLOAD_DIR="/custom/downloads"
XDG_DOCUMENTS_DIR="$HOME"
XDG_MUSIC_DIR="relative"
XDG_PICTURES_DIR="$HOME/Photo\"s"
XDG_VIDEOS_DIR="/literal/$(echo do-not-run)"
XDG_TEMPLATES_DIR="/not-closed
UNRELATED_KEY="/ignored"
`
	override := filepath.Join(home, "override")
	getenv := func(key string) string {
		if key == "XDG_DOWNLOAD_DIR" {
			return override
		}
		return ""
	}
	p.applyEnvironment(getenv)
	p.readPersonalDirectories(strings.NewReader(config), getenv)
	if p.Desktop != filepath.Join(home, "Desk Test") || p.Documents != home || p.Music != filepath.Join(home, "Music") || p.Pictures != filepath.Join(home, "Photo\"s") || p.Videos != "/literal/$(echo do-not-run)" || p.Templates != filepath.Join(home, "Templates") || p.Download != override {
		t.Fatalf("unexpected personal directories: %+v", p)
	}
}

func TestUniqueAbsolutePaths(t *testing.T) {
	home := t.TempDir()
	paths := uniqueAbsolute([]string{home, filepath.Join(home, "child", ".."), "relative", "", filepath.Join(home, "other")})
	if want := []string{home, filepath.Join(home, "other")}; !reflect.DeepEqual(paths, want) {
		t.Fatalf("paths=%v,want %v", paths, want)
	}
}

func BenchmarkResolve(b *testing.B) {
	b.Setenv("XDG_CONFIG_HOME", b.TempDir())
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = Resolve()
	}
}

func FuzzPersonalValue(f *testing.F) {
	f.Add(`"$HOME/Desktop"`)
	f.Add(`"/tmp/a\\b"`)
	f.Add(`"relative"`)
	home := f.TempDir()
	f.Fuzz(func(t *testing.T, value string) {
		if len(value) > 64<<10 {
			t.Skip()
		}
		path, ok := personalValue(value, home)
		if ok && !filepath.IsAbs(path) {
			t.Fatalf("accepted relative path %q", path)
		}
	})
}
