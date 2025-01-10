package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pterm/pterm"

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
		pterm.EnableDebugMessages()
		defer pterm.DisableDebugMessages()
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

	// Initialise file creator.
	var creator config.FileCreator
	if !options.DryRun {
		creator = config.DirCreator(absPath)
	}

	// Start a spinner for progress messages.
	var spinner *pterm.SpinnerPrinter
	if term.IsTerminal() {
		spinner, _ = pterm.DefaultSpinner.Start("Initialising...")
	}

	// Initialise and run generator.
	stats, err := generator.NewGenerator(
		options,
		creator,
		config.DefaultPtermLogger(spinner),
	).Generate(patterns...)

	// Resolve spinner.
	resultMessage := fmt.Sprintf(
		"Processed: %s, %s, %s, %s, %s in %s.",
		pluralise(stats.NumPackages, "Package"),
		pluralise(stats.NumServices, "Service"),
		pluralise(stats.NumMethods, "Method"),
		pluralise(stats.NumEnums, "Enum"),
		pluralise(stats.NumModels, "Model"),
		stats.Elapsed().String(),
	)
	if spinner != nil {
		spinner.Info(resultMessage)
	} else {
		term.Infofln(resultMessage)
	}

	// Report output directory.
	term.Infofln("Output directory: %s", absPath)

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
