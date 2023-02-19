package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

var packageCache = make(map[string]*ParsedPackage)

type packageName = string
type structName = string

type Parameter struct {
	Name      string
	Type      string
	IsStruct  bool
	IsSlice   bool
	IsPointer bool
}

type BoundMethod struct {
	Name       string
	DocComment string
	Inputs     []*Parameter
	Outputs    []*Parameter
}

type Model struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name     string
	Type     string
	IsStruct bool
	IsSlice  bool
}

type ParsedPackage struct {
	Pkg *ast.Package
}

type Project struct {
	Path         string
	BoundMethods map[packageName]map[structName][]*BoundMethod
	Models       map[packageName]map[structName]*Model
}

func ParseProject(projectPath string) (*Project, error) {
	result := &Project{
		BoundMethods: make(map[packageName]map[structName][]*BoundMethod),
		Models:       make(map[packageName]map[structName]*Model),
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
		parsedPackage := &ParsedPackage{Pkg: pkg}
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
								err = p.parseBoundStructMethods(ident.Name, thisPackage)
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
									//// Check if the ident is a package
									//if _, ok := context.packages[ident.Name]; !ok {
									//	externalPackage, err := context.getPackageFromPath(ident.Name)
									//	if err != nil {
									//		println("Error getting package from path: " + err.Error())
									//		return true
									//	}
									//	context.packages[ident.Name] = &parsedPackage{
									//		name:         ident.Name,
									//		pkg:          externalPackage,
									//		boundStructs: make(map[string]*BoundStruct),
									//	}
									//}
									//context.packages[ident.Name].boundStructs[selector.Sel.Name] = &BoundStruct{
									//	Name: selector.Sel.Name,
									//}
									//p.parseStructFromExternalPackage(selector.Sel.Name, ident.Name, thisPackage)
									//p.addBoundStruct(ident.Name, selector.Sel.Name)
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

func (p *Project) addBoundMethods(packageName string, name string, boundMethods []*BoundMethod) {
	_, ok := p.BoundMethods[packageName]
	if !ok {
		p.BoundMethods[packageName] = make(map[structName][]*BoundMethod)
	}
	p.BoundMethods[packageName][name] = boundMethods
}

func (p *Project) parseBoundStructMethods(name string, pkg *ast.Package) error {
	var methods []*BoundMethod
	// Iterate over all files in the package
	for _, file := range pkg.Files {
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
							Inputs:     make([]*Parameter, 0),
							Outputs:    make([]*Parameter, 0),
						}

						method.Inputs = p.parseParameters(funcDecl.Type.Params)
						method.Outputs = p.parseParameters(funcDecl.Type.Results)

						methods = append(methods, method)
					}
				}
			}
		}
	}
	p.addBoundMethods(pkg.Name, name, methods)
	return nil
}

func (p *Project) parseParameters(params *ast.FieldList) []*Parameter {
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
			thisField.Type = getTypeString(field.Type)
			switch t := field.Type.(type) {
			case *ast.StarExpr:
				thisField.IsStruct = isStructType(t.X)
				thisField.IsPointer = true
			case *ast.StructType:
				thisField.IsStruct = true
			case *ast.ArrayType:
				thisField.IsSlice = true
				thisField.IsStruct = isStructType(t.Elt)
			case *ast.MapType:
				thisField.IsSlice = true
				thisField.IsStruct = isStructType(t.Value)
			}
			result = append(result, thisField)
		}
	}
	return result
}

func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return getTypeString(t.X)
	case *ast.ArrayType:
		return getTypeString(t.Elt)
	//case *ast.MapType:
	//	return "map[" + getTypeString(t.Key) + "]" + getTypeString(t.Value)
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
