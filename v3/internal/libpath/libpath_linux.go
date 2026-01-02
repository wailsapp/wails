//go:build linux

package libpath

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

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

// initCache populates the cache with dynamic paths from package managers.
// This is called lazily on first access.
func (c *pathCache) init() {
	c.initOnce.Do(func() {
		c.flatpak = discoverFlatpakLibPaths()
		c.snap = discoverSnapLibPaths()
		c.nix = discoverNixLibPaths()
		c.inited = true
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

// Common library search paths on Linux systems
var defaultLibPaths = []string{
	// Standard paths
	"/usr/lib",
	"/usr/lib64",
	"/lib",
	"/lib64",

	// Debian/Ubuntu multiarch
	"/usr/lib/x86_64-linux-gnu",
	"/usr/lib/aarch64-linux-gnu",
	"/usr/lib/i386-linux-gnu",
	"/usr/lib/arm-linux-gnueabihf",
	"/lib/x86_64-linux-gnu",
	"/lib/aarch64-linux-gnu",

	// Fedora/RHEL/CentOS
	"/usr/lib64/gtk-3.0",
	"/usr/lib64/gtk-4.0",
	"/usr/lib/gcc/x86_64-redhat-linux",
	"/usr/lib/gcc/aarch64-redhat-linux",

	// Arch Linux
	"/usr/lib/webkit2gtk-4.0",
	"/usr/lib/webkit2gtk-4.1",
	"/usr/lib/gtk-3.0",
	"/usr/lib/gtk-4.0",

	// openSUSE
	"/usr/lib64/gcc/x86_64-suse-linux",

	// Local installations
	"/usr/local/lib",
	"/usr/local/lib64",
}

// getFlatpakLibPaths returns cached library paths from installed Flatpak runtimes.
func getFlatpakLibPaths() []string {
	return cache.getFlatpak()
}

// discoverFlatpakLibPaths scans for Flatpak runtime library directories.
// Uses `flatpak --installations` and scans for runtime lib directories.
func discoverFlatpakLibPaths() []string {
	var paths []string

	// Get system and user installation directories
	installDirs := []string{
		"/var/lib/flatpak",                         // System default
		os.ExpandEnv("$HOME/.local/share/flatpak"), // User default
	}

	// Try to get actual installation path from flatpak
	if out, err := exec.Command("flatpak", "--installations").Output(); err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line != "" {
				installDirs = append(installDirs, line)
			}
		}
	}

	// Scan for runtime lib directories
	for _, installDir := range installDirs {
		runtimeDir := filepath.Join(installDir, "runtime")
		if _, err := os.Stat(runtimeDir); err != nil {
			continue
		}

		// Look for lib directories in runtimes
		// Structure: runtime/<name>/<arch>/<version>/<hash>/files/lib
		matches, err := filepath.Glob(filepath.Join(runtimeDir, "*", "*", "*", "*", "files", "lib"))
		if err == nil {
			paths = append(paths, matches...)
		}
	}

	return paths
}

// getSnapLibPaths returns cached library paths from installed Snap packages.
func getSnapLibPaths() []string {
	return cache.getSnap()
}

// discoverSnapLibPaths scans for Snap package library directories.
// Scans /snap/*/current/usr/lib* directories.
func discoverSnapLibPaths() []string {
	var paths []string

	snapDir := "/snap"
	if _, err := os.Stat(snapDir); err != nil {
		return paths
	}

	// Find all snap packages with lib directories
	patterns := []string{
		filepath.Join(snapDir, "*", "current", "usr", "lib"),
		filepath.Join(snapDir, "*", "current", "usr", "lib64"),
		filepath.Join(snapDir, "*", "current", "usr", "lib", "*-linux-gnu"),
		filepath.Join(snapDir, "*", "current", "lib"),
		filepath.Join(snapDir, "*", "current", "lib", "*-linux-gnu"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err == nil {
			paths = append(paths, matches...)
		}
	}

	return paths
}

// getNixLibPaths returns cached library paths for Nix/NixOS installations.
func getNixLibPaths() []string {
	return cache.getNix()
}

// discoverNixLibPaths scans for Nix library paths.
func discoverNixLibPaths() []string {
	var paths []string

	nixProfileLib := os.ExpandEnv("$HOME/.nix-profile/lib")
	if _, err := os.Stat(nixProfileLib); err == nil {
		paths = append(paths, nixProfileLib)
	}

	// System Nix store - packages expose libs through profiles
	nixStoreLib := "/run/current-system/sw/lib"
	if _, err := os.Stat(nixStoreLib); err == nil {
		paths = append(paths, nixStoreLib)
	}

	return paths
}

// searchResult holds the result from a parallel search goroutine.
type searchResult struct {
	path   string
	source string // for debugging: "pkg-config", "ldconfig", "filesystem"
}

// FindLibraryPath attempts to find the path to a library using multiple methods
// in parallel. It searches via pkg-config, ldconfig, and filesystem simultaneously,
// returning as soon as any method finds the library.
//
// The libName should be the pkg-config name (e.g., "gtk+-3.0", "webkit2gtk-4.1").
// Returns the library directory path and any error encountered.
func FindLibraryPath(libName string) (string, error) {
	// Create a context that we'll cancel once we find a result
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to receive results - buffered to avoid goroutine leaks
	results := make(chan searchResult, 3)

	// Launch parallel searches
	var wg sync.WaitGroup
	wg.Add(3)

	// Search via pkg-config
	go func() {
		defer wg.Done()
		if path, err := findWithPkgConfigCtx(ctx, libName); err == nil {
			select {
			case results <- searchResult{path: path, source: "pkg-config"}:
			case <-ctx.Done():
			}
		}
	}()

	// Search via ldconfig
	go func() {
		defer wg.Done()
		if path, err := findWithLdconfigCtx(ctx, libName); err == nil {
			select {
			case results <- searchResult{path: path, source: "ldconfig"}:
			case <-ctx.Done():
			}
		}
	}()

	// Search via filesystem
	go func() {
		defer wg.Done()
		if path, err := findInCommonPathsCtx(ctx, libName); err == nil {
			select {
			case results <- searchResult{path: path, source: "filesystem"}:
			case <-ctx.Done():
			}
		}
	}()

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Return first result or error if none found
	if result, ok := <-results; ok {
		return result.path, nil
	}

	return "", &LibraryNotFoundError{Name: libName}
}

// FindLibraryPathSequential is the original sequential implementation.
// Use this if you need deterministic search order (pkg-config → ldconfig → filesystem).
func FindLibraryPathSequential(libName string) (string, error) {
	// Try pkg-config first (most reliable when available)
	if path, err := findWithPkgConfig(libName); err == nil {
		return path, nil
	}

	// Try ldconfig cache
	if path, err := findWithLdconfig(libName); err == nil {
		return path, nil
	}

	// Fall back to searching common paths
	return findInCommonPaths(libName)
}

// FindLibraryFile finds the full path to a specific library file (e.g., "libgtk-3.so").
func FindLibraryFile(fileName string) (string, error) {
	// Try ldconfig first
	if path, err := findFileWithLdconfig(fileName); err == nil {
		return path, nil
	}

	// Search all paths including dynamic ones
	for _, dir := range GetAllLibPaths() {
		// Check exact match
		fullPath := filepath.Join(dir, fileName)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}

		// Check with .so suffix variations
		matches, err := filepath.Glob(filepath.Join(dir, fileName+"*"))
		if err == nil && len(matches) > 0 {
			return matches[0], nil
		}
	}

	return "", &LibraryNotFoundError{Name: fileName}
}

// findWithPkgConfig uses pkg-config to find library paths.
func findWithPkgConfig(libName string) (string, error) {
	return findWithPkgConfigCtx(context.Background(), libName)
}

// findWithPkgConfigCtx uses pkg-config to find library paths with context support.
func findWithPkgConfigCtx(ctx context.Context, libName string) (string, error) {
	// Check if already cancelled
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	cmd := exec.CommandContext(ctx, "pkg-config", "--libs-only-L", libName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse -L flags from output
	parts := strings.Fields(string(output))
	for _, part := range parts {
		if strings.HasPrefix(part, "-L") {
			path := strings.TrimPrefix(part, "-L")
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}

	// Check context before second command
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// If no -L flag, try --variable=libdir
	cmd = exec.CommandContext(ctx, "pkg-config", "--variable=libdir", libName)
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	path := strings.TrimSpace(string(output))
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", &LibraryNotFoundError{Name: libName}
}

// findWithLdconfig searches the ldconfig cache for library paths.
func findWithLdconfig(libName string) (string, error) {
	return findWithLdconfigCtx(context.Background(), libName)
}

// findWithLdconfigCtx searches the ldconfig cache for library paths with context support.
func findWithLdconfigCtx(ctx context.Context, libName string) (string, error) {
	// Check if already cancelled
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Convert pkg-config name to library name pattern
	// e.g., "gtk+-3.0" -> "libgtk-3", "webkit2gtk-4.1" -> "libwebkit2gtk-4.1"
	searchName := pkgConfigToLibName(libName)

	cmd := exec.CommandContext(ctx, "ldconfig", "-p")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, searchName) {
			// Line format: "	libname.so.X (libc6,x86-64) => /path/to/lib"
			parts := strings.Split(line, "=>")
			if len(parts) == 2 {
				libPath := strings.TrimSpace(parts[1])
				return filepath.Dir(libPath), nil
			}
		}
	}

	return "", &LibraryNotFoundError{Name: libName}
}

// findFileWithLdconfig finds a specific library file using ldconfig.
func findFileWithLdconfig(fileName string) (string, error) {
	cmd := exec.Command("ldconfig", "-p")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	baseName := strings.TrimSuffix(fileName, ".so")
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, baseName) {
			parts := strings.Split(line, "=>")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", &LibraryNotFoundError{Name: fileName}
}

// findInCommonPaths searches common library directories including
// dynamically discovered Flatpak, Snap, and Nix paths.
func findInCommonPaths(libName string) (string, error) {
	return findInCommonPathsCtx(context.Background(), libName)
}

// findInCommonPathsCtx searches common library directories with context support.
func findInCommonPathsCtx(ctx context.Context, libName string) (string, error) {
	searchName := pkgConfigToLibName(libName)

	// Search all paths including dynamic ones
	allPaths := GetAllLibPaths()

	for _, dir := range allPaths {
		// Check if cancelled periodically
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		if _, err := os.Stat(dir); err != nil {
			continue
		}

		// Look for the library file
		pattern := filepath.Join(dir, searchName+"*.so*")
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return dir, nil
		}

		// Also check pkgconfig subdirectory for .pc files
		pcPath := filepath.Join(dir, "pkgconfig", libName+".pc")
		if _, err := os.Stat(pcPath); err == nil {
			return dir, nil
		}
	}

	return "", &LibraryNotFoundError{Name: libName}
}

