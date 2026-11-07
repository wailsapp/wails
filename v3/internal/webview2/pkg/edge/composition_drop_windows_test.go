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
	lockWebView2Thread(t)
	hwnd := newLiveTestWindow(t)
	user32 := windows.NewLazySystemDLL("user32.dll")
	if ok, _, err := user32.NewProc("MoveWindow").Call(hwnd, 200, 150, 640, 480, 0); ok == 0 {
		t.Fatal(err)
	}
	// OLE supplies screen coordinates; the controller must receive client
	// coordinates even when the window is moved away from the desktop origin.
	client := POINT{X: -123, Y: 456}
	screen := client
	if ok, _, err := user32.NewProc("ClientToScreen").Call(hwnd, uintptr(unsafe.Pointer(&screen))); ok == 0 {
		t.Fatal(err)
	}
	point := client.uintptr()
	screenPoint := screen.uintptr()
	if screenPoint == point {
		t.Fatal("test window has no client-to-screen offset")
	}
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
			args = append(args, keys, screenPoint, uintptr(unsafe.Pointer(&effect)))
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

// Use the real OLE registration API with a stub controller so these tests do
// not require an installed WebView2 runtime. The returned count tracks the
// controller reference acquired specifically for the drop target.
func newDropRegistrationChromium(t *testing.T, hwnd uintptr) (*Chromium, *int) {
	t.Helper()
	refs := new(int)
	controller := &compositionController3{Vtbl: &compositionController3Vtbl{}}
	controller.Vtbl.Release = NewComProc(func(uintptr) uintptr {
		*refs--
		return uintptr(*refs)
	})
	base := &ICoreWebView2CompositionController{Vtbl: &ICoreWebView2CompositionControllerVtbl{}}
	base.Vtbl.QueryInterface = NewComProc(func(_ uintptr, iid *windows.GUID, out **compositionController3) uintptr {
		*out = nil
		if iid.Data1 != 0x9570570e {
			return 0x80004002 // E_NOINTERFACE; no base controller to close.
		}
		*refs++
		*out = controller
		return 0
	})
	base.Vtbl.Release = NewComProc(func(uintptr) uintptr { return 0 })
	e := &Chromium{hwnd: hwnd, compositionController: base}
	t.Cleanup(e.Close)
	return e, refs
}

func TestCompositionDropRegistrationLifecycle(t *testing.T) {
	for _, path := range []string{"shutdown", "close", "recovery"} {
		t.Run(path, func(t *testing.T) {
			lockWebView2Thread(t)
			hwnd := newLiveTestWindow(t)
			e, refs := newDropRegistrationChromium(t, hwnd)
			if err := e.registerCompositionDropTarget(); err != nil {
				t.Fatal(err)
			}
			ref := e.compositionDropTarget.object.Ref()
			if *refs != 1 {
				t.Fatalf("registered controller references = %d, want 1", *refs)
			}
			switch path {
			case "shutdown":
				// WM_CLOSE and direct destruction call this before DestroyWindow.
				e.ShuttingDown()
				e.ShuttingDown()
			case "close":
				e.Close()
				e.Close()
			case "recovery":
				e.ReleaseCompositionResources()
				e.ReleaseCompositionResources()
			}
			if e.compositionDropTarget != nil || *refs != 0 {
				t.Fatalf("cleanup retained target/reference: target=%p refs=%d", e.compositionDropTarget, *refs)
			}
			if combridge.Resolve[compositionDropInterface](ref) != nil {
				t.Fatal("drop target still retained by the COM bridge after cleanup")
			}
			// Check before Close (including the cleanup registered by the helper),
			// so it cannot mask broken shutdown or recovery cleanup.
			// Prove the HWND registration was actually revoked, not just the
			// Go field cleared: a replacement must register on the same HWND.
			replacement, replacementRefs := newDropRegistrationChromium(t, hwnd)
			if err := replacement.registerCompositionDropTarget(); err != nil {
				t.Fatalf("register replacement: %v", err)
			}
			replacement.ShuttingDown()
			if *replacementRefs != 0 {
				t.Fatalf("replacement retained %d controller references", *replacementRefs)
			}
		})
	}
}

func TestCompositionDropRegistrationFailureReleasesController(t *testing.T) {
	lockWebView2Thread(t)
	hwnd := newLiveTestWindow(t)
	owner, ownerRefs := newDropRegistrationChromium(t, hwnd)
	if err := owner.registerCompositionDropTarget(); err != nil {
		t.Fatal(err)
	}
	contender, contenderRefs := newDropRegistrationChromium(t, hwnd)
	if err := contender.registerCompositionDropTarget(); err == nil {
		t.Fatal("registering a second target on the same HWND unexpectedly succeeded")
	}
	contender.Close()
	if *contenderRefs != 0 || contender.compositionDropTarget != nil {
		t.Fatalf("failed registration leaked: refs=%d target=%p", *contenderRefs, contender.compositionDropTarget)
	}
	if *ownerRefs != 1 || combridge.Resolve[compositionDropInterface](owner.compositionDropTarget.object.Ref()) == nil {
		t.Fatal("failed registration released the existing owner's target")
	}
	owner.ShuttingDown()
	if *ownerRefs != 0 {
		t.Fatalf("owner retained %d controller references", *ownerRefs)
	}
}
