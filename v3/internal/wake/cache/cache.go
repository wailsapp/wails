package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
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
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charlievieth/fastwalk"
	"github.com/go-git/go-billy/v5/osfs"
	gitignore "github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/zeebo/blake3"
)

const (
	handlerVersion                       = "wake-pipeline-2"
	metadataFastPathWindow time.Duration = 2_000_000_000
)

type FileRecord struct {
	Size      int64  `json:"size"`
	ModTimeNS int64  `json:"mod_time_ns"`
	Mode      uint32 `json:"mode"`
	Identity  string `json:"identity"`
	Digest    string `json:"digest"`
}

type ActionResult struct {
	Artifact       string     `json:"artifact"`
	Output         string     `json:"output"`
	OutputMetadata FileRecord `json:"output_metadata,omitempty"`
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
	legacyIndex  string
	artifactRoot string
	index        indexData
	stats        SnapshotStats
	dirty        bool
	observations map[string]*treeObservation
	shareTrees   bool
}

type SnapshotStats struct {
	Files         int
	BytesRead     int64
	DigestsReused int
	TreesWalked   int
}

type treeObservation struct {
	files []*observedFile
}

type observedFile struct {
	path, relative string
	entry          fs.DirEntry
	info           fs.FileInfo
	infoErr        error
	identity       string
	identityReady  bool
}

func (f *observedFile) fileInfo() (fs.FileInfo, error) {
	if f.info == nil && f.infoErr == nil {
		f.info, f.infoErr = f.entry.Info()
	}
	return f.info, f.infoErr
}

func (f *observedFile) fileIdentity(info fs.FileInfo) string {
	if !f.identityReady {
		f.identity = fileIdentity(info)
		f.identityReady = true
	}
	return f.identity
}

type SnapshotOptions struct {
	Label             string
	Root              string
	IncludeAll        bool
	IncludeNames      []string
	IncludeExtensions []string
	ExcludeDirs       []string
	ExcludeSuffixes   []string
	UseGitIgnore      bool
}

type LookupStatus string

const (
	LookupMiss     LookupStatus = "miss"
	LookupHit      LookupStatus = "hit"
	LookupRestored LookupStatus = "restored"
	LookupDirty    LookupStatus = "dirty"
)

func OpenCache(root string) (*Cache, error) {
	return openCacheWithOperations(root, cacheOpenOperations{
		abs:          filepath.Abs,
		mkdirAll:     os.MkdirAll,
		userCacheDir: os.UserCacheDir,
		readFile:     os.ReadFile,
	})
}

// OpenCacheReadOnly opens existing cache metadata without creating project
// state. Snapshot updates remain in memory and are never saved by inspection.
func OpenCacheReadOnly(root string) (*Cache, error) {
	return openCacheWithOperations(root, cacheOpenOperations{
		abs:          filepath.Abs,
		mkdirAll:     func(string, fs.FileMode) error { return nil },
		userCacheDir: os.UserCacheDir,
		readFile:     os.ReadFile,
	})
}

type cacheOpenOperations struct {
	abs          func(string) (string, error)
	mkdirAll     func(string, fs.FileMode) error
	userCacheDir func() (string, error)
	readFile     func(string) ([]byte, error)
}

