package flags

type ServiceInit struct {
	Name        string `name:"n" description:"Name of service" default:"example_service"`
	Description string `name:"d" description:"Description of service" default:"Example service"`
	PackageName string `name:"p" description:"Package name for service" default:""`
	OutputDir   string `name:"o" description:"Output directory" default:"."`
	Quiet       bool   `name:"q" description:"Suppress output to console"`
	Author      string `name:"a" description:"Author of service" default:""`
	Version     string `name:"v" description:"Version of service" default:""`
	Website     string `name:"w" description:"Website of service" default:""`
	Repository  string `name:"r" description:"Repository of service" default:""`
	License     string `name:"l" description:"License of service" default:""`
}
