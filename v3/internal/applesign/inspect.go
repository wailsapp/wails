package applesign

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// Inspect returns code-directory identities for ticket lookup. It validates
// structural bounds, not certificate trust or the integrity of signed content.
func Inspect(path string) ([]CodeIdentity, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, errors.New("signature input must be a regular file")
	}
	if st.Size() >= 512 {
		var trailer [512]byte
		if _, err := f.ReadAt(trailer[:], st.Size()-512); err != nil {
			return nil, err
		}
		if string(trailer[:4]) == "koly" {
			offset, size := binary.BigEndian.Uint64(trailer[296:]), binary.BigEndian.Uint64(trailer[304:])
			if offset > uint64(st.Size()-512) || size != uint64(st.Size()-512)-offset {
				return nil, errors.New("invalid DMG signature bounds")
			}
			id, err := inspectBlob(io.NewSectionReader(f, int64(offset), int64(size)), int64(size), 0)
			if err != nil {
				return nil, err
			}
			return []CodeIdentity{id}, nil
		}
	}
	parts, err := readSlices(f, st.Size())
	if err != nil {
		return nil, err
	}
	result := make([]CodeIdentity, 0, len(parts))
	for _, part := range parts {
		r := io.NewSectionReader(f, part.offset, part.size)
		m, err := parseMach(r, part.size)
		if err != nil {
			return nil, err
		}
		if !m.hadSignature {
			return nil, errors.New("Mach-O is unsigned")
		}
		id, err := inspectBlob(io.NewSectionReader(r, m.limit, part.size-m.limit), part.size-m.limit, m.cpu)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}

func inspectBlob(r io.ReaderAt, size int64, cpu uint32) (CodeIdentity, error) {
	var zero CodeIdentity
	var header [12]byte
	if _, err := r.ReadAt(header[:], 0); err != nil {
		return zero, err
	}
	be := binary.BigEndian
	length, count := int64(be.Uint32(header[4:])), int(be.Uint32(header[8:]))
	if be.Uint32(header[:]) != 0xfade0cc0 || length > size || length > maxSignatureBlobSize || count > maxSignatureComponents || length < int64(12+8*count) {
		return zero, errors.New("invalid embedded signature bounds")
	}
	data := make([]byte, length)
	if _, err := r.ReadAt(data, 0); err != nil {
		return zero, err
	}
	var directory []byte
	for i := range count {
		slot, offset := be.Uint32(data[12+8*i:]), int(be.Uint32(data[16+8*i:]))
		if offset < 12+8*count || offset > len(data)-8 {
			return zero, errors.New("invalid signature component offset")
		}
		n := int(be.Uint32(data[offset+4:]))
		if n < 8 || n > len(data)-offset {
			return zero, errors.New("invalid signature component length")
		}
		if slot == 0 {
			if directory != nil {
				return zero, errors.New("duplicate primary CodeDirectory")
			}
			directory = data[offset : offset+n]
		}
	}
	if len(directory) < 44 || be.Uint32(directory) != 0xfade0c02 || directory[37] != 2 || directory[36] != 32 {
		return zero, errors.New("signature has no supported SHA256 CodeDirectory")
	}
	offset := int(be.Uint32(directory[20:]))
	if offset < 44 || offset >= len(directory) {
		return zero, errors.New("invalid signature identifier")
	}
	end := bytes.IndexByte(directory[offset:], 0)
	if end < 0 || end > 1024 {
		return zero, errors.New("unterminated signature identifier")
	}
	return identity(cpu, string(directory[offset:offset+end]), directory), nil
}
