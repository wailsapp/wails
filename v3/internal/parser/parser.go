package parser

import (
	"fmt"
	"github.com/wailsapp/wails/v3/internal/flags"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/hash"
)

type packagePath = string
type structName = string

type StructDef struct {
	Name       string
	DocComment string
	Fields     []*Field
}

func (s *StructDef) DefaultValueList() string {
	var allFields []string
	for _, field := range s.Fields {
		thisFieldWithDefaultValue := fmt.Sprintf("%s = %s", field.JSName(), field.DefaultValue())
		allFields = append(allFields, thisFieldWithDefaultValue)
	}
	return strings.Join(allFields, ", ")
}

type ParameterType struct {
	Name      string
	IsStruct  bool
	IsSlice   bool
	IsPointer bool
	IsEnum    bool
	MapKey    *ParameterType
	MapValue  *ParameterType
	Package   string
}

type EnumDef struct {
	Name        string
	Filename    string
	DocComments []string
	Values      []*EnumValue
}

type EnumValue struct {
	Name        string
	Value       string
	DocComments []string
}

type Parameter struct {
	Name string
	Type *ParameterType
}

func (p *Parameter) NamespacedStructType() string {
	var typeName string
	if p.Type.Package != "" {
		parts := strings.Split(p.Type.Package, "/")
		typeName = parts[len(parts)-1] + "."
	}
	return typeName + p.Type.Name
}
func (p *Parameter) NamespacedStructVariable() string {
	var typeName string
	if p.Type.Package != "" {
		parts := strings.Split(p.Type.Package, "/")
		typeName = parts[len(parts)-1]
	}
	return typeName + p.Type.Name
}

func (p *Parameter) JSType() string {
	// Convert type to javascript equivalent type
	var typeName string
	switch p.Type.Name {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
		typeName = "number"
	case "string":
		typeName = "string"
	case "bool":
		typeName = "boolean"
	default:
		typeName = p.Type.Name
	}

	// if the type is a struct, we need to add the package name
	if p.Type.IsStruct || p.Type.IsEnum {
		typeName = p.NamespacedStructType()
		typeName = strings.ReplaceAll(typeName, ".", "")
	}

	// Add slice suffix
	if p.Type.IsSlice {
		typeName += "[]"
	}

	// Add pointer suffix
	if p.Type.IsPointer {
		typeName += " | null"
	}

	return typeName
}

type BoundMethod struct {
	Package    string
	Name       string
	DocComment string
	Inputs     []*Parameter
	Outputs    []*Parameter
	ID         uint32
	Alias      *uint32
}

func (m BoundMethod) IDAsString() string {
	return strconv.Itoa(int(m.ID))
}

type Field struct {
	Name string
	Type *ParameterType
}

func (f *Field) JSName() string {
	return strings.ToLower(f.Name[0:1]) + f.Name[1:]
}

// TSBuild contains the typescript to build a field for a JS object
// via assignment for simple types or constructors for structs
func (f *Field) TSBuild(pkg string) string {
	if !f.Type.IsStruct {
		return fmt.Sprintf("safeSource['%s']", f.JSName())
	}

	if f.Type.Package == "" || f.Type.Package == pkg {
		return fmt.Sprintf("%s.createFrom(source['%s'])", f.Type.Name, f.JSName())
	}

	return fmt.Sprintf("%s.%s.createFrom(source['%s'])", pkgAlias(f.Type.Package), f.Type.Name, f.JSName())
}

func (f *Field) JSDef(pkg string) string {
	var jsType string
	switch f.Type.Name {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
		jsType = "number"
	case "string":
		jsType = "string"
	case "bool":
		jsType = "boolean"
	default:
		jsType = f.Type.Name
	}

	var result string
	isExternalStruct := f.Type.Package != "" && f.Type.Package != pkg && f.Type.IsStruct
	if f.Type.Package == "" || f.Type.Package == pkg || !isExternalStruct {
		result += fmt.Sprintf("%s: %s;", f.JSName(), jsType)
	} else {
		parts := strings.Split(f.Type.Package, "/")
		result += fmt.Sprintf("%s: %s.%s;", f.JSName(), parts[len(parts)-1], jsType)
	}

	if !ast.IsExported(f.Name) {
		result += " // Warning: this is unexported in the Go struct."
	}

	return result
}

