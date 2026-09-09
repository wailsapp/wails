package flags

// ToolPackage represents the options for the package command
type ToolPackage struct {
	Common

	Format          string `name:"format" description:"Package format to generate (deb, rpm, archlinux, dmg)" default:"deb"`
	ExecutableName  string `name:"name" description:"Name of the executable to package" default:"myapp"`
	ConfigPath      string `name:"config" description:"Path to the package configuration file" default:""`
	Out             string `name:"out" description:"Path to the output dir" default:"."`
	BackgroundImage string `name:"background" description:"Path to an optional background image for the DMG" default:""`
	DmgVolumeIcon   string `name:"volume-icon" description:"Path to the icon shown for the mounted DMG volume" default:""`
	DmgFileIcon     string `name:"file-icon" description:"Path to the icon shown for the DMG file in Finder" default:""`
	DmgFiles        string `name:"files" description:"Additional DMG files as name=path pairs separated by commas" default:""`
	DmgWindowWidth  int    `name:"window-width" description:"DMG Finder window width in pixels" default:"540"`
	DmgWindowHeight int    `name:"window-height" description:"DMG Finder window height in pixels" default:"380"`
	CreateDMG       bool   `name:"create-dmg" description:"Create a DMG file (macOS only)" default:"false"`
}
