//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	gdi32            = syscall.NewLazyDLL("gdi32.dll")
	procGetIconInfo  = user32.NewProc("GetIconInfo")
	procDeleteObject = gdi32.NewProc("DeleteObject")
	procGetObject    = gdi32.NewProc("GetObjectW")
	procGetDIBits    = gdi32.NewProc("GetDIBits")
)

// ICONINFO mirrors the Win32 ICONINFO struct
type ICONINFO struct {
	FIcon    int32
	XHotspot uint32
	YHotspot uint32
	HbmMask  uintptr
	HbmColor uintptr
}

// BITMAP mirrors the Win32 BITMAP struct for GetObject
type BITMAP struct {
	Type       int32
	Width      int32
	Height     int32
	WidthBytes int32
	Planes     uint16
	BitsPixel  uint16
	Bits       uintptr
}

// BITMAPINFOHEADER mirrors the Win32 BITMAPINFOHEADER
type BITMAPINFOHEADER struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

// BITMAPINFO wraps a header plus color table
type BITMAPINFO struct {
	Header BITMAPINFOHEADER
	Colors [1]uint32
}

const (
	BI_RGB         = 0
	DIB_RGB_COLORS = 0
)

type Icon struct {
	handle w32.HICON
}

func NewIconFromFile(path string) (*Icon, error) {
	ico := new(Icon)
	var err error
	if ico.handle = w32.LoadIcon(0, syscall.StringToUTF16Ptr(path)); ico.handle == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from %s", path))
	}
	return ico, err
}

func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (*Icon, error) {
	ico := new(Icon)
	var err error
	if ico.handle = w32.LoadIconWithResourceID(instance, resId); ico.handle == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from resource with id %v", resId))
	}
	return ico, err
}

func ExtractIcon(fileName string, index int) (*Icon, error) {
	ico := new(Icon)
	var err error
	if ico.handle = w32.ExtractIcon(fileName, index); ico.handle == 0 || ico.handle == 1 {
		err = errors.New(fmt.Sprintf("Cannot extract icon from %s at index %v", fileName, index))
	}
	return ico, err
}

// SaveHIconAsPNG extracts the color bitmap from an HICON and writes it to a PNG file.
func SaveHIconAsPNG(hIcon w32.HICON, destPath string) error {
	// 1) Get the ICONINFO structure (which contains two HBITMAPs)
	var ii ICONINFO
	if ret, _, err := procGetIconInfo.Call(
		uintptr(hIcon),
		uintptr(unsafe.Pointer(&ii)),
	); ret == 0 {
		return err
	}
	// Make sure we free the bitmaps when done
	defer procDeleteObject.Call(ii.HbmMask)
	defer procDeleteObject.Call(ii.HbmColor)

	// 2) Render the color bitmap (HbmColor) to RGBA pixels and save
	return saveHBitmapAsPNG(w32.HBITMAP(ii.HbmColor), destPath)
}

func saveHBitmapAsPNG(hBmp w32.HBITMAP, destPath string) error {
	// 1) Fetch the BITMAP header
	var bmp BITMAP
	if ret, _, err := procGetObject.Call(
		uintptr(hBmp),
		unsafe.Sizeof(bmp),
		uintptr(unsafe.Pointer(&bmp)),
	); ret == 0 {
		return err
	}
	w, h := int(bmp.Width), int(bmp.Height)

	// 2) Prepare BITMAPINFO for 32-bit, top-down DIB
	var bmi BITMAPINFO
	bmi.Header = BITMAPINFOHEADER{
		Size:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
		Width:       int32(w),
		Height:      -int32(h), // negative = top-down
		Planes:      1,
		BitCount:    32,
		Compression: BI_RGB,
	}

	// 3) Allocate a buffer and pull the bits
	buf := make([]byte, w*h*4)
	if ret, _, err := procGetDIBits.Call(
		0,
		uintptr(hBmp),
		0,
		uintptr(h),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bmi)),
		DIB_RGB_COLORS,
	); ret == 0 {
		return err
	}

	// 4) Copy into an image.RGBA and write PNG
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, buf)

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func (ic *Icon) Destroy() bool {
	return w32.DestroyIcon(ic.handle)
}

func (ic *Icon) Handle() w32.HICON {
	return ic.handle
}
