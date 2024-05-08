package collect

import (
	"sync"

	"golang.org/x/tools/go/packages"
)

// Loader abstracts the process of loading single packages.
// Only syntax is required. In case of errors, Load should return nil.
type Loader interface {
	Load(path string) *packages.Package
}

// LoaderFunc is an adapter to allow the use of ordinary functions as loaders.
type LoaderFunc func(path string) *packages.Package

// Load calls f(path).
func (f LoaderFunc) Load(path string) *packages.Package {
	return f(path)
}

// Collector wraps all bookkeeping data structures that are needed
// to collect data about a set of packages, bindings and models.
type Collector struct {
	loader Loader
	pkgs   sync.Map
}

// NewCollector initialises a new Collector instance.
func NewCollector(loader Loader) *Collector {
	return &Collector{
		loader: loader,
	}
}

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

// Package retrieves from the the unique [PackageInfo] instance
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