// pkgConfigToLibName converts a pkg-config package name to a library name pattern.
func pkgConfigToLibName(pkgName string) string {
	// Common transformations
	name := pkgName

	// Remove version suffix like "-3.0", "-4.1"
	// but keep it for webkit2gtk-4.1 style names
	if strings.HasPrefix(name, "gtk+-") {
		// gtk+-3.0 -> libgtk-3
		name = "libgtk-" + strings.TrimPrefix(name, "gtk+-")
		name = strings.Split(name, ".")[0]
	} else if strings.HasPrefix(name, "webkit2gtk-") {
		// webkit2gtk-4.1 -> libwebkit2gtk-4.1
		name = "lib" + name
	} else if !strings.HasPrefix(name, "lib") {
		name = "lib" + name
	}

	return name
}

// GetAllLibPaths returns all library paths from LD_LIBRARY_PATH, default paths,
// and dynamically discovered paths from Flatpak, Snap, and Nix.
// It does NOT include the current directory for security reasons.
func GetAllLibPaths() []string {
	var paths []string

	// Add LD_LIBRARY_PATH entries first (highest priority)
	if ldPath := os.Getenv("LD_LIBRARY_PATH"); ldPath != "" {
		for _, p := range strings.Split(ldPath, ":") {
			if p != "" {
				paths = append(paths, p)
			}
		}
	}

	// Add default system paths
	paths = append(paths, defaultLibPaths...)

	// Add dynamically discovered paths from package managers
	paths = append(paths, getFlatpakLibPaths()...)
	paths = append(paths, getSnapLibPaths()...)
	paths = append(paths, getNixLibPaths()...)

	return paths
}

