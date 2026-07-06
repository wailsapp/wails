package flags

type Migrate struct {
	Common

	V2Dir         string `name:"d" description:"Path to the Wails v2 project to migrate" default:"."`
	OutputDir     string `name:"o" description:"Directory to write the migrated Wails v3 project to"`
	Force         bool   `name:"f" description:"Write into a non-empty output directory"`
	Quiet         bool   `name:"q" description:"Suppress output to console"`
	SkipGoModTidy bool   `name:"skipgomodtidy" description:"Skip running go mod tidy on the migrated project"`
}
