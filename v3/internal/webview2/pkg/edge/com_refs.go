//go:build windows

package edge

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Native code receives Go-allocated handlers as raw addresses, which do not keep
// them alive. handlerRefs records each native reference keyed by the handler
// pointer, so the handler (and the Chromium it points at) stays reachable until
// WebView2 makes its final Release, even after the owning window drops it.
var handlerRefs = struct {
	sync.Mutex
	counts map[unsafe.Pointer]uint32
}{counts: map[unsafe.Pointer]uint32{}}

// Unbalanced Release calls indicate a native reference that was never counted.
var unbalancedHandlerReleases atomic.Uint64

func handlerAddRef(handler unsafe.Pointer) uintptr {
	handlerRefs.Lock()
	defer handlerRefs.Unlock()
	handlerRefs.counts[handler]++
	return uintptr(handlerRefs.counts[handler])
}

func handlerRelease(handler unsafe.Pointer) uintptr {
	handlerRefs.Lock()
	defer handlerRefs.Unlock()
	count := handlerRefs.counts[handler]
	switch count {
	case 0:
		unbalancedHandlerReleases.Add(1)
		return 0
	case 1:
		delete(handlerRefs.counts, handler)
	default:
		handlerRefs.counts[handler] = count - 1
	}
	return uintptr(count - 1)
}

func handlerRefCount(handler unsafe.Pointer) uint32 {
	handlerRefs.Lock()
	defer handlerRefs.Unlock()
	return handlerRefs.counts[handler]
}
