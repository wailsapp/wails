package flags

type GenerateBindingsOptions struct {
	Silent            bool   `name:"silent" description:"Silent mode"`
	ModelsFilename    string `name:"m" description:"The filename for the models file, excluding the extension" default:"models"`
	TS                bool   `name:"ts" description:"Generate Typescript bindings"`
	TSPrefix          string `description:"The prefix for typescript names" default:""`
	TSSuffix          string `description:"The postfix for typescript names" default:""`
	UseInterfaces     bool   `name:"i" description:"Use interfaces instead of classes"`
	UseBundledRuntime bool   `name:"b" description:"Use the bundled runtime instead of importing the npm package"`
	ProjectDirectory  string `name:"p" description:"The project directory" default:"."`
	UseIDs            bool   `name:"ids" description:"Use IDs instead of names in the binding calls"`
	OutputDirectory   string `name:"d" description:"The output directory" default:"frontend/bindings"`
}
