package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator"
	"github.com/wailsapp/wails/v3/internal/generator/config"
	"github.com/wailsapp/wails/v3/internal/term"
)

func GenerateBindings(options *flags.GenerateBindingsOptions, patterns []string) error {
	DisableFooter = true

	if options.Silent {
		term.DisableOutput()
		defer term.EnableOutput()
	} else if options.Verbose {
		term.EnableDebug()
		defer term.DisableDebug()
	}

	term.Header("Generate Bindings")

	if len(patterns) == 0 {
		// No input pattern, load package from current directory.
		patterns = []string{"."}
	}

	// Compute absolute path of output directory.
	absPath, err := filepath.Abs(options.OutputDirectory)
	if err != nil {
		return err
	}

	// When clean mode is active and we're writing real files, generate bindings
	// into a dot-prefixed sibling temp directory first, then atomically swap it
	// into place with a single os.Rename. This prevents chokidar (used by Vite)
	// from entering a rename-event loop caused by rapid directory delete+recreate,
	// which would otherwise cause the node process to leak memory at ~2-6 MB/s.
	// Dot-prefixed directories are ignored by chokidar's default glob pattern,
	// so no spurious HMR events fire during file generation.
	generationDir := absPath
	if options.Clean && !options.DryRun {
		if err := os.MkdirAll(filepath.Dir(absPath), 0o777); err != nil {
			return fmt.Errorf("failed to create bindings parent directory: %w", err)
		}
		tmpDir, err := os.MkdirTemp(filepath.Dir(absPath), ".bindings-tmp-")
		if err != nil {
			return fmt.Errorf("failed to create temp directory for bindings: %w", err)
		}
		generationDir = tmpDir
		defer func() {
			// Clean up temp dir on error (rename sets generationDir = absPath).
			if generationDir != absPath {
				_ = os.RemoveAll(tmpDir)
			}
		}()
	} else if options.Clean {
		if err := os.RemoveAll(absPath); err != nil {
			return fmt.Errorf("failed to clean output directory: %w", err)
		}
	}

	// Initialise file creator.
	var creator config.FileCreator
	if !options.DryRun {
		creator = config.DirCreator(generationDir)
	}

	// Start a spinner for progress messages.
	spinner := term.StartSpinner("Initialising...")

	// Initialise and run generator.
	stats, err := generator.NewGenerator(
		options,
		creator,
		spinner.Logger(),
	).Generate(patterns...)

	// Stop spinner and print summary.
	term.StopSpinner(spinner)
	term.Infof(
		"Processed: %s, %s, %s, %s, %s, %s in %s.",
		pluralise(stats.NumPackages, "Package"),
		pluralise(stats.NumServices, "Service"),
		pluralise(stats.NumMethods, "Method"),
		pluralise(stats.NumEnums, "Enum"),
		pluralise(stats.NumModels, "Model"),
		pluralise(stats.NumEvents, "Event"),
		stats.Elapsed().String(),
	)

	// Report output directory.
	term.Infof("Output directory: %s", absPath)

	// Process generator error.
	if err != nil {
		var report *generator.ErrorReport
		switch {
		case errors.Is(err, generator.ErrNoPackages):
			// Convert to warning message.
			term.Warning(err)
		case errors.As(err, &report):
			if report.HasErrors() {
				// Report error count.
				return err
			} else if report.HasWarnings() {
				// Report warning count.
				term.Warning(report)
			}
		default:
			// Report error.
			return err
		}
	}

	// Atomically replace the output directory with the temp dir.
	// Move the old directory aside BEFORE installing the new one so that if the
	// rename into place fails we can restore the original bindings. Deleting first
	// would leave no bindings at all on a rename failure.
	if generationDir != absPath {
		var oldDir string
		if _, statErr := os.Stat(absPath); statErr == nil {
			oldDir = absPath + ".old"
			if err := os.Rename(absPath, oldDir); err != nil {
				return fmt.Errorf("failed to move existing bindings aside: %w", err)
			}
		}
		if err := os.Rename(generationDir, absPath); err != nil {
			if oldDir != "" {
				_ = os.Rename(oldDir, absPath) // best-effort restore
			}
			return fmt.Errorf("failed to install new bindings: %w", err)
		}
		if oldDir != "" {
			_ = os.RemoveAll(oldDir)
		}
		// Signal the defer that cleanup is no longer needed.
		generationDir = absPath
	}

	return nil
}

func pluralise(number int, word string) string {
	if number == 1 {
		return fmt.Sprintf("%d %s", number, word)
	}
	return fmt.Sprintf("%d %ss", number, word)
}
