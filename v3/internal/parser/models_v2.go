package parser

import (
	"bytes"
	"go/types"
	"io"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

type VarAnalyzer struct {
	pkg       *Package
	Var       *types.Var
	models    map[*types.Named]bool
	recursive bool
}

func (p *Parameter) Models(pkg *Package, recursive bool) (models map[*types.Named]bool) {
	analyzer := &VarAnalyzer{
		pkg:       pkg,
		Var:       p.Var,
		recursive: recursive,
	}
	return analyzer.FindModels()
}

func (f *Field) Models(pkg *Package) (models map[*types.Named]bool) {
	analyzer := &VarAnalyzer{
		pkg:       pkg,
		Var:       f.Var,
		recursive: false,
	}
	return analyzer.FindModels()
}

func (a *VarAnalyzer) FindModels() (models map[*types.Named]bool) {
	a.models = make(map[*types.Named]bool)
	a.findModels(a.Var.Type())
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
			if a.recursive {
				a.findModelsOfNamed(x)
			}

			return
		case *types.Struct:
			if a.pkg == nil {
				return
			}
			named := types.NewNamed(types.NewTypeName(0, a.pkg.Types, a.pkg.anonymousStructID(x), nil), x, nil)
			a.models[named] = true
			if a.recursive {
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
}

func (a *VarAnalyzer) findModelsOfStruct(s *types.Struct) {
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		a.findModels(field.Type())
	}
}

func (m *BoundMethod) Models(pkg *Package, recursive bool) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	for _, param := range m.JSInputs() {
		maps.Copy(models, param.Models(pkg, recursive))
	}
	for _, param := range m.JSOutputs() {
		maps.Copy(models, param.Models(pkg, recursive))
	}
	return
}

func (s *Service) Models(pkg *Package) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)
	for _, method := range s.Methods {
		maps.Copy(models, method.Models(pkg, true))
	}
	return
}

func (p *Package) Models() (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)

	for _, s := range p.services {
		maps.Copy(models, s.Models(p))
	}
	return
}

func (p *Project) Models() (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)

	for _, pkg := range p.pkgs {
		maps.Copy(models, pkg.Models())
	}
	return
}

type JsonTag struct {
	name     string
	optional bool
	quoted   bool
	visible  bool
}

func parseTag(tag string) *JsonTag {
	tag = reflect.StructTag(tag).Get("json")
	if tag == "-" {
		return &JsonTag{
			"", false, false, false,
		}
	}

	parts := strings.Split(tag, ",")
	jsonTag := &JsonTag{
		name:    parts[0],
		visible: true,
	}

	for _, option := range parts[1:] {
		switch option {
		case "omitempty":
			jsonTag.optional = true
		case "string":
			jsonTag.quoted = true
		}
	}

	return jsonTag
}

type Field struct {
	*types.Var
	index   int
	jsonTag *JsonTag
}

func (f *Field) Name() string {
	var name string
	if len(f.jsonTag.name) > 0 {
		name = f.jsonTag.name
	} else {
		name = f.Var.Name()
	}

	if name == "" || name == "_" {
		return "$" + strconv.Itoa(f.index)
	} else if slices.Contains(reservedWords, name) {
		return "$" + name
	}
	return name
}

func (f *Field) JSType(pkg *Package) string {
	jstype, _ := JSType(f.Type(), pkg)
	return jstype
}

func (f *Field) DefaultValue(pkg *Package, mDef *ModelDefinitions) string {
	return DefaultValue(f.Type(), pkg, mDef)
}

func (f *Field) Exported() bool {
	return f.Var.Exported() && f.jsonTag.visible
}

func (f *Field) Optional() bool {
	return f.jsonTag.optional
}

func (f *Field) Quoted() bool {
	return f.jsonTag.quoted
}

type StructDef struct {
	*types.Struct
	Name string
}

func (s *StructDef) Fields() (fields []*Field) {
	for i := 0; i < s.NumFields(); i++ {
		field := &Field{
			Var:     s.Field(i),
			index:   i,
			jsonTag: parseTag(s.Tag(i)),
		}
		if field.Exported() {
			fields = append(fields, field)
		}
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

func (e *EnumDef) DefaultValue(fieldType types.Type, pkg *Package) string {
	jstype, _ := JSType(fieldType, pkg)

	// FIXME: order of e.Consts is not guaranteed
	// the default value may change between model generations
	return jstype + "." + e.Consts[0].Name
}

func (e *EnumDef) JSType(pkg *Package) string {
	jstype, _ := JSType(e.Type, pkg)
	return jstype
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

func (p *Project) GenerateModels() (result map[string]string, err error) {
	result = make(map[string]string)

	allModels := lo.Keys(p.Models())

	// split models into packages
	pkgModels := make(map[string][]*types.Named)
	for _, model := range allModels {
		pkgName := model.Obj().Pkg().String()
		if models, ok := pkgModels[pkgName]; ok {
			pkgModels[pkgName] = append(models, model)
		} else {
			pkgModels[pkgName] = []*types.Named{model}
		}
	}

	for _, pkg := range p.pkgs {

		models := pkgModels[pkg.Types.String()]

		// split models into structs and enums
		structDefs := make(map[string]*StructDef)
		enumDefs := make(map[string]*EnumDef)

		for _, model := range models {
			modelName := model.Obj().Name()

			switch t := model.Underlying().(type) {
			case *types.Basic:
				consts := []*ConstDef{}
				for name, c := range pkg.constantsOf(model) {
					consts = append(consts, &ConstDef{Name: name, Const: c})
				}

				if len(consts) == 0 {
					continue
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
			Imports: pkg.calculateModelImports(structDefs, p),

			Structs: structDefs,
			Enums:   enumDefs,

			ModelsFilename: p.options.ModelsFilename,
		}, p.options)

		if err != nil {
			return
		}

		packageDir := p.PackageDir(pkg.Types)
		result[packageDir] = buffer.String()
	}

	return
}

func (p *Package) calculateModelImports(m map[string]*StructDef, project *Project) map[string]string {
	result := make(map[string]string)
	pkg := p.Types

	for _, structDef := range m {
		for _, field := range structDef.Fields() {
			models := field.Models(p)
			if len(models) > 1 {
				panic("expected at most one model")
			}
			for model := range models {
				otherPkg := model.Obj().Pkg()
				if otherPkg.String() != pkg.String() {
					result[otherPkg.Name()] = project.RelativePackageDir(pkg, otherPkg)
				}
			}
		}
	}

	return result
}
