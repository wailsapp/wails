package generator

import (
	"go/types"
	"path/filepath"
)

// generateService collects information
// and generates JS/TS binding code
// for the given service type object.
func (generator *Generator) generateService(obj *types.TypeName) {
	generator.logger.Debugf(
		"discovered service type %s from package %s",
		obj.Name(),
		obj.Pkg().Path(),
	)

	success := false
	defer func() {
		if !success {
			generator.logger.Errorf(
				"package %s: type %s: service code generation failed",
				obj.Pkg().Path(),
				obj.Name(),
			)
		}
	}()

	// Collect service information.
	info := generator.collector.Service(obj).Collect()
	if info == nil {
		return
	}

	if info.IsEmpty() {
		if !info.HasInternalMethods {
			generator.logger.Infof(
				"package %s: type %s: service has no valid exported methods, skipping",
				obj.Pkg().Path(),
				obj.Name(),
			)
		}
		success = true
		return
	}

	// Check for standard filename collisions.
	filename := generator.renderer.ServiceFile(info.Name)
	switch filename {
	case generator.renderer.ModelsFile():
		generator.logger.Errorf(
			"package %s: type %s: service filename collides with models filename; please rename the type or choose a different filename for models",
			obj.Pkg().Path(),
			obj.Name(),
		)
		return

	case generator.renderer.IndexFile():
		if !generator.options.NoIndex {
			generator.logger.Errorf(
				"package %s: type %s: service filename collides with JS/TS index filename; please rename the type or choose a different filename for JS/TS indexes",
				obj.Pkg().Path(),
				obj.Name(),
			)
			return
		}
	}

	// Check for upper/lower-case filename collisions.
	path := filepath.Join(info.Imports.Self, filename)
	if other, present := generator.serviceFiles.LoadOrStore(path, obj); present {
		generator.logger.Errorf(
			"package %s: type %s: service filename collides with filename for service %s; please avoid multiple services whose names differ only in case",
			obj.Pkg().Path(),
			obj.Name(),
			other.(*types.TypeName).Name(),
		)
		return
	}

	// Create service file.
	file, err := generator.creator.Create(path)
	if err != nil {
		generator.logger.Errorf("%v", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			generator.logger.Errorf("%v", err)
			success = false
		}
	}()

	// Render service code.
	err = generator.renderer.Service(file, info)
	if err != nil {
		generator.logger.Errorf("%v", err)
		return
	}

	success = true
}
