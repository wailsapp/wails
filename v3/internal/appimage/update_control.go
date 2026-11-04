package appimage

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/md4"
)

// writeUpdateControl implements the documented zsync control format, with full
// MD4 block hashes for compatibility with existing AppImage update clients.
// Format: https://github.com/cph6/zsync/blob/master/SPEC.md
func writeUpdateControl(ctx context.Context, image string) (err error) {
	input, err := os.Open(image)
	if err != nil {
		return err
	}
	defer input.Close()
	name := filepath.Base(image)
	if strings.ContainsAny(name, "\r\n") {
		return fmt.Errorf("invalid update filename %q", name)
	}
	hash := sha1.New() // Legacy clients require the format's whole-file SHA-1 field.
	size, err := io.Copy(contextWriter{ctx, hash}, input)
	if err != nil {
		return err
	}
	if _, err = input.Seek(0, io.SeekStart); err != nil {
		return err
	}
	output, err := os.OpenFile(image+".zsync", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := output.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			os.Remove(output.Name())
		}
	}()
	writer := bufio.NewWriter(output)
	const blockSize = 2048
	_, err = fmt.Fprintf(writer, "zsync: 0.6.2\nFilename: %s\nBlocksize: %d\nLength: %d\nHash-Lengths: 1,4,16\nURL: %s\nSHA-1: %x\n\n", name, blockSize, size, url.PathEscape(name), hash.Sum(nil))
	if err != nil {
		return err
	}
	block := make([]byte, blockSize)
	reader := bufio.NewReaderSize(input, 128<<10)
	strong := md4.New()
	var strongBytes [md4.Size]byte
	var weak [4]byte
	for {
		if err = ctx.Err(); err != nil {
			return err
		}
		clear(block)
		count, readErr := io.ReadFull(reader, block)
		if count == 0 && readErr == io.EOF {
			break
		}
		if readErr != nil && readErr != io.ErrUnexpectedEOF {
			return readErr
		}
		var a, b uint16
		for i, value := range block {
			a += uint16(value)
			b += uint16((blockSize - i) * int(value))
		}
		binary.BigEndian.PutUint16(weak[:2], a)
		binary.BigEndian.PutUint16(weak[2:], b)
		strong.Reset()
		strong.Write(block)
		if _, err = writer.Write(weak[:]); err != nil {
			return err
		}
		if _, err = writer.Write(strong.Sum(strongBytes[:0])); err != nil {
			return err
		}
		if readErr == io.ErrUnexpectedEOF {
			break
		}
	}
	return writer.Flush()
}
