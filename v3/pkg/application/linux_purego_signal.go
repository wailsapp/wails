//go:build linux && !wailscgo && !android && !server

package application

import (
	"syscall"
	"unsafe"
)

// Use the kernel ABI to preserve the handler, restorer and mask without
// depending on glibc's or musl's different userspace sigaction layouts.
// On Linux amd64 and arm64 these four words are handler, flags, restorer, mask.
type kernelSignalAction [4]uintptr

const signalOnStack = 0x08000000

func fixSignal(sig syscall.Signal) {
	var action kernelSignalAction
	_, _, err := syscall.RawSyscall6(syscall.SYS_RT_SIGACTION, uintptr(sig), 0, uintptr(unsafe.Pointer(&action)), 8, 0, 0)
	if err != 0 {
		return
	}
	action[1] |= signalOnStack
	syscall.RawSyscall6(syscall.SYS_RT_SIGACTION, uintptr(sig), uintptr(unsafe.Pointer(&action)), 0, 8, 0, 0)
}

// WebKit can replace Go's handlers without SA_ONSTACK. Preserve its handlers
// while restoring alternate-stack delivery. SIGUSR1 belongs to JavaScriptCore's
// GC suspension protocol and must not be changed (see #5527).
func installSignalHandlers() {
	for _, sig := range []syscall.Signal{
		syscall.SIGCHLD, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGABRT, syscall.SIGFPE, syscall.SIGTERM, syscall.SIGBUS,
		syscall.SIGSEGV, syscall.SIGXCPU, syscall.SIGXFSZ,
	} {
		fixSignal(sig)
	}
}
