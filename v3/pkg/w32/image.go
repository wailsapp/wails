package w32

import (
	"image"
	"syscall"
	"unsafe"
)

func CreateHBITMAPFromImage(img *image.RGBA) (HBITMAP, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create a BITMAPINFO structure for the DIB
	bmi := BITMAPINFO{
		BmiHeader: BITMAPINFOHEADER{
			BiSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
			BiWidth:       int32(width),
			BiHeight:      int32(-height), // negative to indicate top-down bitmap
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: BI_RGB,
			BiSizeImage:   uint32(width * height * 4), // RGBA = 4 bytes
		},
	}

	// Create the DIB section
	var bits unsafe.Pointer

	hbmp := CreateDIBSection(0, &bmi, DIB_RGB_COLORS, &bits, 0, 0)
	if hbmp == 0 {
		return 0, syscall.GetLastError()
	}

	// Copy the pixel data from the Go image to the DIB section
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]
			a := img.Pix[i+3]

			// Write the RGBA pixel data to the DIB section (BGR order)
			offset := y*width*4 + x*4
			*((*uint8)(unsafe.Pointer(uintptr(bits) + uintptr(offset) + 0))) = b
			*((*uint8)(unsafe.Pointer(uintptr(bits) + uintptr(offset) + 1))) = g
			*((*uint8)(unsafe.Pointer(uintptr(bits) + uintptr(offset) + 2))) = r
			*((*uint8)(unsafe.Pointer(uintptr(bits) + uintptr(offset) + 3))) = a
		}
	}

	return hbmp, nil
}
