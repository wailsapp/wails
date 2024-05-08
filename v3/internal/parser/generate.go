package parser

import (
	"errors"
	"go/types"
	"io"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/collect"
	"github.com/wailsapp/wails/v3/internal/parser/config"
	"github.com/wailsapp/wails/v3/internal/parser/render"
	"golang.org/x/tools/go/packages"
)

// ErrNoInitialPackages is returned by [Generator.Generate]
// when [LoadPackages] returns no error and no packages.
var ErrNoInitialPackages = errors.New("the given patterns matched no packages")

// Generator wraps all bookkeeping data structures that are needed
// to generate bindings for a set of packages.
type Generator struct {
	options *flags.GenerateBindingsOptions
	creator config.FileCreator

	collector *collect.Collector
	renderer  *render.Renderer

	controller controller
}

// NewGenerator configures a new generator instance.
// The options argument must not be nil.
// If creator is nil, no output file will be created.
// If logger is not nil, it is used to report messages interactively.
func NewGenerator(options *flags.GenerateBindingsOptions, creator config.FileCreator, logger config.Logger) *Generator {
	if creator == nil {
		creator = config.NullCreator
	}

	report := NewErrorReport(logger)

	return &Generator{
		options: options,
		creator: config.FileCreatorFunc(func(path string) (io.WriteCloser, error) {
			report.Debugf("writing output file %s", path)
			return creator.Create(path)
		}),

		controller: controller{ErrorReport: report},
	}
}

// Generate runs the binding generation process
// for the packages specified by the given patterns.
//
// Concurrent or repeated calls to Generate with the same receiver
// are not allowed.
//
// The stats return field is never nil.
//
// The error return field is nil in case of complete success (no warning).
// Otherwise, it may either report errors that occured while loading
// the initial set of packages, or errors returned by the static analyser,
// or be an [ErrorReport] instance.
//
// If error is an ErrorReport, it may have accumulated no errors, just warnings.
// When this is the case, all bindings have been generated successfully.
//
// Parsing/type-checking errors or errors encountered while writing
// individual files will be printed directly to the [config.Logger] instance
// provided during initialisation.
func (generator *Generator) Generate(patterns ...string) (stats *collect.Stats, err error) {
	stats = &collect.Stats{}
	stats.Start()
	defer stats.Stop()

	buildFlags, err := generator.options.BuildFlags()
	if err != nil {
		return
	}

	// Panic on repeated calls.
	if generator.collector != nil {
		panic("Generate() must not be called more than once on the same receiver")
	}

	// Cache reconstructed build flags.
	generator.controller.buildFlags = buildFlags

	// Initialise components.
	generator.collector = collect.NewCollector(&generator.controller)
	generator.renderer = render.NewRenderer(generator.options, generator.collector)

	// Package loading feedback.
	var lpkgMutex sync.Mutex
	generator.controller.Statusf("Loading packages...")
	go func() {
		time.Sleep(5 * time.Second)
		if lpkgMutex.TryLock() {
			generator.controller.Statusf("Loading packages... (this may take a long time)")
			lpkgMutex.Unlock()
		}
	}()

	// Resolve wails app pkg path.
	wailsAppPkgPaths, err := ResolvePatterns(buildFlags, WailsAppPkgPath)
	if err != nil {
		return
	}

	if len(wailsAppPkgPaths) < 1 {
		err = ErrNoApplicationPackage
		return
	} else if len(wailsAppPkgPaths) > 1 {
		// This should never happen...
		panic("wails application package path matched multiple packages")
	}

	// Load initial packages.
	pkgs, err := LoadPackages(buildFlags, patterns...)

	// Suppress package loading feedback.
	lpkgMutex.Lock()
	defer lpkgMutex.Unlock()

	// Check for loading errors.
	if err != nil {
		return
	}
	if len(patterns) > 0 && len(pkgs) == 0 {
		err = ErrNoInitialPackages
		return
	}

	// Report parsing/type-checking errors.
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			generator.controller.Warningf("%v", err)
		}
	}

	// Warmup collector.
	generator.collector.Preload(pkgs...)

	// Update status.
	generator.controller.Statusf("Looking for bound types...")
	bindingsFound := sync.OnceFunc(func() { generator.controller.Statusf("Generating bindings...") })

	// Run static analysis and schedule binding code generation for each result.
	err = FindServices(pkgs, wailsAppPkgPaths[0], &generator.controller, func(typ *types.TypeName) bool {
		bindingsFound()

		generator.controller.Schedule(func() {
			generator.generateBindings(typ)
		})

		return true
	})

	// Discard initial packages.
	pkgs = nil

	// Wait until all bindings have been generated and all models collected.
	generator.controller.Wait()

	// Check for analyser errors.
	if err != nil {
		return
	}

	// globalImports records all packages
	// that should be added to the global index.
	var globalImports []*collect.PackageInfo

	// Update status.
	if generator.options.NoIndex {
		generator.controller.Statusf("Generating models...")
	} else {
		generator.controller.Statusf("Generating models and index files...")
	}

	// Schedule models and index generation for each package.
	generator.collector.Iterate(func(info *collect.PackageInfo) bool {
		if info.IsEmpty() {
			return true
		}

		if !generator.options.NoIndex {
			globalImports = append(globalImports, info)
		}

		generator.controller.Schedule(func() {
			generator.generateModelsAndIndex(info)
		})
		return true
	})

	// Generate global index and shortcuts.
	if len(globalImports) > 0 {
		generator.generateGlobalIndex(globalImports)
	}

	// Wait until all models and indices have been generated.
	generator.controller.Wait()

	// Populate stats.
	generator.controller.Statusf("Collecting stats...")
	generator.collector.Iterate(func(info *collect.PackageInfo) bool {
		if info.IsEmpty() {
			stats.NumPackages++
		} else {
			stats.Add(info.Stats())
		}
		return true
	})

	// Return non-empty error report.
	if generator.controller.HasErrors() || generator.controller.HasWarnings() {
		err = generator.controller.ErrorReport
	}

	return
}