func (f *Field) JSDocType(pkg string) string {
	var jsType string
	switch f.Type.Name {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
		jsType = "number"
	case "string":
		jsType = "string"
	case "bool":
		jsType = "boolean"
	default:
		jsType = f.Type.Name
	}

	var result string
	isExternalStruct := f.Type.Package != "" && f.Type.Package != pkg && f.Type.IsStruct
	if f.Type.Package == "" || f.Type.Package == pkg || !isExternalStruct {
		if f.Type.IsStruct || f.Type.IsEnum {
			result = fmt.Sprintf("%s.%s", pkg, jsType)
		} else {
			result = jsType
		}
	} else {
		parts := strings.Split(f.Type.Package, "/")
		result += fmt.Sprintf("%s.%s", parts[len(parts)-1], jsType)
	}

	if !ast.IsExported(f.Name) {
		result += " // Warning: this is unexported in the Go struct."
	}

	return result
}

func (f *Field) DefaultValue() string {
	// Return the default value of the typescript version of the type as a string
	switch f.Type.Name {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uintptr", "float32", "float64", "uint64":
		return "0"
	case "string":
		return `""`
	case "bool":
		return "false"
	default:
		return "null"
	}
}

type ConstDef struct {
	Name       string
	DocComment string
	Value      string
}

type TypeDef struct {
	Name           string
	DocComments    []string
	Type           string
	Consts         []*ConstDef
	ShouldGenerate bool
}

func (t *TypeDef) GeneratedName() string {
	return t.Name + "Enum"
}

type ParsedPackage struct {
	Pkg         *ast.Package
	Name        string
	Path        string
	Dir         string
	StructCache map[structName]*StructDef
	TypeCache   map[string]*TypeDef
}

type Project struct {
	packageCache             map[string]*ParsedPackage
	Path                     string
	BoundMethods             map[packagePath]map[structName][]*BoundMethod
	Models                   map[packagePath]map[structName]*StructDef
	Types                    map[packagePath]map[structName]*TypeDef
	anonymousStructIDCounter int
	Stats                    Stats
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

func ParseProject(projectPath string) (*Project, error) {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, err
	}
	result := &Project{
		Path:         absPath,
		BoundMethods: make(map[packagePath]map[structName][]*BoundMethod),
		packageCache: make(map[string]*ParsedPackage),
	}
	result.Stats.StartTime = time.Now()
	pkgs, err := result.parseDirectory(projectPath)
	if err != nil {
		return nil, err
	}
	err = result.findApplicationNewCalls(pkgs)
	if err != nil {
		return nil, err
	}
	for _, pkg := range result.packageCache {
		if len(pkg.StructCache) > 0 {
			if result.Models == nil {
				result.Models = make(map[packagePath]map[structName]*StructDef)
			}
			result.Models[pkg.Path] = pkg.StructCache
		}
	}
	return result, nil
}

