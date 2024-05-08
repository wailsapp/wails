package parser

import (
	"go/types"
	"path/filepath"

	"github.com/pterm/pterm"
)

// generateBindings collects information
// and generates JS/TS bindings for the given type.
func (generator *Generator) generateBindings(typ *types.TypeName) {
	defer generator.wg.Done()

	success := false
	defer func() {
		if !success {
			pterm.Error.Printfln(
				"package %s: type %s: bindings generation failed",
				typ.Pkg().Path(),
				typ.Name(),
			)
		}
	}()

	// Collect bound type information.
	info := generator.collector.BoundType(typ)
	if info == nil {
		return
	}

	if len(info.Methods) == 0 {
		pterm.Info.Printfln(
			"package %s: type %s: bound type has no exported methods, skipping",
			typ.Pkg().Path(),
			typ.Name(),
		)
		success = true
		return
	}

	// Create binding file.
	file, err := generator.creator.Create(filepath.Join(info.Imports.Self, generator.renderer.BindingsFile(info.Name)))
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	defer file.Close()

	// Render bound type.
	err = generator.renderer.Bindings(file, info)
	if err != nil {
		pterm.Error.Println(err)
		return
	}

	success = true
}
