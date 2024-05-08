package parser

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/analyse"
	"github.com/wailsapp/wails/v3/internal/parser/collect"
	"github.com/wailsapp/wails/v3/internal/parser/render"
	"golang.org/x/tools/go/packages"
)

// ErrNoInitialPackages is returned by [Generator.Generate]
// when [LoadPackages] returns no error and no packages.
var ErrNoInitialPackages = errors.New("the given patterns matched no packages")

// FileCreator abstracts away file and directory creation.
// We use this to implement tests cleanly.
//
// The default implementation creates a file with the given path
// by calling first [os.MkdirAll] on the directory part
// then [os.Create] on the full path.
//
// Paths are always relative to the output directory.
//
// A FileCreator must allow concurrent calls to Create transparently.
// Each [io.WriteCloser] instance returned by a call to Create
// will be used by one goroutine at a time; but distinct instances
// must support concurrent use by distinct goroutines.
type FileCreator interface {
	Create(path string) (io.WriteCloser, error)
}

// FileCreatorFunc is an adapter to allow
// the use of ordinary functions as file creators.
type FileCreatorFunc func(path string) (io.WriteCloser, error)

// Create calls f(path).
func (f FileCreatorFunc) Create(path string) (io.WriteCloser, error) {
	return f(path)
}

// defaultCreator implements the default file creation strategy.
// It joins the output directory and the given path,
// calls [os.MkdirAll] on the directory part,
// then [os.Create] on the full path.
func defaultCreator(outputDir string) FileCreator {
	return FileCreatorFunc(func(path string) (io.WriteCloser, error) {
		path = filepath.Join(outputDir, path)

		if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
			return nil, err
		}

		return os.Create(path)
	})
}

// Generator wraps all bookkeeping data structures that are needed
// to generate bindings for a set of packages.
type Generator struct {
	options *flags.GenerateBindingsOptions
	creator FileCreator

	// buildFlags caches parsed build flags from the options struct.
	buildFlags []string

	collector *collect.Collector
	renderer  *render.Renderer

	wg sync.WaitGroup
}

// NewGenerator configures a new generator instance.
// The options argument must not be nil.
// If creator is nil, the default implementation will be used.
func NewGenerator(options *flags.GenerateBindingsOptions, creator FileCreator) *Generator {
	if creator == nil {
		creator = defaultCreator(options.OutputDirectory)
	}

	return &Generator{
		options: options,
		creator: creator,
	}
}

// Generate generates bindings for the packages specified by the given patterns.
// Generate can be called multiple times with different or even overlapping
// sets of packages; however, changes to package files that happen
// between calls to Generate may not be detected.
//
// Concurrent calls to Generate are not allowed.
//
// The return value reports errors that may occur while loading
// the initial set of packages and starting the static analyser.
//
// Parsing/type-checking errors or errors encountered while writing
// individual files will be printed directly to the pterm Error logger.
func (generator *Generator) Generate(patterns ...string) (stats *collect.Stats, err error) {
	stats = &collect.Stats{}
	stats.Start()
	defer stats.Stop()

	buildFlags, err := generator.options.BuildFlags()
	if err != nil {
		return
	}

	// Cache reconstructed build flags.
	generator.buildFlags = buildFlags

	// Initialise collector.
	if generator.collector == nil {
		generator.collector = collect.NewCollector(collect.LoaderFunc(generator.loadAdditionalPackage))
	}

	// Initialise renderer.
	if generator.renderer == nil {
		generator.renderer = render.NewRenderer(generator.options)
	}

	// Load initial packages.
	pkgs, err := LoadPackages(buildFlags, true, patterns...)
	if err != nil {
		return
	}
	if len(patterns) > 0 && len(pkgs) == 0 {
		err = ErrNoInitialPackages
		return
	}

	// Report parsing/type-checking errors and record initial packages.
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			pterm.Warning.Println(err)
		}
	}

	// Warmup collector.
	generator.collector.Preload(pkgs...)

	// Run analyser and schedule bindings generation for each result.
	err = analyse.NewAnalyser(pkgs).Run(func(result analyse.Result) bool {
		generator.wg.Add(1)
		go generator.generateBindings(result)
		return true
	})
	if err != nil {
		return
	}

	// Discard unneeded packages.
	pkgs = nil

	// Wait until all bindings have been generated and all models collected.
	generator.wg.Wait()
	generator.collector.WaitForModels()

	// Record all packages that should be added to the global index.
	var globalImports []*collect.PackageInfo

	// Schedule models and index generation for each package.
	generator.collector.Iterate(func(info *collect.PackageInfo) bool {
		if info.IsEmpty() {
			return true
		}

		if !generator.options.NoIndex {
			globalImports = append(globalImports, info)
		}

		generator.wg.Add(1)
		go generator.generateModelsAndIndex(info)
		return true
	})

	// Generate global index and shortcuts.
	if len(globalImports) > 0 {
		generator.wg.Add(1)
		go generator.generateGlobalIndex(globalImports)
	}

	// Wait until all models and indices have been generated.
	generator.wg.Wait()

	// Populate stats.
	generator.collector.Iterate(func(info *collect.PackageInfo) bool {
		if info.IsEmpty() {
			stats.NumPackages++
		} else {
			stats.Add(info.Stats())
		}
		return true
	})

	return
}

// loadAdditionalPackage loads syntax for the specified package path.
// Errors are printed to the pterm Error logger.
// When an error occurs, loadAdditionalPackage returns nil.
func (generator *Generator) loadAdditionalPackage(path string) *packages.Package {
	pkgs, err := LoadPackages(generator.buildFlags, false, path)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	} else if len(pkgs) < 1 {
		pterm.Error.Printfln("%s: package not found", path)
		return nil
	} else if len(pkgs) > 1 {
		pterm.Error.Printfln("%s: multiple packages loaded for the same path", path)
		return nil
	}

	for _, err := range pkgs[0].Errors {
		pterm.Warning.Println(err)
	}

	return pkgs[0]
}

// generateModelsAndIndex schedules generation of public/private model files
// and if required by the options, of index files.
// for the given package.
func (generator *Generator) generateModelsAndIndex(info *collect.PackageInfo) {
	defer generator.wg.Done()

	index := info.Index()
	empty := len(index.Bindings) == 0

	// Collect package information.
	if !info.Collect() {
		pterm.Error.Printfln("package %s: models and index generation failed", info.Path)
		return
	}

	// Now that Collect has been called, goroutines spawned below
	// can access package information freely.

	if len(index.Models) > 0 {
		empty = false
		generator.wg.Add(1)
		go generator.generateModels(info, index.Models, false)
	}

	if len(index.Internal) > 0 {
		empty = false
		generator.wg.Add(1)
		go generator.generateModels(info, index.Internal, true)
	}

	if !(generator.options.NoIndex || empty) {
		generator.wg.Add(1)
		go generator.generateIndex(index)
		reportDualRoles(index)
	}
}

// reportDualRoles checks for models that are also bound types
// and emits a warning.
func reportDualRoles(index collect.PackageIndex) {
	bindings, models := index.Bindings, index.Models
	for len(bindings) > 0 && len(models) > 0 {
		if bindings[0].Name < models[0].Name {
			bindings = bindings[1:]
		} else if bindings[0].Name > models[0].Name {
			models = models[1:]
		} else {
			pterm.Warning.Printfln(
				"package %s: type %s has been marked both as a bound type and as a model; shadowing between the two may take place when importing generated JS indexes",
				index.Info.Path,
				bindings[0].Name,
			)

			bindings = bindings[1:]
			models = models[1:]
		}
	}
}
