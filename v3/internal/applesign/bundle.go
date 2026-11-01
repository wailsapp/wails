package applesign

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/stagefile"
	"howett.net/plist"
)

// SignBundle copies and signs a canonical macOS Contents-layout app bundle.
// Nested application bundles and Mach-O helpers are signed from the inside out.
// Output must not exist; callers can transactionally publish the completed copy.
func SignBundle(ctx context.Context, input, output string, o Options) (err error) {
	if err = CopyBundle(ctx, input, output); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(output)
		}
	}()
	_, err = signBundleCopy(ctx, output, o, 0)
	return err
}

// CopyBundle creates a fresh copy for signing or stapling, preserving regular
// file modes and relative symlinks while rejecting links outside the input.
// It removes partial output on failure and never overwrites an existing path.
func CopyBundle(ctx context.Context, input, output string) (err error) {
	input, err = filepath.Abs(input)
	if err != nil {
		return err
	}
	output, err = filepath.Abs(output)
	if err != nil {
		return err
	}
	st, err := os.Lstat(input)
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return errors.New("application input must be a directory, not a symlink")
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(output))
	if err != nil {
		return err
	}
	inputReal, err := filepath.EvalSymlinks(input)
	if err != nil {
		return err
	}
	if isWithin(inputReal, filepath.Join(parent, filepath.Base(output))) {
		return errors.New("signed bundle output must be outside its input")
	}
	if err = os.Mkdir(output, 0755); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(output)
		}
	}()
	buffer := make([]byte, 64<<10)
	count := 0
	err = filepath.WalkDir(input, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		count++
		if count > 100000 {
			return errors.New("application exceeds 100000 entries")
		}
		rel, err := filepath.Rel(input, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dest := filepath.Join(output, rel)
		if entry.IsDir() {
			return os.Mkdir(dest, 0755)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			if filepath.IsAbs(link) || !isWithin(input, filepath.Join(filepath.Dir(path), link)) {
				return fmt.Errorf("bundle symlink escapes application: %s", rel)
			}
			return os.Symlink(link, dest)
		}
		return stagefile.Copy(ctx, path, dest, 0, buffer)
	})
	if err != nil {
		return err
	}
	return nil
}

type bundleMetadata struct {
	Executable string `plist:"CFBundleExecutable"`
	Identifier string `plist:"CFBundleIdentifier"`
}

func readBundleInfo(path string) ([]byte, bundleMetadata, error) {
	var metadata bundleMetadata
	info, err := readLimited(path, 1<<20)
	if err != nil {
		return nil, metadata, err
	}
	if _, err := plist.Unmarshal(info, &metadata); err != nil {
		return nil, metadata, fmt.Errorf("bundle Info.plist: %w", err)
	}
	if metadata.Executable == "" || metadata.Executable == "." || metadata.Executable == ".." || strings.ContainsAny(metadata.Executable, "/\\\x00") {
		return nil, metadata, errors.New("invalid CFBundleExecutable")
	}
	if metadata.Identifier == "" {
		return nil, metadata, errors.New("bundle needs CFBundleIdentifier")
	}
	return info, metadata, nil
}

// BundleExecutable locates a canonical app's main executable without following
// symlinked metadata or code directories. The returned path is absolute.
func BundleExecutable(root string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	for _, dir := range []string{root, filepath.Join(root, "Contents"), filepath.Join(root, "Contents/MacOS")} {
		st, err := os.Lstat(dir)
		if err != nil {
			return "", err
		}
		if !st.IsDir() {
			return "", errors.New("bundle code directories must not be symlinks")
		}
	}
	infoPath := filepath.Join(root, "Contents/Info.plist")
	st, err := os.Lstat(infoPath)
	if err != nil {
		return "", err
	}
	if !st.Mode().IsRegular() {
		return "", errors.New("bundle Info.plist must be a regular file")
	}
	_, metadata, err := readBundleInfo(infoPath)
	if err != nil {
		return "", err
	}
	main := filepath.Join(root, "Contents/MacOS", metadata.Executable)
	st, err = os.Lstat(main)
	if err != nil {
		return "", err
	}
	if !st.Mode().IsRegular() {
		return "", errors.New("bundle executable must be a regular file")
	}
	return main, nil
}

func isWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

func signBundleCopy(ctx context.Context, root string, o Options, depth int) ([]CodeIdentity, error) {
	if depth > 16 {
		return nil, errors.New("nested bundle depth exceeds 16")
	}
	if strings.EqualFold(filepath.Ext(root), ".framework") {
		return signFrameworkCopy(ctx, root, o, depth)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	if len(entries) != 1 || entries[0].Name() != "Contents" || !entries[0].IsDir() {
		return nil, errors.New("application bundle must contain only a Contents directory")
	}
	contents := filepath.Join(root, "Contents")
	infoPath := filepath.Join(contents, "Info.plist")
	return signBundleLayout(ctx, contents, infoPath, "MacOS", o, depth)
}

func signBundleLayout(ctx context.Context, contents, infoPath, executables string, o Options, depth int) ([]CodeIdentity, error) {
	info, metadata, err := readBundleInfo(infoPath)
	if err != nil {
		return nil, err
	}
	o.Identifier = metadata.Identifier
	if err := validateOptions(o); err != nil {
		return nil, err
	}
	main := filepath.Join(contents, executables, metadata.Executable)
	st, err := os.Lstat(main)
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, errors.New("bundle executable must be a regular file")
	}
	signatureDir := filepath.Join(contents, "_CodeSignature")
	if err := os.RemoveAll(signatureDir); err != nil {
		return nil, err
	}
	files := make(map[string]any)
	metadataBytes := 0
	err = filepath.WalkDir(contents, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if path == contents || path == main || (path == infoPath && filepath.Dir(infoPath) == contents) || path == filepath.Join(contents, "PkgInfo") {
			return nil
		}
		rel, err := filepath.Rel(contents, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(rel)
		metadataBytes += len(name) + 256
		if metadataBytes > 32<<20 {
			return errors.New("bundle resource metadata exceeds 32 MiB")
		}
		if entry.IsDir() {
			if nestedLocation(name) && strings.Contains(entry.Name(), ".") {
				switch strings.ToLower(filepath.Ext(entry.Name())) {
				case ".app", ".xpc", ".bundle", ".appex", ".framework":
				default:
					return fmt.Errorf("unsupported nested code bundle %s", name)
				}
				ids, err := signBundleCopy(ctx, path, o, depth+1)
				if err != nil {
					return err
				}
				files[name] = nestedSeal(ids)
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			files[name] = map[string]any{"symlink": filepath.ToSlash(link)}
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("unsupported bundle resource %s", name)
		}
		if nestedLocation(name) {
			mach, err := isMachO(path)
			if err != nil {
				return err
			}
			if !mach {
				return fmt.Errorf("non-Mach-O file in nested-code location: %s", name)
			}
			nested := o
			nested.Identifier += "." + strings.ReplaceAll(name, "/", ".")
			nested.Entitlements = nil
			nested.info, nested.resources = nil, nil
			ids, err := signFileCopy(ctx, path, nested)
			if err != nil {
				return err
			}
			files[name] = nestedSeal(ids)
			return nil
		}
		sum, err := fileHash(ctx, path)
		if err != nil {
			return err
		}
		files[name] = map[string]any{"hash2": sum}
		return nil
	})
	if err != nil {
		return nil, err
	}
	rules := map[string]any{"^.*": true, "^Info\\.plist$": map[string]any{"omit": true, "weight": float64(20)}, "^PkgInfo$": map[string]any{"omit": true, "weight": float64(20)}, "^(MacOS|Frameworks|SharedFrameworks|PlugIns|Plug-ins|XPCServices|Helpers|Library/(Automator|Spotlight|LoginItems|LaunchServices))/": map[string]any{"nested": true, "weight": float64(10)}}
	resources, err := plist.Marshal(map[string]any{"files": map[string]any{}, "rules": map[string]any{}, "files2": files, "rules2": rules}, plist.XMLFormat)
	if err != nil {
		return nil, err
	}
	if err := os.Mkdir(signatureDir, 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(signatureDir, "CodeResources"), resources, 0644); err != nil {
		return nil, err
	}
	o.info, o.resources = info, resources
	return signFileCopy(ctx, main, o)
}

func signFrameworkCopy(ctx context.Context, root string, o Options, depth int) ([]CodeIdentity, error) {
	versions := filepath.Join(root, "Versions")
	st, err := os.Lstat(versions)
	if os.IsNotExist(err) {
		return signBundleLayout(ctx, root, filepath.Join(root, "Resources/Info.plist"), "", o, depth)
	}
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, errors.New("framework Versions must be a directory")
	}
	current, err := os.Readlink(filepath.Join(versions, "Current"))
	if err != nil {
		return nil, fmt.Errorf("framework needs Versions/Current symlink: %w", err)
	}
	if current == "" || current == "." || current == ".." || strings.ContainsAny(current, "/\\") {
		return nil, errors.New("framework Current must name one version directory")
	}
	entries, err := os.ReadDir(versions)
	if err != nil {
		return nil, err
	}
	var selected []CodeIdentity
	for _, entry := range entries {
		if entry.Name() == "Current" {
			continue
		}
		if !entry.IsDir() {
			return nil, fmt.Errorf("unsupported framework version entry %s", entry.Name())
		}
		version := filepath.Join(versions, entry.Name())
		ids, err := signBundleLayout(ctx, version, filepath.Join(version, "Resources/Info.plist"), "", o, depth)
		if err != nil {
			return nil, err
		}
		if entry.Name() == current {
			selected = ids
		}
	}
	if selected == nil {
		return nil, errors.New("framework Current version does not exist")
	}
	outer, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range outer {
		if entry.Name() == "Versions" {
			continue
		}
		if entry.Type()&os.ModeSymlink == 0 {
			return nil, fmt.Errorf("versioned framework outer entry must be a symlink: %s", entry.Name())
		}
		link, err := os.Readlink(filepath.Join(root, entry.Name()))
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(filepath.ToSlash(link), "Versions/Current/") {
			return nil, errors.New("framework outer symlinks must target Versions/Current")
		}
	}
	return selected, nil
}

