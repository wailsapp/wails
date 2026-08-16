package cache

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/zeebo/blake3"
)

const (
	handlerVersion         = "wake-pipeline-1"
	metadataFastPathWindow = 2 * time.Second
)

type FileRecord struct {
	Size      int64  `json:"size"`
	ModTimeNS int64  `json:"mod_time_ns"`
	Mode      uint32 `json:"mode"`
	Identity  string `json:"identity"`
	Digest    string `json:"digest"`
}

type ActionResult struct {
	Artifact string `json:"artifact"`
	Output   string `json:"output"`
}

type indexData struct {
	Files    map[string]FileRecord   `json:"files"`
	GoAPI    map[string]FileRecord   `json:"go_api"`
	Actions  map[string]ActionResult `json:"actions"`
	Receipts map[string]time.Time    `json:"receipts"`
}

type Cache struct {
	root         string
	stateDir     string
	indexPath    string
	artifactRoot string
	index        indexData
	stats        SnapshotStats
}

type SnapshotStats struct {
	Files         int
	BytesRead     int64
	DigestsReused int
}

type SnapshotOptions struct {
	Label             string
	Root              string
	IncludeAll        bool
	IncludeNames      []string
	IncludeExtensions []string
	ExcludeDirs       []string
}

type LookupStatus string

const (
	LookupMiss     LookupStatus = "miss"
	LookupHit      LookupStatus = "hit"
	LookupRestored LookupStatus = "restored"
	LookupDirty    LookupStatus = "dirty"
)

func OpenCache(root string) (*Cache, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	stateDir := filepath.Join(absRoot, ".wails")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return nil, fmt.Errorf("wake: create state directory: %w", err)
	}
	userCache, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("wake: locate user cache: %w", err)
	}
	c := &Cache{
		root:         absRoot,
		stateDir:     stateDir,
		indexPath:    filepath.Join(stateDir, "wake-index.json"),
		artifactRoot: filepath.Join(userCache, "wails", "wake", "artifacts"),
		index: indexData{
			Files:    make(map[string]FileRecord),
			GoAPI:    make(map[string]FileRecord),
			Actions:  make(map[string]ActionResult),
			Receipts: make(map[string]time.Time),
		},
	}
	data, err := os.ReadFile(c.indexPath)
	if err == nil {
		if decodeErr := json.Unmarshal(data, &c.index); decodeErr != nil {
			// Cache state is disposable. A truncated or corrupt index is a cold
			// cache, not a build failure that requires manual repair.
			c.index = indexData{}
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("wake: read action index: %w", err)
	}
	if c.index.Files == nil {
		c.index.Files = make(map[string]FileRecord)
	}
	if c.index.Actions == nil {
		c.index.Actions = make(map[string]ActionResult)
	}
	if c.index.GoAPI == nil {
		c.index.GoAPI = make(map[string]FileRecord)
	}
	if c.index.Receipts == nil {
		c.index.Receipts = make(map[string]time.Time)
	}
	return c, nil
}

