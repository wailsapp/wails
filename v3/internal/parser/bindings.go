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

type BindingDefinitions struct {
	Package      string
	Imports      map[string]string
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

func (p *Project) calculateBindingImports(pkg string, methods []*BoundMethod) map[string]string {
	result := make(map[string]string)
	pkgInfo := p.packageCache[pkg]

	for _, method := range methods {
		for _, param := range method.JSInputs() {
			if param.Type.Package != pkg {
				paramPkgInfo := p.packageCache[param.Type.Package]
				// Find the relative path from the source directory to the target directory
				result[paramPkgInfo.Name] = p.RelativeBindingsDir(pkgInfo, paramPkgInfo)
			}
		}

		for _, param := range method.JSOutputs() {
			if param.Type.Package != pkg {
				paramPkgInfo := p.packageCache[param.Type.Package]
				// Find the relative path from the source directory to the target directory
				result[paramPkgInfo.Name] = p.RelativeBindingsDir(pkgInfo, paramPkgInfo)
			}
		}
	}

	return result
}

func (p *Project) calculateBindingLocalImports(pkg string, methods []*BoundMethod) []structName {
	requiredTypes := make(map[structName]bool)

	for _, method := range methods {
		for _, param := range method.JSInputs() {
			if param.Type.Package == pkg && (param.Type.IsStruct || param.Type.IsEnum) {
				requiredTypes[param.Type.Name] = true
			}
		}

		for _, param := range method.JSOutputs() {
			if param.Type.Package == pkg && (param.Type.IsStruct || param.Type.IsEnum) {
				requiredTypes[param.Type.Name] = true
			}
		}
	}

	result := lo.Keys(requiredTypes)
	slices.Sort(result)

	return result
}
