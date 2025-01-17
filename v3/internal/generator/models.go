package generator

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// generateModels generates a file for exported models from the given index information.
func (generator *Generator) generateModels(index *collect.PackageIndex) {
	file, err := generator.creator.Create(filepath.Join(index.Package.Path, generator.renderer.ModelsFile()))
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: exported models generation failed", index.Package.Path)
		return
	}
	defer file.Close()

	err = generator.renderer.Models(file, index)
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: exported models generation failed", index.Package.Path)
	}
}
