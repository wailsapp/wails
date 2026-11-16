package packagesign

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testPassword = "package-test-password"
const testFingerprint = "41FA35E5B4D52D500E9C514B06A24AC943208C69"

func fixtureOptions(format string) Options {
	return Options{Format: format, KeyFile: "testdata/private.asc", Password: testPassword}
}

func TestSignPreservesPackageAndResigns(t *testing.T) {
	for _, format := range []string{"deb", "rpm"} {
		t.Run(format, func(t *testing.T) {
			original := readFile(t, "testdata/unsigned."+format)
			signed := filepath.Join(t.TempDir(), "signed."+format)
			if err := Sign(t.Context(), "testdata/unsigned."+format, signed, fixtureOptions(format)); err != nil {
				t.Fatal(err)
			}
			first := readFile(t, signed)
			assertPayloadPreserved(t, format, original, first)
			resigned := filepath.Join(t.TempDir(), "resigned."+format)
			if err := Sign(t.Context(), signed, resigned, fixtureOptions(format)); err != nil {
				t.Fatal(err)
			}
			second := readFile(t, resigned)
			assertPayloadPreserved(t, format, first, second)
			if format == "deb" && !bytes.Contains(second, []byte("_gpgbuilder1")) {
				t.Fatal("second signature did not retain the first")
			}
		})
	}
}

func assertPayloadPreserved(t *testing.T, format string, before, after []byte) {
	t.Helper()
	if format == "deb" {
		if !bytes.HasPrefix(after, before) {
			t.Fatal("DEB existing members changed")
		}
		return
	}
	offset := func(data []byte) int {
		size := 16 + int(binary.BigEndian.Uint32(data[104:108]))*16 + int(binary.BigEndian.Uint32(data[108:112]))
		return 96 + (size+7)/8*8
	}
	if !bytes.Equal(before[offset(before):], after[offset(after):]) {
		t.Fatal("RPM immutable header or payload changed")
	}
}

func TestSigningKeySelection(t *testing.T) {
	for _, format := range []string{"deb", "rpm"} {
		options := fixtureOptions(format)
		options.KeyFile, options.Password = "testdata/unencrypted.gpg", ""
		if err := Sign(t.Context(), "testdata/unsigned."+format, filepath.Join(t.TempDir(), "signed"), options); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{testFingerprint, "06A24AC943208C69", "0x06a24ac943208c69"} {
		options := fixtureOptions("rpm")
		options.KeyFile, options.KeyID = "testdata/keys.gpg", id
		if err := Sign(t.Context(), "testdata/unsigned.rpm", filepath.Join(t.TempDir(), "signed.rpm"), options); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct{ name, file, password, selector string }{
		{"wrong password", "testdata/private.asc", "wrong", ""},
		{"public only", "testdata/public.asc", "", ""},
		{"unknown key", "testdata/private.asc", testPassword, "0000000000000000"},
		{"ambiguous export", "testdata/keys.gpg", testPassword, ""},
		{"short id", "testdata/private.asc", testPassword, "43208C69"},
		{"non hex", "testdata/private.asc", testPassword, "bad"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "signed.deb")
			err := Sign(t.Context(), "testdata/unsigned.deb", out, Options{Format: "deb", KeyFile: tc.file, Password: tc.password, KeyID: tc.selector})
			if err == nil {
				t.Fatal("expected rejection")
			}
			if _, err := os.Stat(out); !os.IsNotExist(err) {
				t.Fatal("invalid key wrote output")
			}
		})
	}
}

func TestSigningRejectsInvalidInputs(t *testing.T) {
	for _, format := range []string{"deb", "rpm"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			input := filepath.Join(root, "input")
			out := filepath.Join(root, "out")
			original := readFile(t, "testdata/unsigned."+format)
			writeFile(t, input, original)
			if err := Sign(t.Context(), input, input, fixtureOptions(format)); err == nil {
				t.Fatal("same file accepted")
			}
			if !bytes.Equal(readFile(t, input), original) {
				t.Fatal("input changed")
			}
			alias := filepath.Join(root, "alias")
			if err := os.Link(input, alias); err == nil {
				if err := Sign(t.Context(), input, alias, fixtureOptions(format)); err == nil {
					t.Fatal("hardlink accepted")
				}
			}
			writeFile(t, input, original[:len(original)/2])
			if err := Sign(t.Context(), input, out, fixtureOptions(format)); err == nil {
				t.Fatal("truncation accepted")
			}
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			if err := Sign(ctx, input, out, fixtureOptions(format)); !errors.Is(err, context.Canceled) {
				t.Fatalf("canceled: %v", err)
			}
		})
	}
}

