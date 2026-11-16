package universalbinary

import (
	"debug/macho"
	"encoding/binary"
	"fmt"
	"io"
)

// Inspect load-command metadata without loading symbol tables or executable contents.
func machoArchitecture(r io.ReaderAt, size int64) (cpu, sub, align uint32, err error) {
	var header [28]byte
	if _, err = r.ReadAt(header[:], 0); err != nil {
		return
	}
	var order binary.ByteOrder = binary.LittleEndian
	magic := order.Uint32(header[:])
	if magic != macho.Magic32 && magic != macho.Magic64 {
		order = binary.BigEndian
		magic = order.Uint32(header[:])
	}
	if magic != macho.Magic32 && magic != macho.Magic64 {
		err = fmt.Errorf("not a Mach-O binary")
		return
	}
	cpu, sub = order.Uint32(header[4:]), order.Uint32(header[8:])
	align = 16
	if macho.Cpu(cpu) == macho.Cpu386 || macho.Cpu(cpu) == macho.CpuAmd64 {
		align = 12
	}
	if macho.Cpu(cpu) == macho.CpuArm64 {
		align = 14
	}
	if macho.Type(order.Uint32(header[12:])) == macho.TypeObj {
		align = 2
	}
	offset := int64(28)
	if magic == macho.Magic64 {
		offset = 32
	}
	length, count := int64(order.Uint32(header[20:])), int64(order.Uint32(header[16:]))
	if size < offset || length > size-offset || count > length/8 {
		err = fmt.Errorf("invalid Mach-O load-command table")
		return
	}
	end := offset + length
	for i := int64(0); i < count; i++ {
		var command [8]byte
		if _, err = r.ReadAt(command[:], offset); err != nil {
			return
		}
		cmd, n := macho.LoadCmd(order.Uint32(command[:])), int64(order.Uint32(command[4:]))
		if n < 8 || n > end-offset {
			err = fmt.Errorf("invalid Mach-O load command size")
			return
		}
		if cmd == macho.LoadCmdSegment || cmd == macho.LoadCmdSegment64 {
			var sectionAlign uint32
			sectionAlign, err = segmentAlignment(r, order, offset, n, cmd == macho.LoadCmdSegment64)
			if err != nil {
				return
			}
			align = max(align, sectionAlign)
		}
		offset += n
	}
	return
}

func segmentAlignment(r io.ReaderAt, order binary.ByteOrder, offset, length int64, wide bool) (uint32, error) {
	headerSize, sectionSize, countOffset, alignOffset := int64(56), int64(68), 48, int64(44)
	if wide {
		headerSize, sectionSize, countOffset, alignOffset = 72, 80, 64, 52
	}
	if length < headerSize {
		return 0, fmt.Errorf("truncated Mach-O segment")
	}
	var header [72]byte
	if _, err := r.ReadAt(header[:headerSize], offset); err != nil {
		return 0, err
	}
	count := int64(order.Uint32(header[countOffset:]))
	if count > (length-headerSize)/sectionSize {
		return 0, fmt.Errorf("truncated Mach-O sections")
	}
	var align uint32
	for i := int64(0); i < count; i++ {
		var field [4]byte
		if _, err := r.ReadAt(field[:], offset+headerSize+i*sectionSize+alignOffset); err != nil {
			return 0, err
		}
		align = max(align, order.Uint32(field[:]))
	}
	return align, nil
}
