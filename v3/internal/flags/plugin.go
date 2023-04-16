package flags

type PluginInit struct {
	Name        string `name:"n" description:"Name of plugin" default:"example_plugin"`
	Description string `name:"d" description:"Description of plugin" default:"Example plugin"`
	PackageName string `name:"p" description:"Package name for plugin" default:""`
	OutputDir   string `name:"o" description:"Output directory" default:"."`
	Quiet       bool   `name:"q" description:"Suppress output to console"`
}
