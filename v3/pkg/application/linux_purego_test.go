//go:build linux && purego && !gtk3 && !android && !server

package application

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// TestMain loads GTK and WebKitGTK before the tests run, as package init did
// before they loaded on first use: many tests call the bindings directly. A
// failed load is reported by the tests that check linuxLibsErr.
func TestMain(m *testing.M) {
	if os.Getenv("WAILS_TEST_LOAD_FROM_WORKER") == "1" {
		os.Exit(loadFromWorker())
	}
	_ = ensureLinuxLibraries()
	os.Exit(m.Run())
}

// loadFromWorker loads the libraries first from another thread, as an App
// built from a goroutine would, and reports whether the startup thread (which
// TestMain runs on) is still the main thread afterwards.
func loadFromWorker() int {
	result := make(chan int)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		if err := ensureLinuxLibraries(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			result <- 4 // the libraries did not load
			return
		}
		if isOnMainThread() {
			result <- 2 // the worker was taken for the main thread
			return
		}
		result <- 0
	}()
	if r := <-result; r != 0 {
		return r
	}
	if !isOnMainThread() {
		return 3 // the startup thread is no longer the main thread
	}
	return 0
}

// TestMainThreadIsTheStartupThread runs loadFromWorker in a child process,
// where nothing has loaded the libraries yet.
func TestMainThreadIsTheStartupThread(t *testing.T) {
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), "WAILS_TEST_LOAD_FROM_WORKER=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("loading the libraries from a worker moved the main thread: %v\n%s", err, out)
	}
}

// TestPuregoNativeLibraries checks the GTK 4 and WebKitGTK ABIs were loaded.
func TestPuregoNativeLibraries(t *testing.T) {
	// The libraries load on first use, not at init.
	if err := ensureLinuxLibraries(); err != nil {
		t.Fatal(err)
	}
	if gtk_get_major_version() != 4 {
		t.Fatal("wrong GTK ABI loaded")
	}
	if webkit_get_major_version() < 2 {
		t.Fatal("invalid WebKit version")
	}
}
