package generator

import (
	"updater/generator/idlreader"
	"updater/generator/types"
)

func ParseIDL(idlData []byte) ([]*types.GeneratedFile, error) {

	idl, err := idlreader.Parse(idlData)
	if err != nil {
		return nil, err
	}

	err = idl.Process()
	if err != nil {
		return nil, err
	}

	generatedFiles, err := idl.Generate()
	if err != nil {
		return nil, err
	}
	return generatedFiles, nil
}
