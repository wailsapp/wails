// Package readonlyfs writes portable SquashFS v4 filesystems for application
// images. Files use zlib blocks; metadata stays uncompressed for predictable
// references. The implementation follows the published on-disk format.
package readonlyfs

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"math"
	"math/bits"
	"path"
	"sort"
	"strings"
	"time"
)

const metadataSize = 8192

type Entry struct {
	Path string
	Mode fs.FileMode
	Size int64
	Link string
	Open func() (io.ReadCloser, error)
}
type Options struct {
	BlockSize int
	ModTime   time.Time
}

type inode struct {
	Entry
	number          uint32
	reference       uint64
	start           uint32
	blocks          []uint32
	children        []*inode
	parent          *inode
	directoryOffset int
	directorySize   int
}

// Write appends a filesystem at output's current position. Every stored offset
// is relative to that position, so an executable runtime may precede the image.
// The caller owns cleanup/publication when writing fails. Entries and their
// content must remain stable until Write returns. Ownership is root:root.
func Write(ctx context.Context, output io.WriteSeeker, entries []Entry, options Options) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	blockSize := options.BlockSize
	if blockSize == 0 {
		blockSize = 128 << 10
	}
	if blockSize < 4096 || blockSize > 1<<20 || blockSize&(blockSize-1) != 0 {
		return fmt.Errorf("SquashFS block size must be a power of two from 4096 to 1048576")
	}
	timestamp := int64(0)
	if !options.ModTime.IsZero() {
		timestamp = options.ModTime.Unix()
	}
	if timestamp < 0 || timestamp > math.MaxUint32 {
		return fmt.Errorf("SquashFS timestamp is outside its unsigned 32-bit range")
	}
	nodes, err := prepareEntries(entries)
	if err != nil {
		return err
	}
	var blocks, metadataBytes int64
	for _, node := range nodes {
		metadataBytes += int64(60 + len(node.Link) + len(path.Base(node.Path)))
		if node.Mode.IsRegular() {
			blocks += (node.Size + int64(blockSize) - 1) / int64(blockSize)
		}
	}
	if metadataBytes+blocks*4 > 64<<20 {
		return fmt.Errorf("image metadata exceeds 64 MiB")
	}
	base, err := output.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	writer := imageWriter{output: output}
	if err := writer.write(make([]byte, 96)); err != nil {
		return err
	}
	if err := writer.writeFiles(ctx, nodes, blockSize); err != nil {
		return err
	}
	assignReferences(nodes)
	directories := encodeDirectories(nodes)
	inodes := encodeInodes(nodes, uint32(timestamp))
	inodeStart := writer.position
	if err := writer.metadata(inodes); err != nil {
		return err
	}
	directoryStart := writer.position
	if err := writer.metadata(directories); err != nil {
		return err
	}
	idBlock := writer.position
	if err := writer.metadata(make([]byte, 4)); err != nil {
		return err
	}
	idTable := writer.position
	var idPointer [8]byte
	binary.LittleEndian.PutUint64(idPointer[:], uint64(idBlock))
	if err := writer.write(idPointer[:]); err != nil {
		return err
	}
	used := writer.position
	super := make([]byte, 96)
	binary.LittleEndian.PutUint32(super, 0x73717368)
	binary.LittleEndian.PutUint32(super[4:], uint32(len(nodes)))
	binary.LittleEndian.PutUint32(super[8:], uint32(timestamp))
	binary.LittleEndian.PutUint32(super[12:], uint32(blockSize))
	binary.LittleEndian.PutUint16(super[20:], 1)
	binary.LittleEndian.PutUint16(super[22:], uint16(bits.TrailingZeros(uint(blockSize))))
	binary.LittleEndian.PutUint16(super[24:], 0x0001|0x0010|0x0200)
	binary.LittleEndian.PutUint16(super[26:], 1)
	binary.LittleEndian.PutUint16(super[28:], 4)
	binary.LittleEndian.PutUint64(super[32:], nodes[0].reference)
	for offset, value := range map[int]uint64{40: uint64(used), 48: uint64(idTable), 56: math.MaxUint64, 64: uint64(inodeStart), 72: uint64(directoryStart), 80: math.MaxUint64, 88: math.MaxUint64} {
		binary.LittleEndian.PutUint64(super[offset:], value)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := output.Seek(base, io.SeekStart); err != nil {
		return err
	}
	if err := writer.write(super); err != nil {
		return err
	}
	_, err = output.Seek(base+used, io.SeekStart)
	return err
}