func GenerateBindingsAndModels(options *flags.GenerateBindingsOptions) error {
	p, err := ParseProject(options.ProjectDirectory)
	if err != nil {
		return err
	}

	if p.BoundMethods == nil {
		return nil
	}
	err = os.MkdirAll(options.OutputDirectory, 0755)
	if err != nil {
		return err
	}
	p.Stats.NumMethods = len(p.BoundMethods)
	generatedMethods := GenerateBindings(p.BoundMethods, options.UseIDs)
	for pkg, text := range generatedMethods {
		// Write the file
		err = os.WriteFile(filepath.Join(options.OutputDirectory, pkg+".js"), []byte(text), 0644)
		if err != nil {
			return err
		}
	}

	p.Stats.NumModels = len(p.Models)
	p.Stats.NumEnums = len(p.Types)

	// Generate Models
	if len(p.Models) > 0 {
		generatedModels, err := GenerateModels(p.Models, p.Types, options)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(options.OutputDirectory, options.ModelsFilename), []byte(generatedModels), 0644)
		if err != nil {
			return err
		}
	}

	p.Stats.EndTime = time.Now()

	absPath, err := filepath.Abs(options.ProjectDirectory)
	if err != nil {
		return err
	}

	fmt.Printf("Processed: %s, %s, %s, %s, %s in %s.\n",
		pluralise(p.Stats.NumPackages, "Package"),
		pluralise(p.Stats.NumStructs, "Struct"),
		pluralise(p.Stats.NumMethods, "Method"),
		pluralise(p.Stats.NumEnums, "Enum"),
		pluralise(p.Stats.NumModels, "Model"),
		p.Stats.EndTime.Sub(p.Stats.StartTime).String())

	absPath, err = filepath.Abs(options.OutputDirectory)
	if err != nil {
		return err
	}
	println("Output directory: " + absPath)

	return nil
}

func pluralise(number int, word string) string {
	if number == 1 {
		return fmt.Sprintf("%d %s", number, word)
	}
	return fmt.Sprintf("%d %ss", number, word)
}

func (p *Project) parseDirectory(dir string) (map[string]*ParsedPackage, error) {
	if p.packageCache[dir] != nil {
		return map[string]*ParsedPackage{dir: p.packageCache[dir]}, nil
	}
	// Parse the directory
	fset := token.NewFileSet()
	if dir == "." || dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = cwd
	}
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]*ParsedPackage)
	p.Stats.NumPackages = len(pkgs)
	for packageName, pkg := range pkgs {
		parsedPackage := &ParsedPackage{
			Pkg:         pkg,
			Name:        packageName,
			Path:        packageName,
			Dir:         getDirectoryForPackage(pkg),
			StructCache: make(map[structName]*StructDef),
			TypeCache:   make(map[string]*TypeDef),
		}
		p.parseTypes(map[string]*ParsedPackage{packageName: parsedPackage})
		p.packageCache[packageName] = parsedPackage
		result[packageName] = parsedPackage
	}
	return result, nil
}

func (p *Project) findApplicationNewCalls(pkgs map[string]*ParsedPackage) (err error) {

	var callFound bool

	p.parseTypes(pkgs)

	for _, pkg := range pkgs {
		thisPackage := pkg.Pkg
		// Iterate through the package's files
		for _, file := range thisPackage.Files {
			// Use an ast.Inspector to find the calls to application.New
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for const declaration
				genDecl, ok := n.(*ast.GenDecl)
				if ok {
					switch genDecl.Tok {
					case token.TYPE:
						var comments []string
						if genDecl.Doc != nil {
							comments = CommentGroupToText(genDecl.Doc)
						}
						for _, spec := range genDecl.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								p.parseTypeDeclaration(typeSpec, pkg, comments)
							}
						}
					case token.CONST:
						p.parseConstDeclaration(genDecl, pkg)
					default:
					}
					return true

				}

				// Check if the node is a call expression
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				// Check if the function being called is "application.New"
				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if selExpr.Sel.Name != "New" {
					return true
				}
				if id, ok := selExpr.X.(*ast.Ident); !ok || id.Name != "application" {
					return true
				}

				// Check there is only 1 argument
				if len(callExpr.Args) != 1 {
					return true
				}

				// Check argument 1 is a struct literal
				structLit, ok := callExpr.Args[0].(*ast.CompositeLit)
				if !ok {
					return true
				}

				// Check struct literal is of type "application.Options"
				selectorExpr, ok := structLit.Type.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if selectorExpr.Sel.Name != "Options" {
					return true
				}
				if id, ok := selectorExpr.X.(*ast.Ident); !ok || id.Name != "application" {
					return true
				}

				for _, elt := range structLit.Elts {
					// Find the "Bind" field
					kvExpr, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					if id, ok := kvExpr.Key.(*ast.Ident); !ok || id.Name != "Bind" {
						continue
					}
					// Check the value is a slice of interfaces
					sliceExpr, ok := kvExpr.Value.(*ast.CompositeLit)
					if !ok {
						continue
					}
					var arrayType *ast.ArrayType
					if arrayType, ok = sliceExpr.Type.(*ast.ArrayType); !ok {
						continue
					}

					// Check array type is of type "interface{}"
					_, isInterfaceType := arrayType.Elt.(*ast.InterfaceType)
					if !isInterfaceType {
						// Check it's an "any" type
						ident, isAnyType := arrayType.Elt.(*ast.Ident)
						if !isAnyType {
							continue
						}
						if ident.Name != "any" {
							continue
						}
					}
					callFound = true
					// Iterate through the slice elements
					for _, elt := range sliceExpr.Elts {
						result, shouldContinue := p.parseBoundExpression(elt, pkg)
						if shouldContinue {
							continue
						}
						return result

					}
				}

				return true
			})
		}
		p.addTypes(pkg.Path, pkg.TypeCache)
		if !callFound {
			return fmt.Errorf("no Bound structs found")
		}
	}
	return nil
}

