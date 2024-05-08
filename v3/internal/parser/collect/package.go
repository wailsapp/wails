package collect

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

type (
	// PackageInfo records the following information about a package:
	// path, name, declaration groups with their doc comments,
	// type declarations with their doc comments and parent group,
	// constant declarations with their doc comments and parent group,
	// generated bindings and models.
	//
	// A dummy PackageInfo can be initialised with just the path;
	// all other fields will be populated upon calling [PackageInfo.Collect]
	// for the first time.
	//
	// Read accesses to the Path field are safe at any time
	// without any synchronisation.
	// Read accesses to all other exported fields are only safe
	// if a call to [PackageInfo.Collect] has been sequenced before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	//
	// Concurrent write accesses are only allowed through the provided methods.
	PackageInfo struct {
		Path string
		Name string

		Docs []*ast.CommentGroup

		Groups []*GroupInfo
		Types  map[string]*TypeDefInfo
		Consts map[string]*ConstInfo

		mu       sync.Mutex
		bindings []*types.TypeName
		models   []*types.TypeName

		source any
		once   sync.Once
	}

	// GroupInfo records information about a group
	// of type or constant declarations.
	// This may be either a list of distinct specifications
	// wrapped in parentheses, or a single specification
	// declaring multiple constants.
	GroupInfo struct {
		Doc   *ast.CommentGroup
		Group *GroupInfo
	}

	// TypeDefInfo records information about a single type specification.
	TypeDefInfo struct {
		Name    string
		Doc     *ast.CommentGroup
		Group   *GroupInfo
		Alias   bool
		Def     ast.Expr
		Methods map[string]*MethodInfo
	}

	// MethodInfo records information about a method declaration.
	MethodInfo struct {
		Name string
		Doc  *ast.CommentGroup
	}

	// ConstInfo records information about a constant declaration.
	ConstInfo struct {
		Name  string
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

		Bindings []*types.TypeName
		Models   []*types.TypeName
		Internal []*types.TypeName
	}
)

// Preload adds the given package descriptors to the collector,
// so that the loading step may be skipped when collecting information
// about each of those packages.
//
// Preload is safe for concurrent use.
func (collector *Collector) Preload(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		collector.pkgs.LoadOrStore(pkg.PkgPath, NewPackageInfo(pkg.PkgPath, pkg))
	}
}

// Package retrieves the the unique [PackageInfo] instance
// associated to the given path within a Collector.
// If none is present, a new one is initialised.
//
// Package is safe for concurrent use.
func (collector *Collector) Package(path string) *PackageInfo {
	info, _ := collector.pkgs.LoadOrStore(path, NewPackageInfo(path, collector.loader))
	return info.(*PackageInfo)
}

// All calls yield sequentially for each [PackageInfo] instance
// present in the collector. If yield returns false, All stops the iteration.
//
// All may be O(N) with the number of packages in the collector
// even if yield returns false after a constant number of calls.
//
// Package is safe for concurrent use.
func (collector *Collector) All(yield func(pkg *PackageInfo) bool) {
	collector.pkgs.Range(func(key, value any) bool {
		return yield(value.(*PackageInfo))
	})
}

// NewPackageInfo initialises an empty information structure
// for the given package path.
//
// source may be either a pointer to [packages.Package]
// for the given path with syntax information,
// or a [Loader] instance that will be used to load the package.
//
// The cost of this function must be as low as possible
// and it must not perform any significant work,
// as it might be called multiple times for the same package
// and its result might be discarded often.
func NewPackageInfo(path string, source any) *PackageInfo {
	if source == nil {
		panic("source cannot be nil")
	}

	return &PackageInfo{
		Path:   path,
		source: source,
	}
}

// AddBindings adds the given bound types
// to the list of bindings generated for this package.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
func (info *PackageInfo) AddBindings(bindings ...*types.TypeName) {
	info.mu.Lock()
	info.bindings = append(info.bindings, bindings...)
	info.mu.Unlock()
}

// AddModels adds the given model identifiers
// to the list of models generated for this package.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
func (info *PackageInfo) AddModels(models ...*types.TypeName) {
	info.mu.Lock()
	info.models = append(info.models, models...)
	info.mu.Unlock()
}

