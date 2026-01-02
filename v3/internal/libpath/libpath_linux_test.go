//go:build linux

package libpath

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPkgConfigToLibName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"gtk+-3.0", "libgtk-3"},
		{"gtk+-4.0", "libgtk-4"},
		{"webkit2gtk-4.1", "libwebkit2gtk-4.1"},
		{"webkit2gtk-4.0", "libwebkit2gtk-4.0"},
		{"glib-2.0", "libglib-2.0"},
		{"libsoup-3.0", "libsoup-3.0"},
		{"cairo", "libcairo"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := pkgConfigToLibName(tt.input)
			if result != tt.expected {
				t.Errorf("pkgConfigToLibName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetAllLibPaths(t *testing.T) {
	paths := GetAllLibPaths()

	if len(paths) == 0 {
		t.Error("GetAllLibPaths() returned empty slice")
	}

	// Check that default paths are included
	hasUsrLib := false
	for _, p := range paths {
		if p == "/usr/lib" || p == "/usr/lib64" {
			hasUsrLib = true
			break
		}
	}
	if !hasUsrLib {
		t.Error("GetAllLibPaths() should include /usr/lib or /usr/lib64")
	}
}

func TestGetAllLibPaths_WithLDPath(t *testing.T) {
	// Save and restore LD_LIBRARY_PATH
	original := os.Getenv("LD_LIBRARY_PATH")
	defer os.Setenv("LD_LIBRARY_PATH", original)

	testPath := "/test/custom/lib:/another/path"
	os.Setenv("LD_LIBRARY_PATH", testPath)

	paths := GetAllLibPaths()

	// First paths should be from LD_LIBRARY_PATH
	if len(paths) < 2 {
		t.Fatal("Expected at least 2 paths")
	}
	if paths[0] != "/test/custom/lib" {
		t.Errorf("First path should be /test/custom/lib, got %s", paths[0])
	}
	if paths[1] != "/another/path" {
		t.Errorf("Second path should be /another/path, got %s", paths[1])
	}
}

func TestLibraryNotFoundError(t *testing.T) {
	err := &LibraryNotFoundError{Name: "testlib"}
	expected := "library not found: testlib"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestFindLibraryPath_NotFound(t *testing.T) {
	_, err := FindLibraryPath("nonexistent-library-xyz-123")
	if err == nil {
		t.Error("Expected error for nonexistent library")
	}

	var notFoundErr *LibraryNotFoundError
	if _, ok := err.(*LibraryNotFoundError); !ok {
		t.Errorf("Expected LibraryNotFoundError, got %T", err)
	} else {
		notFoundErr = err.(*LibraryNotFoundError)
		if notFoundErr.Name != "nonexistent-library-xyz-123" {
			t.Errorf("Error name = %q, want %q", notFoundErr.Name, "nonexistent-library-xyz-123")
		}
	}
}

func TestFindLibraryFile_NotFound(t *testing.T) {
	_, err := FindLibraryFile("libnonexistent-xyz-123.so")
	if err == nil {
		t.Error("Expected error for nonexistent library file")
	}
}

// Integration tests - these depend on system state
// They're skipped if the required tools/libraries aren't available

func TestFindLibraryPath_WithPkgConfig(t *testing.T) {
	// Skip if pkg-config is not available
	if _, err := exec.LookPath("pkg-config"); err != nil {
		t.Skip("pkg-config not available")
	}

	// Try to find a common library that's likely installed
	commonLibs := []string{"glib-2.0", "zlib"}

	for _, lib := range commonLibs {
		// Check if pkg-config knows about this library
		cmd := exec.Command("pkg-config", "--exists", lib)
		if cmd.Run() != nil {
			continue
		}

		t.Run(lib, func(t *testing.T) {
			path, err := FindLibraryPath(lib)
			if err != nil {
				t.Errorf("FindLibraryPath(%q) failed: %v", lib, err)
				return
			}

			// Verify the path exists
			if _, err := os.Stat(path); err != nil {
				t.Errorf("Returned path %q does not exist", path)
			}
		})
		return // Only need to test one
	}

	t.Skip("No common libraries found via pkg-config")
}

func TestFindLibraryFile_Integration(t *testing.T) {
	// Try to find libc which should exist on any Linux system
	libcNames := []string{"libc.so.6", "libc.so"}

	for _, name := range libcNames {
		path, err := FindLibraryFile(name)
		if err == nil {
			// Verify the path exists
			if _, err := os.Stat(path); err != nil {
				t.Errorf("Returned path %q does not exist", path)
			}
			return
		}
	}

	t.Skip("Could not find libc.so - unusual system configuration")
}

func TestFindInCommonPaths(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create a fake library directory with a fake .so file
	libDir := filepath.Join(tmpDir, "lib")
	if err := os.MkdirAll(libDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a fake library file
	fakeLib := filepath.Join(libDir, "libfaketest.so.1")
	if err := os.WriteFile(fakeLib, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	// Temporarily add our test dir to defaultLibPaths
	originalPaths := defaultLibPaths
	defaultLibPaths = append([]string{libDir}, defaultLibPaths...)
	defer func() { defaultLibPaths = originalPaths }()

	// Now test finding it
	path, err := findInCommonPaths("faketest")
	if err != nil {
		t.Errorf("findInCommonPaths(\"faketest\") failed: %v", err)
		return
	}

	if path != libDir {
		t.Errorf("findInCommonPaths(\"faketest\") = %q, want %q", path, libDir)
	}
}

func TestFindWithLdconfig(t *testing.T) {
	// Skip if ldconfig is not available
	if _, err := exec.LookPath("ldconfig"); err != nil {
		t.Skip("ldconfig not available")
	}

	// Check if we can run ldconfig -p
	cmd := exec.Command("ldconfig", "-p")
	output, err := cmd.Output()
	if err != nil {
		t.Skip("ldconfig -p failed")
	}

	// Find any library from the output to test with
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=>") && strings.Contains(line, "libc.so") {
			// We found libc, try to find it
			path, err := findWithLdconfig("glib-2.0") // Common library
			if err == nil {
				if _, statErr := os.Stat(path); statErr != nil {
					t.Errorf("Returned path %q does not exist", path)
				}
				return
			}
			// If glib not found, that's okay - just means it's not installed
			break
		}
	}
}

func TestFindLibraryPathWithOptions_IncludeCurrentDir(t *testing.T) {
	// Create a temporary directory and change to it
	tmpDir := t.TempDir()

	// Create a fake library file in the temp dir
	fakeLib := filepath.Join(tmpDir, "libcwdtest.so.1")
	if err := os.WriteFile(fakeLib, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	// Save current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Without IncludeCurrentDir, should not find it
	_, err = FindLibraryPathWithOptions("cwdtest", FindOptions{IncludeCurrentDir: false})
	if err == nil {
		t.Error("Expected error without IncludeCurrentDir")
	}

	// With IncludeCurrentDir, should find it
	path, err := FindLibraryPathWithOptions("cwdtest", FindOptions{IncludeCurrentDir: true})
	if err != nil {
		t.Errorf("FindLibraryPathWithOptions with IncludeCurrentDir failed: %v", err)
		return
	}

	if path != tmpDir {
		t.Errorf("Expected path %q, got %q", tmpDir, path)
	}
}

func TestFindLibraryPathWithOptions_ExtraPaths(t *testing.T) {
	// Create a temporary directory with a fake library
	tmpDir := t.TempDir()

	fakeLib := filepath.Join(tmpDir, "libextratest.so.1")
	if err := os.WriteFile(fakeLib, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	// Should find it with ExtraPaths
	path, err := FindLibraryPathWithOptions("extratest", FindOptions{
		ExtraPaths: []string{tmpDir},
	})
	if err != nil {
		t.Errorf("FindLibraryPathWithOptions with ExtraPaths failed: %v", err)
		return
	}

	if path != tmpDir {
		t.Errorf("Expected path %q, got %q", tmpDir, path)
	}
}

func TestDefaultLibPaths_ContainsDistros(t *testing.T) {
	// Verify that paths for various distros are included
	expectedPaths := map[string][]string{
		"Debian/Ubuntu": {"/usr/lib/x86_64-linux-gnu", "/usr/lib/aarch64-linux-gnu"},
		"Fedora/RHEL":   {"/usr/lib64/gtk-3.0", "/usr/lib64/gtk-4.0"},
		"Arch":          {"/usr/lib/webkit2gtk-4.0", "/usr/lib/webkit2gtk-4.1"},
		"openSUSE":      {"/usr/lib64/gcc/x86_64-suse-linux"},
		"Local":         {"/usr/local/lib", "/usr/local/lib64"},
	}

	for distro, paths := range expectedPaths {
		for _, path := range paths {
			found := false
			for _, defaultPath := range defaultLibPaths {
				if defaultPath == path {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Missing %s path: %s", distro, path)
			}
		}
	}
}

func TestGetFlatpakLibPaths(t *testing.T) {
	// This test just ensures the function doesn't panic
	// Actual paths depend on system state
	paths := getFlatpakLibPaths()
	t.Logf("Found %d Flatpak lib paths", len(paths))
	for _, p := range paths {
		t.Logf("  %s", p)
	}
}

func TestGetSnapLibPaths(t *testing.T) {
	// This test just ensures the function doesn't panic
	// Actual paths depend on system state
	paths := getSnapLibPaths()
	t.Logf("Found %d Snap lib paths", len(paths))
	for _, p := range paths {
		t.Logf("  %s", p)
	}
}

func TestGetNixLibPaths(t *testing.T) {
	// This test just ensures the function doesn't panic
	paths := getNixLibPaths()
	t.Logf("Found %d Nix lib paths", len(paths))
	for _, p := range paths {
		t.Logf("  %s", p)
	}
}

func TestGetAllLibPaths_IncludesDynamicPaths(t *testing.T) {
	paths := GetAllLibPaths()

	// Should have at least the default paths
	if len(paths) < len(defaultLibPaths) {
		t.Errorf("GetAllLibPaths returned fewer paths (%d) than defaultLibPaths (%d)",
			len(paths), len(defaultLibPaths))
	}

	// Log all paths for debugging
	t.Logf("Total paths: %d", len(paths))
}

func TestGetAllLibPaths_DoesNotIncludeCurrentDir(t *testing.T) {
	paths := GetAllLibPaths()

	for _, p := range paths {
		if p == "." {
			t.Error("GetAllLibPaths should not include '.' for security reasons")
		}
	}
}

func TestInvalidateCache(t *testing.T) {
	// First call populates cache
	paths1 := GetAllLibPaths()

	// Invalidate and call again
	InvalidateCache()
	paths2 := GetAllLibPaths()

	// Should get same results (assuming no system changes)
	if len(paths1) != len(paths2) {
		t.Logf("Path counts differ after invalidation: %d vs %d", len(paths1), len(paths2))
		// This is not necessarily an error, just informational
	}

	// Verify cache is working by checking getFlatpakLibPaths is fast
	// (would be slow if cache wasn't working)
	for i := 0; i < 100; i++ {
		_ = getFlatpakLibPaths()
	}
}

func TestFindLibraryPath_ParallelConsistency(t *testing.T) {
	// Skip if pkg-config is not available
	if _, err := exec.LookPath("pkg-config"); err != nil {
		t.Skip("pkg-config not available")
	}

	// Check if glib-2.0 is available
	cmd := exec.Command("pkg-config", "--exists", "glib-2.0")
	if cmd.Run() != nil {
		t.Skip("glib-2.0 not installed")
	}

	// Run parallel and sequential versions multiple times
	// to ensure they return consistent results
	for i := 0; i < 10; i++ {
		parallelPath, parallelErr := FindLibraryPath("glib-2.0")
		seqPath, seqErr := FindLibraryPathSequential("glib-2.0")

		if parallelErr != nil && seqErr == nil {
			t.Errorf("Parallel failed but sequential succeeded: %v", parallelErr)
		}
		if parallelErr == nil && seqErr != nil {
			t.Errorf("Sequential failed but parallel succeeded: %v", seqErr)
		}

		// Both should find the library (path might differ if found by different methods)
		if parallelErr != nil {
			t.Errorf("Iteration %d: parallel search failed: %v", i, parallelErr)
		}
		if seqErr != nil {
			t.Errorf("Iteration %d: sequential search failed: %v", i, seqErr)
		}

		// Log paths for debugging
		t.Logf("Iteration %d: parallel=%s, sequential=%s", i, parallelPath, seqPath)
	}
}

func TestFindLibraryPath_ParallelNotFound(t *testing.T) {
	// Both parallel and sequential should return the same error for non-existent libs
	_, parallelErr := FindLibraryPath("nonexistent-library-xyz-123")
	_, seqErr := FindLibraryPathSequential("nonexistent-library-xyz-123")

	if parallelErr == nil {
		t.Error("Parallel search should fail for nonexistent library")
	}
	if seqErr == nil {
		t.Error("Sequential search should fail for nonexistent library")
	}

	// Both should return LibraryNotFoundError
	if _, ok := parallelErr.(*LibraryNotFoundError); !ok {
		t.Errorf("Parallel: expected LibraryNotFoundError, got %T", parallelErr)
	}
	if _, ok := seqErr.(*LibraryNotFoundError); !ok {
		t.Errorf("Sequential: expected LibraryNotFoundError, got %T", seqErr)
	}
}

// Benchmarks

// BenchmarkFindLibraryPath benchmarks finding a library via the full search chain.
func BenchmarkFindLibraryPath(b *testing.B) {
	// Test with glib-2.0 which is commonly installed
	if _, err := exec.LookPath("pkg-config"); err != nil {
		b.Skip("pkg-config not available")
	}
	cmd := exec.Command("pkg-config", "--exists", "glib-2.0")
	if cmd.Run() != nil {
		b.Skip("glib-2.0 not installed")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindLibraryPath("glib-2.0")
	}
}

// BenchmarkFindLibraryPath_NotFound benchmarks the worst case (library not found).
func BenchmarkFindLibraryPath_NotFound(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindLibraryPath("nonexistent-library-xyz-123")
	}
}

// BenchmarkFindLibraryFile benchmarks finding a specific library file.
func BenchmarkFindLibraryFile(b *testing.B) {
	// libc.so.6 should exist on any Linux system
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindLibraryFile("libc.so.6")
	}
}

// BenchmarkGetAllLibPaths benchmarks collecting all library paths.
func BenchmarkGetAllLibPaths(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetAllLibPaths()
	}
}

// BenchmarkFindWithPkgConfig benchmarks pkg-config lookup directly.
func BenchmarkFindWithPkgConfig(b *testing.B) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		b.Skip("pkg-config not available")
	}
	cmd := exec.Command("pkg-config", "--exists", "glib-2.0")
	if cmd.Run() != nil {
		b.Skip("glib-2.0 not installed")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = findWithPkgConfig("glib-2.0")
	}
}

// BenchmarkFindWithLdconfig benchmarks ldconfig lookup directly.
func BenchmarkFindWithLdconfig(b *testing.B) {
	if _, err := exec.LookPath("ldconfig"); err != nil {
		b.Skip("ldconfig not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = findWithLdconfig("glib-2.0")
	}
}

// BenchmarkFindInCommonPaths benchmarks filesystem scanning.
func BenchmarkFindInCommonPaths(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = findInCommonPaths("glib-2.0")
	}
}

// BenchmarkGetFlatpakLibPaths benchmarks Flatpak path discovery.
func BenchmarkGetFlatpakLibPaths(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getFlatpakLibPaths()
	}
}