func nestedLocation(name string) bool {
	for _, prefix := range []string{"MacOS/", "Frameworks/", "SharedFrameworks/", "PlugIns/", "Plug-ins/", "XPCServices/", "Helpers/", "Library/Automator/", "Library/Spotlight/", "Library/LoginItems/", "Library/LaunchServices/"} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func nestedSeal(ids []CodeIdentity) map[string]any {
	terms := make([]string, len(ids))
	for i, id := range ids {
		terms[i] = `cdhash H"` + hex.EncodeToString(id.CDHash[:]) + `"`
	}
	return map[string]any{"cdhash": ids[0].CDHash[:], "requirement": strings.Join(terms, " or ")}
}

func signFileCopy(ctx context.Context, path string, o Options) ([]CodeIdentity, error) {
	work, err := os.MkdirTemp(filepath.Dir(path), ".wails-sign-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(work)
	output := filepath.Join(work, "signed")
	ids, err := SignMachO(ctx, path, output, o)
	if err != nil {
		return nil, err
	}
	return ids, os.Rename(output, path)
}

func isMachO(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	var h [4]byte
	_, err = io.ReadFull(f, h[:])
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	m := binary.BigEndian.Uint32(h[:])
	return m == 0xcffaedfe || m == 0xcafebabe || m == 0xcafebabf, nil
}

func fileHash(ctx context.Context, path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, &contextReader{ctx, f})
	return h.Sum(nil), err
}

func readLimited(path string, limit int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err == nil && int64(len(data)) > limit {
		return nil, fmt.Errorf("%s exceeds size limit", filepath.Base(path))
	}
	return data, err
}
