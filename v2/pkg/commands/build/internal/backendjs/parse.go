package backendjs

import (
	"fmt"
	"go/ast"
	"os"
	"strings"

	"github.com/leaanthony/slicer"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type Parser struct {
	wailsPkgVar string
	appVarName  string

	boundStructLiterals        slicer.StringSlicer
	boundMethods               []string
	boundStructs               map[string]*ParsedStruct
	boundStructPointerLiterals []string
	boundVariables             slicer.StringSlicer

	variableFunctionDecls map[string]string
	variableStructDecls   map[string]string

	internalMethods slicer.StringSlicer

	structCache                map[string]*ParsedStruct
	structPointerFunctionDecls map[string]string
	structFunctionDecls        map[string]string
}

type ParsedParameter struct {
	Name    string
	Type    string
	IsArray bool
}

func (p *ParsedParameter) JSType() string {
	return string(goTypeToJS(p.Type))
}

type ParsedMethod struct {
	Struct   string
	Name     string
	Comments []string
	Inputs   []*ParsedParameter
	Returns  []*ParsedParameter
}

// InputsAsTSText generates a string with the method inputs
// formatted in a way acceptable to Typescript
func (m *ParsedMethod) InputsAsTSText() string {
	var inputs []string

	for _, input := range m.Inputs {
		inputText := fmt.Sprintf("%s: %s", input.Name, goTypeToTS(input))
		inputs = append(inputs, inputText)
	}

	return strings.Join(inputs, ", ")
}

// InputsAsJSText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *ParsedMethod) InputsAsJSText() string {
	var inputs []string

	for _, input := range m.Inputs {
		inputs = append(inputs, input.Name)
	}

	return strings.Join(inputs, ", ")
}

// OutputsAsTSText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *ParsedMethod) OutputsAsTSText() string {

	if len(m.Returns) == 0 {
		return "void"
	}

	var result []string

	for _, output := range m.Returns {
		result = append(result, goTypeToTS(output))
	}
	return strings.Join(result, ", ")
}

type ParsedStruct struct {
	Name    string
	Methods []*ParsedMethod
}

func NewParser() *Parser {
	return &Parser{
		variableFunctionDecls:      make(map[string]string),
		variableStructDecls:        make(map[string]string),
		internalMethods:            *slicer.String([]string{"WailsInit", "WailsShutdown"}),
		structCache:                make(map[string]*ParsedStruct),
		structPointerFunctionDecls: make(map[string]string),
		structFunctionDecls:        make(map[string]string),
		boundStructs:               make(map[string]*ParsedStruct),
	}
}

