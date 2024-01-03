package commands

import (
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser"
	"path/filepath"
)

func GenerateBindings(options *flags.GenerateBindingsOptions) error {

	if options.Silent {
		pterm.DisableOutput()
		defer pterm.EnableOutput()
	}

	project, err := parser.GenerateBindingsAndModels(options)
	if err != nil {
		if errors.Is(err, parser.ErrNoBindingsFound) {
			pterm.Info.Println("No bindings found")
			return nil
		} else {
			return err
		}
	}

	absPath, err := filepath.Abs(options.OutputDirectory)
	if err != nil {
		return err
	}

	pterm.Info.Printf("Processed: %s, %s, %s, %s, %s in %s.\n",
		pluralise(project.Stats.NumPackages, "Package"),
		pluralise(project.Stats.NumStructs, "Struct"),
		pluralise(project.Stats.NumMethods, "Method"),
		pluralise(project.Stats.NumEnums, "Enum"),
		pluralise(project.Stats.NumModels, "Model"),
		project.Stats.EndTime.Sub(project.Stats.StartTime).String())

	pterm.Info.Printf("Output directory: %s\n", absPath)

	return nil
}

func pluralise(number int, word string) string {
	if number == 1 {
		return fmt.Sprintf("%d %s", number, word)
	}
	return fmt.Sprintf("%d %ss", number, word)
}
