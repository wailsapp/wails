package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheRejectsDigestsFromOtherAlgorithms(t *testing.T) {
	for _, algorithm := range []string{"", "blake3", "unknown", cacheHashAlgorithm} {
		t.Run(algorithm, func(t *testing.T) {
			root := t.TempDir()
			state := filepath.Join(root, ".wails")
			if err := os.Mkdir(state, 0755); err != nil {
				t.Fatal(err)
			}
			old := indexData{
				HashAlgorithm: algorithm,
				Files:         map[string]FileRecord{"source": {Digest: "cached-file"}},
				GoAPI:         map[string]FileRecord{"source": {Digest: "cached-api"}},
				Actions:       map[string]ActionResult{"build": {Artifact: "cached-artifact"}},
				Receipts:      map[string]time.Time{"install": time.Now()},
			}
			var encoded bytes.Buffer
			if err := gob.NewEncoder(&encoded).Encode(old); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(state, "wake-index-v2.gob"), encoded.Bytes(), 0644); err != nil {
				t.Fatal(err)
			}
			store, err := OpenCache(root)
			if err != nil {
				t.Fatal(err)
			}
			want := 0
			if algorithm == cacheHashAlgorithm {
				want = 1
			}
			for name, count := range map[string]int{"files": len(store.index.Files), "api": len(store.index.GoAPI), "actions": len(store.index.Actions), "receipts": len(store.index.Receipts)} {
				if count != want {
					t.Errorf("%s entries = %d, want %d", name, count, want)
				}
			}
			if store.index.HashAlgorithm != cacheHashAlgorithm {
				t.Fatalf("algorithm = %q", store.index.HashAlgorithm)
			}
			if err := store.Save(); err != nil {
				t.Fatal(err)
			}
			reopened, err := OpenCache(root)
			if err != nil || reopened.index.HashAlgorithm != cacheHashAlgorithm {
				t.Fatalf("reopening migrated index: %v", err)
			}
		})
	}
}

func TestFileDigestUsesSHA256(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "source")
	data := []byte("Wake cache algorithm regression\n")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	store, err := OpenCache(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := store.fileDigest(path, info)
	want := fmt.Sprintf("%x", sha256.Sum256(data))
	if err != nil || got != want {
		t.Fatalf("file digest = %q, %v; want %q", got, err, want)
	}
}

func BenchmarkCacheSHA256(b *testing.B) {
	for _, size := range []int{1024, 65536, 4 << 20} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			data := bytes.Repeat([]byte{0x5a}, size)
			hash := sha256.New()
			digest := make([]byte, 0, sha256.Size)
			b.SetBytes(int64(size))
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				hash.Reset()
				_, _ = hash.Write(data)
				digest = hash.Sum(digest[:0])
			}
			if len(digest) != sha256.Size {
				b.Fatal("missing digest")
			}
		})
	}
}
