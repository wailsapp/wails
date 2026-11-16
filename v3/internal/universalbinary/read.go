package universalbinary

import (
	"debug/macho"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func readImages(f *os.File, size int64) ([]image, error) {
	var header [8]byte
	if _, err := f.ReadAt(header[:], 0); err != nil {
		return nil, err
	}
	var order binary.ByteOrder = binary.BigEndian
	magic := order.Uint32(header[:4])
	if magic != macho.MagicFat && magic != 0xcafebabf {
		order = binary.LittleEndian
		magic = order.Uint32(header[:4])
	}
	if magic != macho.MagicFat && magic != 0xcafebabf {
		img, err := readImage(f, 0, size)
		if err != nil {
			return nil, err
		}
		return []image{img}, nil
	}
	count := int64(order.Uint32(header[4:]))
	entrySize := int64(20)
	if magic == 0xcafebabf {
		entrySize = 32
	}
	tableEnd := 8 + count*entrySize
	if count == 0 || count > 4096 || tableEnd > size {
		return nil, fmt.Errorf("invalid universal architecture table")
	}
	images := make([]image, 0, count)
	for i := int64(0); i < count; i++ {
		var entry [32]byte
		if _, err := f.ReadAt(entry[:entrySize], 8+i*entrySize); err != nil {
			return nil, err
		}
		offset, length := uint64(order.Uint32(entry[8:])), uint64(order.Uint32(entry[12:]))
		align := order.Uint32(entry[16:])
		if entrySize == 32 {
			offset, length = order.Uint64(entry[8:]), order.Uint64(entry[16:])
			align = order.Uint32(entry[24:])
		}
		if offset < uint64(tableEnd) || offset > uint64(size) || length > uint64(size)-offset || length == 0 || align > 30 || offset%(uint64(1)<<align) != 0 {
			return nil, fmt.Errorf("invalid universal slice bounds or alignment")
		}
		img, err := readImage(f, int64(offset), int64(length))
		if err != nil {
			return nil, err
		}
		if img.cpu != order.Uint32(entry[:4]) || img.sub != order.Uint32(entry[4:8]) {
			return nil, fmt.Errorf("universal slice architecture disagrees with its header")
		}
		img.align = max(img.align, align)
		images = append(images, img)
	}
	return images, nil
}

func readImage(f *os.File, offset, size int64) (image, error) {
	img := image{file: f, offset: offset, size: size}
	var magic [8]byte
	if _, err := io.NewSectionReader(f, offset, size).ReadAt(magic[:], 0); err != nil {
		return img, err
	}
	if string(magic[:]) == "!<arch>\n" {
		cpu, sub, err := archiveArchitecture(io.NewSectionReader(f, offset, size), size)
		img.cpu, img.sub, img.align = cpu, sub, 2
		return img, err
	}
	var err error
	img.cpu, img.sub, img.align, err = machoArchitecture(io.NewSectionReader(f, offset, size), size)
	if err != nil {
		return img, err
	}
	return img, nil
}

func archiveArchitecture(r io.ReaderAt, size int64) (uint32, uint32, error) {
	var cpu, sub uint32
	found := false
	for offset := int64(8); offset < size; {
		var header [60]byte
		if _, err := r.ReadAt(header[:], offset); err != nil {
			return 0, 0, err
		}
		length, err := strconv.ParseInt(strings.TrimSpace(string(header[48:58])), 10, 64)
		data := offset + 60
		if err != nil || string(header[58:]) != "`\n" || length < 0 || length > size-data {
			return 0, 0, fmt.Errorf("invalid archive member")
		}
		next := data + length + (length % 2)
		name := strings.TrimSpace(string(header[:16]))
		if strings.HasPrefix(name, "#1/") {
			nameSize, err := strconv.ParseInt(name[3:], 10, 64)
			if err != nil || nameSize < 0 || nameSize > length {
				return 0, 0, fmt.Errorf("invalid archive extended name")
			}
			data, length = data+nameSize, length-nameSize
		}
		var magic [4]byte
		if length >= 4 {
			if _, err := r.ReadAt(magic[:], data); err != nil {
				return 0, 0, err
			}
			le, be := binary.LittleEndian.Uint32(magic[:]), binary.BigEndian.Uint32(magic[:])
			if le == macho.Magic32 || le == macho.Magic64 || be == macho.Magic32 || be == macho.Magic64 {
				objCPU, objSub, _, err := machoArchitecture(io.NewSectionReader(r, data, length), length)
				if err != nil {
					return 0, 0, err
				}
				if found && (cpu != objCPU || sub != objSub) {
					return 0, 0, fmt.Errorf("archive mixes architectures")
				}
				cpu, sub, found = objCPU, objSub, true
			}
		}
		offset = next
	}
	if !found {
		return 0, 0, fmt.Errorf("archive has no Mach-O objects")
	}
	return cpu, sub, nil
}
