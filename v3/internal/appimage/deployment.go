package appimage

import (
	"context"
	"crypto/sha256"
	"debug/elf"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/stagefile"
)

// runtimeDeployment resolves Linux names against the supplied runtime tree. It
// never uses the packaging host's loader search path for a supplied tree.
type runtimeDeployment struct {
	root    string
	arch    string
	gtk     int
	dirs    []string
	files   map[string]string
	helpers []string
}

func newRuntimeDeployment(root, arch string, gtk int, extra []string) (*runtimeDeployment, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	dirs := []string{"/lib64", "/usr/lib64", "/lib", "/usr/lib", "/usr/local/lib"}
	triplet := map[string]string{"amd64": "x86_64-linux-gnu", "arm64": "aarch64-linux-gnu"}[arch]
	if triplet == "" {
		return nil, fmt.Errorf("unsupported Linux runtime architecture %q", arch)
	}
	dirs = append([]string{"/lib/" + triplet, "/usr/lib/" + triplet}, dirs...)
	for _, dir := range extra {
		if !strings.HasPrefix(dir, "/") {
			return nil, fmt.Errorf("runtime library directory must be absolute within the Linux tree: %s", dir)
		}
	}
	return &runtimeDeployment{root: absolute, arch: arch, gtk: gtk, dirs: append(slices.Clone(extra), dirs...), files: map[string]string{}}, nil
}

// resolve follows target-tree symlinks, including absolute Linux links, without
// allowing them to escape into the build host's filesystem.
func (d *runtimeDeployment) resolve(name string) (string, error) {
	pending := strings.Split(strings.TrimPrefix(path.Clean("/"+name), "/"), "/")
	resolved := []string{}
	links := 0
	for len(pending) > 0 {
		part := pending[0]
		pending = pending[1:]
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			if len(resolved) == 0 {
				return "", fmt.Errorf("runtime path escapes root: %s", name)
			}
			resolved = resolved[:len(resolved)-1]
			continue
		}
		candidate := filepath.Join(append([]string{d.root}, append(resolved, part)...)...)
		info, err := os.Lstat(candidate)
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			links++
			if links > 40 {
				return "", fmt.Errorf("too many runtime symlinks: %s", name)
			}
			target, err := os.Readlink(candidate)
			if err != nil {
				return "", err
			}
			target = filepath.ToSlash(target)
			if strings.HasPrefix(target, "/") {
				resolved = nil
			}
			pending = append(strings.Split(target, "/"), pending...)
		} else {
			resolved = append(resolved, part)
		}
	}
	return filepath.Join(append([]string{d.root}, resolved...)...), nil
}

func (d *runtimeDeployment) library(name string) (string, error) {
	return d.libraryIn(name, d.dirs)
}

func (d *runtimeDeployment) libraryIn(name string, dirs []string) (string, error) {
	if strings.Contains(name, "/") {
		return "", fmt.Errorf("absolute or path-qualified ELF dependency is unsupported: %s", name)
	}
	for _, dir := range dirs {
		filename, err := d.resolve(path.Join(dir, name))
		if err != nil {
			continue
		}
		resolved, err := resolveLibrary(filename, nil, machineForArch(d.arch))
		if err == nil {
			return resolved, nil
		}
	}
	return "", fmt.Errorf("Linux %s runtime library %s not found in %s", d.arch, name, d.root)
}

// These libraries belong to the target's C runtime. Mixing a bundled glibc
// library with the host's ELF interpreter is not a supported deployment model.
func systemRuntimeLibrary(name string) bool {
	switch name {
	case "libc.so.6", "libm.so.6", "libdl.so.2", "libpthread.so.0", "librt.so.1", "libresolv.so.2", "libutil.so.1", "ld-linux-x86-64.so.2", "ld-linux-aarch64.so.1":
		return true
	}
	return false
}

func (d *runtimeDeployment) add(destination, source string) error {
	if previous, ok := d.files[destination]; ok {
		if previous == source {
			return nil
		}
		same, err := sameRuntimeFile(previous, source)
		if err != nil {
			return err
		}
		if !same {
			return fmt.Errorf("runtime files conflict at %s: %s and %s", destination, previous, source)
		}
		return nil
	}
	d.files[destination] = source
	return nil
}

func sameRuntimeFile(a, b string) (bool, error) {
	hash := func(name string) ([32]byte, error) {
		file, err := os.Open(name)
		if err != nil {
			return [32]byte{}, err
		}
		defer file.Close()
		digest := sha256.New()
		if _, err = io.Copy(digest, file); err != nil {
			return [32]byte{}, err
		}
		return [32]byte(digest.Sum(nil)), nil
	}
	first, err := hash(a)
	if err != nil {
		return false, err
	}
	second, err := hash(b)
	return first == second, err
}

