package iconfile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"strconv"
	"testing"
)

func macFrames(t *testing.T, data []byte) map[string]image.Image {
	t.Helper()
	if len(data) < 8 || string(data[:4]) != "icns" || int(binary.BigEndian.Uint32(data[4:8])) != len(data) {
		t.Fatal("invalid ICNS family header")
	}
	frames := map[string]image.Image{}
	for offset := 8; offset < len(data); {
		if len(data)-offset < 8 {
			t.Fatal("truncated ICNS chunk header")
		}
		kind := string(data[offset : offset+4])
		length := int(binary.BigEndian.Uint32(data[offset+4 : offset+8]))
		if length <= 8 || length > len(data)-offset {
			t.Fatalf("invalid %s chunk length %d", kind, length)
		}
		if frames[kind] != nil {
			t.Fatalf("duplicate %s chunk", kind)
		}
		img, err := png.Decode(bytes.NewReader(data[offset+8 : offset+length]))
		if err != nil {
			t.Fatal(err)
		}
		frames[kind] = img
		offset += length
	}
	return frames
}

func TestEncodeMacRepresentations(t *testing.T) {
	for _, size := range []int{32, 64, 100, 128, 256, 512, 1024} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			src := image.NewNRGBA(image.Rect(3, 5, 3+size, 5+size))
			for y := 5; y < 5+size; y++ {
				for x := 3; x < 3+size; x++ {
					src.SetNRGBA(x, y, color.NRGBA{R: 180, G: 70, B: 30, A: 128})
				}
			}
			var output bytes.Buffer
			if err := EncodeMac(&output, src); err != nil {
				t.Fatal(err)
			}
			frames := macFrames(t, output.Bytes())
			wantCount := 0
			for kind, pixels := range map[string]int{"ic11": 32, "ic12": 64, "ic07": 128, "ic08": 256, "ic13": 256, "ic09": 512, "ic14": 512, "ic10": 1024} {
				if pixels > size {
					continue
				}
				wantCount++
				frame := frames[kind]
				if frame == nil {
					t.Fatalf("missing %s", kind)
				}
				if frame.Bounds() != image.Rect(0, 0, pixels, pixels) {
					t.Fatalf("%s has bounds %v", kind, frame.Bounds())
				}
				got := color.NRGBAModel.Convert(frame.At(pixels/2, pixels/2)).(color.NRGBA)
				// The existing Mitchell scaler may round a color channel down by one.
				if got.A != 128 || got.R < 179 || got.R > 180 || got.G < 69 || got.G > 70 || got.B < 29 || got.B > 30 {
					t.Fatalf("%s pixel=%v", kind, got)
				}
			}
			if len(frames) != wantCount {
				t.Fatalf("got %d frames, want %d", len(frames), wantCount)
			}
		})
	}
}

func TestEncodeMacSmallArtwork(t *testing.T) {
	for _, size := range []int{0, 1, 16, 31} {
		var output bytes.Buffer
		if err := EncodeMac(&output, image.NewNRGBA(image.Rect(0, 0, size, size))); err == nil {
			t.Fatalf("accepted %dpx artwork", size)
		}
		if output.Len() != 0 {
			t.Fatal("invalid artwork wrote output")
		}
	}
}

func TestEncodeMacWriteErrors(t *testing.T) {
	failure := errors.New("write failed")
	src := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	if err := EncodeMac(failingWriter{err: failure}, src); !errors.Is(err, failure) {
		t.Fatalf("write error=%v", err)
	}
	if err := EncodeMac(failingWriter{}, src); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short-write error=%v", err)
	}
}

func BenchmarkEncodeMac(b *testing.B) {
	for _, size := range []int{64, 512} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			src := image.NewNRGBA(image.Rect(0, 0, size, size))
			for y := 0; y < size; y++ {
				for x := 0; x < size; x++ {
					src.SetNRGBA(x, y, color.NRGBA{R: uint8(x), G: uint8(y), B: uint8(x * y), A: uint8(x + y)})
				}
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := EncodeMac(io.Discard, src); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