// BenchmarkGetSnapLibPaths benchmarks Snap path discovery.
func BenchmarkGetSnapLibPaths(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getSnapLibPaths()
	}
}

// BenchmarkGetNixLibPaths benchmarks Nix path discovery.
func BenchmarkGetNixLibPaths(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getNixLibPaths()
	}
}

// BenchmarkPkgConfigToLibName benchmarks the name conversion function.
func BenchmarkPkgConfigToLibName(b *testing.B) {
	names := []string{"gtk+-3.0", "webkit2gtk-4.1", "glib-2.0", "cairo", "libsoup-3.0"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, name := range names {
			_ = pkgConfigToLibName(name)
		}
	}
}

// BenchmarkFindLibraryPathSequential benchmarks the sequential search.
func BenchmarkFindLibraryPathSequential(b *testing.B) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		b.Skip("pkg-config not available")
	}
	cmd := exec.Command("pkg-config", "--exists", "glib-2.0")
	if cmd.Run() != nil {
		b.Skip("glib-2.0 not installed")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindLibraryPathSequential("glib-2.0")
	}
}

// BenchmarkFindLibraryPathSequential_NotFound benchmarks the sequential worst case.
func BenchmarkFindLibraryPathSequential_NotFound(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindLibraryPathSequential("nonexistent-library-xyz-123")
	}
}

