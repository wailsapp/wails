package commands

import (
	"bytes"
	"fmt"
	"github.com/ciderapp/wails/v3/internal/parser"
	"os"
	"time"
)

type GenerateBindingsOptions struct {
	ModelsFilename   string `name:"m" description:"The filename for the models file" default:"models.ts"`
	BindingsFilename string `name:"b" description:"The filename for the bindings file" default:"bindings.js"`
	ProjectDirectory string `name:"d" description:"The project directory" default:"."`
}

func GenerateBindings(options *GenerateBindingsOptions) error {

	parserContext, err := parser.ParseDirectory(options.ProjectDirectory)
	if err != nil {
		return fmt.Errorf("error parsing project: %v", err)
	}

	// Generate models
	modelsData, err := parser.GenerateModels(parserContext)
	if err != nil {
		return fmt.Errorf("error generating models: %v", err)
	}

	var modelsFileData bytes.Buffer

	modelsFileData.WriteString(fmt.Sprintf(`// Copyright %v Cider Collective
// This file is automatically generated. DO NOT EDIT

`, time.Now().Format("2006")))

	modelsFileData.Write(modelsData)
	err = os.WriteFile(options.ModelsFilename, modelsFileData.Bytes(), 0755)
	if err != nil {
		return fmt.Errorf("error writing models file: %v", err)
	}
	println("Generated models file '" + options.ModelsFilename + "'")
	return nil
}
