package notaryticket

import (
	"encoding/binary"
	"errors"
	"os"
	"sort"
)

type component struct {
	slot uint32
	data []byte
}

// DMG tickets are signature components, not package trailers. The existing
// CodeDirectory and CMS bytes remain unchanged; only the unsigned container
// indexes and the UDIF signature-length field change.
func attachDMG(path string, ticket []byte) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.Size() < 512 {
		return errors.New("truncated UDIF disk image")
	}
	var footer [512]byte
	if _, err = f.ReadAt(footer[:], info.Size()-512); err != nil {
		return err
	}
	be := binary.BigEndian
	if string(footer[:4]) != "koly" || be.Uint32(footer[4:]) != 4 || be.Uint32(footer[8:]) != 512 {
		return errors.New("unsupported UDIF disk image")
	}
	offset, size := be.Uint64(footer[296:]), be.Uint64(footer[304:])
	if offset == 0 || offset > uint64(info.Size()-512) || size != uint64(info.Size()-512)-offset || size > 64<<20 {
		return errors.New("invalid DMG signature extent")
	}
	data := make([]byte, size)
	if _, err = f.ReadAt(data, int64(offset)); err != nil {
		return err
	}
	parts, err := signatureParts(data)
	if err != nil {
		return err
	}
	blob := make([]byte, 8+len(ticket))
	be.PutUint32(blob, 0xfade0b01)
	be.PutUint32(blob[4:], uint32(len(blob)))
	copy(blob[8:], ticket)
	parts = append(parts, component{0x10002, blob})
	sort.Slice(parts, func(i, j int) bool { return parts[i].slot < parts[j].slot })
	length := 12 + 8*len(parts)
	for _, part := range parts {
		length += len(part.data)
	}
	if length > 64<<20 {
		return errors.New("stapled DMG signature exceeds 64 MiB")
	}
	result := make([]byte, length)
	be.PutUint32(result, 0xfade0cc0)
	be.PutUint32(result[4:], uint32(length))
	be.PutUint32(result[8:], uint32(len(parts)))
	cursor := 12 + 8*len(parts)
	for i, part := range parts {
		be.PutUint32(result[12+8*i:], part.slot)
		be.PutUint32(result[16+8*i:], uint32(cursor))
		copy(result[cursor:], part.data)
		cursor += len(part.data)
	}
	be.PutUint64(footer[304:], uint64(length))
	if _, err = f.WriteAt(result, int64(offset)); err != nil {
		return err
	}
	if _, err = f.WriteAt(footer[:], int64(offset)+int64(length)); err != nil {
		return err
	}
	if err = f.Truncate(int64(offset) + int64(length) + 512); err != nil {
		return err
	}
	return f.Sync()
}

func signatureParts(data []byte) ([]component, error) {
	be := binary.BigEndian
	if len(data) < 12 || be.Uint32(data) != 0xfade0cc0 || uint64(be.Uint32(data[4:])) != uint64(len(data)) {
		return nil, errors.New("invalid DMG signature container")
	}
	count := int(be.Uint32(data[8:]))
	if count < 2 || count > 64 || 12+count*8 > len(data) {
		return nil, errors.New("invalid DMG signature component count")
	}
	seen := map[uint32]bool{}
	var parts []component
	type extent struct{ start, end int }
	var extents []extent
	for i := 0; i < count; i++ {
		slot, offset := be.Uint32(data[12+8*i:]), int(be.Uint32(data[16+8*i:]))
		if seen[slot] || offset < 12+count*8 || offset > len(data)-8 {
			return nil, errors.New("duplicate or out-of-bounds DMG component")
		}
		seen[slot] = true
		size := int(be.Uint32(data[offset+4:]))
		if size < 8 || size > len(data)-offset {
			return nil, errors.New("invalid DMG component length")
		}
		extents = append(extents, extent{offset, offset + size})
		if slot != 0x10002 {
			parts = append(parts, component{slot, data[offset : offset+size]})
		}
	}
	sort.Slice(extents, func(i, j int) bool { return extents[i].start < extents[j].start })
	for i := 1; i < len(extents); i++ {
		if extents[i].start < extents[i-1].end {
			return nil, errors.New("overlapping DMG signature components")
		}
	}
	if !seen[0] || !seen[0x10000] {
		return nil, errors.New("DMG must contain a CodeDirectory and CMS signature")
	}
	return parts, nil
}
