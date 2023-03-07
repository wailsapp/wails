package commands

import "github.com/wailsapp/wails/v3/internal/parser"

type GenerateBindingsOptions struct {
	Silent           bool   `name:"silent" description:"Silent mode"`
	ModelsFilename   string `name:"m" description:"The filename for the models file" default:"models.ts"`
	BindingsFilename string `name:"b" description:"The filename for the bindings file" default:"bindings.js"`
	ProjectDirectory string `name:"p" description:"The project directory" default:"."`
	OutputDirectory  string `name:"d" description:"The output directory" default:"."`
}

func GenerateBindings(options *GenerateBindingsOptions) error {
	return parser.GenerateBindingsAndModels(options.ProjectDirectory, options.OutputDirectory)
}
