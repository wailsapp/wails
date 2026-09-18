//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#include <stdlib.h>
#include "browser_manager_darwin.h"
*/
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// completionWaiters pairs asynchronous NSWorkspace calls (OpenWith,
// SetDefaultHandler) with the goroutine waiting for their completion
// handler. wailsCompletionCallback is the single entry point from
// Objective-C so no drain loop is needed: the callback only closes a
// channel.
var (
	completionWaiters sync.Map
	completionNextID  atomic.Uint64
)

const workspaceCompletionTimeout = 30 * time.Second

func newCompletionWaiter() (uint64, chan error) {
	id := completionNextID.Add(1)
	ch := make(chan error, 1)
	completionWaiters.Store(id, ch)
	return id, ch
}

func awaitCompletion(id uint64, ch chan error, what string) error {
	defer completionWaiters.Delete(id)
	select {
	case err := <-ch:
		return err
	case <-time.After(workspaceCompletionTimeout):
		return fmt.Errorf("%s: timed out waiting for the system", what)
	}
}

//export wailsCompletionCallback
func wailsCompletionCallback(id C.ulonglong, message *C.char) {
	value, ok := completionWaiters.Load(uint64(id))
	if !ok {
		return
	}
	ch := value.(chan error)
	if message == nil {
		ch <- nil
		return
	}
	ch <- errors.New(C.GoString(message))
}

func platformOpenWith(path string, app string) error {
	id, ch := newCompletionWaiter()
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cApp := C.CString(app)
	defer C.free(unsafe.Pointer(cApp))
	InvokeSync(func() {
		C.wailsWorkspaceOpenWith(C.ulonglong(id), cPath, cApp)
	})
	if err := awaitCompletion(id, ch, "browser: OpenWith"); err != nil {
		if err.Error() == "application not found: "+app {
			return fmt.Errorf("%w: %s", ErrApplicationNotFound, app)
		}
		return fmt.Errorf("browser: OpenWith %q: %w", app, err)
	}
	return nil
}

func platformApplicationsForFile(path string) []AppInfo {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	raw := C.wailsWorkspaceApplicationsForFile(cPath)
	defer C.free(unsafe.Pointer(raw))
	var apps []AppInfo
	if err := json.Unmarshal([]byte(C.GoString(raw)), &apps); err != nil {
		globalApplication.error("browser: cannot decode applications for %q: %v", path, err)
		return nil
	}
	return apps
}

func platformActivateApplication(bundleID string) error {
	cBundleID := C.CString(bundleID)
	defer C.free(unsafe.Pointer(cBundleID))
	message := InvokeSyncWithResult(func() string {
		raw := C.wailsWorkspaceActivateApplication(cBundleID)
		if raw == nil {
			return ""
		}
		defer C.free(unsafe.Pointer(raw))
		return C.GoString(raw)
	})
	switch message {
	case "":
		return nil
	case "not running":
		return fmt.Errorf("%w: %s", ErrApplicationNotRunning, bundleID)
	default:
		return fmt.Errorf("browser: ActivateApplication %q: %s", bundleID, message)
	}
}