func (c *Cache) Save() error {
	data, err := json.Marshal(c.index)
	if err != nil {
		return fmt.Errorf("wake: encode action index: %w", err)
	}
	tmp, err := os.CreateTemp(c.stateDir, ".wake-index-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, c.indexPath)
}

func (c *Cache) Stats() SnapshotStats { return c.stats }

func (c *Cache) ResetStats() { c.stats = SnapshotStats{} }

func (c *Cache) Snapshot(options SnapshotOptions) (string, error) {
	root := options.Root
	if !filepath.IsAbs(root) {
		root = filepath.Join(c.root, root)
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	if _, statErr := os.Stat(root); errors.Is(statErr, fs.ErrNotExist) {
		return "", fmt.Errorf("wake: snapshot root %q does not exist", root)
	} else if statErr != nil {
		return "", fmt.Errorf("wake: inspect snapshot root %q: %w", root, statErr)
	}
	includeNames := makeSet(options.IncludeNames)
	includeExts := makeSet(options.IncludeExtensions)
	excludeDirs := makeSet(options.ExcludeDirs)

	type entry struct {
		Path   string
		Mode   fs.FileMode
		Digest string
	}
	var entries []entry
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." && d.IsDir() {
			return nil
		}
		if rel == "." {
			rel = filepath.Base(root)
		}
		if d.IsDir() {
			if excludeDirs[d.Name()] || excludeDirs[rel] {
				return filepath.SkipDir
			}
			return nil
		}
		if !options.IncludeAll && !includeNames[d.Name()] && !includeExts[strings.ToLower(filepath.Ext(d.Name()))] {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		digest, err := c.fileDigest(path, info)
		if err != nil {
			return err
		}
		entries = append(entries, entry{Path: rel, Mode: info.Mode(), Digest: digest})
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	h := blake3.New()
	writePart(h, "wake.snapshot.v1")
	writePart(h, options.Label)
	for _, e := range entries {
		writePart(h, e.Path)
		writePart(h, fmt.Sprintf("%o", relevantMode(e.Mode)))
		writePart(h, e.Digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c *Cache) SnapshotFiles(label string, paths ...string) (string, error) {
	type entry struct{ path, digest string }
	entries := make([]entry, 0, len(paths))
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(c.root, path)
		}
		info, err := os.Lstat(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		digest, err := c.fileDigest(path, info)
		if err != nil {
			return "", err
		}
		logical, err := filepath.Rel(c.root, path)
		if err != nil || strings.HasPrefix(logical, "..") {
			logical = "external/" + filepath.Base(path)
		}
		entries = append(entries, entry{filepath.ToSlash(logical), digest})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].path < entries[j].path })
	h := blake3.New()
	writePart(h, "wake.snapshot-files.v1")
	writePart(h, label)
	for _, e := range entries {
		writePart(h, e.path)
		writePart(h, e.digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SnapshotGoAPI fingerprints binding-relevant Go syntax while ignoring method
// bodies. Free-function bodies remain inputs because main may delegate Wails
// service registration to helpers. Per-file semantic digests use the same
// metadata fast path as content Snapshots, so an unchanged tree is not reparsed.
func (c *Cache) SnapshotGoAPI(label, root string, exclude []string) (string, error) {
	if !filepath.IsAbs(root) {
		root = filepath.Join(c.root, root)
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	excluded := makeSet(exclude)
	type entry struct{ path, digest string }
	var entries []entry
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel != "." && (excluded[d.Name()] || excluded[rel]) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		digest, err := c.goAPIDigest(path, info)
		if err != nil {
			return err
		}
		entries = append(entries, entry{rel, digest})
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].path < entries[j].path })
	h := blake3.New()
	writePart(h, "wake.go-api.v1")
	writePart(h, label)
	for _, entry := range entries {
		writePart(h, entry.path)
		writePart(h, entry.digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c *Cache) goAPIDigest(path string, info fs.FileInfo) (string, error) {
	c.stats.Files++
	identity := fileIdentity(info)
	key := filepath.Clean(path)
	record, ok := c.index.GoAPI[key]
	if metadataFastPathSafe(info) && ok && record.Size == info.Size() && record.ModTimeNS == info.ModTime().UnixNano() && record.Mode == uint32(info.Mode()) && record.Identity == identity {
		c.stats.DigestsReused++
		return record.Digest, nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("wake: parse Go API %s: %w", path, err)
	}
	semantic := &ast.File{Name: file.Name}
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			copy := *fn
			if fn.Recv != nil {
				copy.Body = nil
			}
			semantic.Decls = append(semantic.Decls, &copy)
			continue
		}
		semantic.Decls = append(semantic.Decls, decl)
	}
	var data bytes.Buffer
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if strings.HasPrefix(comment.Text, "//go:build") || strings.HasPrefix(comment.Text, "// +build") {
				data.WriteString(comment.Text)
				data.WriteByte('\n')
			}
		}
	}
	if err := format.Node(&data, fset, semantic); err != nil {
		return "", err
	}
	sum := blake3.Sum256(data.Bytes())
	digest := hex.EncodeToString(sum[:])
	c.stats.BytesRead += int64(data.Len())
	c.index.GoAPI[key] = FileRecord{Size: info.Size(), ModTimeNS: info.ModTime().UnixNano(), Mode: uint32(info.Mode()), Identity: identity, Digest: digest}
	return digest, nil
}

func ActionKey(kind string, spec any, inputs, dependencies []string) (string, error) {
	payload := struct {
		Domain       string   `json:"domain"`
		Handler      string   `json:"handler"`
		Kind         string   `json:"kind"`
		Spec         any      `json:"spec"`
		Inputs       []string `json:"inputs"`
		Dependencies []string `json:"dependencies"`
	}{"wake.action.v1", handlerVersion, kind, spec, inputs, dependencies}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := blake3.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func (c *Cache) Lookup(actionKey, output string) (LookupStatus, string, error) {
	result, ok := c.index.Actions[actionKey]
	if !ok || result.Output != filepath.ToSlash(output) {
		return LookupMiss, "", nil
	}
	absOutput := filepath.Join(c.root, filepath.FromSlash(output))
	_, err := os.Lstat(absOutput)
	if errors.Is(err, fs.ErrNotExist) {
		if err := c.restoreArtifact(result.Artifact, absOutput); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return LookupMiss, "", nil
			}
			return LookupMiss, "", fmt.Errorf("wake: restore artifact %s: %w", result.Artifact, err)
		}
		return LookupRestored, result.Artifact, nil
	}
	if err != nil {
		return LookupMiss, "", err
	}
	digest, err := c.Snapshot(SnapshotOptions{Label: "artifact", Root: absOutput, IncludeAll: true})
	if err != nil {
		return LookupMiss, "", err
	}
	if digest == result.Artifact {
		return LookupHit, digest, nil
	}
	return LookupDirty, "", nil
}

