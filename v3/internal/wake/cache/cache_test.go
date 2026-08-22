package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charlievieth/fastwalk"
	gitignore "github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenCacheFormatMigrationAndFailures(t *testing.T) {
	root := t.TempDir()
	want := errors.New("injected failure")
	base := cacheOpenOperations{
		abs:          func(path string) (string, error) { return path, nil },
		mkdirAll:     func(string, fs.FileMode) error { return nil },
		userCacheDir: func() (string, error) { return "/cache", nil },
		readFile:     func(string) ([]byte, error) { return nil, fs.ErrNotExist },
	}

	operations := base
	operations.abs = func(string) (string, error) { return "", want }
	_, err := openCacheWithOperations(root, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.mkdirAll = func(string, fs.FileMode) error { return want }
	_, err = openCacheWithOperations(root, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.userCacheDir = func() (string, error) { return "", want }
	_, err = openCacheWithOperations(root, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.readFile = func(path string) ([]byte, error) {
		if filepath.Base(path) == "wake-index-v2.gob" {
			return nil, want
		}
		return nil, fs.ErrNotExist
	}
	_, err = openCacheWithOperations(root, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.readFile = func(path string) ([]byte, error) {
		if filepath.Base(path) == "wake-index.json" {
			return nil, want
		}
		return nil, fs.ErrNotExist
	}
	_, err = openCacheWithOperations(root, operations)
	require.ErrorIs(t, err, want)

	var gobData bytes.Buffer
	require.NoError(t, gob.NewEncoder(&gobData).Encode(indexData{}))
	operations = base
	operations.readFile = func(path string) ([]byte, error) {
		if filepath.Base(path) == "wake-index-v2.gob" {
			return gobData.Bytes(), nil
		}
		return nil, fs.ErrNotExist
	}
	loaded, err := openCacheWithOperations(root, operations)
	require.NoError(t, err)
	assert.NotNil(t, loaded.index.Files)
	assert.NotNil(t, loaded.index.GoAPI)
	assert.NotNil(t, loaded.index.Actions)
	assert.NotNil(t, loaded.index.Receipts)

	operations.readFile = func(path string) ([]byte, error) {
		if filepath.Base(path) == "wake-index-v2.gob" {
			return []byte("corrupt"), nil
		}
		return nil, fs.ErrNotExist
	}
	loaded, err = openCacheWithOperations(root, operations)
	require.NoError(t, err)
	assert.Empty(t, loaded.index.Actions)

	legacy, err := json.Marshal(indexData{Actions: map[string]ActionResult{"action": {Artifact: "digest", Output: "bin/app"}}})
	require.NoError(t, err)
	operations.readFile = func(path string) ([]byte, error) {
		if filepath.Base(path) == "wake-index.json" {
			return legacy, nil
		}
		return nil, fs.ErrNotExist
	}
	loaded, err = openCacheWithOperations(root, operations)
	require.NoError(t, err)
	assert.Equal(t, "digest", loaded.index.Actions["action"].Artifact)
	assert.True(t, loaded.dirty, "legacy JSON is rewritten as the binary v2 format")

	operations.readFile = func(path string) ([]byte, error) {
		if filepath.Base(path) == "wake-index.json" {
			return []byte("corrupt"), nil
		}
		return nil, fs.ErrNotExist
	}
	loaded, err = openCacheWithOperations(root, operations)
	require.NoError(t, err)
	assert.Empty(t, loaded.index.Actions)
}

func TestReadOnlyCacheOpenAndPeekNeverCreateOrRestoreProjectState(t *testing.T) {
	root := t.TempDir()
	store, err := OpenCacheReadOnly(root)
	require.NoError(t, err)
	assert.NoDirExists(t, filepath.Join(root, ".wails"))
	status, artifact, err := store.Peek("missing", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
	assert.Empty(t, artifact)
	assert.NoFileExists(t, filepath.Join(root, "bin", "app"))
}

func TestArtifactCacheVerifiesStoredAndRestoredPayloads(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	output := filepath.Join(root, "bin", "app")
	require.NoError(t, os.WriteFile(output, []byte("trusted"), 0o755))
	store, err := OpenCache(root)
	require.NoError(t, err)
	digest, err := store.RecordAction("compile", "bin/app")
	require.NoError(t, err)
	payload := filepath.Join(store.artifactRoot, digest, "payload")
	require.FileExists(t, payload)

	require.NoError(t, os.WriteFile(payload, []byte("corrupt"), 0o755))
	require.NoError(t, os.Remove(output))
	status, _, err := store.Peek("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
	assert.NoFileExists(t, output, "read-only inspection must not restore a corrupt payload")
	status, _, err = store.Lookup("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
	assert.NoFileExists(t, output)
	assert.NoDirExists(t, filepath.Dir(payload), "a corrupt generated cache entry should be discarded")

	require.NoError(t, os.WriteFile(output, []byte("trusted"), 0o755))
	digest, err = store.RecordAction("compile", "bin/app")
	require.NoError(t, err)
	require.NoError(t, os.Remove(output))
	status, restoredDigest, err := store.Lookup("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupRestored, status)
	assert.Equal(t, digest, restoredDigest)
	restored, err := os.ReadFile(output)
	require.NoError(t, err)
	assert.Equal(t, []byte("trusted"), restored)

	require.NoError(t, os.WriteFile(output, []byte("user-modified"), 0o755))
	status, _, err = store.Lookup("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupDirty, status)
}

func TestActionKeysAndReceiptLifecycleAreExact(t *testing.T) {
	base, err := ActionKey("compile", map[string]any{"target": "linux/amd64"}, []string{"source:a"}, []string{"dependency:a"})
	require.NoError(t, err)
	require.Len(t, base, 64)
	for name, candidate := range map[string]func() (string, error){
		"kind": func() (string, error) {
			return ActionKey("package", map[string]any{"target": "linux/amd64"}, []string{"source:a"}, []string{"dependency:a"})
		},
		"spec": func() (string, error) {
			return ActionKey("compile", map[string]any{"target": "linux/arm64"}, []string{"source:a"}, []string{"dependency:a"})
		},
		"input": func() (string, error) {
			return ActionKey("compile", map[string]any{"target": "linux/amd64"}, []string{"source:b"}, []string{"dependency:a"})
		},
		"dependency": func() (string, error) {
			return ActionKey("compile", map[string]any{"target": "linux/amd64"}, []string{"source:a"}, []string{"dependency:b"})
		},
	} {
		t.Run(name, func(t *testing.T) {
			key, err := candidate()
			require.NoError(t, err)
			assert.NotEqual(t, base, key)
		})
	}
	_, err = ActionKey("compile", make(chan int), nil, nil)
	require.Error(t, err)

	root := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	store, err := OpenCache(root)
	require.NoError(t, err)
	marker := filepath.Join("frontend", "node_modules", ".wails-receipt")
	assert.False(t, store.HasReceipt("install", marker))
	require.NoError(t, store.RecordReceipt("install"))
	assert.False(t, store.HasReceipt("install", marker))
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, marker)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, marker), []byte("complete"), 0o644))
	assert.True(t, store.HasReceipt("install", marker))

	reopened, err := OpenCache(root)
	require.NoError(t, err)
	assert.True(t, reopened.HasReceipt("install", marker))
}

func TestPeekCoversHitDirtyRestorableAndMissingPayloadStates(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	output := filepath.Join(root, "bin", "app")
	require.NoError(t, os.WriteFile(output, []byte("trusted"), 0o755))
	store, err := OpenCache(root)
	require.NoError(t, err)
	digest, err := store.RecordAction("compile", "bin/app")
	require.NoError(t, err)

	status, gotDigest, err := store.Peek("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupHit, status)
	assert.Equal(t, digest, gotDigest)
	require.NoError(t, os.WriteFile(output, []byte("dirty"), 0o755))
	status, _, err = store.Peek("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupDirty, status)
	require.NoError(t, os.Remove(output))
	status, gotDigest, err = store.Peek("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupRestored, status)
	assert.Equal(t, digest, gotDigest)
	require.NoError(t, os.RemoveAll(filepath.Join(store.artifactRoot, digest)))
	status, _, err = store.Peek("compile", "bin/app")
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
}

func TestLookupFaultsAndStateTransitionsAreExact(t *testing.T) {
	want := errors.New("injected lookup failure")
	newStore := func() *Cache {
		return &Cache{
			root: "/project", artifactRoot: "/cache/artifacts",
			index: indexData{Actions: map[string]ActionResult{"action": {Artifact: "digest", Output: "bin/app"}}},
		}
	}
	base := cacheLookupOperations{
		lstat:     func(string) (fs.FileInfo, error) { return cacheStaticFileInfo{}, nil },
		restore:   func(string, string) error { return nil },
		removeAll: func(string) error { return nil },
		save:      func() error { return nil },
		digest:    func(string) (string, error) { return "digest", nil },
	}

	store := newStore()
	status, _, err := store.lookupWithOperations("missing", "bin/app", base)
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
	status, _, err = store.lookupWithOperations("action", "bin/other", base)
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)

	operations := base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	status, artifact, err := store.lookupWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupRestored, status)
	assert.Equal(t, "digest", artifact)

	operations.restore = func(string, string) error { return fs.ErrNotExist }
	status, _, err = store.lookupWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)

	operations.restore = func(string, string) error { return want }
	status, _, err = store.lookupWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	assert.Equal(t, LookupMiss, status)

	store = newStore()
	removed, saved := "", false
	operations.restore = func(string, string) error { return errCorruptArtifact }
	operations.removeAll = func(path string) error { removed = path; return want }
	operations.save = func() error { saved = true; return want }
	status, _, err = store.lookupWithOperations("action", "bin/app", operations)
	require.NoError(t, err, "disposable cache cleanup failures must remain a safe miss")
	assert.Equal(t, LookupMiss, status)
	assert.Equal(t, filepath.Join("/cache/artifacts", "digest"), removed)
	assert.True(t, saved)
	assert.True(t, store.dirty)
	assert.NotContains(t, store.index.Actions, "action")

	store = newStore()
	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	_, _, err = store.lookupWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations = base
	operations.digest = func(string) (string, error) { return "", want }
	_, _, err = store.lookupWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations.digest = func(string) (string, error) { return "different", nil }
	status, _, err = store.lookupWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupDirty, status)
	operations.digest = func(string) (string, error) { return "digest", nil }
	status, artifact, err = store.lookupWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupHit, status)
	assert.Equal(t, "digest", artifact)
}

func TestPeekFaultsAndStateTransitionsAreExact(t *testing.T) {
	want := errors.New("injected peek failure")
	store := &Cache{
		root: "/project", artifactRoot: "/cache/artifacts",
		index: indexData{Actions: map[string]ActionResult{"action": {Artifact: "digest", Output: "bin/app"}}},
	}
	base := cachePeekOperations{
		lstat:  func(string) (fs.FileInfo, error) { return cacheStaticFileInfo{}, nil },
		stat:   func(string) (fs.FileInfo, error) { return cacheStaticFileInfo{}, nil },
		digest: func(string) (string, error) { return "digest", nil },
	}
	status, _, err := store.peekWithOperations("missing", "bin/app", base)
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
	status, _, err = store.peekWithOperations("action", "bin/other", base)
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)

	operations := base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	status, artifact, err := store.peekWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupRestored, status)
	assert.Equal(t, "digest", artifact)
	operations.digest = func(string) (string, error) { return "", want }
	status, _, err = store.peekWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	operations.stat = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	status, _, err = store.peekWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupMiss, status)
	operations.stat = func(string) (fs.FileInfo, error) { return nil, want }
	_, _, err = store.peekWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	_, _, err = store.peekWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations = base
	operations.digest = func(string) (string, error) { return "", want }
	_, _, err = store.peekWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations.digest = func(string) (string, error) { return "different", nil }
	status, _, err = store.peekWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	assert.Equal(t, LookupDirty, status)
}

