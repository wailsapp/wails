//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0

#include <stdio.h>
#include "gtk/gtk.h"

extern gboolean invokeCallbacks(void *);

static inline void triggerInvokesOnMainThread() {
    g_idle_add((GSourceFunc)invokeCallbacks, NULL);
}
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	m         sync.Mutex
	mainTid   int
	dispatchq []func()
)

func invokeOnMainThread(f func()) {
	if tryInvokeOnCurrentGoRoutine(f) {
		return
	}

	m.Lock()
	dispatchq = append(dispatchq, f)
	m.Unlock()

	C.triggerInvokesOnMainThread()
}

func tryInvokeOnCurrentGoRoutine(f func()) bool {
	m.Lock()
	mainThreadID := mainTid
	m.Unlock()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if mainThreadID != unix.Gettid() {
		return false
	}
	f()
	return true
}

//export invokeCallbacks
func invokeCallbacks(_ unsafe.Pointer) C.gboolean {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	m.Lock()
	if mainTid == 0 {
		mainTid = unix.Gettid()
	}

	q := append([]func(){}, dispatchq...)
	dispatchq = []func(){}
	m.Unlock()

	for _, v := range q {
		v()
	}
	return C.G_SOURCE_REMOVE
}
