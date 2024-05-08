package parser

import (
	"go/types"
	"path/filepath"
)

// generateBindings collects information
// and generates JS/TS bindings for the given type.
func (generator *Generator) generateBindings(typ *types.TypeName) {
	generator.controller.Debugf(
		"discovered bound type %s from package %s",
		typ.Name(),
		typ.Pkg().Path(),
	)

	success := false
	defer func() {
		if !success {
			generator.controller.Errorf(
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
		generator.controller.Infof(
			"package %s: type %s: bound type has no exported methods, skipping",
			typ.Pkg().Path(),
			typ.Name(),
		)
		success = true
		return
	}

	// Check for naming collisions.
	filename := generator.renderer.BindingsFile(info.Name)
	switch filename {
	case generator.renderer.IndexFile():
		if !generator.options.NoIndex {
			generator.controller.Errorf(
				"package %s: type %s: bound type filename collides with JS/TS index filename; please change the name of the type or choose a different filename for JS/TS indexes",
				typ.Pkg().Path(),
				typ.Name(),
			)
			return
		}

	case generator.renderer.ModelsFile():
		generator.controller.Errorf(
			"package %s: type %s: bound type filename collides with models filename; please change the name of the type or choose a different filename for models",
			typ.Pkg().Path(),
			typ.Name(),
		)
		return

	case generator.renderer.InternalFile():
		generator.controller.Errorf(
			"package %s: type %s: bound type filename collides with internal models filename; please change the name of the type or choose a different filename for internal models",
			typ.Pkg().Path(),
			typ.Name(),
		)
		return
	}

	// Create binding file.
	file, err := generator.creator.Create(filepath.Join(info.Imports.Self, filename))
	if err != nil {
		generator.controller.Errorf("%v", err)
		return
	}
	defer file.Close()

	// Render bound type.
	err = generator.renderer.Bindings(file, info)
	if err != nil {
		generator.controller.Errorf("%v", err)
		return
	}

	success = true
}
