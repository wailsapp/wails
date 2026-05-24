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

// LibraryMatch holds information about a found library.
type LibraryMatch struct {
	Name string
	Path string
}

// FindFirstLibrary is a stub for non-Linux platforms.
func FindFirstLibrary(libNames ...string) (*LibraryMatch, error) {
	if len(libNames) == 0 {
		return nil, &LibraryNotFoundError{Name: "no libraries specified"}
	}
	return nil, &LibraryNotFoundError{Name: libNames[0]}
}

// FindFirstLibraryOrdered is a stub for non-Linux platforms.
func FindFirstLibraryOrdered(libNames ...string) (*LibraryMatch, error) {
	if len(libNames) == 0 {
		return nil, &LibraryNotFoundError{Name: "no libraries specified"}
	}
	return nil, &LibraryNotFoundError{Name: libNames[0]}
}

// FindAllLibraries is a stub for non-Linux platforms.
func FindAllLibraries(libNames ...string) []LibraryMatch {
	return nil
}

// FindLibraryPathSequential is a stub for non-Linux platforms.
func FindLibraryPathSequential(libName string) (string, error) {
	return "", &LibraryNotFoundError{Name: libName}
}
