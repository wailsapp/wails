//go:build windows && !server

package application

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"syscall"
	"testing"
)

const (
	testGR_GDIOBJECTS  = 0
	testGR_USEROBJECTS = 1
)

var (
	testUser32            = syscall.NewLazyDLL("user32.dll")
	testKernel32          = syscall.NewLazyDLL("kernel32.dll")
	testGetGuiResources   = testUser32.NewProc("GetGuiResources")
	testGetCurrentProcess = testKernel32.NewProc("GetCurrentProcess")
)

func getGDIObjectCount(t *testing.T) uint32 {
	t.Helper()
	hProc, _, _ := testGetCurrentProcess.Call()
	n, _, _ := testGetGuiResources.Call(hProc, uintptr(testGR_GDIOBJECTS))
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
