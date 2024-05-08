package collect

import (
	"go/types"
	"strings"
)

type (
	// ImportMap records deduplicated imports by a binding or models module.
	// It computes relative import paths and assigns import names,
	// taking care to avoid collisions.
	ImportMap struct {
		// Self records the path of the importing package.
		Self string

		// ImportModels records whether models from the current package may be needed.
		ImportModels bool
		// ImportInternal records whether internal models from the current package may be needed.
		ImportInternal bool

		// External records information about each imported package,
		// keyed by package path.
		External map[string]ImportInfo

		// counters holds the occurence count for each package name in External.
		counters  map[string]int
		collector *Collector
	}

	// ImportInfo records information about a single import.
	ImportInfo struct {
		Name    string
		Index   int // Progressive number for identically named imports, starting from 0 for each distinct name.
		RelPath string
	}
)

// NewImportMap initialises an import map for the given importer package.
// The argument may be nil, in which case import paths will be relative
// to the root output directory.
func NewImportMap(importer *PackageInfo) *ImportMap {
	var (
		self      string
		collector *Collector
	)
	if importer != nil {
		self = importer.Path
		collector = importer.collector
	}

	return &ImportMap{
		Self: self,

		External: make(map[string]ImportInfo),

		counters:  make(map[string]int),
		collector: collector,
	}
}

// Merge merges the given import map into the receiver.
// The importing package must be the same.
func (imports *ImportMap) Merge(other *ImportMap) {
	if other.Self != imports.Self {
		panic("cannot merge import maps with different importing package")
	}

	if other.ImportModels {
		imports.ImportModels = true
	}
	if other.ImportInternal {
		imports.ImportInternal = true
	}

	for path, info := range other.External {
		if _, ok := imports.External[path]; ok {
			continue
		}

		counter := imports.counters[info.Name]
		imports.counters[info.Name] = counter + 1

		imports.External[path] = ImportInfo{
			Name:    info.Name,
			Index:   counter,
			RelPath: info.RelPath,
		}
	}
}

// Add adds the given package to the import map if not already present,
// choosing import names so as to avoid collisions.
//
// Add does not support unsynchronised concurrent calls
// on the same receiver.
func (imports *ImportMap) Add(pkg *PackageInfo) {
	if pkg.Path == imports.Self {
		// Do not import self.
		return
	}

	if imports.External[pkg.Path].Name != "" {
		// Package already imported.
		return
	}

	name := pkg.Collect().Name

	// Fetch and update counter for name.
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
// It is a runtime error to call AddType on an ImportMap
// created with nil importing package.
//
// Add does not support unsynchronised concurrent calls
// on the same receiver.
func (imports *ImportMap) AddType(typ types.Type) {
	collector := imports.collector
	if collector == nil {
		panic("AddType called on ImportMap with nil importing package")
	}

	for { // Avoid recursion where possible.
		switch t := typ.(type) {
		case *types.Basic:
			if t.Info()&types.IsComplex != 0 {
				// Complex types are not supported by encoding/json
				collector.logger.Warningf("complex types are not supported by encoding/json")
			}
			return

		case *types.Alias:
			if t.Obj().Pkg() == nil {
				// Ignore universe type.
				return
			}

			// Record used types from self.
			if t.Obj().Pkg().Path() == imports.Self {
				if t.Obj().Exported() {
					imports.ImportModels = true
				} else {
					imports.ImportInternal = true
				}
			}

			collector.Model(t.Obj())
			imports.Add(collector.Package(t.Obj().Pkg()))

			// The aliased type might be needed during
			// JS value creation and initialisation.
			typ = types.Unalias(typ)

		case *types.Array:
			typ = t.Elem()

		case *types.Chan:
			collector.logger.Warningf("channel types are not supported by encoding/json")
			return

		case *types.Map:
			if IsMapKey(t.Key()) {
				if IsString(t.Key()) {
					// This model type is always rendered as a string alias,
					// hence we can generate it and use it as a type for JS object keys.
					imports.AddType(t.Key())
				}
			} else {
				collector.logger.Warningf(
					"%s is used as a map key, but does not implement encoding.TextMarshaler: this will likely result in runtime errors",
					types.TypeString(t.Key(), nil),
				)
			}

			typ = t.Elem()

		case *types.Named:
			if t.Obj().Pkg() == nil {
				// Ignore universe type.
				return
			}

			// Record used types from self.
			if t.Obj().Pkg().Path() == imports.Self {
				if t.Obj().Exported() {
					imports.ImportModels = true
				} else {
					imports.ImportInternal = true
				}
			}

			imports.collector.Model(t.Obj())
			imports.Add(collector.Package(t.Obj().Pkg()))

			if IsClass(typ) || IsString(typ) || IsAny(typ) {
				return
			}

			// If named type does not map to a class, string or unknown type,
			// its underlying type may be needed during JS value creation.
			typ = t.Underlying()

		case *types.Pointer:
			typ = t.Elem()

		case *types.Signature:
			collector.logger.Warningf("function types are not supported by encoding/json")
			return

		case *types.Slice:
			typ = t.Elem()

		case *types.Struct:
			if t.NumFields() == 0 {
				// Empty struct.
				return
			}

			// Retrieve struct info and ensure it is complete.
			info := collector.Struct(t).Collect()

			if len(info.Fields) == 0 {
				// No visible fields.
				return
			}

			// Add field dependencies.
			for i, length := 0, len(info.Fields)-1; i < length; i++ {
				imports.AddType(info.Fields[i].Type)
			}

			// Process last field without recursion.
			typ = info.Fields[len(info.Fields)-1].Type

		case *types.Interface, *types.TypeParam:
			// Rendered as any.
			return

		default:
			collector.logger.Warningf("unknown type %s: please report this to Wails maintainers", typ)
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

	if slash == -1 {
		return "./" + imported
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
