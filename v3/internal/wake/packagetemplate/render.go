// Package packagetemplate copies user-owned package replacements into a
// disposable workspace. Replacement contents and names are opaque: Wails does
// not parse, interpolate, merge, sanitise, or normalise them.
package packagetemplate

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type stagedFile interface {
	io.Writer
	Name() string
	Chmod(os.FileMode) error
	Close() error
}

type copyOperations struct {
	stat       func(string) (fs.FileInfo, error)
	open       func(string) (io.ReadCloser, error)
	mkdirAll   func(string, os.FileMode) error
	mkdir      func(string, os.FileMode) error
	mkdirTemp  func(string, string) (string, error)
	createTemp func(string, string) (stagedFile, error)
	openFile   func(string, int, os.FileMode) (io.WriteCloser, error)
	walkDir    func(string, fs.WalkDirFunc) error
	rel        func(string, string) (string, error)
	chmod      func(string, os.FileMode) error
	rename     func(string, string) error
	remove     func(string) error
	removeAll  func(string) error
}

func osCopyOperations() copyOperations {
	return copyOperations{
		stat: os.Stat, open: func(path string) (io.ReadCloser, error) { return os.Open(path) },
		mkdirAll: os.MkdirAll, mkdir: os.Mkdir, mkdirTemp: os.MkdirTemp,
		createTemp: func(directory, pattern string) (stagedFile, error) { return os.CreateTemp(directory, pattern) },
		openFile: func(path string, flag int, mode os.FileMode) (io.WriteCloser, error) {
			return os.OpenFile(path, flag, mode)
		},
		walkDir: filepath.WalkDir, rel: filepath.Rel, chmod: os.Chmod, rename: os.Rename, remove: os.Remove, removeAll: os.RemoveAll,
	}
}

func Copy(source, destination string) error {
	return copyWithOperations(source, destination, osCopyOperations())
}

func copyWithOperations(source, destination string, ops copyOperations) error {
	info, err := ops.stat(source)
	if err != nil {
		return fmt.Errorf("package template: %w", err)
	}
	if info.IsDir() {
		return renderDirectory(source, destination, info, ops)
	}
	return renderFile(source, destination, info, ops)
}

func renderFile(source, destination string, info fs.FileInfo, ops copyOperations) error {
	if err := ops.mkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	temporary, err := ops.createTemp(filepath.Dir(destination), ".package-template-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer ops.remove(temporaryPath)
	input, err := ops.open(source)
	if err != nil {
		_ = temporary.Close()
		return fmt.Errorf("package template: %w", err)
	}
	if _, err := io.Copy(temporary, input); err != nil {
		_ = input.Close()
		_ = temporary.Close()
		return err
	}
	if err := input.Close(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(info.Mode().Perm()); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := replacePath(temporaryPath, destination, ops); err != nil {
		return err
	}
	return nil
}

func renderDirectory(source, destination string, sourceInfo fs.FileInfo, ops copyOperations) error {
	parent := filepath.Dir(destination)
	if err := ops.mkdirAll(parent, 0o755); err != nil {
		return err
	}
	temporary, err := ops.mkdirTemp(parent, ".package-template-*")
	if err != nil {
		return err
	}
	defer ops.removeAll(temporary)
	if err := ops.chmod(temporary, sourceInfo.Mode().Perm()); err != nil {
		return err
	}
	err = ops.walkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := ops.rel(source, path)
		if err != nil || relative == "." {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("package template contains unsupported symlink: %s", relative)
		}
		target := filepath.Join(temporary, relative)
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if err := ops.mkdir(target, info.Mode().Perm()); err != nil {
				return err
			}
			return ops.chmod(target, info.Mode().Perm())
		}
		return copyFile(path, target, ops)
	})
	if err != nil {
		return err
	}
	return replacePath(temporary, destination, ops)
}

func copyFile(source, destination string, ops copyOperations) error {
	info, err := ops.stat(source)
	if err != nil {
		return err
	}
	if err := ops.mkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	in, err := ops.open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := ops.openFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return ops.chmod(destination, info.Mode().Perm())
}

func replacePath(source, destination string, ops copyOperations) error {
	backupDirectory, err := ops.mkdirTemp(filepath.Dir(destination), ".package-template-previous-*")
	if err != nil {
		return err
	}
	defer ops.removeAll(backupDirectory)
	backup := filepath.Join(backupDirectory, "previous")
	hadDestination := false
	if _, err := ops.stat(destination); err == nil {
		if err := ops.rename(destination, backup); err != nil {
			return err
		}
		hadDestination = true
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := ops.rename(source, destination); err != nil {
		if hadDestination {
			_ = ops.rename(backup, destination)
		}
		return err
	}
	return nil
}
