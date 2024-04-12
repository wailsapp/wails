package parser

import (
	"bytes"
	"go/types"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

type BindingDefinitions struct {
	Package      *Package
	Service      *Service
	Imports      map[string]string
	LocalImports []string

	Struct  string
	Methods []*BoundMethod

	ModelsFilename    string
	UseBundledRuntime bool
	UseNames          bool
}

func generateBinding(wr io.Writer, def *BindingDefinitions, options *flags.GenerateBindingsOptions) error {
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

func (p *Project) GenerateBindings() (result map[string]map[string]string, err error) {
	result = make(map[string]map[string]string)

	for _, pkg := range p.pkgs {
		bindings, err := pkg.GenerateBindings(p)
		if err != nil {
			return nil, err
		}
		packageDir := p.PackageDir(pkg.Types)
		result[packageDir] = bindings
	}
	return
}

func (p *Package) GenerateBindings(project *Project) (result map[string]string, err error) {
	result = make(map[string]string)
	options := project.options

	for _, service := range p.services {
		methods := service.Methods
		slices.SortFunc(methods, func(m1, m2 *BoundMethod) int {
			return strings.Compare(m1.Name(), m2.Name())
		})

		var buffer bytes.Buffer
		err = generateBinding(&buffer, &BindingDefinitions{
			Package:      p,
			Service:      service,
			Imports:      service.calculateBindingImports(p, project),
			LocalImports: service.calculateBindingLocalImports(p),

			Methods: methods,

			ModelsFilename:    options.ModelsFilename,
			UseBundledRuntime: options.UseBundledRuntime,
			UseNames:          options.UseNames,
		}, options)

		if err != nil {
			return
		}

		result[service.Name()] = buffer.String()
	}
	return
}

func (s *Service) bindingImportsOf(params []*Parameter, pkg *Package, project *Project) map[string]string {
	result := make(map[string]string)

	for _, param := range params {
		models := param.Models(pkg, false)
		for model := range models {
			if model.Obj() != nil && model.Obj().Pkg() != s.Pkg() {
				otherPkg := model.Obj().Pkg()
				result[otherPkg.Name()] = project.RelativePackageDir(s.Pkg(), otherPkg)
			}
		}
	}
	return result
}

func (s *Service) calculateBindingImports(pkg *Package, project *Project) map[string]string {
	result := make(map[string]string)

	for _, method := range s.Methods {
		maps.Copy(result, s.bindingImportsOf(method.JSInputs(), pkg, project))
		maps.Copy(result, s.bindingImportsOf(method.JSOutputs(), pkg, project))
	}

	return result
}

func (s *Service) bindingLocalImportsOf(params []*Parameter, pkg *Package) map[string]bool {
	requiredTypes := make(map[string]bool)

	for _, param := range params {
		models := param.Models(pkg, false)
		for model := range models {
			if structType, ok := model.Underlying().(*types.Struct); ok && model.Obj() == nil {
				requiredTypes[pkg.anonymousStructID(structType)] = true
			} else if model.Obj().Pkg() == s.Pkg() {
				requiredTypes[model.Obj().Name()] = true
			}
		}
	}
	return requiredTypes
}

func (s *Service) calculateBindingLocalImports(pkg *Package) []string {
	requiredTypes := make(map[string]bool)

	for _, method := range s.Methods {
		maps.Copy(requiredTypes, s.bindingLocalImportsOf(method.JSInputs(), pkg))
		maps.Copy(requiredTypes, s.bindingLocalImportsOf(method.JSOutputs(), pkg))
	}

	result := lo.Keys(requiredTypes)
	slices.Sort(result)

	return result
}
