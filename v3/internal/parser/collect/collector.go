package collect

import (
	"sync"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
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

	// mu protects access to the structs map.
	mu sync.Mutex
	// structs maps struct types to their [StructInfo].
	structs typeutil.Map

	// the omonymous package-level functions wrapped by sync.OnceFunc
	complexWarning func()
	chanWarning    func()
	funcWarning    func()
	genericWarning func()

	// wg is used to wait until concurrent model collection is complete.
	wg sync.WaitGroup
}

// NewCollector initialises a new Collector instance.
func NewCollector(loader Loader) *Collector {
	return &Collector{
		loader: loader,

		complexWarning: sync.OnceFunc(complexWarning),
		chanWarning:    sync.OnceFunc(chanWarning),
		funcWarning:    sync.OnceFunc(funcWarning),
		genericWarning: sync.OnceFunc(genericWarning),
	}
}
