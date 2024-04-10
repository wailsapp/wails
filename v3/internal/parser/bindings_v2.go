package parser

import (
	"bytes"
	"go/types"
	"io"
	"maps"
	"slices"
	"strconv"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

// type ParameterType struct {
// 	Name       string
// 	IsStruct   bool
// 	IsSlice    bool
// 	IsPointer  bool
// 	IsEnum     bool
// 	IsVariadic bool
// 	MapKey     *ParameterType
// 	MapValue   *ParameterType
// }

type Parameter struct {
	*types.Var
	index int

	Parent *BoundMethod
}

func (p *Parameter) Name() (name string) {
	name = p.Var.Name()
	if name == "" || name == "_" {
		return "$" + strconv.Itoa(p.index)
	} else if slices.Contains(reservedWords, name) {
		return "$" + name
	}
	return
}

func (p *Parameter) Optional() bool {
	// TODO
	return false
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
		switch y := x.Underlying().(type) {
		case *types.Struct:
			if x.Obj() != nil {
				return "(new " + x.Obj().Name() + "())"
			} else {
				return "(new " + pkg.anonymousStructID(y) + "())"
			}
		case *types.Basic:
			return DefaultValue(y, pkg)
		}
	case *types.Map:
		return "{}"
	case *types.Pointer:
		return "null"
	case *types.Struct:
		return "(new " + pkg.anonymousStructID(x) + "())"
	}
	return "null"
}

func (p *Parameter) DefaultValue(pkg *Package) string {
	return DefaultValue(p.Type(), pkg)
}

func (p *Parameter) Variadic() bool {
	s := p.Parent.Signature()
	return s.Variadic() && p.index == s.Params().Len()-1
}

func (p *Package) namespaceOf(t *types.TypeName) string {
	if p.Types == t.Pkg() {
		return ""
	}
	return t.Pkg().Name() + "."
}

// JSTypes returns the corresponding javascript type to the given types.Type
// The second return value indicates whether parentheses are needed
func JSType(t types.Type, pkg *Package) (string, bool) {

	switch x := t.(type) {
	case *types.Basic:
		switch x.Kind() {
		case types.String:
			return "string", false
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr, types.Float32, types.Float64:
			return "number", false
		case types.Bool:
			return "boolean", false
		default:
			return "any", false
		}
	case *types.Slice:
		jstype, needsParentheses := JSType(x.Elem(), pkg)
		if needsParentheses {
			return "(" + jstype + ")[]", false
		}
		return jstype + "[]", false
	case *types.Array:
		jstype, needsParentheses := JSType(x.Elem(), pkg)
		if needsParentheses {
			return "(" + jstype + ")[]", false
		}
		return jstype + "[]", false
	case *types.Named:
		return pkg.namespaceOf(x.Obj()) + x.Obj().Name(), false
	case *types.Map:
		jstype, _ := JSType(x.Elem(), pkg)
		return "{ [_: string]: " + jstype + " }", false
	case *types.Pointer:
		jstype, _ := JSType(x.Elem(), pkg)
		return jstype + " | null", true
	case *types.Struct:
		return pkg.anonymousStructID(x), false
	}
	return "any", false
}

func (p *Parameter) JSType(pkg *Package) string {
	jstype, _ := JSType(p.Type(), pkg)
	return jstype
}

type BoundMethod struct {
	*types.Func
	ID  uint32
	FQN string

	Service *Service

	// Name       string
	// DocComment string
	// Inputs     []*Parameter
	// Outputs    []*Parameter
	// ID         uint32
	// Alias      *uint32
}

func (m *BoundMethod) embedTuple(tuple *types.Tuple) (result []*Parameter) {
	if tuple == nil {
		return
	}

	for i := 0; i < tuple.Len(); i++ {
		result = append(result, &Parameter{tuple.At(i), i, m})
	}
	return
}

func (m *BoundMethod) Signature() *types.Signature {
	// Type of *types.Func is always a *types.Signature
	return m.Type().(*types.Signature)
}

func (m *BoundMethod) Params() []*Parameter {
	tuple := m.Signature().Params()
	return m.embedTuple(tuple)
}

func (m *BoundMethod) Results() []*Parameter {
	tuple := m.Signature().Results()
	return m.embedTuple(tuple)
}

func (m *BoundMethod) JSInputs() []*Parameter {
	params := m.Params()

	if len(params) > 0 {
		if named, ok := params[0].Type().(*types.Named); ok && named.Obj() != nil {
			if named.Obj().Name() == "Context" && named.Obj().Pkg().Name() == "context" {
				return params[1:]
			}
		}
	}

	return params
}

func (m *BoundMethod) JSOutputs() (outputs []*Parameter) {
	for _, output := range m.Results() {
		if types.TypeString(output.Var.Type(), nil) == "error" {
			continue
		}
		outputs = append(outputs, output)
	}

	return outputs
}

type BindingDefinitions struct {
	Package      *Package
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

func (p *Project) GenerateBindings(options *flags.GenerateBindingsOptions) (result map[string]map[string]string, err error) {
	result = make(map[string]map[string]string)

	for _, pkg := range p.pkgs {
		bindings, err := pkg.GenerateBindings(options)
		if err != nil {
			return nil, err
		}
		result[pkg.Name] = bindings
	}
	return
}

func (p *Package) GenerateBindings(options *flags.GenerateBindingsOptions) (result map[string]string, err error) {
	result = make(map[string]string)

	for _, service := range p.services {
		structName := service.Name()
		methods := service.Methods()

		var buffer bytes.Buffer
		err = generateBinding(&buffer, &BindingDefinitions{
			Package:      p,
			Imports:      service.calculateBindingImports(p),
			LocalImports: service.calculateBindingLocalImports(p),

			Methods: methods,

			ModelsFilename:    options.ModelsFilename,
			UseBundledRuntime: options.UseBundledRuntime,
			UseNames:          options.UseNames,
		}, options)

		if err != nil {
			return
		}

		result[structName] = buffer.String()
	}
	return
}

func (s *Service) bindingImportsOf(params []*Parameter, pkg *Package) map[string]string {
	result := make(map[string]string)

	for _, param := range params {
		models := param.Models(pkg, false)
		for model := range models {
			if model.Obj() != nil && model.Obj().Pkg() != s.Pkg() {
				otherPkg := model.Obj().Pkg()
				result[otherPkg.Name()] = RelativeBindingsDir(s.Pkg(), otherPkg)
			}
		}
	}
	return result
}

func (s *Service) calculateBindingImports(pkg *Package) map[string]string {
	result := make(map[string]string)

	for _, method := range s.Methods() {
		maps.Copy(result, s.bindingImportsOf(method.JSInputs(), pkg))
		maps.Copy(result, s.bindingImportsOf(method.JSOutputs(), pkg))
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

	for _, method := range s.Methods() {
		maps.Copy(requiredTypes, s.bindingLocalImportsOf(method.JSInputs(), pkg))
		maps.Copy(requiredTypes, s.bindingLocalImportsOf(method.JSOutputs(), pkg))
	}

	result := lo.Keys(requiredTypes)
	slices.Sort(result)

	return result
}
