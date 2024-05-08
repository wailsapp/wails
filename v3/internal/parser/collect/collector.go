package collect

import (
	"sync"

	"github.com/wailsapp/wails/v3/internal/parser/config"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

// A Controller instance provides package loading, task scheduling
// and message logging. All Controller methods may be called
// concurrently by a [Collector].
type Controller interface {
	config.Logger

	// Load should load the given package path in syntax-only mode.
	// In case of errors, it should return nil and handle them internally.
	Load(path string) *packages.Package

	// Schedule should run the given function concurrently.
	// It gives the controller an opportunity
	// to track the progress of collection activities.
	Schedule(task func())
}

// Collector wraps all bookkeeping data structures that are needed
// to collect data about a set of packages, bindings and models.
type Collector struct {
	controller Controller

	// pkgs caches packages that have been registered for collection.
	// The element type must be *PackageInfo.
	pkgs sync.Map

	// mu protects access to the structs map.
	mu sync.Mutex
	// structs maps struct types to their [*StructInfo].
	structs typeutil.Map
}

// NewCollector initialises a new Collector instance.
func NewCollector(controller Controller) *Collector {
	return &Collector{
		controller: controller,
	}
}
