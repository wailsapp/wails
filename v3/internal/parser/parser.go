package parser

import (
	"fmt"
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
)

type packagePath = string
type structName = string

type StructDef struct {
	Name       string
	DocComment string
	Fields     []*Field
}

type ParameterType struct {
	Name      string
	IsStruct  bool
	IsSlice   bool
	IsPointer bool
	MapKey    *ParameterType
	MapValue  *ParameterType
	Package   string
}

type Parameter struct {
	Name string
	Type *ParameterType
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
	if p.Type.IsStruct {
		if p.Type.Package != "" {
			parts := strings.Split(p.Type.Package, "/")
			typeName = parts[len(parts)-1] + "." + typeName
			// TODO: Check if this is a duplicate package name
		}
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
	Name       string
	DocComment string
	Inputs     []*Parameter
	Outputs    []*Parameter
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
		return fmt.Sprintf("source['%s']", f.JSName())
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
	if f.Type.Package == "" || f.Type.Package == pkg {
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

type ParsedPackage struct {
	Pkg         *ast.Package
	Name        string
	Path        string
	Dir         string
	StructCache map[structName]*StructDef
}

type Project struct {
	packageCache             map[string]*ParsedPackage
	Path                     string
	BoundMethods             map[packagePath]map[structName][]*BoundMethod
	Models                   map[packagePath]map[structName]*StructDef
	anonymousStructIDCounter int
}

func ParseProject(projectPath string) (*Project, error) {
	result := &Project{
		BoundMethods: make(map[packagePath]map[structName][]*BoundMethod),
		packageCache: make(map[string]*ParsedPackage),
	}
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

func GenerateBindingsAndModels(projectDir string, outputDir string) error {
	p, err := ParseProject(projectDir)
	if err != nil {
		return err
	}

	if p.BoundMethods == nil {
		return nil
	}
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}
	generatedMethods := GenerateBindings(p.BoundMethods)
	for pkg, text := range generatedMethods {
		// Write the file
		err = os.WriteFile(filepath.Join(outputDir, "bindings_"+pkg+".js"), []byte(text), 0644)
		if err != nil {
			return err
		}
	}

	// Generate Models
	if len(p.Models) > 0 {
		generatedModels, err := GenerateModels(p.Models)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(outputDir, "models.ts"), []byte(generatedModels), 0644)
		if err != nil {
			return err
		}
	}

	absPath, err := filepath.Abs(projectDir)
	if err != nil {
		return err
	}
	println("Generated bindings and models for project: " + absPath)
	absPath, err = filepath.Abs(outputDir)
	if err != nil {
		return err
	}
	println("Output directory: " + absPath)

	return nil
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
	for packageName, pkg := range pkgs {
		parsedPackage := &ParsedPackage{
			Pkg:         pkg,
			Name:        packageName,
			Path:        packageName,
			Dir:         getDirectoryForPackage(pkg),
			StructCache: make(map[structName]*StructDef),
		}
		p.packageCache[packageName] = parsedPackage
		result[packageName] = parsedPackage
	}
	return result, nil
}

func (p *Project) findApplicationNewCalls(pkgs map[string]*ParsedPackage) (err error) {

	var callFound bool

	for _, pkg := range pkgs {
		thisPackage := pkg.Pkg
		// Iterate through the package's files
		for _, file := range thisPackage.Files {
			// Use an ast.Inspector to find the calls to application.New
			ast.Inspect(file, func(n ast.Node) bool {
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
					if _, ok := arrayType.Elt.(*ast.InterfaceType); !ok {
						continue
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
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil {
				// Check if the function has a receiver argument of the struct type
				recvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
				if ok {
					if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == name {
						// Add the method to the list of methods
						method := &BoundMethod{
							Name:       funcDecl.Name.Name,
							DocComment: funcDecl.Doc.Text(),
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
	}
	p.addBoundMethods(pkg.Path, name, methods)
	return nil
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
	result := &ParameterType{}
	result.Name = getTypeString(field.Type)
	switch t := field.Type.(type) {
	case *ast.Ident:
		result.IsStruct = isStructType(t)
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
										Name:       name,
										DocComment: typeDecl.Doc.Text(),
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
				if paramType.Package == "" {
					paramType.Package = pkg.Path
				}
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
				}
				p.packageCache[path] = result
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
	impPkg, err := parser.ParseDir(token.NewFileSet(), packagedir, nil, parser.AllErrors)
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
		return e.Obj != nil && e.Obj.Kind == ast.Typ
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