func TestRecordActionPropagatesDigestStoreAndSaveFailures(t *testing.T) {
	want := errors.New("injected record failure")
	newStore := func() *Cache { return &Cache{root: "/project", index: indexData{Actions: map[string]ActionResult{}}} }
	info := cacheStaticFileInfo{size: 7, mode: 0o755, modTime: time.Unix(1, 2)}
	base := cacheRecordOperations{
		digest: func(string) (string, error) { return "digest", nil },
		store:  func(string, string) error { return nil },
		lstat:  func(string) (fs.FileInfo, error) { return info, nil },
		save:   func() error { return nil },
	}
	operations := base
	operations.digest = func(string) (string, error) { return "", want }
	_, err := newStore().recordActionWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations = base
	operations.store = func(string, string) error { return want }
	_, err = newStore().recordActionWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	_, err = newStore().recordActionWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations = base
	lstatCalls := 0
	operations.lstat = func(string) (fs.FileInfo, error) {
		lstatCalls++
		if lstatCalls == 2 {
			return nil, want
		}
		return info, nil
	}
	_, err = newStore().recordActionWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	operations = base
	operations.save = func() error { return want }
	store := newStore()
	_, err = store.recordActionWithOperations("action", "bin/app", operations)
	require.ErrorIs(t, err, want)
	assert.Equal(t, ActionResult{Artifact: "digest", Output: "bin/app", OutputMetadata: fileRecord(info, "digest")}, store.index.Actions["action"])
	assert.True(t, store.dirty)
	digest, err := newStore().recordActionWithOperations("action", "bin/app", base)
	require.NoError(t, err)
	assert.Equal(t, "digest", digest)
}