func (c *Cache) RecordAction(actionKey, output string) (string, error) {
	absOutput := filepath.Join(c.root, filepath.FromSlash(output))
	digest, err := c.Snapshot(SnapshotOptions{Label: "artifact", Root: absOutput, IncludeAll: true})
	if err != nil {
		return "", err
	}
	if err := c.storeArtifact(digest, absOutput); err != nil {
		return "", err
	}
	c.index.Actions[actionKey] = ActionResult{Artifact: digest, Output: filepath.ToSlash(output)}
	if err := c.Save(); err != nil {
		return "", err
	}
	return digest, nil
}

func (c *Cache) HasReceipt(actionKey, marker string) bool {
	if _, ok := c.index.Receipts[actionKey]; !ok {
		return false
	}
	_, err := os.Stat(filepath.Join(c.root, marker))
	return err == nil
}

func (c *Cache) RecordReceipt(actionKey string) error {
	c.index.Receipts[actionKey] = time.Now()
	return c.Save()
}

func (c *Cache) fileDigest(path string, info fs.FileInfo) (string, error) {
	c.stats.Files++
	identity := fileIdentity(info)
	key := filepath.Clean(path)
	record, ok := c.index.Files[key]
	if metadataFastPathSafe(info) && ok && record.Size == info.Size() && record.ModTimeNS == info.ModTime().UnixNano() &&
		record.Mode == uint32(info.Mode()) && record.Identity == identity {
		c.stats.DigestsReused++
		return record.Digest, nil
	}

	var digest string
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		sum := blake3.Sum256([]byte("symlink\x00" + target))
		digest = hex.EncodeToString(sum[:])
		c.stats.BytesRead += int64(len(target))
	} else {
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		h := blake3.New()
		n, copyErr := io.Copy(h, file)
		closeErr := file.Close()
		if copyErr != nil {
			return "", copyErr
		}
		if closeErr != nil {
			return "", closeErr
		}
		c.stats.BytesRead += n
		digest = hex.EncodeToString(h.Sum(nil))
	}
	c.index.Files[key] = FileRecord{
		Size: info.Size(), ModTimeNS: info.ModTime().UnixNano(), Mode: uint32(info.Mode()),
		Identity: identity, Digest: digest,
	}
	return digest, nil
}

// Files modified very recently are re-read even when their metadata matches
// the index. This closes the correctness hole on filesystems with coarse
// timestamp resolution, while retaining the fast path for stable source trees.
func metadataFastPathSafe(info fs.FileInfo) bool {
	return time.Since(info.ModTime()) >= metadataFastPathWindow
}

func fileIdentity(info fs.FileInfo) string {
	v := reflect.ValueOf(info.Sys())
	if !v.IsValid() {
		return ""
	}
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return fmt.Sprintf("%T", info.Sys())
	}
	parts := []string{fmt.Sprintf("%T", info.Sys())}
	for _, name := range []string{"Dev", "Ino", "FileIndexHigh", "FileIndexLow", "Ctim", "Ctimespec", "ChangeTime"} {
		field := v.FieldByName(name)
		if field.IsValid() && field.CanInterface() {
			parts = append(parts, name+"="+fmt.Sprint(field.Interface()))
		}
	}
	return strings.Join(parts, ";")
}

func (c *Cache) storeArtifact(digest, source string) error {
	destination := filepath.Join(c.artifactRoot, digest)
	if _, err := os.Stat(filepath.Join(destination, "payload")); err == nil {
		return nil
	}
	if err := os.MkdirAll(c.artifactRoot, 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(c.artifactRoot, ".artifact-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := copyPath(source, filepath.Join(tmp, "payload")); err != nil {
		return err
	}
	if err := os.Rename(tmp, destination); err != nil && !errors.Is(err, fs.ErrExist) {
		if _, statErr := os.Stat(filepath.Join(destination, "payload")); statErr != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) restoreArtifact(digest, destination string) error {
	source := filepath.Join(c.artifactRoot, digest, "payload")
	if _, err := os.Stat(source); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(filepath.Dir(destination), ".wake-restore-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	payload := filepath.Join(tmp, "payload")
	if err := copyPath(source, payload); err != nil {
		return err
	}
	return os.Rename(payload, destination)
}

func copyPath(source, destination string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}
		return os.Symlink(target, destination)
	}
	if !info.IsDir() {
		return copyFile(source, destination, info.Mode())
	}
	if err := os.MkdirAll(destination, info.Mode().Perm()); err != nil {
		return err
	}
	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := copyPath(filepath.Join(source, entry.Name()), filepath.Join(destination, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(source, destination string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func writePart(w io.Writer, value string) {
	_, _ = io.WriteString(w, fmt.Sprintf("%d:", len(value)))
	_, _ = io.WriteString(w, value)
}

func makeSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}

func relevantMode(mode fs.FileMode) fs.FileMode {
	return mode & (0o111 | os.ModeType)
}
