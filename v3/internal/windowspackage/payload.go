package windowspackage

import (
	"archive/zip"
	"compress/flate"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"net/url"
	"os"
	"strings"
)

func addPayload(ctx context.Context, archive *zip.Writer, source, name, stage string, level int, buffer []byte) (fileHash, error) {
	var result fileHash
	input, err := os.Open(source)
	if err != nil {
		return result, err
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return result, err
	}
	if !info.Mode().IsRegular() {
		return result, fmt.Errorf("MSIX input is not a regular file: %s", source)
	}
	compressed, err := os.CreateTemp(stage, "payload-*")
	if err != nil {
		return result, err
	}
	defer os.Remove(compressed.Name())
	defer compressed.Close()
	header, result, err := encodePayload(ctx, compressed, input, name, level, info.Size() != 0, buffer)
	if err != nil {
		return result, err
	}
	after, err := input.Stat()
	if err != nil {
		return result, err
	}
	if uint64(info.Size()) != result.Size || info.Size() != after.Size() || !info.ModTime().Equal(after.ModTime()) {
		return result, fmt.Errorf("MSIX input changed while packaging: %s", source)
	}
	if _, err = compressed.Seek(0, io.SeekStart); err != nil {
		return result, err
	}
	writer, err := archive.CreateRaw(header)
	if err != nil {
		return result, err
	}
	_, err = io.CopyBuffer(writer, contextReader{ctx, compressed}, buffer)
	return result, err
}

func encodePayload(ctx context.Context, output io.Writer, input io.Reader, name string, level int, compress bool, buffer []byte) (*zip.FileHeader, fileHash, error) {
	result := fileHash{Name: strings.ReplaceAll(name, "/", "\\")}
	part := (&url.URL{Path: name}).EscapedPath()
	header := &zip.FileHeader{Name: part, Flags: 0x800, ReaderVersion: 20, CreatorVersion: 20, ModifiedDate: 33}
	checksum := crc32.NewIEEE()
	counted := &countWriter{Writer: output}
	var deflater *flate.Writer
	if compress {
		if level == 0 {
			level = 4
		}
		var err error
		deflater, err = flate.NewWriter(counted, level)
		if err != nil {
			return nil, result, err
		}
		header.Method = zip.Deflate
	}
	for {
		n, err := io.ReadFull(contextReader{ctx, input}, buffer)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, result, err
		}
		if n == 0 {
			break
		}
		data := buffer[:n]
		_, _ = checksum.Write(data)
		hash := sha256.Sum256(data)
		block := blockHash{Hash: base64.StdEncoding.EncodeToString(hash[:])}
		before := counted.bytes
		if deflater != nil {
			// Each flushed block starts with a fresh dictionary, allowing block
			// extraction without bytes from a preceding update/download block.
			deflater.Reset(counted)
			if _, err = deflater.Write(data); err == nil {
				err = deflater.Flush()
			}
			block.Size = counted.bytes - before
		} else {
			_, err = counted.Write(data)
		}
		if err != nil {
			return nil, result, err
		}
		result.Blocks = append(result.Blocks, block)
		result.Size += uint64(n)
	}
	if deflater != nil {
		before := counted.bytes
		if err := deflater.Close(); err != nil {
			return nil, result, err
		}
		if len(result.Blocks) > 0 {
			result.Blocks[len(result.Blocks)-1].Size += counted.bytes - before
		}
	}
	header.CRC32 = checksum.Sum32()
	header.UncompressedSize64 = result.Size
	header.CompressedSize64 = counted.bytes
	if result.Size >= 1<<32 || counted.bytes >= 1<<32 {
		header.ReaderVersion = 45
		header.Extra = make([]byte, 20)
		binary.LittleEndian.PutUint16(header.Extra[0:2], 1)
		binary.LittleEndian.PutUint16(header.Extra[2:4], 16)
		binary.LittleEndian.PutUint64(header.Extra[4:12], result.Size)
		binary.LittleEndian.PutUint64(header.Extra[12:20], counted.bytes)
	}
	result.Header = 30 + len(part) + len(header.Extra)
	return header, result, nil
}

type countWriter struct {
	io.Writer
	bytes uint64
}

func (w *countWriter) Write(data []byte) (int, error) {
	n, err := w.Writer.Write(data)
	w.bytes += uint64(n)
	return n, err
}
