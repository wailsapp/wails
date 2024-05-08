package parser

import (
	"bytes"
	"io"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

type BindingDefinitions struct {
	Package      string
	Imports      []*ImportDef
	LocalImports []structName

	Struct  string
	Methods []*BoundMethod

	ModelsFilename    string
	UseBundledRuntime bool
	UseIDs            bool
}

func (p *Project) GenerateBinding(wr io.Writer, def *BindingDefinitions, options *flags.GenerateBindingsOptions) error {
	template := templates.BindingsJS
	if options.TS {
		template = templates.BindingsTS
	}

	if err := template.Execute(wr, def); err != nil {
		println("Problem executing template: " + err.Error())
		return err
	}

	return nil
}

func (p *Project) GenerateBindings(bindings map[packagePath]map[structName][]*BoundMethod, options *flags.GenerateBindingsOptions) (result map[string]map[string]string, err error) {
	result = make(map[string]map[string]string)

	for pkg, structs := range bindings {
		var pkgBindings = make(map[string]string)

		for structName, methods := range structs {
			slices.SortFunc(methods, func(m1, m2 *BoundMethod) int {
				return strings.Compare(m1.Name, m2.Name)
			})

			var buffer bytes.Buffer
			err = p.GenerateBinding(&buffer, &BindingDefinitions{
				Package:      pkg,
				Imports:      p.calculateBindingImports(pkg, methods),
				LocalImports: p.calculateBindingLocalImports(pkg, methods),

				Struct:  pkgAlias(pkg) + "." + structName,
				Methods: methods,

				ModelsFilename:    options.ModelsFilename,
				UseBundledRuntime: options.UseBundledRuntime,
				UseIDs:            options.UseIDs,
			}, options)

			if err != nil {
				return
			}

			pkgBindings[structName] = buffer.String()
		}

		// Get the relative package path
		relativePackageDir := p.RelativePackageDir(pkg)
		result[relativePackageDir] = pkgBindings
	}

	return
}

func (p *Project) calculateBindingImports(pkg string, methods []*BoundMethod) []*ImportDef {
	var result []*ImportDef
	var seen = make(map[string]bool)

	pkgInfo := p.packageCache[pkg]

	var processParameter = func(param *Parameter) {
		if param.Type.Package != pkg {
			// Find the relative path from the source directory to the target directory
			paramPkgInfo := p.packageCache[param.Type.Package]
			relativePath := p.RelativeBindingsDir(pkgInfo, paramPkgInfo)

			// Deduplicate imports
			if _, ok := seen[relativePath]; ok {
				return
			}
			seen[relativePath] = true

			result = append(result, &ImportDef{
				PackageName: paramPkgInfo.Name,
				Path:        relativePath,
			})
		}
	}

	for _, method := range methods {
		for _, param := range method.JSInputs() {
			processParameter(param)
		}
		for _, param := range method.JSOutputs() {
			processParameter(param)
		}
	}

	return result
}

func (p *Project) calculateBindingLocalImports(pkg string, methods []*BoundMethod) []structName {
	var result []structName
	var seen = make(map[string]bool)

	var processParameter = func(param *Parameter) {
		if param.Type.Package == pkg && (param.Type.IsStruct || param.Type.IsEnum) {
			// Deduplicate imports
			if _, ok := seen[param.Type.Name]; ok {
				return
			}
			seen[param.Type.Name] = true

			result = append(result, param.Type.Name)
		}
	}

	for _, method := range methods {
		for _, param := range method.JSInputs() {
			processParameter(param)
		}
		for _, param := range method.JSOutputs() {
			processParameter(param)
		}
	}

	return result
}
