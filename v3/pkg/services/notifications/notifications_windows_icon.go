//go:build windows

package notifications

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

// Windows API constants
const (
	SRCCOPY        = 0x00CC0020
	BI_RGB         = 0
	DIB_RGB_COLORS = 0
)

// Windows structures
type ICONINFO struct {
	FIcon    int32
	XHotspot int32
	YHotspot int32
	HbmMask  syscall.Handle
	HbmColor syscall.Handle
}

type BITMAP struct {
	BmType       int32
	BmWidth      int32
	BmHeight     int32
	BmWidthBytes int32
	BmPlanes     uint16
	BmBitsPixel  uint16
	BmBits       uintptr
}

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

type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors [1]RGBQUAD
}

func saveHIconAsPNG(hIcon w32.HICON, filePath string) error {
	// Load necessary DLLs
	user32 := syscall.NewLazyDLL("user32.dll")
	gdi32 := syscall.NewLazyDLL("gdi32.dll")

	// Get procedures
	getIconInfo := user32.NewProc("GetIconInfo")
	getObject := gdi32.NewProc("GetObjectW")
	createCompatibleDC := gdi32.NewProc("CreateCompatibleDC")
	selectObject := gdi32.NewProc("SelectObject")
	getDIBits := gdi32.NewProc("GetDIBits")
	deleteObject := gdi32.NewProc("DeleteObject")
	deleteDC := gdi32.NewProc("DeleteDC")

	// Get icon info
	var iconInfo ICONINFO
	ret, _, err := getIconInfo.Call(
		uintptr(hIcon),
		uintptr(unsafe.Pointer(&iconInfo)),
	)
	if ret == 0 {
		return err
	}
	defer deleteObject.Call(uintptr(iconInfo.HbmMask))
	defer deleteObject.Call(uintptr(iconInfo.HbmColor))

	// Get bitmap info
	var bmp BITMAP
	ret, _, err = getObject.Call(
		uintptr(iconInfo.HbmColor),
		unsafe.Sizeof(bmp),
		uintptr(unsafe.Pointer(&bmp)),
	)
	if ret == 0 {
		return err
	}

	// Create DC
	hdc, _, _ := createCompatibleDC.Call(0)
	if hdc == 0 {
		return syscall.EINVAL
	}
	defer deleteDC.Call(hdc)

	// Select bitmap into DC
	oldBitmap, _, _ := selectObject.Call(hdc, uintptr(iconInfo.HbmColor))
	defer selectObject.Call(hdc, oldBitmap)

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
	ret, _, err = getDIBits.Call(
		hdc,
		uintptr(iconInfo.HbmColor),
		0,
		uintptr(bmp.BmHeight),
		uintptr(unsafe.Pointer(&bits[0])),
		uintptr(unsafe.Pointer(&bi)),
		DIB_RGB_COLORS,
	)
	if ret == 0 {
		return err
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
