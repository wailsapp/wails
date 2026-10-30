package iconfile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"slices"
	"testing"
)

func TestEncodeBMPArtwork(t *testing.T) {
	artwork, err := os.ReadFile("testdata/artwork.bmp")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := Encode(&output, bytes.NewReader(artwork), []int{16}); err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(output.Bytes()[22:]))
	if err != nil {
		t.Fatal(err)
	}
	got := color.NRGBAModel.Convert(img.At(8, 8)).(color.NRGBA)
	if got != (color.NRGBA{R: 180, G: 70, B: 30, A: 255}) {
		t.Fatalf("BMP pixel = %v", got)
	}
}

func TestEncodeImages(t *testing.T) {
	src := image.NewNRGBA(image.Rect(5, 7, 8, 9))
	for y := 7; y < 9; y++ {
		for x := 5; x < 8; x++ {
			src.SetNRGBA(x, y, color.NRGBA{R: 180, G: 70, B: 30, A: 128})
		}
	}
	var input bytes.Buffer
	if err := png.Encode(&input, src); err != nil {
		t.Fatal(err)
	}
	sizes := []int{1, 256, 16, 48, 16}
	var output bytes.Buffer
	if err := Encode(&output, bytes.NewReader(input.Bytes()), sizes); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(sizes, []int{1, 256, 16, 48, 16}) {
		t.Fatalf("Encode changed caller sizes: %v", sizes)
	}
	data := output.Bytes()
	if binary.LittleEndian.Uint16(data[:2]) != 0 || binary.LittleEndian.Uint16(data[2:4]) != 1 || int(binary.LittleEndian.Uint16(data[4:6])) != len(sizes) {
		t.Fatalf("invalid ICO header: %x", data[:6])
	}
	offset := 6 + 16*len(sizes)
	for i, size := range []int{256, 48, 16, 16, 1} {
		entry := data[6+16*i : 6+16*(i+1)]
		if entry[2] != 0 || entry[3] != 0 || binary.LittleEndian.Uint16(entry[4:6]) != 1 {
			t.Fatalf("invalid color count, reserved byte or planes: %x", entry)
		}
		if int(entry[0]) != size%256 || int(entry[1]) != size%256 || binary.LittleEndian.Uint16(entry[6:8]) != 32 {
			t.Fatalf("invalid directory entry for size %d: %x", size, entry)
		}
		if int(binary.LittleEndian.Uint32(entry[12:16])) != offset {
			t.Fatalf("image %d has incorrect offset", i)
		}
		length := int(binary.LittleEndian.Uint32(entry[8:12]))
		img, err := png.Decode(bytes.NewReader(data[offset : offset+length]))
		if err != nil {
			t.Fatal(err)
		}
		if img.Bounds() != image.Rect(0, 0, size, size) {
			t.Fatalf("image bounds %v, want size %d", img.Bounds(), size)
		}
		got := color.NRGBAModel.Convert(img.At(size/2, size/2)).(color.NRGBA)
		if got.A != 128 || got.R < 178 || got.R > 181 || got.G < 68 || got.G > 71 || got.B < 28 || got.B > 31 {
			t.Fatalf("resized pixel lost color or transparency: %v", got)
		}
		offset += length
	}
	if offset != len(data) {
		t.Fatal("unexpected trailing bytes")
	}
}

func TestEncodeInvalidSizes(t *testing.T) {
	for _, sizes := range [][]int{nil, {}, {0}, {-1}, {257}, {16, 0}, make([]int, 65536)} {
		var output bytes.Buffer
		if err := Encode(&output, bytes.NewReader(nil), sizes); err == nil {
			t.Fatalf("accepted invalid sizes (count %d)", len(sizes))
		}
		if output.Len() != 0 {
			t.Fatal("invalid sizes wrote output")
		}
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write(p []byte) (int, error) { return 0, w.err }

func TestEncodeWriteErrors(t *testing.T) {
	failure := errors.New("write failed")
	for _, want := range []error{failure, io.ErrShortWrite} {
		writer := failingWriter{err: failure}
		if want == io.ErrShortWrite {
			writer.err = nil
		}
		var input bytes.Buffer
		if err := png.Encode(&input, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
			t.Fatal(err)
		}
		err := Encode(writer, &input, []int{16})
		if !errors.Is(err, want) {
			t.Fatalf("error = %v, want %v", err, want)
		}
	}
}

func TestEncodeInvalidImage(t *testing.T) {
	var output bytes.Buffer
	if err := Encode(&output, bytes.NewBufferString("not artwork"), []int{16}); err == nil {
		t.Fatal("accepted invalid artwork")
	}
	if output.Len() != 0 {
		t.Fatal("invalid artwork wrote output")
	}
}

func BenchmarkEncode(b *testing.B) {
	src := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			src.SetNRGBA(x, y, color.NRGBA{R: uint8(x), G: uint8(y), B: uint8(x * y), A: uint8(x + y)})
		}
	}
	var input bytes.Buffer
	if err := png.Encode(&input, src); err != nil {
		b.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		sizes []int
	}{{"small", []int{16}}, {"default", []int{256, 128, 64, 48, 32, 16}}} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(input.Len()))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := Encode(io.Discard, bytes.NewReader(input.Bytes()), tc.sizes); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