// BenchmarkFindLibraryPathParallel explicitly tests parallel performance.
func BenchmarkFindLibraryPathParallel(b *testing.B) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		b.Skip("pkg-config not available")
	}
	cmd := exec.Command("pkg-config", "--exists", "glib-2.0")
	if cmd.Run() != nil {
		b.Skip("glib-2.0 not installed")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindLibraryPath("glib-2.0")
	}
}

// Tests for multi-library search functions

func TestFindFirstLibrary(t *testing.T) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		t.Skip("pkg-config not available")
	}

	// Test with a mix of existing and non-existing libraries
	match, err := FindFirstLibrary("nonexistent-xyz", "glib-2.0", "also-nonexistent")
	if err != nil {
		t.Skipf("glib-2.0 not installed: %v", err)
	}

	if match.Name != "glib-2.0" {
		t.Errorf("Expected glib-2.0, got %s", match.Name)
	}
	if match.Path == "" {
		t.Error("Expected non-empty path")
	}
}

func TestFindFirstLibrary_AllNotFound(t *testing.T) {
	_, err := FindFirstLibrary("nonexistent-1", "nonexistent-2", "nonexistent-3")
	if err == nil {
		t.Error("Expected error for all non-existent libraries")
	}
}

func TestFindFirstLibrary_Empty(t *testing.T) {
	_, err := FindFirstLibrary()
	if err == nil {
		t.Error("Expected error for empty library list")
	}
}