func prepareEntries(entries []Entry) ([]*inode, error) {
	if len(entries) > 100000 {
		return nil, fmt.Errorf("image exceeds 100000 entries")
	}
	byPath := make(map[string]*inode, len(entries)+1)
	for _, entry := range entries {
		if !fs.ValidPath(entry.Path) || len(entry.Path) > 4096 || strings.ContainsAny(entry.Path, "\x00\\") {
			return nil, fmt.Errorf("invalid image path %q", entry.Path)
		}
		if byPath[entry.Path] != nil {
			return nil, fmt.Errorf("duplicate image path %q", entry.Path)
		}
		for _, part := range strings.Split(entry.Path, "/") {
			if len(part) > 256 {
				return nil, fmt.Errorf("image filename exceeds 256 bytes: %s", entry.Path)
			}
		}
		switch entry.Mode.Type() {
		case 0:
			if entry.Size < 0 || entry.Size > math.MaxUint32 || entry.Open == nil {
				return nil, fmt.Errorf("invalid regular file %s (requires reader and size below 4 GiB)", entry.Path)
			}
		case fs.ModeDir:
		case fs.ModeSymlink:
			if entry.Link == "" || len(entry.Link) > 65535 || strings.ContainsRune(entry.Link, 0) {
				return nil, fmt.Errorf("invalid symlink %s", entry.Path)
			}
		default:
			return nil, fmt.Errorf("unsupported image file type: %s", entry.Path)
		}
		byPath[entry.Path] = &inode{Entry: entry}
	}
	if byPath["."] == nil {
		byPath["."] = &inode{Entry: Entry{Path: ".", Mode: fs.ModeDir | 0755}}
	}
	if !byPath["."].Mode.IsDir() {
		return nil, fmt.Errorf("image root must be a directory")
	}
	for _, entry := range entries {
		for parent := path.Dir(entry.Path); parent != "."; parent = path.Dir(parent) {
			if byPath[parent] == nil {
				if len(byPath) >= 100000 {
					return nil, fmt.Errorf("image exceeds 100000 entries")
				}
				byPath[parent] = &inode{Entry: Entry{Path: parent, Mode: fs.ModeDir | 0755}}
			}
		}
	}
	if len(byPath) > 100000 {
		return nil, fmt.Errorf("image exceeds 100000 entries")
	}
	nodes := make([]*inode, 0, len(byPath))
	for _, node := range byPath {
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Path == "." || nodes[j].Path == "." {
			return nodes[i].Path == "."
		}
		return nodes[i].Path < nodes[j].Path
	})
	for i, node := range nodes {
		node.number = uint32(i + 1)
		if node.Path == "." {
			continue
		}
		parent := byPath[path.Dir(node.Path)]
		if !parent.Mode.IsDir() {
			return nil, fmt.Errorf("parent of %s is not a directory", node.Path)
		}
		node.parent = parent
		parent.children = append(parent.children, node)
	}
	return nodes, nil
}

type imageWriter struct {
	output   io.Writer
	position int64
}

