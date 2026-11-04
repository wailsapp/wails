package applesign

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"math"
	"strings"
	"testing"
)

func TestBlobSizeBounds(t *testing.T) {
	for _, test := range []struct {
		base, extra, want int
		valid             bool
	}{
		{8, 0, 8, true},
		{8, maxSignatureBlobSize - 8, maxSignatureBlobSize, true},
		{maxSignatureBlobSize, 0, maxSignatureBlobSize, true},
		{8, maxSignatureBlobSize - 7, 0, false},
		{maxSignatureBlobSize, 1, 0, false},
		{math.MaxInt, 1, 0, false},
		{8, math.MaxInt, 0, false},
		{math.MaxInt, math.MaxInt, 0, false},
		{-1, 8, 0, false},
		{8, -1, 0, false},
	} {
		got, err := addBlobSize(test.base, test.extra)
		if (err == nil) != test.valid || got != test.want {
			t.Fatalf("size(%d, %d) = %d, %v", test.base, test.extra, got, err)
		}
	}
	// Model cumulative superblob payloads without allocating large buffers.
	size, err := blobIndexSize(2)
	if err != nil {
		t.Fatal(err)
	}
	size, err = addBlobSize(size, maxSignatureBlobSize-size-1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := addBlobSize(size, 2); err == nil {
		t.Fatal("combined payload overflow accepted")
	}
	for _, count := range []int{-1, maxSignatureComponents + 1, math.MaxInt} {
		if _, err := blobIndexSize(count); err == nil {
			t.Fatalf("invalid component count %d accepted", count)
		}
	}
	if size, err := blobIndexSize(maxSignatureComponents); err != nil || size != 524 {
		t.Fatalf("maximum index: %d, %v", size, err)
	}
	// The largest supported Mach-O input has at most 2^20 page hashes.
	directory := directorySize(0xffffffff, 12, Options{Identifier: strings.Repeat("i", 1024), TeamID: strings.Repeat("t", 128)}, []component{{7, nil}})
	if directory != 96+1025+129+32*(7+1048576) || directory > maxSignatureBlobSize {
		t.Fatalf("maximum supported code directory size: %d", directory)
	}
}

func TestBlobEncodingAndErrorPropagation(t *testing.T) {
	data := []byte("payload")
	wrapped, err := wrap(0xfade0b01, data)
	if err != nil {
		t.Fatal(err)
	}
	be := binary.BigEndian
	if be.Uint32(wrapped) != 0xfade0b01 || int(be.Uint32(wrapped[4:])) != len(wrapped) || !bytes.Equal(wrapped[8:], data) {
		t.Fatal("incorrect blob envelope")
	}
	encoded, err := superblob(0xfade0cc0, []component{{5, wrapped}, {7, wrapped}})
	if err != nil {
		t.Fatal(err)
	}
	if int(be.Uint32(encoded[4:])) != len(encoded) || be.Uint32(encoded[8:]) != 2 {
		t.Fatal("incorrect superblob header")
	}
	if be.Uint32(encoded[12:]) != 5 || be.Uint32(encoded[16:]) != 28 || be.Uint32(encoded[20:]) != 7 || int(be.Uint32(encoded[24:])) != 28+len(wrapped) {
		t.Fatal("incorrect superblob index")
	}
	if !bytes.Equal(encoded[28:], append(bytes.Clone(wrapped), wrapped...)) {
		t.Fatal("incorrect superblob payload")
	}
	if _, err := superblob(0, make([]component, maxSignatureComponents+1)); err == nil {
		t.Fatal("oversized component index accepted")
	}
	sentinel := errors.New("requirement generation failed")
	_, err = signatureComponents(Options{requirement: func(string) ([]byte, error) { return nil, sentinel }})
	if !errors.Is(err, sentinel) {
		t.Fatalf("requirement error lost: %v", err)
	}
	_, err = signatureBlob(context.Background(), nil, make([]component, maxSignatureComponents-1), Options{SignCMS: func(context.Context, []byte) ([]byte, error) {
		t.Fatal("signer called for an invalid component index")
		return nil, nil
	}})
	if err == nil {
		t.Fatal("signature component overflow accepted")
	}
	if _, err := designatedRequirement(strings.Repeat("a", 1025), "", true, nil); err == nil {
		t.Fatal("unbounded requirement identifier accepted")
	}
}
