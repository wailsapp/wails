package parser

import (
	"bytes"
	"go/types"
	"io"
	"maps"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

func (p *Parameter) Models(pkg *Package, includeFields bool) (models map[*types.Named]bool) {
	analyzer := &VarAnalyzer{
		pkg:           pkg,
		parameter:     p,
		includeFields: includeFields,
	}
	return analyzer.FindModels()
}

type VarAnalyzer struct {
	pkg           *Package
	parameter     *Parameter
	models        map[*types.Named]bool
	includeFields bool
}

func (a *VarAnalyzer) FindModels() (models map[*types.Named]bool) {
	a.models = make(map[*types.Named]bool)
	a.findModels(a.parameter.Type())
	return a.models
}

func (a *VarAnalyzer) findModels(t types.Type) {
	for {
		switch x := t.(type) {
		case *types.Basic:
			return
		case *types.Slice:
			t = x.Elem()
		case *types.Map:
			t = x.Elem()
		case *types.Named:
			if _, ok := a.models[x]; ok {
				return
			}
			a.models[x] = true
			if a.includeFields {
				a.findModelsOfNamed(x)
			}

			return
		case *types.Struct:
			named := types.NewNamed(types.NewTypeName(0, a.pkg.Types, a.pkg.anonymousStructID(x), nil), x, nil)
			a.models[named] = true
			if a.includeFields {
				a.findModelsOfStruct(x)
			}
			return
		case *types.Pointer:
			t = x.Elem()
		default:
			return
		}

	}
}

func (a *VarAnalyzer) findModelsOfNamed(n *types.Named) {
	switch x := n.Underlying().(type) {
	case *types.Struct:
		a.findModelsOfStruct(x)
	}
	return
}

func (a *VarAnalyzer) findModelsOfStruct(s *types.Struct) {
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		a.findModels(field.Type())
	}
	return
}

func (m *BoundMethod) Models(pkg *Package) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	for _, param := range m.JSInputs() {
		maps.Copy(models, param.Models(pkg, true))
	}
	for _, param := range m.JSOutputs() {
		maps.Copy(models, param.Models(pkg, true))
	}
	return
}

func (s *Service) Models(pkg *Package) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	for _, method := range s.Methods() {
		maps.Copy(models, method.Models(pkg))
	}
	return
}

func (p *Package) addService(s *Service) {
	p.services = append(p.services, s)
}

func (p *Package) Models() (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)

	for _, s := range p.services {
		maps.Copy(models, s.Models(p))
	}
	return
}

type StructDef struct {
	*types.Struct
	Name string
}

func (s *StructDef) Fields() (fields []*Parameter) {
	for i := 0; i < s.NumFields(); i++ {
		fields = append(fields, &Parameter{index: i, Var: s.Field(i)})
	}
	return
}

type ConstDef struct {
	*types.Const
	Name string
}

func (c *ConstDef) Value() string {
	return c.Val().String()
}

type EnumDef struct {
	Name   string
	Type   *types.Basic
	Consts []*ConstDef
}

type ModelDefinitions struct {
	Package *Package
	Imports map[string]string

	Structs map[string]*StructDef
	Enums   map[string]*EnumDef

	ModelsFilename string
}

func (p *Project) generateModel(wr io.Writer, def *ModelDefinitions, options *flags.GenerateBindingsOptions) error {
	template := templates.ModelsJS
	if options.TS {
		if options.UseInterfaces {
			template = templates.InterfacesTS
		} else {
			template = templates.ModelsTS
		}
	}

	// TODO
	// Fix up TS names
	// for _, model := range def.Models {
	// 	model.Name = options.TSPrefix + model.Name + options.TSSuffix
	// }

	if err := template.Execute(wr, def); err != nil {
		println("Problem executing template: " + err.Error())
		return err
	}

	return nil
}

func (p *Project) GenerateModels(options *flags.GenerateBindingsOptions) (result map[string]string, err error) {
	result = make(map[string]string)

	for _, pkg := range p.pkgs {

		models := pkg.Models()

		// split models into structs and enums
		structDefs := make(map[string]*StructDef)
		enumDefs := make(map[string]*EnumDef)

		for model := range models {
			modelName := model.Obj().Name()

			switch t := model.Underlying().(type) {
			case *types.Basic:
				consts := []*ConstDef{}
				for name, c := range pkg.constantsOf(model) {
					consts = append(consts, &ConstDef{Name: name, Const: c})
				}

				def := &EnumDef{
					Name:   modelName,
					Type:   t,
					Consts: consts,
				}
				enumDefs[modelName] = def
			case *types.Struct:
				def := &StructDef{
					Name:   modelName,
					Struct: t,
				}
				structDefs[modelName] = def
			}
		}

		// generate model
		var buffer bytes.Buffer
		err = p.generateModel(&buffer, &ModelDefinitions{
			Package: pkg,
			Imports: pkg.calculateModelImports(structDefs),

			Structs: structDefs,
			Enums:   enumDefs,

			ModelsFilename: options.ModelsFilename,
		}, options)

		if err != nil {
			return
		}

		// Get the relative package path
		// TODO: get path of main package
		//relativePackageDir := RelativePackageDir(p.)
		// result[relativePackageDir] = buffer.String()

		result[pkg.Name] = buffer.String()
	}

	return
}

func (p *Package) calculateModelImports(m map[string]*StructDef) map[string]string {
	result := make(map[string]string)

	for _, structDef := range m {
		for i := 0; i < structDef.NumFields(); i++ {
			field := structDef.Field(i)
			if field.Pkg() != p.Types {
				// Find the relative path from the source directory to the target directory
				result[field.Pkg().Name()] = RelativeBindingsDir(p.Types, field.Pkg())
			}
		}
	}

	return result
}
