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

	"github.com/ebitengine/purego"
)

type darwinLock struct {
	file     *os.File
	uniqueID string
	manager  *singleInstanceManager
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
