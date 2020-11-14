package parser

import (
	"fmt"
	"go/ast"
)

func (p *Parser) getFunctionReturnType(pkg *Package, functionName string) (*Struct, error) {

	var result *Struct

	// Iterate through the files in the package looking for the bound structs
	for _, fileAst := range pkg.Gopackage.Syntax {

		var parseError error

		ast.Inspect(fileAst, func(n ast.Node) bool {
			// Parse Call expressions looking for bind calls
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if funcDecl.Name.Name == functionName {
				result, parseError = p.parseFunctionReturnType(fileAst, funcDecl, pkg)
				return false
			}

			return true
		})

		if parseError != nil {
			return nil, parseError
		}

		if result != nil {
			return result, nil
		}
	}

	return result, nil
}

func (p *Parser) parseFunctionReturnType(file *ast.File, funcDecl *ast.FuncDecl, pkg *Package) (*Struct, error) {

	var result *Struct

	if funcDecl.Type.Results == nil {
		return nil, fmt.Errorf("bound function %s has no return values", funcDecl.Name.Name)
	}

	// We expect only 1 return value for a function return
	if len(funcDecl.Type.Results.List) > 1 {
		return nil, fmt.Errorf("bound function %s has more than 1 return value", funcDecl.Name.Name)
	}

	parsedFields, err := p.parseField(file, funcDecl.Type.Results.List[0], pkg)
	if err != nil {
		return nil, err
	}

	if len(parsedFields) > 1 {
		return nil, fmt.Errorf("bound function %s has more than 1 return value", funcDecl.Name.Name)
	}

	result = parsedFields[0].Struct

	return result, nil
}