func parseProject(projectPath string) ([]*Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, projectPath)
	if err != nil {
		return nil, errors.Wrap(err, "Problem loading packages")
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.Wrap(err, "Errors during parsing")
	}

	var result []*Package

	p := NewParser()

	// Iterate the packages
	for _, pkg := range pkgs {

		thisPackage, err := p.parsePackage(pkg)
		if err != nil {
			return nil, err
		}

		for k := range p.structCache {
			thisPackage.Structs = append(thisPackage.Structs, p.structCache[k])
		}

		result = append(result, thisPackage)

	}

	// Resolve links between data
	err = p.Resolve()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Parser) parsePackage(pkg *packages.Package) (*Package, error) {
	result := &Package{Name: pkg.Name}

	for _, file := range pkg.Syntax {
		err := p.parseFile(file)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (p *Parser) parseFile(file *ast.File) error {
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		// Parse import declarations
		case *ast.ImportSpec:
			// Determine what wails has been imported as
			if x.Path.Value == `"github.com/wailsapp/wails/v2"` {
				p.wailsPkgVar = x.Name.Name
			}
		// Parse calls. We are looking for app.Bind() calls
		case *ast.CallExpr:
			f, ok := x.Fun.(*ast.SelectorExpr)
			if ok {
				n, ok := f.X.(*ast.Ident)
				if ok {
					//Check this is the Bind() call associated with the app variable
					if n.Name == p.appVarName && f.Sel.Name == "Bind" {
						if len(x.Args) == 1 {
							ce, ok := x.Args[0].(*ast.CallExpr)
							if ok {
								n, ok := ce.Fun.(*ast.Ident)
								if ok {
									// We found a bind method using a function call
									// EG: app.Bind( newMyStruct() )
									p.boundMethods = append(p.boundMethods, n.Name)
								}
							} else {
								// We also want to check for Bind( &MyStruct{} )
								ue, ok := x.Args[0].(*ast.UnaryExpr)
								if ok {
									if ue.Op.String() == "&" {
										cl, ok := ue.X.(*ast.CompositeLit)
										if ok {
											t, ok := cl.Type.(*ast.Ident)
											if ok {
												// We have found Bind( &MyStruct{} )
												p.boundStructPointerLiterals = append(p.boundStructPointerLiterals, t.Name)
											}
										}
									}
								} else {
									// Let's check when the user binds a struct,
									// rather than a struct pointer: Bind( MyStruct{} )
									// We do this to provide better hints to the user
									cl, ok := x.Args[0].(*ast.CompositeLit)
									if ok {
										t, ok := cl.Type.(*ast.Ident)
										if ok {
											p.boundStructLiterals.Add(t.Name)
										}
									} else {
										// Also check for when we bind a variable
										// myVariable := &MyStruct{}
										// app.Bind( myVariable )
										i, ok := x.Args[0].(*ast.Ident)
										if ok {
											p.boundVariables.Add(i.Name)
										}
									}
								}
							}
						}
					}
				}
			}

		// We scan assignments for a number of reasons:
		//   * Determine the variable containing the main application
		//   * Determine the type of variables that get used in Bind()
		//   * Determine the type of variables that get created with var := &MyStruct{}
		case *ast.AssignStmt:
			for _, rhs := range x.Rhs {
				ce, ok := rhs.(*ast.CallExpr)
				if ok {
					se, ok := ce.Fun.(*ast.SelectorExpr)
					if ok {
						i, ok := se.X.(*ast.Ident)
						if ok {
							// Have we found the wails package name?
							if i.Name == p.wailsPkgVar {
								// Check we are calling a function to create the app
								if se.Sel.Name == "CreateApp" || se.Sel.Name == "CreateAppWithOptions" {
									if len(x.Lhs) == 1 {
										i, ok := x.Lhs[0].(*ast.Ident)
										if ok {
											// Found the app variable name
											p.appVarName = i.Name
										}
									}
								}
							}
						}
					} else {
						// Check for function assignment
						// a := newMyStruct()
						fe, ok := ce.Fun.(*ast.Ident)
						if ok {
							if len(x.Lhs) == 1 {
								i, ok := x.Lhs[0].(*ast.Ident)
								if ok {
									// Store the variable -> Function mapping
									// so we can later resolve the type
									p.variableFunctionDecls[i.Name] = fe.Name
								}
							}
						}
					}
				} else {
					// Check for literal assignment of struct
					// EG: myvar := MyStruct{}
					ue, ok := rhs.(*ast.UnaryExpr)
					if ok {
						cl, ok := ue.X.(*ast.CompositeLit)
						if ok {
							t, ok := cl.Type.(*ast.Ident)
							if ok {
								if len(x.Lhs) == 1 {
									i, ok := x.Lhs[0].(*ast.Ident)
									if ok {
										p.variableStructDecls[i.Name] = t.Name
									}
								}
							}
						}
					}
				}
			}
		// We scan for functions to build up a list of function names
		// for a number of reasons:
		//   * Determine which functions are struct methods that are bound
		//   * Determine
		case *ast.FuncDecl:
			if x.Recv != nil {
				// This is a struct method
				for _, field := range x.Recv.List {
					se, ok := field.Type.(*ast.StarExpr)
					if ok {
						// This is a struct pointer method
						i, ok := se.X.(*ast.Ident)
						if ok {
							// We want to ignore Internal functions
							if p.internalMethods.Contains(x.Name.Name) {
								continue
							}
							// If we haven't already found this struct,
							// Create a placeholder in the cache
							parsedStruct := p.structCache[i.Name]
							if parsedStruct == nil {
								p.structCache[i.Name] = &ParsedStruct{
									Name: i.Name,
								}
								parsedStruct = p.structCache[i.Name]
							}

							// If this method is Public
							if string(x.Name.Name[0]) == strings.ToUpper((string(x.Name.Name[0]))) {
								structMethod := &ParsedMethod{
									Struct: i.Name,
									Name:   x.Name.Name,
								}
								// Check if the method has comments.
								// If so, save it with the parsed method
								if x.Doc != nil {
									for _, comment := range x.Doc.List {
										stringComment := strings.TrimPrefix(comment.Text, "//")
										structMethod.Comments = append(structMethod.Comments, strings.TrimSpace(stringComment))
									}
								}

								// Save the input parameters
								for _, inputField := range x.Type.Params.List {
									t, ok := inputField.Type.(*ast.Ident)
									if !ok {
										continue
									}
									for _, name := range inputField.Names {
										structMethod.Inputs = append(structMethod.Inputs, &ParsedParameter{Name: name.Name, Type: t.Name})
									}
								}

								// Save the output parameters
								if x.Type.Results != nil {

									for _, outputField := range x.Type.Results.List {
										// Check for basic types
										t, ok := outputField.Type.(*ast.Ident)
										if !ok {
											// Check for arrays
											a, ok := outputField.Type.(*ast.ArrayType)
											if ok {
												// spew.Dump(a)
												ident := a.Elt.(*ast.Ident)
												if len(outputField.Names) == 0 {
													structMethod.Returns = append(structMethod.Returns, &ParsedParameter{Type: ident.Name, IsArray: true})
												} else {
													for _, name := range outputField.Names {
														structMethod.Returns = append(structMethod.Returns, &ParsedParameter{Name: name.Name, Type: ident.Name, IsArray: true})
													}
												}

											}

										} else {
											if len(outputField.Names) == 0 {
												structMethod.Returns = append(structMethod.Returns, &ParsedParameter{Type: t.Name})
											} else {
												for _, name := range outputField.Names {
													structMethod.Returns = append(structMethod.Returns, &ParsedParameter{Name: name.Name, Type: t.Name})
												}
											}
										}
									}
								}

								// Append this method to the parsed struct
								parsedStruct.Methods = append(parsedStruct.Methods, structMethod)

							}
						}
					}
				}
			} else {
				// This is a function declaration
				// We care about its name and return type
				// This will allow us to resolve types later
				functionName := x.Name.Name

				// Look for one that returns a single value
				if x.Type != nil && x.Type.Results != nil && x.Type.Results.List != nil {
					if len(x.Type.Results.List) == 1 {
						// Check for *struct
						t, ok := x.Type.Results.List[0].Type.(*ast.StarExpr)
						if ok {
							s, ok := t.X.(*ast.Ident)
							if ok {
								// println("*** Function", functionName, "found which returns: *"+s.Name)
								p.structPointerFunctionDecls[functionName] = s.Name
							}
						} else {
							// Check for functions that return a struct
							// This is to help us provide hints if the user binds a struct
							t, ok := x.Type.Results.List[0].Type.(*ast.Ident)
							if ok {
								// println("*** Function", functionName, "found which returns: "+t.Name)
								p.structFunctionDecls[functionName] = t.Name
							}
						}
					}
				}
			}
		}
		return true
	})
	// spew.Dump(file)

	return nil
}

