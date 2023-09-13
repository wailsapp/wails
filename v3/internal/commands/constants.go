package commands

import (
	"github.com/wailsapp/wails/v3/internal/parser"
	"os"
)

type GenerateConstantsOptions struct {
	ConstantsFilename string `name:"f" description:"The filename for the models file" default:"constants.go"`
	OutputFilename    string `name:"o" description:"The output file" default:"constants.js"`
}

func GenerateConstants(options *GenerateConstantsOptions) error {
	goData, err := os.ReadFile(options.ConstantsFilename)
	if err != nil {
		return err
	}

	result, err := parser.GenerateConstants(goData)
	if err != nil {
		return err
	}

	return os.WriteFile(options.OutputFilename, []byte(result), 0644)
}
