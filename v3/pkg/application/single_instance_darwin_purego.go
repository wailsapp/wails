//go:build darwin && purego && !ios && !server

package application

// CGO-free re-implementation of single_instance_darwin.go.
//
// Everything here mirrors the cgo version 1:1 — the file locking, the launchd
// temp-dir handling and the second-instance data plumbing are byte-for-byte the
// same Go. Only the two Objective-C touch points are re-expressed through the
// purego runtime helpers:
//
//   - NSTemporaryDirectory() to locate the sandbox-safe temp dir for the lock
//     file, and
//   - [[NSDistributedNotificationCenter defaultCenter] postNotificationName:...]
//     to hand the launch data to the already-running first instance.

import (
	"os"
	"sync"
	"syscall"

	"github.com/ebitengine/purego"
)

type darwinLock struct {
	file     *os.File
	uniqueID string
	manager  *singleInstanceManager
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

func (l *darwinLock) release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.file.Name())
		l.file = nil
	}
}

func (l *darwinLock) notify(data string) error {
	// Pass message via object, not userInfo: sandboxed apps cannot post
	// distributed notifications with userInfo (the entire post call is blocked
	// by the sandbox).
	withAutoreleasePool(func() {
		center := class("NSDistributedNotificationCenter").send("defaultCenter")
		center.send("postNotificationName:object:userInfo:deliverImmediately:",
			nsString(l.uniqueID), // name
			nsString(data),       // object (the message)
			id(0),                // userInfo == nil
			true,                 // deliverImmediately:YES
		)
	})

	os.Exit(l.manager.options.ExitCode)
	return nil
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

// handleSecondInstanceData is called by the (purego) application delegate when a
// distributed notification carrying the second instance's launch data arrives.
// The cgo build reaches this via an //export'd C shim; here it is a plain Go
// entry point taking the already-decoded message.
func handleSecondInstanceData(secondInstanceMessage string) {
	secondInstanceBuffer <- secondInstanceMessage
}

// --- Objective-C touch points ----------------------------------------------

// NSTemporaryDirectory is a Foundation C function (not an Objective-C method),
// so we bind it directly rather than messaging a class.
var (
	tempDirOnce sync.Once
	nsTempDir   func() uintptr
)

// macOSTempDir returns NSTemporaryDirectory() as a Go string (with its trailing
// slash preserved, matching the cgo implementation).
func macOSTempDir() string {
	tempDirOnce.Do(func() {
		handle, err := purego.Dlopen(frameworkFoundation, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil && handle != 0 {
			purego.RegisterLibFunc(&nsTempDir, handle, "NSTemporaryDirectory")
		}
	})
	if nsTempDir == nil {
		return os.TempDir() + "/"
	}
	var dir string
	withAutoreleasePool(func() {
		dir = id(nsTempDir()).string()
	})
	return dir
}