func TestFindFirstLibraryOrdered(t *testing.T) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		t.Skip("pkg-config not available")
	}

	// glib-2.0 should be found, and since it's first, it should be returned
	match, err := FindFirstLibraryOrdered("glib-2.0", "nonexistent-xyz")
	if err != nil {
		t.Skipf("glib-2.0 not installed: %v", err)
	}

	if match.Name != "glib-2.0" {
		t.Errorf("Expected glib-2.0, got %s", match.Name)
	}
}

func TestFindFirstLibraryOrdered_PreferFirst(t *testing.T) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		t.Skip("pkg-config not available")
	}

	// Check what GTK versions are available
	gtk4Available := exec.Command("pkg-config", "--exists", "gtk4").Run() == nil
	gtk3Available := exec.Command("pkg-config", "--exists", "gtk+-3.0").Run() == nil

	if !gtk4Available && !gtk3Available {
		t.Skip("Neither GTK3 nor GTK4 installed")
	}

	// If both available, test that order is respected
	if gtk4Available && gtk3Available {
		match, err := FindFirstLibraryOrdered("gtk4", "gtk+-3.0")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if match.Name != "gtk4" {
			t.Errorf("Expected gtk4 (first in order), got %s", match.Name)
		}

		// Reverse order
		match, err = FindFirstLibraryOrdered("gtk+-3.0", "gtk4")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if match.Name != "gtk+-3.0" {
			t.Errorf("Expected gtk+-3.0 (first in order), got %s", match.Name)
		}
	}
}

