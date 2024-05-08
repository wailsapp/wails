package collect

import (
	"go/types"
	"path"
	"strings"

	"github.com/pterm/pterm"
)

type (
	// ImportMap records deduplicated imports by a binding or models module.
	// It computes relative import paths and assigns import names,
	// taking care to avoid collisions.
	ImportMap struct {
		// Self records the path of the importing package.
		Self string

		// Models records required exported models from self.
		// Values are true for typedefs.
		Models map[string]bool

		// Internal records required unexported models from Self.
		// Values are true for typedefs.
		Internal map[string]bool

		// External records information about each imported package,
		// keyed by package path.
		External map[string]ImportInfo

		counters map[string]int
	}

	// ImportInfo records information about a single import.
	ImportInfo struct {
		Name    string
		Index   int // Identically named imports always have distinct indexes.
		RelPath string
	}
)

// NewImportMap initialises an import map for the given importer package.
func NewImportMap(importer *PackageInfo) *ImportMap {
	return &ImportMap{
		Self:     importer.Path,
		External: make(map[string]ImportInfo),
		counters: make(map[string]int),
	}
}

// Add adds the given package to the import map if not already present,
// choosing import names so as to avoid collisions.
//
// Add DOES NOT support unsynchronised concurrent calls.
func (imports *ImportMap) Add(pkg *PackageInfo) {
	if pkg.Path == imports.Self {
		// Do not import self.
		return
	}

	if imports.External[pkg.Path].Name != "" {
		// Package already imported.
		return
	}

	name := path.Base(pkg.Path)
	if pkg.Collect() {
		name = pkg.Name
	}

	counter := imports.counters[name]
	imports.counters[name] = counter + 1

	// Always add counters to
	imports.External[pkg.Path] = ImportInfo{
		Name:    name,
		Index:   counter,
		RelPath: computeImportPath(imports.Self, pkg.Path),
	}
}

// AddType adds all dependencies of the given type to the import map
// and marks all referenced named types as models.
//
// Add DOES NOT support unsynchronised concurrent calls.
func (imports *ImportMap) AddType(collector *Collector, typ types.Type) {
	for {
		switch t := typ.(type) {
		case *types.Basic:
			if t.Info()&types.IsComplex != 0 {
				// Complex types are not supported by encoding/json
				collector.complexWarning()
			}
			return

		case *types.Alias:
			if t.Obj().Pkg() == nil {
				// Universe type
				return
			}

			// Record used types from self.
			if t.Obj().Pkg().Path() == imports.Self {
				if t.Obj().Exported() {
					imports.Models[t.Obj().Name()] = true
				} else {
					imports.Internal[t.Obj().Name()] = true
				}
			}

			pkg := collector.Package(t.Obj().Pkg().Path())
			pkg.AddModels(t.Obj())
			imports.Add(pkg)
			return

		case *types.Array:
			typ = t.Elem()

		case *types.Chan:
			collector.chanWarning()
			return

		case *types.Map:
			if IsMapKey(t.Key()) {
				if IsAlwaysTextMarshaler(t.Key()) && !MaybeJSONMarshaler(t.Key()) {
					// This type is always rendered as a string,
					// hence we can use it safely as an object key type.
					imports.AddType(collector, t.Key())
				}
			} else {
				pterm.Warning.Printfln(
					"%s is used as a map key, but does not implement encoding.TextMarshaler: this will likely result in runtime errors",
					types.TypeString(t.Key(), nil),
				)
			}

			typ = t.Elem()

		case *types.Named:
			if t.Obj().Pkg() == nil {
				// Universe type
				return
			}

			if t.TypeParams() != nil {
				collector.genericWarning()
				return
			}

			// Record used types from self.
			if t.Obj().Pkg().Path() == imports.Self {
				isTypedef := IsAlwaysTextMarshaler(t) || MaybeJSONMarshaler(t)
				if t.Obj().Exported() {
					imports.Models[t.Obj().Name()] = isTypedef
				} else {
					imports.Internal[t.Obj().Name()] = isTypedef
				}
			}

			pkg := collector.Package(t.Obj().Pkg().Path())
			pkg.AddModels(t.Obj())
			imports.Add(pkg)
			return

		case *types.Pointer:
			typ = t.Elem()

		case *types.Signature:
			collector.funcWarning()
			return

		case *types.Slice:
			typ = t.Elem()

		case *types.Struct:
			if t.NumFields() == 0 {
				// Empty struct.
				return
			}

			// Retrieve struct info and ensure it is initialised.
			info := collector.Struct(t)
			info.Collect()

			if len(info.Fields) == 0 {
				// No visible fields.
				return
			}

			// Add field dependencies.
			for i := range len(info.Fields) - 1 {
				imports.AddType(collector, info.Fields[i].Type)
			}

			// Process last field without recursion.
			typ = info.Fields[len(info.Fields)-1].Type

		default:
			// Atomic type
			return
		}
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
