package collect

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

type (
	// PackageInfo records information about a package;
	// in particular, it holds precomputed maps that speed up access
	// to the syntax constructs that declare types, fields,
	// constants and methods.
	//
	// A dummy PackageInfo can be initialised with just the path;
	// all other fields will be populated upon calling [PackageInfo.Collect]
	// for the first time.
	//
	// Read accesses to the Path field are safe at any time
	// without any synchronisation.
	// Read accesses to all other exported fields are only safe
	// if a call to [PackageInfo.Collect] has been completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	//
	// Concurrent write accesses are only allowed through the provided methods.
	PackageInfo struct {
		Path string
		Name string

		Docs []*ast.CommentGroup

		Imports map[*ast.File]*FileImports
		Types   map[string]*TypeDefInfo
		Consts  map[string]*ConstInfo

		// The next two fields record generated bindings for this package.
		// We use a slice behind a lock because it is much faster
		// than [sync.Map] in write-heavy scenarios like this one.
		mu       sync.Mutex
		bindings []*BoundTypeInfo

		// models records the models that have to be generated for this package.
		// We use [sync.Map] for atomic swapping.
		models sync.Map

		// stats caches statistics about this package.
		// Mutex mu must be locked before accessing this field.
		stats *Stats

		// collector holds a pointer to the parent [Collector].
		collector *Collector

		// source holds either a pointer to [packages.Package],
		// or a [Controller] instance that may be used to load the package.
		source any
		once   sync.Once
	}

	// FileImports records information
	// about import declarations in an [ast.File].
	FileImports struct {
		Unnamed []string
		Dot     []string
		Named   map[string]string
	}

	// GroupInfo records information about a group
	// of type, field or constant declarations.
	// This may be either a list of distinct specifications
	// wrapped in parentheses, or a single specification
	// declaring multiple fields or constants.
	GroupInfo struct {
		Pos   token.Pos
		Doc   *ast.CommentGroup
		Group *GroupInfo
	}

	// ConstInfo records information about a constant declaration.
	ConstInfo struct {
		Name  string
		Pos   token.Pos
		Doc   *ast.CommentGroup
		Group *GroupInfo
	}

	// PackageIndex lists all bindings, models and unexported models
	// generated from a package.
	//
	// When obtained through a call to [PackageInfo.Index],
	// each binding and model name appears at most once.
	PackageIndex struct {
		Info *PackageInfo

		Bindings []*BoundTypeInfo
		Models   []*ModelInfo
		Internal []*ModelInfo
	}
)

// Preload adds the given package descriptors to the collector,
// so that the loading step may be skipped when collecting information
// about each of those packages.
//
// Preload is safe for concurrent use.
func (collector *Collector) Preload(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		collector.pkgs.LoadOrStore(pkg.PkgPath, &PackageInfo{
			Path:      pkg.PkgPath,
			collector: collector,
			source:    pkg,
		})
	}
}

// Package retrieves the the unique [PackageInfo] instance
// associated to the given path within a Collector.
// If none is present, a new one is initialised.
//
// Package is safe for concurrent use.
func (collector *Collector) Package(path string) *PackageInfo {
	info, _ := collector.pkgs.LoadOrStore(path, &PackageInfo{
		Path:      path,
		collector: collector,
		source:    collector.controller,
	})
	return info.(*PackageInfo)
}

// All calls yield sequentially for each [PackageInfo] instance
// present in the collector. If yield returns false, All stops the iteration.
//
// All may be O(N) with the number of packages in the collector
// even if yield returns false after a constant number of calls.
//
// Package is safe for concurrent use.
func (collector *Collector) Iterate(yield func(pkg *PackageInfo) bool) {
	collector.pkgs.Range(func(key, value any) bool {
		return yield(value.(*PackageInfo))
	})
}

// IsEmpty retuns true if no bindings and models
// were generated for this package.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
func (info *PackageInfo) IsEmpty() bool {
	info.mu.Lock()
	result := len(info.bindings) == 0
	info.mu.Unlock()

	if !result {
		return false
	}

	// No other way to get the length of a sync.Map...
	info.models.Range(func(key, value any) bool {
		result = false
		return false
	})

	return result
}

