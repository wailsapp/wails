package parser

import (
	"bytes"
	"cmp"
	"go/types"
	"io"
	"maps"
	"reflect"
	"slices"
	"strings"

	"github.com/pterm/pterm"
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
	path    []int
	jsonTag *JsonTag
	origin  *StructDef
}

func (f *Field) JSName() string {
	if len(f.jsonTag.name) > 0 {
		return f.jsonTag.name
	}
	return f.Var.Name()
}

func (f *Field) nameFromTag() bool {
	return len(f.jsonTag.name) > 0
}

func (f *Field) JSType(pkg *Package) string {
	jstype, _ := JSType(f.Type(), pkg)

	// use Typescript template literal types to type encoding/json quoted fields
	if f.Quoted() {
		if jstype == "string" {
			jstype = "`\"${" + jstype + "}\"`"
		} else {
			jstype = "`${" + jstype + "}`"
		}
	}

	return jstype
}

func (f *Field) DefaultValue(pkg *Package, mDef *ModelDefinitions) string {
	value := DefaultValue(f.Type(), pkg, mDef)

	if f.Quoted() {
		value = `"` + value + `"`
	}

	return value
}

func (f *Field) Exported() bool {
	if !f.jsonTag.visible {
		return false
	}
	if f.Embedded() && f.nameFromTag() {
		return true
	}

	return f.Var.Exported()
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

func (s *StructDef) allFields() []*Field {
	fields := []*Field{}

	for i := 0; i < s.NumFields(); i++ {
		field := &Field{
			Var:     s.Field(i),
			path:    []int{i},
			jsonTag: parseTag(s.Tag(i)),
			origin:  s,
		}
		if field.Embedded() && !field.nameFromTag() {
			switch fieldType := field.Type().Underlying().(type) {
			case *types.Struct:
				embDef := &StructDef{
					Name:   field.Type().(*types.Named).Obj().Name(),
					Struct: fieldType,
				}
				embeddedFields := embDef.allFields()
				for _, embeddedField := range embeddedFields {
					embeddedField.path = append([]int{i}, embeddedField.path...)
					fields = append(fields, embeddedField)
				}
			case *types.Basic:
				if field.Exported() {
					fields = append(fields, field)
				}
			case *types.Interface:
				pterm.Warning.Printfln("ignoring interface %v", fieldType)
			}
		} else if field.Exported() {
			fields = append(fields, field)
		}
	}
	return fields
}

func (s *StructDef) Fields() []*Field {
	fields := s.allFields()

	// sort fields
	slices.SortFunc(fields, func(f1 *Field, f2 *Field) int {
		// sort by name first
		if diff := strings.Compare(f1.JSName(), f2.JSName()); diff != 0 {
			return diff
		}

		// break ties by depth of occurrence
		if diff := cmp.Compare(len(f1.path), len(f2.path)); diff != 0 {
			return diff
		}

		// break ties by presence of json tag (prioritize presence)
		if f1.nameFromTag() != f2.nameFromTag() {
			if f1.nameFromTag() {
				return -1
			} else {
				return 1
			}
		}

		// break ties by order of occurrence
		return slices.Compare(f1.path, f2.path)
	})

	count := 0

	// keep for each name the dominant field, drop those for which ties
	// still exist (ignoring order of occurrence)
	for i, j := 0, 1; j <= len(fields); j++ {
		if j < len(fields) && fields[i].JSName() == fields[j].JSName() {
			continue
		}

		// if there is only one field with the current name, or there is a dominant one, keep it
		if i+1 == j || len(fields[i].path) != len(fields[i+1].path) || fields[i].nameFromTag() != fields[i+1].nameFromTag() {
			fields[count] = fields[i]
			count++
		}

		i = j
	}
	result := fields[:count]

	// sort by order of occurrence
	slices.SortFunc(result, func(f1 *Field, f2 *Field) int {
		return slices.Compare(f1.path, f2.path)
	})

	return result
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
	return jstype + "." + e.Consts[0].Name
}

func (e *EnumDef) JSType(pkg *Package) string {
	jstype, _ := JSType(e.Type, pkg)
	return jstype
}

type BasicType struct {
	Type *types.Basic
	Name string
}

func (b *BasicType) JSType(pkg *Package) string {
	jstype, _ := JSType(b.Type, pkg)
	return jstype
}

type ModelDefinitions struct {
	Package *Package
	Imports map[string]string

	Structs    map[string]*StructDef
	Enums      map[string]*EnumDef
	BasicTypes map[string]*BasicType

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

		if len(models) == 0 {
			continue
		}

		// split models into structs, enums and basic types
		structDefs := make(map[string]*StructDef)
		enumDefs := make(map[string]*EnumDef)
		basicTypes := make(map[string]*BasicType)

		for _, model := range models {
			modelName := model.Obj().Name()

			switch t := model.Underlying().(type) {
			case *types.Basic:
				consts := pkg.constantsOf(model)

				if len(consts) == 0 {
					basicTypes[modelName] = &BasicType{
						Name: modelName,
						Type: t,
					}
				} else {
					enumDefs[modelName] = &EnumDef{
						Name:   modelName,
						Type:   t,
						Consts: consts,
					}
				}

			case *types.Struct:
				structDefs[modelName] = &StructDef{
					Name:   modelName,
					Struct: t,
				}
			}
		}

		p.Stats.NumModels = len(structDefs)
		p.Stats.NumEnums = len(enumDefs)

		// generate model
		var buffer bytes.Buffer
		err = p.generateModel(&buffer, &ModelDefinitions{
			Package: pkg,
			Imports: pkg.calculateModelImports(structDefs, p),

			Structs:    structDefs,
			Enums:      enumDefs,
			BasicTypes: basicTypes,

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