func (p *Parser) Resolve() error {
	// Resolve bound Methods
	for _, method := range p.boundMethods {
		s, ok := p.structPointerFunctionDecls[method]
		if !ok {
			s, ok = p.structFunctionDecls[method]
			if !ok {
				return fmt.Errorf("bind statement using " + method + " but cannot find " + method + " declaration")
			} else {
				return fmt.Errorf("cannot bind struct using method `" + method + "` because it returns a struct (" + s + "). Return a pointer to " + s + " instead.")
			}
		}
		structDefinition := p.structCache[s]
		if structDefinition == nil {
			return fmt.Errorf("Fatal: Bind statement using `" + method + "` but cannot find struct " + s + " definition")
		}
		p.boundStructs[s] = structDefinition
	}

	// Resolve bound vars
	for _, structLiteral := range p.boundStructPointerLiterals {
		s, ok := p.structCache[structLiteral]
		if !ok {
			return fmt.Errorf("bind statement using " + structLiteral + " but cannot find " + structLiteral + " declaration")
		}
		p.boundStructs[structLiteral] = s
	}

	var err error

	// Resolve bound variables
	p.boundVariables.Each(func(variable string) {
		v, ok := p.variableStructDecls[variable]
		if !ok {
			method, ok := p.variableFunctionDecls[variable]
			if !ok {
				if err == nil {
					err = fmt.Errorf("bind statement using variable `" + variable + "` which does not resolve to a struct pointer")
				}
			}

			// Resolve function name
			v, ok = p.structPointerFunctionDecls[method]
			if !ok {
				v, ok = p.structFunctionDecls[method]
				if !ok {
					if err == nil {
						err = fmt.Errorf("bind statement using " + method + " but cannot find " + method + " declaration")
					}
				} else {
					if err == nil {
						err = fmt.Errorf("cannot bind variable `" + variable + "` because it resolves to a struct (" + v + "). Return a pointer to " + v + " instead.")
					}
				}
			}
		}

		s, ok := p.structCache[v]
		if !ok {
			println("Fatal: Bind statement using variable `" + variable + "` which resolves to a `" + v + "` but cannot find its declaration")
			os.Exit(1)
		}
		p.boundStructs[v] = s
	})

	// Return first error when resolving bound variables
	if err != nil {
		return err
	}

	// Check for struct literals
	if p.boundStructLiterals.Length() > 0 {
		return fmt.Errorf("cannot bind structs using struct literals. Create a pointer to the struct instead: %s", p.boundStructLiterals.Join(", "))
	}
	return nil
}