func TestActionOutputMetadataFastPathRequiresChangeTrackingIdentity(t *testing.T) {
	t.Run("fallback metadata never replaces content verification", func(t *testing.T) {
		info := cacheStaticFileInfo{size: 7, mode: 0o755, modTime: time.Unix(1, 2)}
		result := ActionResult{Artifact: "original", Output: "bin/app", OutputMetadata: fileRecord(info, "original")}
		store := &Cache{root: "/project", index: indexData{Actions: map[string]ActionResult{"action": result}}}
		digests := 0
		lookup := cacheLookupOperations{
			lstat:     func(string) (fs.FileInfo, error) { return info, nil },
			restore:   func(string, string) error { return nil },
			removeAll: func(string) error { return nil },
			save:      func() error { return nil },
			digest: func(string) (string, error) {
				digests++
				return "changed", nil
			},
		}
		status, _, err := store.lookupWithOperations("action", "bin/app", lookup)
		require.NoError(t, err)
		assert.Equal(t, LookupDirty, status)
		assert.Equal(t, 1, digests)
		peek := cachePeekOperations{lstat: lookup.lstat, stat: lookup.lstat, digest: func(string) (string, error) {
			digests++
			return "original", nil
		}}
		status, _, err = store.peekWithOperations("action", "bin/app", peek)
		require.NoError(t, err)
		assert.Equal(t, LookupHit, status)
		assert.Equal(t, 2, digests)
	})

	t.Run("same size and mtime byte edit is invalidated", func(t *testing.T) {
		root := t.TempDir()
		t.Setenv("XDG_CACHE_HOME", t.TempDir())
		require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
		output := filepath.Join(root, "bin", "app")
		require.NoError(t, os.WriteFile(output, []byte("trusted"), 0o755))
		originalInfo, err := os.Stat(output)
		require.NoError(t, err)
		store, err := OpenCache(root)
		require.NoError(t, err)
		_, err = store.RecordAction("action", "bin/app")
		require.NoError(t, err)

		digests := 0
		lookup := cacheLookupOperations{
			lstat:     os.Lstat,
			restore:   func(string, string) error { return nil },
			removeAll: func(string) error { return nil },
			save:      func() error { return nil },
			digest: func(path string) (string, error) {
				digests++
				return artifactDigest(path)
			},
		}
		status, _, err := store.lookupWithOperations("action", "bin/app", lookup)
		require.NoError(t, err)
		assert.Equal(t, LookupHit, status)
		if platformIdentityTracksChanges() {
			assert.Zero(t, digests, "change-tracking identity should enable the fast path")
		} else {
			assert.Equal(t, 1, digests, "fallback platforms must verify content")
		}
		beforePeekDigests := digests
		peek := cachePeekOperations{lstat: os.Lstat, stat: os.Stat, digest: lookup.digest}
		status, _, err = store.peekWithOperations("action", "bin/app", peek)
		require.NoError(t, err)
		assert.Equal(t, LookupHit, status)
		if platformIdentityTracksChanges() {
			assert.Equal(t, beforePeekDigests, digests)
		} else {
			assert.Equal(t, beforePeekDigests+1, digests)
		}

		require.NoError(t, os.WriteFile(output, []byte("changed"), 0o755))
		require.NoError(t, os.Chtimes(output, originalInfo.ModTime(), originalInfo.ModTime()))
		status, _, err = store.lookupWithOperations("action", "bin/app", lookup)
		require.NoError(t, err)
		assert.Equal(t, LookupDirty, status)
		assert.Greater(t, digests, 0, "same-size bytes with restored mtime must be hashed when identity changed")
	})
}

func TestRecordActionNeverPairsADigestWithChangedOutputMetadata(t *testing.T) {
	before := cacheStaticFileInfo{size: 7, mode: 0o755, modTime: time.Unix(1, 2)}
	after := before
	after.size++
	stats := []fs.FileInfo{before, after}
	operations := cacheRecordOperations{
		digest: func(string) (string, error) { return "original-digest", nil },
		store:  func(string, string) error { return nil },
		lstat: func(string) (fs.FileInfo, error) {
			result := stats[0]
			stats = stats[1:]
			return result, nil
		},
		save: func() error { return nil },
	}
	store := &Cache{root: "/project", index: indexData{Actions: map[string]ActionResult{}}}
	_, err := store.recordActionWithOperations("action", "bin/app", operations)
	require.NoError(t, err)
	result := store.index.Actions["action"]
	assert.Empty(t, result.OutputMetadata, "metadata from a different output snapshot must not enable the fast path")

	lookup := cacheLookupOperations{
		lstat:     func(string) (fs.FileInfo, error) { return after, nil },
		restore:   func(string, string) error { return nil },
		removeAll: func(string) error { return nil },
		save:      func() error { return nil },
		digest:    func(string) (string, error) { return "changed-digest", nil },
	}
	status, _, err := store.lookupWithOperations("action", "bin/app", lookup)
	require.NoError(t, err)
	assert.Equal(t, LookupDirty, status)
}

func TestStoreArtifactHandlesEveryPublicationOutcome(t *testing.T) {
	want := errors.New("injected store failure")
	store := &Cache{artifactRoot: "/cache/artifacts"}
	newBase := func() cacheStoreArtifactOperations {
		return cacheStoreArtifactOperations{
			stat:      func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist },
			digest:    func(string) (string, error) { return "digest", nil },
			removeAll: func(string) error { return nil },
			mkdirAll:  func(string, fs.FileMode) error { return nil },
			mkdirTemp: func(string, string) (string, error) { return "/cache/artifacts/temp", nil },
			copyPath:  func(string, string) error { return nil },
			rename:    func(string, string) error { return nil },
		}
	}

	operations := newBase()
	operations.stat = func(string) (fs.FileInfo, error) { return cacheStaticFileInfo{}, nil }
	require.NoError(t, store.storeArtifactWithOperations("digest", "/source", operations), "an existing verified payload is reused")

	operations = newBase()
	operations.stat = func(string) (fs.FileInfo, error) { return cacheStaticFileInfo{}, nil }
	operations.digest = func(string) (string, error) { return "different", nil }
	operations.removeAll = func(path string) error {
		if path == filepath.Join("/cache/artifacts", "digest") {
			return want
		}
		return nil
	}
	require.ErrorIs(t, store.storeArtifactWithOperations("digest", "/source", operations), want)

	operations = newBase()
	operations.mkdirAll = func(string, fs.FileMode) error { return want }
	require.ErrorIs(t, store.storeArtifactWithOperations("digest", "/source", operations), want)
	operations = newBase()
	operations.mkdirTemp = func(string, string) (string, error) { return "", want }
	require.ErrorIs(t, store.storeArtifactWithOperations("digest", "/source", operations), want)
	operations = newBase()
	operations.copyPath = func(string, string) error { return want }
	require.ErrorIs(t, store.storeArtifactWithOperations("digest", "/source", operations), want)
	operations = newBase()
	operations.rename = func(string, string) error { return fs.ErrExist }
	require.NoError(t, store.storeArtifactWithOperations("digest", "/source", operations))
	operations = newBase()
	operations.rename = func(string, string) error { return want }
	statCalls := 0
	operations.stat = func(string) (fs.FileInfo, error) {
		statCalls++
		if statCalls == 1 {
			return nil, fs.ErrNotExist
		}
		return cacheStaticFileInfo{}, nil
	}
	require.NoError(t, store.storeArtifactWithOperations("digest", "/source", operations), "a concurrent publisher wins safely")
	operations = newBase()
	operations.rename = func(string, string) error { return want }
	require.ErrorIs(t, store.storeArtifactWithOperations("digest", "/source", operations), want)

	operations = newBase()
	removed := ""
	operations.removeAll = func(path string) error { removed = path; return nil }
	require.NoError(t, store.storeArtifactWithOperations("digest", "/source", operations))
	assert.Equal(t, "/cache/artifacts/temp", removed)
}

