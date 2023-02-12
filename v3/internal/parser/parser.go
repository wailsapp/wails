package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/samber/lo"

	"golang.org/x/tools/go/packages"
)

var Debug = false

func debug(msg string, args ...interface{}) {
	if Debug {
		println(fmt.Sprintf(msg, args...))
	}
}

type parsedPackage struct {
	name         string
	pkg          *ast.Package
	boundStructs map[string]*ast.TypeSpec
}

type Context struct {
	packages map[string]*parsedPackage
}

func (c *Context) GetBoundStructs() map[string][]*ast.TypeSpec {
	structs := make(map[string][]*ast.TypeSpec)
	for _, pkg := range c.packages {
		for _, structType := range pkg.boundStructs {
			structs[pkg.name] = append(structs[pkg.name], structType)
		}
	}
	return structs
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
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	context := &Context{
		packages: make(map[string]*parsedPackage),
	}

	// Iterate through the packages
	for _, pkg := range pkgs {
		context.packages[pkg.Name] = &parsedPackage{
			name:         pkg.Name,
			pkg:          pkg,
			boundStructs: make(map[string]*ast.TypeSpec),
		}
	}

	findApplicationNewCalls(context)

	return context, nil
}

func findApplicationNewCalls(context *Context) {
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
				if selectorExpr.Sel.Name != "Application" {
					return true
				}
				if id, ok := selectorExpr.X.(*ast.Ident); !ok || id.Name != "options" {
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
									structTypeSpec := findStructInPackage(thisPackage.pkg, ident.Name)
									thisPackage.boundStructs[ident.Name] = structTypeSpec
									findNestedStructs(structTypeSpec, file, packageName, context)
									continue
								}
								// Check if the ident is a struct type
								if t, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
									thisPackage.boundStructs[ident.Name] = t
									findNestedStructs(t, file, packageName, context)
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
								getStructsFromSelector(selector, file, context)
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

func getStructsFromSelector(selector *ast.SelectorExpr, file *ast.File, context *Context) {
	debug("getStructsFromSelector called with selector '%s' on file '%s.go'", selector.Sel.Name, file.Name.Name)

	// extract package name from selector
	packageName := selector.X.(*ast.Ident).Name

	if context.packages[packageName] == nil {
		context.packages[packageName] = &parsedPackage{
			name:         packageName,
			boundStructs: make(map[string]*ast.TypeSpec),
		}
	}

	// extract struct name from selector
	structName := selector.Sel.Name

	// Find the package name from the imports
	for _, imp := range file.Imports {
		var match bool
		if imp.Name == nil || imp.Name.Name == packageName {
			match = true
		}
		if match == false {
			pathSplit := strings.Split(imp.Path.Value, "/")
			endPath := strings.Trim(pathSplit[len(pathSplit)-1], `"`)
			match = endPath == packageName
		}

		if match {
			// We have the import
			cfg := &packages.Config{
				Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
			}
			pkgs, err := packages.Load(cfg, strings.Trim(imp.Path.Value, `"`))
			if err != nil {
				panic(err)
			}
			foundPackage := pkgs[0]

			// Iterate the files in the package and find struct types
			for _, parsedFile := range foundPackage.Syntax {
				ast.Inspect(parsedFile, func(n ast.Node) bool {
					if n == nil {
						return false
					}
					switch n.(type) {
					case *ast.TypeSpec:
						typeSpec := n.(*ast.TypeSpec)
						if typeSpec.Name.Name == structName {
							if _, ok := context.packages[packageName].boundStructs[structName]; !ok {
								debug("Adding struct '%s' in package '%s'", structName, packageName)
								context.packages[packageName].boundStructs[typeSpec.Name.Name] = typeSpec
								findNestedStructs(typeSpec, parsedFile, packageName, context)
							}
							return false
						}
					}
					return true
				})
			}

			continue
		}
	}

}

func findNestedStructs(t *ast.TypeSpec, parsedFile *ast.File, pkgName string, context *Context) {
	debug("findNestedStructs called with type '%s' on file '%s.go'", t.Name.Name, parsedFile.Name.Name)
	structType, ok := t.Type.(*ast.StructType)
	if !ok {
		return
	}
	for _, field := range structType.Fields.List {
		for _, ident := range field.Names {
			switch t := ident.Obj.Decl.(*ast.Field).Type.(type) {
			case *ast.Ident:
				if t.Obj == nil {
					continue
				}
				if t.Obj.Kind == ast.Typ {
					if _, ok := t.Obj.Decl.(*ast.TypeSpec); ok {
						if _, ok := context.packages[pkgName].boundStructs[t.Name]; !ok {
							debug("Adding nested struct '%s' to package '%s'", t.Name, pkgName)
							context.packages[pkgName].boundStructs[t.Name] = t.Obj.Decl.(*ast.TypeSpec)
							findNestedStructs(t.Obj.Decl.(*ast.TypeSpec), parsedFile, pkgName, context)
						}
					}
				}
			case *ast.SelectorExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					if ident.IsExported() {
						getStructsFromSelector(t, parsedFile, context)
					}
				}
			case *ast.StarExpr:
				if sel, ok := t.X.(*ast.SelectorExpr); ok {
					if _, ok := sel.X.(*ast.Ident); ok {
						if ident.IsExported() {
							getStructsFromSelector(sel, parsedFile, context)
						}
					}
				}
			}
		}
	}
	findStructsInMethods(t.Name.Name, parsedFile, pkgName, context)

}

