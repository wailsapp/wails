package applesign

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

type machImage struct {
	header                                  []byte
	cpu                                     uint32
	limit, textOffset, textSize, linkOffset int64
	linkCommand, signatureCommand           int
	runtime                                 uint32
	main, hadSignature                      bool
}

// parseMach reads only bounded load-command metadata, never symbol/debug tables.
func parseMach(r io.ReaderAt, size int64) (machImage, error) {
	var m machImage
	var header [32]byte
	if _, err := r.ReadAt(header[:], 0); err != nil {
		return m, err
	}
	le := binary.LittleEndian
	if le.Uint32(header[:]) != 0xfeedfacf {
		return m, errors.New("portable signing requires a little-endian 64-bit Mach-O")
	}
	m.cpu = le.Uint32(header[4:])
	if m.cpu != 0x1000007 && m.cpu != 0x100000c {
		return m, errors.New("portable signing supports amd64 and arm64")
	}
	filetype := le.Uint32(header[12:])
	if filetype != 2 && filetype != 6 && filetype != 8 {
		return m, errors.New("Mach-O is not an executable, dylib, or bundle")
	}
	m.main = filetype == 2
	count, length := int(le.Uint32(header[16:])), int(le.Uint32(header[20:]))
	if length > 4<<20 || int64(length)+32 > size || count > length/8 {
		return m, errors.New("invalid Mach-O load-command bounds")
	}
	m.header = make([]byte, 32+length+16)
	if _, err := r.ReadAt(m.header, 0); err != nil {
		return m, err
	}
	m.signatureCommand, m.linkCommand, m.limit = -1, -1, size
	firstSection := size
	for i, off := 0, 32; i < count; i++ {
		if off+8 > 32+length {
			return m, errors.New("truncated Mach-O command")
		}
		cmd, n := le.Uint32(m.header[off:]), int(le.Uint32(m.header[off+4:]))
		if n < 8 || n%8 != 0 || n > 32+length-off {
			return m, errors.New("invalid Mach-O command size")
		}
		data := m.header[off : off+n]
		switch cmd {
		case 0x19: // LC_SEGMENT_64
			if n < 72 {
				return m, errors.New("truncated Mach-O segment")
			}
			start, span := le.Uint64(data[40:]), le.Uint64(data[48:])
			if start > uint64(size) || span > uint64(size)-start {
				return m, errors.New("Mach-O segment outside file")
			}
			name := string(data[8:24])
			if name == "__LINKEDIT\x00\x00\x00\x00\x00\x00" {
				if m.linkCommand >= 0 {
					return m, errors.New("duplicate LINKEDIT segment")
				}
				m.linkCommand, m.linkOffset = off, int64(start)
			}
			if name == "__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" {
				m.textOffset, m.textSize = int64(start), int64(span)
			}
			sections := int(le.Uint32(data[64:]))
			if sections > (n-72)/80 {
				return m, errors.New("truncated Mach-O sections")
			}
			for s := range sections {
				section := data[72+s*80:]
				kind := le.Uint32(section[64:]) & 255
				if kind == 1 || kind == 12 || kind == 18 {
					continue
				} // zero-fill sections have no file data
				if le.Uint64(section[40:]) > 0 {
					firstSection = min(firstSection, int64(le.Uint32(section[48:])))
				}
			}
		case 0x1d: // LC_CODE_SIGNATURE
			if n != 16 || m.signatureCommand >= 0 {
				return m, errors.New("invalid or repeated code-signature command")
			}
			start, span := int64(le.Uint32(data[8:])), int64(le.Uint32(data[12:]))
			if start < int64(32+length) || start > size || span != size-start {
				return m, errors.New("existing signature must occupy the end of its Mach-O slice")
			}
			m.signatureCommand, m.limit = off, start
			m.hadSignature = true
		case 0x32: // LC_BUILD_VERSION
			if n >= 24 {
				if le.Uint32(data[8:]) != 1 || le.Uint32(data[12:]) < 0x000a0c00 {
					return m, errors.New("portable signing requires macOS 10.12 or later")
				}
				m.runtime = le.Uint32(data[16:])
			}
		case 0x24: // LC_VERSION_MIN_MACOSX
			if n == 16 {
				if le.Uint32(data[8:]) < 0x000a0c00 {
					return m, errors.New("portable signing requires macOS 10.12 or later")
				}
				m.runtime = le.Uint32(data[12:])
			}
		case 0x2c: // LC_ENCRYPTION_INFO_64
			if n >= 20 && le.Uint32(data[16:]) != 0 {
				return m, errors.New("encrypted Mach-O cannot be signed")
			}
		}
		if i == count-1 && off+n != 32+length {
			return m, errors.New("Mach-O command count disagrees with table size")
		}
		off += n
	}
	if m.linkCommand < 0 || m.linkOffset > m.limit {
		return m, errors.New("Mach-O needs a final LINKEDIT segment")
	}
	link := m.header[m.linkCommand:]
	if m.linkOffset+int64(le.Uint64(link[48:])) != size {
		return m, errors.New("LINKEDIT must end at the end of the Mach-O slice")
	}
	if m.signatureCommand < 0 {
		if firstSection < int64(len(m.header)) {
			return m, errors.New("Mach-O has no load-command padding for a signature")
		}
		for _, b := range m.header[32+length:] {
			if b != 0 {
				return m, errors.New("Mach-O load-command padding is not empty")
			}
		}
		m.signatureCommand = 32 + length
		le.PutUint32(m.header[16:], uint32(count+1))
		le.PutUint32(m.header[20:], uint32(length+16))
		le.PutUint32(m.header[m.signatureCommand:], 0x1d)
		le.PutUint32(m.header[m.signatureCommand+4:], 16)
	} else {
		m.header = m.header[:32+length]
	}
	return m, nil
}

