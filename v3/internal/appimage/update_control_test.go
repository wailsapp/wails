package appimage

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestUpdateControlContainsPaddedBlockChecksums(t *testing.T) {
	for _, size := range []int{0, 1, 2047, 2048, 2049, 8193} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			data := make([]byte, size)
			for i := range data {
				data[i] = byte(i*17 + 3)
			}
			image := filepath.Join(t.TempDir(), "image with spaces.AppImage")
			if err := os.WriteFile(image, data, 0755); err != nil {
				t.Fatal(err)
			}
			if err := writeUpdateControl(t.Context(), image); err != nil {
				t.Fatal(err)
			}
			if err := verifyZsync(t.Context(), image); err != nil {
				t.Fatal(err)
			}
			control, err := os.ReadFile(image + ".zsync")
			if err != nil {
				t.Fatal(err)
			}
			headers, records, ok := bytes.Cut(control, []byte("\n\n"))
			if !ok {
				t.Fatal("missing record boundary")
			}
			if !strings.Contains(string(headers), "URL: image%20with%20spaces.AppImage") {
				t.Fatalf("bad URL: %s", headers)
			}
			count := (size + 2047) / 2048
			if len(records) != count*20 {
				t.Fatalf("record bytes=%d want=%d", len(records), count*20)
			}
			// Captured from official zsyncmake 0.6.2 with -b 2048. Its
			// shortened records keep the low two rolling-checksum bytes.
			golden := map[int]string{0: "", 1: "180086c983", 2047: "2f0e59fb2eb0", 2048: "30002e52b56b", 2049: "30002e52b5180086c983", 8193: "30002e52b530002e52b530002e52b530002e52b5180086c983"}
			reference, err := hex.DecodeString(golden[size])
			if err != nil {
				t.Fatal(err)
			}
			if count > 0 {
				width := len(reference) / count
				for index := range count {
					record := records[index*20 : (index+1)*20]
					if !bytes.Equal(record[2:2+width], reference[index*width:(index+1)*width]) {
						t.Fatalf("record %d disagrees with official zsyncmake: %x", index, record)
					}
				}
			}
		})
	}
}

func BenchmarkUpdateControl(b *testing.B) {
	root := b.TempDir()
	image := filepath.Join(root, "image.AppImage")
	data := bytes.Repeat([]byte("representative binary contents\x00"), 1<<15)
	if err := os.WriteFile(image, data, 0644); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := writeUpdateControl(b.Context(), image); err != nil {
			b.Fatal(err)
		}
		if err := os.Remove(image + ".zsync"); err != nil {
			b.Fatal(err)
		}
	}
}
