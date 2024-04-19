package parser

import (
	"bytes"
	"go/types"
	"io"
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
		packageDir := p.PackageDir(pkg.Package)
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

		models := lo.Keys(service.FindModels(p, false))

		var buffer bytes.Buffer
		err = generateBinding(&buffer, &BindingDefinitions{
			Package:      p,
			Service:      service,
			Imports:      service.calculateBindingImports(models, project),
			LocalImports: service.calculateBindingLocalImports(models, p),

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

func (s *Service) calculateBindingImports(models []*types.Named, project *Project) map[string]string {
	result := make(map[string]string)

	for _, model := range models {
		if model.Obj() != nil && model.Obj().Pkg() != s.Pkg() {
			otherPkg := model.Obj().Pkg()
			result[otherPkg.Name()] = project.RelativePackageDir(s.Pkg(), otherPkg)
		}
	}

	return result
}

func (s *Service) calculateBindingLocalImports(models []*types.Named, pkg *Package) []string {
	requiredTypes := make(map[string]bool)

	for _, model := range models {
		if structType, ok := model.Underlying().(*types.Struct); ok && model.Obj() == nil {
			requiredTypes[pkg.anonymousStructID(structType)] = true
		} else if model.Obj().Pkg() == s.Pkg() {
			requiredTypes[model.Obj().Name()] = true
		}
	}

	result := lo.Keys(requiredTypes)
	slices.Sort(result)

	return result
}
