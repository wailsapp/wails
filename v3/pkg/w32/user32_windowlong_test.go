//go:build windows

package w32

import (
	"runtime"
	"syscall"
	"testing"
)

func TestWindowLongPtrProceduresResolve(testContext *testing.T) {
	for _, procedure := range []*syscall.LazyProc{procGetWindowLongPtr, procSetWindowLongPtr} {
		if err := procedure.Find(); err != nil {
			testContext.Errorf("%s: %v", procedure.Name, err)
		}
	}
}

func TestWindowLongPtrPreservesUserData(testContext *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	window := CreateWindowEx(0, MustStringToUTF16Ptr("STATIC"), nil, 0, 0, 0, 0, 0, HWND_MESSAGE, 0, 0, nil)
	if window == 0 {
		testContext.Fatal("CreateWindowEx failed")
	}
	defer DestroyWindow(window)

	userData := uintptr(0x89abcdef)
	if previous := SetWindowLongPtr(window, GWLP_USERDATA, userData); previous != 0 {
		testContext.Fatalf("initial user data = %#x, want 0", previous)
	}
	if actual := GetWindowLongPtr(window, GWLP_USERDATA); actual != userData {
		testContext.Fatalf("user data = %#x, want %#x", actual, userData)
	}
	if previous := SetWindowLongPtr(window, GWLP_USERDATA, 0); previous != userData {
		testContext.Fatalf("previous user data = %#x, want %#x", previous, userData)
	}
}
