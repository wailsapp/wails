//go:build windows

package w32_test

import (
	"syscall"
	"testing"

	"github.com/matryer/is"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

// Regression tests for the Windows systray SetMenu crash family.
//
// Before the fix in user32.go, DestroyMenu was calling
// procDestroyMenu.Call(1, uintptr(hMenu), 0, 0), which on x86-64 Windows
// put the literal 1 in RCX (the handle slot) and the real handle in RDX
// (ignored). Every call returned FALSE without freeing anything, which
// meant Win32Menu.Destroy and the leak-fix commits in this series could
// not actually release HMENU handles.
//
// These tests fail on the pre-fix syscall signature:
//   - TestDestroyMenuReturnsTrueForValidHandle fails because the syscall
//     sees handle=1 and returns FALSE.
//   - TestMenuHandleLifecycleDoesNotLeak fails its DestroyMenu assertion
//     on the first iteration for the same reason.
//
// They run on Windows CI via the //go:build windows constraint.

func TestDestroyMenuReturnsTrueForValidHandle(t *testing.T) {
	is := is.New(t)
	hmenu := w32.CreatePopupMenu()
	is.True(hmenu != 0) // CreatePopupMenu returned a valid handle
	is.True(w32.DestroyMenu(w32.HMENU(hmenu)))
}

func TestDestroyMenuReturnsFalseForZeroHandle(t *testing.T) {
	is := is.New(t)
	is.Equal(w32.DestroyMenu(0), false)
}

// TestMenuHandleLifecycleDoesNotLeak exercises the Create+Append+Destroy
// cycle many times and asserts that USER-object handle usage does not
// grow unboundedly. Under the pre-fix DestroyMenu, handle count grows by
// at least one per iteration and the test times out / exhausts the
// process quota; post-fix, the delta stays near zero across thousands
// of iterations.
func TestMenuHandleLifecycleDoesNotLeak(t *testing.T) {
	is := is.New(t)
	const iterations = 2000
	const tolerance = 50 // allow small transient variance

	start := getUserObjectCount(t)

	for i := 0; i < iterations; i++ {
		hmenu := w32.CreatePopupMenu()
		if hmenu == 0 {
			t.Fatalf("CreatePopupMenu failed at iter %d", i)
		}
		label := w32.MustStringToUTF16Ptr("item")
		if !w32.AppendMenu(w32.HMENU(hmenu), w32.MF_STRING, 1, label) {
			t.Fatalf("AppendMenu failed at iter %d", i)
		}
		if !w32.DestroyMenu(w32.HMENU(hmenu)) {
			t.Fatalf("DestroyMenu failed at iter %d — fix 5 regression?", i)
		}
	}

	end := getUserObjectCount(t)
	delta := int64(end) - int64(start)
	t.Logf("USER object delta over %d iterations: start=%d end=%d delta=%d", iterations, start, end, delta)
	is.True(delta < tolerance) // no runaway handle growth across the Create+Append+Destroy cycle
}

// getUserObjectCount returns the current process's USER object count via
// GetGuiResources. Wraps the syscall inline so the test does not depend
// on internal w32 helpers being exported.
func getUserObjectCount(t *testing.T) uint32 {
	t.Helper()
	const GR_USEROBJECTS = 1
	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getGuiResources := user32.NewProc("GetGuiResources")
	getCurrentProcess := kernel32.NewProc("GetCurrentProcess")
	hProc, _, _ := getCurrentProcess.Call()
	ret, _, _ := getGuiResources.Call(hProc, uintptr(GR_USEROBJECTS))
	return uint32(ret)
}
