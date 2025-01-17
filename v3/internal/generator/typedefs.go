package generator

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// generateTypedefs generates a JS/TS typedef file for the given list of models.
// A call to info.Collect must complete before entering generateTypedefs.
func (generator *Generator) generateTypedefs(info *collect.PackageInfo, models []*collect.ModelInfo) {
	// Merge all import maps.
	imports := collect.NewImportMap(info)
	for _, model := range models {
		imports.Merge(model.Imports)
	}

	// Clear irrelevant imports.
	imports.ImportModels = false

	file, err := generator.creator.Create(filepath.Join(info.Path, generator.renderer.InternalFile()))
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: models generation failed", info.Path)
		return
	}
	defer file.Close()

	err = generator.renderer.Typedefs(file, imports, models)
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: models generation failed", info.Path)
	}
}
