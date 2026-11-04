package windowsign

import (
	"context"
	"crypto/sha256"
	"debug/pe"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type peLayout struct {
	checksum, directory, headers, end int64
	sections                          []*pe.Section
}

type contextReader struct {
	ctx context.Context
	io.Reader
}

func (r contextReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.Reader.Read(data)
}

// SignPE replaces an existing terminal certificate table after validating the
// image layout. It preserves the source until the complete output is ready.
func SignPE(ctx context.Context, input, output string, identity Identity, options Options) error {
	source, info, err := snapshotInput(ctx, input, output)
	if err != nil {
		return err
	}
	defer os.Remove(source.Name())
	defer source.Close()
	layout, err := inspectPE(source, info.Size())
	if err != nil {
		return err
	}
	digest, err := hashPE(ctx, source, layout)
	if err != nil {
		return err
	}
	format, err := peFormatData()
	if err != nil {
		return err
	}
	signature, err := signedContent(ctx, identity, oidPEImage, format, digest, options)
	if err != nil {
		return err
	}
	return publishSigned(ctx, output, info.Mode().Perm(), func(stage *os.File) error {
		if _, err := io.Copy(stage, contextReader{ctx, io.NewSectionReader(source, 0, layout.end)}); err != nil {
			return err
		}
		padding := (8 - layout.end%8) % 8
		if _, err := stage.Write(make([]byte, padding)); err != nil {
			return err
		}
		tableOffset := layout.end + padding
		tableSize := (8 + len(signature) + 7) &^ 7
		if tableOffset+int64(tableSize) >= 1<<32 {
			return errors.New("signed PE exceeds 4 GiB")
		}
		certificate := make([]byte, tableSize)
		binary.LittleEndian.PutUint32(certificate, uint32(tableSize))
		binary.LittleEndian.PutUint16(certificate[4:], 0x200)
		binary.LittleEndian.PutUint16(certificate[6:], 2)
		copy(certificate[8:], signature)
		if _, err := stage.Write(certificate); err != nil {
			return err
		}
		directory := make([]byte, 8)
		binary.LittleEndian.PutUint32(directory, uint32(tableOffset))
		binary.LittleEndian.PutUint32(directory[4:], uint32(tableSize))
		if _, err := stage.WriteAt(directory, layout.directory); err != nil {
			return err
		}
		if _, err := stage.WriteAt(make([]byte, 4), layout.checksum); err != nil {
			return err
		}
		checksum, err := imageChecksum(ctx, stage, tableOffset+int64(tableSize))
		if err != nil {
			return err
		}
		data := make([]byte, 4)
		binary.LittleEndian.PutUint32(data, checksum)
		_, err = stage.WriteAt(data, layout.checksum)
		return err
	})
}

func inspectPE(file *os.File, size int64) (peLayout, error) {
	var result peLayout
	parsed, err := pe.NewFile(file)
	if err != nil {
		return result, fmt.Errorf("parse PE: %w", err)
	}
	dos := make([]byte, 64)
	if _, err := file.ReadAt(dos, 0); err != nil {
		return result, err
	}
	optional := int64(binary.LittleEndian.Uint32(dos[0x3c:])) + 24
	var table pe.DataDirectory
	var number uint32
	switch header := parsed.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		result.headers = int64(header.SizeOfHeaders)
		result.directory = optional + 96 + 32
		number = header.NumberOfRvaAndSizes
		table = header.DataDirectory[4]
	case *pe.OptionalHeader64:
		result.headers = int64(header.SizeOfHeaders)
		result.directory = optional + 112 + 32
		number = header.NumberOfRvaAndSizes
		table = header.DataDirectory[4]
	default:
		return result, errors.New("PE image requires PE32 or PE32+ optional header")
	}
	result.checksum = optional + 64
	result.end = size
	if number < 5 || result.directory+8 > optional+int64(parsed.SizeOfOptionalHeader) || result.headers < result.directory+8 || result.headers > size {
		return result, errors.New("invalid PE header bounds")
	}
	if table.VirtualAddress != 0 || table.Size != 0 {
		result.end = int64(table.VirtualAddress)
		if result.end%8 != 0 || table.Size < 8 || result.end < int64(result.headers) || result.end+int64(table.Size) != size {
			return result, errors.New("existing PE certificate table must be aligned at the end of the file")
		}
		if err := validateCertificateTable(file, result.end, int64(table.Size)); err != nil {
			return result, err
		}
	}
	result.sections = append([]*pe.Section(nil), parsed.Sections...)
	sort.Slice(result.sections, func(i, j int) bool { return result.sections[i].Offset < result.sections[j].Offset })
	end := result.headers
	for _, section := range result.sections {
		if section.Size == 0 {
			continue
		}
		start := int64(section.Offset)
		finish := start + int64(section.Size)
		if start < end || finish > result.end {
			return result, errors.New("overlapping or out-of-bounds PE section")
		}
		end = finish
	}
	return result, nil
}

