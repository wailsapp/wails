package assetserver

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// findEmbedRootPath finds the root path in the embed FS. It's the directory which contains all the files.
func findEmbedRootPath(fileSystem embed.FS) (string, error) {
	stopErr := errors.New("files or multiple dirs found")

	fPath := ""
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			fPath = path
			if entries, dErr := fs.ReadDir(fileSystem, path); dErr != nil {
				return dErr
			} else if len(entries) <= 1 {
				return nil
			}
		}

		return stopErr
	})

	if err != nil && !errors.Is(err, stopErr) {
		return "", err
	}

	return fPath, nil
}

func findPathToFile(fileSystem fs.FS, file string) (string, error) {
	stat, _ := fs.Stat(fileSystem, file)
	if stat != nil && !stat.IsDir() {
		return ".", nil
	}
	var selected string
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == file && (selected == "" || len(path) < len(selected)) {
			selected = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if selected != "" {
		path, _ := filepath.Split(selected)
		return path, nil
	}
	return "", fmt.Errorf("%s: %w", file, os.ErrNotExist)
}
