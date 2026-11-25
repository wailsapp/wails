//go:build linux && cgo && !android

package webview

import (
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// TestPipeIsCloseOnExec guards against the response pipe leaking into child
// processes: an inherited write end keeps the body open until the child exits.
func TestPipeIsCloseOnExec(t *testing.T) {
	r, w, err := pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	defer syscall.Close(r)
	defer w.Close()

	for name, fd := range map[string]int{"read": r, "write": int(w.Fd())} {
		flags, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
		if err != nil {
			t.Fatalf("fcntl(%s): %v", name, err)
		}
		if flags&unix.FD_CLOEXEC == 0 {
			t.Errorf("%s end of the response pipe is missing FD_CLOEXEC", name)
		}
	}
}