// Index computes a [PackageIndex] from the list
// of generated bindings and models
// and regenerates cached stats.
//
// Bindings and models appear at most once
// in the returned structure and are sorted by name.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
func (info *PackageInfo) Index() (index PackageIndex) {
	// Init stats
	stats := &Stats{
		NumPackages: 1,
	}

	// Acquire bindings slice.
	info.mu.Lock()

	// Sort bindings by name, then deduplicate.
	// If [Generator.Generate] is called multiple times,
	// there might be distinct objects with the same name,
	// hence we can't just compare pointers.
	slices.SortFunc(info.bindings, func(b1 *BoundTypeInfo, b2 *BoundTypeInfo) int {
		return strings.Compare(b1.Name, b2.Name)
	})
	info.bindings = slices.CompactFunc(info.bindings, func(b1 *BoundTypeInfo, b2 *BoundTypeInfo) bool {
		// If the pointers are equal, so must be the names.
		return b1 == b2 || b1.Name == b2.Name
	})

	// Clone bindings into result.
	index.Bindings = slices.Clone(info.bindings)

	// Release bindings slice.
	info.mu.Unlock()

	// Update bound type stats.
	stats.NumTypes = len(index.Bindings)
	for _, b := range index.Bindings {
		stats.NumMethods += len(b.Methods)
	}

	// Gather models.
	info.models.Range(func(key, value any) bool {
		model := value.(*ModelInfo)
		if len(model.Values) > 0 {
			stats.NumEnums++
		} else {
			stats.NumModels++
		}
		index.Models = append(index.Models, model)
		return true
	})

	// Sort models by exported property (exported first), then by name.
	slices.SortFunc(index.Models, func(m1 *ModelInfo, m2 *ModelInfo) int {
		m1e, m2e := ast.IsExported(m1.Name), ast.IsExported(m2.Name)
		if m1e != m2e {
			if m1e {
				return -1
			} else {
				return 1
			}
		}

		return strings.Compare(m1.Name, m2.Name)
	})

	// Find first unexported model.
	split, _ := slices.BinarySearchFunc(index.Models, struct{}{}, func(m *ModelInfo, _ struct{}) int {
		if ast.IsExported(m.Name) {
			return -1
		} else {
			return 1
		}
	})

	// Separate unexported and exported models.
	index.Internal = index.Models[split:]
	index.Models = index.Models[:split]

	// Store package info.
	index.Info = info

	// Cache stats.
	info.mu.Lock()
	info.stats = stats
	info.mu.Unlock()

	return
}

// Stats returns statistics for this package.
// If no stats have been cached yet, they are generated anew.
// If they are out of date, call [PackageInfo.Index] to regenerate them.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
func (info *PackageInfo) Stats() *Stats {
	var result = &Stats{}
	info.mu.Lock()
	if info.stats == nil {
		info.mu.Unlock()
		info.Index()
		info.mu.Lock()
	}
	*result = *info.stats
	info.mu.Unlock()
	return result
}

