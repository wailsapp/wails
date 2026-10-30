// Package universalbinary combines Mach-O binaries and static archives into universal files.
package universalbinary

import (
	"cmp"
	"debug/macho"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
)

type image struct {
	file            *os.File
	cpu, sub, align uint32
	offset, size    int64
}

// Create combines all architectures from inputs, publishing output only after a complete write.
func Create(output string, inputs []string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("no input binaries")
	}
	var images []image
	var files []*os.File
	defer func() {
		for _, f := range files {
			_ = f.Close()
		}
	}()
	var mode os.FileMode
	existing, statErr := os.Stat(output)
	if statErr != nil && !os.IsNotExist(statErr) {
		return statErr
	}
	for _, path := range inputs {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		files = append(files, f)
		info, err := f.Stat()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("input %q is not a regular file", path)
		}
		if existing != nil && os.SameFile(info, existing) {
			return fmt.Errorf("output aliases input %q", path)
		}
		if len(files) == 1 {
			mode = info.Mode().Perm()
		}
		entries, err := readImages(f, info.Size())
		if err != nil {
			return fmt.Errorf("inspect %q: %w", path, err)
		}
		images = append(images, entries...)
	}
	header, offsets, err := layout(images)
	if err != nil {
		return err
	}
	return publish(output, images, header, offsets, mode)
}

func publish(output string, images []image, header []byte, offsets []int64, mode os.FileMode) error {
	stage, err := os.MkdirTemp(filepath.Dir(output), ".wails-universal-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	path := filepath.Join(stage, "binary")
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err = out.Write(header); err != nil {
		return err
	}
	for i, img := range images {
		if _, err = out.Seek(offsets[i], io.SeekStart); err != nil {
			return err
		}
		if _, err = img.file.Seek(img.offset, io.SeekStart); err != nil {
			return err
		}
		if _, err = io.CopyN(out, img.file, img.size); err != nil {
			return err
		}
	}
	if err = out.Chmod(mode); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}
	return os.Rename(path, output)
}

func layout(images []image) ([]byte, []int64, error) {
	if len(images) == 0 || len(images) > 4096 {
		return nil, nil, fmt.Errorf("invalid architecture count %d", len(images))
	}
	slices.SortFunc(images, func(a, b image) int {
		if c := cmp.Compare(a.cpu, b.cpu); c != 0 {
			return c
		}
		return cmp.Compare(a.sub, b.sub)
	})
	header := make([]byte, 8+20*len(images))
	binary.BigEndian.PutUint32(header, macho.MagicFat)
	binary.BigEndian.PutUint32(header[4:], uint32(len(images)))
	offsets := make([]int64, len(images))
	end := int64(len(header))
	for i, img := range images {
		if i > 0 && img.cpu == images[i-1].cpu && img.sub == images[i-1].sub {
			return nil, nil, fmt.Errorf("duplicate architecture CPU %#x subtype %#x", img.cpu, img.sub)
		}
		if img.align > 30 || img.size <= 0 {
			return nil, nil, fmt.Errorf("invalid architecture size or alignment")
		}
		alignment := int64(1) << img.align
		start := (end + alignment - 1) & -alignment
		if start > math.MaxUint32 || img.size > math.MaxUint32 || img.size > math.MaxUint32-start {
			return nil, nil, fmt.Errorf("universal binary exceeds 32-bit format limits")
		}
		offsets[i], end = start, start+img.size
		entry := header[8+20*i:]
		for j, value := range []uint32{img.cpu, img.sub, uint32(start), uint32(img.size), img.align} {
			binary.BigEndian.PutUint32(entry[j*4:], value)
		}
	}
	return header, offsets, nil
}
