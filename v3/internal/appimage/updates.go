package appimage

import (
	"bufio"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/internal/stagefile"
)

func (options Options) resultName() string {
	if options.UpdateInfo != "" {
		return filepath.Base(options.Output)
	}
	return "result.AppImage"
}

func (options Options) publishResults(ctx context.Context, image string) error {
	if options.UpdateInfo == "" {
		return publish(ctx, image, options.Output)
	}
	if err := verifyZsync(ctx, image); err != nil {
		return fmt.Errorf("AppImage update control: %w", err)
	}
	// Prepare both files beside the destination before publishing either. Renames
	// are atomic per file; readers should fetch the control file after the image.
	stage, err := os.MkdirTemp(filepath.Dir(options.Output), ".wails-appimage-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	for _, target := range options.outputPaths() {
		suffix := strings.TrimPrefix(target, options.Output)
		info, err := os.Lstat(target)
		if err == nil && !info.Mode().IsRegular() {
			return fmt.Errorf("AppImage output must be a regular file: %s", target)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := stagefile.Copy(ctx, image+suffix, filepath.Join(stage, filepath.Base(target)), 0, nil); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, suffix := range []string{"", ".zsync"} {
		if err := os.Rename(filepath.Join(stage, filepath.Base(options.Output)+suffix), options.Output+suffix); err != nil {
			return err
		}
	}
	return nil
}

func verifyZsync(ctx context.Context, image string) error {
	control, err := os.Open(image + ".zsync")
	if err != nil {
		return fmt.Errorf("linuxdeploy must emit the .zsync file (check its zsyncmake support): %w", err)
	}
	defer control.Close()
	reader := bufio.NewReader(io.LimitReader(control, 64*1024))
	fields := map[string]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("incomplete zsync header: %w", err)
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok || fields[key] != "" {
			return fmt.Errorf("invalid zsync header field %q", key)
		}
		fields[key] = strings.TrimSpace(value)
	}
	if fields["zsync"] == "" || fields["Filename"] != filepath.Base(image) {
		return fmt.Errorf("zsync header must describe %s", filepath.Base(image))
	}
	input, err := os.Open(image)
	if err != nil {
		return err
	}
	defer input.Close()
	hash := sha1.New() // The zsync format specifies SHA-1 for content integrity.
	size, err := io.Copy(contextWriter{ctx, hash}, input)
	if err != nil {
		return err
	}
	if fields["Length"] != strconv.FormatInt(size, 10) || !strings.EqualFold(fields["SHA-1"], fmt.Sprintf("%x", hash.Sum(nil))) {
		return fmt.Errorf("zsync length or SHA-1 does not match the AppImage")
	}
	return nil
}

func (options Options) outputPaths() []string {
	paths := []string{options.Output}
	if options.UpdateInfo != "" {
		paths = append(paths, options.Output+".zsync")
	}
	return paths
}
