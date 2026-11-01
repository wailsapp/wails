package applesign

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"slices"
)

type component struct {
	slot uint32
	data []byte
}

// Match the limits accepted by Inspect; these also keep uint32 format lengths
// and host int allocations representable on both 32-bit and 64-bit hosts.
const maxSignatureBlobSize = 64 << 20
const maxSignatureComponents = 64

func addBlobSize(size, extra int) (int, error) {
	if size < 0 || size > maxSignatureBlobSize || extra < 0 || extra > maxSignatureBlobSize-size {
		return 0, errors.New("signature blob exceeds 64 MiB")
	}
	return size + extra, nil
}

func blobIndexSize(count int) (int, error) {
	if count < 0 || count > maxSignatureComponents {
		return 0, errors.New("signature blob exceeds 64 components")
	}
	return 12 + 8*count, nil
}

func wrap(magic uint32, data []byte) ([]byte, error) {
	size, err := addBlobSize(8, len(data))
	if err != nil {
		return nil, err
	}
	result := make([]byte, size)
	binary.BigEndian.PutUint32(result, magic)
	binary.BigEndian.PutUint32(result[4:], uint32(size))
	copy(result[8:], data)
	return result, nil
}

func superblob(magic uint32, parts []component) ([]byte, error) {
	header, err := blobIndexSize(len(parts))
	if err != nil {
		return nil, err
	}
	size := header
	for _, p := range parts {
		size, err = addBlobSize(size, len(p.data))
		if err != nil {
			return nil, err
		}
	}
	result := make([]byte, size)
	be := binary.BigEndian
	be.PutUint32(result, magic)
	be.PutUint32(result[4:], uint32(size))
	be.PutUint32(result[8:], uint32(len(parts)))
	offset := header
	for i, p := range parts {
		be.PutUint32(result[12+8*i:], p.slot)
		be.PutUint32(result[16+8*i:], uint32(offset))
		copy(result[offset:], p.data)
		offset += len(p.data)
	}
	return result, nil
}

func signatureComponents(o Options) ([]component, error) {
	requirements, err := superblob(0xfade0c01, nil)
	if o.requirement != nil {
		requirements, err = o.requirement(o.Identifier)
	}
	if err != nil {
		return nil, err
	}
	parts := []component{{2, requirements}}
	if len(o.Entitlements) != 0 {
		der, err := derEntitlements(o.Entitlements)
		if err != nil {
			return nil, err
		}
		xmlBlob, err := wrap(0xfade7171, o.Entitlements)
		if err != nil {
			return nil, err
		}
		derBlob, err := wrap(0xfade7172, der)
		if err != nil {
			return nil, err
		}
		parts = append(parts, component{5, xmlBlob}, component{7, derBlob})
	}
	return parts, nil
}

func specialSlots(o Options, parts []component) map[uint32][]byte {
	result := make(map[uint32][]byte, len(parts)+2)
	for _, p := range parts {
		result[p.slot] = p.data
	}
	if o.info != nil {
		result[1] = o.info
	}
	if o.resources != nil {
		result[3] = o.resources
	}
	if o.representation != nil {
		result[6] = o.representation
	}
	return result
}

func directorySize(limit int64, page uint8, o Options, parts []component) int {
	maxSlot := uint32(0)
	for slot := range specialSlots(o, parts) {
		maxSlot = max(maxSlot, slot)
	}
	count := int64(1)
	if page != 0 {
		count = (limit + (1 << page) - 1) >> page
	}
	team := 0
	if o.TeamID != "" {
		team = len(o.TeamID) + 1
	}
	return 96 + len(o.Identifier) + 1 + team + 32*(int(maxSlot)+int(count))
}

func codeDirectory(ctx context.Context, r io.ReaderAt, limit int64, page uint8, m machImage, o Options, parts []component) ([]byte, error) {
	special := specialSlots(o, parts)
	maxSlot := uint32(0)
	for slot := range special {
		maxSlot = max(maxSlot, slot)
	}
	d := make([]byte, directorySize(limit, page, o, parts))
	be := binary.BigEndian
	be.PutUint32(d, 0xfade0c02)
	be.PutUint32(d[4:], uint32(len(d)))
	be.PutUint32(d[8:], 0x20500)
	flags := uint32(0)
	if o.SignCMS == nil {
		flags |= 2
	}
	if o.HardenedRuntime {
		flags |= 0x10000
	}
	be.PutUint32(d[12:], flags)
	be.PutUint32(d[20:], 96)
	copy(d[96:], o.Identifier)
	stringsEnd := 96 + len(o.Identifier) + 1
	if o.TeamID != "" {
		be.PutUint32(d[48:], uint32(stringsEnd))
		copy(d[stringsEnd:], o.TeamID)
		stringsEnd += len(o.TeamID) + 1
	}
	hashOffset := stringsEnd + int(maxSlot)*32
	be.PutUint32(d[16:], uint32(hashOffset))
	be.PutUint32(d[24:], maxSlot)
	count := 1
	if page != 0 {
		count = int((limit + (1 << page) - 1) >> page)
	}
	be.PutUint32(d[28:], uint32(count))
	be.PutUint32(d[32:], uint32(limit))
	d[36], d[37], d[39] = 32, 2, page
	be.PutUint64(d[64:], uint64(m.textOffset))
	be.PutUint64(d[72:], uint64(m.textSize))
	if m.main {
		be.PutUint64(d[80:], 1)
	}
	if o.HardenedRuntime {
		be.PutUint32(d[88:], m.runtime)
	}
	for slot, data := range special {
		sum := sha256.Sum256(data)
		copy(d[hashOffset-int(slot)*32:], sum[:])
	}
	block := limit
	if page != 0 {
		block = 1 << page
	}
	buffer := make([]byte, 32<<10)
	for i := range count {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		h := sha256.New()
		start := int64(i) * block
		if _, err := io.CopyBuffer(h, &contextReader{ctx, io.NewSectionReader(r, start, min(block, limit-start))}, buffer); err != nil {
			return nil, err
		}
		copy(d[hashOffset+i*32:], h.Sum(nil))
	}
	return d, nil
}

func signatureBlob(ctx context.Context, directory []byte, parts []component, o Options) ([]byte, error) {
	if len(parts) > maxSignatureComponents-2 {
		return nil, errors.New("signature blob exceeds 64 components")
	}
	var cms []byte
	if o.SignCMS != nil {
		var err error
		cms, err = o.SignCMS(ctx, slices.Clone(directory))
		if err != nil {
			return nil, err
		}
		if len(cms) == 0 || len(cms) > 128<<10 {
			return nil, errors.New("CMS signer returned an empty or oversized signature")
		}
	}
	all := make([]component, 0, len(parts)+2)
	all = append(all, component{0, directory})
	all = append(all, parts...)
	cmsBlob, err := wrap(0xfade0b01, cms)
	if err != nil {
		return nil, err
	}
	all = append(all, component{65536, cmsBlob})
	return superblob(0xfade0cc0, all)
}
