package appimage

import (
	"debug/elf"
	"fmt"
	"io"
	"os"
)

// preferBundledLibraries lets the launcher's LD_LIBRARY_PATH override build
// machine RPATHs without resizing ELF segments or rewriting executable code.
func preferBundledLibraries(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	image, err := elf.NewFile(file)
	if err != nil {
		return nil
	} // Data files need no loader changes.
	for _, program := range image.Progs {
		if program.Type != elf.PT_DYNAMIC {
			continue
		}
		if program.Filesz > 16<<20 {
			return fmt.Errorf("oversized ELF dynamic table: %s", filename)
		}
		data, err := io.ReadAll(program.Open())
		if err != nil {
			return err
		}
		width := 8
		if image.Class == elf.ELFCLASS64 {
			width = 16
		}
		var offsets []int
		for offset := 0; offset+width <= len(data); offset += width {
			var tag uint64
			if width == 16 {
				tag = image.ByteOrder.Uint64(data[offset:])
			} else {
				tag = uint64(image.ByteOrder.Uint32(data[offset:]))
			}
			if tag == uint64(elf.DT_NULL) {
				break
			}
			if tag == uint64(elf.DT_RUNPATH) {
				return nil
			}
			if tag == uint64(elf.DT_RPATH) {
				offsets = append(offsets, offset)
			}
		}
		for _, offset := range offsets {
			if width == 16 {
				image.ByteOrder.PutUint64(data[offset:], uint64(elf.DT_RUNPATH))
			} else {
				image.ByteOrder.PutUint32(data[offset:], uint32(elf.DT_RUNPATH))
			}
			if _, err := file.WriteAt(data[offset:offset+width/2], int64(program.Off)+int64(offset)); err != nil {
				return err
			}
		}
	}
	return nil
}
