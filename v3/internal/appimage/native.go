package appimage

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/readonlyfs"
)

// CreateNative builds a type-2 AppImage without executing target binaries or
// packaging tools. Its target runtime is an explicitly pinned external artifact.
func CreateNative(ctx context.Context, options Options) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	options, err := normalizeNativeOptions(options)
	if err != nil {
		return err
	}
	var runtimePath string
	if err := options.runStage("Prepare pinned AppImage runtime", func() error { var err error; runtimePath, err = options.PrepareRuntime(ctx); return err }); err != nil {
		return err
	}
	base := filepath.Dir(options.Output)
	if options.WorkDir != "" {
		base = options.path(options.WorkDir)
	}
	if err := os.MkdirAll(base, 0755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(base, ".wails-native-appimage-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	var appdir string
	if err := options.runStage("Collect Linux runtime and stage AppImage", func() error { var err error; appdir, err = options.prepareNativeDirectory(ctx, stage); return err }); err != nil {
		return err
	}
	result := filepath.Join(stage, options.resultName())
	if err := options.runStage("Compress AppImage filesystem", func() error {
		output, err := os.OpenFile(result, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
		if err != nil {
			return err
		}
		err = copyRuntime(ctx, runtimePath, output, options.UpdateInfo)
		if err == nil {
			entries, entryErr := imageEntries(appdir)
			if entryErr != nil {
				err = entryErr
			} else {
				err = readonlyfs.Write(ctx, output, entries, readonlyfs.Options{BlockSize: 128 << 10, ModTime: time.Unix(0, 0)})
			}
		}
		closeErr := output.Close()
		if err != nil {
			return err
		}
		return closeErr
	}); err != nil {
		return err
	}
	return options.runStage("Verify and publish AppImage", func() error {
		if err := verify(result, options.Arch); err != nil {
			return err
		}
		if options.UpdateInfo != "" {
			if err := writeUpdateControl(ctx, result); err != nil {
				return err
			}
		}
		if err := os.MkdirAll(filepath.Dir(options.Output), 0755); err != nil {
			return err
		}
		return options.publishResults(ctx, result)
	})
}

func normalizeNativeOptions(options Options) (Options, error) {
	root, err := filepath.Abs(options.Root)
	if err != nil {
		return options, err
	}
	options.Root = root
	if options.Arch == "" {
		options.Arch = runtime.GOARCH
	}
	if _, err := architectureLabel(options.Arch); err != nil {
		return options, err
	}
	if options.GTKVersion != 0 && options.GTKVersion != 3 && options.GTKVersion != 4 {
		return options, fmt.Errorf("AppImage GTK version must be 3 or 4")
	}
	if len(options.Binaries) == 0 || options.Icon == "" || options.DesktopFile == "" {
		return options, fmt.Errorf("AppImage needs binaries, icon and desktop file")
	}
	if options.Name == "" {
		options.Name = options.Binaries[0].Name
	}
	options, err = normalizeOutputOptions(options)
	if err != nil {
		return options, err
	}
	if options.IconName != "" && (options.IconName != filepath.Base(options.IconName) || options.IconName == "." || options.IconName == ".." || options.IconName == ".DirIcon") {
		return options, fmt.Errorf("invalid bundled icon name %q", options.IconName)
	}
	options.Binaries = slices.Clone(options.Binaries)
	seen := map[string]bool{}
	var inputs []string
	for i, binary := range options.Binaries {
		if binary.Name == "" || binary.Name == "." || binary.Name == ".." || strings.ContainsAny(binary.Name, "/\\\x00") || seen[binary.Name] {
			return options, fmt.Errorf("invalid or duplicate AppImage executable %q", binary.Name)
		}
		seen[binary.Name] = true
		options.Binaries[i].Path = options.path(binary.Path)
		inputs = append(inputs, options.Binaries[i].Path)
	}
	for _, input := range []string{options.Icon, options.DesktopFile, options.RuntimePath} {
		if input != "" {
			inputs = append(inputs, options.path(input))
		}
	}
	for _, input := range inputs {
		if slices.Contains(options.outputPaths(), input) {
			return options, fmt.Errorf("AppImage output conflicts with input %s", input)
		}
	}
	return options, nil
}

func imageEntries(root string) ([]readonlyfs.Entry, error) {
	var entries []readonlyfs.Entry
	err := filepath.WalkDir(root, func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filename == root {
			return nil
		}
		relative, err := filepath.Rel(root, filename)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := readonlyfs.Entry{Path: name, Mode: 0644, Size: info.Size()}
		switch {
		case entry.IsDir():
			value.Mode = fs.ModeDir | 0755
		case info.Mode()&fs.ModeSymlink != 0:
			value.Mode = fs.ModeSymlink | 0777
			value.Link, err = os.Readlink(filename)
			if err != nil {
				return err
			}
		case info.Mode().IsRegular():
			if name == "AppRun" || strings.HasPrefix(name, "usr/bin/") || strings.HasPrefix(filepath.Base(name), "WebKit") {
				value.Mode = 0755
			}
			value.Open = func() (io.ReadCloser, error) { return os.Open(filename) }
		default:
			return fmt.Errorf("unsupported AppImage entry %s", filename)
		}
		entries = append(entries, value)
		return nil
	})
	return entries, err
}
