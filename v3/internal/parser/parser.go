package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/hash"
	"golang.org/x/tools/go/packages"
)

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
	if p.Types.String() == t.Pkg().String() {
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

type Service struct {
	*types.TypeName
}

func (s *Service) Methods() (methods []*BoundMethod) {
	if named, ok := s.Type().(*types.Named); ok {
		for i := 0; i < named.NumMethods(); i++ {
			fqn := fmt.Sprintf("%s.%s.%s", s.Pkg().Name(), s.Name(), named.Method(i).Name())

			id, err := hash.Fnv(fqn)
			if err != nil {
				panic("Failed to hash fqn")
			}

			methods = append(methods, &BoundMethod{
				Func:    named.Method(i),
				FQN:     fqn,
				ID:      id,
				Service: s,
			})
		}
	}
	return
}

type Package struct {
	*packages.Package
	services         []*Service
	anonymousStructs map[string]string
	doc              *doc.Package
}

func BuildPackages(buildFlags []string, pkgs []*packages.Package, services []*Service) ([]*Package, error) {
	pPkgMap := make(map[*types.Package]*packages.Package)
	result := make(map[*types.Package]*Package)

	for _, pkg := range pkgs {
		pPkgMap[pkg.Types] = pkg
	}

	for _, service := range services {
		if pkg, ok := result[service.Pkg()]; ok {
			pkg.addService(service)
		} else {
			pPkg, ok := pPkgMap[service.Pkg()]
			if !ok {
				var err error
				pPkg, err = LoadPackage(buildFlags, true, service.Pkg().Path())
				if err != nil {
					return nil, err
				}
				pPkgMap[service.Pkg()] = pPkg
			}

			result[service.Pkg()] = &Package{
				Package:          pPkg,
				services:         []*Service{service},
				anonymousStructs: make(map[string]string),
				doc: NewDoc(pPkg),
			}
		}
	}
	return lo.Values(result), nil
}

func (p *Package) addService(s *Service) {
	p.services = append(p.services, s)
}

func (p *Package) anonymousStructID(s *types.Struct) string {
	key := s.String()

	if _, ok := p.anonymousStructs[key]; !ok {
		p.anonymousStructs[key] = "$anon" + strconv.Itoa(len(p.anonymousStructs)+1)
	}
	return p.anonymousStructs[key]
}

// Credit: https://stackoverflow.com/a/70999797/3140799
func (p *Package) constantsOf(t *types.Named) (values map[string]*types.Const) {
	values = make(map[string]*types.Const)

	for _, file := range p.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, specs := range genDecl.Specs {
				valueSpec, ok := specs.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range valueSpec.Names {
					c := p.TypesInfo.ObjectOf(name).(*types.Const)
					if strings.HasSuffix(c.Type().String(), t.Obj().Name()) {
						values[name.Name] = c
					}
				}
			}
		}
	}
	return
}

type Stats struct {
	NumPackages int
	NumStructs  int
	NumMethods  int
	NumEnums    int
	NumModels   int
	StartTime   time.Time
	EndTime     time.Time
}

type Project struct {
	pkgs  []*Package
	Stats Stats
}

func ParseProject(patterns []string, options *flags.GenerateBindingsOptions) (*Project, error) {
	buildFlags, err := options.BuildFlags()
	if err != nil {
		return nil, err
	}

	pkgs, err := LoadPackages(buildFlags, true,
		append(patterns, WailsAppPkgPath)...,
	)
	if err != nil {
		return nil, err
	}

	services, err := Services(pkgs)
	if err != nil {
		return nil, err
	}

	return &Project{
		pkgs: BuildPackages(pkgs, services),
	}, nil
}

func Services(pkgs []*packages.Package) (services []*Service, err error) {
	var app *packages.Package
	otherPkgs := append(make([]*packages.Package, 0, len(pkgs)), pkgs...)
	if index := slices.IndexFunc(pkgs, func(pkg *packages.Package) bool { return pkg.PkgPath == WailsAppPkgPath }); index >= 0 {
		app = pkgs[index]
		otherPkgs = slices.Delete(otherPkgs, index, index+1)
	}

	if app == nil {
		err = errors.New("LoadPackages() did not load the application package")
		return
	}

	found, err := FindServices(app, otherPkgs)
	if err != nil {
		return
	}

	for _, service := range found {
		services = append(services, &Service{service})
	}
	return
}

func RelativeBindingsDir(base *types.Package, target *types.Package) string {
	if base == target {
		return "."
	}

	basePath := base.Path()
	if base.Name() == "main" {
		basePath = filepath.Join(basePath, "main")
	}

	targetPath := target.Path()
	if target.Name() == "main" {
		targetPath = filepath.Join(targetPath, "main")
	}

	relativePath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		panic(err)
	}

	return filepath.ToSlash(relativePath)
}

var reservedWords = []string{
	"abstract",
	"arguments",
	"await",
	"boolean",
	"break",
	"byte",
	"case",
	"catch",
	"char",
	"class",
	"const",
	"continue",
	"debugger",
	"default",
	"delete",
	"do",
	"double",
	"else",
	"enum",
	"eval",
	"export",
	"extends",
	"false",
	"final",
	"finally",
	"float",
	"for",
	"function",
	"goto",
	"if",
	"implements",
	"import",
	"in",
	"instanceof",
	"int",
	"interface",
	"let",
	"long",
	"native",
	"new",
	"null",
	"package",
	"private",
	"protected",
	"public",
	"return",
	"short",
	"static",
	"super",
	"switch",
	"synchronized",
	"this",
	"throw",
	"throws",
	"transient",
	"true",
	"try",
	"typeof",
	"var",
	"void",
	"volatile",
	"while",
	"with",
	"yield",
	"object",
}
