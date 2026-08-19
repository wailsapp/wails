package cache

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func BenchmarkSharedSemanticAndContentSnapshot(b *testing.B) {
	root := b.TempDir()
	stableTime := time.Now().Add(-metadataFastPathWindow - time.Second)
	for packageIndex := range 100 {
		directory := filepath.Join(root, fmt.Sprintf("package-%03d", packageIndex))
		require.NoError(b, os.MkdirAll(directory, 0o755))
		for fileIndex := range 10 {
			path := filepath.Join(directory, fmt.Sprintf("source-%02d.go", fileIndex))
			source := fmt.Appendf(nil, "package package%03d\n\nfunc value%02d() int { return %d }\n", packageIndex, fileIndex, fileIndex)
			require.NoError(b, os.WriteFile(path, source, 0o644))
			require.NoError(b, os.Chtimes(path, stableTime, stableTime))
		}
	}
	store, err := OpenCache(root)
	require.NoError(b, err)
	store.BeginObservationSession()
	options := SnapshotOptions{Label: "local-source", Root: root, IncludeNames: []string{"go.mod", "go.sum"}, IncludeExtensions: []string{".go"}, ExcludeSuffixes: []string{"_test.go"}}
	store.InvalidateObservations()
	_, err = store.SnapshotGoAPI(options)
	require.NoError(b, err)
	_, err = store.Snapshot(options)
	require.NoError(b, err)
	require.NoError(b, store.Save())

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		store.InvalidateObservations()
		if _, err := store.SnapshotGoAPI(options); err != nil {
			b.Fatal(err)
		}
		if _, err := store.Snapshot(options); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArtifactCacheLifecycle(b *testing.B) {
	root := b.TempDir()
	b.Setenv("XDG_CACHE_HOME", filepath.Join(root, "machine-cache"))
	require.NoError(b, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	output := filepath.Join("bin", "app")
	source := filepath.Join(root, output)
	payload := bytes.Repeat([]byte{0xa5}, 1<<20)
	require.NoError(b, os.WriteFile(source, payload, 0o755))
	store, err := OpenCache(root)
	require.NoError(b, err)
	digest, err := store.RecordAction("build", output)
	require.NoError(b, err)

	b.Run("lookup-hit", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(payload)))
		for range b.N {
			status, _, err := store.Lookup("build", output)
			if err != nil || status != LookupHit {
				b.Fatalf("lookup = %s, %v", status, err)
			}
		}
	})
	b.Run("store", func(b *testing.B) {
		artifact := filepath.Join(store.artifactRoot, digest)
		b.ReportAllocs()
		b.SetBytes(int64(len(payload)))
		for range b.N {
			b.StopTimer()
			if err := os.RemoveAll(artifact); err != nil {
				b.Fatal(err)
			}
			b.StartTimer()
			if err := store.storeArtifact(digest, source); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("restore", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(payload)))
		for range b.N {
			b.StopTimer()
			if err := os.Remove(source); err != nil {
				b.Fatal(err)
			}
			b.StartTimer()
			status, _, err := store.Lookup("build", output)
			if err != nil || status != LookupRestored {
				b.Fatalf("lookup = %s, %v", status, err)
			}
		}
	})
}