func TestRestoreArtifactRejectsEveryIncompleteOrCorruptOutcome(t *testing.T) {
	want := errors.New("injected restore failure")
	store := &Cache{artifactRoot: "/cache/artifacts"}
	newBase := func() cacheRestoreArtifactOperations {
		return cacheRestoreArtifactOperations{
			stat:      func(string) (fs.FileInfo, error) { return cacheStaticFileInfo{}, nil },
			digest:    func(string) (string, error) { return "digest", nil },
			mkdirAll:  func(string, fs.FileMode) error { return nil },
			mkdirTemp: func(string, string) (string, error) { return "/output/temp", nil },
			removeAll: func(string) error { return nil },
			copyPath:  func(string, string) error { return nil },
			rename:    func(string, string) error { return nil },
		}
	}

	operations := newBase()
	operations.stat = func(string) (fs.FileInfo, error) { return nil, want }
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), want)
	operations = newBase()
	operations.digest = func(string) (string, error) { return "", want }
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), want)
	operations = newBase()
	operations.digest = func(string) (string, error) { return "different", nil }
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), errCorruptArtifact)
	operations = newBase()
	operations.mkdirAll = func(string, fs.FileMode) error { return want }
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), want)
	operations = newBase()
	operations.mkdirTemp = func(string, string) (string, error) { return "", want }
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), want)
	operations = newBase()
	operations.copyPath = func(string, string) error { return want }
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), want)
	operations = newBase()
	digests := 0
	operations.digest = func(string) (string, error) {
		digests++
		if digests == 2 {
			return "", want
		}
		return "digest", nil
	}
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), want)
	operations = newBase()
	digests = 0
	operations.digest = func(string) (string, error) {
		digests++
		if digests == 2 {
			return "different", nil
		}
		return "digest", nil
	}
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), errCorruptArtifact)
	operations = newBase()
	operations.rename = func(string, string) error { return want }
	require.ErrorIs(t, store.restoreArtifactWithOperations("digest", "/output/app", operations), want)
	operations = newBase()
	removed := ""
	operations.removeAll = func(path string) error { removed = path; return nil }
	require.NoError(t, store.restoreArtifactWithOperations("digest", "/output/app", operations))
	assert.Equal(t, "/output/temp", removed)
}

func TestArtifactDigestIsIndependentOfRootNameAndIncludesTreeMetadata(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	require.NoError(t, os.WriteFile(first, []byte("same"), 0o644))
	require.NoError(t, os.WriteFile(second, []byte("same"), 0o644))
	firstDigest, err := artifactDigest(first)
	require.NoError(t, err)
	secondDigest, err := artifactDigest(second)
	require.NoError(t, err)
	assert.Equal(t, firstDigest, secondDigest)
	require.NoError(t, os.Chmod(second, 0o755))
	secondDigest, err = artifactDigest(second)
	require.NoError(t, err)
	assert.NotEqual(t, firstDigest, secondDigest)

	firstDirectory := filepath.Join(root, "one")
	secondDirectory := filepath.Join(root, "two")
	for _, directory := range []string{firstDirectory, secondDirectory} {
		require.NoError(t, os.MkdirAll(directory, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(directory, "file"), []byte("content"), 0o644))
		require.NoError(t, os.Symlink("file", filepath.Join(directory, "link")))
	}
	oneDigest, err := artifactDigest(firstDirectory)
	require.NoError(t, err)
	twoDigest, err := artifactDigest(secondDirectory)
	require.NoError(t, err)
	assert.Equal(t, oneDigest, twoDigest)
	require.NoError(t, os.Remove(filepath.Join(secondDirectory, "link")))
	require.NoError(t, os.Symlink("other", filepath.Join(secondDirectory, "link")))
	twoDigest, err = artifactDigest(secondDirectory)
	require.NoError(t, err)
	assert.NotEqual(t, oneDigest, twoDigest)
}

func TestArtifactDigestFaultsAreNeverHidden(t *testing.T) {
	want := errors.New("injected artifact digest failure")
	regular := cacheStaticFileInfo{name: "file", mode: 0o644}
	symlink := cacheStaticFileInfo{name: "link", mode: os.ModeSymlink | 0o777}
	base := cacheArtifactDigestOperations{
		lstat:    func(string) (fs.FileInfo, error) { return regular, nil },
		walkDir:  func(string, fs.WalkDirFunc) error { return nil },
		rel:      filepath.Rel,
		readlink: func(string) (string, error) { return "target", nil },
		open: func(string) (cacheDigestFile, error) {
			return &cacheFaultDigestFile{data: []byte("data")}, nil
		},
	}
	operations := base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	_, err := artifactDigestWithOperations("/artifact", operations)
	require.ErrorIs(t, err, want)
	operations = base
	operations.walkDir = func(string, fs.WalkDirFunc) error { return want }
	_, err = artifactDigestWithOperations("/artifact", operations)
	require.ErrorIs(t, err, want)

	for name, invoke := range map[string]func(fs.WalkDirFunc) error{
		"walk": func(callback fs.WalkDirFunc) error { return callback("/artifact", nil, want) },
		"info": func(callback fs.WalkDirFunc) error {
			return callback("/artifact/file", cacheDirEntry{name: "file", infoErr: want}, nil)
		},
		"relative": func(callback fs.WalkDirFunc) error {
			return callback("/artifact/file", cacheDirEntry{name: "file", info: regular}, nil)
		},
		"readlink": func(callback fs.WalkDirFunc) error {
			return callback("/artifact/link", cacheDirEntry{name: "link", info: symlink}, nil)
		},
		"open": func(callback fs.WalkDirFunc) error {
			return callback("/artifact/file", cacheDirEntry{name: "file", info: regular}, nil)
		},
		"copy": func(callback fs.WalkDirFunc) error {
			return callback("/artifact/file", cacheDirEntry{name: "file", info: regular}, nil)
		},
		"close": func(callback fs.WalkDirFunc) error {
			return callback("/artifact/file", cacheDirEntry{name: "file", info: regular}, nil)
		},
	} {
		t.Run(name, func(t *testing.T) {
			operations := base
			operations.walkDir = func(_ string, callback fs.WalkDirFunc) error { return invoke(callback) }
			switch name {
			case "relative":
				operations.rel = func(string, string) (string, error) { return "", want }
			case "readlink":
				operations.readlink = func(string) (string, error) { return "", want }
			case "open":
				operations.open = func(string) (cacheDigestFile, error) { return nil, want }
			case "copy":
				operations.open = func(string) (cacheDigestFile, error) { return &cacheFaultDigestFile{readErr: want}, nil }
			case "close":
				operations.open = func(string) (cacheDigestFile, error) { return &cacheFaultDigestFile{closeErr: want}, nil }
			}
			_, err := artifactDigestWithOperations("/artifact", operations)
			require.ErrorIs(t, err, want)
		})
	}
}

