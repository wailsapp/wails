package parser

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// generateIndex generates an index file from the given index information.
// A call to index.Info.Collect must complete before entering generateIndex.
func (generator *Generator) generateIndex(index collect.PackageIndex) {
	file, err := generator.creator.Create(filepath.Join(index.Info.Path, generator.renderer.IndexFile()))
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("package %s: index generation failed", index.Info.Path)
		return
	}
	defer file.Close()

	err = generator.renderer.Index(file, &index)
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("package %s: index generation failed", index.Info.Path)
	}
}

// generateGlobalIndex generates a global index file from the given import map,
// as well as shortcut files for each package named by the input patterns.
func (generator *Generator) generateGlobalIndex(imports []*collect.PackageInfo) {
	// Sort imported packages by path
	// to ensure name collisions are resolved deterministically.
	slices.SortFunc(imports, func(p1 *collect.PackageInfo, p2 *collect.PackageInfo) int {
		return strings.Compare(p1.Path, p2.Path)
	})

	// Initialise import map with fake root package.
	importMap := collect.NewImportMap(generator.collector.Package(""))
	for _, pkg := range imports {
		importMap.Add(pkg)
	}

	// Schedule shortcut generation and detect name collisions.
	for _, info := range importMap.External {
		if info.Index == 0 {
			name := info.Name // Declare in loop-local scope
			generator.controller.Schedule(func() {
				generator.generateShortcut(importMap, name)
			})
		} else if info.Index == 1 {
			generator.controller.Warningf(
				"package name '%s' appears more than once in global index; shadowing may take place between identically named exports",
				info.Name,
			)
		}
	}

	file, err := generator.creator.Create(generator.renderer.IndexFile())
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("global index generation failed")
		return
	}
	defer file.Close()

	err = generator.renderer.GlobalIndex(file, importMap)
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("global index generation failed")
	}
}

// generateShortcut generates a shortcut file for the given import.
func (generator *Generator) generateShortcut(imports *collect.ImportMap, name string) {
	filename := generator.renderer.ShortcutFile(name)
	if filename == generator.renderer.IndexFile() {
		generator.controller.Errorf(
			"package name '%s': shortcut filename collides with JS/TS index filename; please rename all such packages or choose a different filename for JS/TS indexes",
			name,
		)

		for _, info := range imports.External {
			if info.Name == name {
				generator.controller.Errorf(
					"package %s: invalid package name '%s'",
					info.Name,
					name,
				)
			}
		}
	}

	file, err := generator.creator.Create(filename)
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("package name '%s': shortcut file generation failed", name)
		return
	}
	defer file.Close()

	err = generator.renderer.Shortcut(file, imports, name)
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("package name '%s': shortcut file generation failed", name)
	}
}