func validateCertificateTable(file *os.File, offset, size int64) error {
	header := make([]byte, 8)
	for cursor := int64(0); cursor < size; {
		if size-cursor < 8 {
			return errors.New("truncated PE certificate table")
		}
		if _, err := file.ReadAt(header, offset+cursor); err != nil {
			return err
		}
		length := int64(binary.LittleEndian.Uint32(header))
		aligned := (length + 7) &^ 7
		if length < 8 || aligned > size-cursor {
			return errors.New("invalid PE certificate entry length")
		}
		cursor += aligned
	}
	return nil
}

func hashPE(ctx context.Context, file *os.File, layout peLayout) ([]byte, error) {
	hash := sha256.New()
	add := func(start, end int64) error {
		if start < 0 || end < start || end > layout.end {
			return errors.New("invalid PE hash range")
		}
		_, err := io.Copy(hash, contextReader{ctx, io.NewSectionReader(file, start, end-start)})
		return err
	}
	for _, span := range [][2]int64{{0, layout.checksum}, {layout.checksum + 4, layout.directory}, {layout.directory + 8, layout.headers}} {
		if err := add(span[0], span[1]); err != nil {
			return nil, err
		}
	}
	sum := layout.headers
	for _, section := range layout.sections {
		if section.Size == 0 {
			continue
		}
		if err := add(int64(section.Offset), int64(section.Offset)+int64(section.Size)); err != nil {
			return nil, err
		}
		sum += int64(section.Size)
	}
	if err := add(sum, layout.end); err != nil {
		return nil, err
	}
	_, _ = hash.Write(make([]byte, (8-layout.end%8)%8))
	return hash.Sum(nil), nil
}

func imageChecksum(ctx context.Context, file *os.File, size int64) (uint32, error) {
	reader := contextReader{ctx, io.NewSectionReader(file, 0, size)}
	buffer := make([]byte, 128*1024)
	var sum uint64
	for {
		n, err := io.ReadFull(reader, buffer)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return 0, err
		}
		for i := 0; i+1 < n; i += 2 {
			sum += uint64(binary.LittleEndian.Uint16(buffer[i:]))
		}
		if n%2 != 0 {
			sum += uint64(buffer[n-1])
		}
		if err != nil {
			break
		}
	}
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return uint32(sum) + uint32(size), nil
}

func publishSigned(ctx context.Context, output string, mode os.FileMode, write func(*os.File) error) (err error) {
	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(output), ".wails-sign-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if err = file.Chmod(mode); err != nil {
		return err
	}
	if err = write(file); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	if err = file.Sync(); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), output)
}

// snapshotInput closes the caller's input before publication, permitting
// in-place replacement on Windows and keeping network signing delays away from
// mutable source bytes.
func snapshotInput(ctx context.Context, input, output string) (*os.File, os.FileInfo, error) {
	original, err := os.Open(input)
	if err != nil {
		return nil, nil, err
	}
	defer original.Close()
	info, err := original.Stat()
	if err != nil {
		return nil, nil, err
	}
	if !info.Mode().IsRegular() || info.Size() >= 1<<32 {
		return nil, nil, errors.New("portable Windows signing requires a regular file smaller than 4 GiB")
	}
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return nil, nil, err
	}
	snapshot, err := os.CreateTemp(filepath.Dir(output), ".wake-sign-input-*")
	if err != nil {
		return nil, nil, err
	}
	success := false
	defer func() {
		if !success {
			snapshot.Close()
			os.Remove(snapshot.Name())
		}
	}()
	size, err := io.Copy(snapshot, contextReader{ctx, original})
	if err != nil {
		return nil, nil, err
	}
	current, err := original.Stat()
	if err != nil {
		return nil, nil, err
	}
	if size != info.Size() || current.Size() != info.Size() || !current.ModTime().Equal(info.ModTime()) {
		return nil, nil, errors.New("signing input changed while it was being copied")
	}
	success = true
	return snapshot, info, nil
}
