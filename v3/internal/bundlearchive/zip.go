// Package bundlearchive writes portable ZIP archives of application bundles.
package bundlearchive

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Create writes a new archive and removes incomplete output on failure.
func Create(ctx context.Context, output, root string) (err error) {
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, file.Close())
		if err != nil {
			_ = os.Remove(output)
		}
	}()
	return Write(ctx, file, root)
}

// Write includes the root directory and preserves Unix modes and relative
// symlinks without following them. The caller owns the output and its cleanup.
func Write(ctx context.Context, output io.Writer, root string) (err error) {
	root = filepath.Clean(root)
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() || filepath.Base(root) == "." || filepath.Base(root) == string(filepath.Separator) {
		return fmt.Errorf("bundle archive requires a named directory: %s", root)
	}
	archive := zip.NewWriter(output)
	defer func() { err = errors.Join(err, archive.Close()) }()
	buffer := make([]byte, 64<<10)
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		name, err := filepath.Rel(filepath.Dir(root), path)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(name)
		header.Method = zip.Deflate
		var link string
		switch {
		case info.IsDir():
			header.Name += "/"
			header.Method = zip.Store
		case info.Mode()&os.ModeSymlink != 0:
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
			resolved, err := filepath.Rel(root, filepath.Join(filepath.Dir(path), link))
			if err != nil || filepath.IsAbs(link) || !filepath.IsLocal(resolved) {
				return fmt.Errorf("bundle symlink escapes archive: %s", name)
			}
			header.Method = zip.Store
		case !info.Mode().IsRegular():
			return fmt.Errorf("unsupported bundle entry: %s (%s)", name, info.Mode())
		}
		writer, err := archive.CreateHeader(header)
		if err != nil || info.IsDir() {
			return err
		}
		if link != "" {
			_, err = io.Copy(writer, strings.NewReader(link))
			return err
		}
		input, err := os.Open(path)
		if err != nil {
			return err
		}
		_, copyErr := io.CopyBuffer(writer, contextReader{ctx, input}, buffer)
		return errors.Join(copyErr, input.Close())
	})
}

type contextReader struct {
	ctx context.Context
	io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.Reader.Read(p)
}