// generateModelsAndIndex schedules generation of public/private model files
// and if required by the options, of index files.
// for the given package.
func (generator *Generator) generateModelsAndIndex(info *collect.PackageInfo) {
	index := info.Index()
	empty := len(index.Bindings) == 0

	// Collect package information.
	if !info.Collect() {
		generator.controller.Errorf("package %s: models and index generation failed", info.Path)
		return
	}

	// Now that Collect has been called, goroutines spawned below
	// can access package information freely.

	if len(index.Models) > 0 {
		empty = false
		generator.controller.Schedule(func() {
			generator.generateModels(info, index.Models, false)
		})
	}

	if len(index.Internal) > 0 {
		empty = false
		generator.controller.Schedule(func() {
			generator.generateModels(info, index.Internal, true)
		})
	}

	if !(generator.options.NoIndex || empty) {
		generator.controller.Schedule(func() {
			generator.generateIndex(index)
		})
		generator.reportDualRoles(index)
	}
}

// reportDualRoles checks for models that are also bound types
// and emits a warning.
func (generator *Generator) reportDualRoles(index collect.PackageIndex) {
	bindings, models := index.Bindings, index.Models
	for len(bindings) > 0 && len(models) > 0 {
		if bindings[0].Name < models[0].Name {
			bindings = bindings[1:]
		} else if bindings[0].Name > models[0].Name {
			models = models[1:]
		} else {
			generator.controller.Warningf(
				"package %s: type %s has been marked both as a bound type and as a model; shadowing between the two may take place when importing generated JS indexes",
				index.Package.Path,
				bindings[0].Name,
			)

			bindings = bindings[1:]
			models = models[1:]
		}
	}
}

// controller provides an implementation of the interface [collect.Controller].
type controller struct {
	sync.WaitGroup

	*ErrorReport

	// buildFlags caches parsed build flags from the options struct.
	buildFlags []string
}

// Load loads the given package path in syntax-only mode.
// In case of errors, it returns nil and adds them to the error report.
func (ctrl *controller) Load(path string) *packages.Package {
	pkgs, err := LoadPackages(ctrl.buildFlags, path)
	if err != nil {
		ctrl.Errorf("%v", err)
		return nil
	} else if len(pkgs) < 1 {
		ctrl.Errorf("%s: package not found", path)
		return nil
	} else if len(pkgs) > 1 {
		ctrl.Errorf("%s: multiple packages loaded for the same path", path)
		return nil
	}

	for _, err := range pkgs[0].Errors {
		ctrl.Warningf("%v", err)
	}

	return pkgs[0]
}

// Schedule runs the given function concurrently,
// tracking it on the controller's wait group.
func (ctrl *controller) Schedule(task func()) {
	ctrl.Add(1)
	go func() {
		defer ctrl.Done()
		task()
	}()
}
