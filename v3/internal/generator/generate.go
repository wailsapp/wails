package generator

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/collect"
	"github.com/wailsapp/wails/v3/internal/generator/config"
	"github.com/wailsapp/wails/v3/internal/generator/render"
)

// Generator wraps all bookkeeping data structures that are needed
// to generate bindings for a set of packages.
type Generator struct {
	options *flags.GenerateBindingsOptions
	creator config.FileCreator

	// serviceFiles maps service file paths to their type object.
	// It is used for lower/upper-case collision detection.
	// Keys are strings, values are *types.TypeName.
	serviceFiles sync.Map

	collector *collect.Collector
	renderer  *render.Renderer

	logger    *ErrorReport
	scheduler scheduler
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

		logger: report,
	}
}

// Generate runs the binding generation process
// for the packages specified by the given patterns.
//
// Concurrent or repeated calls to Generate with the same receiver
// are not allowed.
//
// The stats result field is never nil.
//
// The error result field is nil in case of complete success (no warning).
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

	// Validate file names.
	err = generator.validateFileNames()
	if err != nil {
		return
	}

	// Parse build flags.
	buildFlags, err := generator.options.BuildFlags()
	if err != nil {
		return
	}

	// Start package loading feedback.
	var lpkgMutex sync.Mutex
	generator.logger.Statusf("Loading packages...")
	go func() {
		time.Sleep(5 * time.Second)
		if lpkgMutex.TryLock() {
			generator.logger.Statusf("Loading packages... (this may take a long time)")
			lpkgMutex.Unlock()
		}
	}()

	systemPaths, err := ResolveSystemPaths(buildFlags)
	if err != nil {
		return
	}

	// Load initial packages.
	pkgs, err := LoadPackages(buildFlags, patterns...)

	// Suppress package loading feedback.
	lpkgMutex.Lock()

	// Check for loading errors.
	if err != nil {
		return
	}
	if len(patterns) > 0 && len(pkgs) == 0 {
		err = ErrNoPackages
		return
	}

	// Report parsing/type-checking errors.
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			generator.logger.Warningf("%v", err)
		}
	}

	// Panic on repeated calls.
	if generator.collector != nil {
		panic("Generate() must not be called more than once on the same receiver")
	}

	// Initialise subcomponents.
	generator.collector = collect.NewCollector(pkgs, systemPaths, generator.options, &generator.scheduler, generator.logger)
	generator.renderer = render.NewRenderer(generator.options, generator.collector)

	// Update status.
	generator.logger.Statusf("Looking for services...")
	serviceFound := sync.OnceFunc(func() { generator.logger.Statusf("Generating service bindings...") })

	// Run static analysis.
	services, err := FindServices(pkgs, systemPaths, generator.logger)

	// Check for analyser errors.
	if err != nil {
		return
	}

	// Discard unneeded data.
	pkgs = nil

	// Schedule code generation for each found service.
	for obj := range services {
		serviceFound()
		generator.scheduler.Schedule(func() {
			generator.generateService(obj)
		})
	}

	// Wait until all services have been generated and all models collected.
	generator.scheduler.Wait()

	// Invariants:
	//   - Service files have been generated for all discovered services;
	//   - ModelInfo.Collect has been called on all discovered models, and therefore
	//   - all required models have been discovered.

	// Update status.
	if generator.options.NoIndex {
		generator.logger.Statusf("Generating models...")
	} else {
		generator.logger.Statusf("Generating models and index files...")
	}

	// Schedule models, index and included files generation for each package.
	for info := range generator.collector.Iterate {
		generator.scheduler.Schedule(func() {
			generator.generateModelsIndexIncludes(info)
		})
	}

	// Wait until all models and indices have been generated.
	generator.scheduler.Wait()

	// Populate stats.
	generator.logger.Statusf("Collecting stats...")
	for info := range generator.collector.Iterate {
		stats.Add(info.Stats())
	}

	// Return non-empty error report.
	if generator.logger.HasErrors() || generator.logger.HasWarnings() {
		err = generator.logger
	}

	return
}

// generateModelsIndexIncludes schedules generation of public/private model files,
// included files and, if allowed by the options,
// of an index file for the given package.
func (generator *Generator) generateModelsIndexIncludes(info *collect.PackageInfo) {
	index := info.Index(generator.options.TS)

	// info.Index implies info.Collect: goroutines spawned below
	// can access package information freely.

	if len(index.Models) > 0 {
		generator.scheduler.Schedule(func() {
			generator.generateModels(info, index.Models)
		})
	}

	if len(index.Package.Includes) > 0 {
		generator.scheduler.Schedule(func() {
			generator.generateIncludes(index)
		})
	}

	if !generator.options.NoIndex && !index.IsEmpty() {
		generator.generateIndex(index)
	}
}

// validateFileNames validates user-provided filenames.
func (generator *Generator) validateFileNames() error {
	switch {
	case generator.options.ModelsFilename == "":
		return fmt.Errorf("models filename must not be empty")

	case !generator.options.NoIndex && generator.options.IndexFilename == "":
		return fmt.Errorf("package index filename must not be empty")

	case generator.options.ModelsFilename != strings.ToLower(generator.options.ModelsFilename):
		return fmt.Errorf("models filename must not contain uppercase characters")

	case generator.options.IndexFilename != strings.ToLower(generator.options.IndexFilename):
		return fmt.Errorf("package index filename must not contain uppercase characters")

	case !generator.options.NoIndex && generator.options.ModelsFilename == generator.options.IndexFilename:
		return fmt.Errorf("models and package indexes cannot share the same filename")
	}

	return nil
}

// scheduler provides an implementation of the [collect.Scheduler] interface.
type scheduler struct {
	sync.WaitGroup
}

// Schedule runs the given function concurrently,
// registering it on the scheduler's wait group.
func (sched *scheduler) Schedule(task func()) {
	sched.Add(1)
	go func() {
		defer sched.Done()
		task()
	}()
}