// IsEmpty retuns true if no bindings and models
// were generated for this package.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
func (info *PackageInfo) IsEmpty() bool {
	info.mu.Lock()
	result := len(info.bindings) == 0 && len(info.models) == 0
	info.mu.Unlock()
	return result
}

// Index computes a [PackageIndex] from the list
// of generated bindings and models.
//
// Binding and model names appear at most once
// in the returned structure.
//
// This method is safe to call even if [PackageInfo.Collect]
// has not been called yet.
func (info *PackageInfo) Index() (index PackageIndex) {
	info.mu.Lock()

	// Sort bindings by exported property and name, then deduplicate.
	slices.SortFunc(info.bindings, compareTypes)
	info.bindings = slices.CompactFunc(info.bindings, equateTypes)

	// Sort models by exported property and name, then deduplicate.
	slices.SortFunc(info.models, compareTypes)
	info.models = slices.CompactFunc(info.models, equateTypes)

	// Clone into result.
	index.Bindings = slices.Clone(info.bindings)
	index.Models = slices.Clone(info.models)

	info.mu.Unlock()

	// Find first unexported model.
	split, _ := slices.BinarySearchFunc(index.Models, struct{}{}, partitionTypes)

	// Separate unexported and exported models.
	index.Internal = index.Models[split:]
	index.Models = index.Models[:split]

	// Store package info.
	index.Info = info

	return
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
			pkg = info.source.(Loader).Load(info.Path)
			if pkg == nil {
				return
			}
		}

		// Discard package source.
		info.source = nil

		// Record package name.
		info.Name = pkg.Name

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
			if file.Doc != nil {
				info.Docs = append(info.Docs, file.Doc)
			}

			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					var group *GroupInfo
					empty := true

					if decl.Doc != nil {
						group = &GroupInfo{
							Doc:   decl.Doc,
							Group: nil,
						}
					}

					switch decl.Tok {
					case token.TYPE:
						for _, spec := range decl.Specs {
							tspec, ok := spec.(*ast.TypeSpec)
							if !ok || tspec.Name.Name == "" || tspec.Name.Name == "_" {
								// Ignore blank names and invalid or malformed specs.
								continue
							}

							empty = false
							info.Types[tspec.Name.Name] = &TypeDefInfo{
								Name:  tspec.Name.Name,
								Doc:   tspec.Doc,
								Group: group,
								Alias: tspec.Assign.IsValid(),
								Def:   tspec.Type,
							}
						}

					case token.CONST:
						for _, spec := range decl.Specs {
							vspec, ok := spec.(*ast.ValueSpec)
							if !ok || len(vspec.Names) == 0 {
								// Ignore invalid or malformed specs.
								continue
							}

							sgroup := group
							sempty := true

							doc := vspec.Doc
							if doc != nil && len(vspec.Names) > 1 {
								sgroup = &GroupInfo{
									Doc:   doc,
									Group: group,
								}
								doc = nil
							}

							for _, name := range vspec.Names {
								if !name.IsExported() {
									// Ignore unexported constants.
									continue
								}

								empty = false
								sempty = false
								info.Consts[name.Name] = &ConstInfo{
									Name:  name.Name,
									Doc:   doc,
									Group: sgroup,
								}
							}

							if !sempty && sgroup != group {
								info.Groups = append(info.Groups, sgroup)
							}
						}
					}

					if !empty && group != nil {
						info.Groups = append(info.Groups, group)
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

// compareModelNames compares two types by exported property and name.
// The order is exported before unexported, then lexicographical.
func compareTypes(m1 *types.TypeName, m2 *types.TypeName) int {
	if m1 == m2 {
		return 0
	}

	exp1, exp2 := m1.Exported(), m2.Exported()
	if exp1 && !exp2 {
		return -1
	} else if !exp1 && exp2 {
		return 1
	}

	return strings.Compare(m1.Name(), m2.Name())
}

// equateTypes tests efficiently whether two types have the same name.
func equateTypes(m1 *types.TypeName, m2 *types.TypeName) bool {
	// If the pointers are equal, then the names are too.
	// If the pointers differ, but they come from the same *types.Package,
	// then the names differ.
	return m1 == m2 || (m1.Pkg() != m2.Pkg() && m1.Name() == m2.Name())
}

// partitionTypes returns -1 if m is exported, 1 if unexported.
func partitionTypes(m *types.TypeName, _ struct{}) int {
	if m.Exported() {
		return -1
	}
	return 1
}
