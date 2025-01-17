package collect

import (
	"go/types"
	"path/filepath"
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

	// Fetch and update counter for name.
	counter := imports.counters[pkg.Name]
	imports.counters[pkg.Name] = counter + 1

	// Always add counters to
	imports.External[pkg.Path] = ImportInfo{
		Name:    pkg.Name,
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
// AddType does not support unsynchronised concurrent calls
// on the same receiver.
func (imports *ImportMap) AddType(typ types.Type) {
	imports.addTypeImpl(typ, make(map[*types.TypeName]bool))
}

// addTypeImpl provides the actual implementation of AddType.
// The visited parameter is used to break cycles.
func (imports *ImportMap) addTypeImpl(typ types.Type, visited map[*types.TypeName]bool) {
	collector := imports.collector
	if collector == nil {
		panic("AddType called on ImportMap with nil importing package")
	}

	for { // Avoid recursion where possible.
		switch t := typ.(type) {
		case *types.Alias, *types.Named:
			obj := typ.(interface{ Obj() *types.TypeName }).Obj()
			if visited[obj] {
				return
			}
			visited[obj] = true

			if obj.Pkg() == nil {
				// Ignore universe type.
				return
			}

			if obj.Pkg().Path() == imports.Self {
				imports.ImportModels = true
			}

			// Record model.
			imports.collector.Model(obj)

			// Import parent package.
			imports.Add(collector.Package(obj.Pkg()))

			instance, _ := t.(interface{ TypeArgs() *types.TypeList })
			if instance != nil {
				// Record type argument dependencies.
				if targs := instance.TypeArgs(); targs != nil {
					for i := range targs.Len() {
						imports.addTypeImpl(targs.At(i), visited)
					}
				}
			}

			if collector.options.UseInterfaces {
				// No creation/initialisation code required.
				return
			}

			if _, isAlias := t.(*types.Alias); isAlias {
				// Aliased type might be needed during
				// JS value creation and initialisation.
				typ = types.Unalias(typ)
				break
			}

			if IsClass(typ) || IsString(typ) || IsAny(typ) {
				return
			}

			// If named type does not map to a class, string or unknown type,
			// its underlying type may be needed during JS value creation.
			typ = typ.Underlying()

		case *types.Basic:
			if t.Info()&types.IsComplex != 0 {
				// Complex types are not supported by encoding/json
				collector.logger.Warningf("complex types are not supported by encoding/json")
			}
			return

		case *types.Array, *types.Pointer, *types.Slice:
			typ = typ.(interface{ Elem() types.Type }).Elem()

		case *types.Chan:
			collector.logger.Warningf("channel types are not supported by encoding/json")
			return

		case *types.Map:
			if IsMapKey(t.Key()) {
				if IsString(t.Key()) {
					// This model type is always rendered as a string alias,
					// hence we can generate it and use it as a type for JS object keys.
					imports.addTypeImpl(t.Key(), visited)
				}
			} else {
				collector.logger.Warningf(
					"%s is used as a map key, but does not implement encoding.TextMarshaler: this will likely result in runtime errors",
					types.TypeString(t.Key(), nil),
				)
			}

			typ = t.Elem()

		case *types.Signature:
			collector.logger.Warningf("function types are not supported by encoding/json")
			return

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
			for i := range len(info.Fields) - 1 {
				imports.addTypeImpl(info.Fields[i].Type, visited)
			}

			// Process last field without recursion.
			typ = info.Fields[len(info.Fields)-1].Type

		case *types.Interface, *types.TypeParam:
			// No dependencies.
			return

		default:
			collector.logger.Warningf("unknown type %s: please report this to Wails maintainers", typ)
			return
		}
	}
}

// computeImportPath returns the shortest relative import path
// through which the importer package can reference the imported one.
func computeImportPath(importer string, imported string) string {
	rel, err := filepath.Rel(importer, imported)
	if err != nil {
		panic(err)
	}

	rel = filepath.ToSlash(rel)
	if rel[0] == '.' {
		return rel
	} else {
		return "./" + rel
	}
}