// FindLibraryPathWithOptions finds a library path with additional search options.
type FindOptions struct {
	// IncludeCurrentDir includes "." in the search path.
	// WARNING: This is a security risk and should only be used for development.
	IncludeCurrentDir bool

	// ExtraPaths are additional paths to search before the defaults.
	ExtraPaths []string
}

// FindLibraryPathWithOptions attempts to find the path to a library with custom options.
func FindLibraryPathWithOptions(libName string, opts FindOptions) (string, error) {
	// Try pkg-config first (most reliable when available)
	if path, err := findWithPkgConfig(libName); err == nil {
		return path, nil
	}

	// Try ldconfig cache
	if path, err := findWithLdconfig(libName); err == nil {
		return path, nil
	}

	// Build search paths
	searchPaths := make([]string, 0, len(opts.ExtraPaths)+len(defaultLibPaths)+1)

	if opts.IncludeCurrentDir {
		if cwd, err := os.Getwd(); err == nil {
			searchPaths = append(searchPaths, cwd)
		}
	}

	searchPaths = append(searchPaths, opts.ExtraPaths...)
	searchPaths = append(searchPaths, defaultLibPaths...)

	// Search the paths
	searchName := pkgConfigToLibName(libName)
	for _, dir := range searchPaths {
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		pattern := filepath.Join(dir, searchName+"*.so*")
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return dir, nil
		}

		pcPath := filepath.Join(dir, "pkgconfig", libName+".pc")
		if _, err := os.Stat(pcPath); err == nil {
			return dir, nil
		}
	}

	return "", &LibraryNotFoundError{Name: libName}
}

