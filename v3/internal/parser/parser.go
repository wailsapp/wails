package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
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
	return name
}

func DefaultValue(t types.Type, pkg *Package, mDef *ModelDefinitions) string {
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
			if enum, ok := mDef.Enums[x.Obj().Name()]; ok {
				return enum.DefaultValue(t, pkg)
			} else {
				return DefaultValue(y, pkg, mDef)
			}
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
}

func (m *BoundMethod) embedTuple(tuple *types.Tuple) (result []*Parameter) {
	if tuple == nil {
		return
	}

	for i := 0; i < tuple.Len(); i++ {
		result = append(result, &Parameter{
			Var:    tuple.At(i),
			index:  i,
			Parent: m,
		})
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
	Methods []*BoundMethod
}

func BoundMethods(service *types.TypeName) (methods []*BoundMethod) {
	if named, ok := service.Type().(*types.Named); ok {
		for i := 0; i < named.NumMethods(); i++ {
			fn := named.Method(i)
			if !fn.Exported() {
				continue
			}

			// TODO replace with named.Method(i).String() ???
			fqn := fmt.Sprintf("%s.%s.%s", service.Pkg().Name(), service.Name(), fn.Name())

			id, err := hash.Fnv(fqn)
			if err != nil {
				panic("Failed to hash fqn")
			}

			method := &BoundMethod{
				Func: fn,
				FQN:  fqn,
				ID:   id,
			}

			interfaceFound := false
			for param := range method.Models(nil, false) {
				if types.IsInterface(param.Obj().Type()) {
					interfaceFound = true
					pterm.Warning.Printf("can't bind method %v with interface %v\n", fqn, param.Obj().Name())
				}
			}
			if interfaceFound {
				continue
			}

			methods = append(methods, method)
		}
	}
	return
}

type Package struct {
	*packages.Package
	services         []*Service
	anonymousStructs map[string]string
	doc              *Doc
}

func BuildPackages(buildFlags []string, pkgs []*packages.Package, services []*Service) ([]*Package, error) {
	result := make(map[string]*Package)

	// wrap types.Package
	for _, pPkg := range pkgs {
		result[pPkg.Types.Path()] = &Package{
			Package:          pPkg,
			services:         []*Service{},
			anonymousStructs: make(map[string]string),
			doc:              NewDoc(pPkg),
		}
	}

	// helper function to load missing packages
	loadPackage := func(pkgPath string, services []*Service) (*Package, error) {
		pPkg, err := LoadPackage(buildFlags, true, pkgPath)
		if err != nil {
			return nil, err
		}

		return &Package{
			Package:          pPkg,
			services:         services,
			anonymousStructs: make(map[string]string),
			doc:              NewDoc(pPkg),
		}, nil
	}

	// add services to packages
	for _, service := range services {
		if pkg, ok := result[service.Pkg().Path()]; ok {
			pkg.addService(service)
		} else {
			// load missing packages of service
			pkg, err := loadPackage(service.Pkg().Path(), []*Service{service})
			if err != nil {
				return nil, err
			}
			result[service.Pkg().Path()] = pkg
		}
	}

	// load missing packages of models
	allModels := []*types.Named{}
	for _, pkg := range result {
		allModels = append(allModels, lo.Keys(pkg.Models())...)
	}
	for _, model := range allModels {
		pkgPath := model.Obj().Pkg().Path()
		if _, ok := result[pkgPath]; !ok {
			pkg, err := loadPackage(pkgPath, []*Service{})
			if err != nil {
				return nil, err
			}
			result[pkgPath] = pkg
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
func (p *Package) constantsOf(t *types.Named) []*ConstDef {
	values := []*ConstDef{}

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
						values = append(values, &ConstDef{Name: name.Name, Const: c})
					}
				}
			}
		}
	}
	return values
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
	pkgs    []*Package
	main    *packages.Package
	options *flags.GenerateBindingsOptions
	Stats   Stats
}

