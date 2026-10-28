// Package uuid provides minimal UUID generation (v4 random, v5 name-based SHA1).
package uuid

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
)

// UUID is a 16-byte UUID value.
type UUID [16]byte

// Nil is the all-zeros UUID.
var Nil UUID

// New returns a random UUID (version 4).
func New() UUID {
	var id UUID
	rand.Read(id[:]) //nolint:errcheck // crypto/rand.Read never returns an error in Go 1.20+
	id[6] = (id[6] & 0x0f) | 0x40 // version 4
	id[8] = (id[8] & 0x3f) | 0x80 // variant RFC 4122
	return id
}

// NewSHA1 returns a deterministic name-based UUID (version 5) derived from namespace and data.
func NewSHA1(ns UUID, data []byte) UUID {
	h := sha1.New()
	h.Write(ns[:])
	h.Write(data)
	var digest [sha1.Size]byte
	sum := h.Sum(digest[:0])

	var id UUID
	copy(id[:], sum)
	id[6] = (id[6] & 0x0f) | 0x50 // version 5
	id[8] = (id[8] & 0x3f) | 0x80 // variant RFC 4122
	return id
}

// String returns the standard UUID string representation.
func (id UUID) String() string {
	var text [36]byte
	hex.Encode(text[0:8], id[0:4])
	hex.Encode(text[9:13], id[4:6])
	hex.Encode(text[14:18], id[6:8])
	hex.Encode(text[19:23], id[8:10])
	hex.Encode(text[24:36], id[10:16])
	text[8], text[13], text[18], text[23] = '-', '-', '-', '-'
	return string(text[:])
}
