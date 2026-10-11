package commands

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func lockMacOSOutput(path string) (func(), error) {
	lockDir := filepath.Join(os.TempDir(), "wails-macos-output-locks")
	if err := os.MkdirAll(lockDir, 0o700); err != nil {
		return nil, err
	}
	digest := sha256.Sum256([]byte(path))
	file, err := os.OpenFile(filepath.Join(lockDir, fmt.Sprintf("%x", digest[:16])), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX); err != nil {
		file.Close()
		return nil, err
	}
	return func() {
		_ = unix.Flock(int(file.Fd()), unix.LOCK_UN)
		_ = file.Close()
	}, nil
}

func sameMacOSFilesystem(source, target string) (bool, error) {
	var sourceStat, targetStat unix.Stat_t
	if err := unix.Stat(source, &sourceStat); err != nil {
		return false, err
	}
	if err := unix.Stat(target, &targetStat); err != nil {
		return false, err
	}
	return sourceStat.Dev == targetStat.Dev, nil
}
