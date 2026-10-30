package pcinames

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

const fixture = `# Synthetic PCI names, independently written for these tests.
1234  Example  vendor
	5678  Example device
		1234 0001  Subsystem must not replace device
	9abc  Second device
ABCD  Unicode vendor Ω
	EF01  Unicode device β
C 03  Display controller
	00  VGA controller
X future extension
	1234  Must not become a device
`

func TestReadNames(t *testing.T) {
	for _, data := range [][]byte{[]byte(fixture), compressed(t, []byte(fixture))} {
		db, err := read(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		for _, tt := range []struct{ vendor, device, wantVendor, wantDevice string }{
			{"1234", "5678", "Example  vendor", "Example device"},
			{"1234", "9abc", "Example  vendor", "Second device"},
			{"abcd", "ef01", "Unicode vendor Ω", "Unicode device β"},
			{"ABCD", "EF01", "Unicode vendor Ω", "Unicode device β"},
			{"1234", "0000", "Example  vendor", ""},
			{"1234", "invalid", "Example  vendor", ""},
			{"ffff", "5678", "", ""}, {"invalid", "5678", "", ""},
		} {
			vendor, device := db.Lookup(tt.vendor, tt.device)
			if vendor != tt.wantVendor || device != tt.wantDevice {
				t.Fatalf("%s:%s got %q %q; want %q %q", tt.vendor, tt.device, vendor, device, tt.wantVendor, tt.wantDevice)
			}
		}
	}
}

func TestMalformedRecords(t *testing.T) {
	data := "\t1234  Orphan\n1234  Good\n\t5678  Card\n\t\nzzzz  Invalid vendor\n\t5678  Wrong association\n9999\n\t0000  Wrong association\n5678  \n\t0000  Empty parent\n"
	db, err := read(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(db.vendors) != 1 {
		t.Fatalf("invalid vendors accepted: %d", len(db.vendors))
	}
	_, card := db.Lookup("1234", "5678")
	if card != "Card" {
		t.Fatal(card)
	}
}

type failedReader struct{}

func (failedReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestReadErrors(t *testing.T) {
	if _, err := read(failedReader{}); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("reader error lost: %v", err)
	}
	gzipData := compressed(t, []byte(fixture))
	gzipData[len(gzipData)-1] ^= 0xff
	for _, data := range [][]byte{{0x1f, 0x8b}, gzipData, []byte(strings.Repeat("x", 1024*1024+1))} {
		if _, err := read(bytes.NewReader(data)); err == nil {
			t.Fatal("corrupt or oversized input accepted")
		}
	}
	// Bound decompressed data as well as downloaded bytes.
	if _, err := read(strings.NewReader(strings.Repeat("# ignored\n", maxDownloadBytes/10+1))); err == nil {
		t.Fatal("database size limit ignored")
	}
}

func compressed(t testing.TB, data []byte) []byte {
	t.Helper()
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func BenchmarkParseNames(b *testing.B) {
	var text strings.Builder
	for vendor := 0; vendor < 200; vendor++ {
		fmt.Fprintf(&text, "%04x  Example vendor %d\n", vendor, vendor)
		for device := 0; device < 16; device++ {
			fmt.Fprintf(&text, "\t%04x  Example device %d\n", device, device)
		}
	}
	data := []byte(text.String())
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	for b.Loop() {
		if _, err := read(bytes.NewReader(data)); err != nil {
			b.Fatal(err)
		}
	}
}
func FuzzReadNames(f *testing.F) {
	f.Add([]byte(fixture))
	f.Add([]byte{0x1f, 0x8b})
	f.Add([]byte("C 03  Display\n\t00  VGA\n"))
	f.Fuzz(func(t *testing.T, data []byte) {
		db, err := read(bytes.NewReader(data))
		if err != nil {
			return
		}
		for id, v := range db.vendors {
			if v == nil || v.name == "" {
				t.Fatalf("invalid vendor %x", id)
			}
		}
	})
}
