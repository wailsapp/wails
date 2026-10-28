package commands

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackmordaunt/icns/v2"
)

func TestGenerateMacIconPreservesOutputOnEncodeFailure(t *testing.T) {
	var input bytes.Buffer
	if err := png.Encode(&input, image.NewNRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatal(err)
	}
	for _, existing := range []bool{false, true} {
		name := "new output"
		if existing {
			name = "existing output"
		}
		t.Run(name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "icon.icns")
			const previous = "previous icon"
			if existing {
				if err := os.WriteFile(output, []byte(previous), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := generateMacIcon(input.Bytes(), &IconsOptions{MacFilename: output}); err == nil {
				t.Fatal("expected too-small image error")
			}
			data, err := os.ReadFile(output)
			if !existing {
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
	var input bytes.Buffer
	if err := png.Encode(&input, image.NewNRGBA(image.Rect(0, 0, 64, 64))); err != nil {
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
	decoded, err := icns.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds().Dx() != 64 || decoded.Bounds().Dy() != 64 {
		t.Errorf("decoded bounds=%v, want 64x64", decoded.Bounds())
	}
}
