package parser

import (
	"go/types"
	"path/filepath"
)

// generateService collects information
// and generates JS/TS binding code
// for the given service type object.
func (generator *Generator) generateService(typ *types.TypeName) {
	generator.logger.Debugf(
		"discovered service type %s from package %s",
		typ.Name(),
		typ.Pkg().Path(),
	)

	success := false
	defer func() {
		if !success {
			generator.logger.Errorf(
				"package %s: type %s: service code generation failed",
				typ.Pkg().Path(),
				typ.Name(),
			)
		}
	}()

	// Collect service information.
	info := generator.collector.Service(typ).Collect()
	if info == nil {
		return
	}

	if info.IsEmpty(generator.options.TS) {
		generator.logger.Infof(
			"package %s: type %s: service has no valid exported methods, skipping",
			typ.Pkg().Path(),
			typ.Name(),
		)
		success = true
		return
	}

	// Check for file name collisions.
	filename := generator.renderer.ServiceFile(info.Name)
	switch filename {
	case generator.renderer.ModelsFile():
		generator.logger.Errorf(
			"package %s: type %s: service filename collides with models filename; please rename the type or choose a different filename for models",
			typ.Pkg().Path(),
			typ.Name(),
		)
		return

	case generator.renderer.InternalFile():
		generator.logger.Errorf(
			"package %s: type %s: service filename collides with internal models filename; please rename the type or choose a different filename for internal models",
			typ.Pkg().Path(),
			typ.Name(),
		)
		return

	case generator.renderer.IndexFile():
		if !generator.options.NoIndex {
			generator.logger.Errorf(
				"package %s: type %s: service filename collides with JS/TS index filename; please rename the type or choose a different filename for JS/TS indexes",
				typ.Pkg().Path(),
				typ.Name(),
			)
			return
		}
	}

	// Create service file.
	file, err := generator.creator.Create(filepath.Join(info.Imports.Self, filename))
	if err != nil {
		generator.logger.Errorf("%v", err)
		return
	}
	defer file.Close()

	// Render service code.
	err = generator.renderer.Service(file, info)
	if err != nil {
		generator.logger.Errorf("%v", err)
		return
	}

	success = true
}
