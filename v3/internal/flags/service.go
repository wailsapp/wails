package flags

type ServiceInit struct {
	Name        string `name:"n" description:"Name of plugin" default:"example_plugin"`
	Description string `name:"d" description:"Description of plugin" default:"Example plugin"`
	PackageName string `name:"p" description:"Package name for plugin" default:""`
	OutputDir   string `name:"o" description:"Output directory" default:"."`
	Quiet       bool   `name:"q" description:"Suppress output to console"`
	Author      string `name:"a" description:"Author of plugin" default:""`
	Version     string `name:"v" description:"Version of plugin" default:""`
	Website     string `name:"w" description:"Website of plugin" default:""`
	Repository  string `name:"r" description:"Repository of plugin" default:""`
	License     string `name:"l" description:"License of plugin" default:""`
}
