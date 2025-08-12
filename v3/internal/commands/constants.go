package commands

import (
	"os"

	"github.com/wailsapp/wails/v3/internal/generator"
)

type GenerateConstantsOptions struct {
	ConstantsFilename string `name:"f" description:"The filename for the models file" default:"constants.go"`
	OutputFilename    string `name:"o" description:"The output file" default:"constants.js"`
}

func GenerateConstants(options *GenerateConstantsOptions) error {
	DisableFooter = true
	goData, err := os.ReadFile(options.ConstantsFilename)
	if err != nil {
		return err
	}

	result, err := generator.GenerateConstants(goData)
	if err != nil {
		return err
	}

	return os.WriteFile(options.OutputFilename, []byte(result), 0644)
}
