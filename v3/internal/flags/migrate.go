package flags

type Migrate struct {
	ProjectDir string `name:"d" description:"Path to the Wails v2 project" default:"."`
	Quiet      bool   `name:"q" description:"Suppress output to console"`
}
