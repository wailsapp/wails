//go:build darwin && !ios

package mac

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrNotInAppBundle is returned when the current executable is not inside a
// macOS application bundle.
var ErrNotInAppBundle = errors.New("mac: executable is not inside an application bundle")

// ResourceFS returns a read-only file system rooted at the current
// application's Contents/Resources directory.
//
// Use ResourceFS when resources should be opened or streamed on demand. For
// small resources that should be read entirely into memory, use LoadResource.
func ResourceFS() (fs.FS, error) {
	resourcesPath, err := ResourcePath()
	if err != nil {
		return nil, err
	}

	return os.DirFS(resourcesPath), nil
}

// ResourcePath returns the path to the current application's
// Contents/Resources directory.
func ResourcePath() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("mac: determine executable path: %w", err)
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return "", fmt.Errorf("mac: resolve executable path: %w", err)
	}

	resourcesPath, ok := resourcePathForExecutable(executable)
	if !ok {
		return "", ErrNotInAppBundle
	}

	return resourcesPath, nil
}

// LoadResource reads a resource from the current application's
// Contents/Resources directory.
//
// Resource names use slash-separated paths relative to Contents/Resources.
// Use ResourceFS to open or stream larger resources instead of loading them
// entirely into memory.
func LoadResource(name string) ([]byte, error) {
	resources, err := ResourceFS()
	if err != nil {
		return nil, err
	}

	return fs.ReadFile(resources, name)
}