func TestCopyPathPropagatesEveryFilesystemFailure(t *testing.T) {
	want := errors.New("injected copy path failure")
	regular := cacheStaticFileInfo{name: "file", mode: 0o644}
	directory := cacheStaticFileInfo{name: "directory", mode: os.ModeDir | 0o755}
	symlink := cacheStaticFileInfo{name: "link", mode: os.ModeSymlink | 0o777}
	base := cacheCopyPathOperations{
		lstat:    func(string) (fs.FileInfo, error) { return regular, nil },
		readlink: func(string) (string, error) { return "target", nil },
		symlink:  func(string, string) error { return nil },
		mkdirAll: func(string, fs.FileMode) error { return nil },
		readDir:  func(string) ([]fs.DirEntry, error) { return nil, nil },
		copyFile: func(string, string, fs.FileMode) error { return nil },
	}
	operations := base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	require.ErrorIs(t, copyPathWithOperations("source", "destination", operations), want)
	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return symlink, nil }
	operations.readlink = func(string) (string, error) { return "", want }
	require.ErrorIs(t, copyPathWithOperations("source", "destination", operations), want)
	operations.readlink = func(string) (string, error) { return "target", nil }
	operations.symlink = func(string, string) error { return want }
	require.ErrorIs(t, copyPathWithOperations("source", "destination", operations), want)
	operations = base
	operations.copyFile = func(string, string, fs.FileMode) error { return want }
	require.ErrorIs(t, copyPathWithOperations("source", "destination", operations), want)
	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return directory, nil }
	operations.mkdirAll = func(string, fs.FileMode) error { return want }
	require.ErrorIs(t, copyPathWithOperations("source", "destination", operations), want)
	operations.mkdirAll = func(string, fs.FileMode) error { return nil }
	operations.readDir = func(string) ([]fs.DirEntry, error) { return nil, want }
	require.ErrorIs(t, copyPathWithOperations("source", "destination", operations), want)
	operations = base
	operations.lstat = func(path string) (fs.FileInfo, error) {
		if strings.HasSuffix(path, "child") {
			return regular, nil
		}
		return directory, nil
	}
	operations.readDir = func(string) ([]fs.DirEntry, error) { return []fs.DirEntry{cacheDirEntry{name: "child"}}, nil }
	operations.copyFile = func(string, string, fs.FileMode) error { return want }
	require.ErrorIs(t, copyPathWithOperations("source", "destination", operations), want)
	operations.copyFile = func(string, string, fs.FileMode) error { return nil }
	require.NoError(t, copyPathWithOperations("source", "destination", operations))
}

func TestCopyFilePropagatesEveryStreamFailure(t *testing.T) {
	want := errors.New("injected copy file failure")
	newBase := func() cacheCopyFileOperations {
		return cacheCopyFileOperations{
			mkdirAll: func(string, fs.FileMode) error { return nil },
			open: func(string) (cacheDigestFile, error) {
				return &cacheFaultDigestFile{data: []byte("data")}, nil
			},
			openFile: func(string, int, fs.FileMode) (cacheCopyDestination, error) {
				return &cacheFaultDestination{}, nil
			},
		}
	}
	operations := newBase()
	operations.mkdirAll = func(string, fs.FileMode) error { return want }
	require.ErrorIs(t, copyFileWithOperations("source", "destination", 0o755, operations), want)
	operations = newBase()
	operations.open = func(string) (cacheDigestFile, error) { return nil, want }
	require.ErrorIs(t, copyFileWithOperations("source", "destination", 0o755, operations), want)
	operations = newBase()
	operations.openFile = func(string, int, fs.FileMode) (cacheCopyDestination, error) { return nil, want }
	require.ErrorIs(t, copyFileWithOperations("source", "destination", 0o755, operations), want)
	operations = newBase()
	operations.openFile = func(string, int, fs.FileMode) (cacheCopyDestination, error) {
		return &cacheFaultDestination{writeErr: want}, nil
	}
	require.ErrorIs(t, copyFileWithOperations("source", "destination", 0o755, operations), want)
	operations = newBase()
	operations.openFile = func(string, int, fs.FileMode) (cacheCopyDestination, error) {
		return &cacheFaultDestination{closeErr: want}, nil
	}
	require.ErrorIs(t, copyFileWithOperations("source", "destination", 0o755, operations), want)
	operations = newBase()
	operations.open = func(string) (cacheDigestFile, error) {
		return &cacheFaultDigestFile{data: []byte("data"), closeErr: want}, nil
	}
	require.ErrorIs(t, copyFileWithOperations("source", "destination", 0o755, operations), want)
	require.NoError(t, copyFileWithOperations("source", "destination", 0o755, newBase()))
}

type cacheFaultDestination struct {
	writeErr error
	closeErr error
}

func (destination *cacheFaultDestination) Write(data []byte) (int, error) {
	if destination.writeErr != nil {
		return 0, destination.writeErr
	}
	return len(data), nil
}
func (destination *cacheFaultDestination) Close() error { return destination.closeErr }

func TestSaveIndexFaultsAreAtomicAndRetainDirtyState(t *testing.T) {
	want := errors.New("injected failure")
	newStore := func() *Cache {
		return &Cache{dirty: true, stateDir: "/state", indexPath: "/state/index", index: indexData{Files: map[string]FileRecord{}}}
	}
	newOperations := func(file *cacheFaultFile) cacheSaveOperations {
		return cacheSaveOperations{
			encode:     func(indexData) ([]byte, error) { return []byte("index"), nil },
			createTemp: func(string, string) (cacheTemporaryFile, error) { return file, nil },
			remove:     func(string) error { return nil },
			rename:     func(string, string) error { return nil },
		}
	}

	clean := newStore()
	clean.dirty = false
	require.NoError(t, clean.saveWithOperations(cacheSaveOperations{}))

	store := newStore()
	operations := newOperations(&cacheFaultFile{name: "/state/temp"})
	operations.encode = func(indexData) ([]byte, error) { return nil, want }
	require.ErrorIs(t, store.saveWithOperations(operations), want)

	store = newStore()
	operations = newOperations(&cacheFaultFile{name: "/state/temp"})
	operations.createTemp = func(string, string) (cacheTemporaryFile, error) { return nil, want }
	require.ErrorIs(t, store.saveWithOperations(operations), want)

	for name, file := range map[string]*cacheFaultFile{
		"write": {name: "/state/temp", writeErr: want},
		"short": {name: "/state/temp", shortWrite: true},
		"sync":  {name: "/state/temp", syncErr: want},
		"close": {name: "/state/temp", closeErr: want},
	} {
		t.Run(name, func(t *testing.T) {
			store := newStore()
			err := store.saveWithOperations(newOperations(file))
			require.Error(t, err)
			assert.True(t, store.dirty)
		})
	}

	store = newStore()
	operations = newOperations(&cacheFaultFile{name: "/state/temp"})
	operations.rename = func(string, string) error { return want }
	require.ErrorIs(t, store.saveWithOperations(operations), want)
	assert.True(t, store.dirty)

	store = newStore()
	removed := ""
	renamed := false
	operations = newOperations(&cacheFaultFile{name: "/state/temp"})
	operations.remove = func(path string) error { removed = path; return nil }
	operations.rename = func(source, destination string) error {
		renamed = source == "/state/temp" && destination == "/state/index"
		return nil
	}
	require.NoError(t, store.saveWithOperations(operations))
	assert.True(t, renamed)
	assert.Equal(t, "/state/temp", removed)
	assert.False(t, store.dirty)
}

type cacheFaultFile struct {
	name                        string
	writeErr, syncErr, closeErr error
	shortWrite                  bool
}

func (f *cacheFaultFile) Name() string { return f.name }
func (f *cacheFaultFile) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	if f.shortWrite {
		return len(data) - 1, nil
	}
	return len(data), nil
}
func (f *cacheFaultFile) Sync() error  { return f.syncErr }
func (f *cacheFaultFile) Close() error { return f.closeErr }

