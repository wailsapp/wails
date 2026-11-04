package notaryticket

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/applesign"
	"github.com/wailsapp/wails/v3/internal/universalbinary"
)

func syntheticDMG() []byte {
	be := binary.BigEndian
	directory := make([]byte, 60)
	be.PutUint32(directory, 0xfade0c02)
	be.PutUint32(directory[4:], 60)
	be.PutUint32(directory[20:], 44)
	directory[36], directory[37] = 32, 2
	copy(directory[44:], "owned-fixture\x00")
	cms := []byte{0xfa, 0xde, 0x0b, 0x01, 0, 0, 0, 8}
	sig := make([]byte, 28+len(directory)+len(cms))
	be.PutUint32(sig, 0xfade0cc0)
	be.PutUint32(sig[4:], uint32(len(sig)))
	be.PutUint32(sig[8:], 2)
	be.PutUint32(sig[16:], 28)
	be.PutUint32(sig[20:], 0x10000)
	be.PutUint32(sig[24:], uint32(28+len(directory)))
	copy(sig[28:], directory)
	copy(sig[28+len(directory):], cms)
	image := append(make([]byte, 512), sig...)
	footer := make([]byte, 512)
	copy(footer, "koly")
	be.PutUint32(footer[4:], 4)
	be.PutUint32(footer[8:], 512)
	be.PutUint64(footer[296:], 512)
	be.PutUint64(footer[304:], uint64(len(sig)))
	return append(image, footer...)
}

