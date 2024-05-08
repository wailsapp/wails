package parser

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// generateIncludes copies included files to the package directory
// for the package summarised by the given index.
func (generator *Generator) generateIncludes(index *collect.PackageIndex) {
	for name, info := range index.Package.Includes {
		if !(!generator.options.TS && info.JS) && !(generator.options.TS && info.TS) {
			// Include not enabled for current language.
			continue
		}

		// Validate filename.
		switch name {
		case generator.renderer.ModelsFile():
			if len(index.Models) > 0 {
				generator.logger.Errorf(
					"package %s: included file '%s' collides with models filename; please rename the file or choose a different filename for models",
					index.Package.Path,
					info.Path,
				)
				return
			}

		case generator.renderer.InternalFile():
			if len(index.Internal) > 0 {
				generator.logger.Errorf(
					"package %s: included file '%s' collides with internal models filename; please rename the file or choose a different filename for internal models",
					index.Package.Path,
					info.Path,
				)
				return
			}

		case generator.renderer.IndexFile():
			if !generator.options.NoIndex && !index.IsEmpty() {
				generator.logger.Errorf(
					"package %s: included file '%s' collides with JS/TS index filename; please rename the file or choose a different filename for JS/TS indexes",
					index.Package.Path,
					info.Path,
				)
				return
			}
		}

		// Validate against services.
		service, ok := slices.BinarySearchFunc(index.Services, name, func(service *collect.ServiceInfo, name string) int {
			return strings.Compare(generator.renderer.ServiceFile(service.Name), name)
		})
		if ok {
			generator.logger.Errorf(
				"package %s: included file '%s' collides with filename for service %s; please rename either the file or the service",
				index.Package.Path,
				info.Path,
				index.Services[service].Name,
			)
			return
		}

		func() { // Scoped defer pattern.
			src, err := os.Open(info.Path)
			if err != nil {
				generator.logger.Errorf("%v", err)
				generator.logger.Errorf("package %s: could not read included file '%s'", index.Package.Path, info.Path)
				return
			}
			defer src.Close()

			stat, err := src.Stat()
			if err != nil {
				generator.logger.Errorf("%v", err)
				generator.logger.Errorf("package %s: could not read included file '%s'", index.Package.Path, info.Path)
				return
			}

			if stat.IsDir() {
				generator.logger.Errorf(
					"package %s: included file '%s' is a directory; please glob or list all descendants explicitly",
					index.Package.Path,
					info.Path,
				)
				return
			}

			dst, err := generator.creator.Create(filepath.Join(index.Package.Path, name))
			if err != nil {
				generator.logger.Errorf("%v", err)
				generator.logger.Errorf("package %s: could not write included file '%s'", index.Package.Path, name)
				return
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				generator.logger.Errorf("%v", err)
				generator.logger.Errorf("package %s: could not copy included file '%s'", index.Package.Path, name)
			}
		}()
	}
}
