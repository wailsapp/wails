//go:build linux

package libpath

import "sync"

// pathCache holds cached dynamic library paths to avoid repeated
// expensive filesystem and subprocess operations.
type pathCache struct {
	mu       sync.RWMutex
	flatpak  []string
	snap     []string
	nix      []string
	initOnce sync.Once
	inited   bool
}

var cache pathCache

// init populates the cache with dynamic paths from package managers.
// This is called lazily on first access.
func (c *pathCache) init() {
	c.initOnce.Do(func() {
		// Discover paths without holding the lock
		flatpak := discoverFlatpakLibPaths()
		snap := discoverSnapLibPaths()
		nix := discoverNixLibPaths()

		// Hold lock only while updating the cache
		c.mu.Lock()
		c.flatpak = flatpak
		c.snap = snap
		c.nix = nix
		c.inited = true
		c.mu.Unlock()
	})
}

// getFlatpak returns cached Flatpak library paths.
func (c *pathCache) getFlatpak() []string {
	c.init()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.flatpak
}

// getSnap returns cached Snap library paths.
func (c *pathCache) getSnap() []string {
	c.init()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.snap
}

// getNix returns cached Nix library paths.
func (c *pathCache) getNix() []string {
	c.init()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nix
}

// invalidate clears the cache and forces re-discovery on next access.
func (c *pathCache) invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flatpak = nil
	c.snap = nil
	c.nix = nil
	c.initOnce = sync.Once{} // Reset so init() runs again
	c.inited = false
}

// InvalidateCache clears the cached dynamic library paths.
// Call this if packages are installed or removed during runtime
// and you need to re-discover library paths.
func InvalidateCache() {
	cache.invalidate()
}
