package commands

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"reflect"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestManifestDevWindowsConsoleInput(t *testing.T) {
	if os.Getenv("WAKE_CONSOLE_TEST_HELPER") != "1" {
		executable, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, executable, "-test.run=^TestManifestDevWindowsConsoleInput$", "-test.v")
		cmd.Env = append(os.Environ(), "WAKE_CONSOLE_TEST_HELPER=1")
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.DETACHED_PROCESS}
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("native console helper: %v\n%s", err, output)
		}
		return
	}
	kernel := windows.NewLazySystemDLL("kernel32.dll")
	if ok, _, err := kernel.NewProc("AllocConsole").Call(); ok == 0 {
		t.Fatal(err)
	}
	defer kernel.NewProc("FreeConsole").Call()
	input, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	handle := windows.Handle(input.Fd())
	if err := windows.FlushConsoleInputBuffer(handle); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	var lines []string
	go func() {
		done <- readManifestDevInput(ctx, input, func(line string) bool {
			lines = append(lines, line)
			if len(lines) == 2 {
				cancel()
			}
			return true
		})
	}()
	events := []manifestConsoleEvent{{Type: 1, KeyDown: 1, Repeat: 1, Character: 'r'}, {Type: 1, KeyDown: 1, Repeat: 1, Character: '\r'}, {Type: 1, KeyDown: 1, Repeat: 1, Character: 'q'}, {Type: 1, KeyDown: 1, Repeat: 1, Character: '\r'}}
	var written uint32
	if ok, _, err := kernel.NewProc("WriteConsoleInputW").Call(uintptr(handle), uintptr(unsafe.Pointer(&events[0])), uintptr(len(events)), uintptr(unsafe.Pointer(&written))); ok == 0 || written != uint32(len(events)) {
		t.Fatalf("console write %d: %v", written, err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("console input/cancellation blocked")
	}
	if !reflect.DeepEqual(lines, []string{"r", "q"}) {
		t.Fatal(lines)
	}
	// A new idle reader must stop without waiting for another console event.
	idle, cancelIdle := context.WithCancel(context.Background())
	go func() { done <- readManifestDevInput(idle, input, func(string) bool { return true }) }()
	cancelIdle()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("idle console reader blocked")
	}
}
