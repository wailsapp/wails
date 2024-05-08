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

	// Schedule shortcut generation.
	for _, info := range importMap.External {
		localInfo := info // Redeclare in loop-local scope
		generator.controller.Schedule(func() {
			generator.generateShortcut(localInfo)
		})
	}

	file, err := generator.creator.Create(filepath.Join(generator.renderer.IndexFile()))
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
func (generator *Generator) generateShortcut(info collect.ImportInfo) {
	file, err := generator.creator.Create(filepath.Join(generator.renderer.ShortcutFile(info)))
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("package %s: shortcut file generation failed", info.RelPath[2:])
		return
	}
	defer file.Close()

	err = generator.renderer.Shortcut(file, info)
	if err != nil {
		generator.controller.Errorf("%v", err)
		generator.controller.Errorf("package %s: shortcut file generation failed", info.RelPath[2:])
	}
}
