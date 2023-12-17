package dev

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/internal/fs"

	"github.com/fsnotify/fsnotify"
	gitignore "github.com/sabhiram/go-gitignore"
	"github.com/samber/lo"
)

type Watcher interface {
	Add(name string) error
}

// initialiseWatcher creates the project directory watcher that will trigger recompile
func initialiseWatcher(cwd string) (*fsnotify.Watcher, error) {
	// Ignore dot files, node_modules and build directories by default
	ignoreDirs := getIgnoreDirs(cwd)

	// Get all subdirectories
	dirs, err := fs.GetSubdirectories(cwd)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, dir := range processDirectories(dirs.AsSlice(), ignoreDirs) {
		err := watcher.Add(dir)
		if err != nil {
			return nil, err
		}
	}
	return watcher, nil
}

func getIgnoreDirs(cwd string) []string {
	ignoreDirs := []string{filepath.Join(cwd, "build/*"), ".*", "node_modules"}
	baseDir := filepath.Base(cwd)
	// Read .gitignore into ignoreDirs
	f, err := os.Open(filepath.Join(cwd, ".gitignore"))
	if err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line != baseDir {
				ignoreDirs = append(ignoreDirs, line)
			}
		}
	}

	return lo.Uniq(ignoreDirs)
}

func processDirectories(dirs []string, ignoreDirs []string) []string {
	ignorer := gitignore.CompileIgnoreLines(ignoreDirs...)
	return lo.Filter(dirs, func(dir string, _ int) bool {
		return !ignorer.MatchesPath(dir)
	})
}
