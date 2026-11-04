//go:build !windows

package manifest

import "os"

func replaceFile(source, destination string) error {
	return os.Rename(source, destination)
}