func (w *imageWriter) write(data []byte) error {
	n, err := w.output.Write(data)
	w.position += int64(n)
	if err == nil && n != len(data) {
		err = io.ErrShortWrite
	}
	return err
}
func (w *imageWriter) metadata(data []byte) error {
	for len(data) > 0 {
		count := min(len(data), metadataSize)
		var prefix [2]byte
		binary.LittleEndian.PutUint16(prefix[:], uint16(count)|0x8000)
		if err := w.write(prefix[:]); err != nil {
			return err
		}
		if err := w.write(data[:count]); err != nil {
			return err
		}
		data = data[count:]
	}
	return nil
}
func (w *imageWriter) writeFiles(ctx context.Context, nodes []*inode, blockSize int) error {
	data := make([]byte, blockSize)
	var compressed bytes.Buffer
	compressor := zlib.NewWriter(&compressed)
	defer compressor.Close()
	for _, node := range nodes {
		if !node.Mode.IsRegular() {
			continue
		}
		if w.position > math.MaxUint32 {
			return fmt.Errorf("file data offset exceeds 4 GiB")
		}
		node.start = uint32(w.position)
		input, err := node.Open()
		if err != nil {
			return fmt.Errorf("open %s: %w", node.Path, err)
		}
		if input == nil {
			return fmt.Errorf("open %s returned a nil reader", node.Path)
		}
		err = w.writeFile(ctx, node, input, data, &compressed, compressor)
		closeErr := input.Close()
		if err != nil {
			return fmt.Errorf("write %s: %w", node.Path, err)
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}
func (w *imageWriter) writeFile(ctx context.Context, node *inode, input io.Reader, data []byte, compressed *bytes.Buffer, compressor *zlib.Writer) error {
	for remaining := node.Size; remaining > 0; {
		if err := ctx.Err(); err != nil {
			return err
		}
		size := min(int64(len(data)), remaining)
		if _, err := io.ReadFull(input, data[:size]); err != nil {
			return err
		}
		compressed.Reset()
		compressor.Reset(compressed)
		if _, err := compressor.Write(data[:size]); err != nil {
			return err
		}
		if err := compressor.Close(); err != nil {
			return err
		}
		block := compressed.Bytes()
		encodedSize := uint32(len(block))
		if len(block) >= int(size) {
			block = data[:size]
			encodedSize = uint32(size) | 1<<24
		}
		if err := w.write(block); err != nil {
			return err
		}
		node.blocks = append(node.blocks, encodedSize)
		remaining -= size
	}
	var extra [1]byte
	if n, err := io.ReadFull(input, extra[:]); n != 0 || err != io.EOF {
		if err != nil {
			return err
		}
		return fmt.Errorf("file grew beyond declared size")
	}
	return nil
}

func metadataReference(offset int) uint64 {
	return uint64(offset/metadataSize*(metadataSize+2))<<16 | uint64(offset%metadataSize)
}
func assignReferences(nodes []*inode) {
	offset := 0
	for _, node := range nodes {
		node.reference = metadataReference(offset)
		switch node.Mode.Type() {
		case fs.ModeDir:
			offset += 40
		case fs.ModeSymlink:
			offset += 24 + len(node.Link)
		default:
			offset += 32 + 4*len(node.blocks)
		}
	}
}

func encodeDirectories(nodes []*inode) []byte {
	var data []byte
	for _, node := range nodes {
		if !node.Mode.IsDir() {
			continue
		}
		node.directoryOffset = len(data)
		for begin := 0; begin < len(node.children); {
			first := node.children[begin]
			end := begin + 1
			for end < len(node.children) && end-begin < 256 {
				next := node.children[end]
				if next.reference>>16 != first.reference>>16 || int64(next.number)-int64(first.number) > math.MaxInt16 {
					break
				}
				end++
			}
			data = binary.LittleEndian.AppendUint32(data, uint32(end-begin-1))
			data = binary.LittleEndian.AppendUint32(data, uint32(first.reference>>16))
			data = binary.LittleEndian.AppendUint32(data, first.number)
			for _, child := range node.children[begin:end] {
				name := path.Base(child.Path)
				data = binary.LittleEndian.AppendUint16(data, uint16(child.reference))
				data = binary.LittleEndian.AppendUint16(data, uint16(child.number-first.number))
				data = binary.LittleEndian.AppendUint16(data, basicType(child))
				data = binary.LittleEndian.AppendUint16(data, uint16(len(name)-1))
				data = append(data, name...)
			}
			begin = end
		}
		node.directorySize = len(data) - node.directoryOffset + 3
	}
	return data
}
func basicType(node *inode) uint16 {
	if node.Mode.IsDir() {
		return 1
	}
	if node.Mode&fs.ModeSymlink != 0 {
		return 3
	}
	return 2
}
func encodeInodes(nodes []*inode, mtime uint32) []byte {
	var data []byte
	for _, node := range nodes {
		kind := basicType(node)
		if kind == 1 {
			kind = 8
		}
		mode := uint16(node.Mode.Perm())
		if node.Mode&fs.ModeSetuid != 0 {
			mode |= 04000
		}
		if node.Mode&fs.ModeSetgid != 0 {
			mode |= 02000
		}
		if node.Mode&fs.ModeSticky != 0 {
			mode |= 01000
		}
		data = binary.LittleEndian.AppendUint16(data, kind)
		data = binary.LittleEndian.AppendUint16(data, mode)
		data = binary.LittleEndian.AppendUint32(data, 0)
		data = binary.LittleEndian.AppendUint32(data, mtime)
		data = binary.LittleEndian.AppendUint32(data, node.number)
		switch kind {
		case 8:
			links := uint32(2)
			for _, child := range node.children {
				if child.Mode.IsDir() {
					links++
				}
			}
			ref := metadataReference(node.directoryOffset)
			parent := uint32(0)
			if node.parent != nil {
				parent = node.parent.number
			}
			data = binary.LittleEndian.AppendUint32(data, links)
			data = binary.LittleEndian.AppendUint32(data, uint32(node.directorySize))
			data = binary.LittleEndian.AppendUint32(data, uint32(ref>>16))
			data = binary.LittleEndian.AppendUint32(data, parent)
			data = binary.LittleEndian.AppendUint16(data, 0)
			data = binary.LittleEndian.AppendUint16(data, uint16(ref))
			data = binary.LittleEndian.AppendUint32(data, math.MaxUint32)
		case 3:
			data = binary.LittleEndian.AppendUint32(data, 1)
			data = binary.LittleEndian.AppendUint32(data, uint32(len(node.Link)))
			data = append(data, node.Link...)
		default:
			data = binary.LittleEndian.AppendUint32(data, node.start)
			data = binary.LittleEndian.AppendUint32(data, math.MaxUint32)
			data = binary.LittleEndian.AppendUint32(data, 0)
			data = binary.LittleEndian.AppendUint32(data, uint32(node.Size))
			for _, size := range node.blocks {
				data = binary.LittleEndian.AppendUint32(data, size)
			}
		}
	}
	return data
}
