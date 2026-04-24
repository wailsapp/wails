//go:build windows && !server

package application

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"syscall"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

const testGR_GDIOBJECTS = 0

var (
	testUser32            = syscall.NewLazyDLL("user32.dll")
	testKernel32          = syscall.NewLazyDLL("kernel32.dll")
	testGetGuiResources   = testUser32.NewProc("GetGuiResources")
	testGetCurrentProcess = testKernel32.NewProc("GetCurrentProcess")
)

func getGDIObjectCount(t *testing.T) uint32 {
	t.Helper()
	hProc, _, _ := testGetCurrentProcess.Call()
	n, _, callErr := testGetGuiResources.Call(hProc, uintptr(testGR_GDIOBJECTS))
	// GetGuiResources returns 0 on failure; any live Windows process has
	// non-zero GDI objects from the Go runtime alone, so a 0 here means the
	// probe is broken (locked-down container, WinAPI change, bad handle).
	// Fail loudly rather than silently pass the regression guard.
	if n == 0 {
		t.Fatalf("GetGuiResources returned 0 — probe broken, regression guard cannot measure: %v", callErr)
	}
	return uint32(n)
}

func encodeTestBitmapPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := range 16 {
		for x := range 16 {
			img.Set(x, y, color.RGBA{uint8(x * 16), uint8(y * 16), 0, 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode test bitmap: %v", err)
	}
	return buf.Bytes()
}

// TestMenuRuntimeSetBitmapDoesNotLeak exercises the Menu.Update path that
// this PR fixed: runtime MenuItem.SetBitmap allocates an HBITMAP tracked on
// the windowsMenuItem impl, and Update must release it before the rebuild
// replaces item.impl and makes the old one unreachable.
//
// The test allocates a new HBITMAP on every iteration (via SetBitmap after
// the first build installs impls), rebuilds the menu, and measures the
// GR_GDIOBJECTS delta. A fully leaking path would orphan one HBITMAP per
// item per iteration; a healthy path keeps the delta near zero across the
// full run.
func TestMenuRuntimeSetBitmapDoesNotLeak(t *testing.T) {
	bitmap := encodeTestBitmapPNG(t)

	menu := NewMenu()
	a := menu.Add("Item A")
	b := menu.Add("Item B")
	c := menu.Add("Item C")
	t.Cleanup(menu.Destroy)

	// First build so items have impls — subsequent SetBitmap dispatches
	// through windowsMenuItem.setBitmap and allocates impl.bitmap.
	menu.Update()

	const iters = 500
	start := getGDIObjectCount(t)

	for range iters {
		a.SetBitmap(bitmap)
		b.SetBitmap(bitmap)
		c.SetBitmap(bitmap)
		menu.Update()
	}

	end := getGDIObjectCount(t)
	delta := int64(end) - int64(start)

	// Three items × 500 iterations = 1500 handles if nothing is freed.
	// A handful of handles of churn is expected from GDI internals; 50 is
	// well below a single-item regression and well above normal noise.
	const tolerance = 50
	if delta > tolerance {
		t.Errorf("GDI handle delta %d exceeds tolerance %d across %d iterations — runtime SetBitmap HBITMAPs likely leaking", delta, tolerance, iters)
	}
}

// TestWin32MenuChurnDoesNotLeak mirrors the systray-stress "churn" workload
// at the unit level: allocate a fresh Win32Menu per iteration with
// build-time bitmaps, then Destroy. Exercises the pkg/application-level
// path (freeBitmaps releasing p.bitmaps + DestroyMenu tearing down the
// HMENU tree) that the systray's tray.SetMenu flow hits — the harness
// covers this only via a live Windows session.
//
// NewPopupMenu only stores the parent HWND; it's not dereferenced unless
// ShowAt runs, so w32.GetDesktopWindow is a safe stand-in for a real
// window.
func TestWin32MenuChurnDoesNotLeak(t *testing.T) {
	bitmap := encodeTestBitmapPNG(t)
	parent := w32.GetDesktopWindow()

	const iters = 500
	start := getGDIObjectCount(t)

	for range iters {
		menu := NewMenu()
		menu.Add("Item A").SetBitmap(bitmap)
		menu.Add("Item B").SetBitmap(bitmap)
		menu.Add("Item C").SetBitmap(bitmap)

		win32Menu := NewPopupMenu(parent, menu)
		win32Menu.Destroy()
	}

	end := getGDIObjectCount(t)
	delta := int64(end) - int64(start)

	// Each iteration allocates three SetMenuIcons HBITMAPs plus an HMENU;
	// a leak in freeBitmaps or DestroyMenu would grow the handle count by
	// at least one per iteration. 50 is well below 500 while absorbing
	// GDI-internal churn.
	const tolerance = 50
	if delta > tolerance {
		t.Errorf("GDI handle delta %d exceeds tolerance %d across %d iterations — Win32Menu churn likely leaking", delta, tolerance, iters)
	}
}