func findStructsInMethods(name string, parsedFile *ast.File, pkgName string, context *Context) {
	debug("findStructsInMethods called with type '%s' on file '%s.go'", name, parsedFile.Name.Name)
	// Find the struct declaration for the given name
	var structDecl *ast.TypeSpec
	for _, decl := range parsedFile.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			// check the receiver name is the same as the name given
			if fn.Recv == nil {
				continue
			}
			// Check if the receiver is a pointer
			if starExpr, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
				if ident, ok := starExpr.X.(*ast.Ident); ok {
					if ident.Name != name {
						continue
					}
				}
			} else {
				if ident, ok := fn.Recv.List[0].Type.(*ast.Ident); ok {
					if ident.Name != name {
						continue
					}
				}
			}
			findStructsInMethodParams(fn, parsedFile, pkgName, context)
		}
	}
	if structDecl == nil {
		return
	}
	// Iterate the methods in the struct

}

func findStructsInMethodParams(f *ast.FuncDecl, parsedFile *ast.File, pkgName string, context *Context) {
	debug("findStructsInMethodParams called with type '%s' on file '%s.go'", f.Name.Name, parsedFile.Name.Name)
	if f.Type.Params == nil {
		for _, field := range f.Type.Params.List {
			parseField(field, parsedFile, pkgName, context)
		}
	}
	if f.Type.Results != nil {
		for _, field := range f.Type.Results.List {
			parseField(field, parsedFile, pkgName, context)
		}
	}
}

func parseField(field *ast.Field, parsedFile *ast.File, pkgName string, context *Context) {
	if se, ok := field.Type.(*ast.StarExpr); ok {
		// Check if the star expr is a struct
		if selExp, ok := se.X.(*ast.SelectorExpr); ok {
			getStructsFromSelector(selExp, parsedFile, context)
			return
		}
		if ident, ok := se.X.(*ast.Ident); ok {
			if ident.Obj == nil {
				return
			}
			if ident.Obj.Kind == ast.Typ {
				if _, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
					if _, ok := context.packages[pkgName].boundStructs[ident.Name]; !ok {
						debug("Adding field struct '%s' to package '%s'", ident.Name, pkgName)
						context.packages[pkgName].boundStructs[ident.Name] = ident.Obj.Decl.(*ast.TypeSpec)
						findNestedStructs(ident.Obj.Decl.(*ast.TypeSpec), parsedFile, pkgName, context)
					} else {
						debug("Struct %s already bound", ident.Name)
					}
				}
			}
		}
	}
	if selExp, ok := field.Type.(*ast.SelectorExpr); ok {
		getStructsFromSelector(selExp, parsedFile, context)
		return
	}
	if ident, ok := field.Type.(*ast.Ident); ok {
		if ident.Obj == nil {
			return
		}
		if ident.Obj.Kind == ast.Typ {
			if _, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
				if _, ok := context.packages[pkgName].boundStructs[ident.Name]; !ok {
					debug("Adding field struct '%s' to package '%s'", ident.Name, pkgName)
					context.packages[pkgName].boundStructs[ident.Name] = ident.Obj.Decl.(*ast.TypeSpec)
					findNestedStructs(ident.Obj.Decl.(*ast.TypeSpec), parsedFile, pkgName, context)
				} else {
					debug("Struct %s already bound", ident.Name)
				}
			}
		}
	}
}