func TestDetachedArchiveSignature(t *testing.T) {
	input := "testdata/unsigned.deb"
	before := readFile(t, input)
	output := filepath.Join(t.TempDir(), "archive.sig")
	if err := Sign(t.Context(), input, output, fixtureOptions("archlinux")); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, readFile(t, input)) || len(readFile(t, output)) == 0 {
		t.Fatal("detached signing changed its input or omitted the signature")
	}
}

func TestRPMRejectsTampering(t *testing.T) {
	data := readFile(t, "testdata/unsigned.rpm")
	for _, offset := range []int{len(data) - 1, 104, 108} {
		changed := bytes.Clone(data)
		changed[offset] ^= 0x7f
		input := filepath.Join(t.TempDir(), "changed.rpm")
		writeFile(t, input, changed)
		if err := Sign(t.Context(), input, filepath.Join(t.TempDir(), "signed.rpm"), fixtureOptions("rpm")); err == nil {
			t.Fatalf("tamper at %d accepted", offset)
		}
	}
}

func TestDEBRoleAndMalformedMembers(t *testing.T) {
	for _, role := range []string{"bad\nrole", "../escape", strings.Repeat("a", 11)} {
		options := fixtureOptions("deb")
		options.Role = role
		if err := Sign(t.Context(), "testdata/unsigned.deb", filepath.Join(t.TempDir(), "out"), options); err == nil {
			t.Fatalf("role %q accepted", role)
		}
	}
	data := readFile(t, "testdata/unsigned.deb")
	for _, offset := range []int{0, 8, 56, 66} {
		changed := bytes.Clone(data)
		changed[offset] = 0xff
		if _, err := debMembers(t.Context(), bytes.NewReader(changed)); err == nil {
			t.Fatalf("malformed field at %d accepted", offset)
		}
	}
}

func TestRPMHeaderValueBounds(t *testing.T) {
	for _, tc := range []struct{ kind, offset, count uint32 }{
		{7, 99, 1}, {4, 1, 1}, {5, 0, ^uint32(0)}, {6, 0, 2}, {8, 0, 99}, {9, 0, 1}, {99, 0, 1},
	} {
		if _, err := rpmValue([]byte{1, 2, 3, 4}, tc.kind, tc.offset, tc.count); err == nil {
			t.Fatalf("invalid value accepted: %+v", tc)
		}
	}
	if _, _, err := readRPMHeader(io.LimitReader(bytes.NewReader(readFile(t, "testdata/unsigned.rpm")[96:]), 20)); err == nil {
		t.Fatal("truncated header accepted")
	}
}

func FuzzDEBEnvelope(f *testing.F) {
	data, err := os.ReadFile("testdata/unsigned.deb")
	if err != nil {
		f.Fatal(err)
	}
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) { _, _ = debMembers(t.Context(), bytes.NewReader(data)) })
}
func FuzzRPMHeader(f *testing.F) {
	data, err := os.ReadFile("testdata/unsigned.rpm")
	if err != nil {
		f.Fatal(err)
	}
	f.Add(data[96:])
	f.Fuzz(func(t *testing.T, data []byte) { _, _, _ = readRPMHeader(bytes.NewReader(data)) })
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSignPackage(b *testing.B) {
	for _, format := range []string{"deb", "rpm"} {
		for _, key := range []string{"private.asc", "unencrypted.gpg"} {
			b.Run(format+"/"+key, func(b *testing.B) {
				input := "testdata/unsigned." + format
				out := filepath.Join(b.TempDir(), "signed."+format)
				options := fixtureOptions(format)
				options.KeyFile = "testdata/" + key
				b.ReportAllocs()
				for b.Loop() {
					if err := Sign(b.Context(), input, out, options); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
