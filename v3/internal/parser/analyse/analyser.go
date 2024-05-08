package analyse

import (
	"go/types"
	"sync"

	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

// Analyser instances wrap all bookkeeping data structures that are needed
// to analyse a set of packages and discover bound types.
type Analyser struct {
	// pkgs is the list of packages under analysis.
	pkgs []*packages.Package
	// refs stores a ref map for each package in pkgs (see [RefMap], [BuildRefMap]).
	refs []RefMap

	// paths holds a collection of known field paths used for canonicalisation.
	paths PathSet

	// root holds the root target for the static analyser,
	// as found in the dependency graph rooted at pkgs.
	root RootTarget
	// queue is a list of additional targets that need to be analysed (see [target]).
	queue []target
	// scheduled holds all targets that have been already added to the queue
	// and must be ignored if rediscovered.
	scheduled Set[target]

	// found holds at any time the set of all bound types found so far.
	found Set[Result]
	// yield is invoked immediately (if non-nil) when a new result is found.
	yield func(Result)
}

// New allocates and initialises a static analyser instance
// for the given set of packages.
func NewAnalyser(pkgs []*packages.Package) *Analyser {
	return &Analyser{pkgs: pkgs}
}

// SetHasher sets the type hasher used by the internal [PathSet].
//
// The hasher is used by [PathSet.TypeAssertionStep]
// to compute path elements representing type assertions.
//
// A single Hasher created by may be shared among many Analysers.
// See [typeutil.Map.SetHasher] for more information.
func (analyser *Analyser) SetHasher(hasher typeutil.Hasher) {
	analyser.paths.SetHasher(hasher)
}

// Result returns a slice listing all bound types discovered by the analysis.
func (analyser *Analyser) Results() []Result {
	return lo.Keys(analyser.found)
}

// Run performs the static analysis. If yield is non-nil,
// it is invoked immediately upon discovery of a new bound type.
// If yield returns false, the analyser stops immediately.
// This allows consumers to start processing results
// before the analysis is finished.
//
// During the warm-up phase, Run sorts the Syntax slices of all input packages,
// hence concurrent reads will result in data races.
// After yield is called for the first time, the analyser
// will never again modify any field of the input package structs,
// hence concurrent reads (but not writes) become safe.
// Writes become safe again only after Run returns.
func (analyser *Analyser) Run(yield func(Result) bool) (err error) {
	stop := false
	defer func() {
		if stop {
			recover()
		}
	}()

	if len(analyser.pkgs) == 0 {
		return ErrNoApplicationPackage
	}

	// Setup yield function.
	if yield != nil {
		analyser.yield = func(result Result) {
			stop = !yield(result)
			if stop {
				panic("stop requested by consumer")
			}
		}
	}

	// Each package might spawn analysis tasks on any other package,
	// hence we need all AST files sorted and all ref maps ready
	// _before_ starting the analysis.
	// This work can be performed in parallel if other CPU cores are available.

	// Initialise wait group for concurrent initialisation tasks.
	var wg sync.WaitGroup
	wg.Add(2 * len(analyser.pkgs))

	// Allocate slice for ref maps.
	analyser.refs = make([]RefMap, len(analyser.pkgs))

	for i, pkg := range analyser.pkgs {
		// Instantiate new variables for the closures below.
		ii, ppkg := i, pkg
		go func() {
			SortAstFiles(ppkg)
			wg.Done()
		}()
		go func() {
			analyser.refs[ii] = BuildRefMap(ppkg)
			wg.Done()
		}()
	}

	// The analyser looks for assignments to the struct field
	// application.Options.Bind from the Wails application package.
	// Search the current dependency graph for the corresponding object.
	analyser.root, err = FindRootTarget(lo.Map(analyser.pkgs, func(pkg *packages.Package, _ int) *types.Package { return pkg.Types }))

	// Wait until all concurrent initialization tasks are complete.
	// We must wait even in case of failure to ensure
	// the caller gets back exclusive access to the package list.
	wg.Wait()

	if err != nil {
		// If the root field cannot be found, fail immediately
		return
	}

	// Initialise bookkeeping fields.
	analyser.scheduled = make(Set[target])
	analyser.found = make(Set[Result])

	// Analyse assignments to the root target.
	analyser.processRootTarget()

	// Keep performing additional scheduled work until the queue turns out empty.
	var queue []target
	for len(analyser.queue) > 0 {
		queue, analyser.queue = analyser.queue, queue[:0]
		for _, tgt := range queue {
			analyser.processTarget(tgt)
		}
	}

	return
}
