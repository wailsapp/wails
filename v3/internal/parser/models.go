package parser

import (
	"bytes"
	"io"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

type ModelDefinitions struct {
	Package string
	Imports map[string]string

	Models map[string]*StructDef
	Enums  map[string]*TypeDef

	ModelsFilename string
}

func (p *Project) GenerateModel(wr io.Writer, def *ModelDefinitions, options *flags.GenerateBindingsOptions) error {
	template := templates.ModelsJS
	if options.TS {
		if options.UseInterfaces {
			template = templates.InterfacesTS
		} else {
			template = templates.ModelsTS
		}
	}

	// Fix up TS names
	for _, model := range def.Models {
		model.Name = options.TSPrefix + model.Name + options.TSSuffix
	}

	if err := template.Execute(wr, def); err != nil {
		println("Problem executing template: " + err.Error())
		return err
	}

	return nil
}

type Model struct {
	Package string
}

func (p *Project) GenerateModels(models map[packagePath]map[structName]*StructDef, enums map[packagePath]map[string]*TypeDef, options *flags.GenerateBindingsOptions) (result map[string]string, err error) {
	if len(models) == 0 && len(enums) == 0 {
		return
	}

	result = make(map[string]string)

	// merge package lists by alias (e.g. services) instead of full pkg name (e.g. github.com/wailsapp/wails/somedir/services)

	var keys = lo.Keys(models)
	keys = append(keys, lo.Keys(enums)...)

	slices.SortFunc(keys, func(key1, key2 string) int {
		return strings.Compare(pkgAlias(key1), pkgAlias(key2))
	})
	keys = slices.CompactFunc(keys, func(key1, key2 string) bool {
		return pkgAlias(key1) == pkgAlias(key2)
	})

	for _, pkg := range keys {
		pkgModels, pkgEnums := models[pkg], enums[pkg]

		var buffer bytes.Buffer
		err = p.GenerateModel(&buffer, &ModelDefinitions{
			Package: pkg,
			Imports: p.calculateModelImports(pkg, pkgModels),

			Models: pkgModels,
			Enums:  pkgEnums,

			ModelsFilename: options.ModelsFilename,
		}, options)

		if err != nil {
			return
		}

		// Get the relative package path
		relativePackageDir := p.RelativePackageDir(pkg)
		result[relativePackageDir] = buffer.String()
	}

	return
}

func (p *Project) calculateModelImports(pkg string, m map[structName]*StructDef) map[string]string {
	result := make(map[string]string)
	pkgInfo := p.packageCache[pkg]

	for _, structDef := range m {
		for _, field := range structDef.Fields {
			if field.Type.Package != pkg {
				fieldPkgInfo := p.packageCache[field.Type.Package]
				// Find the relative path from the source directory to the target directory
				result[fieldPkgInfo.Name] = p.RelativeBindingsDir(pkgInfo, fieldPkgInfo)
			}
		}
	}

	return result
}