func (d *runtimeDeployment) dependencies(ctx context.Context, roots []string) error {
	queue := slices.Clone(roots)
	seen := map[string]bool{}
	for len(queue) > 0 {
		if err := ctx.Err(); err != nil {
			return err
		}
		filename := queue[0]
		queue = queue[1:]
		if seen[filename] {
			continue
		}
		seen[filename] = true
		image, err := elf.Open(filename)
		if err != nil {
			return err
		}
		if image.Machine != machineForArch(d.arch) {
			image.Close()
			return fmt.Errorf("runtime ELF target mismatch: %s", filename)
		}
		needed, err := image.ImportedLibraries()
		search := d.searchPaths(image, filename)
		image.Close()
		if err != nil {
			return err
		}
		for _, name := range needed {
			gtk := 0
			if strings.HasPrefix(name, "libgtk-4.so") {
				gtk = 4
			}
			if strings.HasPrefix(name, "libgtk-3.so") {
				gtk = 3
			}
			if gtk != 0 {
				if d.gtk != 0 && d.gtk != gtk {
					return fmt.Errorf("runtime contains incompatible GTK versions")
				}
				d.gtk = gtk
			}
			if systemRuntimeLibrary(name) {
				continue
			}
			source, err := d.libraryIn(name, search)
			if err != nil {
				return fmt.Errorf("%s: %w", filename, err)
			}
			if err = d.add(path.Join("usr/lib", name), source); err != nil {
				return err
			}
			queue = append(queue, source)
		}
	}
	return nil
}

func (d *runtimeDeployment) profile(ctx context.Context) error {
	names := []string{"libglib-2.0.so.0", "libgobject-2.0.so.0", "libgio-2.0.so.0", "libsoup-3.0.so.0"}
	switch d.gtk {
	case 3:
		names = append(names, "libgtk-3.so.0", "libwebkit2gtk-4.1.so.0", "libjavascriptcoregtk-4.1.so.0")
	case 4:
		names = append(names, "libgtk-4.so.1", "libwebkitgtk-6.0.so.4", "libjavascriptcoregtk-6.0.so.1")
	default:
		return fmt.Errorf("WebKit runtime requires GTK version 3 or 4")
	}
	var roots []string
	var dataDirectory string
	for _, name := range names {
		source, err := d.library(name)
		if err != nil {
			return err
		}
		if err = d.add(path.Join("usr/lib", name), source); err != nil {
			return err
		}
		roots = append(roots, source)
		if strings.HasPrefix(name, "libgtk-") {
			relative, err := filepath.Rel(d.root, filepath.Dir(source))
			if err != nil {
				return err
			}
			dataDirectory = "/" + filepath.ToSlash(relative)
		}
	}
	// Global shortcuts load X11 dynamically when available, outside DT_NEEDED.
	if source, err := d.library("libX11.so.6"); err == nil {
		if err = d.add("usr/lib/libX11.so.6", source); err != nil {
			return err
		}
		roots = append(roots, source)
	}
	if err := d.dependencies(ctx, roots); err != nil {
		return err
	}
	version := "webkitgtk-6.0"
	bundle := "libwebkitgtkinjectedbundle.so"
	if d.gtk == 3 {
		version = "webkit2gtk-4.1"
		bundle = "libwebkit2gtkinjectedbundle.so"
	}
	var helperDir string
	for _, dir := range append([]string{"/usr/libexec", "/usr/lib"}, d.dirs...) {
		candidate := path.Join(dir, version)
		if _, err := d.resolve(path.Join(candidate, "WebKitWebProcess")); err == nil {
			helperDir = candidate
			break
		}
	}
	if helperDir == "" {
		return fmt.Errorf("no %s helpers in runtime tree", version)
	}
	d.helpers = []string{helperDir}
	roots = nil
	for _, name := range []string{"WebKitWebProcess", "WebKitNetworkProcess", "WebKitGPUProcess"} {
		source, err := d.resolve(path.Join(helperDir, name))
		if name == "WebKitGPUProcess" && os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if err = d.add(strings.TrimPrefix(path.Join(helperDir, name), "/"), source); err != nil {
			return err
		}
		roots = append(roots, source)
	}
	foundBundle := false
	for _, dir := range d.dirs {
		logical := path.Join(dir, version, "injected-bundle", bundle)
		source, err := d.resolve(logical)
		if err != nil {
			continue
		}
		if err = d.add(strings.TrimPrefix(logical, "/"), source); err != nil {
			return err
		}
		roots = append(roots, source)
		foundBundle = true
		break
	}
	if !foundBundle {
		return fmt.Errorf("no %s injected bundle in runtime tree", version)
	}
	if err := d.dependencies(ctx, roots); err != nil {
		return err
	}
	if _, err := d.resolve("/usr/share/glib-2.0/schemas/gschemas.compiled"); err != nil {
		return fmt.Errorf("runtime tree must contain compiled GLib schemas: %w", err)
	}
	for _, directory := range []string{"/usr/share/glib-2.0/schemas"} {
		if err := d.tree(ctx, directory, strings.TrimPrefix(directory, "/"), false); err != nil {
			return err
		}
	}
	for _, name := range []string{"girepository-1.0", "gio/modules", "gdk-pixbuf-2.0", fmt.Sprintf("gtk-%d.0", d.gtk)} {
		if err := d.tree(ctx, path.Join(dataDirectory, name), path.Join("usr/lib", name), true); err != nil {
			return err
		}
	}
	return d.checkLoaderCaches()
}

