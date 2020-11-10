package parser

import "go/ast"

func (p *Parser) parseStructFields(fileAst *ast.File, structType *ast.StructType, boundStruct *Struct) error {

	// Parse the fields
	for _, field := range structType.Fields.List {
		fields, err := p.parseField(fileAst, field, boundStruct.Package)
		if err != nil {
			return err
		}

		// If this field was a struct, flag that it is used as data
		if len(fields) > 0 {
			if fields[0].Struct != nil {
				fields[0].Struct.IsUsedAsData = true
			}
		}

		boundStruct.Fields = append(boundStruct.Fields, fields...)
	}

	return nil
}
