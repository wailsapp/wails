package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

var packageCache = make(map[string]*ParsedPackage)

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
}

type Parameter struct {
	Name string
	Type *ParameterType
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

type ParsedPackage struct {
	Pkg         *ast.Package
	Name        string
	Path        string
	Dir         string
	StructCache map[structName]*StructDef
}

type Project struct {
	Path         string
	BoundMethods map[packagePath]map[structName][]*BoundMethod
}

func ParseProject(projectPath string) (*Project, error) {
	result := &Project{
		BoundMethods: make(map[packagePath]map[structName][]*BoundMethod),
	}
	pkgs, err := result.parseDirectory(projectPath)
	if err != nil {
		return nil, err
	}
	println("Parsed " + projectPath)
	err = result.findApplicationNewCalls(pkgs)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *Project) parseDirectory(dir string) (map[string]*ParsedPackage, error) {
	println("Parsing directory " + dir)
	if packageCache[dir] != nil {
		println("Found directory in cache!")
		return map[string]*ParsedPackage{dir: packageCache[dir]}, nil
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
		packageCache[dir] = parsedPackage
		result[packageName] = parsedPackage
	}
	return result, nil
}

func (p *Project) findApplicationNewCalls(pkgs map[string]*ParsedPackage) (err error) {

	var callFound bool

	for packageName, pkg := range pkgs {
		thisPackage := pkg.Pkg
		println("  - Looking in package: " + packageName)
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
						// Check the element is a unary expression
						unaryExpr, ok := elt.(*ast.UnaryExpr)
						if ok {
							// Check the unary expression is a composite lit
							boundStructLit, ok := unaryExpr.X.(*ast.CompositeLit)
							if !ok {
								continue
							}
							// Check if the composite lit is a struct
							if _, ok := boundStructLit.Type.(*ast.StructType); ok {
								// Parse struct
								continue
							}
							// Check if the lit is an ident
							ident, ok := boundStructLit.Type.(*ast.Ident)
							if ok {
								err = p.parseBoundStructMethods(ident.Name, pkg)
								if err != nil {
									return true
								}
								continue
								// Check if the ident is a struct type
								//if typeSpec, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
								//	var parsedStruct *StructDefinition
								//	parsedStruct, err = p.parseStruct(typeSpec, thisPackage)
								//	if err != nil {
								//		return true
								//	}
								//	p.addModel(thisPackage.Name, parsedStruct)
								//	p.addBoundStruct(thisPackage.Name, ident.Name)
								//	continue
								//}
								//// Check if the ident is a struct type
								//if _, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
								//	thisPackage.boundStructs[ident.Name] = &BoundStruct{
								//		Name: ident.Name,
								//	}
								//	continue
								//}
								// Check the typespec decl is a struct
								//if _, ok := ident.Obj.Decl.(*ast.StructType); ok {
								//	continue
								//}

							}
							// Check if the lit is a selector
							selector, ok := boundStructLit.Type.(*ast.SelectorExpr)
							if ok {
								// Check if the selector is an ident
								if _, ok := selector.X.(*ast.Ident); ok {
									// Look up the package
									var parsedPackage *ParsedPackage
									parsedPackage, err = p.getParsedPackageFromName(selector.X.(*ast.Ident).Name, pkg)
									if err != nil {
										return true
									}
									err = p.parseBoundStructMethods(selector.Sel.Name, parsedPackage)
									if err != nil {
										return true
									}
									continue
								}
								continue
							}
						}
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
	case *ast.StarExpr:
		result = p.parseParameterType(&ast.Field{Type: t.X}, pkg)
		result.IsPointer = true
		result.IsStruct = isStructType(t.X)
	case *ast.StructType:
		result.IsStruct = true
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
		_, ok := pkg.StructCache[result.Name]
		if !ok {
			p.getStructDef(result.Name, pkg)
		}
	}
	return result
}

func (p *Project) getStructDef(name string, pkg *ParsedPackage) {
	_, ok := pkg.StructCache[name]
	if ok {
		return
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
								}
							}
						}
					}
				}
			}
		}
	}
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
				return &ParsedPackage{
					Pkg:         pkg,
					Name:        packageName,
					Path:        path,
					Dir:         dir,
					StructCache: make(map[string]*StructDef),
				}, nil
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
	default:
		return "any"
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
