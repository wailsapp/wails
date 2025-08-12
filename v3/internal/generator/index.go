package generator

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// generateIndex generates an index file from the given index information.
func (generator *Generator) generateIndex(index *collect.PackageIndex) {
	defer generator.reportDualRoles(index)

	file, err := generator.creator.Create(filepath.Join(index.Package.Path, generator.renderer.IndexFile()))
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: index generation failed", index.Package.Path)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			generator.logger.Errorf("%v", err)
			generator.logger.Errorf("package %s: index generation failed", index.Package.Path)
		}
	}()

	err = generator.renderer.Index(file, index)
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package %s: index generation failed", index.Package.Path)
	}
}

// reportDualRoles checks for models types that are also service types
// and emits a warning.
func (generator *Generator) reportDualRoles(index *collect.PackageIndex) {
	services, models := index.Services, index.Models
	for len(services) > 0 && len(models) > 0 {
		if services[0].Name < models[0].Name {
			services = services[1:]
		} else if services[0].Name > models[0].Name {
			models = models[1:]
		} else {
			generator.logger.Warningf(
				"package %s: type %s has been marked both as a service and as a model; shadowing between the two may take place when importing generated JS indexes",
				index.Package.Path,
				services[0].Name,
			)

			services = services[1:]
			models = models[1:]
		}
	}
}
