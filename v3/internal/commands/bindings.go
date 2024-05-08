package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser"
	"github.com/wailsapp/wails/v3/internal/parser/analyse"
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

	generator := parser.NewGenerator(options, nil)
	stats, err := generator.Generate(patterns...)
	if err != nil {
		switch {
		case errors.Is(err, parser.ErrNoInitialPackages):
			pterm.Info.Println(err)
		case errors.Is(err, analyse.ErrNoApplicationPackage):
			pterm.Info.Println("Input packages do not load the Wails application package")
		default:
			return err
		}
	}

	pterm.Info.Printf("Processed: %s, %s, %s, %s, %s in %s.\n",
		pluralise(stats.NumPackages, "Package"),
		pluralise(stats.NumTypes, "Bound Type"),
		pluralise(stats.NumMethods, "Method"),
		pluralise(stats.NumEnums, "Enum"),
		pluralise(stats.NumModels, "Model"),
		stats.Elapsed().String())

	absPath, err := filepath.Abs(options.OutputDirectory)
	if err != nil {
		return err
	}

	pterm.Info.Printf("Output directory: %s\n", absPath)

	return nil
}

func pluralise(number int, word string) string {
	if number == 1 {
		return fmt.Sprintf("%d %s", number, word)
	}
	return fmt.Sprintf("%d %ss", number, word)
}
