//go:build windows

package w32

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"runtime"
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

const (
	icoHeaderSize       = 6
	icoDirEntrySize     = 16
	icoWidthOffset      = 0
	icoHeightOffset     = 1
	icoColorCountOffset = 2
	icoBitCountOffset   = 6
	icoDataSizeOffset   = 8
	icoDataOffset       = 12
)

type icoImageEntry struct {
	width      int
	height     int
	colorCount byte
	bitCount   uint16
	offset     uint32
	size       uint32
}

// selectICOImage returns the image resource inside an ICO container that best
// matches the requested dimensions
func selectICOImage(fileData []byte, desiredWidth, desiredHeight, displayBitDepth int) ([]byte, error) {
	if len(fileData) < icoHeaderSize || !isICO(fileData) {
		return nil, fmt.Errorf("invalid ICO header")
	}

	imageCount := int(binary.LittleEndian.Uint16(fileData[4:6]))
	if imageCount == 0 {
		return nil, fmt.Errorf("ICO contains no images")
	}
	if imageCount > (len(fileData)-icoHeaderSize)/icoDirEntrySize {
		return nil, fmt.Errorf("invalid ICO directory")
	}
	directoryEnd := icoHeaderSize + imageCount*icoDirEntrySize

	entries := make([]icoImageEntry, 0, imageCount)
	for index := range imageCount {
		entryStart := icoHeaderSize + index*icoDirEntrySize
		entryData := fileData[entryStart : entryStart+icoDirEntrySize]

		width := int(entryData[icoWidthOffset])
		if width == 0 {
			width = 256
		}
		height := int(entryData[icoHeightOffset])
		if height == 0 {
			height = 256
		}

		entry := icoImageEntry{
			width:      width,
			height:     height,
			colorCount: entryData[icoColorCountOffset],
			bitCount:   binary.LittleEndian.Uint16(entryData[icoBitCountOffset : icoBitCountOffset+2]),
			size:       binary.LittleEndian.Uint32(entryData[icoDataSizeOffset : icoDataSizeOffset+4]),
			offset:     binary.LittleEndian.Uint32(entryData[icoDataOffset : icoDataOffset+4]),
		}
		end := uint64(entry.offset) + uint64(entry.size)
		if entry.size == 0 || uint64(entry.offset) < uint64(directoryEnd) || end > uint64(len(fileData)) {
			return nil, fmt.Errorf("ICO image %d is outside the file", index)
		}
		entries = append(entries, entry)
	}

	best := entries[0]
	for _, candidate := range entries[1:] {
		candidateFits := candidate.width <= desiredWidth && candidate.height <= desiredHeight
		bestFits := best.width <= desiredWidth && best.height <= desiredHeight
		candidateDistance := abs(candidate.width-desiredWidth) + abs(candidate.height-desiredHeight)
		bestDistance := abs(best.width-desiredWidth) + abs(best.height-desiredHeight)

		if candidateFits && !bestFits ||
			candidateFits == bestFits && (candidateDistance < bestDistance ||
				candidateDistance == bestDistance && betterBitDepth(
					candidate.effectiveBitDepth(displayBitDepth),
					best.effectiveBitDepth(displayBitDepth),
					displayBitDepth,
				)) {
			best = candidate
		}
	}

	bestEnd := uint64(best.offset) + uint64(best.size)
	return bytes.Clone(fileData[int(best.offset):int(bestEnd)]), nil
}

func (entry icoImageEntry) effectiveBitDepth(displayBitDepth int) int {
	if entry.bitCount != 0 {
		return int(entry.bitCount)
	}
	if entry.colorCount != 0 {
		depth := 0
		for colors := int(entry.colorCount); colors > 1; colors = (colors + 1) / 2 {
			depth++
		}
		return depth
	}
	return displayBitDepth
}