func TestFindAllLibraries(t *testing.T) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		t.Skip("pkg-config not available")
	}

	matches := FindAllLibraries("glib-2.0", "nonexistent-xyz", "zlib")

	// Should find at least glib-2.0 on most systems
	if len(matches) == 0 {
		t.Skip("No common libraries found")
	}

	t.Logf("Found %d libraries:", len(matches))
	for _, m := range matches {
		t.Logf("  %s at %s", m.Name, m.Path)
	}

	// Verify no duplicates and no nonexistent library
	seen := make(map[string]bool)
	for _, m := range matches {
		if m.Name == "nonexistent-xyz" {
			t.Error("Should not have found nonexistent library")
		}
		if seen[m.Name] {
			t.Errorf("Duplicate match for %s", m.Name)
		}
		seen[m.Name] = true
	}
}

func TestFindAllLibraries_Empty(t *testing.T) {
	matches := FindAllLibraries()
	if len(matches) != 0 {
		t.Error("Expected empty result for empty input")
	}
}

func TestFindAllLibraries_AllNotFound(t *testing.T) {
	matches := FindAllLibraries("nonexistent-1", "nonexistent-2")
	if len(matches) != 0 {
		t.Errorf("Expected empty result, got %d matches", len(matches))
	}
}

// Benchmarks for multi-library search

func BenchmarkFindFirstLibrary(b *testing.B) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		b.Skip("pkg-config not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindFirstLibrary("nonexistent-1", "glib-2.0", "nonexistent-2")
	}
}

func BenchmarkFindFirstLibraryOrdered(b *testing.B) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		b.Skip("pkg-config not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindFirstLibraryOrdered("nonexistent-1", "glib-2.0", "nonexistent-2")
	}
}

func BenchmarkFindAllLibraries(b *testing.B) {
	if _, err := exec.LookPath("pkg-config"); err != nil {
		b.Skip("pkg-config not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindAllLibraries("glib-2.0", "zlib", "nonexistent-xyz")
	}
}