func ParseProject(options *flags.GenerateBindingsOptions) (*Project, error) {
	startTime := time.Now()

	buildFlags, err := options.BuildFlags()
	if err != nil {
		return nil, err
	}

	pPkgs, err := LoadPackages(buildFlags, true,
		options.ProjectDirectory, WailsAppPkgPath,
	)
	if err != nil {
		return nil, err
	}
	if n := packages.PrintErrors(pPkgs); n > 0 {
		return nil, errors.New("error while loading packages")
	}

	mainIndex := slices.IndexFunc(pPkgs, func(pkg *packages.Package) bool { return pkg.Name == "main" })
	if mainIndex == -1 {
		return nil, errors.New("application.New() must be inside main package")
	}

	services, err := Services(pPkgs)
	if err != nil {
		return nil, err
	}

	pkgs, err := BuildPackages(buildFlags, pPkgs, services)
	if err != nil {
		return nil, err
	}

	return &Project{
		pkgs:    pkgs,
		main:    pPkgs[mainIndex],
		options: options,
		Stats: Stats{
			StartTime:   startTime,
			NumPackages: len(pkgs),
		},
	}, nil
}

func GenerateBindingsAndModels(options *flags.GenerateBindingsOptions) (*Project, error) {
	p, err := ParseProject(options)
	if err != nil {
		return p, err
	}

	if NumMethods := len(p.BoundMethods()); NumMethods == 0 {
		return p, nil
	} else {
		p.Stats.NumMethods += NumMethods
	}

	err = os.MkdirAll(options.OutputDirectory, 0755)
	if err != nil {
		return p, err
	}

	generatedMethods, err := p.GenerateBindings()
	if err != nil {
		return p, err
	}
	for pkgDir, structs := range generatedMethods {
		// Write the directory
		err = os.MkdirAll(filepath.Join(options.OutputDirectory, pkgDir), 0755)
		if err != nil && !os.IsExist(err) {
			return p, err
		}
		// Write the files
		for structName, text := range structs {
			p.Stats.NumStructs++
			var filename string
			if options.TS {
				filename = structName + ".ts"
			} else {
				filename = structName + ".js"
			}
			err = os.WriteFile(filepath.Join(options.OutputDirectory, pkgDir, filename), []byte(text), 0644)
			if err != nil {
				return p, err
			}
		}
	}

	generatedModels, err := p.GenerateModels()
	if err != nil {
		return p, err
	}
	for pkgDir, text := range generatedModels {
		// Write the directory
		err = os.MkdirAll(filepath.Join(options.OutputDirectory, pkgDir), 0755)
		if err != nil && !os.IsExist(err) {
			return p, err
		}
		// Write the file
		var filename string
		if options.TS {
			filename = options.ModelsFilename + ".ts"
		} else {
			filename = options.ModelsFilename + ".js"
		}
		err = os.WriteFile(filepath.Join(options.OutputDirectory, pkgDir, filename), []byte(text), 0644)
	}
	if err != nil {
		return p, err
	}

	p.Stats.EndTime = time.Now()

	return p, nil
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
		services = append(services, &Service{
			TypeName: service,
			Methods:  BoundMethods(service),
		})
	}
	return
}

func (p *Project) PackageDir(pkg *types.Package) string {
	root := p.main.Types.Path()
	if pkg.Path() == root {
		return "main"
	}

	if strings.HasPrefix(pkg.Path(), root) {
		path, err := filepath.Rel(root, pkg.Path())
		if err != nil {
			panic(err)
		}
		return filepath.ToSlash(path)
	}
	return strings.ReplaceAll(pkg.Path(), "/", "-")
}

func (p *Project) RelativePackageDir(base *types.Package, target *types.Package) string {
	if base == target {
		return "."
	}

	basePath := p.PackageDir(base)
	targetPath := p.PackageDir(target)

	relativePath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		panic(err)
	}

	return filepath.ToSlash(relativePath)
}

func (p *Project) BoundMethods() []*BoundMethod {
	methods := []*BoundMethod{}

	for _, pkg := range p.pkgs {
		for _, service := range pkg.services {
			methods = append(methods, service.Methods...)
		}
	}
	return methods
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
