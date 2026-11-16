// Package applesign writes Apple code signatures without platform executables.
// Signature formats are implementation details; interoperability is tested with
// Apple's codesign verifier, independently of this package's parser.
package applesign

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/universalbinary"
)

type Options struct {
	Identifier      string
	TeamID          string
	Entitlements    []byte
	HardenedRuntime bool
	// SignCMS returns a detached CMS signature over the complete CodeDirectory.
	// A nil callback selects an ad-hoc signature. Key loading, certificate policy
	// and timestamp validation belong to the shared signing-identity layer.
	SignCMS                         func(context.Context, []byte) ([]byte, error)
	requirement                     func(string) ([]byte, error)
	info, resources, representation []byte
}

type CodeIdentity struct {
	CPU        uint32
	Identifier string
	HashType   uint8
	CDHash     [20]byte
}

// SignMachO writes a new signed file. Input remains untouched, output must not
// exist, and partial output is removed on failure. Both native Wails desktop
// architectures are supported, independently signed in a universal executable.
func SignMachO(ctx context.Context, input, output string, options Options) (identities []CodeIdentity, err error) {
	if err = validateOptions(options); err != nil {
		return nil, err
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, errors.New("Mach-O input must be a regular file")
	}
	parts, err := readSlices(f, st.Size())
	if err != nil {
		return nil, err
	}
	// Reserve output before doing work, refusing aliases and pre-existing files.
	out, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, out.Close())
		if err != nil {
			_ = os.Remove(output)
		}
	}()
	if len(parts) == 1 {
		id, e := signThin(ctx, io.NewSectionReader(f, parts[0].offset, parts[0].size), parts[0].size, out, options)
		if e != nil {
			return nil, e
		}
		identities = append(identities, id)
	} else {
		work, e := os.MkdirTemp(filepath.Dir(output), ".wails-apple-slices-*")
		if e != nil {
			return nil, e
		}
		defer os.RemoveAll(work)
		paths := make([]string, 0, len(parts))
		for i, part := range parts {
			path := filepath.Join(work, fmt.Sprintf("slice-%d", i))
			file, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
			if e != nil {
				return nil, e
			}
			id, e := signThin(ctx, io.NewSectionReader(f, part.offset, part.size), part.size, file, options)
			e = errors.Join(e, file.Close())
			if e != nil {
				return nil, e
			}
			identities = append(identities, id)
			paths = append(paths, path)
		}
		combined := filepath.Join(work, "universal")
		if e = universalbinary.Create(combined, paths); e != nil {
			return nil, e
		}
		file, e := os.Open(combined)
		if e != nil {
			return nil, e
		}
		_, e = io.Copy(out, &contextReader{ctx, file})
		e = errors.Join(e, file.Close())
		if e != nil {
			return nil, e
		}
	}
	return identities, out.Chmod(st.Mode().Perm())
}

func validateOptions(o Options) error {
	if o.Identifier == "" || len(o.Identifier) > 1024 || strings.ContainsRune(o.Identifier, 0) {
		return errors.New("signing requires an identifier of at most 1024 bytes without NUL")
	}
	if len(o.TeamID) > 128 || strings.ContainsRune(o.TeamID, 0) {
		return errors.New("invalid signing team identifier")
	}
	if len(o.Entitlements) > 1<<20 {
		return errors.New("entitlements exceed 1 MiB")
	}
	return nil
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

type slice struct{ offset, size int64 }

func readSlices(r io.ReaderAt, size int64) ([]slice, error) {
	var head [8]byte
	if _, err := r.ReadAt(head[:], 0); err != nil {
		return nil, err
	}
	magic := binary.BigEndian.Uint32(head[:])
	if magic != 0xcafebabe && magic != 0xcafebabf {
		return []slice{{0, size}}, nil
	}
	count := int(binary.BigEndian.Uint32(head[4:]))
	width := 20
	if magic == 0xcafebabf {
		width = 32
	}
	end := int64(8 + count*width)
	if count == 0 || count > 32 || end > size {
		return nil, errors.New("invalid universal architecture table")
	}
	result := make([]slice, 0, count)
	for i := range count {
		var data [32]byte
		if _, err := r.ReadAt(data[:width], int64(8+i*width)); err != nil {
			return nil, err
		}
		offset, length := uint64(binary.BigEndian.Uint32(data[8:])), uint64(binary.BigEndian.Uint32(data[12:]))
		align := binary.BigEndian.Uint32(data[16:])
		if width == 32 {
			offset, length = binary.BigEndian.Uint64(data[8:]), binary.BigEndian.Uint64(data[16:])
			align = binary.BigEndian.Uint32(data[24:])
		}
		if offset < uint64(end) || offset > uint64(size) || length < 32 || length > uint64(size)-offset || align > 30 || offset%(uint64(1)<<align) != 0 {
			return nil, errors.New("universal slice outside file")
		}
		var thin [12]byte
		if _, err := r.ReadAt(thin[:], int64(offset)); err != nil {
			return nil, err
		}
		if binary.LittleEndian.Uint32(thin[:]) != 0xfeedfacf || binary.LittleEndian.Uint32(thin[4:]) != binary.BigEndian.Uint32(data[:]) || binary.LittleEndian.Uint32(thin[8:]) != binary.BigEndian.Uint32(data[4:]) {
			return nil, errors.New("universal architecture disagrees with Mach-O slice")
		}
		part := slice{int64(offset), int64(length)}
		for _, other := range result {
			if part.offset < other.offset+other.size && other.offset < part.offset+part.size {
				return nil, errors.New("overlapping universal slices")
			}
		}
		result = append(result, part)
	}
	return result, nil
}

func identity(cpu uint32, identifier string, directory []byte) CodeIdentity {
	h := sha256.Sum256(directory)
	id := CodeIdentity{CPU: cpu, Identifier: identifier, HashType: 2}
	copy(id.CDHash[:], h[:20])
	return id
}
