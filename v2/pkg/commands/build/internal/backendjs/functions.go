package backendjs

import (
	"go/ast"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func (p *Parser) parseField(field *ast.Field, pkg *Package) (string, *StructName) {
	var structName *StructName
	var fieldType string
	switch t := field.Type.(type) {
	case *ast.Ident:
		fieldType = t.Name
	case *ast.StarExpr:
		fieldType = "struct"
		structName = p.parseStructNameFromStarExpr(t)

		// Save external reference if we have it
		if structName.Package == "" {
			pkg.packageReferences.AddUnique(structName.Package)
			pkg.structsUsedAsData.AddUnique(structName.Name)
		} else {
			// Save this reference to the external struct
			referencedPackage := p.Packages[structName.Package]
			if referencedPackage == nil {
				// TODO: Move this to a global reference list instead of bombing out
				println("WARNING: Unknown package referenced by field:", structName.Package)
			}
			referencedPackage.structsUsedAsData.AddUnique(structName.Name)
			pkg.packageReferences.AddUnique(structName.Package)
		}

	default:
		spew.Dump(field)
		println("Unhandled Field type")
		os.Exit(1)
	}
	return fieldType, structName
}

func (p *Parser) parseFunctionDeclaration(funcDecl *ast.FuncDecl, pkg *Package) {
	if funcDecl.Recv != nil {
		// This is a struct method
		for _, field := range funcDecl.Recv.List {
			se, ok := field.Type.(*ast.StarExpr)
			if ok {
				// This is a struct pointer method
				i, ok := se.X.(*ast.Ident)
				if ok {
					// We want to ignore Internal functions
					if p.internalMethods.Contains(funcDecl.Name.Name) {
						continue
					}
					// If we haven't already found this struct,
					// Create a placeholder in the cache
					parsedStruct := pkg.Structs[i.Name]
					if parsedStruct == nil {
						pkg.Structs[i.Name] = &Struct{
							Name: i.Name,
						}
						parsedStruct = pkg.Structs[i.Name]
					}

					// If this method is Public
					if string(funcDecl.Name.Name[0]) == strings.ToUpper((string(funcDecl.Name.Name[0]))) {
						structMethod := &Method{
							Name: funcDecl.Name.Name,
						}
						// Check if the method has comments.
						// If so, save it with the parsed method
						if funcDecl.Doc != nil {
							structMethod.Comments = p.parseComments(funcDecl.Doc)
						}

						// Save the input parameters
						if funcDecl.Type.Params != nil {
							for _, inputField := range funcDecl.Type.Params.List {
								fieldType, structName := p.parseField(inputField, pkg)
								for _, name := range inputField.Names {
									structMethod.Inputs = append(structMethod.Inputs, &Field{
										Name:   name.Name,
										Type:   fieldType,
										Struct: structName,
									})
								}
							}
						}

						// Save the output parameters
						if funcDecl.Type.Results != nil {
							for _, outputField := range funcDecl.Type.Results.List {
								fieldType, structName := p.parseField(outputField, pkg)
								if len(outputField.Names) == 0 {
									structMethod.Returns = append(structMethod.Returns, &Field{
										Type:   fieldType,
										Struct: structName,
									})
								} else {
									for _, name := range outputField.Names {
										structMethod.Returns = append(structMethod.Returns, &Field{
											Name:   name.Name,
											Type:   fieldType,
											Struct: structName,
										})
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
		functionName := funcDecl.Name.Name

		// Look for one that returns a single value
		if funcDecl.Type != nil && funcDecl.Type.Results != nil && funcDecl.Type.Results.List != nil {
			if len(funcDecl.Type.Results.List) == 1 {
				// Check for *struct
				t, ok := funcDecl.Type.Results.List[0].Type.(*ast.StarExpr)
				if ok {
					s, ok := t.X.(*ast.Ident)
					if ok {
						// println("*** Function", functionName, "found which returns: *"+s.Name)
						pkg.functionsThatReturnStructPointers[functionName] = s.Name
					}
				} else {
					// Check for functions that return a struct
					// This is to help us provide hints if the user binds a struct
					t, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident)
					if ok {
						// println("*** Function", functionName, "found which returns: "+t.Name)
						pkg.functionsThatReturnStructs[functionName] = t.Name
					}
				}
			}
		}
	}
}
