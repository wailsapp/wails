//go:build linux

package libpath

import (
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
)

func TestPathCacheConcurrentInvalidation(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("PATH", t.TempDir())
	c := &pathCache{}
	c.getNix()
	start := make(chan struct{})
	var workers sync.WaitGroup
	for i := 0; i < 8; i++ {
		workers.Go(func() {
			<-start
			for j := 0; j < 50; j++ {
				c.getFlatpak()
				c.getSnap()
				c.getNix()
			}
		})
	}
	workers.Go(func() {
		<-start
		for j := 0; j < 50; j++ {
			c.invalidate()
		}
	})
	close(start)
	workers.Wait()
}

func TestPathCacheInvalidationRefreshesPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", t.TempDir())
	c := &pathCache{}
	want := filepath.Join(home, ".nix-profile", "lib")
	if slices.Contains(c.getNix(), want) {
		t.Fatal("cache reported a path before it existed")
	}
	if err := os.MkdirAll(want, 0755); err != nil {
		t.Fatal(err)
	}
	if slices.Contains(c.getNix(), want) {
		t.Fatal("cache unexpectedly rediscovered paths without invalidation")
	}
	c.invalidate()
	if !slices.Contains(c.getNix(), want) {
		t.Fatal("invalidation did not refresh library paths")
	}
}
