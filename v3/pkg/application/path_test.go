package application

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestPathEnvironmentSnapshot(t *testing.T) {
	if os.Getenv("WAILS_USERPATHS_TEST_CHILD") == "1" {
		for selector, key := range map[PathType]string{PathDataHome: "XDG_DATA_HOME", PathConfigHome: "XDG_CONFIG_HOME", PathCacheHome: "XDG_CACHE_HOME", PathDownload: "XDG_DOWNLOAD_DIR"} {
			if got, want := Path(selector), os.Getenv(key); got != want {
				t.Errorf("selector %d=%q, want %q", selector, got, want)
			}
		}
		if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
			if got, want := Path(PathDesktop), os.Getenv("WAILS_USERPATHS_TEST_DESKTOP"); got != want {
				t.Errorf("desktop=%q, want %q", got, want)
			}
		}
		if Path(PathType(-1)) != "" || Paths(PathTypes(-1)) != nil {
			t.Error("invalid selectors changed")
		}
		return
	}
	root := t.TempDir()
	config := filepath.Join(root, "config")
	desktop := filepath.Join(root, "custom desktop")
	if err := os.Mkdir(config, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(config, "user-dirs.dirs"), []byte("XDG_DESKTOP_DIR="+strconv.Quote(desktop)+"\nXDG_DOWNLOAD_DIR=\"/ignored/by/environment\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestPathEnvironmentSnapshot$")
	command.Env = append(os.Environ(), "WAILS_USERPATHS_TEST_CHILD=1", "WAILS_USERPATHS_TEST_DESKTOP="+desktop, "XDG_DESKTOP_DIR=", "XDG_CONFIG_HOME="+config, "XDG_DATA_HOME="+filepath.Join(root, "data"), "XDG_CACHE_HOME="+filepath.Join(root, "cache"), "XDG_DOWNLOAD_DIR="+filepath.Join(root, "downloads"))
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("path integration: %v\n%s", err, output)
	}
	for _, name := range []string{"data", "cache", "downloads", "custom desktop"} {
		if _, err := os.Stat(filepath.Join(root, name)); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("resolving %q created a directory or failed unexpectedly: %v", name, err)
		}
	}
}
