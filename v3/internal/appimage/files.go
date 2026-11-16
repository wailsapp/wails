package appimage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/wailsapp/wails/v3/internal/stagefile"
)

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

func publish(ctx context.Context, source, output string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := os.Rename(source, output)
	if !errors.Is(err, syscall.EXDEV) {
		return err
	}
	stage, err := os.MkdirTemp(filepath.Dir(output), ".wails-appimage-output-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	filename := filepath.Join(stage, filepath.Base(output))
	if err = stagefile.Copy(ctx, source, filename, 0755, nil); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return os.Rename(filename, output)
}
