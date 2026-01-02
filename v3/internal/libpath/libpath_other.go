//go:build !linux

package libpath

// FindLibraryPath is a stub for non-Linux platforms.
func FindLibraryPath(libName string) (string, error) {
	return "", &LibraryNotFoundError{Name: libName}
}

// FindLibraryFile is a stub for non-Linux platforms.
func FindLibraryFile(fileName string) (string, error) {
	return "", &LibraryNotFoundError{Name: fileName}
}

// GetAllLibPaths returns an empty slice on non-Linux platforms.
func GetAllLibPaths() []string {
	return nil
}

// InvalidateCache is a no-op on non-Linux platforms.
func InvalidateCache() {}

// FindOptions controls library search behavior.
type FindOptions struct {
	IncludeCurrentDir bool
	ExtraPaths        []string
}

// FindLibraryPathWithOptions is a stub for non-Linux platforms.
func FindLibraryPathWithOptions(libName string, opts FindOptions) (string, error) {
	return "", &LibraryNotFoundError{Name: libName}
}

// LibraryNotFoundError is returned when a library cannot be found.
type LibraryNotFoundError struct {
	Name string
}

func (e *LibraryNotFoundError) Error() string {
	return "library not found: " + e.Name
}
