package collect

import (
	"cmp"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/tools/go/packages"
)

type (
	// PackageInfo records information about a package.
	//
	// Read accesses to fields Path, Types, TypesInfo, Fset
	// are safe at any time without any synchronisation.
	//
	// Read accesses to all other fields are only safe
	// if a call to [PackageInfo.Collect] has completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	//
	// Concurrent write accesses are only allowed through the provided methods.
	PackageInfo struct {
		// Path holds the canonical path of the described package.
		Path string

		// Name holds the (possibly aliased) import name of the described package.
		Name string

		// Types and TypesInfo hold type information for this package.
		Types     *types.Package
		TypesInfo *types.Info

		// Fset holds the FileSet that was used to parse this package.
		Fset *token.FileSet

		// Files holds parsed files for this package,
		// ordered by start position to support binary search.
		Files []*ast.File

		// Docs holds package doc comments.
		Docs []*ast.CommentGroup

		// Internal is true if the package has been marked as internal.
		Internal bool

		// Includes holds a list of additional files to include
		// with the generated bindings.
		// It maps file names to their paths on disk and additional options.
		Includes map[string]IncludeInfo

		// Injections holds a list of code lines to be injected
		// into the package index file.
		Injections      []InjectionInfo
		NumJSInjections int
		NumTSInjections int

		// models records service types that have to be generated for this package.
		// We rely upon [sync.Map] for atomic swapping support.
		services sync.Map

		// models records model types that have to be generated for this package.
		// We rely upon [sync.Map] for atomic swapping support.
		models sync.Map

		// stats caches statistics about this package.
		stats atomic.Pointer[Stats]

		collector *Collector
		goFiles   []string
		once      sync.Once
	}

	// IncludeInfo records information about an included file.
	IncludeInfo struct {
		Path string
		JS   bool
		TS   bool
	}

	// InjectionInfo records information about an injected line.
	InjectionInfo struct {
		Code string
		JS   bool
		TS   bool
	}
)

func newPackageInfo(pkg *packages.Package, collector *Collector) *PackageInfo {
	return &PackageInfo{
		Path: pkg.PkgPath,
		Name: pkg.Name,

		Types:     pkg.Types,
		TypesInfo: pkg.TypesInfo,

		Fset:  pkg.Fset,
		Files: pkg.Syntax,

		collector: collector,
		goFiles:   pkg.GoFiles,
	}
}

// Package retrieves the the unique [PackageInfo] instance, if any,
// associated to the given package object within a Collector.
//
// Package is safe for concurrent use.
func (collector *Collector) Package(pkg *types.Package) *PackageInfo {
	return collector.pkgs[pkg]
}

// Iterate calls yield sequentially for each [PackageInfo] instance
// registered with the collector. If yield returns false,
// Iterate stops the iteration.
//
// Iterate is safe for concurrent use.
func (collector *Collector) Iterate(yield func(pkg *PackageInfo) bool) {
	for _, pkg := range collector.pkgs {
		if !yield(pkg) {
			return
		}
	}
}

// Stats returns cached statistics for this package.
// If [PackageInfo.Index] has not been called yet, it returns nil.
//
// Stats is safe for unsynchronised concurrent calls.
func (info *PackageInfo) Stats() *Stats {
	return info.stats.Load()
}

