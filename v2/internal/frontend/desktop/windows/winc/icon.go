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
	"image/color"
	"image/png"
	"os"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

var (
	user32                 = syscall.NewLazyDLL("user32.dll")
	gdi32                  = syscall.NewLazyDLL("gdi32.dll")
	procGetIconInfo        = user32.NewProc("GetIconInfo")
	procDeleteObject       = gdi32.NewProc("DeleteObject")
	procGetObject          = gdi32.NewProc("GetObjectW")
	procGetDIBits          = gdi32.NewProc("GetDIBits")
	procCreateCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	procSelectObject       = gdi32.NewProc("SelectObject")
	procDeleteDC           = gdi32.NewProc("DeleteDC")
)

// ICONINFO mirrors the Win32 ICONINFO struct
type ICONINFO struct {
	FIcon    int32
	XHotspot uint32
	YHotspot uint32
	HbmMask  uintptr
	HbmColor uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183376.aspx
type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162938.aspx
type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183375.aspx
type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183371.aspx
type BITMAP struct {
	BmType       int32
	BmWidth      int32
	BmHeight     int32
	BmWidthBytes int32
	BmPlanes     uint16
	BmBitsPixel  uint16
	BmBits       unsafe.Pointer
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

func SaveHIconAsPNG(hIcon w32.HICON, filePath string) error {
	// Get icon info
	var iconInfo ICONINFO
	ret, _, err := procGetIconInfo.Call(
		uintptr(hIcon),
		uintptr(unsafe.Pointer(&iconInfo)),
	)
	if ret == 0 {
		return err
	}
	defer procDeleteObject.Call(uintptr(iconInfo.HbmMask))
	defer procDeleteObject.Call(uintptr(iconInfo.HbmColor))

	// Get bitmap info
	var bmp BITMAP
	ret, _, err = procGetObject.Call(
		uintptr(iconInfo.HbmColor),
		unsafe.Sizeof(bmp),
		uintptr(unsafe.Pointer(&bmp)),
	)
	if ret == 0 {
		return err
	}

	// Create DC
	hdc, _, _ := procCreateCompatibleDC.Call(0)
	if hdc == 0 {
		return syscall.EINVAL
	}
	defer procDeleteDC.Call(hdc)

	// Select bitmap into DC
	oldBitmap, _, _ := procSelectObject.Call(hdc, uintptr(iconInfo.HbmColor))
	defer procSelectObject.Call(hdc, oldBitmap)

	// Prepare bitmap info header
	var bi BITMAPINFO
	bi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bi.BmiHeader))
	bi.BmiHeader.BiWidth = bmp.BmWidth
	bi.BmiHeader.BiHeight = bmp.BmHeight
	bi.BmiHeader.BiPlanes = 1
	bi.BmiHeader.BiBitCount = 32
	bi.BmiHeader.BiCompression = BI_RGB

	// Allocate memory for bitmap bits
	width, height := int(bmp.BmWidth), int(bmp.BmHeight)
	bufferSize := width * height * 4
	bits := make([]byte, bufferSize)

	// Get bitmap bits
	ret, _, err = procGetDIBits.Call(
		hdc,
		uintptr(iconInfo.HbmColor),
		0,
		uintptr(bmp.BmHeight),
		uintptr(unsafe.Pointer(&bits[0])),
		uintptr(unsafe.Pointer(&bi)),
		DIB_RGB_COLORS,
	)
	if ret == 0 {
		return fmt.Errorf("failed to get bitmap bits: %w", err)
	}

	// Create Go image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Convert DIB to RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// DIB is bottom-up, so we need to invert Y
			dibIndex := ((height-1-y)*width + x) * 4

			// BGRA to RGBA
			b := bits[dibIndex]
			g := bits[dibIndex+1]
			r := bits[dibIndex+2]
			a := bits[dibIndex+3]

			// Set pixel in the image
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	// Create output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode and save the image
	return png.Encode(outFile, img)
}

func (ic *Icon) Destroy() bool {
	return w32.DestroyIcon(ic.handle)
}

func (ic *Icon) Handle() w32.HICON {
	return ic.handle
}
