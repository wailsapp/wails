package parser

import (
	"bytes"
	"cmp"
	"go/types"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

type VarAnalyzer struct {
	pkg       *Package
	Var       *types.Var
	models    map[*types.Named]bool
	recursive bool
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
			named := types.NewNamed(types.NewTypeName(0, a.pkg.Package, a.pkg.anonymousStructID(x), nil), x, nil)
			a.models[named] = true
			if a.recursive {
				a.findModelsOfStruct(named.Obj().Name(), x)
			}
			return
		case *types.Pointer:
			t = x.Elem()
		case *types.Alias:
			// Note: at the moment it is not possible to access the RHS of an alias
			// https://github.com/golang/go/issues/66559
			t = aliasToNamed(x)
		default:
			return
		}

	}
}

func (a *VarAnalyzer) findModelsOfNamed(n *types.Named) {
	switch x := n.Underlying().(type) {
	case *types.Struct:
		a.findModelsOfStruct(n.Obj().Name(), x)
	}
}

func (a *VarAnalyzer) findModelsOfStruct(name string, s *types.Struct) {
	structDef := &StructDef{
		Name:   name,
		Struct: s,
	}
	for _, field := range structDef.Fields() {
		a.findModels(field.Type())
	}
}

func (p *Parameter) findModels(found map[*types.Named]bool, pkg *Package, recursive bool) {
	analyzer := &VarAnalyzer{
		pkg:       pkg,
		Var:       p.Var,
		recursive: recursive,
		models:    found,
	}
	analyzer.findModels(p.Var.Type())
}

func (f *Field) FindModels(pkg *Package) map[*types.Named]bool {
	analyzer := &VarAnalyzer{
		pkg:       pkg,
		Var:       f.Var,
		recursive: false,
		models:    make(map[*types.Named]bool),
	}
	analyzer.findModels(f.Var.Type())
	return analyzer.models
}

func (m *BoundMethod) findModels(found map[*types.Named]bool, pkg *Package, recursive bool) {
	for _, param := range m.JSInputs() {
		param.findModels(found, pkg, recursive)
	}
	for _, param := range m.JSOutputs() {
		param.findModels(found, pkg, recursive)
	}
}

func (m *BoundMethod) FindModels(pkg *Package, recursive bool) map[*types.Named]bool {
	models := make(map[*types.Named]bool)
	m.findModels(models, pkg, recursive)
	return models
}

func (s *Service) findModels(found map[*types.Named]bool, pkg *Package, recursive bool) {
	for _, method := range s.Methods {
		method.findModels(found, pkg, recursive)
	}
}

func (s *Service) FindModels(pkg *Package, recursive bool) map[*types.Named]bool {
	models := make(map[*types.Named]bool)
	s.findModels(models, pkg, recursive)
	return models
}