// Collect gathers information for the package described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns true on success, false if the package failed to load.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *PackageInfo) Collect() bool {
	info.once.Do(func() {
		pkg, ok := info.source.(*packages.Package)
		if !ok {
			pkg = info.source.(Controller).Load(info.Path)
			if pkg == nil {
				return
			}
		}

		// Discard package source.
		info.source = nil

		// Record package name.
		info.Name = pkg.Name

		// Initialise maps
		info.Imports = make(map[*ast.File]*FileImports)
		info.Types = make(map[string]*TypeDefInfo)
		info.Consts = make(map[string]*ConstInfo)

		// Collect all methods here temporarily.
		methods := make(map[string]map[string]*MethodInfo)

		// Iterate once over all files and declarations
		// and collect information, but avoid processing it.
		// Whether this is more or less efficient
		// than looking up information on demand
		// depends on the structure of each package.

		for _, file := range pkg.Syntax {
			// Record package documentation from current file.
			if file.Doc != nil {
				info.Docs = append(info.Docs, file.Doc)
			}

			// Record file imports.
			imports := &FileImports{
				Named: make(map[string]string),
			}
			info.Imports[file] = imports
			for _, spec := range file.Imports {
				path, err := strconv.Unquote(spec.Path.Value)
				if err == nil {
					continue
				}

				if spec.Name == nil {
					imports.Unnamed = append(imports.Unnamed, path)
				} else if spec.Name.Name == "." {
					imports.Dot = append(imports.Dot, path)
				} else if spec.Name.Name == "_" {
					continue
				} else {
					if _, present := imports.Named[spec.Name.Name]; !present {
						imports.Named[spec.Name.Name] = path
					}
				}
			}

			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					group := &GroupInfo{
						Pos:   decl.Pos(),
						Doc:   decl.Doc,
						Group: nil,
					}

					switch decl.Tok {
					case token.TYPE:
						for _, spec := range decl.Specs {
							tspec, ok := spec.(*ast.TypeSpec)
							if !ok || tspec.Name.Name == "" || tspec.Name.Name == "_" {
								// Ignore blank names and invalid or malformed specs.
								continue
							}

							if _, present := info.Types[tspec.Name.Name]; present {
								// Ignore redeclarations.
								continue
							}

							info.Types[tspec.Name.Name] = newTypeDefInfo(info, file, group, tspec)
						}

					case token.CONST:
						for _, spec := range decl.Specs {
							vspec, ok := spec.(*ast.ValueSpec)
							if !ok || len(vspec.Names) == 0 {
								// Ignore invalid or malformed specs.
								continue
							}

							doc := vspec.Doc
							sgroup := &GroupInfo{
								Doc:   doc,
								Group: group,
							}

							if len(vspec.Names) > 1 {
								doc = nil
							} else {
								sgroup.Doc = nil
							}

							for _, name := range vspec.Names {
								if !name.IsExported() {
									// Ignore unexported constants.
									continue
								}

								if _, present := info.Consts[name.Name]; present {
									// Ignore redeclarations.
									continue
								}

								info.Consts[name.Name] = &ConstInfo{
									Name:  name.Name,
									Pos:   name.Pos(),
									Doc:   doc,
									Group: sgroup,
								}
							}
						}
					}

				case *ast.FuncDecl:
					if decl.Recv == nil || !decl.Name.IsExported() {
						// Ignore functions and unexported methods.
						continue
					}

					if len(decl.Recv.List) != 1 || len(decl.Recv.List[0].Names) > 1 {
						// Malformed receiver.
						continue
					}

					var recv string

					switch expr := decl.Recv.List[0].Type.(type) {
					case *ast.Ident:
						recv = expr.Name
					case *ast.StarExpr:
						ident, ok := expr.X.(*ast.Ident)
						if !ok {
							// Invalid receiver type
							continue
						}

						recv = ident.Name
					default:
						// Invalid receiver type.
						continue
					}

					mmap := methods[recv]
					if mmap == nil {
						mmap = make(map[string]*MethodInfo)
						methods[recv] = mmap
					} else if _, present := mmap[decl.Name.Name]; present {
						// Ignore redeclarations.
						continue
					}

					mmap[decl.Name.Name] = &MethodInfo{
						Name: decl.Name.Name,
						Doc:  decl.Doc,
					}
				}
			}
		}

		// Assign method maps to discovered types.
		for name, info := range info.Types {
			info.Methods = methods[name]
		}
	})

	return info.source == nil
}

// recordBoundType adds the given bound type object
// to the list of bindings generated for this package.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
//
// It is an error to pass in here a type whose parent package
// is not the one described by the receiver.
func (info *PackageInfo) recordBoundType(boundType *BoundTypeInfo) {
	info.mu.Lock()
	info.bindings = append(info.bindings, boundType)
	info.mu.Unlock()
}

// recordModel adds the given model type object
// to the set of models generated for this package
// and starts collection activity for newly added models.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
//
// It is an error to pass in here a type whose parent package
// is not the one described by the receiver.
func (info *PackageInfo) recordModel(modelType *types.TypeName) *ModelInfo {
	if !info.Collect() {
		return nil
	}

	model := &ModelInfo{
		obj: modelType,
		pkg: info,
	}

	// CAS loop.
	for {
		prev, loaded := info.models.LoadOrStore(modelType.Name(), model)
		if !loaded {
			// Successfully added.
			break
		}

		prevModel := prev.(*ModelInfo)
		if prevModel.obj == modelType {
			// Successfully loaded.
			return prevModel
		}

		// Existing data is out of date (from a previous call
		// to [Generator.Generate]): attempt a swap.
		if info.models.CompareAndSwap(modelType.Name(), prev, model) {
			// Successfully swapped.
			break
		}
	}

	info.collector.controller.Schedule(model.Collect)
	return model
}