func findStructInPackage(pkg *ast.Package, name string) *ast.TypeSpec {
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
				for _, spec := range gen.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if typeSpec.Name.Name == name {
							if _, ok := typeSpec.Type.(*ast.StructType); ok {
								return typeSpec
							}
						}
					}
				}
			}
		}
	}
	return nil
}

type Package struct {
	Name  string
	Specs []*ast.TypeSpec
}

var goToTS = map[string]string{
	"int":     "number",
	"int8":    "number",
	"int16":   "number",
	"int32":   "number",
	"int64":   "number",
	"uint":    "number",
	"uint8":   "number",
	"uint16":  "number",
	"uint32":  "number",
	"uint64":  "number",
	"float32": "number",
	"float64": "number",
	"string":  "string",
	"bool":    "boolean",
}

//func GenerateModels(specs map[string][]*ast.TypeSpec) ([]byte, error) {
//	var buf bytes.Buffer
//	var packages []Package
//	for pkg, pkgSpecs := range specs {
//		packages = append(packages, Package{Name: pkg, Specs: pkgSpecs})
//	}
//	sort.Slice(packages, func(i, j int) bool { return packages[i].Name < packages[j].Name })
//	for _, pkg := range packages {
//		if _, err := fmt.Fprintf(&buf, "namespace %s {\n", pkg.Name); err != nil {
//			return nil, err
//		}
//		sort.Slice(pkg.Specs, func(i, j int) bool { return pkg.Specs[i].Name.Name < pkg.Specs[j].Name.Name })
//		for _, spec := range pkg.Specs {
//			if structType, ok := spec.Type.(*ast.StructType); ok {
//				if _, err := fmt.Fprintf(&buf, "  class %s {\n", spec.Name.Name); err != nil {
//					return nil, err
//				}
//
//				for _, field := range structType.Fields.List {
//
//					// Get the Go type of the field
//					goType := types.ExprString(field.Type)
//					// Look up the corresponding TypeScript type
//					tsType, ok := goToTS[goType]
//					if !ok {
//						tsType = goType
//					}
//
//					if _, err := fmt.Fprintf(&buf, "    %s: %s;\n", field.Names[0].Name, tsType); err != nil {
//						return nil, err
//					}
//				}
//
//				if _, err := fmt.Fprintf(&buf, "  }\n"); err != nil {
//					return nil, err
//				}
//				if _, err := fmt.Fprintf(&buf, "  }\n"); err != nil {
//					return nil, err
//				}
//			}
//		}
//
//		if _, err := fmt.Fprintf(&buf, "}\n"); err != nil {
//			return nil, err
//		}
//	}
//	return buf.Bytes(), nil
//}

type allModels struct {
	known map[string]map[string]struct{}
}

func newAllModels(models map[string][]*ast.TypeSpec) *allModels {
	result := &allModels{known: make(map[string]map[string]struct{})}
	// iterate over all models
	for pkg, pkgSpecs := range models {
		for _, spec := range pkgSpecs {
			result.known[pkg] = make(map[string]struct{})
			result.known[pkg][spec.Name.Name] = struct{}{}
		}
	}
	return result
}

func (k *allModels) exists(name string) bool {
	// Split the name into package and type
	parts := strings.Split(name, ".")
	typ := parts[0]
	pkg := "main"
	if len(parts) == 2 {
		pkg = parts[0]
		typ = parts[1]
	}

	knownPkg, ok := k.known[pkg]
	if !ok {
		return false
	}
	_, ok = knownPkg[typ]
	return ok
}
