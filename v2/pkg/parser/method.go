package parser

import (
	"fmt"
	"go/ast"
	"strings"
)

// Method defines a struct method
type Method struct {
	Name     string
	Comments []string
	Inputs   []*Field
	Returns  []*Field
}

func (p *Parser) parseStructMethods(boundStruct *Struct) error {

	for _, fileAst := range boundStruct.Package.Gopackage.Syntax {

		// Track errors
		var parseError error

		ast.Inspect(fileAst, func(n ast.Node) bool {

			if funcDecl, ok := n.(*ast.FuncDecl); ok {

				if funcDecl.Recv == nil {
					return true
				}

				// This is a struct method
				for _, field := range funcDecl.Recv.List {
					switch f := field.Type.(type) {
					case *ast.StarExpr:
						// This is a struct pointer method
						ident, ok := f.X.(*ast.Ident) // _ ?
						if !ok {
							continue
						}

						// Check this method is for this struct
						if ident.Name != boundStruct.Name {
							continue
						}

						// If this method is not Public, ignore
						if string(funcDecl.Name.Name[0]) != strings.ToUpper((string(funcDecl.Name.Name[0]))) {
							continue
						}

						// Create our struct
						structMethod := &Method{
							Name:     funcDecl.Name.Name,
							Comments: parseComments(funcDecl.Doc),
						}

						// Save the input parameters
						if funcDecl.Type.Params != nil {
							for _, inputField := range funcDecl.Type.Params.List {
								fields, err := p.parseField(fileAst, inputField, boundStruct.Package)
								if err != nil {
									parseError = err
									return false
								}

								// If this field was a struct, flag that it is used as data
								if len(fields) > 0 {
									if fields[0].Struct != nil {
										fields[0].Struct.IsUsedAsData = true
									}
								}

								structMethod.Inputs = append(structMethod.Inputs, fields...)
							}
						}

						// Save the output parameters
						if funcDecl.Type.Results != nil {
							for _, outputField := range funcDecl.Type.Results.List {
								fields, err := p.parseField(fileAst, outputField, boundStruct.Package)
								if err != nil {
									parseError = err
									return false
								}

								// If this field was a struct, flag that it is used as data
								if len(fields) > 0 {
									if fields[0].Struct != nil {
										fields[0].Struct.IsUsedAsData = true
									}
								}

								structMethod.Returns = append(structMethod.Returns, fields...)
							}
						}

						// Append this method to the parsed struct
						boundStruct.Methods = append(boundStruct.Methods, structMethod)

					default:
						// Unsupported
						continue
					}
				}
			}
			return true
		})

		// If we got an error, return it
		if parseError != nil {
			return parseError
		}
	}

	return nil
}

// InputsAsTSText generates a string with the method inputs
// formatted in a way acceptable to Typescript
func (m *Method) InputsAsTSText(pkgName string) string {
	var inputs []string

	for _, input := range m.Inputs {
		inputText := fmt.Sprintf("%s: %s", input.Name, goTypeToTS(input, pkgName))
		inputs = append(inputs, inputText)
	}

	return strings.Join(inputs, ", ")
}

// OutputsAsTSText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *Method) OutputsAsTSText(pkgName string) string {

	if len(m.Returns) == 0 {
		return "void"
	}

	var result []string

	for _, output := range m.Returns {
		result = append(result, goTypeToTS(output, pkgName))
	}
	return strings.Join(result, ", ")
}

// OutputsAsTSDeclarationText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *Method) OutputsAsTSDeclarationText(pkgName string) string {

	if len(m.Returns) == 0 {
		return "void"
	}

	var result []string

	for _, output := range m.Returns {
		result = append(result, goTypeToTSDeclaration(output, pkgName))
	}
	return strings.Join(result, ", ")
}

// InputsAsJSText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *Method) InputsAsJSText() string {
	var inputs []string

	for _, input := range m.Inputs {
		inputs = append(inputs, input.Name)
	}

	return strings.Join(inputs, ", ")
}
