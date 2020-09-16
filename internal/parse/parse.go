package main

import (
	"fmt"
	"go/ast"
	"os"
	"strings"

	"github.com/leaanthony/slicer"
	"golang.org/x/tools/go/packages"
)

var internalMethods = slicer.String([]string{"WailsInit", "Wails Shutdown"})

var structCache = make(map[string]*ParsedStruct)
var boundStructs = make(map[string]*ParsedStruct)
var boundMethods = []string{}
var boundStructPointerLiterals = []string{}
var boundStructLiterals = slicer.StringSlicer{}
var boundVariables = slicer.StringSlicer{}
var app = ""
var structPointerFunctionDecls = make(map[string]string)
var structFunctionDecls = make(map[string]string)
var variableStructDecls = make(map[string]string)
var variableFunctionDecls = make(map[string]string)

type Parameter struct {
	Name string
	Type string
}

type ParsedMethod struct {
	Struct   string
	Name     string
	Comments []string
	Inputs   []*Parameter
	Returns  []*Parameter
}

type ParsedStruct struct {
	Name    string
	Methods []*ParsedMethod
}

type BoundStructs []*ParsedStruct

func ParseProject(projectPath string) (BoundStructs, error) {

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypesInfo}
	pkgs, err := packages.Load(cfg, projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	// Iterate the packages
	for _, pkg := range pkgs {

		// Iterate the files
		for _, file := range pkg.Syntax {

			var wailsPkgVar = ""

			ast.Inspect(file, func(n ast.Node) bool {
				var s string
				switch x := n.(type) {
				// Parse import declarations
				case *ast.ImportSpec:
					// Determine what wails has been imported as
					if x.Path.Value == `"github.com/wailsapp/wails/v2"` {
						wailsPkgVar = x.Name.Name
					}
				// Parse calls. We are looking for app.Bind() calls
				case *ast.CallExpr:
					f, ok := x.Fun.(*ast.SelectorExpr)
					if ok {
						n, ok := f.X.(*ast.Ident)
						if ok {
							//Check this is the Bind() call associated with the app variable
							if n.Name == app && f.Sel.Name == "Bind" {
								if len(x.Args) == 1 {
									ce, ok := x.Args[0].(*ast.CallExpr)
									if ok {
										n, ok := ce.Fun.(*ast.Ident)
										if ok {
											// We found a bind method using a function call
											// EG: app.Bind( newMyStruct() )
											boundMethods = append(boundMethods, n.Name)
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
														boundStructPointerLiterals = append(boundStructPointerLiterals, t.Name)
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
													boundStructLiterals.Add(t.Name)
												}
											} else {
												// Also check for when we bind a variable
												// myVariable := &MyStruct{}
												// app.Bind( myVariable )
												i, ok := x.Args[0].(*ast.Ident)
												if ok {
													boundVariables.Add(i.Name)
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
									if i.Name == wailsPkgVar {
										// Check we are calling a function to create the app
										if se.Sel.Name == "CreateApp" || se.Sel.Name == "CreateAppWithOptions" {
											if len(x.Lhs) == 1 {
												i, ok := x.Lhs[0].(*ast.Ident)
												if ok {
													// Found the app variable name
													app = i.Name
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
											variableFunctionDecls[i.Name] = fe.Name
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
												variableStructDecls[i.Name] = t.Name
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
									if internalMethods.Contains(x.Name.Name) {
										continue
									}
									// If we haven't already found this struct,
									// Create a placeholder in the cache
									parsedStruct := structCache[i.Name]
									if parsedStruct == nil {
										structCache[i.Name] = &ParsedStruct{
											Name: i.Name,
										}
										parsedStruct = structCache[i.Name]
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
												stringComment := comment.Text
												if strings.HasPrefix(stringComment, "//") {
													stringComment = stringComment[2:]
												}
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
												structMethod.Inputs = append(structMethod.Inputs, &Parameter{Name: name.Name, Type: t.Name})
											}
										}

										// Save the output parameters
										for _, outputField := range x.Type.Results.List {
											t, ok := outputField.Type.(*ast.Ident)
											if !ok {
												continue
											}
											if len(outputField.Names) == 0 {
												structMethod.Returns = append(structMethod.Returns, &Parameter{Type: t.Name})
											} else {
												for _, name := range outputField.Names {
													structMethod.Returns = append(structMethod.Returns, &Parameter{Name: name.Name, Type: t.Name})
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
										structPointerFunctionDecls[functionName] = s.Name
									}
								} else {
									// Check for functions that return a struct
									// This is to help us provide hints if the user binds a struct
									t, ok := x.Type.Results.List[0].Type.(*ast.Ident)
									if ok {
										// println("*** Function", functionName, "found which returns: "+t.Name)
										structFunctionDecls[functionName] = t.Name
									}
								}
							}
						}
					}
				}
				return true
			})
			// spew.Dump(file)
		}
	}

	/***** Update bound structs ******/

	// Resolve bound Methods
	for _, method := range boundMethods {
		s, ok := structPointerFunctionDecls[method]
		if !ok {
			s, ok = structFunctionDecls[method]
			if !ok {
				println("Fatal: Bind statement using", method, "but cannot find", method, "declaration")
			} else {
				println("Fatal: Cannot bind struct using method `" + method + "` because it returns a struct (" + s + "). Return a pointer to " + s + " instead.")
			}
			os.Exit(1)
		}
		structDefinition := structCache[s]
		if structDefinition == nil {
			println("Fatal: Bind statement using `"+method+"` but cannot find struct", s, "definition")
			os.Exit(1)
		}
		boundStructs[s] = structDefinition
	}

	// Resolve bound vars
	for _, structLiteral := range boundStructPointerLiterals {
		s, ok := structCache[structLiteral]
		if !ok {
			println("Fatal: Bind statement using", structLiteral, "but cannot find", structLiteral, "declaration")
			os.Exit(1)
		}
		boundStructs[structLiteral] = s
	}

	// Resolve bound variables
	boundVariables.Each(func(variable string) {
		v, ok := variableStructDecls[variable]
		if !ok {
			method, ok := variableFunctionDecls[variable]
			if !ok {
				println("Fatal: Bind statement using variable `" + variable + "` which does not resolve to a struct pointer")
				os.Exit(1)
			}

			// Resolve function name
			v, ok = structPointerFunctionDecls[method]
			if !ok {
				v, ok = structFunctionDecls[method]
				if !ok {
					println("Fatal: Bind statement using", method, "but cannot find", method, "declaration")
				} else {
					println("Fatal: Cannot bind variable `" + variable + "` because it resolves to a struct (" + v + "). Return a pointer to " + v + " instead.")
				}
				os.Exit(1)
			}

		}

		s, ok := structCache[v]
		if !ok {
			println("Fatal: Bind statement using variable `" + variable + "` which resolves to a `" + v + "` but cannot find its declaration")
			os.Exit(1)
		}
		boundStructs[v] = s

	})

	// Check for struct literals
	boundStructLiterals.Each(func(structName string) {
		println("Fatal: Cannot bind struct using struct literal `" + structName + "{}`. Create a pointer to " + structName + " instead.")
		os.Exit(1)
	})

	// Check for bound variables
	// boundVariables.Each(func(varName string) {
	// 	println("Fatal: Cannot bind struct using struct literal `" + structName + "{}`. Create a pointer to " + structName + " instead.")
	// })

	// spew.Dump(boundStructs)
	// os.Exit(0)

	// }
	// Inspect the AST and print all identifiers and literals.

	println("export {")

	noOfStructs := len(boundStructs)
	structCount := 0
	for _, s := range boundStructs {
		structCount++
		println()
		println("  " + s.Name + ": {")
		println()
		noOfMethods := len(s.Methods)
		for methodCount, m := range s.Methods {
			println("   /****************")
			for _, comment := range m.Comments {
				println("    *", comment)
			}
			if len(m.Comments) > 0 {
				println("    *")
			}
			inputNames := ""
			for _, input := range m.Inputs {
				println("    * @param {"+input.Type+"}", input.Name)
				inputNames += input.Name + ", "
			}
			print("    * @return Promise<")
			for _, output := range m.Returns {
				print(output.Type + "|")
			}
			println("Error>")
			println("    *")
			println("    ***/")
			if len(inputNames) > 2 {
				inputNames = inputNames[:len(inputNames)-2]
			}
			println("   ", m.Name+": function("+inputNames+") {")
			println("     return window.backend." + s.Name + "." + m.Name + "(" + inputNames + ");")
			print("    }")
			if methodCount < noOfMethods-1 {
				print(",")
			}
			println()
			println()
		}
		print("  }")
		if structCount < noOfStructs-1 {
			print(",")
		}
		println()
	}
	println()
	println("}")
	println()
}
