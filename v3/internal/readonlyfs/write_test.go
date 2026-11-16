package readonlyfs

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fileEntry(name string, data []byte, mode fs.FileMode) Entry {
	return Entry{Path: name, Mode: mode, Size: int64(len(data)), Open: func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(data)), nil }}
}
func buildImage(t *testing.T, entries []Entry, prefix int) string {
	t.Helper()
	name := filepath.Join(t.TempDir(), "image.sqfs")
	out, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if prefix > 0 {
		if _, err := out.Write(bytes.Repeat([]byte{0xa5}, prefix)); err != nil {
			t.Fatal(err)
		}
	}
	if err := Write(t.Context(), out, entries, Options{BlockSize: 4096, ModTime: time.Unix(1234567890, 0)}); err != nil {
		out.Close()
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
	return name
}
func TestIndependentExtraction(t *testing.T) {
	tool, err := exec.LookPath("unsquashfs")
	if err != nil {
		t.Skip("independent unsquashfs reader unavailable")
	}
	random := make([]byte, 12001)
	rng := rand.New(rand.NewPCG(7, 11))
	for i := range random {
		random[i] = byte(rng.Uint32())
	}
	entries := []Entry{fileEntry("!first", []byte("first"), 0640), fileEntry("usr/bin/app", bytes.Repeat([]byte("compressed data\n"), 3000), 0755), fileEntry("usr/lib/random", random, 0644), fileEntry("empty", nil, 0600), {Path: "emptydir", Mode: fs.ModeDir | 0700}, {Path: "link", Mode: fs.ModeSymlink | 0777, Link: "usr/bin/app"}}
	// Force inode and directory metadata to span several 8KiB blocks, including
	// directory listings that cannot fit the basic 16-bit directory inode.
	for i := range 3200 {
		entries = append(entries, fileEntry(fmt.Sprintf("many/entry-%04d-padding", i), []byte{byte(i)}, 0644))
	}
	for _, prefix := range []int{0, 12345} {
		t.Run(fmt.Sprint(prefix), func(t *testing.T) {
			image := buildImage(t, entries, prefix)
			dest := filepath.Join(t.TempDir(), "extracted")
			command := exec.CommandContext(t.Context(), tool, "-no-progress", "-processors", "2", "-offset", fmt.Sprint(prefix), "-d", dest, image)
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("independent reader: %v\n%s", err, output)
			}
			for _, entry := range entries {
				target := filepath.Join(dest, filepath.FromSlash(entry.Path))
				info, err := os.Lstat(target)
				if err != nil {
					t.Fatal(err)
				}
				if info.Mode().Type() != entry.Mode.Type() {
					t.Fatalf("type %s: %v", entry.Path, info.Mode())
				}
				if info.Mode().Perm() != entry.Mode.Perm() {
					t.Fatalf("mode %s: %v", entry.Path, info.Mode())
				}
				if entry.Mode.IsRegular() {
					r, _ := entry.Open()
					want, _ := io.ReadAll(r)
					r.Close()
					got, err := os.ReadFile(target)
					if err != nil || !bytes.Equal(got, want) {
						t.Fatalf("payload %s differs: %v", entry.Path, err)
					}
				}
				if entry.Mode&fs.ModeSymlink != 0 {
					link, err := os.Readlink(target)
					if err != nil || link != entry.Link {
						t.Fatalf("symlink mismatch %q: %v", link, err)
					}
				}
			}
		})
	}
}
func TestEmptyImageAndDeterminism(t *testing.T) {
	first := buildImage(t, nil, 0)
	data, err := os.ReadFile(first)
	if err != nil {
		t.Fatal(err)
	}
	if binary.LittleEndian.Uint32(data) != 0x73717368 || binary.LittleEndian.Uint64(data[40:]) != uint64(len(data)) {
		t.Fatal("invalid superblock")
	}
	entries := []Entry{fileEntry("b", []byte("B"), 0644), fileEntry("a", []byte("A"), 0755)}
	one := buildImage(t, entries, 0)
	entries[0], entries[1] = entries[1], entries[0]
	two := buildImage(t, entries, 0)
	a, _ := os.ReadFile(one)
	b, _ := os.ReadFile(two)
	if !bytes.Equal(a, b) {
		t.Fatal("input ordering changed filesystem bytes")
	}
}
func TestInvalidInputs(t *testing.T) {
	valid := fileEntry("file", []byte("hello"), 0644)
	cases := [][]Entry{{fileEntry("../escape", nil, 0644)}, {fileEntry("/absolute", nil, 0644)}, {fileEntry("a\\b", nil, 0644)}, {fileEntry(strings.Repeat("a", 257), nil, 0644)}, {valid, valid}, {fileEntry(".", nil, 0644)}, {valid, fileEntry("file/child", nil, 0644)}, {{Path: "fifo", Mode: fs.ModeNamedPipe}}, {{Path: "link", Mode: fs.ModeSymlink}}, {{Path: "file", Size: 5}}}
	for i, entries := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			out, err := os.CreateTemp(t.TempDir(), "image")
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()
			if err := Write(t.Context(), out, entries, Options{}); err == nil {
				t.Fatal("invalid entry accepted")
			}
		})
	}
}
func TestSourceFailuresAndCancellation(t *testing.T) {
	for _, delta := range []int64{-1, 1} {
		entry := fileEntry("file", []byte("hello"), 0644)
		entry.Size += delta
		out, err := os.CreateTemp(t.TempDir(), "image")
		if err != nil {
			t.Fatal(err)
		}
		err = Write(t.Context(), out, []Entry{entry}, Options{})
		out.Close()
		if err == nil {
			t.Fatal("changed file size accepted")
		}
	}
	out, err := os.CreateTemp(t.TempDir(), "image")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := Write(ctx, out, nil, Options{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
	for _, size := range []int{1, 8193, 2 << 20} {
		if err := Write(t.Context(), out, nil, Options{BlockSize: size}); err == nil {
			t.Fatal("invalid block size accepted")
		}
	}
}
func BenchmarkWrite(b *testing.B) {
	payload := make([]byte, 4<<20)
	rng := rand.New(rand.NewPCG(17, 29))
	for i := range payload {
		if i%4 == 0 {
			payload[i] = byte(rng.Uint32())
		} else {
			payload[i] = byte(i % 31)
		}
	}
	entry := fileEntry("usr/bin/app", payload, 0755)
	out, err := os.Create(filepath.Join(b.TempDir(), "image.sqfs"))
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()
	b.SetBytes(int64(len(payload)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := out.Seek(0, io.SeekStart); err != nil {
			b.Fatal(err)
		}
		if err := Write(b.Context(), out, []Entry{entry}, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

type brokenOutput struct {
	seekCalls    int
	seekFailure  int
	writeFailure error
	short        bool
}

func (w *brokenOutput) Seek(offset int64, whence int) (int64, error) {
	w.seekCalls++
	if w.seekCalls == w.seekFailure {
		return 0, errTestOutput
	}
	return offset, nil
}
func (w *brokenOutput) Write(data []byte) (int, error) {
	if w.writeFailure != nil {
		return 0, w.writeFailure
	}
	if w.short {
		return len(data) - 1, nil
	}
	return len(data), nil
}

var errTestOutput = errors.New("test output failure")

func TestOutputErrors(t *testing.T) {
	for _, output := range []*brokenOutput{{seekFailure: 1}, {seekFailure: 2}, {seekFailure: 3}, {writeFailure: errTestOutput}} {
		if err := Write(t.Context(), output, nil, Options{}); !errors.Is(err, errTestOutput) {
			t.Fatalf("output failure hidden: %v", err)
		}
	}
	if err := Write(t.Context(), &brokenOutput{short: true}, nil, Options{}); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short write: %v", err)
	}
}
func TestMetadataBound(t *testing.T) {
	target := strings.Repeat("x", 65535)
	entries := make([]Entry, 1025)
	for i := range entries {
		entries[i] = Entry{Path: fmt.Sprintf("link-%d", i), Mode: fs.ModeSymlink | 0777, Link: target}
	}
	output := &brokenOutput{}
	if err := Write(t.Context(), output, entries, Options{}); err == nil || output.seekCalls != 0 {
		t.Fatalf("metadata limit did not fail before output: %v", err)
	}
}
