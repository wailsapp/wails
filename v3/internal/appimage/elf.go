package appimage

import (
	"bufio"
	"context"
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type elfAnalysis struct {
	gtk  string
	relr bool
}

func analyzeELF(ctx context.Context, binaries, dirs []string) (elfAnalysis, error) {
	var result elfAnalysis
	queue := slices.Clone(binaries)
	seen := map[string]bool{}
	for len(queue) > 0 {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		filename := queue[0]
		queue = queue[1:]
		resolved, err := filepath.EvalSymlinks(filename)
		if err != nil {
			return result, err
		}
		if seen[resolved] {
			continue
		}
		seen[resolved] = true
		image, err := elf.Open(filename)
		if err != nil {
			return result, fmt.Errorf("inspect %s: %w", filename, err)
		}
		libraries, err := image.ImportedLibraries()
		search := elfSearchPaths(image, filename, dirs)
		result.relr = result.relr || image.Section(".relr.dyn") != nil
		machine := image.Machine
		image.Close()
		if err != nil {
			return result, err
		}
		for _, name := range libraries {
			version := ""
			switch {
			case strings.HasPrefix(name, "libgtk-4.so"):
				version = "4"
			case strings.HasPrefix(name, "libgtk-3.so"):
				version = "3"
			case strings.HasPrefix(name, "libgtk-x11-2.0.so"):
				version = "2"
			}
			if version != "" {
				if result.gtk != "" && result.gtk != version {
					return result, fmt.Errorf("AppImage binaries use different GTK major versions")
				}
				result.gtk = version
			}
			dependency, err := resolveLibrary(name, search, machine)
			if err != nil {
				return result, fmt.Errorf("%s: %w", filename, err)
			}
			queue = append(queue, dependency)
		}
	}
	return result, nil
}

func elfSearchPaths(image *elf.File, filename string, dirs []string) []string {
	var search []string
	for _, tag := range []elf.DynTag{elf.DT_RPATH, elf.DT_RUNPATH} {
		entries, _ := image.DynString(tag)
		for _, entry := range entries {
			for _, dir := range strings.Split(entry, ":") {
				dir = strings.ReplaceAll(strings.ReplaceAll(dir, "${ORIGIN}", filepath.Dir(filename)), "$ORIGIN", filepath.Dir(filename))
				if filepath.IsAbs(dir) {
					search = append(search, dir)
				}
			}
		}
	}
	return append(search, dirs...)
}

func resolveLibrary(name string, dirs []string, machine elf.Machine) (string, error) {
	if filepath.IsAbs(name) {
		dirs = []string{""}
	}
	for _, dir := range dirs {
		filename := name
		if dir != "" {
			filename = filepath.Join(dir, name)
		}
		image, err := elf.Open(filename)
		if err != nil {
			continue
		}
		matches := image.Machine == machine
		image.Close()
		if matches {
			return filename, nil
		}
	}
	return "", fmt.Errorf("ELF dependency %q not found for %s; configure libraryDirs for a custom runtime", name, machine)
}

func defaultLibraryDirs(arch string) []string {
	dirs := []string{"/lib64", "/usr/lib64", "/lib", "/usr/lib", "/usr/local/lib"}
	triplet := map[string]string{"amd64": "x86_64-linux-gnu", "arm64": "aarch64-linux-gnu"}[arch]
	if triplet != "" {
		dirs = append([]string{filepath.Join("/lib", triplet), filepath.Join("/usr/lib", triplet)}, dirs...)
	}
	return append(dirs, libraryConfigDirs("/etc/ld.so.conf", map[string]bool{})...)
}

func libraryConfigDirs(filename string, seen map[string]bool) []string {
	if seen[filename] {
		return nil
	}
	seen[filename] = true
	input, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer input.Close()
	var dirs []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line, _, _ := strings.Cut(scanner.Text(), "#")
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == "include" {
			for _, pattern := range fields[1:] {
				if !filepath.IsAbs(pattern) {
					pattern = filepath.Join(filepath.Dir(filename), pattern)
				}
				matches, _ := filepath.Glob(pattern)
				for _, match := range matches {
					dirs = append(dirs, libraryConfigDirs(match, seen)...)
				}
			}
		} else if len(fields) == 1 && filepath.IsAbs(fields[0]) {
			dirs = append(dirs, fields[0])
		}
	}
	return dirs
}

// FindRuntimeFiles searches runtime directories first, with a full-root fallback
// for nonstandard installations. Directory symlinks are not followed.
func FindRuntimeFiles(ctx context.Context, root string, names []string) ([]string, error) {
	roots := []string{}
	for _, directory := range []string{"lib", "lib64", "libexec", "local/lib", "local/lib64", "local/libexec"} {
		candidate := filepath.Join(root, directory)
		info, err := os.Lstat(candidate)
		if err == nil && info.IsDir() {
			roots = append(roots, candidate)
		}
	}
	files, missing, err := findRuntimeFiles(ctx, roots, names)
	if err != nil {
		return nil, err
	}
	if len(missing) != 0 {
		files, missing, err = findRuntimeFiles(ctx, []string{root}, names)
	}
	if err != nil {
		return nil, err
	}
	if len(missing) != 0 {
		return nil, fmt.Errorf("unable to locate required WebKit files: %s", strings.Join(missing, ", "))
	}
	return files, nil
}

func findRuntimeFiles(ctx context.Context, roots, names []string) ([]string, []string, error) {
	requested := map[string]bool{}
	for _, name := range names {
		requested[name] = true
	}
	found := map[string]bool{}
	var result []string
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if cancel := ctx.Err(); cancel != nil {
				return cancel
			}
			if err != nil {
				if os.IsPermission(err) {
					return nil
				}
				return err
			}
			if !entry.IsDir() && requested[entry.Name()] {
				found[entry.Name()] = true
				result = append(result, path)
			}
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}
	var missing []string
	for _, name := range names {
		if !found[name] {
			missing = append(missing, name)
		}
	}
	return result, missing, ctx.Err()
}

// Multiple WebKit stacks often install helpers with identical basenames. Match
// their linked GTK version instead of bundling every installed stack.
func matchingRuntimeFiles(ctx context.Context, files []string, gtk string, dirs []string) ([]string, error) {
	var result []string
	expected := map[string]bool{}
	matched := map[string]bool{}
	for _, filename := range files {
		name := filepath.Base(filename)
		expected[name] = true
		analysis, err := analyzeELF(ctx, []string{filename}, dirs)
		if err != nil {
			return nil, err
		}
		if analysis.gtk != "" && analysis.gtk != gtk {
			continue
		}
		result = append(result, filename)
		matched[name] = true
	}
	for name := range expected {
		if !matched[name] {
			return nil, fmt.Errorf("no GTK%s-compatible WebKit helper %q found", gtk, name)
		}
	}
	return result, nil
}
