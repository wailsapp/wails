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
)

// ErrNoInitialPackages is returned by [Generator.Generate]
// when [LoadPackages] returns no error and no packages.
var ErrNoInitialPackages = errors.New("the given patterns matched no packages")

// CreateFileFunc abstracts away file and directory creation.
// We use this to implement tests cleanly.
//
// The default implementation creates a file with the given path
// by calling first [os.MkdirAll] on the directory part
// then [os.Create] on the full path.
//
// A CreateFileFunc must allow concurrent calls transparently.
// Each [io.WriteCloser] instance returned by a call to CreateFileFunc
// will be used by one goroutine at a time; but distinct instances
// must allow be concurrent use by distinct goroutines.
type CreateFileFunc func(path string) (io.WriteCloser, error)

func defaultCreate(path string) (io.WriteCloser, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// Generator wraps all bookkeeping data structures that are needed
// to generate bindings for a set of packages.
type Generator struct {
	options *flags.GenerateBindingsOptions
	create  CreateFileFunc

	// buildFlags caches build flags for use in the [PackageInfo.Collect] method.
	buildFlags []string

	// pkgs maps package paths to their unique [PackageInfo] instances.
	pkgs PackageMap

	wg sync.WaitGroup
}

// NewGenerator configures a new generator instance.
// The options argument must not be nil.
// If create is nil, it will use the default implementation.
func NewGenerator(options *flags.GenerateBindingsOptions, create CreateFileFunc) *Generator {
	if create == nil {
		create = defaultCreate
	}

	return &Generator{
		options: options,
		create:  create,
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
func (generator *Generator) Generate(patterns ...string) error {
	buildFlags, err := generator.options.BuildFlags()
	if err != nil {
		return err
	}

	// Cache reconstructed build flags.
	generator.buildFlags = buildFlags

	// Load initial packages.
	pkgs, err := LoadPackages(buildFlags, true, patterns...)
	if err != nil {
		return err
	}

	if len(patterns) > 0 && len(pkgs) == 0 {
		return ErrNoInitialPackages
	}

	// Report parsing/type-checking errors
	// and initialise package map.
	for _, pkg := range pkgs {
		generator.pkgs.Preload(pkg)
		for _, err := range pkg.Errors {
			pterm.Warning.Println(err)
		}
	}

	// Run analyser and schedule bindings generation for each result.
	err = analyse.NewAnalyser(pkgs).Run(func(result analyse.Result) bool {
		generator.wg.Add(1)
		go generator.generateBindings(result)
		return true
	})
	if err != nil {
		return err
	}

	// Discard unneeded packages.
	pkgs = nil

	// Wait until all bindings have been generated and all models collected.
	generator.wg.Wait()

	// Schedule models and index generation for each package.
	generator.pkgs.Range(func(info *PackageInfo) bool {
		generator.wg.Add(1)
		go generator.generateModelsAndIndex(info)
		return true
	})

	// Wait until all models have been generated.
	generator.wg.Wait()

	return nil
}

// generateModelsAndIndex schedules generation of public/private model files
// and if required by the options, of index files
// for the given package.
func (generator *Generator) generateModelsAndIndex(info *PackageInfo) {
	index := info.Index()
	writeIndex := !generator.options.NoIndex && len(index.Bindings) > 0 && len(index.Models) > 0

	// Collect package information if it is going to be consumed below.
	if writeIndex || len(index.Models) > 0 || len(index.Internal) > 0 {
		if err := info.Collect(generator.buildFlags); err != nil {
			pterm.Error.Println(err)
			pterm.Error.Printfln("package %s: models and index generation failed", info.Path)
			return
		}
	}

	// Now that Collect has been called, the goroutines spawned below
	// can access package information freely.

	if len(index.Models) > 0 {
		generator.wg.Add(1)
		go generator.generateModels(info, index.Models, false)
	}

	if len(index.Internal) > 0 {
		generator.wg.Add(1)
		go generator.generateModels(info, index.Internal, true)
	}

	if writeIndex {
		generator.generateIndex(index)
	}

	generator.wg.Done()
}