func openCacheWithOperations(root string, operations cacheOpenOperations) (*Cache, error) {
	absRoot, err := operations.abs(root)
	if err != nil {
		return nil, err
	}
	stateDir := filepath.Join(absRoot, ".wails")
	if err := operations.mkdirAll(stateDir, 0o755); err != nil {
		return nil, fmt.Errorf("wake: create state directory: %w", err)
	}
	userCache, err := operations.userCacheDir()
	if err != nil {
		return nil, fmt.Errorf("wake: locate user cache: %w", err)
	}
	c := &Cache{
		root:         absRoot,
		stateDir:     stateDir,
		indexPath:    filepath.Join(stateDir, "wake-index-v2.gob"),
		legacyIndex:  filepath.Join(stateDir, "wake-index.json"),
		artifactRoot: filepath.Join(userCache, "wails", "wake", "artifacts"),
		index: indexData{
			Files:    make(map[string]FileRecord),
			GoAPI:    make(map[string]FileRecord),
			Actions:  make(map[string]ActionResult),
			Receipts: make(map[string]time.Time),
		},
		observations: make(map[string]*treeObservation),
	}
	data, err := operations.readFile(c.indexPath)
	if err == nil {
		if decodeErr := gob.NewDecoder(bytes.NewReader(data)).Decode(&c.index); decodeErr != nil {
			// Cache state is disposable. A truncated or corrupt index is a cold
			// cache, not a build failure that requires manual repair.
			c.index = indexData{}
		}
	} else if errors.Is(err, fs.ErrNotExist) {
		legacy, legacyErr := operations.readFile(c.legacyIndex)
		if legacyErr == nil {
			if decodeErr := json.Unmarshal(legacy, &c.index); decodeErr == nil {
				c.dirty = true
			} else {
				c.index = indexData{}
			}
		} else if !errors.Is(legacyErr, fs.ErrNotExist) {
			return nil, fmt.Errorf("wake: read legacy action index: %w", legacyErr)
		}
	} else {
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
	return c.saveWithOperations(cacheSaveOperations{
		encode: func(index indexData) ([]byte, error) {
			var data bytes.Buffer
			err := gob.NewEncoder(&data).Encode(index)
			return data.Bytes(), err
		},
		createTemp: func(dir, pattern string) (cacheTemporaryFile, error) { return os.CreateTemp(dir, pattern) },
		remove:     os.Remove,
		rename:     os.Rename,
	})
}

type cacheTemporaryFile interface {
	Name() string
	Write([]byte) (int, error)
	Sync() error
	Close() error
}

type cacheSaveOperations struct {
	encode     func(indexData) ([]byte, error)
	createTemp func(string, string) (cacheTemporaryFile, error)
	remove     func(string) error
	rename     func(string, string) error
}

func (c *Cache) saveWithOperations(operations cacheSaveOperations) error {
	if !c.dirty {
		return nil
	}
	data, err := operations.encode(c.index)
	if err != nil {
		return fmt.Errorf("wake: encode action index: %w", err)
	}
	tmp, err := operations.createTemp(c.stateDir, ".wake-index-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer operations.remove(name)
	if written, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	} else if written != len(data) {
		tmp.Close()
		return io.ErrShortWrite
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := operations.rename(name, c.indexPath); err != nil {
		return err
	}
	c.dirty = false
	return nil
}

func (c *Cache) Stats() SnapshotStats { return c.stats }

func (c *Cache) ResetStats() { c.stats = SnapshotStats{} }

// BeginObservationSession lets a finite executor share tree traversal across
// related snapshots. Direct Cache callers retain fresh-per-call semantics.
func (c *Cache) BeginObservationSession() {
	c.shareTrees = true
	c.InvalidateObservations()
}

// InvalidateObservations starts a new filesystem-observation boundary. The
// executor calls this before running a handler, so snapshots within an
// unchanged no-op graph share one consistent tree view while generated or
// user-modified files are always rediscovered after actual work.
func (c *Cache) InvalidateObservations() {
	c.observations = make(map[string]*treeObservation)
}

func (c *Cache) observeTree(root string, options SnapshotOptions, semanticGo bool) (*treeObservation, error) {
	return c.observeTreeWithOperations(root, options, semanticGo, cacheTreeOperations{
		stat:    os.Stat,
		lstat:   os.Lstat,
		ignored: readGitignoreMatcher,
		walk:    observeDirectoryTree,
	})
}

func readGitignoreMatcher(root string) (gitignore.Matcher, error) {
	patterns, err := gitignore.ReadPatterns(osfs.New(root), nil)
	if err != nil {
		return nil, err
	}
	return gitignore.NewMatcher(patterns), nil
}

type cacheTreeOperations struct {
	stat    func(string) (fs.FileInfo, error)
	lstat   func(string) (fs.FileInfo, error)
	ignored func(string) (gitignore.Matcher, error)
	walk    func(string, map[string]bool, gitignore.Matcher, func(string) bool) ([]*observedFile, error)
}

func (c *Cache) observeTreeWithOperations(root string, options SnapshotOptions, semanticGo bool, operations cacheTreeOperations) (*treeObservation, error) {
	excludedValues := append([]string(nil), options.ExcludeDirs...)
	for index := range excludedValues {
		excludedValues[index] = filepath.ToSlash(filepath.Clean(excludedValues[index]))
	}
	sort.Strings(excludedValues)
	options.ExcludeDirs = excludedValues
	key := treeObservationKey(root, options)
	if c.shareTrees {
		if observation := c.observations[key]; observation != nil {
			return observation, nil
		}
	}
	if _, err := operations.stat(root); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("wake: snapshot root %q does not exist", root)
	} else if err != nil {
		return nil, fmt.Errorf("wake: inspect snapshot root %q: %w", root, err)
	}
	rootInfo, err := operations.lstat(root)
	if err != nil {
		return nil, err
	}
	observation := &treeObservation{}
	if !rootInfo.IsDir() {
		observation.files = []*observedFile{{
			path: root, relative: filepath.Base(root), entry: fs.FileInfoToDirEntry(rootInfo), info: rootInfo,
		}}
		if c.shareTrees {
			c.observations[key] = observation
		}
		c.stats.TreesWalked++
		return observation, nil
	}

	excluded := makeSet(excludedValues)
	var ignored gitignore.Matcher
	if options.UseGitIgnore {
		ignored, err = operations.ignored(root)
		if err != nil {
			return nil, fmt.Errorf("wake: read snapshot gitignore: %w", err)
		}
	}
	includeNames := makeSet(options.IncludeNames)
	includeExts := makeSet(options.IncludeExtensions)
	includeFile := func(name string) bool {
		semanticInput := semanticGo && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
		contentInput := options.IncludeAll || includeNames[name] || includeExts[strings.ToLower(filepath.Ext(name))]
		for _, suffix := range options.ExcludeSuffixes {
			if strings.HasSuffix(name, suffix) {
				contentInput = false
				break
			}
		}
		return semanticInput || contentInput
	}
	files, err := operations.walk(root, excluded, ignored, includeFile)
	if err != nil {
		return nil, err
	}
	observation.files = files
	if c.shareTrees {
		c.observations[key] = observation
	}
	c.stats.TreesWalked++
	return observation, nil
}

func treeObservationKey(root string, options SnapshotOptions) string {
	sorted := func(values []string) string {
		values = append([]string(nil), values...)
		sort.Strings(values)
		return strings.Join(values, "\x00")
	}
	return root + "\x00dirs=" + sorted(options.ExcludeDirs) + "\x00names=" + sorted(options.IncludeNames) +
		"\x00extensions=" + sorted(options.IncludeExtensions) + "\x00suffixes=" + sorted(options.ExcludeSuffixes) +
		"\x00all=" + strconv.FormatBool(options.IncludeAll) + "\x00gitignore=" + strconv.FormatBool(options.UseGitIgnore)
}

func observeDirectoryTree(root string, excluded map[string]bool, ignored gitignore.Matcher, eagerInfo func(string) bool) ([]*observedFile, error) {
	workers := fastwalk.DefaultNumWorkers()
	if runtime.GOOS != "windows" {
		workers = min(workers, 8)
	}
	return observeDirectoryTreeWithWorkers(root, excluded, ignored, eagerInfo, workers)
}

func observeDirectoryTreeWithWorkers(root string, excluded map[string]bool, ignored gitignore.Matcher, eagerInfo func(string) bool, workers int) ([]*observedFile, error) {
	return observeDirectoryTreeWithWalk(root, excluded, ignored, eagerInfo, workers, fastwalk.Walk)
}

type cacheWalker func(*fastwalk.Config, string, fs.WalkDirFunc) error

func observeDirectoryTreeWithWalk(root string, excluded map[string]bool, ignored gitignore.Matcher, eagerInfo func(string) bool, workers int, walk cacheWalker) ([]*observedFile, error) {
	var filesMu sync.Mutex
	var files []*observedFile
	config := fastwalk.DefaultConfig.Copy()
	config.NumWorkers = workers
	prefix := root + string(filepath.Separator)
	err := walk(config, root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		relative := filepath.ToSlash(strings.TrimPrefix(path, prefix))
		if ignored != nil && ignored.Match(strings.Split(relative, "/"), entry.IsDir()) {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			if excluded[entry.Name()] || excluded[relative] {
				return fs.SkipDir
			}
			return nil
		}
		if !eagerInfo(entry.Name()) {
			return nil
		}
		observed := &observedFile{path: path, relative: relative, entry: entry}
		observed.info, observed.infoErr = entry.Info()
		if observed.infoErr != nil {
			return observed.infoErr
		}
		filesMu.Lock()
		files = append(files, observed)
		filesMu.Unlock()
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].relative < files[j].relative })
	return files, nil
}

