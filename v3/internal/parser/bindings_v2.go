package parser

import (
	"bytes"
	"go/types"
	"io"
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

	Parent func() *BoundMethod
}

func (p *Parameter) Name() (name string) {
	name = p.Var.Name()
	if name == "" || name == "_" {
		return "$" + strconv.Itoa(p.index)
	}
	return
}

func (p *Parameter) Variadic() bool {
	s := p.Parent().Signature()
	return s.Variadic() && p.index == s.Params().Len()-1
}

func (p *Parameter) js(t types.Type) string {
	metod := p.Parent()
	project := metod.Parent()

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
		return p.js(x.Elem()) + "[]"
	case *types.Named:
		// TODO: add package name for non-local imports, add namespace method
		return x.Obj().Name()
	case *types.Map:
		return "{ [_: string]: " + p.js(x.Elem()) + " }"
	case *types.Pointer:
		return "(" + p.js(x.Elem()) + " | null)"
	case *types.Struct:
		return project.anonymousStructID(x)
	}
	return "any"
}

func (p *Parameter) JSType() string {
	return p.js(p.Type())
}

type BoundMethod struct {
	*types.Func
	ID uint32

	Parent func() *Project

	// Name       string
	// DocComment string
	// Inputs     []*Parameter
	// Outputs    []*Parameter
	// ID         uint32
	// Alias      *uint32
}

func (p *Project) newMethods(service *types.TypeName) (methods []*BoundMethod) {
	// TODO unsafe
	named := service.Type().(*types.Named)

	for i := 0; i < named.NumMethods(); i++ {
		methods = append(methods, &BoundMethod{
			Func: named.Method(i),
			//TODO assign ID
			ID:     0,
			Parent: func() *Project { return p },
		})
	}
	return
}

func (m *BoundMethod) DocComment() string {
	// TODO
	return ""
}

func (m *BoundMethod) newParameters(tuple *types.Tuple) (result []*Parameter) {
	if tuple == nil {
		return
	}

	for i := 0; i < tuple.Len(); i++ {
		result = append(result, &Parameter{tuple.At(i), i, func() *BoundMethod { return m }})
	}
	return
}

func (m *BoundMethod) Signature() *types.Signature {
	// The Type of *types.Func is always a *types.Signature
	return m.Type().(*types.Signature)
}

func (m *BoundMethod) Params() []*Parameter {
	tuple := m.Signature().Params()
	return m.newParameters(tuple)
}

func (m *BoundMethod) Results() []*Parameter {
	tuple := m.Signature().Results()
	return m.newParameters(tuple)
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
	Imports      map[string]string
	LocalImports []string

	Struct  string
	Methods []*BoundMethod

	ModelsFilename    string
	UseBundledRuntime bool
	UseNames          bool
}

func GenerateBinding(wr io.Writer, def *BindingDefinitions, options *flags.GenerateBindingsOptions) error {
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

func (p *Project) GenerateBindings(services []*types.TypeName, options *flags.GenerateBindingsOptions) (result map[string]map[string]string, err error) {
	result = make(map[string]map[string]string)

	for _, service := range services {

		structName := service.Name()
		pkgName := service.Pkg().Name()
		methods := p.newMethods(service)

		var buffer bytes.Buffer
		err = GenerateBinding(&buffer, &BindingDefinitions{
			Imports:      p.calculateBindingImports(service, methods),
			LocalImports: p.calculateBindingLocalImports(service, methods),

			// Struct:  pkgAlias(pkg) + "." + structName,
			Methods: methods,

			ModelsFilename:    options.ModelsFilename,
			UseBundledRuntime: options.UseBundledRuntime,
			UseNames:          options.UseNames,
		}, options)

		if err != nil {
			return
		}

		if _, ok := result[pkgName]; !ok {
			result[pkgName] = make(map[string]string)
		}

		result[pkgName][structName] = buffer.String()
	}
	return
}

func (p *Project) calculateBindingImports(service *types.TypeName, methods []*BoundMethod) map[string]string {
	result := make(map[string]string)

	for _, method := range methods {
		for _, param := range method.JSInputs() {
			if param.Pkg() != service.Pkg() {
				// Find the relative path from the source package to the target package
				result[param.Pkg().Name()] = p.RelativeBindingsDir(service.Pkg(), param.Pkg())
			}
		}

		for _, param := range method.JSOutputs() {
			if param.Pkg() != service.Pkg() {
				// Find the relative path from the source package to the target package
				result[param.Pkg().Name()] = p.RelativeBindingsDir(service.Pkg(), param.Pkg())
			}
		}
	}

	return result
}

func (p *Project) calculateBindingLocalImports(service *types.TypeName, methods []*BoundMethod) []string {
	requiredTypes := make(map[string]bool)

	for _, method := range methods {
		for _, param := range method.JSInputs() {
			if param.Pkg() == service.Pkg() {
				models := param.Models()
				for _, model := range models {
					if s, ok := model.Underlying().(*types.Struct); ok && model.Obj() == nil {
						requiredTypes[p.anonymousStructID(s)] = true
					} else {
						requiredTypes[model.Obj().Name()] = true
					}
				}
			}
		}

		for _, param := range method.JSOutputs() {
			if param.Pkg() == service.Pkg() {
				models := param.Models()
				for _, model := range models {
					if s, ok := model.Underlying().(*types.Struct); ok && model.Obj() == nil {
						requiredTypes[p.anonymousStructID(s)] = true
					} else {
						requiredTypes[model.Obj().Name()] = true
					}
				}
			}
		}
	}

	result := lo.Keys(requiredTypes)
	slices.Sort(result)

	return result
}
