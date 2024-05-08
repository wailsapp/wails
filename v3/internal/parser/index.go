package parser

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// generateIndex generates an index file from the given index information.
// A call to index.Info.Collect must complete before entering generateIndex.
func (generator *Generator) generateIndex(index collect.PackageIndex) {
	defer generator.wg.Done()

	file, err := generator.creator.Create(filepath.Join(index.Info.Path, generator.renderer.IndexFile()))
	if err != nil {
		pterm.Error.Println(err)
		pterm.Error.Printfln("package %s: index generation failed", index.Info.Path)
		return
	}
	defer file.Close()

	err = generator.renderer.Index(file, &index)
	if err != nil {
		pterm.Error.Println(err)
		pterm.Error.Printfln("package %s: index generation failed", index.Info.Path)
	}
}

// generateGlobalIndex generates a global index file from the given import map.
func (generator *Generator) generateGlobalIndex(imports *collect.ImportMap) {
	defer generator.wg.Done()

	file, err := generator.creator.Create(filepath.Join(generator.renderer.IndexFile()))
	if err != nil {
		pterm.Error.Println(err)
		pterm.Error.Printfln("global index generation failed")
		return
	}
	defer file.Close()

	err = generator.renderer.GlobalIndex(file, imports)
	if err != nil {
		pterm.Error.Println(err)
		pterm.Error.Printfln("global index generation failed")
	}
}

// generateShortcut generates a shortcut file for the given import.
func (generator *Generator) generateShortcut(info collect.ImportInfo) {
	defer generator.wg.Done()

	file, err := generator.creator.Create(filepath.Join(generator.renderer.ShortcutFile(info)))
	if err != nil {
		pterm.Error.Println(err)
		pterm.Error.Printfln("package %s: shortcut file generation failed", info.RelPath[2:])
		return
	}
	defer file.Close()

	err = generator.renderer.Shortcut(file, info)
	if err != nil {
		pterm.Error.Println(err)
		pterm.Error.Printfln("package %s: shortcut file generation failed", info.RelPath[2:])
	}
}
