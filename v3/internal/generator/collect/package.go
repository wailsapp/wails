package collect

import (
	"cmp"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/tools/go/packages"
)

// PackageInfo records information about a package.
//
// Read accesses to fields Path, Name, Types, TypesInfo, Fset
// are safe at any time without any synchronisation.
//
// Read accesses to all other fields are only safe
// if a call to [PackageInfo.Collect] has completed before the access,
// for example by calling it in the accessing goroutine
// or before spawning the accessing goroutine.
//
// Concurrent write accesses are only allowed through the provided methods.
type PackageInfo struct {
	// Path holds the canonical path of the described package.
	Path string

	// Name holds the import name of the described package.
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

	// Includes holds a list of additional files to include
	// with the generated bindings.
	// It maps file names to their paths on disk.
	Includes map[string]string

	// Injections holds a list of code lines to be injected
	// into the package index file.
	Injections []string

	// services records service types that have to be generated for this package.
	// We rely upon [sync.Map] for atomic swapping support.
	// Keys are *types.TypeName, values are *ServiceInfo.
	services sync.Map

	// models records model types that have to be generated for this package.
	// We rely upon [sync.Map] for atomic swapping support.
	// Keys are *types.TypeName, values are *ModelInfo.
	models sync.Map

	// stats caches statistics about this package.
	stats atomic.Pointer[Stats]

	collector *Collector
	once      sync.Once
}

func newPackageInfo(pkg *packages.Package, collector *Collector) *PackageInfo {
	return &PackageInfo{
		Path: pkg.PkgPath,
		Name: pkg.Name,

		Types:     pkg.Types,
		TypesInfo: pkg.TypesInfo,

		Fset:  pkg.Fset,
		Files: pkg.Syntax,

		collector: collector,
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

		// Collect docs and parse directives.
		for _, file := range info.Files {
			if file.Doc == nil {
				continue
			}

			info.Docs = append(info.Docs, file.Doc)

			// Retrieve file directory.
			pos := info.Fset.Position(file.Pos())
			if !pos.IsValid() {
				collector.logger.Errorf(
					"package %s: found AST file with unknown path: `wails:include` directives from that file will be ignored",
					info.Path,
				)
			}
			dir := filepath.Dir(pos.Filename)

			// Parse directives.
			if info.Includes == nil {
				info.Includes = make(map[string]string)
			}
			for _, comment := range file.Doc.List {
				switch {
				case IsDirective(comment.Text, "inject"):
					// Check condition.
					line, cond, err := ParseCondition(ParseDirective(comment.Text, "inject"))
					if err != nil {
						collector.logger.Errorf(
							"%s: in `wails:inject` directive: %v",
							info.Fset.Position(comment.Pos()),
							err,
						)
						continue
					}

					if !cond.Satisfied(collector.options) {
						continue
					}

					// Record injected line.
					info.Injections = append(info.Injections, line)

				case pos.IsValid() && IsDirective(comment.Text, "include"):
					// Check condition.
					pattern, cond, err := ParseCondition(ParseDirective(comment.Text, "include"))
					if err != nil {
						collector.logger.Errorf(
							"%s: in `wails:include` directive: %v",
							info.Fset.Position(comment.Pos()),
							err,
						)
						continue
					}

					if !cond.Satisfied(collector.options) {
						continue
					}

					// Collect matching files.
					paths, err := filepath.Glob(filepath.Join(dir, pattern))
					if err != nil {
						collector.logger.Errorf(
							"%s: invalid pattern '%s' in `wails:include` directive: %v",
							info.Fset.Position(comment.Pos()),
							pattern,
							err,
						)
						continue
					} else if len(paths) == 0 {
						collector.logger.Warningf(
							"%s: pattern '%s' in `wails:include` directive matched no files",
							info.Fset.Position(comment.Pos()),
							pattern,
						)
						continue
					}

					// Announce and record matching files.
					for _, path := range paths {
						name := strings.ToLower(filepath.Base(path))
						if old, ok := info.Includes[name]; ok {
							collector.logger.Errorf(
								"%s: duplicate included file name '%s' in package %s; old path: '%s'; new path: '%s'",
								info.Fset.Position(comment.Pos()),
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

						info.Includes[name] = path
					}
				}
			}
		}
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
