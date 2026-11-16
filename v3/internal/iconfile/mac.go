package iconfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

// EncodeMac writes PNG representations in an ICNS icon family.
// Artwork is squared using Mitchell-Netravali scaling; representations larger
// than its longest side are omitted. At least 32 pixels are required.
func EncodeMac(w io.Writer, src image.Image) error {
	limit := max(src.Bounds().Dx(), src.Bounds().Dy())
	if limit < 32 {
		return fmt.Errorf("input image is too small to generate a macOS icon: at least 32 pixels required")
	}
	encoder := png.Encoder{}
	// Artwork at least 64px emits two or more PNG sizes.
	if limit >= 64 {
		encoder.BufferPool = &pngBuffers{}
	}
	var output bytes.Buffer
	output.WriteString("icns")
	output.Write(make([]byte, 4))
	// Retina and standard representations may share a pixel size and PNG payload.
	for _, representation := range []struct {
		size  int
		types []string
	}{
		{1024, []string{"ic10"}},
		{512, []string{"ic14", "ic09"}},
		{256, []string{"ic13", "ic08"}},
		{128, []string{"ic07"}},
		{64, []string{"ic12"}},
		{32, []string{"ic11"}},
	} {
		if representation.size > limit {
			continue
		}
		resized := resize.Resize(uint(representation.size), uint(representation.size), src, resize.MitchellNetravali)
		var payload bytes.Buffer
		if err := encoder.Encode(&payload, resized); err != nil {
			return err
		}
		for _, kind := range representation.types {
			output.WriteString(kind)
			var length [4]byte
			binary.BigEndian.PutUint32(length[:], uint32(payload.Len()+8))
			output.Write(length[:])
			output.Write(payload.Bytes())
		}
	}
	binary.BigEndian.PutUint32(output.Bytes()[4:8], uint32(output.Len()))
	n, err := w.Write(output.Bytes())
	if err == nil && n != output.Len() {
		return io.ErrShortWrite
	}
	return err
}
