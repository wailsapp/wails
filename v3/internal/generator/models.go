package generator

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// generateModels generates a JS/TS model file for the given list of models.
//
// If internal is true, the generated file is named by Renderer.InternalFile
// and the types declared therein are not exported by the package index file.
//
// A call to info.Collect must complete before entering generateModels.
func (generator *Generator) generateModels(info *collect.PackageInfo, models []*collect.ModelInfo, internal bool) {
	// Merge all import maps.
	imports := collect.NewImportMap(info)
	for _, model := range models {
		imports.Merge(model.Imports)
	}

	// Clear irrelevant imports.
	if internal {
		imports.ImportInternal = false
	} else {
		imports.ImportModels = false
	}

	// Retrieve appropriate file name.
	var filename string
	if internal {
		filename = generator.renderer.InternalFile()
	} else {
		filename = generator.renderer.ModelsFile()
	}

	file, err := generator.creator.Create(filepath.Join(info.Path, filename))
	if err != nil {
		generator.logger.Errorf("%v", err)

		var prefix string
		if internal {
			prefix = "internal "
		}
		generator.logger.Errorf("package %s: %smodels generation failed", info.Path, prefix)
		return
	}
	defer file.Close()

	err = generator.renderer.Models(file, imports, models)
	if err != nil {
		generator.logger.Errorf("%v", err)

		var prefix string
		if internal {
			prefix = "internal "
		}
		generator.logger.Errorf("package %s: %smodels generation failed", info.Path, prefix)
	}
}
