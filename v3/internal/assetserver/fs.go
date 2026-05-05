package assetserver

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
	if stat != nil {
		return ".", nil
	}
	var indexFiles []string
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, file) {
			indexFiles = append(indexFiles, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if len(indexFiles) > 1 {
		selected := indexFiles[0]
		for _, f := range indexFiles {
			if len(f) < len(selected) {
				selected = f
			}
		}
		path, _ := filepath.Split(selected)
		return path, nil
	}
	if len(indexFiles) > 0 {
		path, _ := filepath.Split(indexFiles[0])
		return path, nil
	}
	return "", fmt.Errorf("%s: %w", file, os.ErrNotExist)
}
