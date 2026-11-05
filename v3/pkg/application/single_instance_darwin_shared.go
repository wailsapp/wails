//go:build darwin && !ios && !server

package application

import (
	"os"
	"syscall"
)

func (l *darwinLock) release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.file.Name())
		l.file = nil
	}
}

// CreateLockFile tries to create a file with given name and acquire an
// exclusive lock on it. If the file already exists AND is still locked, it will
// fail.
func createLockFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &darwinLock{
		manager: manager,
	}, nil
}

func (l *darwinLock) acquire(uniqueID string) error {
	l.uniqueID = uniqueID
	// NSTemporaryDirectory() returns the correct sandbox-container temp dir
	// regardless of $TMPDIR, which may point outside the sandbox when launched from a terminal.
	lockFilePath := macOSTempDir()
	lockFileName := uniqueID + ".lock"
	var err error
	l.file, err = createLockFile(lockFilePath + lockFileName)
	if err != nil {
		return alreadyRunningError
	}
	return nil
}