func signThin(ctx context.Context, source io.ReaderAt, size int64, out *os.File, o Options) (CodeIdentity, error) {
	m, err := parseMach(source, size)
	if err != nil {
		return CodeIdentity{}, err
	}
	components, err := signatureComponents(o)
	if err != nil {
		return CodeIdentity{}, err
	}
	limit := (m.limit + 15) &^ 15
	if limit > 0xffffffff {
		return CodeIdentity{}, errors.New("Mach-O signing limit exceeds 4 GiB")
	}
	// The reserved CMS region keeps load commands stable while a timestamp
	// service returns a variable-sized token. It is outside the signed pages.
	dirSize := directorySize(limit, 12, o, components)
	signatureSize, err := blobIndexSize(len(components) + 2)
	if err != nil {
		return CodeIdentity{}, err
	}
	for _, n := range []int{dirSize, 8} {
		signatureSize, err = addBlobSize(signatureSize, n)
		if err != nil {
			return CodeIdentity{}, err
		}
	}
	for _, c := range components {
		signatureSize, err = addBlobSize(signatureSize, len(c.data))
		if err != nil {
			return CodeIdentity{}, err
		}
	}
	if o.SignCMS != nil {
		signatureSize, err = addBlobSize(signatureSize, 128<<10)
		if err != nil {
			return CodeIdentity{}, err
		}
	}
	signatureSize, err = addBlobSize(signatureSize, (-signatureSize)&15)
	if err != nil {
		return CodeIdentity{}, err
	}
	le := binary.LittleEndian
	le.PutUint32(m.header[m.signatureCommand+8:], uint32(limit))
	le.PutUint32(m.header[m.signatureCommand+12:], uint32(signatureSize))
	linkSize := limit + int64(signatureSize) - m.linkOffset
	le.PutUint64(m.header[m.linkCommand+48:], uint64(linkSize))
	page := int64(4096)
	if m.cpu == 0x100000c {
		page = 16384
	}
	le.PutUint64(m.header[m.linkCommand+32:], uint64((linkSize+page-1)&-page))
	if _, err := io.CopyN(out, &contextReader{ctx, io.NewSectionReader(source, 0, m.limit)}, m.limit); err != nil {
		return CodeIdentity{}, err
	}
	if err := out.Truncate(limit + int64(signatureSize)); err != nil {
		return CodeIdentity{}, err
	}
	if _, err := out.WriteAt(m.header, 0); err != nil {
		return CodeIdentity{}, err
	}
	directory, err := codeDirectory(ctx, out, limit, 12, m, o, components)
	if err != nil {
		return CodeIdentity{}, err
	}
	signature, err := signatureBlob(ctx, directory, components, o)
	if err != nil {
		return CodeIdentity{}, err
	}
	if len(signature) > signatureSize {
		return CodeIdentity{}, fmt.Errorf("CMS signature exceeds reserved size")
	}
	if _, err := out.WriteAt(signature, limit); err != nil {
		return CodeIdentity{}, err
	}
	return identity(m.cpu, o.Identifier, directory), ctx.Err()
}