func (c *Cache) Snapshot(options SnapshotOptions) (string, error) {
	root := options.Root
	if !filepath.IsAbs(root) {
		root = filepath.Join(c.root, root)
	}
	root = filepath.Clean(root)
	observation, err := c.observeTree(root, options, false)
	if err != nil {
		return "", err
	}
	h := blake3.New()
	writePart(h, "wake.snapshot.v1")
	writePart(h, options.Label)
	for _, file := range observation.files {
		info, err := file.fileInfo()
		if err != nil {
			return "", err
		}
		digest, err := c.fileDigestWithIdentity(file.path, info, file.fileIdentity(info))
		if err != nil {
			return "", err
		}
		writePart(h, file.relative)
		writePart(h, strconv.FormatUint(uint64(relevantMode(info.Mode())), 8))
		writePart(h, digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c *Cache) SnapshotFiles(label string, paths ...string) (string, error) {
	return c.snapshotFilesWithOperations(label, paths, cacheSnapshotFileOperations{
		lstat: os.Lstat,
		digest: func(path string, info fs.FileInfo) (string, error) {
			if info.IsDir() {
				return artifactDigest(path)
			}
			return c.fileDigest(path, info)
		},
		rel: filepath.Rel,
	})
}

type cacheSnapshotFileOperations struct {
	lstat  func(string) (fs.FileInfo, error)
	digest func(string, fs.FileInfo) (string, error)
	rel    func(string, string) (string, error)
}

func (c *Cache) snapshotFilesWithOperations(label string, paths []string, operations cacheSnapshotFileOperations) (string, error) {
	type entry struct {
		path, digest string
		mode         fs.FileMode
	}
	entries := make([]entry, 0, len(paths))
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(c.root, path)
		}
		info, err := operations.lstat(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		digest, err := operations.digest(path, info)
		if err != nil {
			return "", err
		}
		logical, err := operations.rel(c.root, path)
		if err != nil || strings.HasPrefix(logical, "..") {
			logical = "external/" + filepath.Base(path)
		}
		entries = append(entries, entry{path: filepath.ToSlash(logical), digest: digest, mode: info.Mode()})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].path < entries[j].path })
	h := blake3.New()
	writePart(h, "wake.snapshot-files.v1")
	writePart(h, label)
	for _, e := range entries {
		writePart(h, e.path)
		writePart(h, strconv.FormatUint(uint64(relevantMode(e.mode)), 8))
		writePart(h, e.digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SnapshotGoAPI fingerprints binding-relevant Go syntax while ignoring method
// bodies. Free-function bodies remain inputs because main may delegate Wails
// service registration to helpers. Per-file semantic digests use the same
// metadata fast path as content Snapshots, so an unchanged tree is not reparsed.
func (c *Cache) SnapshotGoAPI(options SnapshotOptions) (string, error) {
	root := options.Root
	if !filepath.IsAbs(root) {
		root = filepath.Join(c.root, root)
	}
	root = filepath.Clean(root)
	observation, err := c.observeTree(root, options, true)
	if err != nil {
		return "", err
	}
	h := blake3.New()
	writePart(h, "wake.go-api.v1")
	writePart(h, options.Label)
	for _, observed := range observation.files {
		if !strings.HasSuffix(observed.entry.Name(), ".go") || strings.HasSuffix(observed.entry.Name(), "_test.go") {
			continue
		}
		info, err := observed.fileInfo()
		if err != nil {
			return "", err
		}
		digest, err := c.goAPIDigestWithIdentity(observed.path, info, observed.fileIdentity(info))
		if err != nil {
			return "", err
		}
		writePart(h, observed.relative)
		writePart(h, digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c *Cache) goAPIDigestWithIdentity(path string, info fs.FileInfo, identity string) (string, error) {
	return c.goAPIDigestWithFormatter(path, info, identity, format.Node)
}

type cacheFormatNode func(io.Writer, *token.FileSet, any) error

func (c *Cache) goAPIDigestWithFormatter(path string, info fs.FileInfo, identity string, formatNode cacheFormatNode) (string, error) {
	c.stats.Files++
	key := path
	record, ok := c.index.GoAPI[key]
	if metadataFastPathSafe(info, identity) && ok && record.Size == info.Size() && record.ModTimeNS == info.ModTime().UnixNano() && record.Mode == uint32(info.Mode()) && record.Identity == identity {
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
	if err := formatNode(&data, fset, semantic); err != nil {
		return "", err
	}
	sum := blake3.Sum256(data.Bytes())
	digest := hex.EncodeToString(sum[:])
	c.stats.BytesRead += int64(data.Len())
	c.index.GoAPI[key] = FileRecord{Size: info.Size(), ModTimeNS: info.ModTime().UnixNano(), Mode: uint32(info.Mode()), Identity: identity, Digest: digest}
	c.dirty = true
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
	return c.lookupWithOperations(actionKey, output, cacheLookupOperations{
		lstat: os.Lstat,
		restore: func(digest, destination string) error {
			return c.restoreArtifact(digest, destination)
		},
		removeAll: os.RemoveAll,
		save:      c.Save,
		digest:    artifactDigest,
	})
}

type cacheLookupOperations struct {
	lstat     func(string) (fs.FileInfo, error)
	restore   func(string, string) error
	removeAll func(string) error
	save      func() error
	digest    func(string) (string, error)
}

func (c *Cache) lookupWithOperations(actionKey, output string, operations cacheLookupOperations) (LookupStatus, string, error) {
	result, ok := c.index.Actions[actionKey]
	if !ok || result.Output != filepath.ToSlash(output) {
		return LookupMiss, "", nil
	}
	absOutput := filepath.Join(c.root, filepath.FromSlash(output))
	info, err := operations.lstat(absOutput)
	if errors.Is(err, fs.ErrNotExist) {
		if err := operations.restore(result.Artifact, absOutput); err != nil {
			if errors.Is(err, errCorruptArtifact) {
				_ = operations.removeAll(filepath.Join(c.artifactRoot, result.Artifact))
				delete(c.index.Actions, actionKey)
				c.dirty = true
				_ = operations.save()
				return LookupMiss, "", nil
			}
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
	if actionOutputMetadataMatches(result, info) {
		return LookupHit, result.Artifact, nil
	}
	digest, err := operations.digest(absOutput)
	if err != nil {
		return LookupMiss, "", err
	}
	if digest == result.Artifact {
		return LookupHit, digest, nil
	}
	return LookupDirty, "", nil
}

// Peek reports what Lookup would do without restoring or modifying files.
func (c *Cache) Peek(actionKey, output string) (LookupStatus, string, error) {
	return c.peekWithOperations(actionKey, output, cachePeekOperations{
		lstat:  os.Lstat,
		stat:   os.Stat,
		digest: artifactDigest,
	})
}

type cachePeekOperations struct {
	lstat  func(string) (fs.FileInfo, error)
	stat   func(string) (fs.FileInfo, error)
	digest func(string) (string, error)
}

func (c *Cache) peekWithOperations(actionKey, output string, operations cachePeekOperations) (LookupStatus, string, error) {
	result, ok := c.index.Actions[actionKey]
	if !ok || result.Output != filepath.ToSlash(output) {
		return LookupMiss, "", nil
	}
	absOutput := filepath.Join(c.root, filepath.FromSlash(output))
	info, err := operations.lstat(absOutput)
	if errors.Is(err, fs.ErrNotExist) {
		payload := filepath.Join(c.artifactRoot, result.Artifact, "payload")
		if _, artifactErr := operations.stat(payload); artifactErr == nil {
			if digest, digestErr := operations.digest(payload); digestErr != nil || digest != result.Artifact {
				return LookupMiss, "", nil
			}
			return LookupRestored, result.Artifact, nil
		} else if errors.Is(artifactErr, fs.ErrNotExist) {
			return LookupMiss, "", nil
		} else {
			return LookupMiss, "", artifactErr
		}
	}
	if err != nil {
		return LookupMiss, "", err
	}
	if actionOutputMetadataMatches(result, info) {
		return LookupHit, result.Artifact, nil
	}
	digest, err := operations.digest(absOutput)
	if err != nil {
		return LookupMiss, "", err
	}
	if digest == result.Artifact {
		return LookupHit, digest, nil
	}
	return LookupDirty, "", nil
}

func (c *Cache) RecordAction(actionKey, output string) (string, error) {
	return c.recordActionWithOperations(actionKey, output, cacheRecordOperations{
		digest: artifactDigest,
		store:  c.storeArtifact,
		lstat:  os.Lstat,
		save:   c.Save,
	})
}

type cacheRecordOperations struct {
	digest func(string) (string, error)
	store  func(string, string) error
	lstat  func(string) (fs.FileInfo, error)
	save   func() error
}

func (c *Cache) recordActionWithOperations(actionKey, output string, operations cacheRecordOperations) (string, error) {
	absOutput := filepath.Join(c.root, filepath.FromSlash(output))
	before, err := operations.lstat(absOutput)
	if err != nil {
		return "", err
	}
	digest, err := operations.digest(absOutput)
	if err != nil {
		return "", err
	}
	if err := operations.store(digest, absOutput); err != nil {
		return "", err
	}
	after, err := operations.lstat(absOutput)
	if err != nil {
		return "", err
	}
	result := ActionResult{Artifact: digest, Output: filepath.ToSlash(output)}
	if after.Mode().IsRegular() && sameFileMetadata(before, after) {
		result.OutputMetadata = fileRecord(after, digest)
	}
	c.index.Actions[actionKey] = result
	c.dirty = true
	if err := operations.save(); err != nil {
		return "", err
	}
	return digest, nil
}

func sameFileMetadata(left, right fs.FileInfo) bool {
	return left.Size() == right.Size() && left.ModTime().UnixNano() == right.ModTime().UnixNano() &&
		left.Mode() == right.Mode() && fileIdentity(left) == fileIdentity(right)
}

func actionOutputMetadataMatches(result ActionResult, info fs.FileInfo) bool {
	record := result.OutputMetadata
	if !info.Mode().IsRegular() || record.Digest == "" || record.Digest != result.Artifact {
		return false
	}
	identity := fileIdentity(info)
	// Final artifacts are user-visible outputs. Trust metadata without hashing
	// only when the platform identity includes an OS-maintained content-change
	// counter (ctime on Linux and macOS). Age-based mtime reuse is adequate for
	// source snapshots, but is not a strong enough integrity check here.
	return identity != "" && platformIdentityTracksChanges() && record.Size == info.Size() && record.ModTimeNS == info.ModTime().UnixNano() &&
		record.Mode == uint32(info.Mode()) && record.Identity == identity
}

func fileRecord(info fs.FileInfo, digest string) FileRecord {
	return FileRecord{
		Size: info.Size(), ModTimeNS: info.ModTime().UnixNano(), Mode: uint32(info.Mode()),
		Identity: fileIdentity(info), Digest: digest,
	}
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
	c.dirty = true
	return c.Save()
}

func (c *Cache) fileDigest(path string, info fs.FileInfo) (string, error) {
	return c.fileDigestWithIdentity(path, info, fileIdentity(info))
}

func (c *Cache) fileDigestWithIdentity(path string, info fs.FileInfo, identity string) (string, error) {
	return c.fileDigestWithOperations(path, info, identity, cacheDigestOperations{
		readlink: os.Readlink,
		open:     func(path string) (cacheDigestFile, error) { return os.Open(path) },
	})
}

type cacheDigestFile interface {
	io.Reader
	Close() error
}

type cacheDigestOperations struct {
	readlink func(string) (string, error)
	open     func(string) (cacheDigestFile, error)
}

func (c *Cache) fileDigestWithOperations(path string, info fs.FileInfo, identity string, operations cacheDigestOperations) (string, error) {
	c.stats.Files++
	key := path
	record, ok := c.index.Files[key]
	if metadataFastPathSafe(info, identity) && ok && record.Size == info.Size() && record.ModTimeNS == info.ModTime().UnixNano() &&
		record.Mode == uint32(info.Mode()) && record.Identity == identity {
		c.stats.DigestsReused++
		return record.Digest, nil
	}

	var digest string
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := operations.readlink(path)
		if err != nil {
			return "", err
		}
		sum := blake3.Sum256([]byte("symlink\x00" + target))
		digest = hex.EncodeToString(sum[:])
		c.stats.BytesRead += int64(len(target))
	} else {
		file, err := operations.open(path)
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
	c.dirty = true
	return digest, nil
}

// Files modified very recently are re-read even when their metadata matches
// the index. This closes the correctness hole on filesystems with coarse
// timestamp resolution, while retaining the fast path for stable source trees.
func metadataFastPathSafe(info fs.FileInfo, identity string) bool {
	return (identity != "" && platformIdentityTracksChanges()) || time.Since(info.ModTime()) >= metadataFastPathWindow
}

func fileIdentity(info fs.FileInfo) string {
	if identity, ok := platformFileIdentity(info); ok {
		return identity
	}
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

func encodeFileIdentity(first, second, third, fourth uint64) string {
	var raw [32]byte
	binary.LittleEndian.PutUint64(raw[0:8], first)
	binary.LittleEndian.PutUint64(raw[8:16], second)
	binary.LittleEndian.PutUint64(raw[16:24], third)
	binary.LittleEndian.PutUint64(raw[24:32], fourth)
	return hex.EncodeToString(raw[:])
}

func (c *Cache) storeArtifact(digest, source string) error {
	return c.storeArtifactWithOperations(digest, source, cacheStoreArtifactOperations{
		stat:      os.Stat,
		digest:    artifactDigest,
		removeAll: os.RemoveAll,
		mkdirAll:  os.MkdirAll,
		mkdirTemp: os.MkdirTemp,
		copyPath:  copyPath,
		rename:    os.Rename,
	})
}

type cacheStoreArtifactOperations struct {
	stat      func(string) (fs.FileInfo, error)
	digest    func(string) (string, error)
	removeAll func(string) error
	mkdirAll  func(string, fs.FileMode) error
	mkdirTemp func(string, string) (string, error)
	copyPath  func(string, string) error
	rename    func(string, string) error
}

func (c *Cache) storeArtifactWithOperations(digest, source string, operations cacheStoreArtifactOperations) error {
	destination := filepath.Join(c.artifactRoot, digest)
	if _, err := operations.stat(filepath.Join(destination, "payload")); err == nil {
		if actual, digestErr := operations.digest(filepath.Join(destination, "payload")); digestErr == nil && actual == digest {
			return nil
		}
		if err := operations.removeAll(destination); err != nil {
			return err
		}
	}
	if err := operations.mkdirAll(c.artifactRoot, 0o755); err != nil {
		return err
	}
	tmp, err := operations.mkdirTemp(c.artifactRoot, ".artifact-*")
	if err != nil {
		return err
	}
	defer operations.removeAll(tmp)
	if err := operations.copyPath(source, filepath.Join(tmp, "payload")); err != nil {
		return err
	}
	if err := operations.rename(tmp, destination); err != nil && !errors.Is(err, fs.ErrExist) {
		if _, statErr := operations.stat(filepath.Join(destination, "payload")); statErr != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) restoreArtifact(digest, destination string) error {
	return c.restoreArtifactWithOperations(digest, destination, cacheRestoreArtifactOperations{
		stat:      os.Stat,
		digest:    artifactDigest,
		mkdirAll:  os.MkdirAll,
		mkdirTemp: os.MkdirTemp,
		removeAll: os.RemoveAll,
		copyPath:  copyPath,
		rename:    os.Rename,
	})
}

type cacheRestoreArtifactOperations struct {
	stat      func(string) (fs.FileInfo, error)
	digest    func(string) (string, error)
	mkdirAll  func(string, fs.FileMode) error
	mkdirTemp func(string, string) (string, error)
	removeAll func(string) error
	copyPath  func(string, string) error
	rename    func(string, string) error
}

func (c *Cache) restoreArtifactWithOperations(digest, destination string, operations cacheRestoreArtifactOperations) error {
	source := filepath.Join(c.artifactRoot, digest, "payload")
	if _, err := operations.stat(source); err != nil {
		return err
	}
	if actual, err := operations.digest(source); err != nil {
		return err
	} else if actual != digest {
		return fmt.Errorf("%w: cache payload %s has digest %s", errCorruptArtifact, digest, actual)
	}
	if err := operations.mkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	tmp, err := operations.mkdirTemp(filepath.Dir(destination), ".wake-restore-*")
	if err != nil {
		return err
	}
	defer operations.removeAll(tmp)
	payload := filepath.Join(tmp, "payload")
	if err := operations.copyPath(source, payload); err != nil {
		return err
	}
	if actual, err := operations.digest(payload); err != nil {
		return err
	} else if actual != digest {
		return fmt.Errorf("%w: copied payload %s has digest %s", errCorruptArtifact, digest, actual)
	}
	return operations.rename(payload, destination)
}

var errCorruptArtifact = errors.New("corrupt cached artifact")

func artifactDigest(root string) (string, error) {
	return artifactDigestWithOperations(root, cacheArtifactDigestOperations{
		lstat:    os.Lstat,
		walkDir:  filepath.WalkDir,
		rel:      filepath.Rel,
		readlink: os.Readlink,
		open:     func(path string) (cacheDigestFile, error) { return os.Open(path) },
	})
}

type cacheArtifactDigestOperations struct {
	lstat    func(string) (fs.FileInfo, error)
	walkDir  func(string, fs.WalkDirFunc) error
	rel      func(string, string) (string, error)
	readlink func(string) (string, error)
	open     func(string) (cacheDigestFile, error)
}

func artifactDigestWithOperations(root string, operations cacheArtifactDigestOperations) (string, error) {
	root = filepath.Clean(root)
	if _, err := operations.lstat(root); err != nil {
		return "", err
	}
	hash := blake3.New()
	writePart(hash, "wake.artifact.v2")
	err := operations.walkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		relative, err := operations.rel(root, path)
		if err != nil {
			return err
		}
		writePart(hash, filepath.ToSlash(relative))
		writePart(hash, strconv.FormatUint(uint64(relevantMode(info.Mode())), 8))
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := operations.readlink(path)
			if err != nil {
				return err
			}
			writePart(hash, target)
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		file, err := operations.open(path)
		if err != nil {
			return err
		}
		content := blake3.New()
		_, copyErr := io.Copy(content, file)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		writePart(hash, hex.EncodeToString(content.Sum(nil)))
		return nil
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func copyPath(source, destination string) error {
	return copyPathWithOperations(source, destination, cacheCopyPathOperations{
		lstat:    os.Lstat,
		readlink: os.Readlink,
		symlink:  os.Symlink,
		mkdirAll: os.MkdirAll,
		readDir:  os.ReadDir,
		copyFile: copyFile,
	})
}

type cacheCopyPathOperations struct {
	lstat    func(string) (fs.FileInfo, error)
	readlink func(string) (string, error)
	symlink  func(string, string) error
	mkdirAll func(string, fs.FileMode) error
	readDir  func(string) ([]fs.DirEntry, error)
	copyFile func(string, string, fs.FileMode) error
}

func copyPathWithOperations(source, destination string, operations cacheCopyPathOperations) error {
	info, err := operations.lstat(source)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := operations.readlink(source)
		if err != nil {
			return err
		}
		return operations.symlink(target, destination)
	}
	if !info.IsDir() {
		return operations.copyFile(source, destination, info.Mode())
	}
	if err := operations.mkdirAll(destination, info.Mode().Perm()); err != nil {
		return err
	}
	entries, err := operations.readDir(source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := copyPathWithOperations(filepath.Join(source, entry.Name()), filepath.Join(destination, entry.Name()), operations); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(source, destination string, mode fs.FileMode) error {
	return copyFileWithOperations(source, destination, mode, cacheCopyFileOperations{
		mkdirAll: os.MkdirAll,
		open:     func(path string) (cacheDigestFile, error) { return os.Open(path) },
		openFile: func(path string, flag int, permission fs.FileMode) (cacheCopyDestination, error) {
			return os.OpenFile(path, flag, permission)
		},
	})
}

type cacheCopyDestination interface {
	io.Writer
	Close() error
}

type cacheCopyFileOperations struct {
	mkdirAll func(string, fs.FileMode) error
	open     func(string) (cacheDigestFile, error)
	openFile func(string, int, fs.FileMode) (cacheCopyDestination, error)
}

func copyFileWithOperations(source, destination string, mode fs.FileMode, operations cacheCopyFileOperations) error {
	if err := operations.mkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	in, err := operations.open(source)
	if err != nil {
		return err
	}
	out, err := operations.openFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm())
	if err != nil {
		_ = in.Close()
		return err
	}
	_, copyErr := io.Copy(out, in)
	outCloseErr := out.Close()
	inCloseErr := in.Close()
	if copyErr != nil {
		return copyErr
	}
	if outCloseErr != nil {
		return outCloseErr
	}
	return inCloseErr
}

func writePart(w *blake3.Hasher, value string) {
	var encodedLength [24]byte
	length := strconv.AppendInt(encodedLength[:0], int64(len(value)), 10)
	_, _ = w.Write(length)
	_, _ = w.WriteString(":")
	_, _ = w.WriteString(value)
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
