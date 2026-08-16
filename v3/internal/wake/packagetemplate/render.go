// Package packagetemplate renders user-owned package inputs from a stable,
// format-neutral model. It deliberately exposes manifest intent and resolved
// paths, never Pipeline Nodes or process environment.
package packagetemplate

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type Model struct {
	Version      int
	Project      manifest.Project
	Target       Target
	Package      Package
	Paths        Paths
	Associations []manifest.Association
	Protocols    []manifest.Protocol
	Options      map[string]any
}

type Target struct {
	OS, Arch, Variant, MinimumVersion string
	Capabilities                      []string
}

type Package struct {
	Format string
}

type Paths struct {
	Project, Binary, Output, Assets, Icon, Workspace string
}

func Render(source, destination string, model Model) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("package template: %w", err)
	}
	if info.IsDir() {
		return renderDirectory(source, destination, model)
	}
	return renderFile(source, destination, model)
}

func renderFile(source, destination string, model Model) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("package template: %w", err)
	}
	contents, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("package template: %w", err)
	}
	tmpl, err := template.New(filepath.Base(source)).Option("missingkey=error").Parse(string(contents))
	if err != nil {
		return fmt.Errorf("parse package template %s: %w", source, err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(destination), ".package-template-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := tmpl.Execute(temporary, model); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("render package template %s: %w", source, err)
	}
	if err := temporary.Chmod(info.Mode().Perm()); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := replacePath(temporaryPath, destination); err != nil {
		return err
	}
	return nil
}

func renderDirectory(source, destination string, model Model) error {
	parent := filepath.Dir(destination)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	temporary, err := os.MkdirTemp(parent, ".package-template-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temporary)
	err = filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil || relative == "." {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("package template contains unsupported symlink: %s", relative)
		}
		target := filepath.Join(temporary, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if strings.HasSuffix(target, ".tmpl") {
			target = strings.TrimSuffix(target, ".tmpl")
			return renderFile(path, target, model)
		}
		return copyFile(path, target)
	})
	if err != nil {
		return err
	}
	return replacePath(temporary, destination)
}

func copyFile(source, destination string) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func replacePath(source, destination string) error {
	backupDirectory, err := os.MkdirTemp(filepath.Dir(destination), ".package-template-previous-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(backupDirectory)
	backup := filepath.Join(backupDirectory, "previous")
	hadDestination := false
	if _, err := os.Stat(destination); err == nil {
		if err := os.Rename(destination, backup); err != nil {
			return err
		}
		hadDestination = true
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(source, destination); err != nil {
		if hadDestination {
			_ = os.Rename(backup, destination)
		}
		return err
	}
	return nil
}