// Collect gathers information about the package described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *PackageInfo) Collect() *PackageInfo {
	if info == nil {
		return nil
	}

	info.once.Do(func() {
		collector := info.collector

		// Sort files by source position.
		if !slices.IsSortedFunc(info.Files, compareAstFiles) {
			info.Files = slices.Clone(info.Files)
			slices.SortFunc(info.Files, compareAstFiles)
		}

		var packageNameFound bool
		fset := token.NewFileSet()

		// Collect docs and parse packageName/include directives.
		// Here we can't use info.Files because
		//   - they are loaded from CompiledGoFiles, hence some file might be missing;
		//   - there is no robust way to retrieve their file path on disk;
		// instead, we range over GoFiles and parse their package clause only,
		// which should be fast enough.
		for _, goFile := range info.goFiles {
			file, err := parser.ParseFile(fset, goFile, nil, parser.PackageClauseOnly|parser.ParseComments)
			if err != nil {
				// Ignore failures silently:
				// they have been already reported by packages.Load.
				continue
			}

			if file.Doc == nil {
				continue
			}

			info.Docs = append(info.Docs, file.Doc)

			// Parse directives.
			dir := filepath.Dir(goFile)
			info.Includes = make(map[string]IncludeInfo)
			for _, comment := range file.Doc.List {
				switch {
				case IsDirective(comment.Text, "internal"):
					info.Internal = true

				case IsDirective(comment.Text, "inject"):
					// Record injected line.
					info.NumJSInjections++
					info.NumTSInjections++
					info.Injections = append(info.Injections, InjectionInfo{
						Code: ParseDirective(comment.Text, "inject"),
						JS:   true,
						TS:   true,
					})

				case IsDirective(comment.Text, "inject:js"):
					// Record injected line.
					info.NumJSInjections++
					info.Injections = append(info.Injections, InjectionInfo{
						Code: ParseDirective(comment.Text, "inject"),
						JS:   true,
						TS:   false,
					})

				case IsDirective(comment.Text, "inject:ts"):
					// Record injected line.
					info.NumTSInjections++
					info.Injections = append(info.Injections, InjectionInfo{
						Code: ParseDirective(comment.Text, "inject"),
						JS:   false,
						TS:   true,
					})

				case IsDirective(comment.Text, "include"):
					// Collect matching files.
					pattern := ParseDirective(comment.Text, "include")

					// Detect conditions.
					includeWithJS, includeWithTS := true, true
					if cond, rest, hasCond := strings.Cut(pattern, ":"); hasCond {
						pattern = rest
						switch cond {
						case "js":
							includeWithTS = false
						case "ts":
							includeWithJS = false
						default:
							collector.logger.Errorf(
								"%s: invalid condition '%s:' in `wails:include` directive: expected either 'js:' or 'ts:'",
								fset.Position(comment.Pos()),
								cond,
								err,
							)
							continue
						}
					}

					// Match files.
					paths, err := filepath.Glob(filepath.Join(dir, pattern))
					if err != nil {
						collector.logger.Errorf(
							"%s: invalid pattern '%s' in `wails:include` directive: %v",
							fset.Position(comment.Pos()),
							pattern,
							err,
						)
						continue
					} else if len(paths) == 0 {
						collector.logger.Warningf(
							"%s: pattern '%s' in `wails:include` directive matched no files",
							fset.Position(comment.Pos()),
							pattern,
						)
						continue
					}

					// Announce and record matching files.
					for _, path := range paths {
						name := filepath.Base(path)
						if old, ok := info.Includes[name]; ok {
							collector.logger.Errorf(
								"%s: duplicate included file name '%s' in package %s; old path: '%s'; new path: '%s'",
								fset.Position(comment.Pos()),
								name,
								info.Path,
								old,
								path,
							)
							continue
						}

						collector.logger.Debugf(
							"including file '%s' as '%s' in package %s",
							path,
							name,
							info.Path,
						)

						info.Includes[name] = IncludeInfo{
							Path: path,
							JS:   includeWithJS,
							TS:   includeWithTS,
						}
					}

				case !packageNameFound && IsDirective(comment.Text, "name"):
					packageName := ParseDirective(comment.Text, "name")
					if !token.IsIdentifier(packageName) {
						collector.logger.Errorf(
							"%s: invalid value in `wails:name` directive: '%s': expected a valid Go identifier",
							fset.Position(comment.Pos()),
							packageName,
						)
						continue
					}

					// Announce and record alias.
					collector.logger.Infof(
						"package %s: default package name '%s' replaced by '%s'",
						info.Path,
						info.Name,
						packageName,
					)
					info.Name = packageName
					packageNameFound = true
				}
			}
		}

		// Discard file path list.
		info.goFiles = nil
	})

	return info
}

// recordService adds the given service type object
// to the set of bindings generated for this package.
// It returns the unique [ServiceInfo] instance associated
// with the given type object.
//
// It is an error to pass in here a type whose parent package
// is not the one described by the receiver.
//
// recordService is safe for unsynchronised concurrent calls.
func (info *PackageInfo) recordService(obj *types.TypeName) *ServiceInfo {
	// Fetch current value, then add if not already present.
	service, _ := info.services.Load(obj)
	if service == nil {
		service, _ = info.services.LoadOrStore(obj, newServiceInfo(info.collector, obj))
	}
	return service.(*ServiceInfo)
}

// recordModel adds the given model type object
// to the set of models generated for this package.
// It returns the unique [ModelInfo] instance associated
// with the given type object. The present result is true
// if the model was already registered.
//
// It is an error to pass in here a type whose parent package
// is not the one described by the receiver.
//
// recordModel is safe for unsynchronised concurrent calls.
func (info *PackageInfo) recordModel(obj *types.TypeName) (model *ModelInfo, present bool) {
	// Fetch current value, then add if not already present.
	imodel, present := info.models.Load(obj)
	if imodel == nil {
		imodel, present = info.models.LoadOrStore(obj, newModelInfo(info.collector, obj))
	}
	return imodel.(*ModelInfo), present
}

// compareAstFiles compares two AST files by starting position.
func compareAstFiles(f1 *ast.File, f2 *ast.File) int {
	return cmp.Compare(f1.FileStart, f2.FileStart)
}