func FindModels(pkgs []*Package) (models map[*types.Named]bool) {
	models = make(map[*types.Named]bool)

	for _, pkg := range pkgs {
		for _, s := range pkg.services {
			for _, method := range s.Methods {
				method.findModels(models, pkg, true)
			}
		}
	}
	return models
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

func DefaultValue(t types.Type, pkg *Package) string {
	switch x := t.(type) {
	case *types.Basic:
		switch x.Kind() {
		case types.String:
			return "\"\""
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr, types.Float32, types.Float64:
			return "0"
		case types.Bool:
			return "false"
		default:
			return "null"
		}
	case *types.Slice, *types.Array:
		return "[]"
	case *types.Named:
		return pkg.DefaultValue(x)
	case *types.Map:
		return "{}"
	case *types.Pointer:
		return "null"
	case *types.Struct:
		return "(new " + pkg.anonymousStructID(x) + "())"
	case *types.Alias:
		return DefaultValue(aliasToNamed(x), pkg)
	}
	return "null"
}

func (f *Field) DefaultValue(pkg *Package) string {
	value := DefaultValue(f.Type(), pkg)

	if f.Quoted() {
		value = "`" + value + "`"
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

func (s *StructDef) DefaultValue(pkg *Package) string {
	return "(new " + s.Name + "())"
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
				filteredWarning.Printfln("interface as field: ignoring field %s.%s of type %s", s.Name, field.Name(), field.Type().String())
			}
		} else if field.Exported() {
			switch field.Type().Underlying().(type) {
			case *types.Interface:
				filteredWarning.Printfln("interface as field: ignoring field %s.%s of type %s", s.Name, field.Name(), field.Type().String())
			default:
				fields = append(fields, field)
			}

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
	Value string
	Name  string
}

type EnumDef struct {
	Name   string
	Type   *types.Basic
	Consts []*ConstDef
}

func (e *EnumDef) DefaultValue(pkg *Package) string {
	return e.Name + "." + e.Consts[0].Name
}

func (e *EnumDef) JSType(pkg *Package) string {
	jstype, _ := JSType(e.Type, pkg)
	return jstype
}

type AliasDef struct {
	Name string
	Type types.Type
}

func (a *AliasDef) DefaultValue(pkg *Package) string {
	return DefaultValue(a.Type.Underlying(), pkg)
}

func (a *AliasDef) JSType(pkg *Package) string {
	jstype, _ := JSType(a.Type.Underlying(), pkg)
	return jstype
}

type Model interface {
	DefaultValue(pkg *Package) string
}

type Models struct {
	Structs map[string]*StructDef
	Enums   map[string]*EnumDef
	Aliases map[string]*AliasDef
}

func NewModels() *Models {
	return &Models{
		Structs: make(map[string]*StructDef),
		Enums:   make(map[string]*EnumDef),
		Aliases: make(map[string]*AliasDef),
	}
}

func (m *Models) Length() int {
	return len(m.Structs) + len(m.Enums) + len(m.Aliases)
}

func (p *Package) addModel(model *types.Named, marshaler, textMarshaler *types.Interface) {

	modelName := model.Obj().Name()
	models := p.models

	if types.Implements(model, marshaler) {
		pterm.Warning.Printfln("Generator can not predict json keys of model %s, because it implements json.Marshaler", model.String())
	}

	if types.Implements(model, textMarshaler) {
		models.Aliases[modelName] = &AliasDef{
			Name: modelName,
			Type: types.Typ[types.String],
		}
		return
	}

	switch t := model.Underlying().(type) {
	case *types.Basic:
		consts := p.constantsOf(model)

		if len(consts) == 0 {
			models.Aliases[modelName] = &AliasDef{
				Name: modelName,
				Type: model,
			}
		} else {
			models.Enums[modelName] = &EnumDef{
				Name:   modelName,
				Type:   t,
				Consts: consts,
			}
		}

	case *types.Struct:
		models.Structs[modelName] = &StructDef{
			Name:   modelName,
			Struct: t,
		}
	default:
		models.Aliases[modelName] = &AliasDef{
			Name: modelName,
			Type: model,
		}
	}

}

func (p *Package) DefaultValue(named *types.Named) string {
	pkgPath := named.Obj().Pkg().Path()
	modelName := named.Obj().Name()
	model := p.project.getModel(pkgPath, modelName)
	if model != nil {
		return model.DefaultValue(p)
	}

	return DefaultValue(named.Underlying(), p)
}

type ModelDefinitions struct {
	Package *Package
	Imports map[string]string

	Structs map[string]*StructDef
	Enums   map[string]*EnumDef
	Aliases map[string]*AliasDef

	ModelsFilename string
}

func (p *Project) getModel(pkgPath, modelName string) Model {
	for _, pkg := range p.pkgs {
		if pkg.Package.Path() == pkgPath {
			if structDef, ok := pkg.models.Structs[modelName]; ok {
				return structDef
			}
			if enum, ok := pkg.models.Enums[modelName]; ok {
				return enum
			}
			if alias, ok := pkg.models.Aliases[modelName]; ok {
				return alias
			}
		}
	}
	return nil
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

	if err := template.Execute(wr, def); err != nil {
		println("Problem executing template: " + err.Error())
		return err
	}

	return nil
}

func (p *Project) GenerateModels() (result map[string]string, err error) {
	result = make(map[string]string)

	for _, pkg := range p.pkgs {

		models := pkg.models

		if models.Length() == 0 {
			continue
		}

		// update stats
		p.Stats.NumModels += len(models.Structs)
		p.Stats.NumEnums += len(models.Enums)
		p.Stats.NumAliases += len(models.Aliases)

		// generate model
		var buffer bytes.Buffer
		err = p.generateModel(&buffer, &ModelDefinitions{
			Package: pkg,
			Imports: pkg.calculateModelImports(models.Structs, p),

			Structs: models.Structs,
			Enums:   models.Enums,
			Aliases: models.Aliases,

			ModelsFilename: p.options.ModelsFilename,
		}, p.options)

		if err != nil {
			return
		}

		packageDir := p.PackageDir(pkg.Package)
		result[packageDir] = buffer.String()
	}

	return
}

func (p *Package) calculateModelImports(m map[string]*StructDef, project *Project) map[string]string {
	result := make(map[string]string)
	pkg := p.Package

	for _, structDef := range m {
		for _, field := range structDef.Fields() {
			models := field.FindModels(p)
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
