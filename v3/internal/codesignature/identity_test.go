package codesignature

import (
	"bytes"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"os"
	"testing"
)

func TestDecodeOpenSSLPFX(t *testing.T) {
	for _, name := range []string{"modern", "legacy-3des"} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile("testdata/" + name + ".pfx")
			if err != nil {
				t.Fatal(err)
			}
			before := sha256.Sum256(data)
			identity, err := DecodePFX(data, "fixture-only")
			if err != nil {
				t.Fatal(err)
			}
			if identity.Certificate.Subject.CommonName != "Wails Portable Signing Test Only" {
				t.Fatal("wrong certificate")
			}
			if before != sha256.Sum256(data) {
				t.Fatal("PFX input mutated")
			}
			if _, err := DecodePFX(data, "wrong"); err == nil {
				t.Fatal("wrong password accepted")
			}
			altered := bytes.Clone(data)
			altered[len(altered)-50] ^= 1
			if _, err := DecodePFX(altered, "fixture-only"); err == nil {
				t.Fatal("tampered PFX accepted")
			}
		})
	}
}

func TestPFXAuthenticatedContentsAndWorkBounds(t *testing.T) {
	data, err := os.ReadFile("testdata/modern.pfx")
	if err != nil {
		t.Fatal(err)
	}
	var document pfxDocument
	if err := decodeDER(data, &document); err != nil {
		t.Fatal(err)
	}
	var safe []byte
	if err := decodeDER(document.Safe.Content.Bytes, &safe); err != nil {
		t.Fatal(err)
	}
	safe[len(safe)/2] ^= 1
	encoded, err := asn1.Marshal(safe)
	if err != nil {
		t.Fatal(err)
	}
	document.Safe.Content = explicitContent(encoded)
	altered, err := asn1.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodePFX(altered, "fixture-only"); err == nil {
		t.Fatal("modified authenticated safe accepted")
	}
	decoder := pfxDecoder{budget: 2_000_000}
	for _, count := range []int{0, -1, 1_000_001, int(^uint(0) >> 1)} {
		if err := decoder.useIterations(count); err == nil {
			t.Fatalf("accepted excessive iteration count %d", count)
		}
	}
}

func BenchmarkDecodePFX(b *testing.B) {
	data, err := os.ReadFile("testdata/modern.pfx")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := DecodePFX(data, "fixture-only"); err != nil {
			b.Fatal(err)
		}
	}
}

func TestUnsupportedPFXEncryptionRequiresValidIntegrity(t *testing.T) {
	data, err := os.ReadFile("testdata/legacy-rc2.pfx")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodePFX(data, "fixture-only"); !errors.Is(err, ErrUnsupportedPFXEncryption) {
		t.Fatalf("RC2 classification: %v", err)
	}
	if _, err := DecodePFX(data, "wrong"); err == nil || errors.Is(err, ErrUnsupportedPFXEncryption) {
		t.Fatalf("wrong password allowed fallback: %v", err)
	}
	data[len(data)-50] ^= 1
	if _, err := DecodePFX(data, "fixture-only"); err == nil || errors.Is(err, ErrUnsupportedPFXEncryption) {
		t.Fatalf("tampering allowed fallback: %v", err)
	}
}