func TestSnapshotFilesIncludesExecutableMode(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "hook.sh")
	require.NoError(t, os.WriteFile(path, []byte("#!/bin/sh\n"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)
	plain, err := store.SnapshotFiles("hook", "hook.sh")
	require.NoError(t, err)

	require.NoError(t, os.Chmod(path, 0o755))
	executable, err := store.SnapshotFiles("hook", "hook.sh")
	require.NoError(t, err)
	assert.NotEqual(t, plain, executable)
}

func TestSnapshotFilesFingerprintsDeclaredDirectoriesRecursively(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "assets")
	require.NoError(t, os.MkdirAll(filepath.Join(assets, "nested"), 0o755))
	file := filepath.Join(assets, "nested", "icon.png")
	require.NoError(t, os.WriteFile(file, []byte("first"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)
	first, err := store.SnapshotFiles("assets", "assets")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(file, []byte("second"), 0o644))
	second, err := store.SnapshotFiles("assets", "assets")
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}

func TestSnapshotFilesFaultsMissingAndExternalPaths(t *testing.T) {
	want := errors.New("injected snapshot file failure")
	store := &Cache{root: "/project"}
	info := cacheStaticFileInfo{name: "source", mode: 0o644, size: 4, modTime: time.Now()}
	base := cacheSnapshotFileOperations{
		lstat:  func(string) (fs.FileInfo, error) { return info, nil },
		digest: func(string, fs.FileInfo) (string, error) { return "digest", nil },
		rel:    filepath.Rel,
	}

	operations := base
	operations.lstat = func(path string) (fs.FileInfo, error) {
		if strings.Contains(path, "missing") {
			return nil, fs.ErrNotExist
		}
		return info, nil
	}
	digest, err := store.snapshotFilesWithOperations("files", []string{"missing", "z", "a"}, operations)
	require.NoError(t, err)
	assert.NotEmpty(t, digest)

	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	_, err = store.snapshotFilesWithOperations("files", []string{"source"}, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.digest = func(string, fs.FileInfo) (string, error) { return "", want }
	_, err = store.snapshotFilesWithOperations("files", []string{"source"}, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.rel = func(string, string) (string, error) { return "", want }
	fromRelError, err := store.snapshotFilesWithOperations("files", []string{"/outside/source"}, operations)
	require.NoError(t, err)
	operations = base
	fromEscape, err := store.snapshotFilesWithOperations("files", []string{"/outside/source"}, operations)
	require.NoError(t, err)
	assert.Equal(t, fromEscape, fromRelError, "relative-path failures and explicit escapes use the same stable external identity")
}

func TestSnapshotStableMetadataFastPathAndSymlinkContent(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "source.txt")
	require.NoError(t, os.WriteFile(path, []byte("source"), 0o644))
	stable := time.Now().Add(-metadataFastPathWindow - time.Second)
	require.NoError(t, os.Chtimes(path, stable, stable))
	require.NoError(t, os.Symlink("source.txt", filepath.Join(root, "source.link")))
	store, err := OpenCache(root)
	require.NoError(t, err)
	first, err := store.Snapshot(SnapshotOptions{Label: "source", Root: ".", IncludeAll: true})
	require.NoError(t, err)
	store.ResetStats()
	second, err := store.Snapshot(SnapshotOptions{Label: "source", Root: ".", IncludeAll: true})
	require.NoError(t, err)
	assert.Equal(t, first, second)
	wantReused := 1
	if platformIdentityTracksChanges() {
		wantReused = 2
	}
	assert.Equal(t, wantReused, store.Stats().DigestsReused)

	info, err := os.Lstat(filepath.Join(root, "source.link"))
	require.NoError(t, err)
	require.NoError(t, os.Remove(filepath.Join(root, "source.link")))
	_, err = store.fileDigestWithIdentity(filepath.Join(root, "source.link"), info, "changed")
	require.Error(t, err)
}

func TestRecentMetadataFastPathUsesChangeTrackingIdentityWithoutMissingEdits(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "source.txt")
	require.NoError(t, os.WriteFile(path, []byte("first"), 0o644))
	original, err := os.Stat(path)
	require.NoError(t, err)
	store, err := OpenCache(root)
	require.NoError(t, err)
	first, err := store.SnapshotFiles("source", "source.txt")
	require.NoError(t, err)
	store.ResetStats()
	second, err := store.SnapshotFiles("source", "source.txt")
	require.NoError(t, err)
	assert.Equal(t, first, second)
	if platformIdentityTracksChanges() {
		assert.Equal(t, 1, store.Stats().DigestsReused)
	}

	require.NoError(t, os.WriteFile(path, []byte("other"), 0o644))
	require.NoError(t, os.Chtimes(path, original.ModTime(), original.ModTime()))
	third, err := store.SnapshotFiles("source", "source.txt")
	require.NoError(t, err)
	assert.NotEqual(t, first, third, "same-size edits with restored mtimes must invalidate through ctime or conservative re-reading")

	recent := cacheStaticFileInfo{modTime: time.Now()}
	assert.False(t, metadataFastPathSafe(recent, ""))
	stable := cacheStaticFileInfo{modTime: time.Now().Add(-metadataFastPathWindow - time.Second)}
	assert.True(t, metadataFastPathSafe(stable, ""))
}

func TestFileIdentityFallbacks(t *testing.T) {
	assert.Empty(t, fileIdentity(cacheStaticFileInfo{}))
	value := 1
	assert.NotEmpty(t, fileIdentity(cacheStaticFileInfo{sys: &value}))
	assert.NotEmpty(t, fileIdentity(cacheStaticFileInfo{sys: cacheFakeSys{Dev: 1, Ino: 2}}))
	assert.NotEmpty(t, fileIdentity(cacheStaticFileInfo{sys: (*int)(nil)}))
	identity, ok := platformFileIdentity(cacheStaticFileInfo{})
	assert.False(t, ok)
	assert.Empty(t, identity)
}

func TestDigestFaultAdapters(t *testing.T) {
	want := errors.New("injected failure")
	store := &Cache{index: indexData{Files: map[string]FileRecord{}, GoAPI: map[string]FileRecord{}}}
	regular := cacheStaticFileInfo{modTime: time.Now()}
	symlink := cacheStaticFileInfo{mode: os.ModeSymlink, modTime: time.Now()}

	_, err := store.fileDigestWithOperations("link", symlink, "identity", cacheDigestOperations{
		readlink: func(string) (string, error) { return "", want },
	})
	require.ErrorIs(t, err, want)

	_, err = store.fileDigestWithOperations("file", regular, "identity", cacheDigestOperations{
		open: func(string) (cacheDigestFile, error) { return nil, want },
	})
	require.ErrorIs(t, err, want)

	_, err = store.fileDigestWithOperations("file", regular, "identity", cacheDigestOperations{
		open: func(string) (cacheDigestFile, error) { return &cacheFaultDigestFile{readErr: want}, nil },
	})
	require.ErrorIs(t, err, want)

	_, err = store.fileDigestWithOperations("file", regular, "identity", cacheDigestOperations{
		open: func(string) (cacheDigestFile, error) { return &cacheFaultDigestFile{closeErr: want}, nil },
	})
	require.ErrorIs(t, err, want)

	root := t.TempDir()
	path := filepath.Join(root, "valid.go")
	require.NoError(t, os.WriteFile(path, []byte("package valid\n"), 0o644))
	info, err := os.Stat(path)
	require.NoError(t, err)
	_, err = store.goAPIDigestWithFormatter(path, info, "identity", func(io.Writer, *token.FileSet, any) error { return want })
	require.ErrorIs(t, err, want)
}

type cacheFaultDigestFile struct {
	readErr, closeErr error
	read              bool
	data              []byte
	offset            int
}

func (f *cacheFaultDigestFile) Read(buffer []byte) (int, error) {
	if !f.read && f.readErr != nil {
		f.read = true
		return 0, f.readErr
	}
	f.read = true
	if f.offset >= len(f.data) {
		return 0, io.EOF
	}
	n := copy(buffer, f.data[f.offset:])
	f.offset += n
	return n, nil
}
func (f *cacheFaultDigestFile) Close() error { return f.closeErr }

type cacheFakeSys struct {
	Dev uint64
	Ino uint64
}

type cacheStaticFileInfo struct {
	name    string
	size    int64
	sys     any
	mode    fs.FileMode
	modTime time.Time
}

func (i cacheStaticFileInfo) Name() string {
	if i.name == "" {
		return "file"
	}
	return i.name
}
func (i cacheStaticFileInfo) Size() int64        { return i.size }
func (i cacheStaticFileInfo) Mode() fs.FileMode  { return i.mode }
func (i cacheStaticFileInfo) ModTime() time.Time { return i.modTime }
func (i cacheStaticFileInfo) IsDir() bool        { return i.mode.IsDir() }
func (i cacheStaticFileInfo) Sys() any           { return i.sys }

func TestSnapshotExcludesConfiguredFileSuffixes(t *testing.T) {
	root := t.TempDir()
	mainPath := filepath.Join(root, "main.go")
	testPath := filepath.Join(root, "main_test.go")
	require.NoError(t, os.WriteFile(mainPath, []byte("package main\n"), 0o644))
	require.NoError(t, os.WriteFile(testPath, []byte("package main\n"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)
	options := SnapshotOptions{Label: "go-build", Root: root, IncludeExtensions: []string{".go"}, ExcludeSuffixes: []string{"_test.go"}}
	first, err := store.Snapshot(options)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(testPath, []byte("package main\n// test-only edit\n"), 0o644))
	second, err := store.Snapshot(options)
	require.NoError(t, err)
	assert.Equal(t, first, second)

	require.NoError(t, os.WriteFile(mainPath, []byte("package main\n// build edit\n"), 0o644))
	third, err := store.Snapshot(options)
	require.NoError(t, err)
	assert.NotEqual(t, second, third)
}

func TestSnapshotHonoursGitIgnoreForDevelopmentInputs(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "ignored"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n"), 0o644))
	ignoredPath := filepath.Join(root, "ignored", "service.go")
	require.NoError(t, os.WriteFile(ignoredPath, []byte("package ignored\nvar Value = 1\n"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)
	options := SnapshotOptions{Label: "dev-go", Root: root, IncludeExtensions: []string{".go"}, UseGitIgnore: true}
	before, err := store.Snapshot(options)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(ignoredPath, []byte("package ignored\nvar Value = 2\n"), 0o644))
	afterIgnoredEdit, err := store.Snapshot(options)
	require.NoError(t, err)
	assert.Equal(t, before, afterIgnoredEdit)

	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nvar Value = 2\n"), 0o644))
	afterSourceEdit, err := store.Snapshot(options)
	require.NoError(t, err)
	assert.NotEqual(t, before, afterSourceEdit)

	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), nil, 0o644))
	afterUnignore, err := store.Snapshot(options)
	require.NoError(t, err)
	assert.NotEqual(t, afterSourceEdit, afterUnignore)

	matcher, err := readGitignoreMatcher(root)
	require.NoError(t, err)
	assert.NotNil(t, matcher)
	_, err = readGitignoreMatcher(filepath.Join(root, "missing"))
	require.Error(t, err)
}

