package parser

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// generateModels generates a JS/TS model file for the given list of models.
//
// If internal is true, the generated file is named "internal"
// and the types declared therein are not exported by the package index file.
//
// A call to index.Info.Collect must complete before entering generateModels.
func (generator *Generator) generateModels(info *collect.PackageInfo, models []*collect.ModelInfo, internal bool) {
	defer generator.wg.Done()

	// Merge all import maps.
	imports := collect.NewImportMap(info)
	for _, model := range models {
		imports.Merge(model.Imports)
	}

	if internal {
		clear(imports.Internal)
	} else {
		clear(imports.Models)
	}

	var filename string
	if internal {
		filename = generator.renderer.InternalFile()
	} else {
		filename = generator.renderer.ModelsFile()
	}

	file, err := generator.creator.Create(filepath.Join(info.Path, filename))
	if err != nil {
		pterm.Error.Println(err)

		var prefix string
		if internal {
			prefix = "internal "
		}
		pterm.Error.Printfln("package %s: %smodels generation failed", info.Path, prefix)
		return
	}
	defer file.Close()

	err = generator.renderer.Models(file, imports, models)
	if err != nil {
		pterm.Error.Println(err)

		var prefix string
		if internal {
			prefix = "internal "
		}
		pterm.Error.Printfln("package %s: %smodels generation failed", info.Path, prefix)
	}
}
