// Package libpath provides utilities for finding native library paths on Linux.
//
// # Overview
//
// This package helps locate shared libraries (.so files) on Linux systems,
// supporting multiple distributions and package managers. It's particularly
// useful for applications that need to link against libraries like GTK,
// WebKit2GTK, or other system libraries at runtime.
//
// # Search Strategy
//
// The package uses a multi-tier search strategy, trying each method in order
// until a library is found:
//
//  1. pkg-config: Queries the pkg-config database for library paths
//  2. ldconfig: Searches the dynamic linker cache
//  3. Filesystem: Scans common library directories
//
// # Supported Distributions
//
// The package includes default search paths for:
//
//   - Debian/Ubuntu (multiarch paths like /usr/lib/x86_64-linux-gnu)
//   - Fedora/RHEL/CentOS (/usr/lib64, /usr/lib64/gtk-*)
//   - Arch Linux (/usr/lib/webkit2gtk-*, /usr/lib/gtk-*)
//   - openSUSE (/usr/lib64/gcc/x86_64-suse-linux)
//   - NixOS and Nix package manager
//
// # Package Manager Support
//
// Dynamic paths are discovered from:
//
//   - Flatpak: Scans runtime directories via `flatpak --installations`
//   - Snap: Globs /snap/*/current/usr/lib* directories
//   - Nix: Checks ~/.nix-profile/lib and /run/current-system/sw/lib
//
// # Caching
//
// Dynamic path discovery (Flatpak, Snap, Nix) is cached for performance.
// The cache is populated on first access and persists for the process lifetime.
// Use [InvalidateCache] to force re-discovery if packages are installed/removed
// during runtime.
//
// # Security
//
// The current directory (".") is never included in search paths by default,
// as this is a security risk. Use [FindLibraryPathWithOptions] with
// IncludeCurrentDir if you explicitly need this behavior (not recommended
// for production).
//
// # Performance
//
// Typical lookup times (cached):
//
//   - Found via pkg-config: ~2ms (spawns external process)
//   - Found via ldconfig: ~1.3ms (spawns external process)
//   - Found via filesystem: ~0.1ms (uses cached paths)
//   - Not found (worst case): ~20ms (searches all paths)
//
// # Example Usage
//
//	// Find a library by its pkg-config name
//	path, err := libpath.FindLibraryPath("webkit2gtk-4.1")
//	if err != nil {
//	    log.Fatal("WebKit2GTK not found:", err)
//	}
//	fmt.Println("Found at:", path)
//
//	// Find a specific .so file
//	soPath, err := libpath.FindLibraryFile("libgtk-3.so")
//	if err != nil {
//	    log.Fatal("GTK3 library file not found:", err)
//	}
//	fmt.Println("Library file:", soPath)
//
//	// Get all library search paths
//	for _, p := range libpath.GetAllLibPaths() {
//	    fmt.Println(p)
//	}
//
// # Multi-Library Search
//
// When you don't know which version of a library is installed, use the
// multi-library search functions:
//
//	// Find any available WebKit2GTK version (first found wins)
//	match, err := libpath.FindFirstLibrary("webkit2gtk-4.1", "webkit2gtk-4.0", "webkit2gtk-6.0")
//	if err != nil {
//	    log.Fatal("No WebKit2GTK found")
//	}
//	fmt.Printf("Found %s at %s\n", match.Name, match.Path)
//
//	// Prefer newer versions (ordered search)
//	match, err := libpath.FindFirstLibraryOrdered("gtk4", "gtk+-3.0")
//
//	// Discover all available versions
//	matches := libpath.FindAllLibraries("gtk+-3.0", "gtk4", "webkit2gtk-4.0", "webkit2gtk-4.1")
//	for _, m := range matches {
//	    fmt.Printf("Available: %s at %s\n", m.Name, m.Path)
//	}
//
// On non-Linux platforms, stub implementations are provided that always
// return [LibraryNotFoundError].
package libpath