func (p *Project) parseBoundUnaryExpression(unaryExpr *ast.UnaryExpr, pkg *ParsedPackage) (bool, bool) {
	// Check the unary expression is a composite lit

	switch t := unaryExpr.X.(type) {
	case *ast.CompositeLit:
		return p.parseBoundCompositeLit(t, pkg)
	}
	return false, true

}

func (p *Project) addBoundMethods(packagePath string, name string, boundMethods []*BoundMethod) {
	_, ok := p.BoundMethods[packagePath]
	if !ok {
		p.BoundMethods[packagePath] = make(map[structName][]*BoundMethod)
	}
	p.BoundMethods[packagePath][name] = boundMethods
}

func (p *Project) parseBoundStructMethods(name string, pkg *ParsedPackage) error {
	var methods []*BoundMethod
	// Iterate over all files in the package
	for _, file := range pkg.Pkg.Files {
		// Iterate over all declarations in the file
		for _, decl := range file.Decls {
			// Check if the declaration is a type declaration
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil && funcDecl.Name.IsExported() {
				var ident *ast.Ident
				var ok bool

				switch funcDecl.Recv.List[0].Type.(type) {
				case *ast.StarExpr:
					recv := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
					ident, ok = recv.X.(*ast.Ident)
				case *ast.Ident:
					ident, ok = funcDecl.Recv.List[0].Type.(*ast.Ident)
				}
				if ok && ident.Name == name {
					fqn := fmt.Sprintf("%s.%s.%s", pkg.Path, name, funcDecl.Name.Name)

					var alias *uint32
					var err error
					// Check for the text `wails:methodID <integer>`
					if funcDecl.Doc != nil {
						for _, docstring := range funcDecl.Doc.List {
							if strings.Contains(docstring.Text, "//wails:methodID") {
								idString := strings.TrimSpace(strings.TrimPrefix(docstring.Text, "//wails:methodID"))
								parsedID, err := strconv.ParseUint(idString, 10, 32)
								if err != nil {
									return fmt.Errorf("invalid value in `wails:methodID` directive: '%s'. Expected a valid uint32 value", idString)
								}
								alias = lo.ToPtr(uint32(parsedID))
								break
							}
						}
					}
					id, err := hash.Fnv(fqn)
					if err != nil {
						return err
					}

					method := &BoundMethod{
						Package:    pkg.Path,
						ID:         id,
						Name:       funcDecl.Name.Name,
						DocComment: strings.TrimSpace(funcDecl.Doc.Text()),
						Alias:      alias,
					}

					if funcDecl.Type.Params != nil {
						method.Inputs = p.parseParameters(funcDecl.Type.Params, pkg)
					}
					if funcDecl.Type.Results != nil {
						method.Outputs = p.parseParameters(funcDecl.Type.Results, pkg)
					}

					methods = append(methods, method)
				}

			}
		}
	}
	p.addBoundMethods(pkg.Path, name, methods)
	return nil
}

