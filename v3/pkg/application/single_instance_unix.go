//go:build unix

package application

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

type unixLock struct {
	file     *os.File
	uniqueID string
}

func newPlatformLock(_ *singleInstanceManager) (platformLock, error) {
	return &unixLock{}, nil
}

func (l *unixLock) acquire(uniqueID string) error {
	if uniqueID == "" {
		return fmt.Errorf("UniqueID is required for single instance lock")
	}

	l.uniqueID = uniqueID
	lockPath := getLockPath(uniqueID)
	file, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open lock file: %w", err)
	}

	// Try to acquire an exclusive lock
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		if strings.Contains(err.Error(), "resource temporarily unavailable") {
			return fmt.Errorf("another instance is already running")
		}
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	// Write PID to lock file
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := file.Truncate(0); err != nil {
		syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		file.Close()
		return fmt.Errorf("failed to truncate lock file: %w", err)
	}
	if _, err := file.WriteString(pid); err != nil {
		syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		file.Close()
		return fmt.Errorf("failed to write PID to lock file: %w", err)
	}

	l.file = file
	return nil
}

func (l *unixLock) release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.file.Name())
		l.file = nil
	}
}

func (l *unixLock) notify(data string) error {
	// For Unix systems, we need to use some form of IPC
	lockPath := getLockPath(l.uniqueID)
	pidBytes, err := os.ReadFile(lockPath)
	if err != nil {
		return fmt.Errorf("failed to read PID from lock file: %w", err)
	}

	pid := strings.TrimSpace(string(pidBytes))
	process, err := os.FindProcess(parseInt(pid))
	if err != nil {
		return fmt.Errorf("failed to find first instance process: %w", err)
	}

	// Send SIGUSR1 to notify the first instance
	err = process.Signal(syscall.SIGUSR1)
	if err != nil {
		return fmt.Errorf("failed to signal first instance: %w", err)
	}

	// Write data to a temporary file that the first instance can read
	dataFile := lockPath + ".data"
	err = os.WriteFile(dataFile, []byte(data), 0600)
	if err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
}

func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
