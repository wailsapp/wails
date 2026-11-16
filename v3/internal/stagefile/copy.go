// Package stagefile copies regular files into caller-owned staging directories.
package stagefile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Copy creates destination exclusively and preserves source permissions when
// mode is zero. A reusable buffer avoids per-file allocations. The caller owns
// cleanup of staged output if copying fails or the context is cancelled.
func Copy(ctx context.Context, source, destination string, mode os.FileMode, buffer []byte) (err error) {
	if err = ctx.Err(); err != nil {
		return err
	}
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("source is not a regular file: %s", source)
	}
	if err = os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, output.Close()) }()
	_, err = io.CopyBuffer(contextWriter{ctx, output}, struct{ io.Reader }{input}, buffer)
	if err == nil {
		if mode == 0 {
			mode = info.Mode().Perm()
		}
		err = output.Chmod(mode)
	}
	return err
}

type contextWriter struct {
	ctx    context.Context
	writer io.Writer
}

func (w contextWriter) Write(data []byte) (int, error) {
	if err := w.ctx.Err(); err != nil {
		return 0, err
	}
	return w.writer.Write(data)
}
