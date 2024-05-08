package collect

import (
	"fmt"
	"path"
	"strings"
)

type (
	// ImportMap records deduplicated imports by a binding or models module.
	// It computes relative import paths and assigns import names
	// taking care to avoid collisions.
	ImportMap struct {
		// Self records the path of the importing package.
		Self string

		// Map records information about each imported package,
		// keyed by package path.
		Map map[string]ImportInfo

		counters map[string]int
	}

	// ImportInfo records information about a single import.
	ImportInfo struct {
		Name    string
		RelPath string
	}
)

// NewImportMap initialises an import map for the given importer package.
func NewImportMap(importer *PackageInfo) *ImportMap {
	return &ImportMap{
		Self:     importer.Path,
		Map:      make(map[string]ImportInfo),
		counters: make(map[string]int),
	}
}

// Add adds the given package to the import map if not already present,
// choosing import names so as to avoid collisions.
func (imports *ImportMap) Add(pkg *PackageInfo) {
	if imports.Map[pkg.Path].Name != "" {
		// Package already imported.
		return
	}

	name := path.Base(pkg.Path)
	if pkg.Collect() {
		name = pkg.Name
	}

	counter := imports.counters[name]
	imports.counters[name] = counter + 1

	if counter > 0 {
		name = fmt.Sprintf("%s$%d", name, counter)
	}

	imports.Map[pkg.Path] = ImportInfo{
		Name:    name,
		RelPath: computeImportPath(imports.Self, pkg.Path),
	}
}

// computeImportPath returns the shortest relative import path
// through which the importer package can reference the imported one.
//
// We provide a custom implementation to work around
// the fact that filepath.Rel may change separators
// plus we know that loaded package paths are well-formed.
func computeImportPath(importer string, imported string) string {
	// Find longest common prefix.
	i, slash := 0, -1
	for ; i < len(importer) && i < len(imported); i++ {
		if importer[i] != imported[i] {
			break
		}

		if importer[i] == '/' {
			slash = i
		}
	}

	// One path is a prefix of the other, seen as strings:
	// check if the extension starts with a slash.
	if (i == len(importer) && i < len(imported) && imported[i] == '/') || (i == len(imported) && i < len(importer) && importer[i] == '/') {
		slash = i
	}

	// Build path from the right number of parent steps plus suffix.
	var builder strings.Builder

	back := strings.Count(importer[slash:], "/")
	if back == 0 {
		builder.WriteByte('.')
	} else {
		builder.WriteString("..")
		for back--; back > 0; back-- {
			builder.WriteString("/..")
		}
	}

	builder.WriteString(imported[slash:])

	return builder.String()
}
