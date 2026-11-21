//go:build windows

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	gdi32                  = windows.NewLazySystemDLL("gdi32.dll")
	procGetDC              = user32.NewProc("GetDC")
	procReleaseDC          = user32.NewProc("ReleaseDC")
	procGetSystemMetrics   = user32.NewProc("GetSystemMetrics")
	procCreateCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	procCreateDIBSection   = gdi32.NewProc("CreateDIBSection")
	procGdiFlush           = gdi32.NewProc("GdiFlush")
	procSelectObject       = gdi32.NewProc("SelectObject")
	procBitBlt             = gdi32.NewProc("BitBlt")
	procDeleteObject       = gdi32.NewProc("DeleteObject")
	procDeleteDC           = gdi32.NewProc("DeleteDC")
)

const (
	smXVirtualScreen  = 76
	smYVirtualScreen  = 77
	smCXVirtualScreen = 78
	smCYVirtualScreen = 79
	srcCopy           = 0x00CC0020
	dibRGBColors      = 0
)

type bitmapInfoHeader struct {
	size                      uint32
	width, height             int32
	planes, bitCount          uint16
	compression, sizeImage    uint32
	xPelsPerMeter, yPelsPerMe int32
	clrUsed, clrImportant     uint32
}

func metric(i uintptr) int32 {
	v, _, _ := procGetSystemMetrics.Call(i)
	return int32(v)
}

// screenshot writes the whole virtual screen to file as PNG, with GDI.
func screenshot(file string) error {
	x, y := metric(smXVirtualScreen), metric(smYVirtualScreen)
	w, h := metric(smCXVirtualScreen), metric(smCYVirtualScreen)
	if w <= 0 || h <= 0 {
		return fmt.Errorf("no screen (%dx%d)", w, h)
	}
	screen, _, _ := procGetDC.Call(0)
	if screen == 0 {
		return fmt.Errorf("GetDC failed")
	}
	defer procReleaseDC.Call(0, screen)
	mem, _, _ := procCreateCompatibleDC.Call(screen)
	if mem == 0 {
		return fmt.Errorf("CreateCompatibleDC failed")
	}
	defer procDeleteDC.Call(mem)
	// A top-down (negative height) 32-bit DIB section: BitBlt writes the
	// screen into memory this process can read directly, as BGRA.
	header := bitmapInfoHeader{width: w, height: -h, planes: 1, bitCount: 32}
	header.size = uint32(unsafe.Sizeof(header))
	var bits unsafe.Pointer
	bmp, _, err := procCreateDIBSection.Call(screen, uintptr(unsafe.Pointer(&header)), dibRGBColors,
		uintptr(unsafe.Pointer(&bits)), 0, 0)
	if bmp == 0 || bits == nil {
		return fmt.Errorf("CreateDIBSection: %w", err)
	}
	defer procDeleteObject.Call(bmp)
	old, _, _ := procSelectObject.Call(mem, bmp)
	ok, _, err := procBitBlt.Call(mem, 0, 0, uintptr(w), uintptr(h), screen, uintptr(x), uintptr(y), srcCopy)
	procGdiFlush.Call()
	procSelectObject.Call(mem, old)
	if ok == 0 {
		return fmt.Errorf("BitBlt: %w", err)
	}
	img := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	copy(img.Pix, unsafe.Slice((*byte)(bits), len(img.Pix)))
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i], img.Pix[i+2], img.Pix[i+3] = img.Pix[i+2], img.Pix[i], 0xFF
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
