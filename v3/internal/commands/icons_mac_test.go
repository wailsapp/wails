package commands

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestGenerateMacIconPreservesOutputOnEncodeFailure(t *testing.T) {
	for _, tc := range []struct {
		size     int
		existing bool
	}{
		{1, false}, {1, true}, {16, false}, {16, true}, {31, false}, {31, true},
	} {
		t.Run(strconv.Itoa(tc.size)+"/existing="+strconv.FormatBool(tc.existing), func(t *testing.T) {
			var input bytes.Buffer
			if err := png.Encode(&input, image.NewNRGBA(image.Rect(0, 0, tc.size, tc.size))); err != nil {
				t.Fatal(err)
			}
			output := filepath.Join(t.TempDir(), "icon.icns")
			const previous = "previous icon"
			if tc.existing {
				if err := os.WriteFile(output, []byte(previous), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := generateMacIcon(input.Bytes(), &IconsOptions{MacFilename: output}); err == nil {
				t.Fatal("expected too-small image error")
			}
			data, err := os.ReadFile(output)
			if !tc.existing {
				if !os.IsNotExist(err) {
					t.Fatalf("failed encoding created output: data=%q, err=%v", data, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != previous {
				t.Fatalf("existing icon changed to %q", data)
			}
		})
	}
}

func TestGenerateMacIconSuccess(t *testing.T) {
	for _, size := range []int{32, 64} {
		var input bytes.Buffer
		if err := png.Encode(&input, image.NewNRGBA(image.Rect(0, 0, size, size))); err != nil {
			t.Fatal(err)
		}
		output := filepath.Join(t.TempDir(), "icon.icns")
		if err := os.WriteFile(output, []byte("old icon"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := generateMacIcon(input.Bytes(), &IconsOptions{MacFilename: output}); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(output)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) < 16 || string(data[:4]) != "icns" || int(binary.BigEndian.Uint32(data[4:8])) != len(data) {
			t.Fatalf("invalid ICNS family header: %x", data)
		}
		length := int(binary.BigEndian.Uint32(data[12:16]))
		if length < 8 || length > len(data)-8 {
			t.Fatalf("invalid ICNS image length: %d", length)
		}
		decoded, err := png.Decode(bytes.NewReader(data[16 : 8+length]))
		if err != nil {
			t.Fatal(err)
		}
		if decoded.Bounds().Dx() != size || decoded.Bounds().Dy() != size {
			t.Errorf("decoded bounds=%v, want %dx%d", decoded.Bounds(), size, size)
		}
	}
}
