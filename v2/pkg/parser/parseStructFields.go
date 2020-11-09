package parser

import "go/ast"

func (p *Parser) parseStructFields(fileAst *ast.File, structType *ast.StructType, boundStruct *Struct) error {

	// Parse the fields
	for _, field := range structType.Fields.List {
		fields, err := p.parseField(fileAst, field, boundStruct.Package)
		if err != nil {
			return err
		}
		boundStruct.Fields = append(boundStruct.Fields, fields...)
	}

	return nil
}
