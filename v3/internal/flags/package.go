package flags

// ToolPackage represents the options for the package command
type ToolPackage struct {
	Common

	Format     string `name:"format" description:"Package format to generate (deb, rpm, archlinux)" default:"deb"`
	ConfigPath string `name:"config" description:"Path to the package configuration file" default:""`
}
