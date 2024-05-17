package collect

import (
	"go/ast"
	"go/types"
	"sync"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/config"
	"golang.org/x/tools/go/packages"
)

// Scheduler instances provide task scheduling
// for collection activities.
type Scheduler interface {
	// Schedule should run the given function according
	// to the implementation's preferred strategy.
	//
	// Scheduled tasks may call Schedule again;
	// therefore, if tasks run concurrently,
	// the implementation must support concurrent calls.
	Schedule(task func())
}

// Info instances provide information about either
// a type-checker object, a struct type or a group of declarations.
type Info interface {
	Object() types.Object
	Type() types.Type
	Node() ast.Node
}

// Collector wraps all bookkeeping data structures that are needed
// to collect data about a set of packages, bindings and models.
type Collector struct {
	// pkgs caches packages that have been registered for collection.
	pkgs map[*types.Package]*PackageInfo

	// cache caches collected information about type-checker objects
	// and declaration groups. Elements are [Info] instances.
	cache sync.Map

	systemPaths *config.SystemPaths
	options     *flags.GenerateBindingsOptions
	scheduler   Scheduler
	logger      config.Logger
}

// NewCollector initialises a new Collector instance for the given package set.
func NewCollector(pkgs []*packages.Package, systemPaths *config.SystemPaths, options *flags.GenerateBindingsOptions, scheduler Scheduler, logger config.Logger) *Collector {
	collector := &Collector{
		pkgs: make(map[*types.Package]*PackageInfo, len(pkgs)),

		systemPaths: systemPaths,
		options:     options,
		scheduler:   scheduler,
		logger:      logger,
	}

	// Register packages.
	for _, pkg := range pkgs {
		collector.pkgs[pkg.Types] = newPackageInfo(pkg, collector)
	}

	return collector
}

// fromCache returns the cached Info instance associated
// to the given type-checker object or declaration group.
// If none exists, a new one is created.
func (collector *Collector) fromCache(objectOrGroup any) Info {
	entry, ok := collector.cache.Load(objectOrGroup)
	info, _ := entry.(Info)

	if !ok {
		switch x := objectOrGroup.(type) {
		case *ast.GenDecl, *ast.ValueSpec, *ast.Field:
			info = newGroupInfo(x.(ast.Node))

		case *types.Const:
			info = newConstInfo(collector, x)

		case *types.Func:
			info = newMethodInfo(collector, x.Origin())

		case *types.TypeName:
			info = newTypeInfo(collector, x)

		case *types.Var:
			if !x.IsField() {
				panic("cache lookup for invalid object kind")
			}

			info = newFieldInfo(collector, x.Origin())

		case *types.Struct:
			info = newStructInfo(collector, x)

		default:
			panic("cache lookup for invalid object kind")
		}

		entry, _ = collector.cache.LoadOrStore(objectOrGroup, info)
		info, _ = entry.(Info)
	}

	return info
}
