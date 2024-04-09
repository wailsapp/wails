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
	}
	return
}

func (p *Parameter) Optional() bool {
	// TODO
	return false
}

func (p *Parameter) DefaultValue() string {
	// TODO
	return "null"
}

func (p *Parameter) Variadic() bool {
	s := p.Parent.Signature()
	return s.Variadic() && p.index == s.Params().Len()-1
}

func JSType(t types.Type, pkg *Package) string {

	switch x := t.(type) {
	case *types.Basic:
		switch x.Kind() {
		case types.String:
			return "string"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr, types.Float32, types.Float64:
			return "number"
		case types.Bool:
			return "boolean"
		default:
			return "any"
		}
	case *types.Slice:
		return JSType(x.Elem(), pkg) + "[]"
	case *types.Named:
		// TODO: add package name for non-local imports, add namespace method
		return x.Obj().Name()
	case *types.Map:
		return "{ [_: string]: " + JSType(x.Elem(), pkg) + " }"
	case *types.Pointer:
		return "(" + JSType(x.Elem(), pkg) + " | null)"
	case *types.Struct:
		return pkg.anonymousStructID(x)
	}
	return "any"
}

func (p *Parameter) JSType(pkg *Package) string {
	return JSType(p.Type(), pkg)
}

type BoundMethod struct {
	*types.Func
	ID uint32

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

	// TODO
	// if len(params) > 0 {
	// 	if firstArg := params[0]; firstArg.Type.Package.Path == "context" && firstArg.Type.Name == "Context" {
	// 		return params[1:]
	// 	}
	// }
	return params
}

func (m *BoundMethod) JSOutputs() (outputs []*Parameter) {

	// TODO
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
			Imports:      service.calculateBindingImports(),
			LocalImports: service.calculateBindingLocalImports(p),

			// Struct:  pkgAlias(pkg) + "." + structName,
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

func (s *Service) bindingImportsOf(params []*Parameter) map[string]string {
	result := make(map[string]string)

	for _, param := range params {
		if param.Pkg() != s.Pkg() {
			// Find the relative path from the source package to the target package
			result[param.Pkg().Name()] = RelativeBindingsDir(s.Pkg(), param.Pkg())
		}
	}
	return result
}

func (s *Service) calculateBindingImports() map[string]string {
	result := make(map[string]string)

	for _, method := range s.Methods() {
		maps.Copy(result, s.bindingImportsOf(method.JSInputs()))
		maps.Copy(result, s.bindingImportsOf(method.JSOutputs()))
	}

	return result
}

func (s *Service) bindingLocalImportsOf(params []*Parameter, pkg *Package) map[string]bool {
	requiredTypes := make(map[string]bool)

	for _, param := range params {
		if param.Pkg() == s.Pkg() {
			models := param.Models(pkg)
			for model := range models {
				if s, ok := model.Underlying().(*types.Struct); ok && model.Obj() == nil {
					requiredTypes[pkg.anonymousStructID(s)] = true
				} else {
					requiredTypes[model.Obj().Name()] = true
				}
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
