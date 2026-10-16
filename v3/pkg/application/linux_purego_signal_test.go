//go:build linux && !wailscgo && !android && !server

package application

import (
	"syscall"
	"testing"
	"unsafe"
)

func TestSignalFixPreservesNativeHandler(t *testing.T) {
	read := func(sig syscall.Signal) kernelSignalAction {
		t.Helper()
		var action kernelSignalAction
		_, _, err := syscall.RawSyscall6(syscall.SYS_RT_SIGACTION, uintptr(sig), 0, uintptr(unsafe.Pointer(&action)), 8, 0, 0)
		if err != 0 {
			t.Fatal(err)
		}
		return action
	}
	before := read(syscall.SIGTERM)
	usr1 := read(syscall.SIGUSR1)
	installSignalHandlers()
	after := read(syscall.SIGTERM)
	want := before
	want[1] |= signalOnStack
	if after != want {
		t.Fatalf("handler, restorer or mask changed: before %x, after %x", before, after)
	}
	if got := read(syscall.SIGUSR1); got != usr1 {
		t.Fatalf("JavaScriptCore signal modified: before %x, after %x", usr1, got)
	}
}
