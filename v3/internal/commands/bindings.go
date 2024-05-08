package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser"
	"github.com/wailsapp/wails/v3/internal/parser/config"
)

func GenerateBindings(options *flags.GenerateBindingsOptions, patterns []string) error {
	if options.Silent {
		pterm.DisableOutput()
		defer pterm.EnableOutput()
	} else if options.Verbose {
		pterm.EnableDebugMessages()
		defer pterm.DisableDebugMessages()
	}

	if len(patterns) == 0 {
		// No input pattern, load package from current directory.
		patterns = []string{"."}
	}

	// Compute absolute path of output directory.
	absPath, err := filepath.Abs(options.OutputDirectory)
	if err != nil {
		return err
	}

	// Initialise file creator.
	var creator config.FileCreator
	if !options.DryRun {
		creator = config.DirCreator(absPath)
	}

	// Start a spinner for progress messages.
	spinner, _ := pterm.DefaultSpinner.Start("Initialising...")

	// Initialise and run generator.
	stats, err := parser.NewGenerator(
		options,
		creator,
		config.DefaultPtermLogger(spinner),
	).Generate(patterns...)

	// Resolve spinner.
	spinner.Info(fmt.Sprintf(
		"Processed: %s, %s, %s, %s, %s in %s.",
		pluralise(stats.NumPackages, "Package"),
		pluralise(stats.NumServices, "Service"),
		pluralise(stats.NumMethods, "Method"),
		pluralise(stats.NumEnums, "Enum"),
		pluralise(stats.NumModels, "Model"),
		stats.Elapsed().String(),
	))

	// Report output directory.
	pterm.Info.Printfln("Output directory: %s", absPath)

	// Process generator error.
	if err != nil {
		var report *parser.ErrorReport
		switch {
		case errors.Is(err, parser.ErrNoPackages):
			// Convert to warning message.
			pterm.Warning.Println(err)
		case errors.As(err, &report):
			if report.HasErrors() {
				// Report error count.
				return err
			} else if report.HasWarnings() {
				// Report warning count.
				pterm.Warning.Println(report)
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