func TestObserveDirectoryTreeCallbackContracts(t *testing.T) {
	root := string(filepath.Separator) + "project"
	want := errors.New("injected failure")
	ignored := gitignore.NewMatcher([]gitignore.Pattern{gitignore.ParsePattern("ignored/", nil), gitignore.ParsePattern("ignored.txt", nil)})

	_, err := observeDirectoryTreeWithWalk(root, map[string]bool{}, nil, func(string) bool { return false }, 2, func(_ *fastwalk.Config, _ string, callback fs.WalkDirFunc) error {
		return callback(root, cacheDirEntry{name: "project", dir: true}, want)
	})
	require.ErrorIs(t, err, want)

	files, err := observeDirectoryTreeWithWalk(root, map[string]bool{"excluded": true, "nested/path": true}, ignored, func(string) bool { return true }, 2, func(_ *fastwalk.Config, _ string, callback fs.WalkDirFunc) error {
		require.NoError(t, callback(root, cacheDirEntry{name: "project", dir: true}, nil))
		assert.ErrorIs(t, callback(filepath.Join(root, "ignored"), cacheDirEntry{name: "ignored", dir: true}, nil), fs.SkipDir)
		require.NoError(t, callback(filepath.Join(root, "ignored.txt"), cacheDirEntry{name: "ignored.txt"}, nil))
		assert.ErrorIs(t, callback(filepath.Join(root, "excluded"), cacheDirEntry{name: "excluded", dir: true}, nil), fs.SkipDir)
		assert.ErrorIs(t, callback(filepath.Join(root, "nested", "path"), cacheDirEntry{name: "path", dir: true}, nil), fs.SkipDir)
		require.NoError(t, callback(filepath.Join(root, "included"), cacheDirEntry{name: "included", dir: true}, nil))
		require.NoError(t, callback(filepath.Join(root, "z.txt"), cacheDirEntry{name: "z.txt"}, nil))
		require.NoError(t, callback(filepath.Join(root, "a.txt"), cacheDirEntry{name: "a.txt"}, nil))
		return nil
	})
	require.NoError(t, err)
	require.Len(t, files, 2)
	assert.Equal(t, "a.txt", files[0].relative)
	assert.Equal(t, "z.txt", files[1].relative)

	_, err = observeDirectoryTreeWithWalk(root, map[string]bool{}, nil, func(string) bool { return true }, 2, func(_ *fastwalk.Config, _ string, callback fs.WalkDirFunc) error {
		return callback(filepath.Join(root, "bad.go"), cacheDirEntry{name: "bad.go", infoErr: want}, nil)
	})
	require.ErrorIs(t, err, want)

	_, err = observeDirectoryTreeWithWalk(root, map[string]bool{}, nil, func(string) bool { return false }, 2, func(*fastwalk.Config, string, fs.WalkDirFunc) error {
		return want
	})
	require.ErrorIs(t, err, want)
}

type cacheDirEntry struct {
	name    string
	dir     bool
	info    fs.FileInfo
	infoErr error
}

func (e cacheDirEntry) Name() string               { return e.name }
func (e cacheDirEntry) IsDir() bool                { return e.dir }
func (e cacheDirEntry) Type() fs.FileMode          { return 0 }
func (e cacheDirEntry) Info() (fs.FileInfo, error) { return e.info, e.infoErr }

func TestObservedFileCachesSuccessfulAndFailedInfoLookups(t *testing.T) {
	info := cacheStaticFileInfo{name: "source", mode: 0o644}
	success := &observedFile{entry: cacheDirEntry{name: "source", info: info}}
	first, err := success.fileInfo()
	require.NoError(t, err)
	second, err := success.fileInfo()
	require.NoError(t, err)
	assert.Equal(t, first, second)

	want := errors.New("injected info failure")
	failure := &observedFile{entry: cacheDirEntry{name: "source", infoErr: want}}
	_, err = failure.fileInfo()
	require.ErrorIs(t, err, want)
	_, err = failure.fileInfo()
	require.ErrorIs(t, err, want)
}