func betterBitDepth(candidate, best, display int) bool {
	if candidate == display {
		return best != display
	}
	if best == display {
		return false
	}

	candidateFits := candidate <= display
	bestFits := best <= display
	if candidateFits != bestFits {
		return candidateFits
	}
	if candidateFits {
		return candidate > best
	}
	return candidate < best
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func iconResourceData(fileData []byte, desiredWidth, desiredHeight int) ([]byte, error) {
	if len(fileData) < 8 {
		return nil, fmt.Errorf("invalid file format")
	}
	if isPNG(fileData) {
		return bytes.Clone(fileData), nil
	}
	if isICO(fileData) {
		return selectICOImage(fileData, desiredWidth, desiredHeight, currentDisplayBitDepth())
	}
	return nil, fmt.Errorf("unsupported file format")
}

func currentDisplayBitDepth() int {
	hdc := GetDC(0)
	if hdc == 0 {
		return 32
	}
	defer ReleaseDC(0, hdc)

	depth := GetDeviceCaps(hdc, BITSPIXEL) * GetDeviceCaps(hdc, PLANES)
	if depth <= 0 {
		return 32
	}
	return depth
}

// CreateSmallHIconFromImage creates a HICON from a PNG or ICO file
func CreateSmallHIconFromImage(fileData []byte) (HICON, error) {
	iconWidth := GetSystemMetrics(SM_CXSMICON)
	iconHeight := GetSystemMetrics(SM_CYSMICON)
	resourceData, err := iconResourceData(fileData, iconWidth, iconHeight)
	if err != nil {
		return 0, err
	}
	icon, err := CreateIconFromResourceEx(
		uintptr(unsafe.Pointer(&resourceData[0])),
		uint32(len(resourceData)),
		true,
		0x00030000,
		iconWidth,
		iconHeight,
		LR_DEFAULTSIZE)
	runtime.KeepAlive(resourceData)
	return HICON(icon), err
}

// CreateLargeHIconFromImage creates a HICON from a PNG or ICO file
func CreateLargeHIconFromImage(fileData []byte) (HICON, error) {
	iconWidth := GetSystemMetrics(SM_CXICON)
	iconHeight := GetSystemMetrics(SM_CYICON)
	resourceData, err := iconResourceData(fileData, iconWidth, iconHeight)
	if err != nil {
		return 0, err
	}
	icon, err := CreateIconFromResourceEx(
		uintptr(unsafe.Pointer(&resourceData[0])),
		uint32(len(resourceData)),
		true,
		0x00030000,
		iconWidth,
		iconHeight,
		LR_DEFAULTSIZE)
	runtime.KeepAlive(resourceData)
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

// SetMenuIcons decodes the given PNG bytes, installs the bitmaps on the
// identified menu item, and returns the allocated HBITMAP handles. The caller
// owns the returned handles and must DeleteObject them once the menu item no
// longer needs them (for example, in the owning menu's Destroy path). When
// checked is nil a single HBITMAP is allocated and reused for both slots, so
// the returned slice has one entry; otherwise two.
func SetMenuIcons(parentMenu HMENU, itemID int, unchecked []byte, checked []byte) ([]HBITMAP, error) {
	if unchecked == nil {
		return nil, fmt.Errorf("invalid unchecked bitmap")
	}
	uncheckedImage, err := pngToImage(unchecked)
	if err != nil {
		return nil, err
	}
	uncheckedIcon, err := CreateHBITMAPFromImage(uncheckedImage)
	if err != nil {
		return nil, err
	}
	handles := []HBITMAP{uncheckedIcon}
	checkedIcon := uncheckedIcon
	if checked != nil {
		checkedImage, err := pngToImage(checked)
		if err != nil {
			DeleteObject(HGDIOBJ(uncheckedIcon))
			return nil, err
		}
		checkedIcon, err = CreateHBITMAPFromImage(checkedImage)
		if err != nil {
			DeleteObject(HGDIOBJ(uncheckedIcon))
			return nil, err
		}
		handles = append(handles, checkedIcon)
	}
	if err := SetMenuItemBitmaps(parentMenu, uint32(itemID), MF_BYCOMMAND, checkedIcon, uncheckedIcon); err != nil {
		for _, h := range handles {
			DeleteObject(HGDIOBJ(h))
		}
		return nil, err
	}
	return handles, nil
}
