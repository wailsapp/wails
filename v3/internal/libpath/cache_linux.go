//go:build linux

package libpath

import "sync"

// pathCache holds cached dynamic library paths to avoid repeated
// expensive filesystem and subprocess operations.
type pathCache struct {
	mu      sync.Mutex
	flatpak []string
	snap    []string
	nix     []string
	inited  bool
}

var cache pathCache

// initLocked populates the cache on first access. The caller must hold c.mu.
func (c *pathCache) initLocked() {
	if c.inited {
		return
	}
	c.flatpak = discoverFlatpakLibPaths()
	c.snap = discoverSnapLibPaths()
	c.nix = discoverNixLibPaths()
	c.inited = true
}

// getFlatpak returns cached Flatpak library paths.
func (c *pathCache) getFlatpak() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.initLocked()
	return c.flatpak
}

// getSnap returns cached Snap library paths.
func (c *pathCache) getSnap() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.initLocked()
	return c.snap
}

// getNix returns cached Nix library paths.
func (c *pathCache) getNix() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.initLocked()
	return c.nix
}

// invalidate clears the cache and forces re-discovery on next access.
func (c *pathCache) invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flatpak = nil
	c.snap = nil
	c.nix = nil
	c.inited = false
}

// InvalidateCache clears the cached dynamic library paths.
// Call this if packages are installed or removed during runtime
// and you need to re-discover library paths.
func InvalidateCache() {
	cache.invalidate()
}