func TestDMGWriterPreservesSignedComponents(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "input.dmg")
	output := filepath.Join(root, "output.dmg")
	original := syntheticDMG()
	if err := os.WriteFile(input, original, 0644); err != nil {
		t.Fatal(err)
	}
	hashes, err := Digests(input)
	if err != nil {
		t.Fatal(err)
	}
	// The container writer uses an owned synthetic CodeDirectory. Authenticity
	// tests separately use real Apple-signed envelopes, never a local trust bypass.
	data, _ := fixture(t)
	ticket := &Ticket{data: data, hashes: map[Digest]bool{hashes[0]: true}}
	if err := Staple(t.Context(), input, output, ticket); err != nil {
		t.Fatal(err)
	}
	before := sha256.Sum256(original)
	after, _ := os.ReadFile(input)
	if before != sha256.Sum256(after) {
		t.Fatal("source changed")
	}
	result, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	parts, err := signatureParts(result[512 : len(result)-512])
	if err != nil {
		t.Fatal(err)
	}
	oldParts, err := signatureParts(original[512 : len(original)-512])
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != len(oldParts) {
		t.Fatal("signed component count changed")
	}
	for i := range parts {
		if parts[i].slot != oldParts[i].slot || !bytes.Equal(parts[i].data, oldParts[i].data) {
			t.Fatal("signed component altered")
		}
	}
	got, err := Digests(output)
	if err != nil || len(got) != 1 || got[0] != hashes[0] {
		t.Fatalf("code identity changed: %v %v", got, err)
	}
	if err := Staple(t.Context(), input, output, ticket); err == nil {
		t.Fatal("overwrote existing output")
	}
	unchanged, _ := os.ReadFile(output)
	if !bytes.Equal(unchanged, result) {
		t.Fatal("existing output changed")
	}
	wrong, err := Verify(data, []Digest{{Algorithm: 2, Value: [20]byte{0x04, 0x3b, 0xaa, 0x30, 0x56, 0xb5, 0x6e, 0x8b, 0x21, 0xcf, 0xaa, 0x9d, 0x52, 0x2d, 0x56, 0x37, 0x27, 0x75, 0x4d, 0x1d}}})
	if err != nil {
		t.Fatal(err)
	}
	rejected := filepath.Join(root, "wrong.dmg")
	if err := Staple(t.Context(), input, rejected, wrong); err == nil {
		t.Fatal("stapled unrelated ticket")
	}
	if _, err := os.Lstat(rejected); !os.IsNotExist(err) {
		t.Fatal("failed output leaked")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := Staple(ctx, input, rejected, ticket); err == nil {
		t.Fatal("ignored cancellation")
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 2 {
		t.Fatal("temporary files leaked")
	}
}

func TestDMGRejectsMalformedComponents(t *testing.T) {
	image := syntheticDMG()
	data := image[512 : len(image)-512]
	for _, mutate := range []func([]byte){
		func(b []byte) { binary.BigEndian.PutUint32(b[8:], 65) },
		func(b []byte) { binary.BigEndian.PutUint32(b[24:], 28) },
		func(b []byte) { binary.BigEndian.PutUint32(b[20:], 0) },
		func(b []byte) { binary.BigEndian.PutUint32(b[16:], 0xffffffff) },
		func(b []byte) { binary.BigEndian.PutUint32(b[32:], 0xffffffff) },
	} {
		b := bytes.Clone(data)
		mutate(b)
		if _, err := signatureParts(b); err == nil {
			t.Fatal("accepted malformed container")
		}
	}
}

func TestUniversalBundleStapling(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "Input.app")
	main := filepath.Join(input, "Contents", "MacOS", "probe")
	if err := os.MkdirAll(filepath.Dir(main), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(input, "Contents", "Info.plist"), []byte(`<?xml version="1.0"?><plist version="1.0"><dict><key>CFBundleIdentifier</key><string>org.wails.ticketprobe</string><key>CFBundleExecutable</key><string>probe</string></dict></plist>`), 0644); err != nil {
		t.Fatal(err)
	}
	raw := filepath.Join(root, "universal")
	if err := universalbinary.Create(raw, []string{"../applesign/testdata/amd64.macho", "../applesign/testdata/arm64.macho"}); err != nil {
		t.Fatal(err)
	}
	if _, err := applesign.SignMachO(t.Context(), raw, main, applesign.Options{Identifier: "org.wails.ticketprobe"}); err != nil {
		t.Fatal(err)
	}
	hashes, err := Digests(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(hashes) != 2 {
		t.Fatal("fixture is not universal")
	}
	data, _ := fixture(t)
	ticket := &Ticket{data: data, hashes: map[Digest]bool{hashes[0]: true, hashes[1]: true}}
	output := filepath.Join(root, "Output.app")
	if err := Staple(t.Context(), input, output, ticket); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(input, "Contents", "CodeResources")); !os.IsNotExist(err) {
		t.Fatal("modified source bundle")
	}
	got, err := os.ReadFile(filepath.Join(output, "Contents", "CodeResources"))
	if err != nil || !bytes.Equal(got, data) {
		t.Fatalf("missing output ticket: %v", err)
	}
	if err := Staple(t.Context(), input, output, ticket); err == nil {
		t.Fatal("replaced existing bundle")
	}
	got, err = os.ReadFile(filepath.Join(output, "Contents", "CodeResources"))
	if err != nil || !bytes.Equal(got, data) {
		t.Fatal("existing output changed")
	}
	partial := &Ticket{data: data, hashes: map[Digest]bool{hashes[0]: true}}
	failed := filepath.Join(root, "Failed.app")
	if err := Staple(t.Context(), input, failed, partial); err == nil {
		t.Fatal("accepted ticket missing one architecture")
	}
	if _, err := os.Lstat(failed); !os.IsNotExist(err) {
		t.Fatal("failed bundle leaked")
	}
	restapled := filepath.Join(root, "Restapled.app")
	if err := Staple(t.Context(), output, restapled, ticket); err != nil {
		t.Fatal(err)
	}
	t.Run("ticket-symlink", func(t *testing.T) {
		seal := filepath.Join(input, "Contents", "owned-seal")
		if err := os.WriteFile(seal, []byte("unchanged"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("owned-seal", filepath.Join(input, "Contents", "CodeResources")); err != nil {
			t.Skipf("host cannot create symlinks: %v", err)
		}
		if err := Staple(t.Context(), input, failed, ticket); err == nil {
			t.Fatal("followed existing ticket symlink")
		}
		if _, err := os.Lstat(failed); !os.IsNotExist(err) {
			t.Fatal("failed bundle leaked")
		}
		got, err := os.ReadFile(seal)
		if err != nil || string(got) != "unchanged" {
			t.Fatal("modified source symlink target")
		}
	})
}
