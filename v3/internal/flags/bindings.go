package flags

type GenerateBindingsOptions struct {
	Silent           bool   `name:"silent" description:"Silent mode"`
	ModelsFilename   string `name:"m" description:"The filename for the models file" default:"models.ts"`
	TS               bool   `name:"ts" description:"Generate Typescript bindings"`
	TSPrefix         string `description:"The prefix for the typescript names" default:""`
	TSSuffix         string `description:"The postfix for the typescript names" default:""`
	UseInterfaces    bool   `name:"i" description:"Use interfaces instead of classes"`
	ProjectDirectory string `name:"p" description:"The project directory" default:"."`
	OutputDirectory  string `name:"d" description:"The output directory" default:"frontend/bindings"`
}
