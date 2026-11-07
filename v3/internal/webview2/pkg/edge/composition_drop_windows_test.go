//go:build windows

package edge

import (
	"runtime"
	"syscall"
	"testing"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/webview2/pkg/combridge"
	"golang.org/x/sys/windows"
)

// Exercise the actual COM vtable trampolines: POINTL packing and the different
// IDropTarget/controller method orders must work on both amd64 and arm64.
func TestCompositionDropForwarding(t *testing.T) {
	hwnd, _, _ := windows.NewLazySystemDLL("user32.dll").NewProc("GetDesktopWindow").Call()
	point := POINT{X: -123, Y: 456}.uintptr()
	const data = uintptr(0x12345678)
	const keys = uintptr(9)
	var calls []string
	check := func(name string, gotData, gotKeys, gotPoint, effect uintptr) uintptr {
		calls = append(calls, name)
		if gotData != data || gotKeys != keys || gotPoint != point {
			t.Errorf("%s: data=%x keys=%x point=%x", name, gotData, gotKeys, gotPoint)
		}
		p := (*uint32)(unsafe.Pointer(effect))
		if *p != 7 {
			t.Errorf("%s lost incoming allowed effects: %d", name, *p)
		}
		*p = 1
		return 1 // S_FALSE is returned unchanged.
	}
	controller := &compositionController3{Vtbl: &compositionController3Vtbl{
		DragEnter: NewComProc(func(_, d, k, p, e uintptr) uintptr { return check("enter", d, k, p, e) }),
		DragOver:  NewComProc(func(_, k, p, e uintptr) uintptr { return check("over", data, k, p, e) }),
		DragLeave: NewComProc(func(uintptr) uintptr { calls = append(calls, "leave"); return 1 }),
		Drop:      NewComProc(func(_, d, k, p, e uintptr) uintptr { return check("drop", d, k, p, e) }),
	}}
	target := &compositionDropTarget{hwnd: hwnd, controller: controller}
	object := combridge.New[compositionDropInterface](target)
	defer object.Close()
	ref := object.Ref()
	vt := unsafe.Slice(*(**uintptr)(unsafe.Pointer(ref)), 7)
	for _, index := range []int{3, 4, 5, 6} {
		effect := uint32(7)
		// The uintptr is stored in args before SyscallN. Pin the Go allocation
		// across callbacks, which may grow and relocate the caller's stack.
		var pin runtime.Pinner
		pin.Pin(&effect)
		defer pin.Unpin()
		args := []uintptr{ref}
		if index == 3 || index == 6 {
			args = append(args, data)
		}
		if index != 5 {
			args = append(args, keys, point, uintptr(unsafe.Pointer(&effect)))
		}
		hr, _, _ := syscall.SyscallN(vt[index], args...)
		if hr != 1 {
			t.Errorf("slot %d HRESULT=%x", index, hr)
		}
		if index != 5 && effect != 1 {
			t.Errorf("slot %d did not return effect", index)
		}
	}
	if len(calls) != 4 || calls[0] != "enter" || calls[1] != "over" || calls[2] != "leave" || calls[3] != "drop" {
		t.Fatalf("wrong callback order: %v", calls)
	}
}
