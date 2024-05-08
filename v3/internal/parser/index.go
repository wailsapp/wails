package parser

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
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
	defer file.Close()

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

// generateGlobalIndex generates a shortcut file for each package
// in the given slice, then a global index file that exports all shortcuts.
func (generator *Generator) generateGlobalIndex(imports []*collect.PackageInfo) {
	// Sort imported packages by path
	// to ensure name collisions are resolved deterministically.
	slices.SortFunc(imports, func(p1 *collect.PackageInfo, p2 *collect.PackageInfo) int {
		return strings.Compare(p1.Path, p2.Path)
	})

	// Initialise import map for root package.
	importMap := collect.NewImportMap(nil)
	for _, pkg := range imports {
		importMap.Add(pkg)
	}

	// Schedule shortcut generation and detect name collisions.
	for _, info := range importMap.External {
		if info.Index == 0 {
			name := info.Name // Declare in loop-local scope
			generator.scheduler.Schedule(func() {
				generator.generateShortcut(importMap, name)
			})
		} else if info.Index == 1 {
			generator.logger.Warningf(
				"package name '%s' appears more than once in global index; shadowing may take place between identically named exports",
				info.Name,
			)
		}
	}

	file, err := generator.creator.Create(generator.renderer.IndexFile())
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("global index generation failed")
		return
	}
	defer file.Close()

	err = generator.renderer.GlobalIndex(file, importMap)
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("global index generation failed")
	}
}

// generateShortcut generates a shortcut file for the given import.
func (generator *Generator) generateShortcut(imports *collect.ImportMap, name string) {
	filename := generator.renderer.ShortcutFile(name)

	// Validate filename.
	if filename == generator.renderer.IndexFile() {
		generator.logger.Errorf(
			"package name '%s': shortcut filename collides with JS/TS index filename; please rename all such packages or choose a different filename for JS/TS indexes",
			name,
		)

		for _, info := range imports.External {
			if info.Name == name {
				generator.logger.Errorf(
					"package %s: invalid package name '%s'",
					info.Name,
					name,
				)
			}
		}

		return
	}

	file, err := generator.creator.Create(filename)
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package name '%s': shortcut file generation failed", name)
		return
	}
	defer file.Close()

	err = generator.renderer.Shortcut(file, imports, name)
	if err != nil {
		generator.logger.Errorf("%v", err)
		generator.logger.Errorf("package name '%s': shortcut file generation failed", name)
	}
}
