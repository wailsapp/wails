//go:build windows

package w32

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"syscall"
	"unsafe"
)

func CreateIconFromResourceEx(presbits uintptr, dwResSize uint32, isIcon bool, version uint32, cxDesired int, cyDesired int, flags uint) (uintptr, error) {
	icon := 0
	if isIcon {
		icon = 1
	}
	r, _, err := procCreateIconFromResourceEx.Call(
		presbits,
		uintptr(dwResSize),
		uintptr(icon),
		uintptr(version),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(flags),
	)

	if r == 0 {
		return 0, err
	}
	return r, nil
}

func isPNG(fileData []byte) bool {
	if len(fileData) < 4 {
		return false
	}
	return string(fileData[:4]) == "\x89PNG"
}

func isICO(fileData []byte) bool {
	if len(fileData) < 4 {
		return false
	}
	return string(fileData[:4]) == "\x00\x00\x01\x00"
}

// CreateSmallHIconFromImage creates a HICON from a PNG or ICO file
func CreateSmallHIconFromImage(fileData []byte) (HICON, error) {
	if len(fileData) < 8 {
		return 0, fmt.Errorf("invalid file format")
	}

	if !isPNG(fileData) && !isICO(fileData) {
		return 0, fmt.Errorf("unsupported file format")
	}
	iconWidth := GetSystemMetrics(SM_CXSMICON)
	iconHeight := GetSystemMetrics(SM_CYSMICON)
	icon, err := CreateIconFromResourceEx(
		uintptr(unsafe.Pointer(&fileData[0])),
		uint32(len(fileData)),
		true,
		0x00030000,
		iconWidth,
		iconHeight,
		LR_DEFAULTSIZE)
	return HICON(icon), err
}

// CreateLargeHIconFromImage creates a HICON from a PNG or ICO file
func CreateLargeHIconFromImage(fileData []byte) (HICON, error) {
	if len(fileData) < 8 {
		return 0, fmt.Errorf("invalid file format")
	}

	if !isPNG(fileData) && !isICO(fileData) {
		return 0, fmt.Errorf("unsupported file format")
	}
	iconWidth := GetSystemMetrics(SM_CXICON)
	iconHeight := GetSystemMetrics(SM_CXICON)
	icon, err := CreateIconFromResourceEx(
		uintptr(unsafe.Pointer(&fileData[0])),
		uint32(len(fileData)),
		true,
		0x00030000,
		iconWidth,
		iconHeight,
		LR_DEFAULTSIZE)
	return HICON(icon), err
}

type ICONINFO struct {
	FIcon    int32
	XHotspot int32
	YHotspot int32
	HbmMask  syscall.Handle
	HbmColor syscall.Handle
}

func SaveHIconAsPNG(hIcon HICON, filePath string) error {
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

func SetWindowIcon(hwnd HWND, icon HICON) {
	SendMessage(hwnd, WM_SETICON, ICON_SMALL, uintptr(icon))
}

func pngToImage(data []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba, nil
}

func SetMenuIcons(parentMenu HMENU, itemID int, unchecked []byte, checked []byte) error {
	if unchecked == nil {
		return fmt.Errorf("invalid unchecked bitmap")
	}
	var err error
	var uncheckedIcon, checkedIcon HBITMAP
	var uncheckedImage, checkedImage *image.RGBA
	uncheckedImage, err = pngToImage(unchecked)
	if err != nil {
		return err
	}
	uncheckedIcon, err = CreateHBITMAPFromImage(uncheckedImage)
	if err != nil {
		return err
	}
	if checked != nil {
		checkedImage, err = pngToImage(checked)
		if err != nil {
			return err
		}
		checkedIcon, err = CreateHBITMAPFromImage(checkedImage)
		if err != nil {
			return err
		}
	} else {
		checkedIcon = uncheckedIcon
	}
	return SetMenuItemBitmaps(parentMenu, uint32(itemID), MF_BYCOMMAND, checkedIcon, uncheckedIcon)
}
