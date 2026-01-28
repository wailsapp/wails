package flags

type Init struct {
	Common

	PackageName        string `name:"p" description:"Package name" default:"main"`
	TemplateName       string `name:"t" description:"Name of built-in template to use, path to template or template url" default:"vanilla"`
	ProjectName        string `name:"n" description:"Name of project" default:""`
	ProjectDir         string `name:"d" description:"Project directory" default:"."`
	Quiet              bool   `name:"q" description:"Suppress output to console"`
	List               bool   `name:"l" description:"List templates"`
	SkipGoModTidy      bool   `name:"skipgomodtidy" description:"Skip running go mod tidy"`
	Git                string `name:"git" description:"Git repository URL to initialize (e.g. github.com/username/project)"`
	ModulePath         string `name:"mod" description:"The Go module path for the project. Will be computed from the Git URL if unspecified."`
	ProductName        string `description:"The name of the product" default:"My Product"`
	ProductDescription string `description:"The description of the product" default:"My Product Description"`
	ProductVersion     string `description:"The version of the product" default:"0.1.0"`
	ProductCompany     string `description:"The company of the product" default:"My Company"`
	ProductCopyright   string `description:"The copyright notice" default:"\u00a9 now, My Company"`
	ProductComments    string `description:"Comments to add to the generated files" default:"This is a comment"`
	ProductIdentifier  string `description:"The product identifier, e.g com.mycompany.myproduct"`
	SkipWarning        bool   `name:"s" description:"Skips the warning message when using remote templates"`
	UseInterfaces      bool   `description:"Generate TypeScript interfaces instead of classes for bindings" default:"true"`

	TemplateFromDefaults      bool
	UseInterfacesFromDefaults bool
}
