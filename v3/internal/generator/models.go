package generator

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// generateModels generates a JS/TS models file for the given list of models.
// A call to info.Collect must complete before entering generateModels.
func (generator *Generator) generateModels(info *collect.PackageInfo, models []*collect.ModelInfo) {
	// Merge all import maps.
	imports := collect.NewImportMap(info)
	for _, model := range models {
		imports.Merge(model.Imports)
	}

	// Clear irrelevant imports.
	imports.ImportModels = false

	file, err := generator.creator.Create(filepath.Join(info.Path, generator.renderer.ModelsFile()))
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: models generation failed", info.Path)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			generator.logger.Errorf("%v", err)
			generator.logger.Errorf("package %s: models generation failed", info.Path)
		}
	}()

	err = generator.renderer.Models(file, imports, models)
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: models generation failed", info.Path)
	}
}
