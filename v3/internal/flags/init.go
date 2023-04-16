package flags

type Init struct {
	Common

	PackageName  string `name:"p" description:"Package name" default:"main"`
	TemplateName string `name:"t" description:"Name of built-in template to use, path to template or template url" default:"vanilla"`
	ProjectName  string `name:"n" description:"Name of project" default:""`
	ProjectDir   string `name:"d" description:"Project directory" default:"."`
	Quiet        bool   `name:"q" description:"Suppress output to console"`
	List         bool   `name:"l" description:"List templates"`
}
