//go:build windows

package single_instance

import (
	"os"
	"strconv"
)

// CreateLockFile tries to create a file with given name and acquire an
// exclusive lock on it. If the file already exists AND is still locked, it will
// fail.
func CreateLockFile(filename string, PID int) (*os.File, error) {
	if _, err := os.Stat(filename); err == nil {
		// If the file exists, we first try to remove it
		if err = os.Remove(filename); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	// Write PID to lock file
	_, err = file.WriteString(strconv.Itoa(PID))
	if err != nil {
		return nil, err
	}

	return file, nil
}
