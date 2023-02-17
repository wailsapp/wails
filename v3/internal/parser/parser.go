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

	"github.com/samber/lo"
)

var Debug = false

func debug(msg string, args ...interface{}) {
	if Debug {
		println(fmt.Sprintf(msg, args...))
	}
}

type BoundStruct struct {
	Name     string
	Methods  map[string]*FuncSignature
	Comments []string
}

type parsedPackage struct {
	name         string
	pkg          *ast.Package
	boundStructs map[string]*BoundStruct
}

type Context struct {
	packages map[string]*parsedPackage
	dir      string
}

func (c *Context) findImportPackage(pkgName string, pkg *ast.Package) (*ast.Package, error) {
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return nil, err
			}
			if imp.Name != nil && imp.Name.Name == pkgName {
				return c.getPackageFromPath(path)
			} else {
				_, pkgName := filepath.Split(path)
				if pkgName == pkgName {
					return c.getPackageFromPath(path)
				}
			}
		}
	}
	return nil, fmt.Errorf("package '%s' not found in %s", pkgName, pkg.Name)
}

func (c *Context) getPackageFromPath(path string) (*ast.Package, error) {
	dir, err := filepath.Abs(c.dir)
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(path) {
		dir = filepath.Join(dir, path)
	} else {
		impPkgDir, err := build.Import(path, dir, build.ImportMode(0))
		if err != nil {
			return nil, err
		}
		dir = impPkgDir.Dir
	}
	impPkg, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	for impName, impPkg := range impPkg {
		if impName == "main" {
			continue
		}
		return impPkg, nil
	}
	return nil, fmt.Errorf("Package not found in imported package %s", path)
}

func ParseDirectory(dir string) (*Context, error) {
	// Parse the directory
	fset := token.NewFileSet()
	if dir == "." || dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = cwd
	}
	println("Parsing directory " + dir)
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	context := &Context{
		dir:      dir,
		packages: make(map[string]*parsedPackage),
	}

	// Iterate through the packages
	for _, pkg := range pkgs {
		context.packages[pkg.Name] = &parsedPackage{
			name:         pkg.Name,
			pkg:          pkg,
			boundStructs: make(map[string]*BoundStruct),
		}
	}

	findApplicationNewCalls(context)
	err = findStructDefinitions(context)
	if err != nil {
		return nil, err
	}

	return context, nil
}

func findStructDefinitions(context *Context) error {
	// iterate over the packages
	for _, pkg := range context.packages {
		// iterate the struct names
		for structName, _ := range pkg.boundStructs {
			structSpec, methods, comments := getStructTypeSpec(pkg.pkg, structName)
			if structSpec == nil {
				return fmt.Errorf("unable to find struct %s in package %s", structName, pkg.name)
			}
			pkg.boundStructs[structName] = &BoundStruct{
				Name:     structName,
				Comments: comments,
			}
			if pkg.boundStructs[structName].Methods == nil {
				pkg.boundStructs[structName].Methods = make(map[string]*FuncSignature)
			}
			for _, method := range methods {
				pkg.boundStructs[structName].Methods[method.Name] = FuncTypeToSignature(method.Type)
				pkg.boundStructs[structName].Methods[method.Name].Comments = method.Comments
			}
		}
	}
	return nil
}

func findApplicationNewCalls(context *Context) {
	println("Finding application.New calls")
	// Iterate through the packages
	currentPackages := lo.Keys(context.packages)

	for _, packageName := range currentPackages {
		thisPackage := context.packages[packageName]
		debug("Parsing package: %s", packageName)
		// Iterate through the package's files
		for _, file := range thisPackage.pkg.Files {
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
								if ident.Obj == nil {
									thisPackage.boundStructs[ident.Name] = &BoundStruct{
										Name: ident.Name,
									}
									continue
								}
								// Check if the ident is a struct type
								if _, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
									thisPackage.boundStructs[ident.Name] = &BoundStruct{
										Name: ident.Name,
									}
									continue
								}
								// Check the typespec decl is a struct
								if _, ok := ident.Obj.Decl.(*ast.StructType); ok {
									continue
								}

							}
							// Check if the lit is a selector
							selector, ok := boundStructLit.Type.(*ast.SelectorExpr)
							if ok {
								// Check if the selector is an ident
								if ident, ok := selector.X.(*ast.Ident); ok {
									// Check if the ident is a package
									if _, ok := context.packages[ident.Name]; !ok {
										externalPackage, err := context.getPackageFromPath(ident.Name)
										if err != nil {
											println("Error getting package from path: " + err.Error())
											return true
										}
										context.packages[ident.Name] = &parsedPackage{
											name:         ident.Name,
											pkg:          externalPackage,
											boundStructs: make(map[string]*BoundStruct),
										}
									}
									context.packages[ident.Name].boundStructs[selector.Sel.Name] = &BoundStruct{
										Name: selector.Sel.Name,
									}
								}
								continue
							}
						}
					}
				}

				return true
			})
		}
	}
}

