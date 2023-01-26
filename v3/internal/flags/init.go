package flags

type Init struct {
	Common

	TemplateName string `name:"t" description:"Name of built-in template to use, path to template or template url" default:"vanilla"`
	ProjectName  string `name:"n" description:"Name of project" default:"."`
	//CIMode       bool   `name:"ci" description:"CI Mode"`
	ProjectDir string `name:"d" description:"Project directory" default:"."`
	Quiet      bool   `name:"q" description:"Suppress output to console"`
	//InitGit    bool   `name:"g" description:"Initialise git repository"`
	//IDE        string `name:"ide" description:"Generate IDE project files"`
	List bool `name:"l" description:"List templates"`
}