func (p *Project) addTypes(packagePath string, types map[string]*TypeDef) {
	if types == nil || len(types) == 0 {
		return
	}
	if p.Types == nil {
		p.Types = make(map[string]map[string]*TypeDef)
	}
	_, ok := p.Types[packagePath]
	if !ok {
		p.Types[packagePath] = make(map[string]*TypeDef)
	}
	p.Types[packagePath] = types
}

func (p *Project) parseParameters(params *ast.FieldList, pkg *ParsedPackage) []*Parameter {
	var result []*Parameter
	for _, field := range params.List {
		var theseFields []*Parameter
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				theseFields = append(theseFields, &Parameter{
					Name: name.Name,
				})
			}
		} else {
			theseFields = append(theseFields, &Parameter{
				Name: "",
			})
		}
		// loop over fields
		for _, thisField := range theseFields {
			thisField.Type = p.parseParameterType(field, pkg)
			result = append(result, thisField)
		}
	}
	return result
}

func (p *Project) parseParameterType(field *ast.Field, pkg *ParsedPackage) *ParameterType {
	result := &ParameterType{
		Package: pkg.Path,
	}
	result.Name = getTypeString(field.Type)
	switch t := field.Type.(type) {
	case *ast.Ident:
		result.IsStruct = isStructType(t)
		if !result.IsStruct {
			// Check if it's a type alias
			typeDef, ok := pkg.TypeCache[t.Name]
			if ok {
				typeDef.ShouldGenerate = true
				result.IsEnum = true
			}
		}
	case *ast.StarExpr:
		result = p.parseParameterType(&ast.Field{Type: t.X}, pkg)
		result.IsPointer = true
	case *ast.StructType:
		result.IsStruct = true
		if result.Name == "" {
			// Anonymous struct
			result.Name = p.anonymousStructID()
			// Create a new struct definition
			result := &StructDef{
				Name: result.Name,
			}
			pkg.StructCache[result.Name] = result
			// Parse the fields
			result.Fields = p.parseStructFields(&ast.StructType{
				Fields: t.Fields,
			}, pkg)
			_ = result
		}
	case *ast.SelectorExpr:
		extPackage, err := p.getParsedPackageFromName(t.X.(*ast.Ident).Name, pkg)
		if err != nil {
			log.Fatal(err)
		}
		result.IsStruct = p.getStructDef(t.Sel.Name, extPackage)
		if !result.IsStruct {
			// Check if it's a type alias
			typeDef, ok := extPackage.TypeCache[t.Sel.Name]
			if ok {
				typeDef.ShouldGenerate = true
				result.IsEnum = true
			}
		}
		result.Package = extPackage.Path
	case *ast.ArrayType:
		result.IsSlice = true
		result.IsStruct = isStructType(t.Elt)
	case *ast.MapType:
		tempfield := &ast.Field{Type: t.Key}
		result.MapKey = p.parseParameterType(tempfield, pkg)
		tempfield.Type = t.Value
		result.MapValue = p.parseParameterType(tempfield, pkg)
	default:
	}
	if result.IsStruct {
		p.getStructDef(result.Name, pkg)
		if result.Package == "" {
			result.Package = pkg.Path
		}
	}
	return result
}