func TestSnapshotsShareOneTreeObservationUntilInvalidated(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/app\n"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)
	store.BeginObservationSession()
	options := SnapshotOptions{Root: root, IncludeExtensions: []string{".go"}, ExcludeDirs: []string{"node_modules"}}

	_, err = store.SnapshotGoAPI(options)
	require.NoError(t, err)
	_, err = store.Snapshot(options)
	require.NoError(t, err)
	assert.Equal(t, 1, store.Stats().TreesWalked)

	store.InvalidateObservations()
	_, err = store.Snapshot(options)
	require.NoError(t, err)
	assert.Equal(t, 2, store.Stats().TreesWalked)
}

func TestObserveTreeFaultsRootFilesAndEagerSelection(t *testing.T) {
	root := t.TempDir()
	want := errors.New("injected failure")
	base := cacheTreeOperations{
		stat:    os.Stat,
		lstat:   os.Lstat,
		ignored: func(string) (gitignore.Matcher, error) { return nil, nil },
		walk: func(string, map[string]bool, gitignore.Matcher, func(string) bool) ([]*observedFile, error) {
			return nil, nil
		},
	}

	store := &Cache{observations: map[string]*treeObservation{}}
	operations := base
	operations.stat = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	_, err := store.observeTreeWithOperations(root, SnapshotOptions{}, false, operations)
	require.ErrorContains(t, err, "does not exist")

	operations = base
	operations.stat = func(string) (fs.FileInfo, error) { return nil, want }
	_, err = store.observeTreeWithOperations(root, SnapshotOptions{}, false, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	_, err = store.observeTreeWithOperations(root, SnapshotOptions{}, false, operations)
	require.ErrorIs(t, err, want)

	file := filepath.Join(root, "input.txt")
	require.NoError(t, os.WriteFile(file, []byte("input"), 0o644))
	store.shareTrees = true
	observation, err := store.observeTreeWithOperations(file, SnapshotOptions{}, false, base)
	require.NoError(t, err)
	require.Len(t, observation.files, 1)
	assert.Equal(t, "input.txt", observation.files[0].relative)

	operations = base
	operations.ignored = func(string) (gitignore.Matcher, error) { return nil, want }
	_, err = store.observeTreeWithOperations(root, SnapshotOptions{UseGitIgnore: true}, false, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.walk = func(string, map[string]bool, gitignore.Matcher, func(string) bool) ([]*observedFile, error) {
		return nil, want
	}
	_, err = store.observeTreeWithOperations(root, SnapshotOptions{}, false, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.walk = func(_ string, excluded map[string]bool, _ gitignore.Matcher, eager func(string) bool) ([]*observedFile, error) {
		assert.True(t, excluded["node_modules"])
		assert.True(t, eager("main.go"))
		assert.False(t, eager("main_test.go"))
		assert.False(t, eager("README.md"))
		return nil, nil
	}
	_, err = store.observeTreeWithOperations(root, SnapshotOptions{ExcludeDirs: []string{"node_modules"}}, true, operations)
	require.NoError(t, err)

	operations.walk = func(_ string, _ map[string]bool, _ gitignore.Matcher, eager func(string) bool) ([]*observedFile, error) {
		assert.True(t, eager("go.mod"))
		assert.True(t, eager("main.go"))
		assert.False(t, eager("main_test.go"))
		assert.False(t, eager("README.md"))
		return nil, nil
	}
	_, err = store.observeTreeWithOperations(root, SnapshotOptions{IncludeNames: []string{"go.mod"}, IncludeExtensions: []string{".go"}, ExcludeSuffixes: []string{"_test.go"}}, false, operations)
	require.NoError(t, err)

	store.shareTrees = true
	store.observations = map[string]*treeObservation{}
	operations.walk = func(string, map[string]bool, gitignore.Matcher, func(string) bool) ([]*observedFile, error) {
		return []*observedFile{}, nil
	}
	first, err := store.observeTreeWithOperations(root, SnapshotOptions{}, false, operations)
	require.NoError(t, err)
	operations.walk = func(string, map[string]bool, gitignore.Matcher, func(string) bool) ([]*observedFile, error) {
		t.Fatal("cached observation walked twice")
		return nil, nil
	}
	second, err := store.observeTreeWithOperations(root, SnapshotOptions{}, false, operations)
	require.NoError(t, err)
	assert.Same(t, first, second)
}

func TestSnapshotAndGoAPIReturnObservationAndDigestFailures(t *testing.T) {
	root := t.TempDir()
	newStore := func(file *observedFile) *Cache {
		observation := &treeObservation{files: []*observedFile{file}}
		return &Cache{
			root: root, shareTrees: true,
			observations: map[string]*treeObservation{
				treeObservationKey(root, SnapshotOptions{IncludeAll: true}): observation,
				treeObservationKey(root, SnapshotOptions{}):                 observation,
			},
			index: indexData{Files: map[string]FileRecord{}, GoAPI: map[string]FileRecord{}},
		}
	}
	want := errors.New("injected failure")
	entry := cacheDirEntry{name: "source.go"}

	store := newStore(&observedFile{path: filepath.Join(root, "source.go"), relative: "source.go", entry: entry, infoErr: want})
	_, err := store.Snapshot(SnapshotOptions{Root: root, IncludeAll: true})
	require.ErrorIs(t, err, want)
	_, err = store.SnapshotGoAPI(SnapshotOptions{Root: root})
	require.ErrorIs(t, err, want)

	path := filepath.Join(root, "removed.go")
	require.NoError(t, os.WriteFile(path, []byte("package removed\n"), 0o644))
	info, err := os.Stat(path)
	require.NoError(t, err)
	require.NoError(t, os.Remove(path))
	store = newStore(&observedFile{path: path, relative: "removed.go", entry: cacheDirEntry{name: "removed.go"}, info: info})
	_, err = store.Snapshot(SnapshotOptions{Root: root, IncludeAll: true})
	require.Error(t, err)

	malformed := filepath.Join(root, "malformed.go")
	require.NoError(t, os.WriteFile(malformed, []byte("package"), 0o644))
	store, err = OpenCache(root)
	require.NoError(t, err)
	_, err = store.SnapshotGoAPI(SnapshotOptions{Root: root})
	require.ErrorContains(t, err, "parse Go API")

	_, err = store.Snapshot(SnapshotOptions{Root: filepath.Join(root, "missing"), IncludeAll: true})
	require.Error(t, err)
	_, err = store.SnapshotGoAPI(SnapshotOptions{Root: filepath.Join(root, "missing")})
	require.Error(t, err)
}

func TestSaveSkipsAStableIndex(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "source.txt")
	require.NoError(t, os.WriteFile(path, []byte("content"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)

	require.NoError(t, store.Save())
	assert.NoFileExists(t, store.indexPath)

	_, err = store.SnapshotFiles("source", "source.txt")
	require.NoError(t, err)
	require.NoError(t, store.Save())
	assert.FileExists(t, store.indexPath)

	store.indexPath = filepath.Join(root, "missing", "index.json")
	require.NoError(t, store.Save(), "a clean cache must not touch the filesystem")
	require.NoError(t, os.WriteFile(path, []byte("changed"), 0o644))
	_, err = store.SnapshotFiles("source", "source.txt")
	require.NoError(t, err)
	require.Error(t, store.Save(), "a changed cache must persist its new records")
}
