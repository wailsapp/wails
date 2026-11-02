// Package appbundle assembles fresh macOS application bundles without shell tools.
package appbundle

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/stagefile"
)

type Options struct {
	Output     string
	InfoPlist  string
	Icon       string
	Assets     string
	Binaries   map[string]string
	InfoValues map[string]string
	// Prepare can sign and verify the staged bundle before replacing an old one.
	Prepare func(context.Context, string) error
}

func Create(ctx context.Context, options Options) (err error) {
	if err := validate(options); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	output, err := filepath.Abs(options.Output)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(filepath.Dir(output), ".wails-app-*")
	if err != nil {
		return err
	}
	backup := filepath.Join(stage, "previous.app")
	defer func() {
		if err != nil {
			if _, statErr := os.Lstat(backup); statErr == nil {
				err = fmt.Errorf("%w; previous bundle preserved at %s", err, backup)
				return
			}
		}
		_ = os.RemoveAll(stage)
	}()
	bundle := filepath.Join(stage, filepath.Base(output))
	files := map[string]string{"Contents/Info.plist": options.InfoPlist}
	if options.Icon != "" {
		files["Contents/Resources/icons.icns"] = options.Icon
	}
	if options.Assets != "" {
		files["Contents/Resources/Assets.car"] = options.Assets
	}
	for name, source := range options.Binaries {
		files[filepath.Join("Contents", "MacOS", name)] = source
	}
	buffer := make([]byte, 32<<10)
	for name, source := range files {
		if err = stagefile.Copy(ctx, source, filepath.Join(bundle, name), 0, buffer); err != nil {
			return fmt.Errorf("assemble %s: %w", name, err)
		}
	}
	// Metadata must be finalized before the bundle is signed.
	if err = overrideInfo(filepath.Join(bundle, "Contents", "Info.plist"), options.InfoValues); err != nil {
		return err
	}
	if options.Prepare != nil {
		if err = options.Prepare(ctx, bundle); err != nil {
			return err
		}
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return replace(bundle, output, backup)
}

func validate(options Options) error {
	if options.Output == "" || !strings.HasSuffix(filepath.Base(options.Output), ".app") {
		return fmt.Errorf("bundle output must end in .app")
	}
	if options.InfoPlist == "" || len(options.Binaries) == 0 {
		return fmt.Errorf("bundle needs Info.plist and binaries")
	}
	for name := range options.Binaries {
		if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\\x00") {
			return fmt.Errorf("invalid bundle executable name %q", name)
		}
	}
	return nil
}

func replace(source, output, backup string) error {
	_, err := os.Lstat(output)
	if os.IsNotExist(err) {
		return os.Rename(source, output)
	}
	if err != nil {
		return err
	}
	if err = os.Rename(output, backup); err != nil {
		return err
	}
	if err = os.Rename(source, output); err != nil {
		return errors.Join(err, os.Rename(backup, output))
	}
	return nil
}