func (p *Project) getStructDef(name string, pkg *ParsedPackage) bool {
	_, ok := pkg.StructCache[name]
	if ok {
		return true
	}
	// Iterate over all files in the package
	for _, file := range pkg.Pkg.Files {
		// Iterate over all declarations in the file
		for _, decl := range file.Decls {
			// Check if the declaration is a type declaration
			if typeDecl, ok := decl.(*ast.GenDecl); ok {
				// Check if the type declaration is a struct type
				if typeDecl.Tok == token.TYPE {
					for _, spec := range typeDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if structType, ok := typeSpec.Type.(*ast.StructType); ok {
								if typeSpec.Name.Name == name {
									result := &StructDef{
										Name: name,
										//TODO DocComment: CommentGroupToText(typeDecl.Doc),
									}
									pkg.StructCache[name] = result
									result.Fields = p.parseStructFields(structType, pkg)
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (p *Project) parseStructFields(structType *ast.StructType, pkg *ParsedPackage) []*Field {
	var result []*Field
	for _, field := range structType.Fields.List {
		var theseFields []*Field
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				theseFields = append(theseFields, &Field{
					Name: name.Name,
				})
			}
		} else {
			theseFields = append(theseFields, &Field{
				Name: "",
			})
		}
		// loop over fields
		for _, thisField := range theseFields {
			paramType := p.parseParameterType(field, pkg)
			if paramType.IsStruct {
				_, ok := pkg.StructCache[paramType.Name]
				if !ok {
					p.getStructDef(paramType.Name, pkg)
				}
			}
			if paramType.Package == "" {
				paramType.Package = pkg.Path
			}
			thisField.Type = paramType
			result = append(result, thisField)
		}
	}
	return result
}

func (p *Project) getParsedPackageFromName(packageName string, currentPackage *ParsedPackage) (*ParsedPackage, error) {
	for _, file := range currentPackage.Pkg.Files {
		for _, imp := range file.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return nil, err
			}
			_, lastPathElement := filepath.Split(path)
			if imp.Name != nil && imp.Name.Name == packageName || lastPathElement == packageName {
				// Get the directory for the package
				dir, err := getPackageDir(path)
				if err != nil {
					return nil, err
				}
				pkg, err := p.getPackageFromPath(dir, path)
				if err != nil {
					return nil, err
				}
				result := &ParsedPackage{
					Pkg:         pkg,
					Name:        packageName,
					Path:        path,
					Dir:         dir,
					StructCache: make(map[string]*StructDef),
					TypeCache:   make(map[string]*TypeDef),
				}
				p.packageCache[path] = result

				// Parse types
				p.parseTypes(map[string]*ParsedPackage{path: result})

				return result, nil
			}
		}
	}
	return nil, fmt.Errorf("package %s not found in %s", packageName, currentPackage.Name)
}

func getPackageDir(importPath string) (string, error) {
	pkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

func (p *Project) getPackageFromPath(packagedir string, packagepath string) (*ast.Package, error) {
	impPkg, err := parser.ParseDir(token.NewFileSet(), packagedir, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for impName, impPkg := range impPkg {
		if impName == "main" {
			continue
		}
		return impPkg, nil
	}
	return nil, fmt.Errorf("package not found in imported package %s", packagepath)
}

func (p *Project) anonymousStructID() string {
	p.anonymousStructIDCounter++
	return fmt.Sprintf("anon%d", p.anonymousStructIDCounter)
}

func (p *Project) parseBoundExpression(elt ast.Expr, pkg *ParsedPackage) (bool, bool) {

	switch t := elt.(type) {
	case *ast.UnaryExpr:
		return p.parseBoundUnaryExpression(t, pkg)
	case *ast.Ident:
		return p.parseBoundIdent(t, pkg)
	case *ast.CallExpr:
		return p.parseBoundCallExpression(t, pkg)
	default:
		println("unhandled expression type", reflect.TypeOf(t).String())
	}

	return false, false
}

func (p *Project) parseBoundIdent(ident *ast.Ident, pkg *ParsedPackage) (bool, bool) {
	if ident.Obj == nil {
		return false, true
	}
	switch t := ident.Obj.Decl.(type) {
	//case *ast.StructType:
	//	return p.parseBoundStruct(t, pkg)
	case *ast.TypeSpec:
		return p.parseBoundTypeSpec(t, pkg)
	case *ast.AssignStmt:
		return p.parseBoundAssignment(t, pkg)
	default:
		println("unhandled ident type", reflect.TypeOf(t).String())
	}
	return false, false
}

func (p *Project) parseBoundAssignment(assign *ast.AssignStmt, pkg *ParsedPackage) (bool, bool) {
	return p.parseBoundExpression(assign.Rhs[0], pkg)
}

func (p *Project) parseBoundCompositeLit(lit *ast.CompositeLit, pkg *ParsedPackage) (bool, bool) {

	switch t := lit.Type.(type) {
	case *ast.StructType:
		//return p.parseBoundStructType(t, pkg)
		return false, true
	case *ast.Ident:
		err := p.parseBoundStructMethods(t.Name, pkg)
		if err != nil {
			return true, false
		}
		return false, true
	case *ast.SelectorExpr:
		return p.parseBoundSelectorExpression(t, pkg)
	}

	return false, true
}

func (p *Project) parseBoundSelectorExpression(selector *ast.SelectorExpr, pkg *ParsedPackage) (bool, bool) {

	switch t := selector.X.(type) {
	case *ast.Ident:
		// Look up the package
		var parsedPackage *ParsedPackage
		parsedPackage, err := p.getParsedPackageFromName(t.Name, pkg)
		if err != nil {
			return true, false
		}
		err = p.parseBoundStructMethods(selector.Sel.Name, parsedPackage)
		if err != nil {
			return true, false
		}
		return false, true
	default:
		println("unhandled selector type", reflect.TypeOf(t).String())
	}
	return false, true
}

func (p *Project) parseBoundCallExpression(callExpr *ast.CallExpr, pkg *ParsedPackage) (bool, bool) {

	// Check if this call returns a struct pointer
	switch t := callExpr.Fun.(type) {
	case *ast.Ident:
		if t.Obj == nil {
			return false, true
		}
		switch t := t.Obj.Decl.(type) {
		case *ast.FuncDecl:
			return p.parseBoundFuncDecl(t, pkg)
		}

	case *ast.SelectorExpr:
		// Get package for selector
		var parsedPackage *ParsedPackage
		parsedPackage, err := p.getParsedPackageFromName(t.X.(*ast.Ident).Name, pkg)
		if err != nil {
			return true, false
		}
		// Get function from package
		var extFundDecl *ast.FuncDecl
		extFundDecl, err = p.getFunctionFromName(t.Sel.Name, parsedPackage)
		if err != nil {
			return true, false
		}
		return p.parseBoundFuncDecl(extFundDecl, parsedPackage)
	default:
		println("unhandled call type", reflect.TypeOf(t).String())
	}

	return false, true
}

func (p *Project) parseBoundFuncDecl(t *ast.FuncDecl, pkg *ParsedPackage) (bool, bool) {
	if t.Type.Results == nil {
		return false, true
	}
	if len(t.Type.Results.List) != 1 {
		return false, true
	}
	switch t := t.Type.Results.List[0].Type.(type) {
	case *ast.StarExpr:
		return p.parseBoundExpression(t.X, pkg)
	default:
		println("Unhandled funcdecl type", reflect.TypeOf(t).String())
	}
	return false, false
}

func (p *Project) parseBoundTypeSpec(typeSpec *ast.TypeSpec, pkg *ParsedPackage) (bool, bool) {
	switch t := typeSpec.Type.(type) {
	case *ast.StructType:
		err := p.parseBoundStructMethods(typeSpec.Name.Name, pkg)
		if err != nil {
			return true, false
		}
	default:
		println("unhandled type spec type", reflect.TypeOf(t).String())
	}
	return false, true
}

func (p *Project) getFunctionFromName(name string, parsedPackage *ParsedPackage) (*ast.FuncDecl, error) {
	for _, f := range parsedPackage.Pkg.Files {
		for _, decl := range f.Decls {
			switch t := decl.(type) {
			case *ast.FuncDecl:
				if t.Name.Name == name {
					return t, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("function not found")
}

func (p *Project) parseTypeDeclaration(decl *ast.TypeSpec, pkg *ParsedPackage, comments []string) {
	switch t := decl.Type.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64",
			"uintptr", "float32", "float64", "string", "bool":
			// Store this in the type cache
			pkg.TypeCache[decl.Name.Name] = &TypeDef{
				Name:        decl.Name.Name,
				Type:        t.Name,
				DocComments: TrimSlice(comments),
			}
		}
	}
}

func (p *Project) parseConstDeclaration(decl *ast.GenDecl, pkg *ParsedPackage) {
	// Check if the type of the constant is in the type cache and if it doesn't exist, return
	var latestValues []ast.Expr

	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		// Extract the Type
		typeString := getTypeString(valueSpec.Type)
		if typeString == "" {
			continue
		}
		// Check if the type is in the type cache
		typeDecl, ok := pkg.TypeCache[typeString]
		if !ok {
			continue
		}

		// Get the latest values
		if len(valueSpec.Values) > 0 {
			latestValues = valueSpec.Values
		}

		// Iterate over the names
		for index, name := range valueSpec.Names {
			constDecl := &ConstDef{
				Name:  name.Name,
				Value: typeString,
			}

			// Get the value
			if len(latestValues) > 0 {
				switch t := latestValues[index].(type) {
				case *ast.BasicLit:
					constDecl.Value = t.Value
				case *ast.Ident:
					constDecl.Value = t.Name
				}
			}

			if valueSpec.Doc != nil {
				constDecl.DocComment = strings.TrimSpace(valueSpec.Doc.Text())
			}
			typeDecl.Consts = append(typeDecl.Consts, constDecl)
		}
	}
}

func (p *Project) RelativePackageDir(path string) string {
	return strings.TrimPrefix(path, p.Path)
}

func (p *Project) parseTypes(pkgs map[string]*ParsedPackage) {
	for _, pkg := range pkgs {
		thisPackage := pkg.Pkg
		// Iterate through the package's files
		for _, file := range thisPackage.Files {
			// Use an ast.Inspector to find the calls to application.New
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for const declaration
				genDecl, ok := n.(*ast.GenDecl)
				if ok {
					switch genDecl.Tok {
					case token.TYPE:
						var comments []string
						if genDecl.Doc != nil {
							comments = CommentGroupToText(genDecl.Doc)
						}
						for _, spec := range genDecl.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								p.parseTypeDeclaration(typeSpec, pkg, comments)
							}
						}
					case token.CONST:
						p.parseConstDeclaration(genDecl, pkg)
					default:
					}
					return true

				}

				return true
			})
		}
		p.addTypes(pkg.Path, pkg.TypeCache)
	}
}

func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return getTypeString(t.X)
	case *ast.ArrayType:
		return getTypeString(t.Elt)
	case *ast.MapType:
		return "map"
	case *ast.SelectorExpr:
		return getTypeString(t.Sel)
	default:
		return ""
	}
}

func isStructType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.StructType:
		return true
	case *ast.StarExpr:
		return isStructType(e.X)
	case *ast.SelectorExpr:
		return isStructType(e.Sel)
	case *ast.ArrayType:
		return isStructType(e.Elt)
	case *ast.SliceExpr:
		return isStructType(e.X)
	case *ast.Ident:
		if e.Obj != nil && e.Obj.Kind == ast.Typ {
			return isStructType(e.Obj.Decl.(*ast.TypeSpec).Type)
		}
		return false
	default:
		return false
	}
}

func getDirectoryForPackage(pkg *ast.Package) string {
	for filename := range pkg.Files {
		path := filepath.Dir(filename)
		abs, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		return abs
	}
	return ""
}

func CommentGroupToText(comments *ast.CommentGroup) []string {
	if comments == nil {
		return nil
	}
	var result []string
	for _, comment := range comments.List {
		result = append(result, strings.TrimSpace(comment.Text))
	}
	return result
}

func TrimSlice(comments []string) []string {
	for i, comment := range comments {
		comments[i] = strings.TrimSpace(comment)
	}
	return comments
}
