package parser

import (
	"bytes"
	"go/types"
	"io"
	"maps"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

func (p *Parameter) Models() map[*types.Named]bool {
	return p.modelsIn(p.Type())
}

func (p *Parameter) modelsIn(t types.Type) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	method := p.Parent
	project := method.Parent.Parent

	for {
		switch x := t.(type) {
		case *types.Basic:
			return
		case *types.Slice:
			t = x.Elem()
		case *types.Map:
			t = x.Elem()
		case *types.Named:
			if _, ok := models[x]; ok {
				return
			}
			models[x] = true
			maps.Copy(models, p.modelsInNamed(x))
			return
		case *types.Struct:
			named := types.NewNamed(types.NewTypeName(0, method.Pkg(), project.anonymousStructID(x), nil), x, nil)
			models[named] = true
			maps.Copy(models, p.modelsInStruct(x))
			return
		case *types.Pointer:
			t = x.Elem()
		default:
			return
		}

	}
}

func (p *Parameter) modelsInNamed(n *types.Named) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	switch x := n.Underlying().(type) {
	case *types.Struct:
		maps.Copy(models, p.modelsInStruct(x))
	}
	return
}

func (p *Parameter) modelsInStruct(s *types.Struct) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		maps.Copy(models, p.modelsIn(field.Type()))
	}
	return
}

func (m *BoundMethod) Models() (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	for _, param := range m.JSInputs() {
		maps.Copy(models, param.Models())
	}
	for _, param := range m.JSOutputs() {
		maps.Copy(models, param.Models())
	}
	return
}

func (s *Service) Models() (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	for _, method := range s.Methods() {
		maps.Copy(models, method.Models())
	}
	return
}

type Package struct {
	*types.Package
	services []*Service
}

func (p *Package) addService(s *Service) {
	p.services = append(p.services, s)
}

func (p *Package) Models() (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)

	for _, s := range p.services {
		maps.Copy(models, s.Models())
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

func (s *StructDef) DocComment() string {
	// TODO
	return ""
}

type ConstDef struct {
	*types.Const
	Name string
}

func (c *ConstDef) Value() string {
	return c.Val().String()
}

func (c *ConstDef) DocComment() string {
	// TODO
	return ""
}

type EnumDef struct {
	Name   string
	Type   *types.Basic
	Consts []*ConstDef
}

func (e *EnumDef) DocComment() string {
	// TODO
	return ""
}

type ModelDefinitions struct {
	Project *Project
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

	services, err := p.Services()
	if err != nil {
		return
	}

	pkgs := make(map[*types.Package]*Package)
	for _, service := range services {
		if pkg, ok := pkgs[service.Pkg()]; ok {
			pkg.addService(service)
		} else {
			pkgs[service.Pkg()] = &Package{
				service.Pkg(),
				[]*Service{service},
			}
		}
	}

	for _, pkg := range pkgs {

		models := pkg.Models()

		// split models into structs and enums
		structDefs := make(map[string]*StructDef)
		enumDefs := make(map[string]*EnumDef)

		for model := range models {
			modelName := model.Obj().Name()

			switch t := model.Underlying().(type) {
			case *types.Basic:
				consts := []*ConstDef{}
				for name, c := range p.constantsOf(model) {
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
			Project: p,
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

		result[pkg.Name()] = buffer.String()
	}

	return
}

func (p *Package) calculateModelImports(m map[string]*StructDef) map[string]string {
	result := make(map[string]string)

	for _, structDef := range m {
		for i := 0; i < structDef.NumFields(); i++ {
			field := structDef.Field(i)
			if field.Pkg() != p.Package {
				// Find the relative path from the source directory to the target directory
				result[field.Pkg().Name()] = RelativeBindingsDir(p.Package, field.Pkg())
			}
		}
	}

	return result
}
