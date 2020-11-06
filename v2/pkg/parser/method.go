package parser

import (
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
	for _, fileAst := range boundStruct.Package.Syntax {

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

						// We want to ignore Internal functions
						if p.internalMethods.Contains(funcDecl.Name.Name) {
							continue
						}

						// If this method is not Public, ignore
						if string(funcDecl.Name.Name[0]) != strings.ToUpper((string(funcDecl.Name.Name[0]))) {
							continue
						}

						// Create our struct
						structMethod := &Method{
							Name:     funcDecl.Name.Name,
							Comments: p.parseComments(funcDecl.Doc),
						}

						// Save the input parameters
						if funcDecl.Type.Params != nil {
							for _, inputField := range funcDecl.Type.Params.List {
								fields, err := p.parseField(inputField, boundStruct.Package.Name)
								if err != nil {
									parseError = err
									return false
								}

								structMethod.Inputs = append(structMethod.Inputs, fields...)
							}
						}

						// Save the output parameters
						if funcDecl.Type.Results != nil {
							for _, outputField := range funcDecl.Type.Results.List {
								fields, err := p.parseField(outputField, boundStruct.Package.Name)
								if err != nil {
									parseError = err
									return false
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
