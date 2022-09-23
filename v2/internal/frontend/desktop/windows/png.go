package windows

import (
	"bufio"
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/png"
)

func ResizePNG(in []byte, size int) ([]byte, error) {
	imagedata, _, err := image.Decode(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	// Scale image
	rect := image.Rect(0, 0, size, size)
	rawdata := image.NewRGBA(rect)
	scale := draw.CatmullRom
	scale.Scale(rawdata, rect, imagedata, imagedata.Bounds(), draw.Over, nil)

	// Convert back to PNG
	icondata := new(bytes.Buffer)
	writer := bufio.NewWriter(icondata)
	err = png.Encode(writer, rawdata)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	// Save image data
	return icondata.Bytes(), nil
}
