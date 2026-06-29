// Package uuid provides minimal UUID generation (v4 random, v5 name-based SHA1).
package uuid

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
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
	sum := h.Sum(nil)

	var id UUID
	copy(id[:], sum)
	id[6] = (id[6] & 0x0f) | 0x50 // version 5
	id[8] = (id[8] & 0x3f) | 0x80 // variant RFC 4122
	return id
}

// String returns the standard UUID string representation.
func (id UUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		id[0:4], id[4:6], id[6:8], id[8:10], id[10:16])
}
