//go:build windows

package w32

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
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