func (d *runtimeDeployment) tree(ctx context.Context, logical, destination string, libraries bool) error {
	root, err := d.resolve(logical)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var modules []string
	err = filepath.WalkDir(root, func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, filename)
		if err != nil {
			return err
		}
		source, err := d.resolve(path.Join(logical, filepath.ToSlash(relative)))
		if err != nil {
			return err
		}
		info, err := os.Stat(source)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		target := path.Join(destination, filepath.ToSlash(relative))
		if err = d.add(target, source); err != nil {
			return err
		}
		if libraries && strings.Contains(entry.Name(), ".so") {
			image, err := elf.Open(source)
			if err != nil {
				return nil
			}
			image.Close()
			if err = d.add(path.Join("usr/lib", entry.Name()), source); err != nil {
				return err
			}
			modules = append(modules, source)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return d.dependencies(ctx, modules)
}

func (d *runtimeDeployment) stage(ctx context.Context, appdir string) error {
	buffer := make([]byte, 128<<10)
	for destination, source := range d.files {
		mode := os.FileMode(0644)
		if strings.HasPrefix(path.Base(destination), "WebKit") {
			mode = 0755
		}
		target := filepath.Join(appdir, filepath.FromSlash(destination))
		if err := stagefile.Copy(ctx, source, target, mode, buffer); err != nil {
			return err
		}
		if err := preferBundledLibraries(target); err != nil {
			return err
		}
		if strings.HasSuffix(destination, "loaders.cache") || strings.HasSuffix(destination, "immodules.cache") {
			data, err := os.ReadFile(target)
			if err != nil {
				return err
			}
			lines := strings.Split(string(data), "\n")
			for i, line := range lines {
				if strings.HasPrefix(line, "\"/") {
					end := strings.Index(line[1:], "\"") + 1
					if end > 1 && strings.Contains(line[:end], ".so") {
						lines[i] = "\"" + path.Base(line[1:end]) + line[end:]
					}
				}
			}
			if err = os.WriteFile(target, []byte(strings.Join(lines, "\n")), mode); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *runtimeDeployment) searchPaths(image *elf.File, filename string) []string {
	search := slices.Clone(d.dirs)
	relative, relativeErr := filepath.Rel(d.root, filename)
	for _, tag := range []elf.DynTag{elf.DT_RPATH, elf.DT_RUNPATH} {
		values, _ := image.DynString(tag)
		for _, value := range values {
			for _, dir := range strings.Split(value, ":") {
				if strings.Contains(dir, "$ORIGIN") || strings.Contains(dir, "${ORIGIN}") {
					if relativeErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
						continue
					}
					origin := path.Dir("/" + filepath.ToSlash(relative))
					dir = strings.ReplaceAll(strings.ReplaceAll(dir, "${ORIGIN}", origin), "$ORIGIN", origin)
				}
				if strings.HasPrefix(dir, "/") {
					search = append([]string{dir}, search...)
				}
			}
		}
	}
	return search
}

func (d *runtimeDeployment) checkLoaderCaches() error {
	for destination := range d.files {
		if !strings.HasSuffix(destination, ".so") {
			continue
		}
		directory := path.Dir(destination)
		var cache string
		switch {
		case strings.HasPrefix(destination, "usr/lib/gdk-pixbuf-2.0/") && path.Base(directory) == "loaders":
			cache = path.Join(path.Dir(directory), "loaders.cache")
		case strings.HasPrefix(destination, "usr/lib/gtk-3.0/") && path.Base(directory) == "immodules":
			cache = path.Join(path.Dir(directory), "immodules.cache")
		}
		if cache != "" && d.files[cache] == "" {
			return fmt.Errorf("runtime tree must include pre-generated %s for its loader modules", cache)
		}
	}
	return nil
}