//type Method struct {
//	Name string
//	Type *ast.FuncType
//}

//func getStructTypeSpec(pkg *ast.Package, structName string) (*ast.TypeSpec, []Method) {
//	var typeSpec *ast.TypeSpec
//	var methods []Method
//
//	// Iterate over all files in the package
//	for _, file := range pkg.Files {
//		// Iterate over all declarations in the file
//		for _, decl := range file.Decls {
//			// Check if the declaration is a type declaration
//			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
//				// Iterate over all type specifications in the type declaration
//				for _, spec := range genDecl.Specs {
//					// Check if the type specification is a struct type specification
//					if tSpec, ok := spec.(*ast.TypeSpec); ok && tSpec.Name.Name == structName {
//						// Check if the type specification is a struct type
//						if _, ok := tSpec.Type.(*ast.StructType); ok {
//							typeSpec = tSpec
//						}
//					}
//				}
//			} else if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil {
//				// Check if the function has a receiver argument of the struct type
//				recvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
//				if ok {
//					if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == structName {
//						// Add the method to the list of methods
//						method := Method{
//							Name: funcDecl.Name.Name,
//							Type: funcDecl.Type,
//						}
//						methods = append(methods, method)
//					}
//				}
//			}
//		}
//	}
//
//	return typeSpec, methods
//}

type Arg struct {
	Name string
	Type string
}

type FuncSignature struct {
	Comments []string
	Inputs   []Arg
	Outputs  []Arg
}

func FuncTypeToSignature(ft *ast.FuncType) *FuncSignature {
	sig := &FuncSignature{}

	// process input arguments
	if ft.Params != nil {
		for _, field := range ft.Params.List {
			arg := Arg{}
			for _, name := range field.Names {
				arg.Name = name.Name
			}
			arg.Type = tokenToString(field.Type)
			sig.Inputs = append(sig.Inputs, arg)
		}
	}

	// process output arguments
	if ft.Results != nil {
		for _, field := range ft.Results.List {
			arg := Arg{}
			arg.Type = tokenToString(field.Type)
			sig.Outputs = append(sig.Outputs, arg)
		}
	}

	return sig
}

func tokenToString(t ast.Expr) string {
	switch t := t.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + tokenToString(t.X)
	case *ast.SelectorExpr:
		return tokenToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + tokenToString(t.Elt)
	case *ast.StructType:
		return "struct"
	default:
		return ""
	}
}

type Method struct {
	Name     string
	Type     *ast.FuncType
	Comments []string // Add a field to capture comments for the method
}

func getStructTypeSpec(pkg *ast.Package, structName string) (*ast.TypeSpec, []Method, []string) {
	var typeSpec *ast.TypeSpec
	var methods []Method
	var structComments []string

	// Iterate over all files in the package
	for _, file := range pkg.Files {
		// Iterate over all declarations in the file
		for _, decl := range file.Decls {
			// Check if the declaration is a type declaration
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				// Iterate over all type specifications in the type declaration
				for _, spec := range genDecl.Specs {
					// Check if the type specification is a struct type specification
					if tSpec, ok := spec.(*ast.TypeSpec); ok && tSpec.Name.Name == structName {
						// Check if the type specification is a struct type
						if _, ok := tSpec.Type.(*ast.StructType); ok {
							// Get comments associated with the struct
							if genDecl.Doc != nil {
								for _, comment := range genDecl.Doc.List {
									structComments = append(structComments, comment.Text)
								}
							}
							typeSpec = tSpec
						}
					}
				}
			} else if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil {
				// Check if the function has a receiver argument of the struct type
				recvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
				if ok {
					if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == structName {
						// Get comments associated with the method
						if funcDecl.Doc != nil {
							var comments []string
							for _, comment := range funcDecl.Doc.List {
								comments = append(comments, comment.Text)
							}
							// Add the method to the list of methods
							method := Method{
								Name:     funcDecl.Name.Name,
								Type:     funcDecl.Type,
								Comments: comments,
							}
							methods = append(methods, method)
						}
					}
				}
			}
		}
	}

	return typeSpec, methods, structComments
}
