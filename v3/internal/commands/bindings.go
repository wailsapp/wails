package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser"
)

func GenerateBindings(options *flags.GenerateBindingsOptions) error {

	if options.Silent {
		pterm.DisableOutput()
		defer pterm.EnableOutput()
	}

	project, err := parser.GenerateBindingsAndModels(options)
	if err != nil {
		if errors.Is(err, parser.ErrBadApplicationOptions) {
			pterm.Warning.Println(err.Error())
			return nil
		} else {
			return err
		}
	}

	absPath, err := filepath.Abs(options.OutputDirectory)
	if err != nil {
		return err
	}

	pterm.Info.Printf("Processed: %s, %s, %s, %s, %s, %s in %s.\n",
		pluralise(project.Stats.NumPackages, "Package"),
		pluralise(project.Stats.NumServices, "Service"),
		pluralise(project.Stats.NumMethods, "Method"),
		pluralise(project.Stats.NumEnums, "Enum"),
		pluralise(project.Stats.NumModels, "Model"),
		pluralise(project.Stats.NumAliases, "Alias", true),
		project.Stats.EndTime.Sub(project.Stats.StartTime).String())

	pterm.Info.Printf("Output directory: %s\n", absPath)

	return nil
}

func pluralise(number int, word string, es ...bool) string {
	if number == 1 {
		return fmt.Sprintf("%d %s", number, word)
	}
	if len(es) > 0 && es[0] {
		return fmt.Sprintf("%d %ses", number, word)
	}
	return fmt.Sprintf("%d %ss", number, word)
}
