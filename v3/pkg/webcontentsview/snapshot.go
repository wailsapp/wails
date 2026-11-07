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
	if ch, ok := snapshotCallbacks.Load(id); ok {
		ch.(chan string) <- data
		snapshotCallbacks.Delete(id)
	}
}
