// Package iconfile generates Windows and macOS icon files from application artwork.
package iconfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	stddraw "image/draw"
	"image/png"
	"io"
	"math"
	"slices"

	_ "golang.org/x/image/bmp" // Register BMP artwork for image.Decode.
	"golang.org/x/image/draw"
)

// Encode reads application artwork and writes a PNG-based ICO file.
// Sizes must be between 1 and 256 pixels. The directory is ordered largest first.
func Encode(w io.Writer, r io.Reader, sizes []int) error {
	if len(sizes) == 0 || len(sizes) > math.MaxUint16 {
		return fmt.Errorf("icon requires between 1 and %d images", math.MaxUint16)
	}
	for _, size := range sizes {
		if size < 1 || size > 256 {
			return fmt.Errorf("icon size must be between 1 and 256 pixels: %d", size)
		}
	}
	src, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	// ICO stores a six-byte header, sixteen-byte directory entries, then images.
	var output bytes.Buffer
	output.Write(make([]byte, 6+16*len(sizes)))
	binary.LittleEndian.PutUint16(output.Bytes()[2:4], 1)
	binary.LittleEndian.PutUint16(output.Bytes()[4:6], uint16(len(sizes)))
	ordered := slices.Clone(sizes)
	slices.SortFunc(ordered, func(a, b int) int { return b - a })
	encoder := png.Encoder{}
	if len(sizes) > 1 {
		encoder.BufferPool = &pngBuffers{}
	}
	for i, size := range ordered {
		dst := image.NewRGBA(image.Rect(0, 0, size, size))
		draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
		straight := image.NewNRGBA(dst.Bounds())
		stddraw.Draw(straight, straight.Bounds(), dst, dst.Bounds().Min, stddraw.Src)
		offset := output.Len()
		if err := encoder.Encode(&output, straight); err != nil {
			return err
		}
		if uint64(output.Len()) > math.MaxUint32 {
			return fmt.Errorf("icon file exceeds the 32-bit ICO size limit")
		}
		entry := output.Bytes()[6+16*i : 6+16*(i+1)]
		entry[0], entry[1] = byte(size%256), byte(size%256)
		binary.LittleEndian.PutUint16(entry[4:6], 1)
		binary.LittleEndian.PutUint16(entry[6:8], 32)
		binary.LittleEndian.PutUint32(entry[8:12], uint32(output.Len()-offset))
		binary.LittleEndian.PutUint32(entry[12:16], uint32(offset))
	}
	n, err := w.Write(output.Bytes())
	if err == nil && n != output.Len() {
		return io.ErrShortWrite
	}
	return err
}
