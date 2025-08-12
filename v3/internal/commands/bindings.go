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

	// Clean the output directory if clean option is enabled
	if options.Clean {
		if err := os.RemoveAll(absPath); err != nil {
			return fmt.Errorf("failed to clean output directory: %w", err)
		}
	}

	// Initialise file creator.
	var creator config.FileCreator
	if !options.DryRun {
		creator = config.DirCreator(absPath)
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
		"Processed: %s, %s, %s, %s, %s in %s.",
		pluralise(stats.NumPackages, "Package"),
		pluralise(stats.NumServices, "Service"),
		pluralise(stats.NumMethods, "Method"),
		pluralise(stats.NumEnums, "Enum"),
		pluralise(stats.NumModels, "Model"),
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

	return nil
}

func pluralise(number int, word string) string {
	if number == 1 {
		return fmt.Sprintf("%d %s", number, word)
	}
	return fmt.Sprintf("%d %ss", number, word)
}