// LibraryNotFoundError is returned when a library cannot be found.
type LibraryNotFoundError struct {
	Name string
}

func (e *LibraryNotFoundError) Error() string {
	return "library not found: " + e.Name
}

// LibraryMatch holds information about a found library.
type LibraryMatch struct {
	// Name is the pkg-config name that was searched for.
	Name string
	// Path is the directory containing the library.
	Path string
}

// FindFirstLibrary searches for multiple libraries in parallel and returns
// the first one found. This is useful when you don't know the exact version
// of a library installed (e.g., gtk+-3.0 vs gtk+-4.0).
//
// The search order among candidates is non-deterministic - whichever is found
// first wins. If you need a specific preference order, list preferred libraries
// first and use FindFirstLibraryOrdered instead.
//
// Example:
//
//	match, err := FindFirstLibrary("webkit2gtk-4.1", "webkit2gtk-4.0", "webkit2gtk-6.0")
//	if err != nil {
//	    log.Fatal("No WebKit2GTK found")
//	}
//	fmt.Printf("Found %s at %s\n", match.Name, match.Path)
func FindFirstLibrary(libNames ...string) (*LibraryMatch, error) {
	if len(libNames) == 0 {
		return nil, &LibraryNotFoundError{Name: "no libraries specified"}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan *LibraryMatch, len(libNames))
	var wg sync.WaitGroup

	for _, name := range libNames {
		wg.Add(1)
		go func(libName string) {
			defer wg.Done()
			if path, err := findLibraryPathCtx(ctx, libName); err == nil {
				select {
				case results <- &LibraryMatch{Name: libName, Path: path}:
				case <-ctx.Done():
				}
			}
		}(name)
	}

	// Close results when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	if result := <-results; result != nil {
		return result, nil
	}

	return nil, &LibraryNotFoundError{Name: strings.Join(libNames, ", ")}
}

// FindFirstLibraryOrdered searches for libraries in order of preference,
// returning the first one found. Unlike FindFirstLibrary, this respects
// the order of candidates - earlier entries are preferred.
//
// This is useful when you want to prefer newer library versions:
//
//	match, err := FindFirstLibraryOrdered("gtk+-4.0", "gtk+-3.0")
//	// Will return gtk+-4.0 if available, otherwise gtk+-3.0
func FindFirstLibraryOrdered(libNames ...string) (*LibraryMatch, error) {
	if len(libNames) == 0 {
		return nil, &LibraryNotFoundError{Name: "no libraries specified"}
	}

	for _, name := range libNames {
		if path, err := FindLibraryPath(name); err == nil {
			return &LibraryMatch{Name: name, Path: path}, nil
		}
	}

	return nil, &LibraryNotFoundError{Name: strings.Join(libNames, ", ")}
}

// FindAllLibraries searches for multiple libraries in parallel and returns
// all that are found. This is useful for discovering which library versions
// are available on the system.
//
// Example:
//
//	matches := FindAllLibraries("gtk+-3.0", "gtk+-4.0", "webkit2gtk-4.0", "webkit2gtk-4.1")
//	for _, m := range matches {
//	    fmt.Printf("Found %s at %s\n", m.Name, m.Path)
//	}
func FindAllLibraries(libNames ...string) []LibraryMatch {
	if len(libNames) == 0 {
		return nil
	}

	results := make(chan *LibraryMatch, len(libNames))
	var wg sync.WaitGroup

	for _, name := range libNames {
		wg.Add(1)
		go func(libName string) {
			defer wg.Done()
			if path, err := FindLibraryPath(libName); err == nil {
				results <- &LibraryMatch{Name: libName, Path: path}
			}
		}(name)
	}

	// Close results when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	var matches []LibraryMatch
	for result := range results {
		matches = append(matches, *result)
	}

	return matches
}

// findLibraryPathCtx is FindLibraryPath with context support.
func findLibraryPathCtx(ctx context.Context, libName string) (string, error) {
	// Create a child context for this search
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan searchResult, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		if path, err := findWithPkgConfigCtx(ctx, libName); err == nil {
			select {
			case results <- searchResult{path: path, source: "pkg-config"}:
			case <-ctx.Done():
			}
		}
	}()

	go func() {
		defer wg.Done()
		if path, err := findWithLdconfigCtx(ctx, libName); err == nil {
			select {
			case results <- searchResult{path: path, source: "ldconfig"}:
			case <-ctx.Done():
			}
		}
	}()

	go func() {
		defer wg.Done()
		if path, err := findInCommonPathsCtx(ctx, libName); err == nil {
			select {
			case results <- searchResult{path: path, source: "filesystem"}:
			case <-ctx.Done():
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	if result, ok := <-results; ok {
		return result.path, nil
	}

	return "", &LibraryNotFoundError{Name: libName}
}
