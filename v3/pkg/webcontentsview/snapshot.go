package webcontentsview

import (
	"sync"
	"sync/atomic"
)

var snapshotCallbacks sync.Map
var snapshotCallbackID uintptr

func registerSnapshotCallback(ch chan string) uintptr {
	id := atomic.AddUintptr(&snapshotCallbackID, 1)
	snapshotCallbacks.Store(id, ch)
	return id
}

func dispatchSnapshotResult(id uintptr, data string) {
	if ch, ok := snapshotCallbacks.LoadAndDelete(id); ok {
		// A native browser can terminate before completing a snapshot. Do not
		// let a late completion block the platform callback after the caller has
		// already given up waiting.
		select {
		case ch.(chan string) <- data:
		default:
		}
	}
}

func unregisterSnapshotCallback(id uintptr) {
	snapshotCallbacks.Delete(id)
}
