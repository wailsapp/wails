package parser

import (
	"go/ast"

	"github.com/pkg/errors"
)

func (p *Parser) parseStructFields(fileAst *ast.File, structType *ast.StructType, boundStruct *Struct) error {

	// Parse the fields
	for _, field := range structType.Fields.List {
		fields, err := p.parseField(fileAst, field, boundStruct.Package)
		if err != nil {
			return errors.Wrap(err, "error parsing struct "+boundStruct.Name)
		}

		// If this field was a struct, flag that it is used as data
		if len(fields) > 0 {
			if fields[0].Struct != nil {
				fields[0].Struct.IsUsedAsData = true
			}
		}

		// If this field name is lowercase, it won't be exported
		for _, field := range fields {
			if !startsWithLowerCaseLetter(field.Name) {
				boundStruct.Fields = append(boundStruct.Fields, field)
			}
		}

	}

	return nil
}
