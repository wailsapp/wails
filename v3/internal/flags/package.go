package flags

// ToolPackage represents the options for the package command
type ToolPackage struct {
	Common

	Format            string `name:"format" description:"Package format to generate (deb, rpm, archlinux, dmg, pkg)" default:"deb"`
	ExecutableName    string `name:"name" description:"Name of the executable to package" default:"myapp"`
	ConfigPath        string `name:"config" description:"Path to the package configuration file" default:""`
	Out               string `name:"out" description:"Path to the output dir" default:"."`
	BackgroundImage   string `name:"background" description:"Path to an optional background image for the DMG" default:""`
	CreateDMG         bool   `name:"create-dmg" description:"Create a DMG file (macOS only)" default:"false"`
	CreatePKG         bool   `name:"create-pkg" description:"Create a PKG installer (macOS only)" default:"false"`
	GenerateTemplate  bool   `name:"generate-template" description:"Generate a sample configuration template" default:"false"`
	SkipNotarization  bool   `name:"skip-notarization" description:"Skip notarization step for PKG" default:"false"`
	ValidateOnly      bool   `name:"validate-only" description:"Only validate configuration without building" default:"false"`
}
